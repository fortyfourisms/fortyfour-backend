package handlers

import (
	"bytes"
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/testhelpers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupAuthHandler() (*AuthHandler, *testhelpers.MockRedisClient) {
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := services.NewTokenService(redis, "test-secret")
	authService := services.NewAuthService(userRepo, tokenService)
	handler := NewAuthHandler(authService, tokenService)

	return handler, redis
}

func TestAuthHandler_Register_Success(t *testing.T) {
	// Arrange
	handler, _ := setupAuthHandler()

	reqBody := dto.RegisterRequest{
		Username: "testuser",
		Password: "P@sJord121",
		Email:    "test@example.com",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.Register(w, req)

	// Assert
	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var response dto.AuthResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.User.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", response.User.Username)
	}

	if response.AccessToken == "" {
		t.Error("expected access token in response")
	}

	if response.RefreshToken == "" {
		t.Error("expected refresh token in response")
	}

	if response.ExpiresAt == "" {
		t.Error("expected expires_at in response")
	}
}

func TestAuthHandler_Login_Success(t *testing.T) {
	// Arrange
	handler, _ := setupAuthHandler()

	// Register user first
	registerBody := map[string]string{
		"username": "testuser",
		"password": "P@ssj0rd121",
		"email":    "test@example.com",
	}
	body, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Register(w, req)

	// Now login
	loginBody := map[string]string{
		"username": "testuser",
		"password": "P@ssj0rd121",
	}
	body, _ = json.Marshal(loginBody)
	req = httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	// Act
	handler.Login(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response dto.AuthResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.AccessToken == "" {
		t.Error("expected access token in response")
	}

	if response.RefreshToken == "" {
		t.Error("expected refresh token in response")
	}
}

func TestAuthHandler_RefreshToken_Success(t *testing.T) {
	// Arrange
	handler, _ := setupAuthHandler()

	// Register user and get tokens
	registerBody := map[string]string{
		"username": "testuser",
		"password": "P@ssj0rd121",
		"email":    "test@example.com",
	}
	body, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Register(w, req)

	var registerResponse dto.AuthResponse
	json.NewDecoder(w.Body).Decode(&registerResponse)

	// Now refresh token
	refreshBody := map[string]string{
		"refresh_token": registerResponse.RefreshToken,
	}
	body, _ = json.Marshal(refreshBody)
	req = httptest.NewRequest(http.MethodPost, "/api/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	// Act
	handler.RefreshToken(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)

	if response["access_token"] == "" {
		t.Error("expected new access token in response")
	}

	if response["refresh_token"] == "" {
		t.Error("expected new refresh token in response")
	}
}

func TestAuthHandler_RefreshToken_InvalidToken(t *testing.T) {
	// Arrange
	handler, _ := setupAuthHandler()

	refreshBody := map[string]string{
		"refresh_token": "invalid-token",
	}
	body, _ := json.Marshal(refreshBody)
	req := httptest.NewRequest(http.MethodPost, "/api/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.RefreshToken(w, req)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthHandler_Logout_Success(t *testing.T) {
	// Arrange
	handler, _ := setupAuthHandler()

	// Register user and get tokens
	registerBody := map[string]string{
		"username": "testuser",
		"password": "P@ssj0rd121",
		"email":    "test@example.com",
	}
	body, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Register(w, req)

	var registerResponse dto.AuthResponse
	json.NewDecoder(w.Body).Decode(&registerResponse)

	// Now logout
	logoutBody := map[string]string{
		"refresh_token": registerResponse.RefreshToken,
	}
	body, _ = json.Marshal(logoutBody)
	req = httptest.NewRequest(http.MethodPost, "/api/logout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	// Act
	handler.Logout(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Try to use the revoked token
	refreshBody := map[string]string{
		"refresh_token": registerResponse.RefreshToken,
	}
	body, _ = json.Marshal(refreshBody)
	req = httptest.NewRequest(http.MethodPost, "/api/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.RefreshToken(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Error("expected revoked token to be rejected")
	}
}
