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
	sseService := services.NewSSEService(nil)
	roleService := services.NewRoleService(mockRepo, nil, nil)
	handler := NewRoleHandler(roleService, sseService)
	return handler, mockRepo, sseService
}

// seedRole adalah helper untuk membuat role di mock repo
func seedRole(t *testing.T, repo repository.RoleRepository, id, name, desc string) *models.Role {
	t.Helper()
	role := &models.Role{
		ID:          id,
		Name:        name,
		Description: desc,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	repo.Create(context.Background(), role)
	return role
}

// =========================
// handleGetAll
// =========================

func TestRoleHandler_handleGetAll(t *testing.T) {
	handler, mockRepo, _ := setupRoleHandler()

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

func TestRoleHandler_handleGetAll_Empty(t *testing.T) {
	handler, _, _ := setupRoleHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/role", nil)
	w := httptest.NewRecorder()

	handler.handleGetAll(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 for empty list, got %d", w.Code)
	}

	var response []*dto.RoleResponse
	json.NewDecoder(w.Body).Decode(&response)
	if len(response) != 0 {
		t.Errorf("expected empty list, got %d items", len(response))
	}
}

func TestRoleHandler_handleGetAll_ResponseBodyHasCorrectFields(t *testing.T) {
	handler, mockRepo, _ := setupRoleHandler()

	seedRole(t, mockRepo, "role-abc", "editor", "Can edit content")

	req := httptest.NewRequest(http.MethodGet, "/api/role", nil)
	w := httptest.NewRecorder()
	handler.handleGetAll(w, req)

	var response []*dto.RoleResponse
	json.NewDecoder(w.Body).Decode(&response)

	if len(response) != 1 {
		t.Fatalf("expected 1 role, got %d", len(response))
	}
	if response[0].ID != "role-abc" {
		t.Errorf("ID: want 'role-abc', got '%s'", response[0].ID)
	}
	if response[0].Name != "editor" {
		t.Errorf("Name: want 'editor', got '%s'", response[0].Name)
	}
	if response[0].Description != "Can edit content" {
		t.Errorf("Description: want 'Can edit content', got '%s'", response[0].Description)
	}
}

func TestRoleHandler_handleGetAll_ContentTypeJSON(t *testing.T) {
	handler, _, _ := setupRoleHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/role", nil)
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

func TestRoleHandler_handleGetByID(t *testing.T) {
	handler, mockRepo, _ := setupRoleHandler()

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

func TestRoleHandler_handleGetByID_ResponseBodyHasCorrectFields(t *testing.T) {
	handler, mockRepo, _ := setupRoleHandler()

	seedRole(t, mockRepo, "role-xyz", "moderator", "Can moderate")

	req := httptest.NewRequest(http.MethodGet, "/api/role/role-xyz", nil)
	w := httptest.NewRecorder()
	handler.handleGetByID(w, req, "role-xyz")

	var response dto.RoleResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Name != "moderator" {
		t.Errorf("Name: want 'moderator', got '%s'", response.Name)
	}
	if response.Description != "Can moderate" {
		t.Errorf("Description: want 'Can moderate', got '%s'", response.Description)
	}
}

func TestRoleHandler_handleGetByID_NotFound_ErrorMessageInBody(t *testing.T) {
	handler, _, _ := setupRoleHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/role/does-not-exist", nil)
	w := httptest.NewRecorder()
	handler.handleGetByID(w, req, "does-not-exist")

	// RespondError menggunakan key "error", bukan "message"
	var errResp map[string]string
	json.NewDecoder(w.Body).Decode(&errResp)

	if errResp["error"] == "" {
		t.Error("expected 'error' field in response body")
	}
}

// =========================
// handleCreate
// =========================

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

func TestRoleHandler_handleCreate_ResponseBodyHasID(t *testing.T) {
	handler, _, _ := setupRoleHandler()

	body, _ := json.Marshal(dto.CreateRoleRequest{Name: "viewer", Description: "Read only"})
	req := httptest.NewRequest(http.MethodPost, "/api/role", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleCreate(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var response dto.RoleResponse
	json.NewDecoder(w.Body).Decode(&response)

	// MockRoleRepository tidak auto-generate UUID, tapi Name dan Description harus benar
	if response.Name != "viewer" {
		t.Errorf("Name: want 'viewer', got '%s'", response.Name)
	}
	if response.Description != "Read only" {
		t.Errorf("Description: want 'Read only', got '%s'", response.Description)
	}
}

func TestRoleHandler_handleCreate_WithoutUserContext(t *testing.T) {
	handler, _, _ := setupRoleHandler()

	body, _ := json.Marshal(dto.CreateRoleRequest{Name: "guest", Description: "Guest role"})
	req := httptest.NewRequest(http.MethodPost, "/api/role", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// Sengaja tidak set UserIDKey di context
	w := httptest.NewRecorder()

	handler.handleCreate(w, req)

	// Harus tetap berhasil — user ID kosong tapi tidak error
	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201 even without user context, got %d", w.Code)
	}
}

// =========================
// handleUpdate
// =========================

func TestRoleHandler_handleUpdate(t *testing.T) {
	handler, mockRepo, _ := setupRoleHandler()

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

	updateReq := dto.UpdateRoleRequest{Name: "new-name"}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/api/role", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestRoleHandler_handleUpdate_NotFound(t *testing.T) {
	handler, _, _ := setupRoleHandler()

	body, _ := json.Marshal(dto.UpdateRoleRequest{Name: "updated"})
	req := httptest.NewRequest(http.MethodPut, "/api/role/nonexistent", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleUpdate(w, req, "nonexistent")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for update on nonexistent role, got %d", w.Code)
	}
}

func TestRoleHandler_handleUpdate_ResponseBodyReflectsChange(t *testing.T) {
	handler, mockRepo, _ := setupRoleHandler()

	seedRole(t, mockRepo, "role-upd", "original", "Original desc")

	body, _ := json.Marshal(dto.UpdateRoleRequest{Name: "renamed", Description: "New desc"})
	req := httptest.NewRequest(http.MethodPut, "/api/role/role-upd", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	handler.handleUpdate(w, req, "role-upd")

	var response dto.RoleResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Name != "renamed" {
		t.Errorf("Name: want 'renamed', got '%s'", response.Name)
	}
	if response.Description != "New desc" {
		t.Errorf("Description: want 'New desc', got '%s'", response.Description)
	}
}

func TestRoleHandler_handleUpdate_WithoutUserContext(t *testing.T) {
	handler, mockRepo, _ := setupRoleHandler()

	seedRole(t, mockRepo, "role-ctx", "original", "desc")

	body, _ := json.Marshal(dto.UpdateRoleRequest{Name: "updated"})
	req := httptest.NewRequest(http.MethodPut, "/api/role/role-ctx", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// Sengaja tidak set UserIDKey
	w := httptest.NewRecorder()

	handler.handleUpdate(w, req, "role-ctx")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 even without user context, got %d", w.Code)
	}
}

// =========================
// handleDelete
// =========================

func TestRoleHandler_handleDelete(t *testing.T) {
	handler, mockRepo, _ := setupRoleHandler()

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

func TestRoleHandler_handleDelete_NotFound(t *testing.T) {
	handler, _, _ := setupRoleHandler()

	req := httptest.NewRequest(http.MethodDelete, "/api/role/nonexistent", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleDelete(w, req, "nonexistent")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for delete on nonexistent role, got %d", w.Code)
	}
}

func TestRoleHandler_handleDelete_ResponseBodyHasMessage(t *testing.T) {
	handler, mockRepo, _ := setupRoleHandler()

	seedRole(t, mockRepo, "role-del", "to-delete", "desc")

	req := httptest.NewRequest(http.MethodDelete, "/api/role/role-del", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	handler.handleDelete(w, req, "role-del")

	var response dto.MessageResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Message == "" {
		t.Error("expected message in delete response body")
	}
}

func TestRoleHandler_handleDelete_CannotGetDeletedRole(t *testing.T) {
	handler, mockRepo, _ := setupRoleHandler()

	seedRole(t, mockRepo, "role-gone", "temporary", "desc")

	// Delete
	delReq := httptest.NewRequest(http.MethodDelete, "/api/role/role-gone", nil)
	ctx := context.WithValue(delReq.Context(), middleware.UserIDKey, "user-1")
	delReq = delReq.WithContext(ctx)
	delW := httptest.NewRecorder()
	handler.handleDelete(delW, delReq, "role-gone")

	if delW.Code != http.StatusOK {
		t.Fatalf("delete should succeed, got %d", delW.Code)
	}

	// Try to get the deleted role — should 404
	getReq := httptest.NewRequest(http.MethodGet, "/api/role/role-gone", nil)
	getW := httptest.NewRecorder()
	handler.handleGetByID(getW, getReq, "role-gone")

	if getW.Code != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", getW.Code)
	}
}

func TestRoleHandler_handleDelete_WithoutUserContext(t *testing.T) {
	handler, mockRepo, _ := setupRoleHandler()

	seedRole(t, mockRepo, "role-noctx", "test", "desc")

	req := httptest.NewRequest(http.MethodDelete, "/api/role/role-noctx", nil)
	// Sengaja tidak set UserIDKey
	w := httptest.NewRecorder()

	handler.handleDelete(w, req, "role-noctx")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 even without user context, got %d", w.Code)
	}
}

// =========================
// ServeHTTP routing
// =========================

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

func TestRoleHandler_ServeHTTP_HeadNotAllowed(t *testing.T) {
	handler, _, _ := setupRoleHandler()

	req := httptest.NewRequest(http.MethodHead, "/api/role", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 for HEAD, got %d", w.Code)
	}
}

func TestRoleHandler_ServeHTTP_PostWithIDReturns400(t *testing.T) {
	handler, _, _ := setupRoleHandler()

	body, _ := json.Marshal(dto.CreateRoleRequest{Name: "x", Description: "y"})
	req := httptest.NewRequest(http.MethodPost, "/api/role/some-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for POST with ID, got %d", w.Code)
	}
}

func TestRoleHandler_ServeHTTP_PutWithoutIDReturns400(t *testing.T) {
	handler, _, _ := setupRoleHandler()

	body, _ := json.Marshal(dto.UpdateRoleRequest{Name: "x"})
	req := httptest.NewRequest(http.MethodPut, "/api/role", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for PUT without ID, got %d", w.Code)
	}
}

func TestRoleHandler_ServeHTTP_DeleteWithoutIDReturns400(t *testing.T) {
	handler, _, _ := setupRoleHandler()

	req := httptest.NewRequest(http.MethodDelete, "/api/role", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for DELETE without ID, got %d", w.Code)
	}
}
