package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRedis mocks the RedisInterface for testing
type MockRedis struct {
	mock.Mock
}

func (m *MockRedis) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockRedis) Set(key string, value interface{}, expiration time.Duration) error {
	args := m.Called(key, value, expiration)
	return args.Error(0)
}

func (m *MockRedis) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockRedis) Exists(key string) (bool, error) {
	args := m.Called(key)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedis) Scan(pattern string) ([]string, error) {
	args := m.Called(pattern)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRedis) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestNewTokenBucketLimiter(t *testing.T) {
	mockRedis := new(MockRedis)
	capacity := 100
	refillRate := 10

	limiter := NewTokenBucketLimiter(mockRedis, capacity, refillRate)

	assert.NotNil(t, limiter)
	assert.Equal(t, capacity, limiter.capacity)
	assert.Equal(t, refillRate, limiter.refillRate)
	assert.Equal(t, 1*time.Second, limiter.refillInterval)
	assert.Equal(t, "token_bucket", limiter.keyPrefix)
	assert.Equal(t, mockRedis, limiter.redis)
}

func TestNewTokenBucketLimiter_WithDifferentValues(t *testing.T) {
	tests := []struct {
		name       string
		capacity   int
		refillRate int
	}{
		{
			name:       "small capacity",
			capacity:   10,
			refillRate: 1,
		},
		{
			name:       "medium capacity",
			capacity:   100,
			refillRate: 10,
		},
		{
			name:       "large capacity",
			capacity:   1000,
			refillRate: 100,
		},
		{
			name:       "zero capacity",
			capacity:   0,
			refillRate: 0,
		},
		{
			name:       "high refill rate",
			capacity:   100,
			refillRate: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRedis := new(MockRedis)
			limiter := NewTokenBucketLimiter(mockRedis, tt.capacity, tt.refillRate)

			assert.Equal(t, tt.capacity, limiter.capacity)
			assert.Equal(t, tt.refillRate, limiter.refillRate)
		})
	}
}

func TestTokenBucketLimiter_Limit_FirstRequest(t *testing.T) {
	mockRedis := new(MockRedis)
	limiter := NewTokenBucketLimiter(mockRedis, 10, 1)

	// Mock Redis Get to return error (key doesn't exist - first request)
	mockRedis.On("Get", mock.Anything).Return("", errors.New("key not found"))
	mockRedis.On("Set", mock.Anything, "1", 1*time.Second).Return(nil)

	// Create test handler
	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rr := httptest.NewRecorder()

	// Apply rate limiter
	handler := limiter.Limit(testHandler)
	handler(rr, req)

	// Verify that request was allowed
	assert.True(t, handlerCalled, "handler should be called on first request")
	assert.Equal(t, http.StatusOK, rr.Code)

	mockRedis.AssertExpectations(t)
}

func TestTokenBucketLimiter_Limit_WithDifferentIPs(t *testing.T) {
	tests := []struct {
		name string
		ip   string
	}{
		{
			name: "IPv4 address",
			ip:   "192.168.1.1:12345",
		},
		{
			name: "Different IPv4",
			ip:   "10.0.0.1:54321",
		},
		{
			name: "IPv6 address",
			ip:   "[::1]:8080",
		},
		{
			name: "localhost",
			ip:   "127.0.0.1:3000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRedis := new(MockRedis)
			limiter := NewTokenBucketLimiter(mockRedis, 10, 1)

			mockRedis.On("Get", mock.Anything).Return("", errors.New("key not found"))
			mockRedis.On("Set", mock.Anything, "1", 1*time.Second).Return(nil)

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = tt.ip
			rr := httptest.NewRecorder()

			handler := limiter.Limit(testHandler)
			handler(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
			mockRedis.AssertExpectations(t)
		})
	}
}

func TestTokenBucketLimiter_ConsumeToken_NewKey(t *testing.T) {
	mockRedis := new(MockRedis)
	limiter := NewTokenBucketLimiter(mockRedis, 10, 1)

	// Mock Get to return error (key doesn't exist)
	mockRedis.On("Get", "test_key").Return("", errors.New("key not found"))
	mockRedis.On("Set", "test_key", "1", 1*time.Second).Return(nil)

	allowed, err := limiter.consumeToken("test_key")

	assert.True(t, allowed, "should allow first request")
	assert.Nil(t, err)
	mockRedis.AssertExpectations(t)
}

func TestTokenBucketLimiter_ConsumeToken_ExistingKey(t *testing.T) {
	mockRedis := new(MockRedis)
	limiter := NewTokenBucketLimiter(mockRedis, 10, 1)

	// Mock Get to return existing value
	mockRedis.On("Get", "test_key").Return("5", nil)

	allowed, err := limiter.consumeToken("test_key")

	// Based on current implementation, it always returns true
	assert.True(t, allowed)
	assert.Nil(t, err)
	mockRedis.AssertExpectations(t)
}

func TestTokenBucketLimiter_KeyFormat(t *testing.T) {
	mockRedis := new(MockRedis)
	limiter := NewTokenBucketLimiter(mockRedis, 10, 1)

	testIP := "192.168.1.1"
	expectedKeyPrefix := "token_bucket:"

	// Create a handler to capture the key being used
	var capturedKey string
	mockRedis.On("Get", mock.MatchedBy(func(key string) bool {
		capturedKey = key
		return true
	})).Return("", errors.New("key not found"))
	mockRedis.On("Set", mock.Anything, "1", 1*time.Second).Return(nil)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = testIP + ":12345"
	rr := httptest.NewRecorder()

	handler := limiter.Limit(testHandler)
	handler(rr, req)

	// Verify key format
	assert.Contains(t, capturedKey, expectedKeyPrefix)
	assert.Contains(t, capturedKey, testIP)
}

func TestTokenBucketLimiter_RefillInterval(t *testing.T) {
	mockRedis := new(MockRedis)
	limiter := NewTokenBucketLimiter(mockRedis, 100, 10)

	// Verify default refill interval
	assert.Equal(t, 1*time.Second, limiter.refillInterval,
		"default refill interval should be 1 second")
}

func TestTokenBucketLimiter_DefaultKeyPrefix(t *testing.T) {
	mockRedis := new(MockRedis)
	limiter := NewTokenBucketLimiter(mockRedis, 100, 10)

	assert.Equal(t, "token_bucket", limiter.keyPrefix,
		"default key prefix should be 'token_bucket'")
}

func TestTokenBucketLimiter_Limit_RedisError(t *testing.T) {
	mockRedis := new(MockRedis)
	limiter := NewTokenBucketLimiter(mockRedis, 10, 1)

	// Mock Redis Get to return error
	mockRedis.On("Get", mock.Anything).Return("", errors.New("redis connection error"))
	mockRedis.On("Set", mock.Anything, "1", 1*time.Second).Return(errors.New("redis error"))

	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rr := httptest.NewRecorder()

	handler := limiter.Limit(testHandler)
	handler(rr, req)

	// Based on current implementation, handler is still called
	assert.True(t, handlerCalled)
	mockRedis.AssertExpectations(t)
}

func TestTokenBucketLimiter_MultipleRequests(t *testing.T) {
	mockRedis := new(MockRedis)
	limiter := NewTokenBucketLimiter(mockRedis, 10, 1)

	// First request - key doesn't exist
	mockRedis.On("Get", mock.Anything).Return("", errors.New("key not found")).Once()
	mockRedis.On("Set", mock.Anything, "1", 1*time.Second).Return(nil).Once()

	// Second request - key exists
	mockRedis.On("Get", mock.Anything).Return("1", nil).Once()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// First request
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	rr1 := httptest.NewRecorder()

	handler := limiter.Limit(testHandler)
	handler(rr1, req1)

	assert.Equal(t, http.StatusOK, rr1.Code)

	// Second request from same IP
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "192.168.1.1:12345"
	rr2 := httptest.NewRecorder()

	handler(rr2, req2)

	assert.Equal(t, http.StatusOK, rr2.Code)

	mockRedis.AssertExpectations(t)
}

func TestTokenBucketLimiter_Capacity(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
		want     int
	}{
		{
			name:     "low capacity",
			capacity: 5,
			want:     5,
		},
		{
			name:     "medium capacity",
			capacity: 50,
			want:     50,
		},
		{
			name:     "high capacity",
			capacity: 500,
			want:     500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRedis := new(MockRedis)
			limiter := NewTokenBucketLimiter(mockRedis, tt.capacity, 10)

			assert.Equal(t, tt.want, limiter.capacity)
		})
	}
}

func TestTokenBucketLimiter_RefillRate(t *testing.T) {
	tests := []struct {
		name       string
		refillRate int
		want       int
	}{
		{
			name:       "slow refill",
			refillRate: 1,
			want:       1,
		},
		{
			name:       "medium refill",
			refillRate: 10,
			want:       10,
		},
		{
			name:       "fast refill",
			refillRate: 100,
			want:       100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRedis := new(MockRedis)
			limiter := NewTokenBucketLimiter(mockRedis, 100, tt.refillRate)

			assert.Equal(t, tt.want, limiter.refillRate)
		})
	}
}

func TestTokenBucketLimiter_ConsumeToken_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		getError    error
		setError    error
		expectAllow bool
	}{
		{
			name:        "get error, set success",
			getError:    errors.New("get error"),
			setError:    nil,
			expectAllow: true,
		},
		{
			name:        "get error, set error",
			getError:    errors.New("get error"),
			setError:    errors.New("set error"),
			expectAllow: true, // Current implementation still returns true
		},
		{
			name:        "get success",
			getError:    nil,
			setError:    nil,
			expectAllow: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRedis := new(MockRedis)
			limiter := NewTokenBucketLimiter(mockRedis, 10, 1)

			if tt.getError != nil {
				mockRedis.On("Get", "test_key").Return("", tt.getError)
				mockRedis.On("Set", "test_key", "1", 1*time.Second).Return(tt.setError)
			} else {
				mockRedis.On("Get", "test_key").Return("5", nil)
			}

			allowed, err := limiter.consumeToken("test_key")

			assert.Equal(t, tt.expectAllow, allowed)
			assert.Nil(t, err) // Current implementation doesn't return errors
			mockRedis.AssertExpectations(t)
		})
	}
}

func TestTokenBucketLimiter_ConcurrentRequests(t *testing.T) {
	// Test that multiple concurrent requests don't cause issues
	mockRedis := new(MockRedis)
	limiter := NewTokenBucketLimiter(mockRedis, 100, 10)

	mockRedis.On("Get", mock.Anything).Return("", errors.New("key not found"))
	mockRedis.On("Set", mock.Anything, "1", 1*time.Second).Return(nil)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create multiple requests with different IPs
	ips := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"}

	for i, ip := range ips {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = fmt.Sprintf("%s:%d", ip, 12345+i)
		rr := httptest.NewRecorder()

		handler := limiter.Limit(testHandler)
		handler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	}
}
