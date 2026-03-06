package handlers

import (
	"encoding/json"
	"fmt"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/services"
	"log"
	"net/http"
	"strings"
)

type ChatHandler struct {
	service *services.ChatService
}

func NewChatHandler(s *services.ChatService) *ChatHandler {
	return &ChatHandler{service: s}
}

// Stream godoc
// @Summary      Chat dengan AI (SSE streaming)
// @Description  Mengirim pesan ke AI dan menerima respons secara streaming via Server-Sent Events
// @Tags         Chat
// @Security     BearerAuth
// @Accept       json
// @Produce      text/event-stream
// @Param        request  body  dto.ChatRequest  true  "Chat request"
// @Success      200  {string}  string  "SSE stream"
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/chat [post]
func (h *ChatHandler) Stream(w http.ResponseWriter, r *http.Request) {
	// SSE headers
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

	log.Printf("User Question: %s", req.Message)

	// STEP 1: Generate SQL Query
	sqlQuery, err := h.service.GenerateSQLQuery(req.Message)
	if err != nil {
		log.Printf("SQL Generation Error: %v", err)

		// PROTECT AI RESPONSE (NOT ALLOWED OPERATION)
		safePrompt := fmt.Sprintf(`
		Pengguna mengajukan permintaan yang berpotensi membahayakan keamanan atau integritas data:

		"%s"

		Tolak permintaan tersebut dengan bahasa Indonesia yang sopan, singkat, dan profesional.
		JANGAN menyebutkan aturan teknis, jenis query, SQL, database, atau mekanisme internal apa pun.
		Cukup jelaskan bahwa permintaan tidak dapat diproses demi keamanan dan perlindungan data,
		lalu arahkan pengguna untuk mengajukan pertanyaan yang bersifat informatif.
		`, req.Message)

		answer, genErr := h.service.GetGemini().Generate(safePrompt)
		if genErr != nil || strings.TrimSpace(answer) == "" {
			answer = "Maaf, permintaan tersebut tidak dapat diproses demi keamanan dan perlindungan data. Silakan ajukan pertanyaan lain yang bersifat informatif."
		}

		// Stream jawaban AI
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
		log.Printf("Saved rejection response for session: %s", req.SessionID)

		return
	}

	log.Printf("Generated SQL: %s", sqlQuery)

	// STEP 2: Execute SQL Query
	results, err := h.service.ExecuteQuery(sqlQuery)
	if err != nil {
		log.Printf("Query Execution Error: %v", err)
		sendSSEError(w, flusher, "Gagal menjalankan query database")
		return
	}

	log.Printf("Query Results: %d rows", len(results))

	// STEP 3: Format hasil query
	answer, err := h.service.FormatQueryResults(req.Message, results)
	if err != nil {
		log.Printf("Formatting Error: %v", err)
		sendSSEError(w, flusher, "Gagal memformat jawaban")
		return
	}

	// STEP 4: Stream response word by word
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
	log.Printf("Saved to history for session: %s", req.SessionID)
}

// DeleteSession godoc
// @Summary      Hapus chat session
// @Description  Menghapus riwayat percakapan berdasarkan session ID
// @Tags         Chat
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      object{session_id=string}  true  "Session ID"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      404      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /api/chat/delete-session [delete]
func (h *ChatHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
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

// Helper functions
func sendSSEEvent(w http.ResponseWriter, flusher http.Flusher, data interface{}) {
	jsonData, _ := json.Marshal(data)
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	flusher.Flush()
}

func sendSSEError(w http.ResponseWriter, flusher http.Flusher, message string) {
	event := map[string]interface{}{
		"type":    "error",
		"content": message,
		"done":    true,
	}
	sendSSEEvent(w, flusher, event)
}
