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

func setupGulihHandler() (*GulihHandler, repository.GulihRepositoryInterface, *services.SSEService) {
	mockRepo := testhelpers.NewMockGulihRepository()
	sseService := services.NewSSEService()
	gulihService := services.NewGulihService(mockRepo)
	handler := NewGulihHandler(gulihService, sseService)
	return handler, mockRepo, sseService
}

func TestGulihHandler_handleGetAll(t *testing.T) {
	handler, _, _ := setupGulihHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/gulih", nil)
	w := httptest.NewRecorder()

	handler.handleGetAll(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestGulihHandler_handleGetByID(t *testing.T) {
	handler, mockRepo, _ := setupGulihHandler()

	mockRepo.Create(dto.CreateGulihRequest{
		NilaiGulih:      4.2,
		NilaiSubdomain1: 4.0,
		NilaiSubdomain2: 4.5,
		NilaiSubdomain3: 4.1,
		NilaiSubdomain4: 3.9,
	}, "test-id")

	req := httptest.NewRequest(http.MethodGet, "/api/gulih/test-id", nil)
	w := httptest.NewRecorder()

	handler.handleGetByID(w, req, "test-id")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestGulihHandler_handleGetByID_NotFound(t *testing.T) {
	handler, _, _ := setupGulihHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/gulih/nonexistent", nil)
	w := httptest.NewRecorder()

	handler.handleGetByID(w, req, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestGulihHandler_handleCreate(t *testing.T) {
	handler, _, _ := setupGulihHandler()

	reqBody := dto.CreateGulihRequest{
		NilaiGulih:      4.2,
		NilaiSubdomain1: 4.0,
		NilaiSubdomain2: 4.5,
		NilaiSubdomain3: 4.1,
		NilaiSubdomain4: 3.9,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/gulih", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleCreate(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestGulihHandler_handleCreate_InvalidBody(t *testing.T) {
	handler, _, _ := setupGulihHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/gulih", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.handleCreate(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestGulihHandler_handleUpdate(t *testing.T) {
	handler, mockRepo, _ := setupGulihHandler()

	mockRepo.Create(dto.CreateGulihRequest{
		NilaiGulih:      4.0,
		NilaiSubdomain1: 4.0,
		NilaiSubdomain2: 4.0,
		NilaiSubdomain3: 4.0,
		NilaiSubdomain4: 4.0,
	}, "test-id")

	updateReq := dto.UpdateGulihRequest{
		NilaiGulih: floatPtr(4.5),
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/api/gulih/test-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleUpdate(w, req, "test-id")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestGulihHandler_handleDelete(t *testing.T) {
	handler, mockRepo, _ := setupGulihHandler()

	mockRepo.Create(dto.CreateGulihRequest{
		NilaiGulih:      4.2,
		NilaiSubdomain1: 4.0,
		NilaiSubdomain2: 4.5,
		NilaiSubdomain3: 4.1,
		NilaiSubdomain4: 3.9,
	}, "test-id")

	req := httptest.NewRequest(http.MethodDelete, "/api/gulih/test-id", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleDelete(w, req, "test-id")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestGulihHandler_ServeHTTP(t *testing.T) {
	handler, _, _ := setupGulihHandler()

	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"GET all", http.MethodGet, "/api/gulih", http.StatusOK},
		{"GET by ID", http.MethodGet, "/api/gulih/test-id", http.StatusNotFound},
		{"POST create with ID", http.MethodPost, "/api/gulih/test-id", http.StatusBadRequest},
		{"PUT update without ID", http.MethodPut, "/api/gulih", http.StatusBadRequest},
		{"DELETE without ID", http.MethodDelete, "/api/gulih", http.StatusBadRequest},
		{"Method not allowed", http.MethodPatch, "/api/gulih", http.StatusMethodNotAllowed},
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

