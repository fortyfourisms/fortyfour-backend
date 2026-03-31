package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/utils"
)

// =========================
// MOCK REDIS
// =========================

type mockTokenRedis struct {
	store      map[string]string
	failSet    bool
	failGet    bool
	failDelete bool
	failScan   bool
}

func newMockTokenRedis() *mockTokenRedis {
	return &mockTokenRedis{store: make(map[string]string)}
}

func (m *mockTokenRedis) Set(key string, value interface{}, ttl time.Duration) error {
	if m.failSet {
		return errors.New("redis set error")
	}
	str, ok := value.(string)
	if !ok {
		return errors.New("value must be string")
	}
	m.store[key] = str
	return nil
}

func (m *mockTokenRedis) Get(key string) (string, error) {
	if m.failGet {
		return "", errors.New("redis get error")
	}
	val, ok := m.store[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return val, nil
}

func (m *mockTokenRedis) Delete(key string) error {
	if m.failDelete {
		return errors.New("redis delete error")
	}
	delete(m.store, key)
	return nil
}

func (m *mockTokenRedis) Exists(key string) (bool, error) {
	_, ok := m.store[key]
	return ok, nil
}

func (m *mockTokenRedis) Close() error { return nil }

func (m *mockTokenRedis) Scan(pattern string) ([]string, error) {
	if m.failScan {
		return nil, errors.New("redis scan error")
	}
	prefix := strings.TrimSuffix(pattern, "*")
	keys := make([]string, 0, len(m.store))
	for k := range m.store {
		if strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

// =========================
// HELPERS
// =========================

const tokenTestJWTSecret = "test-secret-key-for-unit-testing"

func newTestTokenSvc(r *mockTokenRedis) *TokenService {
	return NewTokenService(r, tokenTestJWTSecret, false, "localhost")
}

func newTestTokenSvcProd(r *mockTokenRedis) *TokenService {
	return NewTokenService(r, tokenTestJWTSecret, true, "example.com")
}

// generateValidPair membuat token pair yang valid dan tersimpan di redis mock.
func generateValidPair(t *testing.T, svc *TokenService, userID, username, role, company string) *models.TokenPair {
	t.Helper()
	pair, err := svc.GenerateTokenPair(userID, username, role, company)
	if err != nil {
		t.Fatalf("generateValidPair: %v", err)
	}
	return pair
}

// =========================
// GenerateTokenPair
// =========================

func TestTokenService_GenerateTokenPair_ReturnsNonEmptyTokens(t *testing.T) {
	svc := newTestTokenSvc(newMockTokenRedis())

	pair, err := svc.GenerateTokenPair("user-1", "alice", "admin", "cmp-1")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if pair.AccessToken == "" {
		t.Error("access token should not be empty")
	}
	if pair.RefreshToken == "" {
		t.Error("refresh token should not be empty")
	}
	if pair.ExpiresAt.IsZero() {
		t.Error("ExpiresAt should not be zero")
	}
}

func TestTokenService_GenerateTokenPair_StoresDataInRedis(t *testing.T) {
	redis := newMockTokenRedis()
	svc := newTestTokenSvc(redis)

	pair, err := svc.GenerateTokenPair("user-1", "alice", "admin", "cmp-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	key := fmt.Sprintf("refresh_token:%s", pair.RefreshToken)
	raw, ok := redis.store[key]
	if !ok {
		t.Fatal("refresh token data not found in redis")
	}

	var data models.RefreshTokenData
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		t.Fatalf("failed to parse stored token data: %v", err)
	}

	fields := []struct{ got, want, name string }{
		{data.UserID, "user-1", "UserID"},
		{data.Username, "alice", "Username"},
		{data.Role, "admin", "Role"},
		{data.IDPerusahaan, "cmp-1", "IDPerusahaan"},
	}
	for _, f := range fields {
		if f.got != f.want {
			t.Errorf("%s: want '%s', got '%s'", f.name, f.want, f.got)
		}
	}
}

func TestTokenService_GenerateTokenPair_FailsWhenRedisDown(t *testing.T) {
	redis := newMockTokenRedis()
	redis.failSet = true
	svc := newTestTokenSvc(redis)

	_, err := svc.GenerateTokenPair("user-1", "alice", "admin", "cmp-1")
	if err == nil {
		t.Error("expected error when redis Set fails")
	}
}

func TestTokenService_GenerateTokenPair_AccessTokenContainsCorrectClaims(t *testing.T) {
	svc := newTestTokenSvc(newMockTokenRedis())

	pair, err := svc.GenerateTokenPair("user-42", "bob", "viewer", "cmp-99")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	claims, err := utils.ValidateAccessToken(pair.AccessToken, tokenTestJWTSecret)
	if err != nil {
		t.Fatalf("access token failed validation: %v", err)
	}

	if claims.UserID != "user-42" {
		t.Errorf("UserID: want 'user-42', got '%s'", claims.UserID)
	}
	if claims.Username != "bob" {
		t.Errorf("Username: want 'bob', got '%s'", claims.Username)
	}
	if claims.Role != "viewer" {
		t.Errorf("Role: want 'viewer', got '%s'", claims.Role)
	}
	if claims.IDPerusahaan != "cmp-99" {
		t.Errorf("IDPerusahaan: want 'cmp-99', got '%s'", claims.IDPerusahaan)
	}
}

func TestTokenService_GenerateTokenPair_UniqueRefreshTokensEachCall(t *testing.T) {
	svc := newTestTokenSvc(newMockTokenRedis())

	p1, _ := svc.GenerateTokenPair("user-1", "alice", "admin", "cmp-1")
	p2, _ := svc.GenerateTokenPair("user-1", "alice", "admin", "cmp-1")

	if p1.RefreshToken == p2.RefreshToken {
		t.Error("expected unique refresh tokens on each call")
	}
}

// =========================
// SetAuthCookies
// =========================

func TestTokenService_SetAuthCookies_SetsBothCookies(t *testing.T) {
	svc := newTestTokenSvc(newMockTokenRedis())
	w := httptest.NewRecorder()

	pair := &models.TokenPair{
		AccessToken:  "at-value",
		RefreshToken: "rt-value",
		ExpiresAt:    time.Now().Add(15 * time.Minute),
	}
	svc.SetAuthCookies(w, pair)

	cookieMap := make(map[string]*http.Cookie)
	for _, c := range w.Result().Cookies() {
		cookieMap[c.Name] = c
	}

	at, ok := cookieMap["access_token"]
	if !ok {
		t.Fatal("access_token cookie not set")
	}
	if at.Value != "at-value" {
		t.Errorf("access_token value: want 'at-value', got '%s'", at.Value)
	}
	if !at.HttpOnly {
		t.Error("access_token must be HttpOnly")
	}
	if at.Secure {
		t.Error("access_token should not be Secure in non-production")
	}

	rt, ok := cookieMap["refresh_token"]
	if !ok {
		t.Fatal("refresh_token cookie not set")
	}
	if rt.Value != "rt-value" {
		t.Errorf("refresh_token value: want 'rt-value', got '%s'", rt.Value)
	}
	if rt.Path != "/api/refresh" {
		t.Errorf("refresh_token path: want '/api/refresh', got '%s'", rt.Path)
	}
}

func TestTokenService_SetAuthCookies_SecureFlagInProduction(t *testing.T) {
	svc := newTestTokenSvcProd(newMockTokenRedis())
	w := httptest.NewRecorder()

	pair := &models.TokenPair{AccessToken: "at", RefreshToken: "rt", ExpiresAt: time.Now()}
	svc.SetAuthCookies(w, pair)

	for _, c := range w.Result().Cookies() {
		if !c.Secure {
			t.Errorf("cookie '%s' should be Secure in production mode", c.Name)
		}
	}
}

func TestTokenService_SetAuthCookies_MaxAgeValues(t *testing.T) {
	svc := newTestTokenSvc(newMockTokenRedis())
	w := httptest.NewRecorder()

	svc.SetAuthCookies(w, &models.TokenPair{AccessToken: "at", RefreshToken: "rt"})

	for _, c := range w.Result().Cookies() {
		switch c.Name {
		case "access_token":
			if c.MaxAge != 15*60 {
				t.Errorf("access_token MaxAge: want %d, got %d", 15*60, c.MaxAge)
			}
		case "refresh_token":
			if c.MaxAge != 7*24*60*60 {
				t.Errorf("refresh_token MaxAge: want %d, got %d", 7*24*60*60, c.MaxAge)
			}
		}
	}
}

// =========================
// ClearAuthCookies
// =========================

func TestTokenService_ClearAuthCookies_SetsNegativeMaxAge(t *testing.T) {
	svc := newTestTokenSvc(newMockTokenRedis())
	w := httptest.NewRecorder()

	svc.ClearAuthCookies(w)

	cookieMap := make(map[string]*http.Cookie)
	for _, c := range w.Result().Cookies() {
		cookieMap[c.Name] = c
	}

	for _, name := range []string{"access_token", "refresh_token"} {
		c, ok := cookieMap[name]
		if !ok {
			t.Fatalf("cookie '%s' not found after ClearAuthCookies", name)
		}
		if c.MaxAge != -1 {
			t.Errorf("cookie '%s' MaxAge: want -1, got %d", name, c.MaxAge)
		}
		if c.Value != "" {
			t.Errorf("cookie '%s' should be empty after clear, got '%s'", name, c.Value)
		}
	}
}

// =========================
// GetAccessTokenFromCookie
// =========================

func TestTokenService_GetAccessTokenFromCookie_Found(t *testing.T) {
	svc := newTestTokenSvc(newMockTokenRedis())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "tok-123"})

	got, err := svc.GetAccessTokenFromCookie(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "tok-123" {
		t.Errorf("want 'tok-123', got '%s'", got)
	}
}

func TestTokenService_GetAccessTokenFromCookie_NotFound(t *testing.T) {
	svc := newTestTokenSvc(newMockTokenRedis())
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	_, err := svc.GetAccessTokenFromCookie(req)
	if err == nil {
		t.Error("expected error when access_token cookie is absent")
	}
}

// =========================
// GetRefreshTokenFromCookie
// =========================

func TestTokenService_GetRefreshTokenFromCookie_Found(t *testing.T) {
	svc := newTestTokenSvc(newMockTokenRedis())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "ref-456"})

	got, err := svc.GetRefreshTokenFromCookie(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "ref-456" {
		t.Errorf("want 'ref-456', got '%s'", got)
	}
}

func TestTokenService_GetRefreshTokenFromCookie_NotFound(t *testing.T) {
	svc := newTestTokenSvc(newMockTokenRedis())
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	_, err := svc.GetRefreshTokenFromCookie(req)
	if err == nil {
		t.Error("expected error when refresh_token cookie is absent")
	}
}

// =========================
// RefreshAccessToken
// =========================

func TestTokenService_RefreshAccessToken_Success(t *testing.T) {
	redis := newMockTokenRedis()
	svc := newTestTokenSvc(redis)

	oldPair := generateValidPair(t, svc, "user-1", "alice", "admin", "cmp-1")

	newPair, err := svc.RefreshAccessToken(oldPair.RefreshToken)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if newPair.AccessToken == "" {
		t.Error("new access token should not be empty")
	}
}

func TestTokenService_RefreshAccessToken_RotatesToken(t *testing.T) {
	redis := newMockTokenRedis()
	svc := newTestTokenSvc(redis)

	oldPair := generateValidPair(t, svc, "user-1", "alice", "admin", "cmp-1")

	newPair, err := svc.RefreshAccessToken(oldPair.RefreshToken)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if newPair.RefreshToken == oldPair.RefreshToken {
		t.Error("refresh token should be rotated — new token must differ from old")
	}
}

func TestTokenService_RefreshAccessToken_OldTokenIsRevoked(t *testing.T) {
	redis := newMockTokenRedis()
	svc := newTestTokenSvc(redis)

	oldPair := generateValidPair(t, svc, "user-1", "alice", "admin", "cmp-1")

	// Pakai sekali
	if _, err := svc.RefreshAccessToken(oldPair.RefreshToken); err != nil {
		t.Fatalf("first refresh failed: %v", err)
	}

	// Replay attack — harus ditolak
	_, err := svc.RefreshAccessToken(oldPair.RefreshToken)
	if err == nil {
		t.Error("expected error on replay: old refresh token should be revoked after rotation")
	}
}

func TestTokenService_RefreshAccessToken_InvalidToken(t *testing.T) {
	svc := newTestTokenSvc(newMockTokenRedis())

	_, err := svc.RefreshAccessToken("this-token-does-not-exist")
	if err == nil {
		t.Error("expected error for unknown refresh token")
	}
}

func TestTokenService_RefreshAccessToken_CorruptData(t *testing.T) {
	redis := newMockTokenRedis()
	svc := newTestTokenSvc(redis)

	redis.store["refresh_token:corrupt"] = "{not valid json"

	_, err := svc.RefreshAccessToken("corrupt")
	if err == nil {
		t.Error("expected error when redis contains malformed token data")
	}
}

func TestTokenService_RefreshAccessToken_PreservesUserClaims(t *testing.T) {
	redis := newMockTokenRedis()
	svc := newTestTokenSvc(redis)

	oldPair := generateValidPair(t, svc, "user-77", "charlie", "manager", "cmp-77")

	newPair, err := svc.RefreshAccessToken(oldPair.RefreshToken)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	claims, err := utils.ValidateAccessToken(newPair.AccessToken, tokenTestJWTSecret)
	if err != nil {
		t.Fatalf("new access token validation failed: %v", err)
	}
	if claims.UserID != "user-77" {
		t.Errorf("UserID: want 'user-77', got '%s'", claims.UserID)
	}
	if claims.Role != "manager" {
		t.Errorf("Role: want 'manager', got '%s'", claims.Role)
	}
}

// =========================
// RevokeRefreshToken
// =========================

func TestTokenService_RevokeRefreshToken_DeletesFromRedis(t *testing.T) {
	redis := newMockTokenRedis()
	svc := newTestTokenSvc(redis)

	pair := generateValidPair(t, svc, "user-1", "alice", "admin", "cmp-1")
	key := fmt.Sprintf("refresh_token:%s", pair.RefreshToken)

	if _, ok := redis.store[key]; !ok {
		t.Fatal("token should exist before revoking")
	}

	if err := svc.RevokeRefreshToken(pair.RefreshToken); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := redis.store[key]; ok {
		t.Error("token should be deleted after RevokeRefreshToken")
	}
}

func TestTokenService_RevokeRefreshToken_RedisDeleteFails(t *testing.T) {
	redis := newMockTokenRedis()
	redis.failDelete = true
	svc := newTestTokenSvc(redis)

	err := svc.RevokeRefreshToken("any-token")
	if err == nil {
		t.Error("expected error when redis Delete fails")
	}
}

// =========================
// RevokeAllUserTokens
// =========================

func TestTokenService_RevokeAllUserTokens_RevokesAllTokensForUser(t *testing.T) {
	redis := newMockTokenRedis()
	svc := newTestTokenSvc(redis)

	// 2 session untuk user-1
	generateValidPair(t, svc, "user-1", "alice", "admin", "cmp-1")
	generateValidPair(t, svc, "user-1", "alice", "admin", "cmp-1")

	if err := svc.RevokeAllUserTokens("user-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, v := range redis.store {
		var data models.RefreshTokenData
		if err := json.Unmarshal([]byte(v), &data); err == nil {
			if data.UserID == "user-1" {
				t.Error("found user-1 token still in redis after RevokeAllUserTokens")
			}
		}
	}
}

func TestTokenService_RevokeAllUserTokens_OnlyRevokesTargetUser(t *testing.T) {
	redis := newMockTokenRedis()
	svc := newTestTokenSvc(redis)

	generateValidPair(t, svc, "user-1", "alice", "admin", "cmp-1")
	user2Pair := generateValidPair(t, svc, "user-2", "bob", "viewer", "cmp-2")

	if err := svc.RevokeAllUserTokens("user-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	key := fmt.Sprintf("refresh_token:%s", user2Pair.RefreshToken)
	if _, ok := redis.store[key]; !ok {
		t.Error("user-2 token should NOT be revoked when revoking user-1")
	}
}

func TestTokenService_RevokeAllUserTokens_ScanError(t *testing.T) {
	redis := newMockTokenRedis()
	redis.failScan = true
	svc := newTestTokenSvc(redis)

	err := svc.RevokeAllUserTokens("user-1")
	if err == nil {
		t.Error("expected error when redis Scan fails")
	}
}

func TestTokenService_RevokeAllUserTokens_NoOpForUserWithNoTokens(t *testing.T) {
	svc := newTestTokenSvc(newMockTokenRedis())

	if err := svc.RevokeAllUserTokens("user-with-no-tokens"); err != nil {
		t.Errorf("expected no error for user with no tokens, got: %v", err)
	}
}

// =========================
// ValidateAndRefreshIfNeeded
// =========================

func TestTokenService_ValidateAndRefreshIfNeeded_ValidAccessToken(t *testing.T) {
	redis := newMockTokenRedis()
	svc := newTestTokenSvc(redis)

	pair := generateValidPair(t, svc, "user-1", "alice", "admin", "cmp-1")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: pair.AccessToken})
	w := httptest.NewRecorder()

	claims, err := svc.ValidateAndRefreshIfNeeded(w, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claims.UserID != "user-1" {
		t.Errorf("UserID: want 'user-1', got '%s'", claims.UserID)
	}
}

func TestTokenService_ValidateAndRefreshIfNeeded_InvalidAccessUsesRefresh(t *testing.T) {
	redis := newMockTokenRedis()
	svc := newTestTokenSvc(redis)

	validPair := generateValidPair(t, svc, "user-1", "alice", "admin", "cmp-1")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "invalid-or-expired"})
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: validPair.RefreshToken})
	w := httptest.NewRecorder()

	claims, err := svc.ValidateAndRefreshIfNeeded(w, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claims.UserID != "user-1" {
		t.Errorf("UserID: want 'user-1', got '%s'", claims.UserID)
	}
	if len(w.Result().Cookies()) == 0 {
		t.Error("expected new auth cookies to be set after refresh")
	}
}

func TestTokenService_ValidateAndRefreshIfNeeded_NoCookies(t *testing.T) {
	svc := newTestTokenSvc(newMockTokenRedis())

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	_, err := svc.ValidateAndRefreshIfNeeded(w, req)
	if err == nil {
		t.Error("expected error when no cookies are present")
	}
}

func TestTokenService_ValidateAndRefreshIfNeeded_BothTokensInvalid(t *testing.T) {
	svc := newTestTokenSvc(newMockTokenRedis())

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "bad-access"})
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "bad-refresh"})
	w := httptest.NewRecorder()

	_, err := svc.ValidateAndRefreshIfNeeded(w, req)
	if err == nil {
		t.Error("expected error when both access and refresh tokens are invalid")
	}
}
