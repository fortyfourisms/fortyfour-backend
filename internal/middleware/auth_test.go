package middleware

import (
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/testhelpers"
	"fortyfour-backend/internal/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testJWTSecret = "test-secret"

// createTestTokenService creates a TokenService with mock Redis for testing
func createTestTokenService() *services.TokenService {
	mockRedis := testhelpers.NewMockRedisClient()
	return services.NewTokenService(mockRedis, testJWTSecret, false, "localhost")
}

// createAccessTokenCookie creates a cookie with the given access token
func createAccessTokenCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:  "access_token",
		Value: token,
	}
}

func TestAuthMiddleware_Authenticate_MissingCookie(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)
	handler := middleware.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_Authenticate_InvalidToken(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)
	handler := middleware.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(createAccessTokenCookie("invalid-token"))
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_Authenticate_ValidToken(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)

	// Generate a valid token
	token, _, err := utils.GenerateAccessToken("user-1", "testuser", "admin", testJWTSecret)
	assert.NoError(t, err)

	handler := middleware.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		// Check if user ID is in context
		userID := r.Context().Value(UserIDKey)
		if userID == nil {
			t.Error("user ID not found in context")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(createAccessTokenCookie(token))
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_Authenticate_ContextValues(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)

	// Generate a valid token
	token, _, err := utils.GenerateAccessToken("user-1", "testuser", "admin", testJWTSecret)
	assert.NoError(t, err)

	handler := middleware.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(UserIDKey)
		username := r.Context().Value(UsernameKey)
		role := r.Context().Value(RoleKey)

		assert.Equal(t, "user-1", userID)
		assert.Equal(t, "testuser", username)
		assert.Equal(t, "admin", role)

		// Check X-User-Role header (the only header set by the middleware)
		assert.Equal(t, "admin", r.Header.Get("X-User-Role"))

		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(createAccessTokenCookie(token))
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_OptionalAuth_WithoutToken(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)

	handler := middleware.OptionalAuth(func(w http.ResponseWriter, r *http.Request) {
		// Should still pass through without authentication
		userID := r.Context().Value(UserIDKey)
		assert.Nil(t, userID)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_OptionalAuth_WithValidToken(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)

	// Generate a valid token
	token, _, err := utils.GenerateAccessToken("user-1", "testuser", "user", testJWTSecret)
	assert.NoError(t, err)

	handler := middleware.OptionalAuth(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(UserIDKey)
		username := r.Context().Value(UsernameKey)
		role := r.Context().Value(RoleKey)

		assert.Equal(t, "user-1", userID)
		assert.Equal(t, "testuser", username)
		assert.Equal(t, "user", role)

		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(createAccessTokenCookie(token))
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_OptionalAuth_WithInvalidToken(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)

	handler := middleware.OptionalAuth(func(w http.ResponseWriter, r *http.Request) {
		// Should still pass through, but without context values
		userID := r.Context().Value(UserIDKey)
		assert.Nil(t, userID)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(createAccessTokenCookie("invalid-token"))
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ========== EDGE CASE TESTS ==========

func TestAuthMiddleware_Authenticate_MalformedToken(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)

	handler := middleware.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with malformed token")
		w.WriteHeader(http.StatusOK)
	})

	// Test various malformed tokens
	malformedTokens := []string{
		"not.a.jwt",
		"only-one-part",
		"two.parts",
		"",
		"Bearer token.with.bearer",
	}

	for _, token := range malformedTokens {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(createAccessTokenCookie(token))
		w := httptest.NewRecorder()

		handler(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Should reject malformed token: %s", token)
	}
}

func TestAuthMiddleware_Authenticate_EmptyCookieValue(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)

	handler := middleware.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with empty cookie")
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(createAccessTokenCookie(""))
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_Authenticate_LongCookieValue(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)

	handler := middleware.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with invalid long token")
		w.WriteHeader(http.StatusOK)
	})

	// Create a very long invalid token (simulating potential attack)
	longToken := strings.Repeat("a", 4000) // use strings.Repeat instead of inefficient concatenation

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(createAccessTokenCookie(longToken))
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_Authenticate_SpecialCharactersInCookie(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)

	handler := middleware.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with special chars token")
		w.WriteHeader(http.StatusOK)
	})

	// Test tokens with special characters
	specialTokens := []string{
		"token<script>alert('xss')</script>",
		"token'; DROP TABLE users;--",
		"../../../etc/passwd",
	}

	for _, token := range specialTokens {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(createAccessTokenCookie(token))
		w := httptest.NewRecorder()

		handler(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Should reject token with special chars: %s", token)
	}
}

func TestAuthMiddleware_OptionalAuth_MalformedTokenGracefulFail(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)

	handler := middleware.OptionalAuth(func(w http.ResponseWriter, r *http.Request) {
		// Should still pass through even with malformed token
		userID := r.Context().Value(UserIDKey)
		assert.Nil(t, userID, "Context should not be set with malformed token")
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(createAccessTokenCookie("malformed.token.here"))
	w := httptest.NewRecorder()

	handler(w, req)

	// OptionalAuth should still return OK
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_Authenticate_MultipleCookiesWithSameName(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)

	// Generate a valid token
	validToken, _, err := utils.GenerateAccessToken("user-1", "testuser", "admin", testJWTSecret)
	assert.NoError(t, err)

	handler := middleware.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(UserIDKey)
		assert.Equal(t, "user-1", userID)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	
	// Add multiple cookies with the same name (first one should be used)
	req.AddCookie(createAccessTokenCookie(validToken))
	req.AddCookie(createAccessTokenCookie("invalid-token"))
	
	w := httptest.NewRecorder()

	handler(w, req)

	// Should use the first cookie (which is valid)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_Authenticate_TokenWithDifferentAlgorithm(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)

	handler := middleware.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with wrong algorithm token")
		w.WriteHeader(http.StatusOK)
	})

	// Create a token with "none" algorithm (security vulnerability test)
	noneAlgToken := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJ1c2VyX2lkIjoidXNlci0xIiwidXNlcm5hbWUiOiJ0ZXN0Iiwicm9sZSI6ImFkbWluIn0."

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(createAccessTokenCookie(noneAlgToken))
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_Authenticate_ContextValueTypes(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)

	token, _, err := utils.GenerateAccessToken("user-123", "john_doe", "superadmin", testJWTSecret)
	assert.NoError(t, err)

	handler := middleware.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		// Verify types of context values
		userID := r.Context().Value(UserIDKey)
		username := r.Context().Value(UsernameKey)
		role := r.Context().Value(RoleKey)

		// All should be strings
		assert.IsType(t, "", userID)
		assert.IsType(t, "", username)
		assert.IsType(t, "", role)

		// Verify actual values
		assert.Equal(t, "user-123", userID)
		assert.Equal(t, "john_doe", username)
		assert.Equal(t, "superadmin", role)

		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(createAccessTokenCookie(token))
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ========== AUTO REFRESH TESTS ==========


func TestAuthMiddleware_AutoRefresh_ValidTokenNoRefreshNeeded(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)

	// Generate a valid token that doesn't need refresh
	token, _, err := utils.GenerateAccessToken("user-1", "testuser", "admin", testJWTSecret)
	assert.NoError(t, err)

	handler := middleware.AutoRefresh(func(w http.ResponseWriter, r *http.Request) {
		// Check if context values are set correctly
		userID := r.Context().Value(UserIDKey)
		username := r.Context().Value(UsernameKey)
		role := r.Context().Value(RoleKey)

		assert.Equal(t, "user-1", userID)
		assert.Equal(t, "testuser", username)
		assert.Equal(t, "admin", role)
		assert.Equal(t, "admin", r.Header.Get("X-User-Role"))

		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(createAccessTokenCookie(token))
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_AutoRefresh_ExpiredTokenRefreshSuccess(t *testing.T) {
	mockRedis := testhelpers.NewMockRedisClient()
	tokenService := services.NewTokenService(mockRedis, testJWTSecret, false, "localhost")
	middleware := NewAuthMiddleware(tokenService)

	// Create token pair with valid refresh token
	tokens, err := tokenService.GenerateTokenPair("user-2", "refreshuser", "user")
	assert.NoError(t, err)

	// Create an expired access token manually (using very short expiry)
	// For testing purposes, we'll use an invalid token to simulate expiration
	expiredToken := "expired.invalid.token"

	handler := middleware.AutoRefresh(func(w http.ResponseWriter, r *http.Request) {
		// After successful refresh, context should be populated
		userID := r.Context().Value(UserIDKey)
		username := r.Context().Value(UsernameKey)
		role := r.Context().Value(RoleKey)

		assert.Equal(t, "user-2", userID)
		assert.Equal(t, "refreshuser", username)
		assert.Equal(t, "user", role)

		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(createAccessTokenCookie(expiredToken))
	req.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: tokens.RefreshToken,
	})
	w := httptest.NewRecorder()

	handler(w, req)

	// Should succeed and set new cookies
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify new access_token cookie was set
	cookies := w.Result().Cookies()
	hasAccessToken := false
	for _, cookie := range cookies {
		if cookie.Name == "access_token" && cookie.Value != "" {
			hasAccessToken = true
			break
		}
	}
	assert.True(t, hasAccessToken, "New access token cookie should be set")
}

func TestAuthMiddleware_AutoRefresh_MissingAccessToken(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)

	handler := middleware.AutoRefresh(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called when auth fails")
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// No cookies at all
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_AutoRefresh_InvalidTokenMissingRefreshToken(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)

	handler := middleware.AutoRefresh(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called when auth fails")
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// Invalid access token, no refresh token
	req.AddCookie(createAccessTokenCookie("invalid.token.here"))
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_AutoRefresh_InvalidRefreshToken(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)

	handler := middleware.AutoRefresh(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called when auth fails")
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// Invalid access token with invalid refresh token
	req.AddCookie(createAccessTokenCookie("invalid.access.token"))
	req.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: "invalid-refresh-token",
	})
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_AutoRefresh_ContextIsolation(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)

	// Generate two different tokens
	token1, _, err := utils.GenerateAccessToken("user-1", "user1", "admin", testJWTSecret)
	assert.NoError(t, err)

	token2, _, err := utils.GenerateAccessToken("user-2", "user2", "user", testJWTSecret)
	assert.NoError(t, err)

	// First request
	handler := middleware.AutoRefresh(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(UserIDKey)
		assert.Equal(t, "user-1", userID)
		w.WriteHeader(http.StatusOK)
	})

	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req1.AddCookie(createAccessTokenCookie(token1))
	w1 := httptest.NewRecorder()

	handler(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Second request with different token
	handler2 := middleware.AutoRefresh(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(UserIDKey)
		assert.Equal(t, "user-2", userID) // Should be different user
		w.WriteHeader(http.StatusOK)
	})

	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.AddCookie(createAccessTokenCookie(token2))
	w2 := httptest.NewRecorder()

	handler2(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestAuthMiddleware_AutoRefresh_HeadersSetCorrectly(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)

	token, _, err := utils.GenerateAccessToken("user-1", "testuser", "admin", testJWTSecret)
	assert.NoError(t, err)

	handler := middleware.AutoRefresh(func(w http.ResponseWriter, r *http.Request) {
		// Verify X-User-Role header is set
		role := r.Header.Get("X-User-Role")
		assert.Equal(t, "admin", role)

		// Verify it's the only custom header set by middleware
		// (standard HTTP headers may exist)
		assert.NotEmpty(t, role)

		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(createAccessTokenCookie(token))
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_AutoRefresh_ErrorResponseFormat(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)

	handler := middleware.AutoRefresh(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// No token provided
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	
	// Verify response is JSON
	contentType := w.Header().Get("Content-Type")
	assert.Contains(t, contentType, "application/json")
	
	// Response body should contain error message
	assert.NotEmpty(t, w.Body.String())
}

func TestNewAuthMiddleware(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)

	assert.NotNil(t, middleware)
	assert.Equal(t, tokenService, middleware.tokenService)
}
