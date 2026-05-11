package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// RECOVERY
func TestRecovery_Panic(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	handler := Recovery(h)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

// LOGGER
func TestLogger(t *testing.T) {
	called := false

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	handler := Logger(h)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if !called {
		t.Error("handler not called")
	}
}

// CORS
func TestCORS_Headers(t *testing.T) {
	called := false

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	handler := CORS(h)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if !called {
		t.Error("handler not called")
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("CORS header missing")
	}
}

// CONTEXT
func TestContext_SetGet(t *testing.T) {
	ctx := context.Background()

	ctx = context.WithValue(ctx, UserIDKey, "123")
	ctx = context.WithValue(ctx, RoleKey, "admin")
	ctx = context.WithValue(ctx, PerusahaanIDKey, "456")

	if GetUserID(ctx) != "123" {
		t.Error("user id mismatch")
	}

	if GetRole(ctx) != "admin" {
		t.Error("role mismatch")
	}

	if GetPerusahaanID(ctx) != "456" {
		t.Error("perusahaan id mismatch")
	}
}

// AUTH
func TestAuthMiddleware_Success(t *testing.T) {
	internalKey := "secret-key"
	authM := NewAuthMiddleware(internalKey)

	var userID, role string

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID = GetUserID(r.Context())
		role = GetRole(r.Context())
	})

	handler := authM.Authenticate(h)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Internal-Key", internalKey)
	req.Header.Set("X-User-ID", "42")
	req.Header.Set("X-User-Role", "admin")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	if userID != "42" {
		t.Errorf("expected userID 42, got %s", userID)
	}

	if role != "admin" {
		t.Errorf("expected role admin, got %s", role)
	}
}

func TestAuthMiddleware_Unauthorized_InvalidInternalKey(t *testing.T) {
	internalKey := "secret-key"
	authM := NewAuthMiddleware(internalKey)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	handler := authM.Authenticate(h)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Internal-Key", "wrong-key")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_Unauthorized_MissingUserHeaders(t *testing.T) {
	internalKey := "secret-key"
	authM := NewAuthMiddleware(internalKey)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	handler := authM.Authenticate(h)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Internal-Key", internalKey)
	// Missing X-User-ID and X-User-Role

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}
