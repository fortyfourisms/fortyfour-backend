package services

import (
	"fortyfour-backend/internal/testhelpers"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register_Success(t *testing.T) {
	// Arrange
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret")
	authService := NewAuthService(userRepo, tokenService)

	// Act
	user, tokens, err := authService.Register("testuser", "password123", "test@example.com")

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("expected user to be created")
	}

	if user.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", user.Username)
	}

	if tokens == nil {
		t.Fatal("expected tokens to be returned")
	}

	if tokens.AccessToken == "" {
		t.Error("expected access token to be generated")
	}

	if tokens.RefreshToken == "" {
		t.Error("expected refresh token to be generated")
	}

	// Verify password is hashed
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
	if err != nil {
		t.Error("password should be hashed correctly")
	}
}

func TestAuthService_Register_DuplicateUsername(t *testing.T) {
	// Arrange
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret")
	authService := NewAuthService(userRepo, tokenService)

	authService.Register("testuser", "password123", "test@example.com")

	// Act
	_, _, err := authService.Register("testuser", "newpassword", "another@example.com")

	// Assert
	if err == nil {
		t.Fatal("expected error for duplicate username")
	}

	if err.Error() != "username already exists" {
		t.Errorf("expected 'username already exists' error, got '%s'", err.Error())
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	// Arrange
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret")
	authService := NewAuthService(userRepo, tokenService)

	authService.Register("testuser", "password123", "test@example.com")

	// Act
	user, tokens, err := authService.Login("testuser", "password123")

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("expected user to be returned")
	}

	if tokens == nil {
		t.Fatal("expected tokens to be returned")
	}

	if tokens.AccessToken == "" {
		t.Error("expected access token to be generated")
	}

	if tokens.RefreshToken == "" {
		t.Error("expected refresh token to be generated")
	}
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	// Arrange
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret")
	authService := NewAuthService(userRepo, tokenService)

	authService.Register("testuser", "password123", "test@example.com")

	// Act
	_, _, err := authService.Login("testuser", "wrongpassword")

	// Assert
	if err == nil {
		t.Fatal("expected error for invalid password")
	}

	if err.Error() != "invalid credentials" {
		t.Errorf("expected 'invalid credentials', got '%s'", err.Error())
	}
}

func TestAuthService_Logout_Success(t *testing.T) {
	// Arrange
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret")
	authService := NewAuthService(userRepo, tokenService)

	// Register and get tokens
	_, tokens, _ := authService.Register("testuser", "password123", "test@example.com")

	// Act
	err := authService.Logout(tokens.RefreshToken)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify token is revoked
	_, err = tokenService.RefreshAccessToken(tokens.RefreshToken)
	if err == nil {
		t.Error("expected error when using revoked token")
	}
}
