package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/testhelpers"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
)

func setupAuthHandler() (*AuthHandler, *testhelpers.MockRedisClient) {
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := services.NewTokenService(redis, "test-secret", false, "localhost")
	authService := services.NewAuthService(userRepo, testhelpers.NewMockRoleRepositoryWithDefaults(), tokenService, services.NewNotificationService(redis))
	perusahaanService := testhelpers.NewMockPerusahaanService()
	handler := NewAuthHandler(authService, tokenService, perusahaanService, nil, "")

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
		Identifier: "testuser",
		Password:   "P@ssj0rd121",
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
		Identifier: "testuser",
		Password:   "P@ssj0rd121",
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
	// Arrange
	handler, _ := setupAuthHandler()

	// 1. Register user
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

	// 2. Login to get setup_token
	loginBody := dto.LoginRequest{
		Identifier: "testuser",
		Password:   "P@ssj0rd121",
	}
	body, _ = json.Marshal(loginBody)
	req = httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.Login(w, req)

	var loginResponse map[string]interface{}
	json.NewDecoder(w.Body).Decode(&loginResponse)
	setupToken := loginResponse["setup_token"].(string)

	// 3. Setup MFA to get secret
	setupBody := map[string]string{
		"setup_token": setupToken,
	}
	body, _ = json.Marshal(setupBody)
	req = httptest.NewRequest(http.MethodPost, "/api/mfa/setup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.SetupMFA(w, req)

	var setupResponse map[string]string
	json.NewDecoder(w.Body).Decode(&setupResponse)
	secret := setupResponse["secret"]

	// 4. Generate valid TOTP code
	code := generateValidTOTPCode(secret)

	// 5. Enable MFA with valid code
	enableBody := map[string]interface{}{
		"code":        code,
		"setup_token": setupToken,
	}
	body, _ = json.Marshal(enableBody)
	req = httptest.NewRequest(http.MethodPost, "/api/mfa/enable", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	// Act
	handler.EnableMFA(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify response
	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)

	if response["message"] != "MFA enabled successfully" {
		t.Errorf("expected success message, got %v", response["message"])
	}

	// Verify cookies are set
	cookies := w.Result().Cookies()
	var hasAccessToken, hasRefreshToken bool
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			hasAccessToken = true
		}
		if cookie.Name == "refresh_token" {
			hasRefreshToken = true
		}
	}

	if !hasAccessToken {
		t.Error("expected access_token cookie")
	}
	if !hasRefreshToken {
		t.Error("expected refresh_token cookie")
	}
}

func TestAuthHandler_VerifyMFA_Success(t *testing.T) {
	// Arrange
	handler, _ := setupAuthHandler()

	// 1. Register user
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

	// 2. Login to get setup_token
	loginBody := dto.LoginRequest{
		Identifier: "testuser",
		Password:   "P@ssj0rd121",
	}
	body, _ = json.Marshal(loginBody)
	req = httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.Login(w, req)

	var loginResponse map[string]interface{}
	json.NewDecoder(w.Body).Decode(&loginResponse)
	setupToken := loginResponse["setup_token"].(string)

	// 3. Setup MFA
	setupBody := map[string]string{
		"setup_token": setupToken,
	}
	body, _ = json.Marshal(setupBody)
	req = httptest.NewRequest(http.MethodPost, "/api/mfa/setup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.SetupMFA(w, req)

	var setupResponse map[string]string
	json.NewDecoder(w.Body).Decode(&setupResponse)
	secret := setupResponse["secret"]

	// 4. Enable MFA
	code := generateValidTOTPCode(secret)
	enableBody := map[string]interface{}{
		"code":        code,
		"setup_token": setupToken,
	}
	body, _ = json.Marshal(enableBody)
	req = httptest.NewRequest(http.MethodPost, "/api/mfa/enable", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.EnableMFA(w, req)

	// 5. Login again (user now has MFA enabled) to get mfa_token
	body, _ = json.Marshal(loginBody)
	req = httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.Login(w, req)

	var mfaResponse map[string]interface{}
	json.NewDecoder(w.Body).Decode(&mfaResponse)
	mfaToken := mfaResponse["mfa_token"].(string)

	// 6. Generate new valid TOTP code for verification
	verifyCode := generateValidTOTPCode(secret)

	// 7. Verify MFA
	verifyBody := map[string]string{
		"mfa_token": mfaToken,
		"code":      verifyCode,
	}
	body, _ = json.Marshal(verifyBody)
	req = httptest.NewRequest(http.MethodPost, "/api/mfa/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	// Act
	handler.VerifyMFA(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify response
	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)

	if response["message"] != "MFA verification successful" {
		t.Errorf("expected success message, got %v", response["message"])
	}

	// Verify cookies are set
	cookies := w.Result().Cookies()
	var hasAccessToken, hasRefreshToken bool
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			hasAccessToken = true
		}
		if cookie.Name == "refresh_token" {
			hasRefreshToken = true
		}
	}

	if !hasAccessToken {
		t.Error("expected access_token cookie")
	}
	if !hasRefreshToken {
		t.Error("expected refresh_token cookie")
	}
}

/* ===================== MFA ERROR TESTS ===================== */

func TestAuthHandler_EnableMFA_InvalidCode(t *testing.T) {
	// Arrange
	handler, _ := setupAuthHandler()

	// Setup: Register, Login, SetupMFA
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

	loginBody := dto.LoginRequest{
		Identifier: "testuser",
		Password:   "P@ssj0rd121",
	}
	body, _ = json.Marshal(loginBody)
	req = httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.Login(w, req)

	var loginResponse map[string]interface{}
	json.NewDecoder(w.Body).Decode(&loginResponse)
	setupToken := loginResponse["setup_token"].(string)

	setupBody := map[string]string{"setup_token": setupToken}
	body, _ = json.Marshal(setupBody)
	req = httptest.NewRequest(http.MethodPost, "/api/mfa/setup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.SetupMFA(w, req)

	// Act: Try to enable MFA with invalid code
	enableBody := map[string]interface{}{
		"code":        "000000", // Invalid code
		"setup_token": setupToken,
	}
	body, _ = json.Marshal(enableBody)
	req = httptest.NewRequest(http.MethodPost, "/api/mfa/enable", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.EnableMFA(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_EnableMFA_ExpiredSetupToken(t *testing.T) {
	// Arrange
	handler, redis := setupAuthHandler()

	// Setup: Register and Login
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

	loginBody := dto.LoginRequest{
		Identifier: "testuser",
		Password:   "P@ssj0rd121",
	}
	body, _ = json.Marshal(loginBody)
	req = httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.Login(w, req)

	var loginResponse map[string]interface{}
	json.NewDecoder(w.Body).Decode(&loginResponse)
	setupToken := loginResponse["setup_token"].(string)

	// Setup MFA first to create the mfa_setup key
	setupBody := map[string]string{"setup_token": setupToken}
	body, _ = json.Marshal(setupBody)
	req = httptest.NewRequest(http.MethodPost, "/api/mfa/setup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.SetupMFA(w, req)

	// Simulate expired setup token by deleting the mfa_setup_token key from Redis
	redis.Delete("mfa_setup_token:" + setupToken)

	// Act: Try to enable MFA with expired token
	enableBody := map[string]interface{}{
		"code":        "123456",
		"setup_token": setupToken,
	}
	body, _ = json.Marshal(enableBody)
	req = httptest.NewRequest(http.MethodPost, "/api/mfa/enable", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.EnableMFA(w, req)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthHandler_EnableMFA_NoToken(t *testing.T) {
	// Arrange
	handler, _ := setupAuthHandler()

	// Act: Try to enable MFA without setup_token or authentication
	enableBody := map[string]string{
		"code": "123456",
	}
	body, _ := json.Marshal(enableBody)
	req := httptest.NewRequest(http.MethodPost, "/api/mfa/enable", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.EnableMFA(w, req)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthHandler_EnableMFA_InvalidBody(t *testing.T) {
	// Arrange
	handler, _ := setupAuthHandler()

	// Act: Send invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/api/mfa/enable", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.EnableMFA(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_VerifyMFA_InvalidCode(t *testing.T) {
	// Arrange
	handler, redis := setupAuthHandler()

	// Simulate a user with MFA enabled and mfa_pending token
	mfaToken := "test-mfa-token"
	userID := "user-123"
	redis.Set("mfa_pending:"+mfaToken, userID, 0)

	// Act: Try to verify with invalid code
	verifyBody := map[string]string{
		"mfa_token": mfaToken,
		"code":      "000000", // Invalid code
	}
	body, _ := json.Marshal(verifyBody)
	req := httptest.NewRequest(http.MethodPost, "/api/mfa/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.VerifyMFA(w, req)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthHandler_VerifyMFA_ExpiredMFAToken(t *testing.T) {
	// Arrange
	handler, _ := setupAuthHandler()

	// Act: Try to verify with non-existent/expired mfa_token
	verifyBody := map[string]string{
		"mfa_token": "expired-token",
		"code":      "123456",
	}
	body, _ := json.Marshal(verifyBody)
	req := httptest.NewRequest(http.MethodPost, "/api/mfa/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.VerifyMFA(w, req)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthHandler_VerifyMFA_InvalidBody(t *testing.T) {
	// Arrange
	handler, _ := setupAuthHandler()

	// Act: Send invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/api/mfa/verify", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.VerifyMFA(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_VerifyMFA_MFANotConfigured(t *testing.T) {
	// Arrange
	handler, redis := setupAuthHandler()

	// Register a user without MFA
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

	// Get user ID from response
	var registerResponse map[string]interface{}
	json.NewDecoder(w.Body).Decode(&registerResponse)
	userMap := registerResponse["user"].(map[string]interface{})
	userID := userMap["id"].(string)

	// Simulate mfa_pending token for user without MFA configured
	mfaToken := "test-mfa-token"
	redis.Set("mfa_pending:"+mfaToken, userID, 0)

	// Act: Try to verify MFA for user without MFA configured
	verifyBody := map[string]string{
		"mfa_token": mfaToken,
		"code":      "123456",
	}
	body, _ = json.Marshal(verifyBody)
	req = httptest.NewRequest(http.MethodPost, "/api/mfa/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.VerifyMFA(w, req)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

/* ===================== HELPER FUNCTIONS ===================== */

// generateValidTOTPCode generates a valid TOTP code for testing
func generateValidTOTPCode(secret string) string {
	code, _ := totp.GenerateCode(secret, time.Now())
	return code
}

/*
=====================================
 SETUP HELPER WITH USER SERVICE
=====================================
*/

func withMeUserContext(req *http.Request, userID string) *http.Request {
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	return req.WithContext(ctx)
}

/*
=====================================
 TEST GET ME
=====================================
*/

func TestAuthHandler_GetMe_Success(t *testing.T) {
	uploadPath := t.TempDir()
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := services.NewTokenService(redis, "test-secret", false, "localhost")
	authService := services.NewAuthService(userRepo, testhelpers.NewMockRoleRepositoryWithDefaults(), tokenService, services.NewNotificationService(redis))
	userService := services.NewUserService(userRepo, uploadPath, nil)
	handler := NewAuthHandler(authService, tokenService, testhelpers.NewMockPerusahaanService(), userService, uploadPath)

	user := testhelpers.CreateTestUser("user-1", "testuser", "test@test.com")
	_ = userRepo.Create(user)

	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req = withMeUserContext(req, "user-1")
	w := httptest.NewRecorder()

	handler.GetMe(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["username"] != "testuser" {
		t.Errorf("expected username 'testuser', got %v", resp["username"])
	}
}

func TestAuthHandler_GetMe_Unauthorized(t *testing.T) {
	handler, _ := setupAuthHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	// Tidak ada userID di context
	w := httptest.NewRecorder()

	handler.GetMe(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthHandler_GetMe_UserNotFound(t *testing.T) {
	uploadPath := t.TempDir()
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := services.NewTokenService(redis, "test-secret", false, "localhost")
	authService := services.NewAuthService(userRepo, testhelpers.NewMockRoleRepositoryWithDefaults(), tokenService, services.NewNotificationService(redis))
	userService := services.NewUserService(userRepo, uploadPath, nil)
	handler := NewAuthHandler(authService, tokenService, testhelpers.NewMockPerusahaanService(), userService, uploadPath)

	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req = withMeUserContext(req, "nonexistent-id")
	w := httptest.NewRecorder()

	handler.GetMe(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

/*
=====================================
 TEST UPDATE ME
=====================================
*/

func TestAuthHandler_UpdateMe_Success(t *testing.T) {
	uploadPath := t.TempDir()
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := services.NewTokenService(redis, "test-secret", false, "localhost")
	authService := services.NewAuthService(userRepo, testhelpers.NewMockRoleRepositoryWithDefaults(), tokenService, services.NewNotificationService(redis))
	userService := services.NewUserService(userRepo, uploadPath, nil)
	handler := NewAuthHandler(authService, tokenService, testhelpers.NewMockPerusahaanService(), userService, uploadPath)

	user := testhelpers.CreateTestUser("user-1", "oldname", "old@test.com")
	_ = userRepo.Create(user)

	newDisplayName := "New Display Name"
	body, _ := json.Marshal(dto.UpdateMeRequest{DisplayName: &newDisplayName})
	req := httptest.NewRequest(http.MethodPut, "/api/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = withMeUserContext(req, "user-1")
	w := httptest.NewRecorder()

	handler.UpdateMe(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["display_name"] != "New Display Name" {
		t.Errorf("expected display_name 'New Display Name', got %v", resp["display_name"])
	}
}

func TestAuthHandler_UpdateMe_WithIDJabatan(t *testing.T) {
	uploadPath := t.TempDir()
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := services.NewTokenService(redis, "test-secret", false, "localhost")
	authService := services.NewAuthService(userRepo, testhelpers.NewMockRoleRepositoryWithDefaults(), tokenService, services.NewNotificationService(redis))
	userService := services.NewUserService(userRepo, uploadPath, nil)
	handler := NewAuthHandler(authService, tokenService, testhelpers.NewMockPerusahaanService(), userService, uploadPath)

	user := testhelpers.CreateTestUser("user-1", "testuser", "test@test.com")
	_ = userRepo.Create(user)

	jabatanID := "jabatan-uuid-123"
	body, _ := json.Marshal(dto.UpdateMeRequest{IDJabatan: &jabatanID})
	req := httptest.NewRequest(http.MethodPut, "/api/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = withMeUserContext(req, "user-1")
	w := httptest.NewRecorder()

	handler.UpdateMe(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_UpdateMe_Unauthorized(t *testing.T) {
	handler, _ := setupAuthHandler()

	newDisplayName := "newname"
	body, _ := json.Marshal(dto.UpdateMeRequest{DisplayName: &newDisplayName})
	req := httptest.NewRequest(http.MethodPut, "/api/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.UpdateMe(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthHandler_UpdateMe_InvalidBody(t *testing.T) {
	handler, _ := setupAuthHandler()

	req := httptest.NewRequest(http.MethodPut, "/api/me", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	req = withMeUserContext(req, "user-1")
	w := httptest.NewRecorder()

	handler.UpdateMe(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestAuthHandler_UpdateMe_RoleIDNotUpdatable(t *testing.T) {
	uploadPath := t.TempDir()
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := services.NewTokenService(redis, "test-secret", false, "localhost")
	authService := services.NewAuthService(userRepo, testhelpers.NewMockRoleRepositoryWithDefaults(), tokenService, services.NewNotificationService(redis))
	userService := services.NewUserService(userRepo, uploadPath, nil)
	handler := NewAuthHandler(authService, tokenService, testhelpers.NewMockPerusahaanService(), userService, uploadPath)

	roleUser := "role-user"
	user := testhelpers.CreateTestUser("user-1", "testuser", "test@test.com")
	user.RoleID = &roleUser
	_ = userRepo.Create(user)

	// Kirim JSON dengan role_id — seharusnya diabaikan
	rawBody := `{"username": "newname", "role_id": "role-admin"}`
	req := httptest.NewRequest(http.MethodPut, "/api/me", strings.NewReader(rawBody))
	req.Header.Set("Content-Type", "application/json")
	req = withMeUserContext(req, "user-1")
	w := httptest.NewRecorder()

	handler.UpdateMe(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	// RoleID tidak boleh berubah
	updatedUser, _ := userRepo.FindByID("user-1")
	if updatedUser.RoleID == nil || *updatedUser.RoleID != "role-user" {
		t.Errorf("expected role_id unchanged 'role-user', got '%v'", updatedUser.RoleID)
	}
}

/*
=====================================
 TEST UPDATE ME PASSWORD
=====================================
*/

func TestAuthHandler_UpdateMePassword_Success(t *testing.T) {
	uploadPath := t.TempDir()
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := services.NewTokenService(redis, "test-secret", false, "localhost")
	authService := services.NewAuthService(userRepo, testhelpers.NewMockRoleRepositoryWithDefaults(), tokenService, services.NewNotificationService(redis))
	userService := services.NewUserService(userRepo, uploadPath, nil)
	handler := NewAuthHandler(authService, tokenService, testhelpers.NewMockPerusahaanService(), userService, uploadPath)

	// Register supaya password di-hash dengan benar
	reqBody := dto.RegisterRequest{
		Username: "testuser",
		Password: "Xk9#mP2$qL7!",
		Email:    "test@test.com",
	}
	body, _ := json.Marshal(reqBody)
	regReq := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
	regReq.Header.Set("Content-Type", "application/json")
	regW := httptest.NewRecorder()
	handler.Register(regW, regReq)

	if regW.Code != http.StatusCreated {
		t.Fatalf("register failed: %d %s", regW.Code, regW.Body.String())
	}

	var regResp map[string]interface{}
	json.NewDecoder(regW.Body).Decode(&regResp)
	userSection, ok := regResp["user"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected 'user' in register response, got: %v", regResp)
	}
	userID, ok := userSection["id"].(string)
	if !ok || userID == "" {
		t.Fatalf("expected user ID in register response, got: %v", userSection)
	}

	// Update password
	pwBody, _ := json.Marshal(dto.UpdateUserPasswordRequest{
		OldPassword:        "Xk9#mP2$qL7!",
		NewPassword:        "Rz4@wN8&vB3^",
		ConfirmNewPassword: "Rz4@wN8&vB3^",
	})
	req := httptest.NewRequest(http.MethodPut, "/api/me/password", bytes.NewBuffer(pwBody))
	req.Header.Set("Content-Type", "application/json")
	req = withMeUserContext(req, userID)
	w := httptest.NewRecorder()

	handler.UpdateMePassword(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["message"] != "Password berhasil diubah" {
		t.Errorf("unexpected message: %v", resp["message"])
	}
}

func TestAuthHandler_UpdateMePassword_Unauthorized(t *testing.T) {
	handler, _ := setupAuthHandler()

	body, _ := json.Marshal(dto.UpdateUserPasswordRequest{
		OldPassword: "old", NewPassword: "new12345", ConfirmNewPassword: "new12345",
	})
	req := httptest.NewRequest(http.MethodPut, "/api/me/password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.UpdateMePassword(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthHandler_UpdateMePassword_InvalidBody(t *testing.T) {
	handler, _ := setupAuthHandler()

	req := httptest.NewRequest(http.MethodPut, "/api/me/password", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	req = withMeUserContext(req, "user-1")
	w := httptest.NewRecorder()

	handler.UpdateMePassword(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestAuthHandler_UpdateMePassword_WrongOldPassword(t *testing.T) {
	uploadPath := t.TempDir()
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := services.NewTokenService(redis, "test-secret", false, "localhost")
	authService := services.NewAuthService(userRepo, testhelpers.NewMockRoleRepositoryWithDefaults(), tokenService, services.NewNotificationService(redis))
	userService := services.NewUserService(userRepo, uploadPath, nil)
	handler := NewAuthHandler(authService, tokenService, testhelpers.NewMockPerusahaanService(), userService, uploadPath)

	reqBody := dto.RegisterRequest{Username: "pwdtestuser", Password: "Xk9#mP2$qL7!", Email: "t@t.com"}
	body, _ := json.Marshal(reqBody)
	regReq := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
	regReq.Header.Set("Content-Type", "application/json")
	regW := httptest.NewRecorder()
	handler.Register(regW, regReq)

	if regW.Code != http.StatusCreated {
		t.Fatalf("register failed: %d %s", regW.Code, regW.Body.String())
	}

	var regResp map[string]interface{}
	json.NewDecoder(regW.Body).Decode(&regResp)
	userSection, ok := regResp["user"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected 'user' in register response, got: %v", regResp)
	}
	userID := userSection["id"].(string)

	pwBody, _ := json.Marshal(dto.UpdateUserPasswordRequest{
		OldPassword:        "Wr0ng#Pass!99",
		NewPassword:        "Rz4@wN8&vB3^",
		ConfirmNewPassword: "Rz4@wN8&vB3^",
	})
	req := httptest.NewRequest(http.MethodPut, "/api/me/password", bytes.NewBuffer(pwBody))
	req.Header.Set("Content-Type", "application/json")
	req = withMeUserContext(req, userID)
	w := httptest.NewRecorder()

	handler.UpdateMePassword(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

/*
=====================================
 TEST UPDATE ME MEDIA
=====================================
*/

func TestAuthHandler_UpdateMeMedia_Unauthorized(t *testing.T) {
	handler, _ := setupAuthHandler()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, "/api/me/media", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	handler.UpdateMeMedia(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthHandler_UpdateMeMedia_NoFileProvided(t *testing.T) {
	uploadPath := t.TempDir()
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := services.NewTokenService(redis, "test-secret", false, "localhost")
	authService := services.NewAuthService(userRepo, testhelpers.NewMockRoleRepositoryWithDefaults(), tokenService, services.NewNotificationService(redis))
	userService := services.NewUserService(userRepo, uploadPath, nil)
	handler := NewAuthHandler(authService, tokenService, testhelpers.NewMockPerusahaanService(), userService, uploadPath)

	user := testhelpers.CreateTestUser("user-1", "testuser", "test@test.com")
	_ = userRepo.Create(user)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, "/api/me/media", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = withMeUserContext(req, "user-1")
	w := httptest.NewRecorder()

	handler.UpdateMeMedia(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["error"] == "" {
		t.Error("expected error message")
	}
}

func TestAuthHandler_UpdateMeMedia_InvalidPhotoFormat(t *testing.T) {
	uploadPath := t.TempDir()
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := services.NewTokenService(redis, "test-secret", false, "localhost")
	authService := services.NewAuthService(userRepo, testhelpers.NewMockRoleRepositoryWithDefaults(), tokenService, services.NewNotificationService(redis))
	userService := services.NewUserService(userRepo, uploadPath, nil)
	handler := NewAuthHandler(authService, tokenService, testhelpers.NewMockPerusahaanService(), userService, uploadPath)

	user := testhelpers.CreateTestUser("user-1", "testuser", "test@test.com")
	_ = userRepo.Create(user)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("profile_photo", "photo.gif") // format tidak valid
	io.WriteString(part, "dummy content")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/me/media", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = withMeUserContext(req, "user-1")
	w := httptest.NewRecorder()

	handler.UpdateMeMedia(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["error"] == "" || !strings.Contains(resp["error"], "format") {
		t.Errorf("expected format error, got: %v", resp["error"])
	}
}

func TestAuthHandler_UpdateMeMedia_UploadProfilePhoto_Success(t *testing.T) {
	uploadPath := t.TempDir()
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := services.NewTokenService(redis, "test-secret", false, "localhost")
	authService := services.NewAuthService(userRepo, testhelpers.NewMockRoleRepositoryWithDefaults(), tokenService, services.NewNotificationService(redis))
	userService := services.NewUserService(userRepo, uploadPath, nil)
	handler := NewAuthHandler(authService, tokenService, testhelpers.NewMockPerusahaanService(), userService, uploadPath)

	user := testhelpers.CreateTestUser("user-1", "testuser", "test@test.com")
	_ = userRepo.Create(user)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("profile_photo", "photo.jpg")
	io.WriteString(part, "dummy image content")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/me/media", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = withMeUserContext(req, "user-1")
	w := httptest.NewRecorder()

	handler.UpdateMeMedia(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_UpdateMeMedia_UploadBanner_Success(t *testing.T) {
	uploadPath := t.TempDir()
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := services.NewTokenService(redis, "test-secret", false, "localhost")
	authService := services.NewAuthService(userRepo, testhelpers.NewMockRoleRepositoryWithDefaults(), tokenService, services.NewNotificationService(redis))
	userService := services.NewUserService(userRepo, uploadPath, nil)
	handler := NewAuthHandler(authService, tokenService, testhelpers.NewMockPerusahaanService(), userService, uploadPath)

	user := testhelpers.CreateTestUser("user-1", "testuser", "test@test.com")
	_ = userRepo.Create(user)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("banner", "banner.png")
	io.WriteString(part, "dummy banner content")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/me/media", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = withMeUserContext(req, "user-1")
	w := httptest.NewRecorder()

	handler.UpdateMeMedia(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

/*
=====================================
 TEST LOGIN (missing cases)
=====================================
*/

func TestAuthHandler_Login_InvalidBody(t *testing.T) {
	handler, _ := setupAuthHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestAuthHandler_Login_ValidationFails(t *testing.T) {
	handler, _ := setupAuthHandler()

	// Identifier kosong — seharusnya gagal validasi
	loginBody := dto.LoginRequest{
		Identifier: "",
		Password:   "",
	}
	body, _ := json.Marshal(loginBody)
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestAuthHandler_Login_WrongCredentials(t *testing.T) {
	handler, _ := setupAuthHandler()

	// Register user dulu
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

	// Login dengan password salah
	loginBody := dto.LoginRequest{
		Identifier: "testuser",
		Password:   "Wr0ng#Password!",
	}
	body, _ = json.Marshal(loginBody)
	req = httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthHandler_Login_MFAEnabled_ReturnsMFAToken(t *testing.T) {
	handler, _ := setupAuthHandler()

	// 1. Register
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

	// 2. Login → dapat setup_token
	loginBody := dto.LoginRequest{Identifier: "testuser", Password: "P@ssj0rd121"}
	body, _ = json.Marshal(loginBody)
	req = httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.Login(w, req)
	var loginResp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&loginResp)
	setupToken := loginResp["setup_token"].(string)

	// 3. Setup MFA → dapat secret
	setupBody := map[string]string{"setup_token": setupToken}
	body, _ = json.Marshal(setupBody)
	req = httptest.NewRequest(http.MethodPost, "/api/mfa/setup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.SetupMFA(w, req)
	var setupResp map[string]string
	json.NewDecoder(w.Body).Decode(&setupResp)
	secret := setupResp["secret"]

	// 4. Enable MFA
	code := generateValidTOTPCode(secret)
	enableBody := map[string]interface{}{"code": code, "setup_token": setupToken}
	body, _ = json.Marshal(enableBody)
	req = httptest.NewRequest(http.MethodPost, "/api/mfa/enable", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.EnableMFA(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("enable MFA failed: %d %s", w.Code, w.Body.String())
	}

	// 5. Login kembali — MFA sudah aktif, harus dapat mfa_token bukan access_token
	body, _ = json.Marshal(loginBody)
	req = httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)

	if mfaRequired, ok := resp["mfa_required"].(bool); !ok || !mfaRequired {
		t.Error("expected mfa_required: true in response")
	}
	if mfaToken, ok := resp["mfa_token"].(string); !ok || mfaToken == "" {
		t.Error("expected non-empty mfa_token in response")
	}
	// Pastikan tidak ada access_token di cookie
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == "access_token" && cookie.MaxAge >= 0 {
			t.Error("expected no access_token cookie when MFA is required")
		}
	}
}

/*
=====================================
 TEST LOGOUT ALL
=====================================
*/

// func setupAuthHandlerWithUserService() (*AuthHandler, *testhelpers.MockRedisClient) {
// 	uploadPath := os.TempDir()
// 	userRepo := testhelpers.NewMockUserRepository()
// 	redis := testhelpers.NewMockRedisClient()
// 	tokenService := services.NewTokenService(redis, "test-secret", false, "localhost")
// 	authService := services.NewAuthService(userRepo, tokenService, services.NewNotificationService(redis))
// 	userService := services.NewUserService(userRepo, uploadPath, nil)
// 	handler := NewAuthHandler(authService, tokenService, testhelpers.NewMockPerusahaanService(), userService, uploadPath)
// 	return handler, redis
// }

func TestAuthHandler_LogoutAll_Success(t *testing.T) {
	handler, _ := setupAuthHandler()

	// Register user untuk dapat token valid
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

	var regResp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&regResp)
	userMap := regResp["user"].(map[string]interface{})
	userID := userMap["id"].(string)

	// Logout all dengan user ID di context
	req = httptest.NewRequest(http.MethodPost, "/api/logout-all", nil)
	req = withMeUserContext(req, userID)
	w = httptest.NewRecorder()

	handler.LogoutAll(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["message"] == "" {
		t.Error("expected message in response")
	}

	// Cookies harus di-clear
	for _, cookie := range w.Result().Cookies() {
		if cookie.MaxAge >= 0 {
			t.Errorf("expected cookie %s to be expired/cleared after LogoutAll", cookie.Name)
		}
	}
}

func TestAuthHandler_LogoutAll_Unauthorized(t *testing.T) {
	handler, _ := setupAuthHandler()

	// Request tanpa user ID di context
	req := httptest.NewRequest(http.MethodPost, "/api/logout-all", nil)
	w := httptest.NewRecorder()

	handler.LogoutAll(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthHandler_LogoutAll_RevokesAllTokens(t *testing.T) {
	handler, redis := setupAuthHandler()

	// Register dan dapat beberapa token
	registerBody := dto.RegisterRequest{
		Username: "testuser",
		Password: "P@ssj0rd121",
		Email:    "test@example.com",
	}
	body, _ := json.Marshal(registerBody)

	// Register → token pertama
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Register(w, req)
	var regResp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&regResp)
	userMap := regResp["user"].(map[string]interface{})
	userID := userMap["id"].(string)

	// Ambil refresh token pertama dari cookie
	var firstRefreshToken string
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == "refresh_token" {
			firstRefreshToken = cookie.Value
		}
	}

	// Simulasi token kedua langsung di Redis
	redis.Set("refresh_token:token-sesi-2:"+userID, userID, 0)

	// Logout all
	req = httptest.NewRequest(http.MethodPost, "/api/logout-all", nil)
	req = withMeUserContext(req, userID)
	w = httptest.NewRecorder()
	handler.LogoutAll(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("LogoutAll failed: %d", w.Code)
	}

	// Refresh token pertama tidak boleh bisa dipakai lagi
	req = httptest.NewRequest(http.MethodPost, "/api/refresh", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: firstRefreshToken})
	w = httptest.NewRecorder()
	handler.Refresh(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected revoked refresh token to be rejected (401), got %d", w.Code)
	}
}

/*
=====================================
 TEST ME ROUTER
=====================================
*/

func TestAuthHandler_MeRouter_GET_RoutesToGetMe(t *testing.T) {
	uploadPath := t.TempDir()
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := services.NewTokenService(redis, "test-secret", false, "localhost")
	authService := services.NewAuthService(userRepo, testhelpers.NewMockRoleRepositoryWithDefaults(), tokenService, services.NewNotificationService(redis))
	userService := services.NewUserService(userRepo, uploadPath, nil)
	handler := NewAuthHandler(authService, tokenService, testhelpers.NewMockPerusahaanService(), userService, uploadPath)

	user := testhelpers.CreateTestUser("user-1", "testuser", "test@test.com")
	_ = userRepo.Create(user)

	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req = withMeUserContext(req, "user-1")
	w := httptest.NewRecorder()

	handler.MeRouter(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["username"] != "testuser" {
		t.Errorf("expected username 'testuser', got %v", resp["username"])
	}
}

func TestAuthHandler_MeRouter_PUT_RoutesToUpdateMe(t *testing.T) {
	uploadPath := t.TempDir()
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := services.NewTokenService(redis, "test-secret", false, "localhost")
	authService := services.NewAuthService(userRepo, testhelpers.NewMockRoleRepositoryWithDefaults(), tokenService, services.NewNotificationService(redis))
	userService := services.NewUserService(userRepo, uploadPath, nil)
	handler := NewAuthHandler(authService, tokenService, testhelpers.NewMockPerusahaanService(), userService, uploadPath)

	user := testhelpers.CreateTestUser("user-1", "oldname", "old@test.com")
	_ = userRepo.Create(user)

	newDisplayName := "Updated Name"
	body, _ := json.Marshal(dto.UpdateMeRequest{DisplayName: &newDisplayName})
	req := httptest.NewRequest(http.MethodPut, "/api/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = withMeUserContext(req, "user-1")
	w := httptest.NewRecorder()

	handler.MeRouter(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["display_name"] != "Updated Name" {
		t.Errorf("expected display_name 'Updated Name', got %v", resp["display_name"])
	}
}

func TestAuthHandler_MeRouter_PUT_password_RoutesToUpdateMePassword(t *testing.T) {
	uploadPath := t.TempDir()
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := services.NewTokenService(redis, "test-secret", false, "localhost")
	authService := services.NewAuthService(userRepo, testhelpers.NewMockRoleRepositoryWithDefaults(), tokenService, services.NewNotificationService(redis))
	userService := services.NewUserService(userRepo, uploadPath, nil)
	handler := NewAuthHandler(authService, tokenService, testhelpers.NewMockPerusahaanService(), userService, uploadPath)

	// Register untuk mendapat user dengan password yang di-hash dengan benar
	reqBody := dto.RegisterRequest{Username: "meuser", Password: "Xk9#mP2$qL7!", Email: "me@test.com"}
	body, _ := json.Marshal(reqBody)
	regReq := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
	regReq.Header.Set("Content-Type", "application/json")
	regW := httptest.NewRecorder()
	handler.Register(regW, regReq)
	var regResp map[string]interface{}
	json.NewDecoder(regW.Body).Decode(&regResp)
	userMap := regResp["user"].(map[string]interface{})
	userID := userMap["id"].(string)

	pwBody, _ := json.Marshal(dto.UpdateUserPasswordRequest{
		OldPassword:        "Xk9#mP2$qL7!",
		NewPassword:        "Rz4@wN8&vB3^",
		ConfirmNewPassword: "Rz4@wN8&vB3^",
	})
	req := httptest.NewRequest(http.MethodPut, "/api/me/password", bytes.NewBuffer(pwBody))
	req.Header.Set("Content-Type", "application/json")
	req = withMeUserContext(req, userID)
	w := httptest.NewRecorder()

	handler.MeRouter(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_MeRouter_POST_media_RoutesToUpdateMeMedia(t *testing.T) {
	uploadPath := t.TempDir()
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := services.NewTokenService(redis, "test-secret", false, "localhost")
	authService := services.NewAuthService(userRepo, testhelpers.NewMockRoleRepositoryWithDefaults(), tokenService, services.NewNotificationService(redis))
	userService := services.NewUserService(userRepo, uploadPath, nil)
	handler := NewAuthHandler(authService, tokenService, testhelpers.NewMockPerusahaanService(), userService, uploadPath)

	user := testhelpers.CreateTestUser("user-1", "testuser", "test@test.com")
	_ = userRepo.Create(user)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("profile_photo", "photo.jpg")
	io.WriteString(part, "dummy image content")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/me/media", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = withMeUserContext(req, "user-1")
	w := httptest.NewRecorder()

	handler.MeRouter(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_MeRouter_UnknownPath_Returns404(t *testing.T) {
	handler, _ := setupAuthHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/me/unknown-sub-path", nil)
	req = withMeUserContext(req, "user-1")
	w := httptest.NewRecorder()

	handler.MeRouter(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestAuthHandler_MeRouter_MethodNotAllowed(t *testing.T) {
	handler, _ := setupAuthHandler()

	// Method DELETE pada /api/me tidak diizinkan
	req := httptest.NewRequest(http.MethodDelete, "/api/me", nil)
	req = withMeUserContext(req, "user-1")
	w := httptest.NewRecorder()

	handler.MeRouter(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

/*
=====================================
 TEST REGISTER (missing error cases)
=====================================
*/

func TestAuthHandler_Register_InvalidBody(t *testing.T) {
	handler, _ := setupAuthHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader("bukan json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestAuthHandler_Register_ValidationFails_EmptyFields(t *testing.T) {
	handler, _ := setupAuthHandler()

	// Semua field kosong
	reqBody := dto.RegisterRequest{}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["error"] == "" {
		t.Error("expected error message in response")
	}
}

func TestAuthHandler_Register_ValidationFails_InvalidEmail(t *testing.T) {
	handler, _ := setupAuthHandler()

	reqBody := dto.RegisterRequest{
		Username: "validuser",
		Password: "P@ssj0rd121",
		Email:    "bukan-email-valid",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestAuthHandler_Register_DuplicateUsername(t *testing.T) {
	handler, _ := setupAuthHandler()

	reqBody := dto.RegisterRequest{
		Username: "sameuser",
		Password: "P@ssj0rd121",
		Email:    "first@example.com",
	}
	body, _ := json.Marshal(reqBody)

	// Register pertama — harus sukses
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Register(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("first register failed: %d %s", w.Code, w.Body.String())
	}

	// Register kedua dengan username sama — harus gagal
	reqBody2 := dto.RegisterRequest{
		Username: "sameuser",
		Password: "P@ssj0rd121",
		Email:    "second@example.com",
	}
	body2, _ := json.Marshal(reqBody2)
	req2 := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	handler.Register(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for duplicate username, got %d: %s", w2.Code, w2.Body.String())
	}
}
