package middleware

import (
	"time"
)

// RateLimitConfigs provides different rate limiting strategies
type RateLimitConfigs struct {
	// Strict: For sensitive endpoints (login, register)
	Strict RateLimiterConfig

	// Moderate: For authenticated API endpoints
	Moderate RateLimiterConfig

	// Lenient: For public read-only endpoints
	Lenient RateLimiterConfig
}

func GetRateLimitConfigs() RateLimitConfigs {
	return RateLimitConfigs{
		Strict: RateLimiterConfig{
			RequestsPerWindow: 10,
			WindowDuration:    1 * time.Minute,
			KeyPrefix:         "rate_limit_strict",
		},
		Moderate: RateLimiterConfig{
			RequestsPerWindow: 100,
			WindowDuration:    1 * time.Minute,
			KeyPrefix:         "rate_limit_moderate",
		},
		Lenient: RateLimiterConfig{
			RequestsPerWindow: 1000,
			WindowDuration:    1 * time.Minute,
			KeyPrefix:         "rate_limit_lenient",
		},
	}
}
