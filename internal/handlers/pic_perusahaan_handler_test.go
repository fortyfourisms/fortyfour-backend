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

func setupPICHandler() (*PICHandler, repository.PICRepositoryInterface, *services.SSEService) {
	mockRepo := testhelpers.NewMockPICRepository()
	sseService := services.NewSSEService()
	picService := services.NewPICService(mockRepo)
	handler := NewPICHandler(picService, sseService)
	return handler, mockRepo, sseService
}

func TestPICHandler_handleGetAll(t *testing.T) {
	handler, _, _ := setupPICHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/pic", nil)
	w := httptest.NewRecorder()

	handler.handleGetAll(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestPICHandler_handleGetByID(t *testing.T) {
	handler, mockRepo, _ := setupPICHandler()

	mockRepo.Create(dto.CreatePICRequest{
		Nama:    stringPtr("Test PIC"),
		Telepon: stringPtr("081234567890"),
	}, "test-id")

	req := httptest.NewRequest(http.MethodGet, "/api/pic/test-id", nil)
	w := httptest.NewRecorder()

	handler.handleGetByID(w, req, "test-id")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestPICHandler_handleGetByID_NotFound(t *testing.T) {
	handler, _, _ := setupPICHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/pic/nonexistent", nil)
	w := httptest.NewRecorder()

	handler.handleGetByID(w, req, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestPICHandler_handleCreate(t *testing.T) {
	handler, _, _ := setupPICHandler()

	reqBody := dto.CreatePICRequest{
		Nama:    stringPtr("New PIC"),
		Telepon: stringPtr("081234567890"),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/pic", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleCreate(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestPICHandler_handleCreate_InvalidBody(t *testing.T) {
	handler, _, _ := setupPICHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/pic", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.handleCreate(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestPICHandler_handleUpdate(t *testing.T) {
	handler, mockRepo, _ := setupPICHandler()

	mockRepo.Create(dto.CreatePICRequest{
		Nama:    stringPtr("Old Name"),
		Telepon: stringPtr("081234567890"),
	}, "test-id")

	updateReq := dto.UpdatePICRequest{
		Nama: stringPtr("New Name"),
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/api/pic/test-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleUpdate(w, req, "test-id")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestPICHandler_handleDelete(t *testing.T) {
	handler, mockRepo, _ := setupPICHandler()

	mockRepo.Create(dto.CreatePICRequest{
		Nama:    stringPtr("Test PIC"),
		Telepon: stringPtr("081234567890"),
	}, "test-id")

	req := httptest.NewRequest(http.MethodDelete, "/api/pic/test-id", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleDelete(w, req, "test-id")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestPICHandler_ServeHTTP(t *testing.T) {
	handler, _, _ := setupPICHandler()

	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"GET all", http.MethodGet, "/api/pic", http.StatusOK},
		{"GET by ID", http.MethodGet, "/api/pic/test-id", http.StatusNotFound},
		{"POST create with ID", http.MethodPost, "/api/pic/test-id", http.StatusBadRequest},
		{"PUT update without ID", http.MethodPut, "/api/pic", http.StatusBadRequest},
		{"DELETE without ID", http.MethodDelete, "/api/pic", http.StatusBadRequest},
		{"Method not allowed", http.MethodPatch, "/api/pic", http.StatusMethodNotAllowed},
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
