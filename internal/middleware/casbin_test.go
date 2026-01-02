package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casbin/casbin/v2"
)

func setupCasbinMiddleware() (*CasbinMiddleware, *casbin.Enforcer) {
	enforcer, _ := casbin.NewEnforcer("../../casbin_model.conf", false)
	middleware := NewCasbinMiddleware(enforcer)
	return middleware, enforcer
}

func TestCasbinMiddleware_Authorize_NoRole(t *testing.T) {
	middleware, _ := setupCasbinMiddleware()

	handler := middleware.Authorize(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestCasbinMiddleware_Authorize_WithRole(t *testing.T) {
	middleware, enforcer := setupCasbinMiddleware()

	// Add policy for testing
	enforcer.AddPolicy("admin", "/api/test", "GET")

	handler := middleware.Authorize(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("X-User-Role", "admin")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestCasbinMiddleware_Authorize_InsufficientPermissions(t *testing.T) {
	middleware, _ := setupCasbinMiddleware()

	handler := middleware.Authorize(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("X-User-Role", "user")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", w.Code)
	}
}

