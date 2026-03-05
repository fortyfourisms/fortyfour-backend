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

func TestJabatanHandler_handleGetAll(t *testing.T) {
	handler, _, _ := setupJabatanHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/jabatan", nil)
	w := httptest.NewRecorder()

	handler.handleGetAll(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestJabatanHandler_handleGetByID(t *testing.T) {
	handler, mockRepo, _ := setupJabatanHandler()

	// Create test jabatan
	// jabatan := &dto.JabatanResponse{
	// 	ID:          "test-id",
	// 	NamaJabatan: "Test Jabatan",
	// }
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

func TestJabatanHandler_handleUpdate(t *testing.T) {
	handler, mockRepo, _ := setupJabatanHandler()

	// Create test jabatan
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

func TestJabatanHandler_handleDelete(t *testing.T) {
	handler, mockRepo, _ := setupJabatanHandler()

	// Create test jabatan
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

// ============================================================
// TAMBAHAN: error path & response body
// ============================================================

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

	// Handler mengembalikan 400 untuk semua error termasuk not found
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestJabatanHandler_handleDelete_NotFound(t *testing.T) {
	handler, _, _ := setupJabatanHandler()

	req := httptest.NewRequest(http.MethodDelete, "/api/jabatan/tidak-ada", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	handler.handleDelete(w, req, "tidak-ada")

	// Handler mengembalikan 400 untuk semua error termasuk not found
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
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
