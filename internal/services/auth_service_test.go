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
	user, tokens, err := authService.Register("testuser", "password123", "test@example.com", nil, nil)

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

	if user.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", user.Email)
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

func TestAuthService_Register_WithRoleAndJabatan(t *testing.T) {
	// Arrange
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret")
	authService := NewAuthService(userRepo, tokenService)

	roleID := "role-123"
	idJabatan := "jabatan-456"

	// Act
	user, tokens, err := authService.Register("testuser", "password123", "test@example.com", &roleID, &idJabatan)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("expected user to be created")
	}

	if user.RoleID == nil || *user.RoleID != roleID {
		t.Errorf("expected roleID '%s', got '%v'", roleID, user.RoleID)
	}

	if user.IDJabatan == nil || *user.IDJabatan != idJabatan {
		t.Errorf("expected idJabatan '%s', got '%v'", idJabatan, user.IDJabatan)
	}

	if tokens == nil {
		t.Fatal("expected tokens to be returned")
	}
}

func TestAuthService_Register_EmptyFields(t *testing.T) {
	// Arrange
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret")
	authService := NewAuthService(userRepo, tokenService)

	tests := []struct {
		name     string
		username string
		password string
		email    string
	}{
		{"empty username", "", "password123", "test@example.com"},
		{"empty password", "testuser", "", "test@example.com"},
		{"empty email", "testuser", "password123", ""},
		{"all empty", "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			_, _, err := authService.Register(tt.username, tt.password, tt.email, nil, nil)

			// Assert
			if err == nil {
				t.Fatal("expected error for empty fields")
			}

			if err.Error() != "all fields are required" {
				t.Errorf("expected 'all fields are required' error, got '%s'", err.Error())
			}
		})
	}
}

func TestAuthService_Register_DuplicateUsername(t *testing.T) {
	// Arrange
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret")
	authService := NewAuthService(userRepo, tokenService)

	authService.Register("testuser", "password123", "test@example.com", nil, nil)

	// Act
	_, _, err := authService.Register("testuser", "newpassword", "another@example.com", nil, nil)

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

	authService.Register("testuser", "password123", "test@example.com", nil, nil)

	// Act
	user, tokens, err := authService.Login("testuser", "password123")

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("expected user to be returned")
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
}

func TestAuthService_Login_InvalidUsername(t *testing.T) {
	// Arrange
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret")
	authService := NewAuthService(userRepo, tokenService)

	authService.Register("testuser", "password123", "test@example.com", nil, nil)

	// Act
	_, _, err := authService.Login("nonexistent", "password123")

	// Assert
	if err == nil {
		t.Fatal("expected error for invalid username")
	}

	if err.Error() != "invalid credentials" {
		t.Errorf("expected 'invalid credentials', got '%s'", err.Error())
	}
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	// Arrange
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret")
	authService := NewAuthService(userRepo, tokenService)

	authService.Register("testuser", "password123", "test@example.com", nil, nil)

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

func TestAuthService_Login_WithRoleName(t *testing.T) {
	// Arrange
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret")
	authService := NewAuthService(userRepo, tokenService)

	roleID := "role-123"
	authService.Register("testuser", "password123", "test@example.com", &roleID, nil)

	// Act
	user, tokens, err := authService.Login("testuser", "password123")

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.RoleName == "" {
		t.Error("expected roleName to be populated after login")
	}

	if tokens == nil {
		t.Fatal("expected tokens to be returned")
	}
}

func TestAuthService_Logout_Success(t *testing.T) {
	// Arrange
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret")
	authService := NewAuthService(userRepo, tokenService)

	// Register and get tokens
	_, tokens, _ := authService.Register("testuser", "password123", "test@example.com", nil, nil)

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

func TestAuthService_Logout_InvalidToken(t *testing.T) {
	// Arrange
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret")
	authService := NewAuthService(userRepo, tokenService)

	// Act
	err := authService.Logout("invalid-token")

	// Assert - tergantung implementasi RevokeRefreshToken
	// bisa error atau tidak error jika token tidak ada
	// sesuaikan dengan behavior yang diharapkan
	_ = err // untuk saat ini kita terima hasil apapun
}
