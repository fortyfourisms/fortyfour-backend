package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ═══════════════════════════════════════════════════════════════════════════
// TEST: RabbitMQConfig.GetURL
// ═══════════════════════════════════════════════════════════════════════════

func TestRabbitMQConfig_GetURL_DefaultVhost(t *testing.T) {
	cfg := RabbitMQConfig{
		Host:     "localhost",
		Port:     "5672",
		User:     "guest",
		Password: "guest",
		Vhost:    "/",
	}

	expected := "amqp://guest:guest@localhost:5672/"
	assert.Equal(t, expected, cfg.GetURL())
}

func TestRabbitMQConfig_GetURL_CustomVhost(t *testing.T) {
	cfg := RabbitMQConfig{
		Host:     "mq.example.com",
		Port:     "5672",
		User:     "admin",
		Password: "secret",
		Vhost:    "/myapp",
	}

	expected := "amqp://admin:secret@mq.example.com:5672/myapp"
	assert.Equal(t, expected, cfg.GetURL())
}

func TestRabbitMQConfig_GetURL_VhostWithoutSlash(t *testing.T) {
	cfg := RabbitMQConfig{
		Host:     "localhost",
		Port:     "5672",
		User:     "user",
		Password: "pass",
		Vhost:    "myvhost",
	}

	expected := "amqp://user:pass@localhost:5672/myvhost"
	assert.Equal(t, expected, cfg.GetURL())
}

func TestRabbitMQConfig_GetURL_EmptyVhost(t *testing.T) {
	cfg := RabbitMQConfig{
		Host:     "localhost",
		Port:     "5672",
		User:     "guest",
		Password: "guest",
		Vhost:    "",
	}

	expected := "amqp://guest:guest@localhost:5672/"
	assert.Equal(t, expected, cfg.GetURL())
}

func TestRabbitMQConfig_GetURL_EmptyFields(t *testing.T) {
	cfg := RabbitMQConfig{}

	// Should not panic even with all empty fields
	url := cfg.GetURL()
	assert.Contains(t, url, "amqp://")
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: Load() — RabbitMQ env vars
// ═══════════════════════════════════════════════════════════════════════════

func TestLoad_RabbitMQDefaults(t *testing.T) {
	clearEnvVars()
	clearRabbitMQEnvVars()

	cfg := Load()

	assert.Equal(t, "localhost", cfg.RabbitMQ.Host)
	assert.Equal(t, "5672", cfg.RabbitMQ.Port)
	assert.Equal(t, "guest", cfg.RabbitMQ.User)
	assert.Equal(t, "guest", cfg.RabbitMQ.Password)
	assert.Equal(t, "/", cfg.RabbitMQ.Vhost)
}

func TestLoad_RabbitMQFromEnv(t *testing.T) {
	clearEnvVars()
	clearRabbitMQEnvVars()

	os.Setenv("RABBITMQ_HOST", "mq.prod.com")
	os.Setenv("RABBITMQ_PORT", "5673")
	os.Setenv("RABBITMQ_USER", "produser")
	os.Setenv("RABBITMQ_PASSWORD", "prodpass")
	os.Setenv("RABBITMQ_VHOST", "/production")
	defer clearRabbitMQEnvVars()

	cfg := Load()

	assert.Equal(t, "mq.prod.com", cfg.RabbitMQ.Host)
	assert.Equal(t, "5673", cfg.RabbitMQ.Port)
	assert.Equal(t, "produser", cfg.RabbitMQ.User)
	assert.Equal(t, "prodpass", cfg.RabbitMQ.Password)
	assert.Equal(t, "/production", cfg.RabbitMQ.Vhost)
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: Load() — Gemini, InternalGateway, Turnstile env vars
// ═══════════════════════════════════════════════════════════════════════════

func TestLoad_GeminiAPIKeyDefault(t *testing.T) {
	clearEnvVars()
	os.Unsetenv("GEMINI_API_KEY")

	cfg := Load()
	assert.Equal(t, "", cfg.GeminiAPIKey)
}

func TestLoad_GeminiAPIKeyFromEnv(t *testing.T) {
	clearEnvVars()
	os.Setenv("GEMINI_API_KEY", "test-gemini-key")
	defer os.Unsetenv("GEMINI_API_KEY")

	cfg := Load()
	assert.Equal(t, "test-gemini-key", cfg.GeminiAPIKey)
}

func TestLoad_InternalGatewayKeyDefault(t *testing.T) {
	clearEnvVars()
	os.Unsetenv("INTERNAL_GATEWAY_KEY")

	cfg := Load()
	assert.Equal(t, "", cfg.InternalGatewayKey)
}

func TestLoad_InternalGatewayKeyFromEnv(t *testing.T) {
	clearEnvVars()
	os.Setenv("INTERNAL_GATEWAY_KEY", "gw-secret")
	defer os.Unsetenv("INTERNAL_GATEWAY_KEY")

	cfg := Load()
	assert.Equal(t, "gw-secret", cfg.InternalGatewayKey)
}

func TestLoad_TurnstileSecretKeyDefault(t *testing.T) {
	clearEnvVars()
	os.Unsetenv("TURNSTILE_SECRET_KEY")

	cfg := Load()
	assert.Equal(t, "", cfg.TurnstileSecretKey)
}

func TestLoad_TurnstileSecretKeyFromEnv(t *testing.T) {
	clearEnvVars()
	os.Setenv("TURNSTILE_SECRET_KEY", "ts-key-123")
	defer os.Unsetenv("TURNSTILE_SECRET_KEY")

	cfg := Load()
	assert.Equal(t, "ts-key-123", cfg.TurnstileSecretKey)
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: RabbitMQConfig.GetURL — end-to-end with Load()
// ═══════════════════════════════════════════════════════════════════════════

func TestLoad_RabbitMQGetURL_Integration(t *testing.T) {
	clearEnvVars()
	clearRabbitMQEnvVars()

	cfg := Load()

	url := cfg.RabbitMQ.GetURL()
	assert.Equal(t, "amqp://guest:guest@localhost:5672/", url)
}

// ═══════════════════════════════════════════════════════════════════════════
// HELPER
// ═══════════════════════════════════════════════════════════════════════════

func clearRabbitMQEnvVars() {
	vars := []string{"RABBITMQ_HOST", "RABBITMQ_PORT", "RABBITMQ_USER", "RABBITMQ_PASSWORD", "RABBITMQ_VHOST"}
	for _, key := range vars {
		os.Unsetenv(key)
	}
	os.Unsetenv("GEMINI_API_KEY")
	os.Unsetenv("INTERNAL_GATEWAY_KEY")
	os.Unsetenv("TURNSTILE_SECRET_KEY")
}
