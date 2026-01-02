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

func setupDeteksiHandler() (*DeteksiHandler, repository.DeteksiRepositoryInterface, *services.SSEService) {
	mockRepo := testhelpers.NewMockDeteksiRepository()
	sseService := services.NewSSEService()
	deteksiService := services.NewDeteksiService(mockRepo)
	handler := NewDeteksiHandler(deteksiService, sseService)
	return handler, mockRepo, sseService
}

func TestDeteksiHandler_handleGetAll(t *testing.T) {
	handler, _, _ := setupDeteksiHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/deteksi", nil)
	w := httptest.NewRecorder()

	handler.handleGetAll(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestDeteksiHandler_handleGetByID(t *testing.T) {
	handler, mockRepo, _ := setupDeteksiHandler()

	mockRepo.Create(dto.CreateDeteksiRequest{
		NilaiDeteksi:    4.2,
		NilaiSubdomain1:   4.0,
		NilaiSubdomain2:   4.5,
		NilaiSubdomain3:   4.1,
	}, "test-id")

	req := httptest.NewRequest(http.MethodGet, "/api/deteksi/test-id", nil)
	w := httptest.NewRecorder()

	handler.handleGetByID(w, req, "test-id")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestDeteksiHandler_handleGetByID_NotFound(t *testing.T) {
	handler, _, _ := setupDeteksiHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/deteksi/nonexistent", nil)
	w := httptest.NewRecorder()

	handler.handleGetByID(w, req, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestDeteksiHandler_handleCreate(t *testing.T) {
	handler, _, _ := setupDeteksiHandler()

	reqBody := dto.CreateDeteksiRequest{
		NilaiDeteksi:    4.2,
		NilaiSubdomain1: 4.0,
		NilaiSubdomain2: 4.5,
		NilaiSubdomain3: 4.1,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/deteksi", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleCreate(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestDeteksiHandler_handleCreate_InvalidBody(t *testing.T) {
	handler, _, _ := setupDeteksiHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/deteksi", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.handleCreate(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestDeteksiHandler_handleUpdate(t *testing.T) {
	handler, mockRepo, _ := setupDeteksiHandler()

	mockRepo.Create(dto.CreateDeteksiRequest{
		NilaiDeteksi:    4.0,
		NilaiSubdomain1: 4.0,
		NilaiSubdomain2: 4.0,
		NilaiSubdomain3: 4.0,
	}, "test-id")

	updateReq := dto.UpdateDeteksiRequest{
		NilaiDeteksi: floatPtr(4.5),
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/api/deteksi/test-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleUpdate(w, req, "test-id")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestDeteksiHandler_handleDelete(t *testing.T) {
	handler, mockRepo, _ := setupDeteksiHandler()

	mockRepo.Create(dto.CreateDeteksiRequest{
		NilaiDeteksi:    4.2,
		NilaiSubdomain1: 4.0,
		NilaiSubdomain2: 4.5,
		NilaiSubdomain3: 4.1,
	}, "test-id")

	req := httptest.NewRequest(http.MethodDelete, "/api/deteksi/test-id", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleDelete(w, req, "test-id")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestDeteksiHandler_ServeHTTP(t *testing.T) {
	handler, _, _ := setupDeteksiHandler()

	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"GET all", http.MethodGet, "/api/deteksi", http.StatusOK},
		{"GET by ID", http.MethodGet, "/api/deteksi/test-id", http.StatusNotFound},
		{"POST create with ID", http.MethodPost, "/api/deteksi/test-id", http.StatusBadRequest},
		{"PUT update without ID", http.MethodPut, "/api/deteksi", http.StatusBadRequest},
		{"DELETE without ID", http.MethodDelete, "/api/deteksi", http.StatusBadRequest},
		{"Method not allowed", http.MethodPatch, "/api/deteksi", http.StatusMethodNotAllowed},
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

