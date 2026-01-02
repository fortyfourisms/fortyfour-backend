package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/testhelpers"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

func setupRoleHandler() (*RoleHandler, repository.RoleRepository, *services.SSEService) {
	mockRepo := testhelpers.NewMockRoleRepository()
	sseService := services.NewSSEService()
	roleService := services.NewRoleService(mockRepo)
	handler := NewRoleHandler(roleService, sseService)
	return handler, mockRepo, sseService
}

func TestRoleHandler_handleGetAll(t *testing.T) {
	handler, mockRepo, _ := setupRoleHandler()

	// Create test roles
	role1 := &models.Role{
		ID:          uuid.New().String(),
		Name:        "admin",
		Description: "Administrator",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	role2 := &models.Role{
		ID:          uuid.New().String(),
		Name:        "user",
		Description: "Regular User",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockRepo.Create(context.Background(), role1)
	mockRepo.Create(context.Background(), role2)

	req := httptest.NewRequest(http.MethodGet, "/api/role", nil)
	w := httptest.NewRecorder()

	handler.handleGetAll(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response []*dto.RoleResponse
	json.NewDecoder(w.Body).Decode(&response)
	if len(response) != 2 {
		t.Errorf("expected 2 roles, got %d", len(response))
	}
}

func TestRoleHandler_handleGetByID(t *testing.T) {
	handler, mockRepo, _ := setupRoleHandler()

	// Create test role
	role := &models.Role{
		ID:          "test-role-id",
		Name:        "test-role",
		Description: "Test Description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockRepo.Create(context.Background(), role)

	req := httptest.NewRequest(http.MethodGet, "/api/role/test-role-id", nil)
	w := httptest.NewRecorder()

	handler.handleGetByID(w, req, "test-role-id")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response dto.RoleResponse
	json.NewDecoder(w.Body).Decode(&response)
	if response.ID != "test-role-id" {
		t.Errorf("expected ID 'test-role-id', got '%s'", response.ID)
	}
}

func TestRoleHandler_handleGetByID_NotFound(t *testing.T) {
	handler, _, _ := setupRoleHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/role/nonexistent", nil)
	w := httptest.NewRecorder()

	handler.handleGetByID(w, req, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestRoleHandler_handleCreate(t *testing.T) {
	handler, _, _ := setupRoleHandler()

	reqBody := dto.CreateRoleRequest{
		Name:        "new-role",
		Description: "New Role Description",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/role", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleCreate(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var response dto.RoleResponse
	json.NewDecoder(w.Body).Decode(&response)
	if response.Name != "new-role" {
		t.Errorf("expected name 'new-role', got '%s'", response.Name)
	}
}

func TestRoleHandler_handleCreate_InvalidBody(t *testing.T) {
	handler, _, _ := setupRoleHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/role", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.handleCreate(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestRoleHandler_handleCreate_WithID(t *testing.T) {
	handler, _, _ := setupRoleHandler()

	reqBody := dto.CreateRoleRequest{
		Name:        "new-role",
		Description: "New Role Description",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/role/test-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestRoleHandler_handleUpdate(t *testing.T) {
	handler, mockRepo, _ := setupRoleHandler()

	// Create test role
	role := &models.Role{
		ID:          "test-role-id",
		Name:        "old-name",
		Description: "Old Description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockRepo.Create(context.Background(), role)

	updateReq := dto.UpdateRoleRequest{
		Name:        "new-name",
		Description: "New Description",
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/api/role/test-role-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleUpdate(w, req, "test-role-id")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestRoleHandler_handleUpdate_InvalidBody(t *testing.T) {
	handler, _, _ := setupRoleHandler()

	req := httptest.NewRequest(http.MethodPut, "/api/role/test-id", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.handleUpdate(w, req, "test-id")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestRoleHandler_handleUpdate_WithoutID(t *testing.T) {
	handler, _, _ := setupRoleHandler()

	updateReq := dto.UpdateRoleRequest{
		Name: "new-name",
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/api/role", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestRoleHandler_handleDelete(t *testing.T) {
	handler, mockRepo, _ := setupRoleHandler()

	// Create test role
	role := &models.Role{
		ID:          "test-role-id",
		Name:        "test-role",
		Description: "Test Description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockRepo.Create(context.Background(), role)

	req := httptest.NewRequest(http.MethodDelete, "/api/role/test-role-id", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleDelete(w, req, "test-role-id")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestRoleHandler_handleDelete_WithoutID(t *testing.T) {
	handler, _, _ := setupRoleHandler()

	req := httptest.NewRequest(http.MethodDelete, "/api/role", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestRoleHandler_ServeHTTP(t *testing.T) {
	handler, _, _ := setupRoleHandler()

	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"GET all", http.MethodGet, "/api/role", http.StatusOK},
		{"GET by ID", http.MethodGet, "/api/role/test-id", http.StatusNotFound},
		{"POST create", http.MethodPost, "/api/role", http.StatusBadRequest},
		{"PUT update", http.MethodPut, "/api/role/test-id", http.StatusBadRequest},
		{"DELETE", http.MethodDelete, "/api/role/test-id", http.StatusBadRequest},
		{"Method not allowed", http.MethodPatch, "/api/role", http.StatusMethodNotAllowed},
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
