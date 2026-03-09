package services

import (
	"encoding/json"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/testhelpers"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTokenService_GenerateTokenPair_Success(t *testing.T) {
	// Arrange
	redis := testhelpers.NewMockRedisClient()
	service := NewTokenService(redis, "test-secret", false, "localhost")

	// Act
	tokens, err := service.GenerateTokenPair("1", "testuser", "admin")

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if tokens.AccessToken == "" {
		t.Error("expected access token to be generated")
	}

	if tokens.RefreshToken == "" {
		t.Error("expected refresh token to be generated")
	}

	if tokens.ExpiresAt.IsZero() {
		t.Error("expected expires_at to be set")
	}

	// Verify refresh token is stored in Redis
	key := "refresh_token:" + tokens.RefreshToken
	exists, err := redis.Exists(key)
	if err != nil {
		t.Fatalf("error checking Redis: %v", err)
	}
	if !exists {
		t.Error("expected refresh token to be stored in Redis")
	}

	// Verify token data in Redis contains role
	data, err := redis.Get(key)
	if err != nil {
		t.Fatalf("error getting token data: %v", err)
	}

	var tokenData models.RefreshTokenData
	if err := json.Unmarshal([]byte(data), &tokenData); err != nil {
		t.Fatalf("error unmarshaling token data: %v", err)
	}

	if tokenData.UserID != "1" {
		t.Errorf("expected userID '1', got '%s'", tokenData.UserID)
	}

	if tokenData.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", tokenData.Username)
	}

	if tokenData.Role != "admin" {
		t.Errorf("expected role 'admin', got '%s'", tokenData.Role)
	}
}

func TestTokenService_RefreshAccessToken_Success(t *testing.T) {
	// Arrange
	redis := testhelpers.NewMockRedisClient()
	service := NewTokenService(redis, "test-secret", false, "localhost")

	// Generate initial token pair
	initialTokens, err := service.GenerateTokenPair("1", "testuser", "admin")
	if err != nil {
		t.Fatalf("failed to generate initial tokens: %v", err)
	}

	// Wait lebih lama untuk memastikan timestamp berbeda (1 detik)
	// JWT exp claim rounded ke seconds, jadi perlu delay untuk token berbeda
	time.Sleep(1 * time.Second)

	// Act
	newTokens, err := service.RefreshAccessToken(initialTokens.RefreshToken)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if newTokens == nil {
		t.Fatal("expected new tokens to be returned")
	}

	if newTokens.AccessToken == "" {
		t.Error("expected new access token to be generated")
	}

	if newTokens.RefreshToken == "" {
		t.Error("expected new refresh token to be generated")
	}

	// ACCESS TOKEN: Bisa sama atau beda tergantung timing
	// JWT tokens bisa sama jika dibuat di detik yang sama karena `exp` claim rounded ke seconds
	if newTokens.AccessToken == initialTokens.AccessToken {
		t.Log("Note: Access tokens are the same - this can happen if generated in same second (JWT exp is in seconds)")
	}

	// REFRESH TOKEN ROTATION: Check apakah token lama di-revoke
	oldKey := "refresh_token:" + initialTokens.RefreshToken
	exists, _ := redis.Exists(oldKey)
	if exists {
		t.Error("expected old refresh token to be revoked (token rotation)")
	}

	// Check token baru bisa digunakan
	if newTokens.RefreshToken == initialTokens.RefreshToken {
		t.Log("Note: Refresh token not rotated - same token reused")
	} else {
		// Verify new refresh token is usable
		_, err = service.RefreshAccessToken(newTokens.RefreshToken)
		if err != nil {
			t.Error("expected new refresh token to be usable")
		}
	}
}

func TestTokenService_RefreshAccessToken_InvalidToken(t *testing.T) {
	// Arrange
	redis := testhelpers.NewMockRedisClient()
	service := NewTokenService(redis, "test-secret", false, "localhost")

	// Act
	_, err := service.RefreshAccessToken("invalid-token")

	// Assert
	if err == nil {
		t.Fatal("expected error for invalid refresh token")
	}

	if err.Error() != "invalid or expired refresh token" {
		t.Errorf("expected 'invalid or expired refresh token', got '%s'", err.Error())
	}
}

func TestTokenService_RevokeRefreshToken_Success(t *testing.T) {
	// Arrange
	redis := testhelpers.NewMockRedisClient()
	service := NewTokenService(redis, "test-secret", false, "localhost")

	// Generate token pair
	tokens, _ := service.GenerateTokenPair("1", "testuser", "admin")

	// Act
	err := service.RevokeRefreshToken(tokens.RefreshToken)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify token is removed from Redis
	key := "refresh_token:" + tokens.RefreshToken
	exists, _ := redis.Exists(key)
	if exists {
		t.Error("expected refresh token to be removed from Redis")
	}

	// Try to use revoked token
	_, err = service.RefreshAccessToken(tokens.RefreshToken)
	if err == nil {
		t.Error("expected error when using revoked token")
	}
}

// ============================================================
// TestTokenService_SetAuthCookies
// ============================================================

func TestTokenService_SetAuthCookies_SetsBothCookies(t *testing.T) {
	redis := testhelpers.NewMockRedisClient()
	service := NewTokenService(redis, "test-secret", false, "localhost")

	tokens := &models.TokenPair{
		AccessToken:  "access-token-value",
		RefreshToken: "refresh-token-value",
		ExpiresAt:    time.Now().Add(15 * time.Minute),
	}

	w := httptest.NewRecorder()
	service.SetAuthCookies(w, tokens)

	resp := w.Result()
	cookieMap := make(map[string]*http.Cookie)
	for _, c := range resp.Cookies() {
		cookieMap[c.Name] = c
	}

	// access_token cookie
	ac, ok := cookieMap["access_token"]
	if !ok {
		t.Fatal("access_token cookie tidak ditemukan")
	}
	if ac.Value != "access-token-value" {
		t.Errorf("expected value 'access-token-value', got '%s'", ac.Value)
	}
	if ac.Path != "/" {
		t.Errorf("expected path '/', got '%s'", ac.Path)
	}
	if ac.MaxAge != 15*60 {
		t.Errorf("expected MaxAge %d, got %d", 15*60, ac.MaxAge)
	}
	if !ac.HttpOnly {
		t.Error("expected HttpOnly true")
	}
	if ac.Secure {
		t.Error("expected Secure false di non-production")
	}

	// refresh_token cookie
	rc, ok := cookieMap["refresh_token"]
	if !ok {
		t.Fatal("refresh_token cookie tidak ditemukan")
	}
	if rc.Value != "refresh-token-value" {
		t.Errorf("expected value 'refresh-token-value', got '%s'", rc.Value)
	}
	if rc.Path != "/api/refresh" {
		t.Errorf("expected path '/api/refresh', got '%s'", rc.Path)
	}
	if rc.MaxAge != 7*24*60*60 {
		t.Errorf("expected MaxAge %d, got %d", 7*24*60*60, rc.MaxAge)
	}
	if !rc.HttpOnly {
		t.Error("expected HttpOnly true")
	}
}

func TestTokenService_SetAuthCookies_SecureInProduction(t *testing.T) {
	redis := testhelpers.NewMockRedisClient()
	service := NewTokenService(redis, "test-secret", true, "example.com")

	tokens := &models.TokenPair{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresAt:    time.Now().Add(15 * time.Minute),
	}

	w := httptest.NewRecorder()
	service.SetAuthCookies(w, tokens)

	for _, c := range w.Result().Cookies() {
		if !c.Secure {
			t.Errorf("cookie '%s' harus Secure di production", c.Name)
		}
	}
}

func TestTokenService_SetAuthCookies_DomainDiSet(t *testing.T) {
	redis := testhelpers.NewMockRedisClient()
	service := NewTokenService(redis, "test-secret", false, "myapp.example.com")

	tokens := &models.TokenPair{
		AccessToken:  "access",
		RefreshToken: "refresh",
		ExpiresAt:    time.Now().Add(15 * time.Minute),
	}

	w := httptest.NewRecorder()
	service.SetAuthCookies(w, tokens)

	for _, c := range w.Result().Cookies() {
		if c.Domain != "myapp.example.com" {
			t.Errorf("cookie '%s': expected domain 'myapp.example.com', got '%s'", c.Name, c.Domain)
		}
	}
}

// ============================================================
// TestTokenService_GetAccessTokenFromCookie
// ============================================================

func TestTokenService_GetAccessTokenFromCookie_Success(t *testing.T) {
	service := NewTokenService(testhelpers.NewMockRedisClient(), "secret", false, "localhost")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "my-access-token"})

	token, err := service.GetAccessTokenFromCookie(req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token != "my-access-token" {
		t.Errorf("expected 'my-access-token', got '%s'", token)
	}
}

func TestTokenService_GetAccessTokenFromCookie_NoCookie(t *testing.T) {
	service := NewTokenService(testhelpers.NewMockRedisClient(), "secret", false, "localhost")

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	token, err := service.GetAccessTokenFromCookie(req)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "access token cookie not found" {
		t.Errorf("expected 'access token cookie not found', got '%s'", err.Error())
	}
	if token != "" {
		t.Errorf("expected empty token, got '%s'", token)
	}
}

func TestTokenService_GetAccessTokenFromCookie_WrongCookieName(t *testing.T) {
	service := NewTokenService(testhelpers.NewMockRedisClient(), "secret", false, "localhost")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "some-value"})

	_, err := service.GetAccessTokenFromCookie(req)

	if err == nil {
		t.Error("expected error saat cookie access_token tidak ada")
	}
}

// ============================================================
// TestTokenService_GetRefreshTokenFromCookie
// ============================================================

func TestTokenService_GetRefreshTokenFromCookie_Success(t *testing.T) {
	service := NewTokenService(testhelpers.NewMockRedisClient(), "secret", false, "localhost")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "my-refresh-token"})

	token, err := service.GetRefreshTokenFromCookie(req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token != "my-refresh-token" {
		t.Errorf("expected 'my-refresh-token', got '%s'", token)
	}
}

func TestTokenService_GetRefreshTokenFromCookie_NoCookie(t *testing.T) {
	service := NewTokenService(testhelpers.NewMockRedisClient(), "secret", false, "localhost")

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	token, err := service.GetRefreshTokenFromCookie(req)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "refresh token cookie not found" {
		t.Errorf("expected 'refresh token cookie not found', got '%s'", err.Error())
	}
	if token != "" {
		t.Errorf("expected empty token, got '%s'", token)
	}
}

// ============================================================
// TestTokenService_ClearAuthCookies
// ============================================================

func TestTokenService_ClearAuthCookies_MaxAgeMinusOne(t *testing.T) {
	service := NewTokenService(testhelpers.NewMockRedisClient(), "secret", false, "localhost")

	w := httptest.NewRecorder()
	service.ClearAuthCookies(w)

	cookieMap := make(map[string]*http.Cookie)
	for _, c := range w.Result().Cookies() {
		cookieMap[c.Name] = c
	}

	ac, ok := cookieMap["access_token"]
	if !ok {
		t.Fatal("access_token cookie tidak ditemukan")
	}
	if ac.MaxAge != -1 {
		t.Errorf("expected MaxAge -1, got %d", ac.MaxAge)
	}
	if ac.Value != "" {
		t.Errorf("expected Value kosong, got '%s'", ac.Value)
	}

	rc, ok := cookieMap["refresh_token"]
	if !ok {
		t.Fatal("refresh_token cookie tidak ditemukan")
	}
	if rc.MaxAge != -1 {
		t.Errorf("expected MaxAge -1, got %d", rc.MaxAge)
	}
	if rc.Value != "" {
		t.Errorf("expected Value kosong, got '%s'", rc.Value)
	}
}

func TestTokenService_ClearAuthCookies_PathBenar(t *testing.T) {
	service := NewTokenService(testhelpers.NewMockRedisClient(), "secret", false, "localhost")

	w := httptest.NewRecorder()
	service.ClearAuthCookies(w)

	cookieMap := make(map[string]*http.Cookie)
	for _, c := range w.Result().Cookies() {
		cookieMap[c.Name] = c
	}

	if cookieMap["access_token"].Path != "/" {
		t.Errorf("access_token path: expected '/', got '%s'", cookieMap["access_token"].Path)
	}
	if cookieMap["refresh_token"].Path != "/api/refresh" {
		t.Errorf("refresh_token path: expected '/api/refresh', got '%s'", cookieMap["refresh_token"].Path)
	}
}

// ============================================================
// TestTokenService_RevokeAllUserTokens
// ============================================================

func TestTokenService_RevokeAllUserTokens_HapusSemuaTokenUser(t *testing.T) {
	redis := testhelpers.NewMockRedisClient()
	service := NewTokenService(redis, "test-secret", false, "localhost")

	// User 1 login dari 3 device
	t1, _ := service.GenerateTokenPair("user-1", "alice", "admin")
	t2, _ := service.GenerateTokenPair("user-1", "alice", "admin")
	t3, _ := service.GenerateTokenPair("user-1", "alice", "admin")

	// User 2 punya token tersendiri
	t4, _ := service.GenerateTokenPair("user-2", "bob", "user")

	err := service.RevokeAllUserTokens("user-1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Semua token user-1 harus terhapus
	for i, tok := range []*models.TokenPair{t1, t2, t3} {
		exists, _ := redis.Exists("refresh_token:" + tok.RefreshToken)
		if exists {
			t.Errorf("token user-1 ke-%d seharusnya sudah dihapus", i+1)
		}
	}

	// Token user-2 harus tetap ada
	exists, _ := redis.Exists("refresh_token:" + t4.RefreshToken)
	if !exists {
		t.Error("token user-2 tidak boleh ikut terhapus")
	}
}

func TestTokenService_RevokeAllUserTokens_UserTanpaToken(t *testing.T) {
	redis := testhelpers.NewMockRedisClient()
	service := NewTokenService(redis, "test-secret", false, "localhost")

	// Revoke user yang tidak punya token → tidak error
	err := service.RevokeAllUserTokens("user-tidak-ada")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestTokenService_RevokeAllUserTokens_TokenTidakBisaDigunakan(t *testing.T) {
	redis := testhelpers.NewMockRedisClient()
	service := NewTokenService(redis, "test-secret", false, "localhost")

	tokens, _ := service.GenerateTokenPair("user-1", "alice", "admin")

	_ = service.RevokeAllUserTokens("user-1")

	_, err := service.RefreshAccessToken(tokens.RefreshToken)
	if err == nil {
		t.Error("token yang sudah direvoke seharusnya tidak bisa digunakan")
	}
}

// ============================================================
// TestTokenService_ValidateAndRefreshIfNeeded
// ============================================================

func TestTokenService_ValidateAndRefreshIfNeeded_AccessTokenValid(t *testing.T) {
	redis := testhelpers.NewMockRedisClient()
	service := NewTokenService(redis, "test-secret", false, "localhost")

	tokens, _ := service.GenerateTokenPair("user-1", "alice", "admin")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: tokens.AccessToken})
	w := httptest.NewRecorder()

	claims, err := service.ValidateAndRefreshIfNeeded(w, req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if claims.UserID != "user-1" {
		t.Errorf("expected userID 'user-1', got '%s'", claims.UserID)
	}
	if claims.Username != "alice" {
		t.Errorf("expected username 'alice', got '%s'", claims.Username)
	}
	if claims.Role != "admin" {
		t.Errorf("expected role 'admin', got '%s'", claims.Role)
	}
}

func TestTokenService_ValidateAndRefreshIfNeeded_TanpaCookie_ReturnsError(t *testing.T) {
	redis := testhelpers.NewMockRedisClient()
	service := NewTokenService(redis, "test-secret", false, "localhost")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	claims, err := service.ValidateAndRefreshIfNeeded(w, req)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "authentication required" {
		t.Errorf("expected 'authentication required', got '%s'", err.Error())
	}
	if claims != nil {
		t.Error("expected nil claims")
	}
}

func TestTokenService_ValidateAndRefreshIfNeeded_AccessInvalid_RefreshValid_IssuesNewTokens(t *testing.T) {
	redis := testhelpers.NewMockRedisClient()
	service := NewTokenService(redis, "test-secret", false, "localhost")

	tokens, _ := service.GenerateTokenPair("user-1", "alice", "admin")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "access-token-tidak-valid"})
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: tokens.RefreshToken})
	w := httptest.NewRecorder()

	claims, err := service.ValidateAndRefreshIfNeeded(w, req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if claims.UserID != "user-1" {
		t.Errorf("expected userID 'user-1', got '%s'", claims.UserID)
	}

	// Cookie baru harus di-set di response
	cookieMap := make(map[string]*http.Cookie)
	for _, c := range w.Result().Cookies() {
		cookieMap[c.Name] = c
	}
	if _, ok := cookieMap["access_token"]; !ok {
		t.Error("access_token cookie baru harus di-set")
	}
	if _, ok := cookieMap["refresh_token"]; !ok {
		t.Error("refresh_token cookie baru harus di-set")
	}
}

func TestTokenService_ValidateAndRefreshIfNeeded_BothTokensInvalid_ReturnsError(t *testing.T) {
	redis := testhelpers.NewMockRedisClient()
	service := NewTokenService(redis, "test-secret", false, "localhost")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "access-tidak-valid"})
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "refresh-tidak-valid"})
	w := httptest.NewRecorder()

	claims, err := service.ValidateAndRefreshIfNeeded(w, req)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if claims != nil {
		t.Error("expected nil claims")
	}
}
