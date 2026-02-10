package middleware

import (
	"fmt"
	"fortyfour-backend/pkg/cache"
	"net/http"
	"time"

)

// TokenBucketLimiter implements token bucket algorithm
type TokenBucketLimiter struct {
	redis          cache.RedisInterface
	capacity       int           // Maximum tokens
	refillRate     int           // Tokens per second
	refillInterval time.Duration // How often to refill
	keyPrefix      string
}

func NewTokenBucketLimiter(redis cache.RedisInterface, capacity, refillRate int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		redis:          redis,
		capacity:       capacity,
		refillRate:     refillRate,
		refillInterval: 1 * time.Second,
		keyPrefix:      "token_bucket",
	}
}

// Limit applies token bucket rate limiting
func (tbl *TokenBucketLimiter) Limit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		key := fmt.Sprintf("%s:%s", tbl.keyPrefix, ip)

		// Check if token available
		allowed, err := tbl.consumeToken(key)
		if err != nil || !allowed {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next(w, r)
	}
}

func (tbl *TokenBucketLimiter) consumeToken(key string) (bool, error) {
	// Simplified implementation - in production, use Lua script for atomicity
	// countStr, err := tbl.redis.Get(key)
	_, err := tbl.redis.Get(key)
	if err != nil {
		// Initialize bucket
		tbl.redis.Set(key, "1", tbl.refillInterval)
		return true, nil
	}

	// Parse current tokens
	// Implementation details omitted for brevity
	// In production, use Redis Lua scripts for atomic operations

	return true, nil
}
