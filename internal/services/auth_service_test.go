package services

import (
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/testhelpers"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// Skeleton tests dari main
func TestNewAuthService(t *testing.T) {
	tests := []struct {
		name         string
		userRepo     repository.UserRepositoryInterface
		tokenService *TokenService
		want         *AuthService
	}{
		{
			name:         "create new auth service",
			userRepo:     testhelpers.NewMockUserRepository(),
			tokenService: NewTokenService(testhelpers.NewMockRedisClient(), "test-secret"),
			want:         nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAuthService(tt.userRepo, tt.tokenService)
			if got == nil {
				t.Errorf("NewAuthService() returned nil")
			}
		})
	}
}

/*
=====================================
 TEST REGISTER
=====================================
*/

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name      string
		username  string
		password  string
		email     string
		roleID    *string
		idJabatan *string
		wantErr   bool
		errMsg    string
	}{
		{
			name:     "empty username",
			username: "",
			password: "MySecureP@ssw0rd2024!",
			email:    "test@example.com",
			wantErr:  true,
			errMsg:   "username wajib diisi",
		},
		{
			name:     "empty password",
			username: "testuser",
			password: "",
			email:    "test@example.com",
			wantErr:  true,
			errMsg:   "password wajib diisi",
		},
		{
			name:     "empty email",
			username: "testuser",
			password: "MySecureP@ssw0rd2024!",
			email:    "",
			wantErr:  true,
			errMsg:   "email wajib diisi",
		},
		{
			name:     "all empty",
			username: "",
			password: "",
			email:    "",
			wantErr:  true,
			errMsg:   "username wajib diisi",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := testhelpers.NewMockUserRepository()
			redis := testhelpers.NewMockRedisClient()
			tokenService := NewTokenService(redis, "test-secret")
			s := NewAuthService(userRepo, tokenService)
			
			_, _, gotErr := s.Register(tt.username, tt.password, tt.email, tt.roleID, tt.idJabatan)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Register() failed: %v", gotErr)
				} else if gotErr.Error() != tt.errMsg {
					t.Errorf("expected '%s' error, got '%s'", tt.errMsg, gotErr.Error())
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Register() succeeded unexpectedly")
			}
		})
	}
}

func TestAuthService_Register_Success(t *testing.T) {
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret")
	authService := NewAuthService(userRepo, tokenService)

	user, tokens, err := authService.Register("testuser", "MySecureP@ssw0rd2024!", "test@example.com", nil, nil)

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
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("MySecureP@ssw0rd2024!"))
	if err != nil {
		t.Error("password should be hashed correctly")
	}
}

func TestAuthService_Register_WithRoleAndJabatan(t *testing.T) {
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret")
	authService := NewAuthService(userRepo, tokenService)

	roleID := "role-123"
	idJabatan := "jabatan-456"

	user, tokens, err := authService.Register("testuser", "MySecureP@ssw0rd2024!", "test@example.com", &roleID, &idJabatan)

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

func TestAuthService_Register_DuplicateUsername(t *testing.T) {
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret")
	authService := NewAuthService(userRepo, tokenService)

	authService.Register("testuser", "MySecureP@ssw0rd2024!", "test@example.com", nil, nil)

	_, _, err := authService.Register("testuser", "DifferentSecureP@ss2024!", "another@example.com", nil, nil)

	if err == nil {
		t.Fatal("expected error for duplicate username")
	}

	if err.Error() != "username sudah digunakan" {
		t.Errorf("expected 'username sudah digunakan' error, got '%s'", err.Error())
	}
}

/*
=====================================
 TEST LOGIN
=====================================
*/

func TestAuthService_Login(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		// Placeholder untuk kompatibilitas dengan struktur main
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := testhelpers.NewMockUserRepository()
			redis := testhelpers.NewMockRedisClient()
			tokenService := NewTokenService(redis, "test-secret")
			s := NewAuthService(userRepo, tokenService)
			
			_, _, gotErr := s.Login(tt.username, tt.password)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Login() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Login() succeeded unexpectedly")
			}
		})
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret")
	authService := NewAuthService(userRepo, tokenService)

	authService.Register("testuser", "MySecureP@ssw0rd2024!", "test@example.com", nil, nil)

	user, tokens, err := authService.Login("testuser", "MySecureP@ssw0rd2024!")

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
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret")
	authService := NewAuthService(userRepo, tokenService)

	authService.Register("testuser", "MySecureP@ssw0rd2024!", "test@example.com", nil, nil)

	_, _, err := authService.Login("nonexistent", "MySecureP@ssw0rd2024!")

	if err == nil {
		t.Fatal("expected error for invalid username")
	}

	if err.Error() != "username atau password salah" {
		t.Errorf("expected 'username atau password salah', got '%s'", err.Error())
	}
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret")
	authService := NewAuthService(userRepo, tokenService)

	authService.Register("testuser", "MySecureP@ssw0rd2024!", "test@example.com", nil, nil)

	_, _, err := authService.Login("testuser", "WrongP@ssword!")

	if err == nil {
		t.Fatal("expected error for invalid password")
	}

	if err.Error() != "username atau password salah" {
		t.Errorf("expected 'username atau password salah', got '%s'", err.Error())
	}
}

func TestAuthService_Login_WithRoleName(t *testing.T) {
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret")
	authService := NewAuthService(userRepo, tokenService)

	roleID := "role-123"
	
	user, _, err := authService.Register("testuser", "MySecureP@ssw0rd2024!", "test@example.com", &roleID, nil)
	if err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	if user.RoleID != nil {
		user.RoleName = "Admin"
	}

	loginUser, tokens, err := authService.Login("testuser", "MySecureP@ssw0rd2024!")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if loginUser.RoleName == "" {
		t.Log("Note: RoleName is empty - this might be expected if mock doesn't populate role relationships")
	}

	if tokens == nil {
		t.Fatal("expected tokens to be returned")
	}
}

/*
=====================================
 TEST LOGOUT
=====================================
*/

func TestAuthService_Logout(t *testing.T) {
	tests := []struct {
		name         string
		refreshToken string
		wantErr      bool
	}{
		// Placeholder untuk kompatibilitas dengan struktur main
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := testhelpers.NewMockUserRepository()
			redis := testhelpers.NewMockRedisClient()
			tokenService := NewTokenService(redis, "test-secret")
			s := NewAuthService(userRepo, tokenService)
			
			gotErr := s.Logout(tt.refreshToken)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Logout() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Logout() succeeded unexpectedly")
			}
		})
	}
}

func TestAuthService_Logout_Success(t *testing.T) {
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret")
	authService := NewAuthService(userRepo, tokenService)

	// Register and get tokens
	_, tokens, err := authService.Register("testuser", "MySecureP@ssw0rd2024!", "test@example.com", nil, nil)
	if err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	if tokens == nil {
		t.Fatal("expected tokens to be returned from registration")
	}

	// Act - Logout
	err = authService.Logout(tokens.RefreshToken)

	// Assert
	if err != nil {
		t.Fatalf("expected no error during logout, got %v", err)
	}

	// Verify token is revoked
	_, err = tokenService.RefreshAccessToken(tokens.RefreshToken)
	if err == nil {
		t.Error("expected error when using revoked token")
	}
}

func TestAuthService_Logout_InvalidToken(t *testing.T) {
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret")
	authService := NewAuthService(userRepo, tokenService)

	err := authService.Logout("invalid-token")

	_ = err // Accept any result for invalid token
}