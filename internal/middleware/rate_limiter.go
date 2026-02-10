package middleware

import (
	"fmt"
	"fortyfour-backend/internal/utils"
	"fortyfour-backend/pkg/cache"
	"net/http"
	"strconv"
	"time"

)

type RateLimiterConfig struct {
	RequestsPerWindow int           // Number of requests allowed
	WindowDuration    time.Duration // Time window (e.g., 1 minute)
	KeyPrefix         string        // Redis key prefix
}

type RateLimiter struct {
	redis  cache.RedisInterface
	config RateLimiterConfig
}

func NewRateLimiter(redis cache.RedisInterface, config RateLimiterConfig) *RateLimiter {
	// Set defaults if not provided
	if config.RequestsPerWindow == 0 {
		config.RequestsPerWindow = 100 // Default: 100 requests
	}
	if config.WindowDuration == 0 {
		config.WindowDuration = 1 * time.Minute // Default: per minute
	}
	if config.KeyPrefix == "" {
		config.KeyPrefix = "rate_limit"
	}

	return &RateLimiter{
		redis:  redis,
		config: config,
	}
}

// LimitByIP rate limits based on client IP address
func (rl *RateLimiter) LimitByIP(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		key := fmt.Sprintf("%s:ip:%s", rl.config.KeyPrefix, ip)

		allowed, remaining, resetTime, err := rl.checkLimit(key)
		if err != nil {
			// On error, log and allow request (fail open)
			// In production, you might want to fail closed instead
			next(w, r)
			return
		}

		// Set rate limit headers
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.config.RequestsPerWindow))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		if !allowed {
			w.Header().Set("Retry-After", strconv.Itoa(int(time.Until(resetTime).Seconds())))
			utils.RespondError(w, http.StatusTooManyRequests, "Rate limit exceeded. Please try again later.")
			return
		}

		next(w, r)
	}
}

// LimitByUser rate limits based on authenticated user ID
func (rl *RateLimiter) LimitByUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			// If no user ID, fall back to IP-based limiting
			rl.LimitByIP(next)(w, r)
			return
		}

		key := fmt.Sprintf("%s:user:%s", rl.config.KeyPrefix, userID)

		allowed, remaining, resetTime, err := rl.checkLimit(key)
		if err != nil {
			next(w, r)
			return
		}

		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.config.RequestsPerWindow))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		if !allowed {
			w.Header().Set("Retry-After", strconv.Itoa(int(time.Until(resetTime).Seconds())))
			utils.RespondError(w, http.StatusTooManyRequests, "Rate limit exceeded. Please try again later.")
			return
		}

		next(w, r)
	}
}

// LimitByAPIKey rate limits based on API key
func (rl *RateLimiter) LimitByAPIKey(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			utils.RespondError(w, http.StatusUnauthorized, "API key required")
			return
		}

		key := fmt.Sprintf("%s:apikey:%s", rl.config.KeyPrefix, apiKey)

		allowed, remaining, resetTime, err := rl.checkLimit(key)
		if err != nil {
			next(w, r)
			return
		}

		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.config.RequestsPerWindow))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		if !allowed {
			w.Header().Set("Retry-After", strconv.Itoa(int(time.Until(resetTime).Seconds())))
			utils.RespondError(w, http.StatusTooManyRequests, "Rate limit exceeded. Please try again later.")
			return
		}

		next(w, r)
	}
}

// checkLimit implements the sliding window rate limiting algorithm
func (rl *RateLimiter) checkLimit(key string) (allowed bool, remaining int, resetTime time.Time, err error) {
	now := time.Now()
	// windowStart := now.Add(-rl.config.WindowDuration)

	// Get current count
	countStr, err := rl.redis.Get(key)
	if err != nil {
		// Key doesn't exist, this is the first request
		count := 1
		if err := rl.redis.Set(key, strconv.Itoa(count), rl.config.WindowDuration); err != nil {
			return false, 0, now, err
		}
		resetTime = now.Add(rl.config.WindowDuration)
		return true, rl.config.RequestsPerWindow - count, resetTime, nil
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		return false, 0, now, err
	}

	// Check if limit exceeded
	if count >= rl.config.RequestsPerWindow {
		resetTime = now.Add(rl.config.WindowDuration)
		return false, 0, resetTime, nil
	}

	// Increment counter
	count++
	if err := rl.redis.Set(key, strconv.Itoa(count), rl.config.WindowDuration); err != nil {
		return false, 0, now, err
	}

	resetTime = now.Add(rl.config.WindowDuration)
	remaining = rl.config.RequestsPerWindow - count
	return true, remaining, resetTime, nil
}

// getClientIP extracts the real client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (proxy/load balancer)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Take the first IP if there are multiple
		return forwarded
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}
