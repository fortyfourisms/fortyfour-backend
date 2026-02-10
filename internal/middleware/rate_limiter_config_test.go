package middleware

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetRateLimitConfigs(t *testing.T) {
	configs := GetRateLimitConfigs()

	// Verify that configs are returned
	assert.NotNil(t, configs)

	// Verify Strict config
	t.Run("Strict config", func(t *testing.T) {
		assert.Equal(t, 5, configs.Strict.RequestsPerWindow, "Strict should allow 5 requests")
		assert.Equal(t, 1*time.Minute, configs.Strict.WindowDuration, "Strict window should be 1 minute")
		assert.Equal(t, "rate_limit_strict", configs.Strict.KeyPrefix, "Strict key prefix should be correct")
	})

	// Verify Moderate config
	t.Run("Moderate config", func(t *testing.T) {
		assert.Equal(t, 100, configs.Moderate.RequestsPerWindow, "Moderate should allow 100 requests")
		assert.Equal(t, 1*time.Minute, configs.Moderate.WindowDuration, "Moderate window should be 1 minute")
		assert.Equal(t, "rate_limit_moderate", configs.Moderate.KeyPrefix, "Moderate key prefix should be correct")
	})

	// Verify Lenient config
	t.Run("Lenient config", func(t *testing.T) {
		assert.Equal(t, 1000, configs.Lenient.RequestsPerWindow, "Lenient should allow 1000 requests")
		assert.Equal(t, 1*time.Minute, configs.Lenient.WindowDuration, "Lenient window should be 1 minute")
		assert.Equal(t, "rate_limit_lenient", configs.Lenient.KeyPrefix, "Lenient key prefix should be correct")
	})
}

func TestRateLimitConfigs_Hierarchy(t *testing.T) {
	configs := GetRateLimitConfigs()

	// Verify that strict < moderate < lenient in terms of allowed requests
	assert.Less(t, configs.Strict.RequestsPerWindow, configs.Moderate.RequestsPerWindow,
		"Strict should allow fewer requests than Moderate")
	assert.Less(t, configs.Moderate.RequestsPerWindow, configs.Lenient.RequestsPerWindow,
		"Moderate should allow fewer requests than Lenient")
}

func TestRateLimitConfigs_AllHaveSameWindow(t *testing.T) {
	configs := GetRateLimitConfigs()

	// All configs should use the same time window for consistency
	expectedWindow := 1 * time.Minute

	assert.Equal(t, expectedWindow, configs.Strict.WindowDuration)
	assert.Equal(t, expectedWindow, configs.Moderate.WindowDuration)
	assert.Equal(t, expectedWindow, configs.Lenient.WindowDuration)
}

func TestRateLimitConfigs_KeyPrefixesAreUnique(t *testing.T) {
	configs := GetRateLimitConfigs()

	// Key prefixes should be unique to avoid conflicts
	assert.NotEqual(t, configs.Strict.KeyPrefix, configs.Moderate.KeyPrefix)
	assert.NotEqual(t, configs.Strict.KeyPrefix, configs.Lenient.KeyPrefix)
	assert.NotEqual(t, configs.Moderate.KeyPrefix, configs.Lenient.KeyPrefix)
}

func TestRateLimitConfigs_KeyPrefixesNotEmpty(t *testing.T) {
	configs := GetRateLimitConfigs()

	assert.NotEmpty(t, configs.Strict.KeyPrefix, "Strict key prefix should not be empty")
	assert.NotEmpty(t, configs.Moderate.KeyPrefix, "Moderate key prefix should not be empty")
	assert.NotEmpty(t, configs.Lenient.KeyPrefix, "Lenient key prefix should not be empty")
}

func TestRateLimitConfigs_PositiveValues(t *testing.T) {
	configs := GetRateLimitConfigs()

	// All request counts should be positive
	assert.Greater(t, configs.Strict.RequestsPerWindow, 0, "Strict requests should be > 0")
	assert.Greater(t, configs.Moderate.RequestsPerWindow, 0, "Moderate requests should be > 0")
	assert.Greater(t, configs.Lenient.RequestsPerWindow, 0, "Lenient requests should be > 0")

	// All window durations should be positive
	assert.Greater(t, configs.Strict.WindowDuration, time.Duration(0), "Strict window should be > 0")
	assert.Greater(t, configs.Moderate.WindowDuration, time.Duration(0), "Moderate window should be > 0")
	assert.Greater(t, configs.Lenient.WindowDuration, time.Duration(0), "Lenient window should be > 0")
}

func TestRateLimitConfigs_Idempotency(t *testing.T) {
	// Calling GetRateLimitConfigs multiple times should return same values
	configs1 := GetRateLimitConfigs()
	configs2 := GetRateLimitConfigs()

	assert.Equal(t, configs1.Strict.RequestsPerWindow, configs2.Strict.RequestsPerWindow)
	assert.Equal(t, configs1.Moderate.RequestsPerWindow, configs2.Moderate.RequestsPerWindow)
	assert.Equal(t, configs1.Lenient.RequestsPerWindow, configs2.Lenient.RequestsPerWindow)

	assert.Equal(t, configs1.Strict.WindowDuration, configs2.Strict.WindowDuration)
	assert.Equal(t, configs1.Moderate.WindowDuration, configs2.Moderate.WindowDuration)
	assert.Equal(t, configs1.Lenient.WindowDuration, configs2.Lenient.WindowDuration)

	assert.Equal(t, configs1.Strict.KeyPrefix, configs2.Strict.KeyPrefix)
	assert.Equal(t, configs1.Moderate.KeyPrefix, configs2.Moderate.KeyPrefix)
	assert.Equal(t, configs1.Lenient.KeyPrefix, configs2.Lenient.KeyPrefix)
}

func TestRateLimitConfigs_StrictForSensitiveEndpoints(t *testing.T) {
	configs := GetRateLimitConfigs()

	// Strict should be very restrictive (suitable for login/register)
	assert.LessOrEqual(t, configs.Strict.RequestsPerWindow, 10,
		"Strict config should be very restrictive for sensitive endpoints")
}

func TestRateLimitConfigs_ModerateForAPI(t *testing.T) {
	configs := GetRateLimitConfigs()

	// Moderate should be reasonable for API usage
	assert.GreaterOrEqual(t, configs.Moderate.RequestsPerWindow, 50,
		"Moderate should allow reasonable API usage")
	assert.LessOrEqual(t, configs.Moderate.RequestsPerWindow, 500,
		"Moderate should still have limits")
}

func TestRateLimitConfigs_LenientForPublic(t *testing.T) {
	configs := GetRateLimitConfigs()

	// Lenient should allow high traffic for public endpoints
	assert.GreaterOrEqual(t, configs.Lenient.RequestsPerWindow, 500,
		"Lenient should allow high traffic for public endpoints")
}

func TestRateLimitConfigs_WindowDurationReasonable(t *testing.T) {
	configs := GetRateLimitConfigs()

	// Window should be reasonable (not too short, not too long)
	minWindow := 1 * time.Second
	maxWindow := 1 * time.Hour

	assert.GreaterOrEqual(t, configs.Strict.WindowDuration, minWindow)
	assert.LessOrEqual(t, configs.Strict.WindowDuration, maxWindow)

	assert.GreaterOrEqual(t, configs.Moderate.WindowDuration, minWindow)
	assert.LessOrEqual(t, configs.Moderate.WindowDuration, maxWindow)

	assert.GreaterOrEqual(t, configs.Lenient.WindowDuration, minWindow)
	assert.LessOrEqual(t, configs.Lenient.WindowDuration, maxWindow)
}

func TestRateLimitConfigs_KeyPrefixFormat(t *testing.T) {
	configs := GetRateLimitConfigs()

	// Key prefixes should follow a consistent format
	assert.Contains(t, configs.Strict.KeyPrefix, "rate_limit")
	assert.Contains(t, configs.Moderate.KeyPrefix, "rate_limit")
	assert.Contains(t, configs.Lenient.KeyPrefix, "rate_limit")
}

func TestRateLimitConfigs_StrictValues(t *testing.T) {
	configs := GetRateLimitConfigs()

	// Test exact values for Strict config (to catch unintended changes)
	strictConfig := configs.Strict
	
	assert.Equal(t, 5, strictConfig.RequestsPerWindow,
		"Strict requests per window should be 5")
	assert.Equal(t, 1*time.Minute, strictConfig.WindowDuration,
		"Strict window duration should be 1 minute")
	assert.Equal(t, "rate_limit_strict", strictConfig.KeyPrefix,
		"Strict key prefix should be 'rate_limit_strict'")
}

func TestRateLimitConfigs_ModerateValues(t *testing.T) {
	configs := GetRateLimitConfigs()

	// Test exact values for Moderate config
	moderateConfig := configs.Moderate
	
	assert.Equal(t, 100, moderateConfig.RequestsPerWindow,
		"Moderate requests per window should be 100")
	assert.Equal(t, 1*time.Minute, moderateConfig.WindowDuration,
		"Moderate window duration should be 1 minute")
	assert.Equal(t, "rate_limit_moderate", moderateConfig.KeyPrefix,
		"Moderate key prefix should be 'rate_limit_moderate'")
}

func TestRateLimitConfigs_LenientValues(t *testing.T) {
	configs := GetRateLimitConfigs()

	// Test exact values for Lenient config
	lenientConfig := configs.Lenient
	
	assert.Equal(t, 1000, lenientConfig.RequestsPerWindow,
		"Lenient requests per window should be 1000")
	assert.Equal(t, 1*time.Minute, lenientConfig.WindowDuration,
		"Lenient window duration should be 1 minute")
	assert.Equal(t, "rate_limit_lenient", lenientConfig.KeyPrefix,
		"Lenient key prefix should be 'rate_limit_lenient'")
}

func TestRateLimitConfigs_StructureComplete(t *testing.T) {
	configs := GetRateLimitConfigs()

	// Ensure all three tiers are present
	assert.NotNil(t, configs.Strict)
	assert.NotNil(t, configs.Moderate)
	assert.NotNil(t, configs.Lenient)

	// Ensure no zero values
	assert.NotZero(t, configs.Strict.RequestsPerWindow)
	assert.NotZero(t, configs.Strict.WindowDuration)
	assert.NotZero(t, configs.Strict.KeyPrefix)

	assert.NotZero(t, configs.Moderate.RequestsPerWindow)
	assert.NotZero(t, configs.Moderate.WindowDuration)
	assert.NotZero(t, configs.Moderate.KeyPrefix)

	assert.NotZero(t, configs.Lenient.RequestsPerWindow)
	assert.NotZero(t, configs.Lenient.WindowDuration)
	assert.NotZero(t, configs.Lenient.KeyPrefix)
}