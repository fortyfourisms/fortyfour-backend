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

func setupIdentifikasiHandler() (*IdentifikasiHandler, repository.IdentifikasiRepositoryInterface, *services.SSEService) {
	mockRepo := testhelpers.NewMockIdentifikasiRepository()
	sseService := services.NewSSEService()
	identifikasiService := services.NewIdentifikasiService(mockRepo)
	handler := NewIdentifikasiHandler(identifikasiService, sseService)
	return handler, mockRepo, sseService
}

func TestIdentifikasiHandler_handleGetAll(t *testing.T) {
	handler, _, _ := setupIdentifikasiHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/identifikasi", nil)
	w := httptest.NewRecorder()

	handler.handleGetAll(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestIdentifikasiHandler_handleGetByID(t *testing.T) {
	handler, mockRepo, _ := setupIdentifikasiHandler()

	// Create test identifikasi
	mockRepo.Create(dto.CreateIdentifikasiRequest{
		NilaiIdentifikasi: 4.2,
		NilaiSubdomain1:   4.0,
		NilaiSubdomain2:   4.5,
		NilaiSubdomain3:   4.1,
		NilaiSubdomain4:   3.9,
		NilaiSubdomain5:   4.0,
	}, "test-id")

	req := httptest.NewRequest(http.MethodGet, "/api/identifikasi/test-id", nil)
	w := httptest.NewRecorder()

	handler.handleGetByID(w, req, "test-id")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestIdentifikasiHandler_handleGetByID_NotFound(t *testing.T) {
	handler, _, _ := setupIdentifikasiHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/identifikasi/nonexistent", nil)
	w := httptest.NewRecorder()

	handler.handleGetByID(w, req, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestIdentifikasiHandler_handleCreate(t *testing.T) {
	handler, _, _ := setupIdentifikasiHandler()

	reqBody := dto.CreateIdentifikasiRequest{
		NilaiIdentifikasi: 4.2,
		NilaiSubdomain1:   4.0,
		NilaiSubdomain2:   4.5,
		NilaiSubdomain3:   4.1,
		NilaiSubdomain4:   3.9,
		NilaiSubdomain5:   4.0,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/identifikasi", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleCreate(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestIdentifikasiHandler_handleCreate_InvalidBody(t *testing.T) {
	handler, _, _ := setupIdentifikasiHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/identifikasi", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.handleCreate(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestIdentifikasiHandler_handleUpdate(t *testing.T) {
	handler, mockRepo, _ := setupIdentifikasiHandler()

	// Create test identifikasi
	mockRepo.Create(dto.CreateIdentifikasiRequest{
		NilaiIdentifikasi: 4.0,
		NilaiSubdomain1:   4.0,
		NilaiSubdomain2:   4.0,
		NilaiSubdomain3:   4.0,
		NilaiSubdomain4:   4.0,
		NilaiSubdomain5:   4.0,
	}, "test-id")

	updateReq := dto.UpdateIdentifikasiRequest{
		NilaiIdentifikasi: floatPtr(4.5),
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/api/identifikasi/test-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleUpdate(w, req, "test-id")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestIdentifikasiHandler_handleUpdate_InvalidBody(t *testing.T) {
	handler, _, _ := setupIdentifikasiHandler()

	req := httptest.NewRequest(http.MethodPut, "/api/identifikasi/test-id", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.handleUpdate(w, req, "test-id")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestIdentifikasiHandler_handleDelete(t *testing.T) {
	handler, mockRepo, _ := setupIdentifikasiHandler()

	// Create test identifikasi
	mockRepo.Create(dto.CreateIdentifikasiRequest{
		NilaiIdentifikasi: 4.2,
		NilaiSubdomain1:   4.0,
		NilaiSubdomain2:   4.5,
		NilaiSubdomain3:   4.1,
		NilaiSubdomain4:   3.9,
		NilaiSubdomain5:   4.0,
	}, "test-id")

	req := httptest.NewRequest(http.MethodDelete, "/api/identifikasi/test-id", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleDelete(w, req, "test-id")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestIdentifikasiHandler_ServeHTTP(t *testing.T) {
	handler, _, _ := setupIdentifikasiHandler()

	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"GET all", http.MethodGet, "/api/identifikasi", http.StatusOK},
		{"GET by ID", http.MethodGet, "/api/identifikasi/test-id", http.StatusNotFound},
		{"POST create with ID", http.MethodPost, "/api/identifikasi/test-id", http.StatusBadRequest},
		{"PUT update without ID", http.MethodPut, "/api/identifikasi", http.StatusBadRequest},
		{"DELETE without ID", http.MethodDelete, "/api/identifikasi", http.StatusBadRequest},
		{"Method not allowed", http.MethodPatch, "/api/identifikasi", http.StatusMethodNotAllowed},
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

func floatPtr(f float64) *float64 {
	return &f
}
