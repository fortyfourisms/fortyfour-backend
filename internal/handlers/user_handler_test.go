package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/testhelpers"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func strPtr(s string) *string {
	return &s
}

func setupUserHandler() (*UserHandler, *testhelpers.MockUserRepository, *services.SSEService) {
	mockRepo := testhelpers.NewMockUserRepository()
	uploadPath := "./test_uploads"
	os.MkdirAll(uploadPath, os.ModePerm)
	sseService := services.NewSSEService()
	userService := services.NewUserService(mockRepo, uploadPath)
	handler := NewUserHandler(userService, uploadPath, sseService)
	return handler, mockRepo, sseService
}

func TestUserHandler_handleGetAll(t *testing.T) {
	handler, _, _ := setupUserHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	w := httptest.NewRecorder()

	handler.handleGetAll(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestUserHandler_handleGetByID(t *testing.T) {
	handler, mockRepo, _ := setupUserHandler()

	// Create test user
	user := &models.User{
		ID:        "test-id",
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.Create(user)

	req := httptest.NewRequest(http.MethodGet, "/api/users/test-id", nil)
	w := httptest.NewRecorder()

	handler.handleGetByID(w, req, "test-id")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response dto.UserResponse
	json.NewDecoder(w.Body).Decode(&response)
	if response.ID != "test-id" {
		t.Errorf("expected ID 'test-id', got '%s'", response.ID)
	}
}

func TestUserHandler_handleGetByID_NotFound(t *testing.T) {
	handler, _, _ := setupUserHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/users/nonexistent", nil)
	w := httptest.NewRecorder()

	handler.handleGetByID(w, req, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestUserHandler_handleCreate(t *testing.T) {
	handler, _, _ := setupUserHandler()

	reqBody := dto.CreateUserRequest{
		Username: "newuser",
		Password: "P@sJ0rd121!",
		Email:    "newuser@example.com",
		RoleID:   strPtr("role-user"),
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "admin-id")
	ctx = context.WithValue(ctx, middleware.RoleKey, "admin")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.handleCreate(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var response dto.UserResponse
	json.NewDecoder(w.Body).Decode(&response)
	if response.Username != "newuser" {
		t.Errorf("expected username 'newuser', got '%s'", response.Username)
	}
}

func TestUserHandler_handleCreate_InvalidBody(t *testing.T) {
	handler, _, _ := setupUserHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.handleCreate(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestUserHandler_handleUpdate(t *testing.T) {
	handler, mockRepo, _ := setupUserHandler()

	// Create test user
	user := &models.User{
		ID:        "id1",
		Username:  "testuser",
		Password:  "P@sJ0rd121!",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.Create(user)

	updateReq := dto.UpdateUserRequest{
		Username: stringPtr("updateduser"),
		Email:    stringPtr("updated@example.com"),
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/api/users/id1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	ctx = context.WithValue(ctx, middleware.RoleKey, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleUpdate(w, req, "id1")

	// Assert
	assert.Equal(t, http.StatusOK, w.Code, "Response body: %s", w.Body.String())
}

func TestUserHandler_handleUpdate_Unauthorized(t *testing.T) {
	handler, mockRepo, _ := setupUserHandler()

	// Create test user
	user := &models.User{
		ID:        "test-id",
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.Create(user)

	updateReq := dto.UpdateUserRequest{
		Username: stringPtr("updateduser"),
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/api/users/test-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "other-user-id")
	ctx = context.WithValue(ctx, middleware.RoleKey, "user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleUpdate(w, req, "test-id")

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", w.Code)
	}
}

func TestUserHandler_handleUpdatePassword(t *testing.T) {
	handler, mockRepo, _ := setupUserHandler()

	// Create test user with hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("0ld!Pass#2023"), bcrypt.DefaultCost)
	user := &models.User{
		ID:        "test-id",
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.Create(user)

	updateReq := dto.UpdateUserPasswordRequest{
		OldPassword:        "0ld!Pass#2023",
		NewPassword:        "N3w@Strong$Pass2025",
		ConfirmNewPassword: "N3w@Strong$Pass2025",
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/users/test-id/password",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	// Set context userID & role
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "test-id")
	ctx = context.WithValue(ctx, middleware.RoleKey, "user")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.handleUpdatePassword(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestUserHandler_handleUpdatePassword_WrongUser(t *testing.T) {
	handler, mockRepo, _ := setupUserHandler()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("0ld!Pass#2023"), bcrypt.DefaultCost)
	user := &models.User{
		ID:       "test-id",
		Password: string(hashedPassword),
	}
	mockRepo.Create(user)

	updateReq := dto.UpdateUserPasswordRequest{
		OldPassword:        "0ld!Pass#2023",
		NewPassword:        "N3w@Strong$Pass2025",
		ConfirmNewPassword: "N3w@Strong$Pass2025",
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/users/test-id/password",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "other-user-id")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.handleUpdatePassword(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", w.Code)
	}
}

func TestUserHandler_handleDelete(t *testing.T) {
	handler, mockRepo, _ := setupUserHandler()

	// Create test user
	user := &models.User{
		ID:        "test-id",
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.Create(user)

	req := httptest.NewRequest(http.MethodDelete, "/api/users/test-id", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "admin-id")
	ctx = context.WithValue(ctx, middleware.RoleKey, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleDelete(w, req, "test-id")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestUserHandler_handleDelete_NotAdmin(t *testing.T) {
	handler, _, _ := setupUserHandler()

	req := httptest.NewRequest(http.MethodDelete, "/api/users/test-id", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-id")
	ctx = context.WithValue(ctx, middleware.RoleKey, "user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleDelete(w, req, "test-id")

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", w.Code)
	}
}

func TestUserHandler_ServeHTTP(t *testing.T) {
	handler, _, _ := setupUserHandler()

	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"GET all", http.MethodGet, "/api/users", http.StatusOK},
		{"GET by ID", http.MethodGet, "/api/users/test-id", http.StatusNotFound},
		{"POST create", http.MethodPost, "/api/users", http.StatusBadRequest},
		{"PUT update", http.MethodPut, "/api/users/test-id", http.StatusForbidden},
		{"DELETE", http.MethodDelete, "/api/users/test-id", http.StatusForbidden},
		{"Method not allowed", http.MethodPatch, "/api/users", http.StatusMethodNotAllowed},
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

func TestUserHandler_isValidImageType(t *testing.T) {
	handler, _, _ := setupUserHandler()

	testCases := []struct {
		filename string
		expected bool
	}{
		{"test.jpg", true},
		{"test.jpeg", true},
		{"test.png", true},
		{"test.JPG", true},
		{"test.PNG", true},
		{"test.gif", false},
		{"test.pdf", false},
		{"test", false},
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			result := handler.isValidImageType(tc.filename)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
