package services

import (
	"fortyfour-backend/internal/testhelpers"
	"testing"
	"time"
)

func TestTokenService_GenerateTokenPair_Success(t *testing.T) {
	// Arrange
	redis := testhelpers.NewMockRedisClient()
	service := NewTokenService(redis, "test-secret")

	// Act
	tokens, err := service.GenerateTokenPair(1, "testuser")

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if tokens.AccessToken == "" {
		t.Error("expected access token to be generated")
	}

	if tokens.RefreshToken == "" {
		t.Error("expected refresh token to be generated")
	}

	if tokens.ExpiresAt.IsZero() {
		t.Error("expected expires_at to be set")
	}

	// Verify refresh token is stored in Redis
	key := "refresh_token:" + tokens.RefreshToken
	exists, err := redis.Exists(key)
	if err != nil {
		t.Fatalf("error checking Redis: %v", err)
	}
	if !exists {
		t.Error("expected refresh token to be stored in Redis")
	}
}

func TestTokenService_RefreshAccessToken_Success(t *testing.T) {
	// Arrange
	redis := testhelpers.NewMockRedisClient()
	service := NewTokenService(redis, "test-secret")

	// Generate initial token pair
	initialTokens, _ := service.GenerateTokenPair(1, "testuser")

	// Wait a moment to ensure new token is different
	time.Sleep(10 * time.Millisecond)

	// Act
	newTokens, err := service.RefreshAccessToken(initialTokens.RefreshToken)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if newTokens.AccessToken == initialTokens.AccessToken {
		t.Error("expected new access token to be different")
	}

	if newTokens.RefreshToken == initialTokens.RefreshToken {
		t.Error("expected new refresh token to be different (token rotation)")
	}
}

func TestTokenService_RefreshAccessToken_InvalidToken(t *testing.T) {
	// Arrange
	redis := testhelpers.NewMockRedisClient()
	service := NewTokenService(redis, "test-secret")

	// Act
	_, err := service.RefreshAccessToken("invalid-token")

	// Assert
	if err == nil {
		t.Fatal("expected error for invalid refresh token")
	}

	if err.Error() != "invalid or expired refresh token" {
		t.Errorf("expected 'invalid or expired refresh token', got '%s'", err.Error())
	}
}

func TestTokenService_RevokeRefreshToken_Success(t *testing.T) {
	// Arrange
	redis := testhelpers.NewMockRedisClient()
	service := NewTokenService(redis, "test-secret")

	// Generate token pair
	tokens, _ := service.GenerateTokenPair(1, "testuser")

	// Act
	err := service.RevokeRefreshToken(tokens.RefreshToken)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify token is removed from Redis
	key := "refresh_token:" + tokens.RefreshToken
	exists, _ := redis.Exists(key)
	if exists {
		t.Error("expected refresh token to be removed from Redis")
	}

	// Try to use revoked token
	_, err = service.RefreshAccessToken(tokens.RefreshToken)
	if err == nil {
		t.Error("expected error when using revoked token")
	}
}
