package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"fortyfour-backend/internal/testhelpers"
)

func setupRateLimiter() (*RateLimiter, *testhelpers.MockRedisClient) {
	mockRedis := testhelpers.NewMockRedisClient()
	config := RateLimiterConfig{
		RequestsPerWindow: 5,
		WindowDuration:    1 * time.Minute,
		KeyPrefix:         "test_rate_limit",
	}
	limiter := NewRateLimiter(mockRedis, config)
	return limiter, mockRedis
}

func TestRateLimiter_LimitByIP_Success(t *testing.T) {
	limiter, _ := setupRateLimiter()

	handler := limiter.LimitByIP(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Check headers
	if w.Header().Get("X-RateLimit-Limit") == "" {
		t.Error("expected X-RateLimit-Limit header")
	}
	if w.Header().Get("X-RateLimit-Remaining") == "" {
		t.Error("expected X-RateLimit-Remaining header")
	}
}

func TestRateLimiter_LimitByIP_Exceeded(t *testing.T) {
	limiter, _ := setupRateLimiter()

	handler := limiter.LimitByIP(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	// w := httptest.NewRecorder()

	// Make requests up to the limit
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		handler(w, req)
	}

	// This should be rate limited
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected status 429, got %d", w.Code)
	}
}

func TestRateLimiter_LimitByUser_Success(t *testing.T) {
	limiter, _ := setupRateLimiter()

	handler := limiter.LimitByUser(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("X-User-ID", "user-1")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestRateLimiter_LimitByUser_NoUserID(t *testing.T) {
	limiter, _ := setupRateLimiter()

	handler := limiter.LimitByUser(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	w := httptest.NewRecorder()

	handler(w, req)

	// Should fall back to IP-based limiting
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestRateLimiter_NewRateLimiter_Defaults(t *testing.T) {
	mockRedis := testhelpers.NewMockRedisClient()
	config := RateLimiterConfig{}
	limiter := NewRateLimiter(mockRedis, config)

	if limiter.config.RequestsPerWindow == 0 {
		t.Error("expected default RequestsPerWindow to be set")
	}
	if limiter.config.WindowDuration == 0 {
		t.Error("expected default WindowDuration to be set")
	}
	if limiter.config.KeyPrefix == "" {
		t.Error("expected default KeyPrefix to be set")
	}
}

// ================================================================
// LimitByAPIKey
// ================================================================

func TestRateLimiter_LimitByAPIKey_Success(t *testing.T) {
	limiter, _ := setupRateLimiter()

	handler := limiter.LimitByAPIKey(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
	req.Header.Set("X-API-Key", "key-abc-123")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if w.Header().Get("X-RateLimit-Limit") == "" {
		t.Error("expected X-RateLimit-Limit header")
	}
	if w.Header().Get("X-RateLimit-Remaining") == "" {
		t.Error("expected X-RateLimit-Remaining header")
	}
}

func TestRateLimiter_LimitByAPIKey_MissingKey_Returns401(t *testing.T) {
	limiter, _ := setupRateLimiter()

	handler := limiter.LimitByAPIKey(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
	// Tidak ada X-API-Key header
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 untuk missing API key, got %d", w.Code)
	}
}

func TestRateLimiter_LimitByAPIKey_Exceeded_Returns429(t *testing.T) {
	limiter, _ := setupRateLimiter() // limit = 5

	handler := limiter.LimitByAPIKey(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
	req.Header.Set("X-API-Key", "limited-key")

	// Habiskan limit
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		handler(w, req)
	}

	// Request ke-6 harus di-rate limit
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", w.Code)
	}
	if w.Header().Get("Retry-After") == "" {
		t.Error("expected Retry-After header saat rate limit tercapai")
	}
}

func TestRateLimiter_LimitByAPIKey_DifferentKeys_IndependentCounters(t *testing.T) {
	limiter, _ := setupRateLimiter()

	handler := limiter.LimitByAPIKey(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Habiskan limit untuk key-A
	reqA := httptest.NewRequest(http.MethodGet, "/api/data", nil)
	reqA.Header.Set("X-API-Key", "key-A")
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		handler(w, reqA)
	}

	// key-B masih harus bisa request
	reqB := httptest.NewRequest(http.MethodGet, "/api/data", nil)
	reqB.Header.Set("X-API-Key", "key-B")
	w := httptest.NewRecorder()
	handler(w, reqB)

	if w.Code != http.StatusOK {
		t.Errorf("key-B seharusnya masih diizinkan, got %d", w.Code)
	}
}

// ================================================================
// getClientIP
// ================================================================

func TestGetClientIP_XForwardedFor(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.5")

	ip := getClientIP(req)

	if ip != "203.0.113.5" {
		t.Errorf("expected '203.0.113.5', got '%s'", ip)
	}
}

func TestGetClientIP_XRealIP_FallbackWhenNoForwarded(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "198.51.100.7")

	ip := getClientIP(req)

	if ip != "198.51.100.7" {
		t.Errorf("expected '198.51.100.7', got '%s'", ip)
	}
}

func TestGetClientIP_RemoteAddr_FallbackWhenNoHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.0.2.1:54321"

	ip := getClientIP(req)

	if ip != "192.0.2.1:54321" {
		t.Errorf("expected '192.0.2.1:54321', got '%s'", ip)
	}
}

func TestGetClientIP_XForwardedFor_TakesPriorityOverXRealIP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.1")
	req.Header.Set("X-Real-IP", "198.51.100.2")

	ip := getClientIP(req)

	if ip != "203.0.113.1" {
		t.Errorf("X-Forwarded-For harus diprioritaskan, got '%s'", ip)
	}
}