package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/testhelpers"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"strings"
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
	userService := services.NewUserService(mockRepo, uploadPath, nil)
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

/* =========================
   HELPERS
========================= */

func createUserMultipartRequest(method, url, fieldName, filename string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Buat file field dengan content-type image/jpeg
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, filename))
	h.Set("Content-Type", "image/jpeg")
	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, err
	}
	// Tulis dummy JPEG bytes (minimal valid)
	io.WriteString(part, "fakeimagecontent")
	writer.Close()

	req := httptest.NewRequest(method, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func withUserCtx(req *http.Request, userID, role string) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, middleware.UserIDKey, userID)
	ctx = context.WithValue(ctx, middleware.RoleKey, role)
	return req.WithContext(ctx)
}

/* =========================
   TEST handleUpdateProfilePhoto
========================= */

func TestUserHandler_handleUpdateProfilePhoto_Success(t *testing.T) {
	handler, mockRepo, _ := setupUserHandler()
	defer os.RemoveAll("./test_uploads")

	user := &models.User{
		ID:        "user-1",
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.Create(user)

	req, err := createUserMultipartRequest(http.MethodPost, "/api/users/user-1/profile-photo", "profile_photo", "photo.jpg")
	assert.NoError(t, err)
	req = withUserCtx(req, "user-1", "user")

	w := httptest.NewRecorder()
	handler.handleUpdateProfilePhoto(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_handleUpdateProfilePhoto_Forbidden(t *testing.T) {
	handler, _, _ := setupUserHandler()
	defer os.RemoveAll("./test_uploads")

	req, err := createUserMultipartRequest(http.MethodPost, "/api/users/user-other/profile-photo", "profile_photo", "photo.jpg")
	assert.NoError(t, err)
	req = withUserCtx(req, "user-1", "user") // user-1 coba update user-other

	w := httptest.NewRecorder()
	handler.handleUpdateProfilePhoto(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUserHandler_handleUpdateProfilePhoto_AdminCanUpdateOthers(t *testing.T) {
	handler, mockRepo, _ := setupUserHandler()
	defer os.RemoveAll("./test_uploads")

	user := &models.User{
		ID:        "user-target",
		Username:  "targetuser",
		Email:     "target@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.Create(user)

	req, err := createUserMultipartRequest(http.MethodPost, "/api/users/user-target/profile-photo", "profile_photo", "photo.jpg")
	assert.NoError(t, err)
	req = withUserCtx(req, "admin-1", "admin") // admin update user lain

	w := httptest.NewRecorder()
	handler.handleUpdateProfilePhoto(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_handleUpdateProfilePhoto_NoFile(t *testing.T) {
	handler, _, _ := setupUserHandler()
	defer os.RemoveAll("./test_uploads")

	// Request multipart tapi tanpa file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/users/user-1/profile-photo", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = withUserCtx(req, "user-1", "user")

	w := httptest.NewRecorder()
	handler.handleUpdateProfilePhoto(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_handleUpdateProfilePhoto_InvalidFormat(t *testing.T) {
	handler, _, _ := setupUserHandler()
	defer os.RemoveAll("./test_uploads")

	// File .gif - bukan format yang diizinkan
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="profile_photo"; filename="photo.gif"`)
	h.Set("Content-Type", "image/gif")
	part, _ := writer.CreatePart(h)
	io.WriteString(part, "fakegifcontent")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/users/user-1/profile-photo", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = withUserCtx(req, "user-1", "user")

	w := httptest.NewRecorder()
	handler.handleUpdateProfilePhoto(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Format file")
}

/* =========================
   TEST handleUpdateBanner
========================= */

func TestUserHandler_handleUpdateBanner_Success(t *testing.T) {
	handler, mockRepo, _ := setupUserHandler()
	defer os.RemoveAll("./test_uploads")

	user := &models.User{
		ID:        "user-1",
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.Create(user)

	req, err := createUserMultipartRequest(http.MethodPost, "/api/users/user-1/banner", "banner", "banner.jpg")
	assert.NoError(t, err)
	req = withUserCtx(req, "user-1", "user")

	w := httptest.NewRecorder()
	handler.handleUpdateBanner(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_handleUpdateBanner_Forbidden(t *testing.T) {
	handler, _, _ := setupUserHandler()
	defer os.RemoveAll("./test_uploads")

	req, err := createUserMultipartRequest(http.MethodPost, "/api/users/user-other/banner", "banner", "banner.jpg")
	assert.NoError(t, err)
	req = withUserCtx(req, "user-1", "user")

	w := httptest.NewRecorder()
	handler.handleUpdateBanner(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUserHandler_handleUpdateBanner_NoFile(t *testing.T) {
	handler, _, _ := setupUserHandler()
	defer os.RemoveAll("./test_uploads")

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/users/user-1/banner", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = withUserCtx(req, "user-1", "user")

	w := httptest.NewRecorder()
	handler.handleUpdateBanner(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_handleUpdateBanner_InvalidFormat(t *testing.T) {
	handler, _, _ := setupUserHandler()
	defer os.RemoveAll("./test_uploads")

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="banner"; filename="banner.gif"`)
	h.Set("Content-Type", "image/gif")
	part, _ := writer.CreatePart(h)
	io.WriteString(part, "fakegifcontent")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/users/user-1/banner", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = withUserCtx(req, "user-1", "user")

	w := httptest.NewRecorder()
	handler.handleUpdateBanner(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Format file")
}

/* =========================
   TEST handleUpdateStatus
========================= */

func TestUserHandler_handleUpdateStatus_Success(t *testing.T) {
	handler, mockRepo, _ := setupUserHandler()

	user := &models.User{
		ID:        "user-1",
		Username:  "targetuser",
		Email:     "target@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.Create(user)

	body := strings.NewReader(`{"status":"Suspend"}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/users/user-1/status", body)
	req = withUserCtx(req, "admin-1", "admin")

	w := httptest.NewRecorder()
	handler.handleUpdateStatus(w, req, "user-1")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Suspend")
}

func TestUserHandler_handleUpdateStatus_Forbidden(t *testing.T) {
	handler, _, _ := setupUserHandler()

	body := strings.NewReader(`{"status":"Suspend"}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/users/user-1/status", body)
	req = withUserCtx(req, "user-1", "user") // bukan admin

	w := httptest.NewRecorder()
	handler.handleUpdateStatus(w, req, "user-1")

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUserHandler_handleUpdateStatus_InvalidBody(t *testing.T) {
	handler, _, _ := setupUserHandler()

	req := httptest.NewRequest(http.MethodPatch, "/api/users/user-1/status", strings.NewReader("invalid json"))
	req = withUserCtx(req, "admin-1", "admin")

	w := httptest.NewRecorder()
	handler.handleUpdateStatus(w, req, "user-1")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_handleUpdateStatus_InvalidStatus(t *testing.T) {
	handler, mockRepo, _ := setupUserHandler()

	user := &models.User{
		ID:        "user-1",
		Username:  "targetuser",
		Email:     "target@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.Create(user)

	body := strings.NewReader(`{"status":"InvalidStatus"}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/users/user-1/status", body)
	req = withUserCtx(req, "admin-1", "admin")

	w := httptest.NewRecorder()
	handler.handleUpdateStatus(w, req, "user-1")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

/* =========================
   TEST ServeHTTP — routing guards
========================= */

// POST /api/users/{id} → 400 karena ID tidak diperlukan untuk create
func TestUserHandler_ServeHTTP_PostWithID_Returns400(t *testing.T) {
	handler, _, _ := setupUserHandler()

	body := strings.NewReader(`{"username":"newuser","email":"new@example.com","password":"P@ss123!"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/users/some-id", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "ID tidak diperlukan")
}

// PUT /api/users (tanpa ID) → 400
func TestUserHandler_ServeHTTP_PutWithoutID_Returns400(t *testing.T) {
	handler, _, _ := setupUserHandler()

	body := strings.NewReader(`{"username":"updateduser"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "ID wajib")
}

// DELETE /api/users (tanpa ID) → 400
func TestUserHandler_ServeHTTP_DeleteWithoutID_Returns400(t *testing.T) {
	handler, _, _ := setupUserHandler()

	req := httptest.NewRequest(http.MethodDelete, "/api/users", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "ID wajib")
}

/* =========================
   TEST handleUpdatePassword — cabang yang belum dicakup
========================= */

// Method selain PUT ke /password → 405
func TestUserHandler_handleUpdatePassword_MethodNotAllowed(t *testing.T) {
	handler, _, _ := setupUserHandler()

	for _, method := range []string{http.MethodPost, http.MethodGet, http.MethodDelete, http.MethodPatch} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/users/user-1/password", nil)
			req = withUserCtx(req, "user-1", "user")
			w := httptest.NewRecorder()

			handler.handleUpdatePassword(w, req)

			assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		})
	}
}

// Body tidak valid JSON → 400
func TestUserHandler_handleUpdatePassword_InvalidBody(t *testing.T) {
	handler, _, _ := setupUserHandler()

	req := httptest.NewRequest(http.MethodPut, "/api/users/user-1/password", strings.NewReader("not-json"))
	req.Header.Set("Content-Type", "application/json")
	req = withUserCtx(req, "user-1", "user")
	w := httptest.NewRecorder()

	handler.handleUpdatePassword(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}

// Body valid JSON tapi field kosong → validasi gagal → 400
func TestUserHandler_handleUpdatePassword_ValidationFails(t *testing.T) {
	handler, _, _ := setupUserHandler()

	// Kirim body dengan password kosong (setelah di-trim akan gagal validasi)
	body := strings.NewReader(`{"old_password":"","new_password":"","confirm_new_password":""}`)
	req := httptest.NewRequest(http.MethodPut, "/api/users/user-1/password", body)
	req.Header.Set("Content-Type", "application/json")
	req = withUserCtx(req, "user-1", "user")
	w := httptest.NewRecorder()

	handler.handleUpdatePassword(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

/* =========================
   TEST handleUpdate — cabang yang belum dicakup
========================= */

// Body tidak valid JSON → 400 (setelah lolos cek ownership)
func TestUserHandler_handleUpdate_InvalidBody(t *testing.T) {
	handler, _, _ := setupUserHandler()

	req := httptest.NewRequest(http.MethodPut, "/api/users/user-1", strings.NewReader("not-json"))
	req.Header.Set("Content-Type", "application/json")
	req = withUserCtx(req, "user-1", "user") // user update dirinya sendiri → lolos ownership check
	w := httptest.NewRecorder()

	handler.handleUpdate(w, req, "user-1")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}

// Non-admin mencoba mengubah RoleID → 403
func TestUserHandler_handleUpdate_NonAdminChangeRole_Returns403(t *testing.T) {
	handler, mockRepo, _ := setupUserHandler()

	user := &models.User{
		ID:        "user-1",
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.Create(user)

	roleID := "role-admin"
	updateReq := dto.UpdateUserRequest{
		RoleID: &roleID,
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/api/users/user-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = withUserCtx(req, "user-1", "user") // bukan admin
	w := httptest.NewRecorder()

	handler.handleUpdate(w, req, "user-1")

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "tidak bisa mengubah role")
}

/* =========================
   TEST handleUpdateProfilePhoto / handleUpdateBanner — Method Not Allowed
========================= */

func TestUserHandler_handleUpdateProfilePhoto_MethodNotAllowed(t *testing.T) {
	handler, _, _ := setupUserHandler()

	for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPatch} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/users/user-1/profile-photo", nil)
			req = withUserCtx(req, "user-1", "user")
			w := httptest.NewRecorder()

			handler.handleUpdateProfilePhoto(w, req)

			assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		})
	}
}

func TestUserHandler_handleUpdateBanner_MethodNotAllowed(t *testing.T) {
	handler, _, _ := setupUserHandler()

	for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPatch} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/users/user-1/banner", nil)
			req = withUserCtx(req, "user-1", "user")
			w := httptest.NewRecorder()

			handler.handleUpdateBanner(w, req)

			assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		})
	}
}

/* =========================
   TEST handleUpdateStatus — ID kosong
========================= */

// ID kosong → 400
func TestUserHandler_handleUpdateStatus_EmptyID_Returns400(t *testing.T) {
	handler, _, _ := setupUserHandler()

	body := strings.NewReader(`{"status":"Suspend"}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/users//status", body)
	req = withUserCtx(req, "admin-1", "admin")
	w := httptest.NewRecorder()

	handler.handleUpdateStatus(w, req, "") // id="" disimulasikan langsung

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "ID user wajib diisi")
}

/* =========================
   TEST handleDelete — service error
========================= */

// Delete user yang tidak ada → service kembalikan error → 400
func TestUserHandler_handleDelete_ServiceError(t *testing.T) {
	handler, _, _ := setupUserHandler()
	// mockRepo kosong, tidak ada user → Delete akan error
	req := httptest.NewRequest(http.MethodDelete, "/api/users/nonexistent-id", nil)
	req = withUserCtx(req, "admin-1", "admin")
	w := httptest.NewRecorder()

	handler.handleDelete(w, req, "nonexistent-id")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}