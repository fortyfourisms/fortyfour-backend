package handlers

import (
	"context"
	"fortyfour-backend/internal/middleware"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProxyHandler(t *testing.T) {
	t.Run("valid URL", func(t *testing.T) {
		h := NewProxyHandler("http://localhost:8081", "secret-key")
		require.NotNil(t, h)
		assert.Equal(t, "localhost:8081", h.target.Host)
		assert.Equal(t, "http", h.target.Scheme)
		assert.Equal(t, "secret-key", h.internalKey)
		assert.NotNil(t, h.proxy)
	})
}

func TestProxyHandler_ServeHTTP_InjectsHeaders(t *testing.T) {
	// Spin up a fake backend that captures the incoming request headers.
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Echo the headers we care about back as response headers so the test
		// can inspect them without coupling to the request object.
		w.Header().Set("X-Got-User-ID", r.Header.Get("X-User-ID"))
		w.Header().Set("X-Got-User-Role", r.Header.Get("X-User-Role"))
		w.Header().Set("X-Got-Internal-Key", r.Header.Get("X-Internal-Key"))
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	h := NewProxyHandler(backend.URL, "my-internal-key")

	t.Run("with auth context values", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/some-path", nil)
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-123")
		ctx = context.WithValue(ctx, middleware.RoleKey, "admin")
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "user-123", rr.Header().Get("X-Got-User-ID"))
		assert.Equal(t, "admin", rr.Header().Get("X-Got-User-Role"))
		assert.Equal(t, "my-internal-key", rr.Header().Get("X-Got-Internal-Key"))
	})

	t.Run("without auth context - headers are empty strings", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/other-path", nil)

		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "", rr.Header().Get("X-Got-User-ID"))
		assert.Equal(t, "", rr.Header().Get("X-Got-User-Role"))
		assert.Equal(t, "my-internal-key", rr.Header().Get("X-Got-Internal-Key"))
	})
}

func TestProxyHandler_ServeHTTP_ForwardsRequest(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Got-Method", r.Method)
		w.Header().Set("X-Got-Path", r.URL.Path)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer backend.Close()

	h := NewProxyHandler(backend.URL, "key")

	t.Run("preserves method and path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/resource", nil)

		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, "POST", rr.Header().Get("X-Got-Method"))
		assert.Equal(t, "/api/v1/resource", rr.Header().Get("X-Got-Path"))
		assert.Equal(t, `{"status":"ok"}`, rr.Body.String())
	})

	t.Run("forwards query parameters", func(t *testing.T) {
		backend2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Got-Query", r.URL.RawQuery)
			w.WriteHeader(http.StatusOK)
		}))
		defer backend2.Close()

		h2 := NewProxyHandler(backend2.URL, "key")
		req := httptest.NewRequest(http.MethodGet, "/search?q=test&page=1", nil)

		rr := httptest.NewRecorder()
		h2.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "q=test&page=1", rr.Header().Get("X-Got-Query"))
	})
}

func TestProxyHandler_ServeHTTP_BackendDown(t *testing.T) {
	// Point to a URL that is not listening.
	h := NewProxyHandler("http://127.0.0.1:19999", "key")

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	// ReverseProxy returns 502 Bad Gateway when the backend is unreachable.
	assert.Equal(t, http.StatusBadGateway, rr.Code)
}

func TestProxyHandler_ServeHTTP_ForwardsRequestBody(t *testing.T) {
	// Verifikasi body POST/PUT diteruskan ke backend dengan utuh
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, 1024)
		n, _ := r.Body.Read(buf)
		w.Header().Set("X-Got-Body", string(buf[:n]))
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	h := NewProxyHandler(backend.URL, "key")

	bodyPayload := `{"nama":"test","value":42}`
	req := httptest.NewRequest(http.MethodPost, "/api/resource", strings.NewReader(bodyPayload))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, bodyPayload, rr.Header().Get("X-Got-Body"))
}

func TestProxyHandler_ServeHTTP_ForwardsCustomRequestHeaders(t *testing.T) {
	// Verifikasi header custom dari client diteruskan ke backend
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Got-Accept", r.Header.Get("Accept"))
		w.Header().Set("X-Got-Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	h := NewProxyHandler(backend.URL, "key")

	req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("X-Got-Accept"))
	assert.Equal(t, "application/json", rr.Header().Get("X-Got-Content-Type"))
}

func TestProxyHandler_ServeHTTP_OnlyUserIDInContext(t *testing.T) {
	// Verifikasi partial context: hanya UserID ada, Role kosong
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Got-User-ID", r.Header.Get("X-User-ID"))
		w.Header().Set("X-Got-User-Role", r.Header.Get("X-User-Role"))
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	h := NewProxyHandler(backend.URL, "partial-key")

	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-only-999")
	// RoleKey tidak diset — harus menghasilkan header kosong
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "user-only-999", rr.Header().Get("X-Got-User-ID"))
	assert.Equal(t, "", rr.Header().Get("X-Got-User-Role"),
		"Role header harus kosong jika RoleKey tidak ada di context")
}

func TestProxyHandler_ServeHTTP_InternalKeyAlwaysInjected(t *testing.T) {
	// Verifikasi X-Internal-Key selalu dikirim terlepas dari ada tidaknya context
	internalKey := "super-internal-secret"
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Got-Internal-Key", r.Header.Get("X-Internal-Key"))
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	h := NewProxyHandler(backend.URL, internalKey)

	cases := []struct {
		name string
		ctx  context.Context
	}{
		{"dengan context auth", func() context.Context {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			ctx := context.WithValue(req.Context(), middleware.UserIDKey, "u1")
			return context.WithValue(ctx, middleware.RoleKey, "admin")
		}()},
		{"tanpa context auth", context.Background()},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/check", nil)
			req = req.WithContext(c.ctx)

			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)

			assert.Equal(t, internalKey, rr.Header().Get("X-Got-Internal-Key"),
				"X-Internal-Key harus selalu dikirim")
		})
	}
}
