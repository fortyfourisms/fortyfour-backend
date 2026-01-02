package middleware

import (
	"fortyfour-backend/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware_Authenticate_MissingHeader(t *testing.T) {
	middleware := NewAuthMiddleware("test-secret")
	handler := middleware.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_Authenticate_InvalidFormat(t *testing.T) {
	middleware := NewAuthMiddleware("test-secret")
	handler := middleware.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_Authenticate_InvalidToken(t *testing.T) {
	middleware := NewAuthMiddleware("test-secret")
	handler := middleware.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_Authenticate_ValidToken(t *testing.T) {
	middleware := NewAuthMiddleware("test-secret")
	
	// Generate a valid token
	token, _, err := utils.GenerateAccessToken("user-1", "testuser", "admin", "test-secret")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

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
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestAuthMiddleware_Authenticate_ContextValues(t *testing.T) {
	middleware := NewAuthMiddleware("test-secret")
	
	// Generate a valid token
	token, _, err := utils.GenerateAccessToken("user-1", "testuser", "admin", "test-secret")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	handler := middleware.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(UserIDKey)
		username := r.Context().Value(Username)
		role := r.Context().Value(Role)

		if userID != "user-1" {
			t.Errorf("expected user ID 'user-1', got '%v'", userID)
		}
		if username != "testuser" {
			t.Errorf("expected username 'testuser', got '%v'", username)
		}
		if role != "admin" {
			t.Errorf("expected role 'admin', got '%v'", role)
		}

		// Check headers
		if r.Header.Get("X-User-ID") != "user-1" {
			t.Errorf("expected X-User-ID header 'user-1', got '%s'", r.Header.Get("X-User-ID"))
		}
		if r.Header.Get("X-Username") != "testuser" {
			t.Errorf("expected X-Username header 'testuser', got '%s'", r.Header.Get("X-Username"))
		}
		if r.Header.Get("X-User-Role") != "admin" {
			t.Errorf("expected X-User-Role header 'admin', got '%s'", r.Header.Get("X-User-Role"))
		}

		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

