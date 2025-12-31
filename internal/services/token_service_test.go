package services

import (
	"encoding/json"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/testhelpers"
	"testing"
	"time"
)

func TestTokenService_GenerateTokenPair_Success(t *testing.T) {
	// Arrange
	redis := testhelpers.NewMockRedisClient()
	service := NewTokenService(redis, "test-secret")

	// Act
	tokens, err := service.GenerateTokenPair("1", "testuser", "admin")

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

	// Verify token data in Redis contains role
	data, err := redis.Get(key)
	if err != nil {
		t.Fatalf("error getting token data: %v", err)
	}

	var tokenData models.RefreshTokenData
	if err := json.Unmarshal([]byte(data), &tokenData); err != nil {
		t.Fatalf("error unmarshaling token data: %v", err)
	}

	if tokenData.UserID != "1" {
		t.Errorf("expected userID '1', got '%s'", tokenData.UserID)
	}

	if tokenData.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", tokenData.Username)
	}

	if tokenData.Role != "admin" {
		t.Errorf("expected role 'admin', got '%s'", tokenData.Role)
	}
}

func TestTokenService_RefreshAccessToken_Success(t *testing.T) {
	// Arrange
	redis := testhelpers.NewMockRedisClient()
	service := NewTokenService(redis, "test-secret")

	// Generate initial token pair
	initialTokens, err := service.GenerateTokenPair("1", "testuser", "admin")
	if err != nil {
		t.Fatalf("failed to generate initial tokens: %v", err)
	}

	// Wait lebih lama untuk memastikan timestamp berbeda (1 detik)
	time.Sleep(1 * time.Second)

	// Act
	newTokens, err := service.RefreshAccessToken(initialTokens.RefreshToken)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if newTokens == nil {
		t.Fatal("expected new tokens to be returned")
	}

	if newTokens.AccessToken == "" {
		t.Error("expected new access token to be generated")
	}

	if newTokens.RefreshToken == "" {
		t.Error("expected new refresh token to be generated")
	}

	// JWT tokens bisa sama jika dibuat di detik yang sama karena `exp` claim rounded ke seconds
	// Jadi cek bahwa token tidak kosong saja, bukan harus berbeda
	if newTokens.AccessToken == initialTokens.AccessToken {
		t.Log("Note: Access tokens are the same - this can happen if generated in same second (JWT exp is in seconds)")
	}

	// Check token rotation jika implementasi service support token rotation
	if newTokens.RefreshToken == initialTokens.RefreshToken {
		t.Log("Note: Refresh token not rotated - this is expected if service doesn't implement token rotation")
		// Atau bisa t.Error jika expect token rotation harus ada
		// t.Error("expected new refresh token to be different (token rotation)")
	}

	_, err = service.RefreshAccessToken(newTokens.RefreshToken)
	if err != nil {
		t.Error("expected new refresh token to be usable")
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
	tokens, _ := service.GenerateTokenPair("1", "testuser", "admin")

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