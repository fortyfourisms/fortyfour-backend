package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/testhelpers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupJabatanHandler() (*JabatanHandler, repository.JabatanRepositoryInterface, *services.SSEService) {
	mockRepo := testhelpers.NewMockJabatanRepository()
	sseService := services.NewSSEService()
	jabatanService := services.NewJabatanService(mockRepo, nil)
	handler := NewJabatanHandler(jabatanService, sseService)
	return handler, mockRepo, sseService
}

// seedJabatan adalah helper untuk membuat jabatan di mock repo
func seedJabatan(t *testing.T, repo repository.JabatanRepositoryInterface, id, nama string) {
	t.Helper()
	repo.Create(dto.CreateJabatanRequest{NamaJabatan: stringPtr(nama)}, id)
}

// =========================
// handleGetAll
// =========================

func TestJabatanHandler_handleGetAll(t *testing.T) {
	handler, _, _ := setupJabatanHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/jabatan", nil)
	w := httptest.NewRecorder()

	handler.handleGetAll(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestJabatanHandler_handleGetAll_Empty(t *testing.T) {
	handler, _, _ := setupJabatanHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/jabatan", nil)
	w := httptest.NewRecorder()
	handler.handleGetAll(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for empty list, got %d", w.Code)
	}
	var result []dto.JabatanResponse
	json.NewDecoder(w.Body).Decode(&result)
	if len(result) != 0 {
		t.Errorf("expected empty list, got %d items", len(result))
	}
}

func TestJabatanHandler_handleGetAll_ResponseBody(t *testing.T) {
	handler, mockRepo, _ := setupJabatanHandler()
	mockRepo.Create(dto.CreateJabatanRequest{NamaJabatan: stringPtr("Manajer")}, "id-1")
	mockRepo.Create(dto.CreateJabatanRequest{NamaJabatan: stringPtr("Staff")}, "id-2")

	req := httptest.NewRequest(http.MethodGet, "/api/jabatan", nil)
	w := httptest.NewRecorder()
	handler.handleGetAll(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var result []dto.JabatanResponse
	json.NewDecoder(w.Body).Decode(&result)
	if len(result) != 2 {
		t.Errorf("expected 2 jabatan, got %d", len(result))
	}
}

func TestJabatanHandler_handleGetAll_ContentTypeJSON(t *testing.T) {
	handler, _, _ := setupJabatanHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/jabatan", nil)
	w := httptest.NewRecorder()
	handler.handleGetAll(w, req)

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Content-Type: want 'application/json', got '%s'", ct)
	}
}

// =========================
// handleGetByID
// =========================

func TestJabatanHandler_handleGetByID(t *testing.T) {
	handler, mockRepo, _ := setupJabatanHandler()

	mockRepo.Create(dto.CreateJabatanRequest{NamaJabatan: stringPtr("Test Jabatan")}, "test-id")

	req := httptest.NewRequest(http.MethodGet, "/api/jabatan/test-id", nil)
	w := httptest.NewRecorder()

	handler.handleGetByID(w, req, "test-id")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response dto.JabatanResponse
	json.NewDecoder(w.Body).Decode(&response)
	if response.ID != "test-id" {
		t.Errorf("expected ID 'test-id', got '%s'", response.ID)
	}
}

func TestJabatanHandler_handleGetByID_NotFound(t *testing.T) {
	handler, _, _ := setupJabatanHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/jabatan/nonexistent", nil)
	w := httptest.NewRecorder()

	handler.handleGetByID(w, req, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestJabatanHandler_handleGetByID_ResponseBodyFields(t *testing.T) {
	handler, mockRepo, _ := setupJabatanHandler()

	seedJabatan(t, mockRepo, "jab-abc", "Supervisor")

	req := httptest.NewRequest(http.MethodGet, "/api/jabatan/jab-abc", nil)
	w := httptest.NewRecorder()
	handler.handleGetByID(w, req, "jab-abc")

	var response dto.JabatanResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.ID != "jab-abc" {
		t.Errorf("ID: want 'jab-abc', got '%s'", response.ID)
	}
	if response.NamaJabatan != "Supervisor" {
		t.Errorf("NamaJabatan: want 'Supervisor', got '%s'", response.NamaJabatan)
	}
}

func TestJabatanHandler_handleGetByID_NotFound_ErrorInBody(t *testing.T) {
	handler, _, _ := setupJabatanHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/jabatan/tidak-ada", nil)
	w := httptest.NewRecorder()
	handler.handleGetByID(w, req, "tidak-ada")

	// RespondError memakai key "error"
	var errResp map[string]string
	json.NewDecoder(w.Body).Decode(&errResp)
	if errResp["error"] == "" {
		t.Error("expected 'error' field in response body on 404")
	}
}

// =========================
// handleCreate
// =========================

func TestJabatanHandler_handleCreate(t *testing.T) {
	handler, _, _ := setupJabatanHandler()

	reqBody := dto.CreateJabatanRequest{
		NamaJabatan: stringPtr("New Jabatan"),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/jabatan", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleCreate(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestJabatanHandler_handleCreate_InvalidBody(t *testing.T) {
	handler, _, _ := setupJabatanHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/jabatan", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.handleCreate(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestJabatanHandler_handleCreate_ResponseBody(t *testing.T) {
	handler, _, _ := setupJabatanHandler()

	reqBody := dto.CreateJabatanRequest{NamaJabatan: stringPtr("Direktur")}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/jabatan", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	handler.handleCreate(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
	var result dto.JabatanResponse
	json.NewDecoder(w.Body).Decode(&result)
	if result.NamaJabatan != "Direktur" {
		t.Errorf("expected NamaJabatan 'Direktur', got '%s'", result.NamaJabatan)
	}
}

func TestJabatanHandler_handleCreate_WithoutUserContext(t *testing.T) {
	handler, _, _ := setupJabatanHandler()

	body, _ := json.Marshal(dto.CreateJabatanRequest{NamaJabatan: stringPtr("Guest Jabatan")})
	req := httptest.NewRequest(http.MethodPost, "/api/jabatan", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// Sengaja tidak set UserIDKey
	w := httptest.NewRecorder()

	handler.handleCreate(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 even without user context, got %d", w.Code)
	}
}

// =========================
// handleUpdate
// =========================

func TestJabatanHandler_handleUpdate(t *testing.T) {
	handler, mockRepo, _ := setupJabatanHandler()

	mockRepo.Create(dto.CreateJabatanRequest{NamaJabatan: stringPtr("Old Name")}, "test-id")

	updateReq := dto.UpdateJabatanRequest{
		NamaJabatan: stringPtr("New Name"),
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/api/jabatan/test-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleUpdate(w, req, "test-id")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestJabatanHandler_handleUpdate_InvalidBody(t *testing.T) {
	handler, _, _ := setupJabatanHandler()

	req := httptest.NewRequest(http.MethodPut, "/api/jabatan/test-id", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.handleUpdate(w, req, "test-id")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestJabatanHandler_handleUpdate_NotFound(t *testing.T) {
	handler, _, _ := setupJabatanHandler()

	updateReq := dto.UpdateJabatanRequest{NamaJabatan: stringPtr("Baru")}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/api/jabatan/non-existent", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	handler.handleUpdate(w, req, "non-existent")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for update on nonexistent jabatan, got %d", w.Code)
	}
}

func TestJabatanHandler_handleUpdate_ResponseBody(t *testing.T) {
	handler, mockRepo, _ := setupJabatanHandler()
	mockRepo.Create(dto.CreateJabatanRequest{NamaJabatan: stringPtr("Lama")}, "upd-id")

	updateReq := dto.UpdateJabatanRequest{NamaJabatan: stringPtr("Baru")}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/api/jabatan/upd-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	handler.handleUpdate(w, req, "upd-id")

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var result dto.JabatanResponse
	json.NewDecoder(w.Body).Decode(&result)
	if result.NamaJabatan != "Baru" {
		t.Errorf("expected NamaJabatan 'Baru', got '%s'", result.NamaJabatan)
	}
}

func TestJabatanHandler_handleUpdate_WithoutUserContext(t *testing.T) {
	handler, mockRepo, _ := setupJabatanHandler()

	seedJabatan(t, mockRepo, "jab-ctx", "Original")

	body, _ := json.Marshal(dto.UpdateJabatanRequest{NamaJabatan: stringPtr("Updated")})
	req := httptest.NewRequest(http.MethodPut, "/api/jabatan/jab-ctx", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// Sengaja tidak set UserIDKey
	w := httptest.NewRecorder()

	handler.handleUpdate(w, req, "jab-ctx")

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 even without user context, got %d", w.Code)
	}
}

// =========================
// handleDelete
// =========================

func TestJabatanHandler_handleDelete(t *testing.T) {
	handler, mockRepo, _ := setupJabatanHandler()

	mockRepo.Create(dto.CreateJabatanRequest{NamaJabatan: stringPtr("Test Jabatan")}, "test-id")

	req := httptest.NewRequest(http.MethodDelete, "/api/jabatan/test-id", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleDelete(w, req, "test-id")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestJabatanHandler_handleDelete_NotFound(t *testing.T) {
	handler, _, _ := setupJabatanHandler()

	req := httptest.NewRequest(http.MethodDelete, "/api/jabatan/tidak-ada", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	handler.handleDelete(w, req, "tidak-ada")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for delete on nonexistent jabatan, got %d", w.Code)
	}
}

func TestJabatanHandler_handleDelete_ResponseBodyHasMessage(t *testing.T) {
	handler, mockRepo, _ := setupJabatanHandler()

	seedJabatan(t, mockRepo, "jab-del", "To Delete")

	req := httptest.NewRequest(http.MethodDelete, "/api/jabatan/jab-del", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	handler.handleDelete(w, req, "jab-del")

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	if response["message"] == "" {
		t.Error("expected 'message' field in delete response body")
	}
}

func TestJabatanHandler_handleDelete_CannotGetDeletedJabatan(t *testing.T) {
	handler, mockRepo, _ := setupJabatanHandler()

	seedJabatan(t, mockRepo, "jab-gone", "Temporary")

	// Delete
	delReq := httptest.NewRequest(http.MethodDelete, "/api/jabatan/jab-gone", nil)
	ctx := context.WithValue(delReq.Context(), middleware.UserIDKey, "user-1")
	delReq = delReq.WithContext(ctx)
	delW := httptest.NewRecorder()
	handler.handleDelete(delW, delReq, "jab-gone")

	if delW.Code != http.StatusOK {
		t.Fatalf("delete should succeed, got %d", delW.Code)
	}

	// Coba get — harus 404
	getReq := httptest.NewRequest(http.MethodGet, "/api/jabatan/jab-gone", nil)
	getW := httptest.NewRecorder()
	handler.handleGetByID(getW, getReq, "jab-gone")

	if getW.Code != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", getW.Code)
	}
}

func TestJabatanHandler_handleDelete_WithoutUserContext(t *testing.T) {
	handler, mockRepo, _ := setupJabatanHandler()

	seedJabatan(t, mockRepo, "jab-noctx", "Test")

	req := httptest.NewRequest(http.MethodDelete, "/api/jabatan/jab-noctx", nil)
	// Sengaja tidak set UserIDKey
	w := httptest.NewRecorder()

	handler.handleDelete(w, req, "jab-noctx")

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 even without user context, got %d", w.Code)
	}
}

// =========================
// ServeHTTP routing
// =========================

func TestJabatanHandler_ServeHTTP(t *testing.T) {
	handler, _, _ := setupJabatanHandler()

	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"GET all", http.MethodGet, "/api/jabatan", http.StatusOK},
		{"GET by ID", http.MethodGet, "/api/jabatan/test-id", http.StatusNotFound},
		{"POST create with ID", http.MethodPost, "/api/jabatan/test-id", http.StatusBadRequest},
		{"PUT update without ID", http.MethodPut, "/api/jabatan", http.StatusBadRequest},
		{"DELETE without ID", http.MethodDelete, "/api/jabatan", http.StatusBadRequest},
		{"Method not allowed", http.MethodPatch, "/api/jabatan", http.StatusMethodNotAllowed},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}
		})
	}
}

func TestJabatanHandler_ServeHTTP_HeadNotAllowed(t *testing.T) {
	handler, _, _ := setupJabatanHandler()

	req := httptest.NewRequest(http.MethodHead, "/api/jabatan", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 for HEAD, got %d", w.Code)
	}
}

func TestJabatanHandler_ServeHTTP_DeleteWithID(t *testing.T) {
	handler, mockRepo, _ := setupJabatanHandler()

	seedJabatan(t, mockRepo, "del-via-http", "Via ServeHTTP")

	req := httptest.NewRequest(http.MethodDelete, "/api/jabatan/del-via-http", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for DELETE with ID via ServeHTTP, got %d", w.Code)
	}
}

func TestJabatanHandler_ServeHTTP_PutWithID(t *testing.T) {
	handler, mockRepo, _ := setupJabatanHandler()

	seedJabatan(t, mockRepo, "put-via-http", "Via ServeHTTP")

	body, _ := json.Marshal(dto.UpdateJabatanRequest{NamaJabatan: stringPtr("Updated")})
	req := httptest.NewRequest(http.MethodPut, "/api/jabatan/put-via-http", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for PUT with ID via ServeHTTP, got %d", w.Code)
	}
}