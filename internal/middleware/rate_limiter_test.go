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
	w := httptest.NewRecorder()

	// Make requests up to the limit
	for i := 0; i < 5; i++ {
		w = httptest.NewRecorder()
		handler(w, req)
	}

	// This should be rate limited
	w = httptest.NewRecorder()
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
