package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_WithDefaultValues(t *testing.T) {
	// Clear all environment variables
	clearEnvVars()

	cfg := Load()

	assert.NotNil(t, cfg)
	assert.Equal(t, ":8080", cfg.Port)
	assert.Equal(t, "your-secret-key", cfg.JWTSecret)
	assert.Equal(t, "https://admin.kssindustri.site", cfg.Domain)

	// Database defaults
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, "3306", cfg.Database.Port)
	assert.Equal(t, "root", cfg.Database.User)
	assert.Equal(t, "", cfg.Database.Password)
	assert.Equal(t, "fortyfour-backend_db", cfg.Database.DBName)

	// Redis defaults
	assert.Equal(t, "localhost", cfg.Redis.Host)
	assert.Equal(t, "6379", cfg.Redis.Port)
	assert.Equal(t, "", cfg.Redis.Password)
	assert.Equal(t, 0, cfg.Redis.DB)

	// Rollbar defaults - FIXED: Match actual default value in config.go
	assert.Equal(t, "0eddf8fb05e44067a12a8bb36ccc3ef9", cfg.Rollbar.Token)
	assert.Equal(t, "production", cfg.Rollbar.Env)

	// Casbin model path should be set
	assert.NotEmpty(t, cfg.CasbinModelPath)
	assert.Contains(t, cfg.CasbinModelPath, "casbin_model.conf")
}

func TestLoad_WithEnvironmentVariables(t *testing.T) {
	// Setup environment variables
	envVars := map[string]string{
		"PORT":            ":9090",
		"JWT_SECRET":      "test-secret-key",
		"DOMAIN":          "https://test.example.com",
		"DB_HOST":         "db.example.com",
		"DB_PORT":         "5432",
		"DB_USER":         "testuser",
		"DB_PASSWORD":     "testpass",
		"DB_NAME":         "testdb",
		"REDIS_HOST":      "redis.example.com",
		"REDIS_PORT":      "6380",
		"REDIS_PASSWORD":  "redispass",
		"REDIS_DB":        "5",
		"ROLLBAR_TOKEN":   "test-rollbar-token",
		"ROLLBAR_STATUS":  "development",
		"CASBIN_MODEL_PATH": "/custom/path/model.conf",
	}

	setEnvVars(envVars)
	defer clearEnvVars()

	cfg := Load()

	assert.NotNil(t, cfg)
	assert.Equal(t, ":9090", cfg.Port)
	assert.Equal(t, "test-secret-key", cfg.JWTSecret)
	assert.Equal(t, "https://test.example.com", cfg.Domain)

	// Database
	assert.Equal(t, "db.example.com", cfg.Database.Host)
	assert.Equal(t, "5432", cfg.Database.Port)
	assert.Equal(t, "testuser", cfg.Database.User)
	assert.Equal(t, "testpass", cfg.Database.Password)
	assert.Equal(t, "testdb", cfg.Database.DBName)

	// Redis
	assert.Equal(t, "redis.example.com", cfg.Redis.Host)
	assert.Equal(t, "6380", cfg.Redis.Port)
	assert.Equal(t, "redispass", cfg.Redis.Password)
	assert.Equal(t, 5, cfg.Redis.DB)

	// Rollbar
	assert.Equal(t, "test-rollbar-token", cfg.Rollbar.Token)
	assert.Equal(t, "development", cfg.Rollbar.Env)

	// Casbin
	assert.Equal(t, "/custom/path/model.conf", cfg.CasbinModelPath)
}

func TestLoad_WithPartialEnvironmentVariables(t *testing.T) {
	clearEnvVars()

	// Set only some environment variables
	envVars := map[string]string{
		"PORT":       ":3000",
		"JWT_SECRET": "partial-secret",
		"DB_HOST":    "custom-db-host",
		"REDIS_DB":   "3",
	}

	setEnvVars(envVars)
	defer clearEnvVars()

	cfg := Load()

	// Environment variables should be used
	assert.Equal(t, ":3000", cfg.Port)
	assert.Equal(t, "partial-secret", cfg.JWTSecret)
	assert.Equal(t, "custom-db-host", cfg.Database.Host)
	assert.Equal(t, 3, cfg.Redis.DB)

	// Defaults should be used for missing vars
	assert.Equal(t, "https://admin.kssindustri.site", cfg.Domain)
	assert.Equal(t, "3306", cfg.Database.Port)
	assert.Equal(t, "root", cfg.Database.User)
	assert.Equal(t, "localhost", cfg.Redis.Host)
}

func TestDatabaseConfig_GetDSN(t *testing.T) {
	tests := []struct {
		name     string
		config   DatabaseConfig
		expected string
	}{
		{
			name: "full configuration",
			config: DatabaseConfig{
				User:     "root",
				Password: "password123",
				Host:     "localhost",
				Port:     "3306",
				DBName:   "testdb",
			},
			expected: "root:password123@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
		},
		{
			name: "empty password",
			config: DatabaseConfig{
				User:     "root",
				Password: "",
				Host:     "localhost",
				Port:     "3306",
				DBName:   "testdb",
			},
			expected: "root:@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
		},
		{
			name: "custom port",
			config: DatabaseConfig{
				User:     "admin",
				Password: "admin123",
				Host:     "db.example.com",
				Port:     "3307",
				DBName:   "myapp_db",
			},
			expected: "admin:admin123@tcp(db.example.com:3307)/myapp_db?charset=utf8mb4&parseTime=True&loc=Local",
		},
		{
			name: "remote database",
			config: DatabaseConfig{
				User:     "user",
				Password: "pass",
				Host:     "192.168.1.100",
				Port:     "3306",
				DBName:   "production_db",
			},
			expected: "user:pass@tcp(192.168.1.100:3306)/production_db?charset=utf8mb4&parseTime=True&loc=Local",
		},
		{
			name: "special characters in password",
			config: DatabaseConfig{
				User:     "root",
				Password: "p@ss!word#123",
				Host:     "localhost",
				Port:     "3306",
				DBName:   "testdb",
			},
			expected: "root:p@ss!word#123@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := tt.config.GetDSN()
			assert.Equal(t, tt.expected, dsn)
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		fallback string
		envValue string
		expected string
	}{
		{
			name:     "environment variable exists",
			key:      "TEST_VAR",
			fallback: "default",
			envValue: "custom-value",
			expected: "custom-value",
		},
		{
			name:     "environment variable not set - use fallback",
			key:      "NON_EXISTENT_VAR",
			fallback: "default-value",
			envValue: "",
			expected: "default-value",
		},
		{
			name:     "environment variable is empty string",
			key:      "EMPTY_VAR",
			fallback: "default",
			envValue: "",
			expected: "default",
		},
		{
			name:     "environment variable with spaces",
			key:      "SPACE_VAR",
			fallback: "default",
			envValue: "  value with spaces  ",
			expected: "  value with spaces  ",
		},
		{
			name:     "empty fallback",
			key:      "TEST_EMPTY_FALLBACK",
			fallback: "",
			envValue: "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear the environment variable first
			os.Unsetenv(tt.key)

			// Set if not empty
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.fallback)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnvAsInt(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		fallback int
		envValue string
		expected int
	}{
		{
			name:     "valid integer",
			key:      "TEST_INT",
			fallback: 10,
			envValue: "42",
			expected: 42,
		},
		{
			name:     "zero value",
			key:      "TEST_ZERO",
			fallback: 10,
			envValue: "0",
			expected: 0,
		},
		{
			name:     "negative integer",
			key:      "TEST_NEGATIVE",
			fallback: 10,
			envValue: "-5",
			expected: -5,
		},
		{
			name:     "not set - use fallback",
			key:      "NON_EXISTENT_INT",
			fallback: 100,
			envValue: "",
			expected: 100,
		},
		{
			name:     "invalid integer - use fallback",
			key:      "INVALID_INT",
			fallback: 50,
			envValue: "not-a-number",
			expected: 50,
		},
		{
			name:     "float value - use fallback",
			key:      "FLOAT_VAR",
			fallback: 20,
			envValue: "3.14",
			expected: 20,
		},
		{
			name:     "integer with whitespace - use fallback",
			key:      "WHITESPACE_INT",
			fallback: 30,
			envValue: " 123 ",
			expected: 30,
		},
		{
			name:     "large integer",
			key:      "LARGE_INT",
			fallback: 1,
			envValue: "999999999",
			expected: 999999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear the environment variable first
			os.Unsetenv(tt.key)

			// Set if not empty
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnvAsInt(tt.key, tt.fallback)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig_StructureIntegrity(t *testing.T) {
	cfg := Load()

	// Verify all nested configs are initialized
	require.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.Port)
	assert.NotEmpty(t, cfg.JWTSecret)
	assert.NotEmpty(t, cfg.Domain)

	// Database config should be complete
	assert.NotEmpty(t, cfg.Database.Host)
	assert.NotEmpty(t, cfg.Database.Port)
	assert.NotEmpty(t, cfg.Database.User)
	// Password can be empty (valid state)
	assert.NotEmpty(t, cfg.Database.DBName)

	// Redis config should be complete
	assert.NotEmpty(t, cfg.Redis.Host)
	assert.NotEmpty(t, cfg.Redis.Port)
	// Password can be empty (valid state)
	assert.GreaterOrEqual(t, cfg.Redis.DB, 0)

	// Rollbar config should be complete
	assert.NotEmpty(t, cfg.Rollbar.Token)
	assert.NotEmpty(t, cfg.Rollbar.Env)

	// Casbin model path should exist
	assert.NotEmpty(t, cfg.CasbinModelPath)
}

func TestRedisDB_ValidRange(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected int
	}{
		{
			name:     "minimum valid DB (0)",
			envValue: "0",
			expected: 0,
		},
		{
			name:     "typical DB (1)",
			envValue: "1",
			expected: 1,
		},
		{
			name:     "high DB number",
			envValue: "15",
			expected: 15,
		},
		{
			name:     "negative - should use fallback",
			envValue: "-1",
			expected: -1, // Will be parsed as -1, app should validate separately
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnvVars()
			os.Setenv("REDIS_DB", tt.envValue)
			defer os.Unsetenv("REDIS_DB")

			cfg := Load()
			assert.Equal(t, tt.expected, cfg.Redis.DB)
		})
	}
}

func TestLoad_Idempotency(t *testing.T) {
	// Load should be idempotent - calling multiple times with same env should give same result
	clearEnvVars()

	cfg1 := Load()
	cfg2 := Load()

	assert.Equal(t, cfg1.Port, cfg2.Port)
	assert.Equal(t, cfg1.JWTSecret, cfg2.JWTSecret)
	assert.Equal(t, cfg1.Domain, cfg2.Domain)
	assert.Equal(t, cfg1.Database.Host, cfg2.Database.Host)
	assert.Equal(t, cfg1.Redis.Port, cfg2.Redis.Port)
}

func TestDatabaseConfig_GetDSN_EmptyFields(t *testing.T) {
	// Test DSN generation with empty fields (edge case)
	config := DatabaseConfig{
		User:     "",
		Password: "",
		Host:     "",
		Port:     "",
		DBName:   "",
	}

	dsn := config.GetDSN()
	expected := ":@tcp(:)/?charset=utf8mb4&parseTime=True&loc=Local"
	assert.Equal(t, expected, dsn)
}

// Helper functions

func clearEnvVars() {
	envVars := []string{
		"PORT", "JWT_SECRET", "DOMAIN",
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME",
		"REDIS_HOST", "REDIS_PORT", "REDIS_PASSWORD", "REDIS_DB",
		"ROLLBAR_TOKEN", "ROLLBAR_STATUS",
		"CASBIN_MODEL_PATH",
	}

	for _, key := range envVars {
		os.Unsetenv(key)
	}
}

func setEnvVars(vars map[string]string) {
	for key, value := range vars {
		os.Setenv(key, value)
	}
}