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
	"time"

	"github.com/pquerna/otp/totp"
)

func setupAuthHandler() (*AuthHandler, *testhelpers.MockRedisClient) {
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := services.NewTokenService(redis, "test-secret", false, "localhost")
	authService := services.NewAuthService(userRepo, tokenService, services.NewNotificationService(redis))
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
