package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/services"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ================================================================
// MOCK: ChatServiceInterface
// ================================================================

type mockChatService struct {
	GenerateSQLQueryFn  func(userQuestion string) (string, error)
	ExecuteQueryFn      func(sqlQuery string) ([]map[string]interface{}, error)
	FormatQueryResultFn func(userQuestion string, results []map[string]interface{}) (string, error)
	repo                repository.ChatRepository
	gemini              services.GeminiGenerator
}

func (m *mockChatService) GenerateSQLQuery(userQuestion string) (string, error) {
	if m.GenerateSQLQueryFn != nil {
		return m.GenerateSQLQueryFn(userQuestion)
	}
	return "SELECT 1", nil
}

func (m *mockChatService) ExecuteQuery(sqlQuery string) ([]map[string]interface{}, error) {
	if m.ExecuteQueryFn != nil {
		return m.ExecuteQueryFn(sqlQuery)
	}
	return []map[string]interface{}{{"col": "val"}}, nil
}

func (m *mockChatService) FormatQueryResults(userQuestion string, results []map[string]interface{}) (string, error) {
	if m.FormatQueryResultFn != nil {
		return m.FormatQueryResultFn(userQuestion, results)
	}
	return "hasil formatted", nil
}

func (m *mockChatService) Repo() repository.ChatRepository {
	return m.repo
}

func (m *mockChatService) GetGemini() services.GeminiGenerator {
	return m.gemini
}

// compile-time check
var _ services.ChatServiceInterface = (*mockChatService)(nil)

// ================================================================
// MOCK: GeminiGenerator
// ================================================================

type mockGemini struct {
	GenerateFn func(prompt string) (string, error)
}

func (m *mockGemini) Generate(prompt string) (string, error) {
	if m.GenerateFn != nil {
		return m.GenerateFn(prompt)
	}
	return "jawaban dari gemini", nil
}

// ================================================================
// HELPER: ChatHandler dengan ChatServiceInterface
// ================================================================

// chatHandlerTestable adalah versi handler yang menerima interface
type chatHandlerTestable struct {
	service services.ChatServiceInterface
}

func (h *chatHandlerTestable) Stream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	var req dto.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendSSEError(w, flusher, "Invalid request body")
		return
	}

	// STEP 1: Generate SQL Query
	sqlQuery, err := h.service.GenerateSQLQuery(req.Message)
	if err != nil {
		safePrompt := fmt.Sprintf(`Pengguna mengajukan permintaan yang tidak diizinkan: "%s". Tolak dengan sopan.`, req.Message)

		answer, genErr := h.service.GetGemini().Generate(safePrompt)
		if genErr != nil || strings.TrimSpace(answer) == "" {
			answer = "Maaf, permintaan tersebut tidak dapat diproses demi keamanan dan perlindungan data."
		}

		words := strings.Fields(answer)
		fullAnswer := ""
		for i, word := range words {
			if i > 0 {
				fullAnswer += " "
			}
			fullAnswer += word
			sendSSEEvent(w, flusher, map[string]interface{}{
				"type":    "chunk",
				"content": word + " ",
				"done":    false,
			})
		}
		sendSSEEvent(w, flusher, map[string]interface{}{
			"type":    "done",
			"content": fullAnswer,
			"done":    true,
		})
		_ = h.service.Repo().Save(req.SessionID, req.Message, fullAnswer)
		return
	}

	// STEP 2: Execute SQL
	results, err := h.service.ExecuteQuery(sqlQuery)
	if err != nil {
		sendSSEError(w, flusher, "Gagal menjalankan query database")
		return
	}

	// STEP 3: Format hasil
	answer, err := h.service.FormatQueryResults(req.Message, results)
	if err != nil {
		sendSSEError(w, flusher, "Gagal memformat jawaban")
		return
	}

	// STEP 4: Stream
	words := strings.Fields(answer)
	fullAnswer := ""
	for i, word := range words {
		if i > 0 {
			fullAnswer += " "
		}
		fullAnswer += word
		sendSSEEvent(w, flusher, map[string]interface{}{
			"type":    "chunk",
			"content": word + " ",
			"done":    false,
		})
	}
	sendSSEEvent(w, flusher, map[string]interface{}{
		"type":    "done",
		"content": fullAnswer,
		"done":    true,
	})

	// STEP 5: Save history
	_ = h.service.Repo().Save(req.SessionID, req.Message, fullAnswer)
}

func (h *chatHandlerTestable) DeleteSession(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SessionID string `json:"session_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if repo, ok := h.service.Repo().(*repository.InMemoryChatRepo); ok {
		if err := repo.DeleteSession(req.SessionID); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": fmt.Sprintf("Session %s deleted", req.SessionID),
		})
		return
	}

	http.Error(w, "unsupported repository type", http.StatusInternalServerError)
}

// ================================================================
// HELPER: builder handler + in-memory repo
// ================================================================

func newTestChatHandler(svc services.ChatServiceInterface) *chatHandlerTestable {
	return &chatHandlerTestable{service: svc}
}

func newChatServiceWithInMemoryRepo(gemini services.GeminiGenerator) (*mockChatService, *repository.InMemoryChatRepo) {
	repo := repository.NewInMemoryChatRepo()
	svc := &mockChatService{
		repo:   repo,
		gemini: gemini,
	}
	return svc, repo
}

// ================================================================
// HELPER: parse SSE body menjadi slice of event map
// ================================================================

func parseSSEEvents(body string) []map[string]interface{} {
	var events []map[string]interface{}
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		raw := strings.TrimPrefix(line, "data: ")
		var ev map[string]interface{}
		if err := json.Unmarshal([]byte(raw), &ev); err == nil {
			events = append(events, ev)
		}
	}
	return events
}

// ================================================================
// TEST: Stream — alur sukses
// ================================================================

func TestChatHandler_Stream_Success(t *testing.T) {
	svc, _ := newChatServiceWithInMemoryRepo(&mockGemini{})
	svc.GenerateSQLQueryFn = func(_ string) (string, error) { return "SELECT 1", nil }
	svc.ExecuteQueryFn = func(_ string) ([]map[string]interface{}, error) {
		return []map[string]interface{}{{"nama": "Test"}}, nil
	}
	svc.FormatQueryResultFn = func(_ string, _ []map[string]interface{}) (string, error) {
		return "Ada satu data yaitu Test", nil
	}

	handler := newTestChatHandler(svc)
	body, _ := json.Marshal(dto.ChatRequest{SessionID: "s1", Message: "tampilkan data"})
	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Stream(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/event-stream", w.Header().Get("Content-Type"))

	events := parseSSEEvents(w.Body.String())
	require.NotEmpty(t, events)

	// Event terakhir harus bertipe "done" dan done=true
	lastEvent := events[len(events)-1]
	assert.Equal(t, "done", lastEvent["type"])
	assert.Equal(t, true, lastEvent["done"])

	// Semua event sebelum "done" harus bertipe "chunk"
	for _, ev := range events[:len(events)-1] {
		assert.Equal(t, "chunk", ev["type"])
		assert.Equal(t, false, ev["done"])
	}
}

// ================================================================
// TEST: Stream — SSE header selalu di-set
// ================================================================

func TestChatHandler_Stream_SetsSSEHeaders(t *testing.T) {
	svc, _ := newChatServiceWithInMemoryRepo(&mockGemini{})
	svc.GenerateSQLQueryFn = func(_ string) (string, error) { return "SELECT 1", nil }
	svc.ExecuteQueryFn = func(_ string) ([]map[string]interface{}, error) { return nil, nil }
	svc.FormatQueryResultFn = func(_ string, _ []map[string]interface{}) (string, error) { return "oke", nil }

	handler := newTestChatHandler(svc)
	body, _ := json.Marshal(dto.ChatRequest{Message: "test"})
	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Stream(w, req)

	assert.Equal(t, "text/event-stream", w.Header().Get("Content-Type"))
	assert.Equal(t, "no-cache", w.Header().Get("Cache-Control"))
	assert.Equal(t, "keep-alive", w.Header().Get("Connection"))
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

// ================================================================
// TEST: Stream — request body tidak valid (bukan JSON)
// ================================================================

func TestChatHandler_Stream_InvalidBody(t *testing.T) {
	svc, _ := newChatServiceWithInMemoryRepo(&mockGemini{})
	handler := newTestChatHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/chat", strings.NewReader("bukan json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Stream(w, req)

	// Harus mengirim SSE error event
	body := w.Body.String()
	assert.Contains(t, body, "data:")
	events := parseSSEEvents(body)
	require.NotEmpty(t, events)
	assert.Equal(t, "error", events[0]["type"])
	assert.Equal(t, true, events[0]["done"])
}

// ================================================================
// TEST: Stream — GenerateSQLQuery gagal → fallback ke Gemini
// ================================================================

func TestChatHandler_Stream_SQLGenerationError_FallbackToGemini(t *testing.T) {
	gemini := &mockGemini{
		GenerateFn: func(prompt string) (string, error) {
			return "Maaf permintaan ini tidak dapat diproses", nil
		},
	}
	svc, _ := newChatServiceWithInMemoryRepo(gemini)
	svc.GenerateSQLQueryFn = func(_ string) (string, error) {
		return "", errors.New("query berbahaya: DROP TABLE")
	}

	handler := newTestChatHandler(svc)
	body, _ := json.Marshal(dto.ChatRequest{SessionID: "s-reject", Message: "hapus semua data"})
	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Stream(w, req)

	events := parseSSEEvents(w.Body.String())
	require.NotEmpty(t, events)

	// Event terakhir harus "done"
	last := events[len(events)-1]
	assert.Equal(t, "done", last["type"])
	assert.Equal(t, true, last["done"])

	// Konten harus berisi jawaban dari gemini
	fullContent := last["content"].(string)
	assert.Contains(t, fullContent, "Maaf")
}

// ================================================================
// TEST: Stream — GenerateSQLQuery gagal + Gemini juga gagal
//       → fallback ke pesan default hardcoded
// ================================================================

func TestChatHandler_Stream_SQLGenerationError_GeminiFallbackToDefault(t *testing.T) {
	gemini := &mockGemini{
		GenerateFn: func(prompt string) (string, error) {
			return "", errors.New("gemini juga down")
		},
	}
	svc, _ := newChatServiceWithInMemoryRepo(gemini)
	svc.GenerateSQLQueryFn = func(_ string) (string, error) {
		return "", errors.New("query tidak diizinkan")
	}

	handler := newTestChatHandler(svc)
	body, _ := json.Marshal(dto.ChatRequest{Message: "query berbahaya"})
	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Stream(w, req)

	events := parseSSEEvents(w.Body.String())
	require.NotEmpty(t, events)

	last := events[len(events)-1]
	assert.Equal(t, "done", last["type"])
	fullContent := last["content"].(string)
	// Harus pakai pesan default karena Gemini gagal
	assert.Contains(t, fullContent, "keamanan")
}

// ================================================================
// TEST: Stream — ExecuteQuery gagal → SSE error
// ================================================================

func TestChatHandler_Stream_ExecuteQueryError(t *testing.T) {
	svc, _ := newChatServiceWithInMemoryRepo(&mockGemini{})
	svc.GenerateSQLQueryFn = func(_ string) (string, error) { return "SELECT 1", nil }
	svc.ExecuteQueryFn = func(_ string) ([]map[string]interface{}, error) {
		return nil, errors.New("koneksi DB terputus")
	}

	handler := newTestChatHandler(svc)
	body, _ := json.Marshal(dto.ChatRequest{Message: "data user"})
	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Stream(w, req)

	events := parseSSEEvents(w.Body.String())
	require.NotEmpty(t, events)
	assert.Equal(t, "error", events[0]["type"])
	assert.Contains(t, events[0]["content"].(string), "query database")
}

// ================================================================
// TEST: Stream — FormatQueryResults gagal → SSE error
// ================================================================

func TestChatHandler_Stream_FormatError(t *testing.T) {
	svc, _ := newChatServiceWithInMemoryRepo(&mockGemini{})
	svc.GenerateSQLQueryFn = func(_ string) (string, error) { return "SELECT 1", nil }
	svc.ExecuteQueryFn = func(_ string) ([]map[string]interface{}, error) {
		return []map[string]interface{}{{"x": "y"}}, nil
	}
	svc.FormatQueryResultFn = func(_ string, _ []map[string]interface{}) (string, error) {
		return "", errors.New("format gagal")
	}

	handler := newTestChatHandler(svc)
	body, _ := json.Marshal(dto.ChatRequest{Message: "pertanyaan apapun"})
	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Stream(w, req)

	events := parseSSEEvents(w.Body.String())
	require.NotEmpty(t, events)
	assert.Equal(t, "error", events[0]["type"])
	assert.Contains(t, events[0]["content"].(string), "jawaban")
}

// ================================================================
// TEST: Stream — histori disimpan ke repo setelah berhasil
// ================================================================

func TestChatHandler_Stream_SavesHistoryOnSuccess(t *testing.T) {
	svc, repo := newChatServiceWithInMemoryRepo(&mockGemini{})
	svc.GenerateSQLQueryFn = func(_ string) (string, error) { return "SELECT 1", nil }
	svc.ExecuteQueryFn = func(_ string) ([]map[string]interface{}, error) {
		return []map[string]interface{}{{"col": "val"}}, nil
	}
	svc.FormatQueryResultFn = func(_ string, _ []map[string]interface{}) (string, error) {
		return "jawaban tersimpan", nil
	}

	handler := newTestChatHandler(svc)
	body, _ := json.Marshal(dto.ChatRequest{SessionID: "session-save", Message: "pertanyaan saya"})
	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Stream(w, req)

	// Verifikasi data tersimpan di repo
	history, err := repo.GetHistory("session-save")
	require.NoError(t, err)
	require.Len(t, history, 1)
	assert.Equal(t, "pertanyaan saya", history[0].User)
	assert.Equal(t, "jawaban tersimpan", history[0].Bot)
}

// ================================================================
// TEST: Stream — histori disimpan saat fallback rejection
// ================================================================

func TestChatHandler_Stream_SavesHistoryOnRejection(t *testing.T) {
	gemini := &mockGemini{
		GenerateFn: func(_ string) (string, error) {
			return "Permintaan ditolak karena keamanan", nil
		},
	}
	svc, repo := newChatServiceWithInMemoryRepo(gemini)
	svc.GenerateSQLQueryFn = func(_ string) (string, error) {
		return "", errors.New("query berbahaya")
	}

	handler := newTestChatHandler(svc)
	body, _ := json.Marshal(dto.ChatRequest{SessionID: "session-reject", Message: "hapus data"})
	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Stream(w, req)

	history, err := repo.GetHistory("session-reject")
	require.NoError(t, err)
	require.Len(t, history, 1)
	assert.Equal(t, "hapus data", history[0].User)
	assert.NotEmpty(t, history[0].Bot)
}

// ================================================================
// TEST: Stream — query menghasilkan 0 baris → tetap streaming
// ================================================================

func TestChatHandler_Stream_EmptyQueryResult(t *testing.T) {
	svc, _ := newChatServiceWithInMemoryRepo(&mockGemini{})
	svc.GenerateSQLQueryFn = func(_ string) (string, error) { return "SELECT 1", nil }
	svc.ExecuteQueryFn = func(_ string) ([]map[string]interface{}, error) {
		return []map[string]interface{}{}, nil // 0 rows
	}
	svc.FormatQueryResultFn = func(_ string, _ []map[string]interface{}) (string, error) {
		return "Tidak ada data ditemukan", nil
	}

	handler := newTestChatHandler(svc)
	body, _ := json.Marshal(dto.ChatRequest{Message: "cari data tidak ada"})
	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Stream(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	events := parseSSEEvents(w.Body.String())
	require.NotEmpty(t, events)
	last := events[len(events)-1]
	assert.Equal(t, "done", last["type"])
}

// ================================================================
// TEST: DeleteSession — sukses
// ================================================================

func TestChatHandler_DeleteSession_Success(t *testing.T) {
	svc, repo := newChatServiceWithInMemoryRepo(&mockGemini{})

	// Isi repo dengan satu session
	require.NoError(t, repo.Save("ses-del", "halo", "oke"))

	handler := newTestChatHandler(svc)
	body, _ := json.Marshal(map[string]string{"session_id": "ses-del"})
	req := httptest.NewRequest(http.MethodDelete, "/api/chat/delete-session", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.DeleteSession(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Contains(t, resp["message"], "ses-del")

	// Verifikasi session sudah terhapus
	history, err := repo.GetHistory("ses-del")
	require.NoError(t, err)
	assert.Empty(t, history)
}

// ================================================================
// TEST: DeleteSession — session tidak ditemukan → 404
// ================================================================

func TestChatHandler_DeleteSession_NotFound(t *testing.T) {
	svc, _ := newChatServiceWithInMemoryRepo(&mockGemini{})
	handler := newTestChatHandler(svc)

	body, _ := json.Marshal(map[string]string{"session_id": "tidak-ada"})
	req := httptest.NewRequest(http.MethodDelete, "/api/chat/delete-session", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.DeleteSession(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ================================================================
// TEST: DeleteSession — body tidak valid → 400
// ================================================================

func TestChatHandler_DeleteSession_InvalidBody(t *testing.T) {
	svc, _ := newChatServiceWithInMemoryRepo(&mockGemini{})
	handler := newTestChatHandler(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/chat/delete-session", strings.NewReader("bukan json"))
	w := httptest.NewRecorder()

	handler.DeleteSession(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ================================================================
// TEST: DeleteSession — repo bukan InMemoryChatRepo → 500
// ================================================================

func TestChatHandler_DeleteSession_UnsupportedRepo(t *testing.T) {
	// Gunakan mockChatRepo biasa (bukan *repository.InMemoryChatRepo)
	svc := &mockChatService{
		repo:   &mockChatRepoForDelete{},
		gemini: &mockGemini{},
	}
	handler := newTestChatHandler(svc)

	body, _ := json.Marshal(map[string]string{"session_id": "s1"})
	req := httptest.NewRequest(http.MethodDelete, "/api/chat/delete-session", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.DeleteSession(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ================================================================
// HELPER mock repo kecil untuk test DeleteSession unsupported
// ================================================================

type mockChatRepoForDelete struct{}

func (m *mockChatRepoForDelete) GetHistory(sessionID string) ([]dto.ChatHistory, error) {
	return nil, nil
}
func (m *mockChatRepoForDelete) Save(sessionID, userMsg, botMsg string) error { return nil }

var _ repository.ChatRepository = (*mockChatRepoForDelete)(nil)

// ================================================================
// TEST: sendSSEEvent — format output benar
// ================================================================

func TestSendSSEEvent_Format(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]interface{}{
		"type":    "chunk",
		"content": "halo",
		"done":    false,
	}
	sendSSEEvent(w, w, data)

	body := w.Body.String()
	assert.True(t, strings.HasPrefix(body, "data: "), "harus diawali 'data: '")
	assert.True(t, strings.HasSuffix(body, "\n\n"), "harus diakhiri double newline")

	// JSON di dalamnya harus bisa di-parse kembali
	raw := strings.TrimPrefix(strings.TrimSpace(body), "data: ")
	var out map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(raw), &out))
	assert.Equal(t, "chunk", out["type"])
	assert.Equal(t, "halo", out["content"])
}

// ================================================================
// TEST: sendSSEError — format error event benar
// ================================================================

func TestSendSSEError_Format(t *testing.T) {
	w := httptest.NewRecorder()
	sendSSEError(w, w, "terjadi kesalahan sistem")

	events := parseSSEEvents(w.Body.String())
	require.Len(t, events, 1)
	assert.Equal(t, "error", events[0]["type"])
	assert.Equal(t, "terjadi kesalahan sistem", events[0]["content"])
	assert.Equal(t, true, events[0]["done"])
}
