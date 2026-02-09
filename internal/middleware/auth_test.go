package middleware

import (
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/testhelpers"
	"fortyfour-backend/internal/utils"
	"net/http"
	"net/http/httptest"
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

func TestNewAuthMiddleware(t *testing.T) {
	tokenService := createTestTokenService()
	middleware := NewAuthMiddleware(tokenService)

	assert.NotNil(t, middleware)
	assert.Equal(t, tokenService, middleware.tokenService)
}
