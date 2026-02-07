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
	tokenService := services.NewTokenService(redis, "test-secret", false, "localhost")
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
	if w.Code != http.StatusOK { // Changed to 200 based on implementation
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify User in Response
	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)

	userMap, ok := response["user"].(map[string]interface{})
	if !ok {
		t.Fatal("expected user object in response")
	}

	if userMap["username"] != "testuser" {
		t.Errorf("expected username 'testuser', got '%v'", userMap["username"])
	}

	// Verify Cookies
	cookies := w.Result().Cookies()
	var accessToken, refreshToken string
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			accessToken = cookie.Value
		}
		if cookie.Name == "refresh_token" {
			refreshToken = cookie.Value
		}
	}

	if accessToken == "" {
		t.Error("expected access_token cookie")
	}
	if refreshToken == "" {
		t.Error("expected refresh_token cookie")
	}
}

func TestAuthHandler_Login_Success(t *testing.T) {
	// Arrange
	handler, _ := setupAuthHandler()

	// Register user first
	registerBody := dto.RegisterRequest{
		Username: "testuser",
		Password: "P@ssj0rd121",
		Email:    "test@example.com",
	}
	body, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Register(w, req)

	// Now login
	loginBody := dto.LoginRequest{
		Username: "testuser",
		Password: "P@ssj0rd121",
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

	// Verify Cookies
	cookies := w.Result().Cookies()
	var accessToken, refreshToken string
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			accessToken = cookie.Value
		}
		if cookie.Name == "refresh_token" {
			refreshToken = cookie.Value
		}
	}

	if accessToken == "" {
		t.Error("expected access_token cookie")
	}
	if refreshToken == "" {
		t.Error("expected refresh_token cookie")
	}
}

func TestAuthHandler_RefreshToken_Success(t *testing.T) {
	// Arrange
	handler, _ := setupAuthHandler()

	// Register user and get tokens
	registerBody := dto.RegisterRequest{
		Username: "testuser",
		Password: "P@ssj0rd121",
		Email:    "test@example.com",
	}
	body, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Register(w, req)

	// Extract refresh token from cookie
	cookies := w.Result().Cookies()
	var refreshToken string
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" {
			refreshToken = cookie.Value
		}
	}

	if refreshToken == "" {
		t.Fatal("failed to get refresh token from register response")
	}

	// Now refresh token
	req = httptest.NewRequest(http.MethodPost, "/api/refresh", nil) // No body needed
	// Add refresh token cookie
	req.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: refreshToken,
	})

	w = httptest.NewRecorder()

	// Act
	handler.Refresh(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Check new cookies
	cookies = w.Result().Cookies()
	var newAccessToken string
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			newAccessToken = cookie.Value
		}
	}

	if newAccessToken == "" {
		t.Error("expected new access_token cookie")
	}
}

func TestAuthHandler_RefreshToken_InvalidToken(t *testing.T) {
	// Arrange
	handler, _ := setupAuthHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/refresh", nil)
	// Add invalid refresh token cookie
	req.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: "invalid-token",
	})
	w := httptest.NewRecorder()

	// Act
	handler.Refresh(w, req)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthHandler_Logout_Success(t *testing.T) {
	// Arrange
	handler, _ := setupAuthHandler()

	// Register user
	registerBody := dto.RegisterRequest{
		Username: "testuser",
		Password: "P@ssj0rd121",
		Email:    "test@example.com",
	}
	body, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Register(w, req)

	// Extract refresh token from cookie
	cookies := w.Result().Cookies()
	var refreshToken string
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" {
			refreshToken = cookie.Value
		}
	}

	if refreshToken == "" {
		t.Fatal("failed to get refresh token")
	}

	// Now logout
	req = httptest.NewRequest(http.MethodPost, "/api/logout", nil)
	req.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: refreshToken,
	})
	w = httptest.NewRecorder()

	// Act
	handler.Logout(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify cookies are cleared (MaxAge < 0)
	logoutCookies := w.Result().Cookies()
	for _, cookie := range logoutCookies {
		if cookie.MaxAge >= 0 {
			t.Errorf("expected cookie %s to be expired/cleared", cookie.Name)
		}
	}

	// Try to use the revoked token
	req = httptest.NewRequest(http.MethodPost, "/api/refresh", nil)
	req.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: refreshToken,
	})
	w = httptest.NewRecorder()
	handler.Refresh(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Error("expected revoked token to be rejected")
	}
}
