package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	ctx := SetUserID(context.TODO(), "123")

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
	req.Header.Set("X-User-Role", "admin")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if userID != "42" {
		t.Errorf("expected userID 42, got %s", userID)
	}

	if role != "admin" {
		t.Errorf("expected role admin, got %s", role)
	}
}

func TestAuthMiddleware_FallsBackToGatewayHeaders(t *testing.T) {
	var userID, role string

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID = GetUserID(r.Context())
		role = GetRole(r.Context())
	})

	handler := Auth(h)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer real-jwt-from-gateway")
	req.Header.Set("X-User-ID", "user-123")
	req.Header.Set("X-User-Role", "user_pic")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	if userID != "user-123" {
		t.Errorf("expected userID user-123, got %s", userID)
	}

	if role != "user_pic" {
		t.Errorf("expected role user_pic, got %s", role)
	}
}

func TestAuthMiddleware_AcceptsJWTAccessToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "user-456",
		"role":    "user_pic",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	var userID, role string

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID = GetUserID(r.Context())
		role = GetRole(r.Context())
	})

	handler := Auth(h)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	if userID != "user-456" {
		t.Errorf("expected userID user-456, got %s", userID)
	}

	if role != "user_pic" {
		t.Errorf("expected role user_pic, got %s", role)
	}
}

func TestAuthMiddleware_AllowsRequestWithoutIdentity(t *testing.T) {
	called := false

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		if GetUserID(r.Context()) != "" {
			t.Errorf("expected empty user id")
		}

		if GetRole(r.Context()) != "" {
			t.Errorf("expected empty role")
		}
	})

	handler := Auth(h)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if !called {
		t.Error("handler not called")
	}

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
