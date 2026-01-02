package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casbin/casbin/v2"
)

func setupAuthorizationMiddleware() (*AuthorizationMiddleware, *casbin.Enforcer) {
	enforcer, _ := casbin.NewEnforcer("../../casbin_model.conf", false)
	middleware := NewAuthorizationMiddleware(enforcer)
	return middleware, enforcer
}

func TestAuthorizationMiddleware_Authorize_NoRole(t *testing.T) {
	middleware, _ := setupAuthorizationMiddleware()

	handler := middleware.Authorize(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", w.Code)
	}
}

func TestAuthorizationMiddleware_Authorize_WithRole(t *testing.T) {
	middleware, enforcer := setupAuthorizationMiddleware()

	// Add policy for testing
	enforcer.AddPolicy("admin", "/api/test", "GET")

	handler := middleware.Authorize(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	ctx := context.WithValue(req.Context(), Role, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestAuthorizationMiddleware_Authorize_InsufficientPermissions(t *testing.T) {
	middleware, _ := setupAuthorizationMiddleware()

	handler := middleware.Authorize(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	ctx := context.WithValue(req.Context(), Role, "user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", w.Code)
	}
}

func TestAuthorizationMiddleware_extractResourcePath(t *testing.T) {
	middleware, _ := setupAuthorizationMiddleware()

	testCases := []struct {
		path     string
		expected string
	}{
		{"/api/posts/123", "/api/posts"},
		{"/api/users", "/api/users"},
		{"/api/test/456/789", "/api/test"},
		{"/api", "/api"},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			result := middleware.extractResourcePath(tc.path)
			if result != tc.expected {
				t.Errorf("expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}
