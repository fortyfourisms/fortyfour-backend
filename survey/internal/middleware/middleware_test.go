package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ADAPT
func TestAdaptHandler(t *testing.T) {
	called := false

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	handler := AdaptHandler(h)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if !called {
		t.Error("handler not called")
	}
}

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

func TestCORS_Options(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not be called")
	})

	handler := CORS(h)

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200")
	}
}

// CONTEXT
func TestContext_SetGet(t *testing.T) {
	ctx := context.Background()

	ctx = SetUserID(ctx, "123")
	ctx = SetRole(ctx, "admin")

	if GetUserID(ctx) != "123" {
		t.Error("user id mismatch")
	}

	if GetRole(ctx) != "admin" {
		t.Error("role mismatch")
	}
}

func TestContext_NilSafe(t *testing.T) {
	ctx := SetUserID(nil, "123")

	if GetUserID(ctx) != "123" {
		t.Error("should handle nil context safely")
	}
}

// AUTH
func TestAuthMiddleware(t *testing.T) {
	var userID, role string

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID = GetUserID(r.Context())
		role = GetRole(r.Context())
	})

	handler := Auth(h)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-User-ID", "42")
	req.Header.Set("X-Role", "admin")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if userID != "42" {
		t.Errorf("expected userID 42, got %s", userID)
	}

	if role != "admin" {
		t.Errorf("expected role admin, got %s", role)
	}
}