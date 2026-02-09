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
	perusahaanService := testhelpers.NewMockPerusahaanService()
	handler := NewAuthHandler(authService, tokenService, perusahaanService)

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

	// Verify Cookies (tokens should be in cookies, not response body)
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

func TestAuthHandler_Login_Success_WithMFASetupRequired(t *testing.T) {
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

	// Now login (user baru harus setup MFA dulu)
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

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)

	// User baru harus setup MFA dulu
	if mfaRequired, ok := response["mfa_setup_required"].(bool); !ok || !mfaRequired {
		t.Error("expected mfa_setup_required in response for new user")
	}

	if setupToken, ok := response["setup_token"].(string); !ok || setupToken == "" {
		t.Error("expected setup_token in response")
	}

	if message, ok := response["message"].(string); !ok || message == "" {
		t.Error("expected message in response")
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

/* ===================== MFA TESTS ===================== */

func TestAuthHandler_SetupMFA_WithSetupToken(t *testing.T) {
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

	// Login to get setup token
	loginBody := dto.LoginRequest{
		Username: "testuser",
		Password: "P@ssj0rd121",
	}
	body, _ = json.Marshal(loginBody)
	req = httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.Login(w, req)

	var loginResponse map[string]interface{}
	json.NewDecoder(w.Body).Decode(&loginResponse)
	setupToken := loginResponse["setup_token"].(string)

	// Setup MFA
	setupBody := map[string]string{
		"setup_token": setupToken,
	}
	body, _ = json.Marshal(setupBody)
	req = httptest.NewRequest(http.MethodPost, "/api/mfa/setup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	// Act
	handler.SetupMFA(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)

	if response["provisioning_uri"] == "" {
		t.Error("expected provisioning_uri in response")
	}
	if response["secret"] == "" {
		t.Error("expected secret in response")
	}
}

func TestAuthHandler_EnableMFA_Success(t *testing.T) {
	// This test would require mocking TOTP validation
	// Skipping detailed implementation for brevity
	t.Skip("Requires TOTP mocking")
}

func TestAuthHandler_VerifyMFA_Success(t *testing.T) {
	// This test would require mocking TOTP validation and MFA flow
	// Skipping detailed implementation for brevity
	t.Skip("Requires TOTP mocking and full MFA flow")
}