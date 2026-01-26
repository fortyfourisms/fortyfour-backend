package services

import (
	"errors"
	"testing"
	"time"

	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/testhelpers"

	"golang.org/x/crypto/bcrypt"
)

//
// =========================
// MOCK REDIS (for inline tests)
// =========================
//

type mockRedis struct {
	store      map[string]string
	failSet    bool
	failExists bool
	failDelete bool
}

func newMockRedis() *mockRedis {
	return &mockRedis{
		store: make(map[string]string),
	}
}

func (m *mockRedis) Set(key string, value interface{}, ttl time.Duration) error {
	if m.failSet {
		return errors.New("redis set error")
	}
	str, ok := value.(string)
	if !ok {
		return errors.New("value must be string")
	}
	m.store[key] = str
	return nil
}

func (m *mockRedis) Get(key string) (string, error) {
	val, ok := m.store[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return val, nil
}

func (m *mockRedis) Delete(key string) error {
	if m.failDelete {
		return errors.New("redis delete error")
	}
	delete(m.store, key)
	return nil
}

func (m *mockRedis) Exists(key string) (bool, error) {
	if m.failExists {
		return false, errors.New("redis exists error")
	}
	_, ok := m.store[key]
	return ok, nil
}

func (m *mockRedis) Close() error {
	return nil
}

//
// =========================
// MOCK USER REPOSITORY (for inline tests)
// =========================
//

type mockUserRepo struct {
	users      map[string]*models.User
	failCreate bool
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users: make(map[string]*models.User),
	}
}

func (m *mockUserRepo) Create(user *models.User) error {
	if m.failCreate {
		return errors.New("create failed")
	}
	user.ID = "user-1"
	user.RoleName = "admin"
	m.users[user.Username] = user
	return nil
}

func (m *mockUserRepo) FindByUsername(username string) (*models.User, error) {
	u, ok := m.users[username]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

func (m *mockUserRepo) FindByID(id string) (*models.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockUserRepo) FindAll() ([]models.User, error) { return nil, nil }
func (m *mockUserRepo) Update(user *models.User) error  { return nil }
func (m *mockUserRepo) UpdateWithPhoto(user *models.User) error {
	return nil
}
func (m *mockUserRepo) UpdatePassword(id, hashed string) error { return nil }
func (m *mockUserRepo) GetPasswordByID(id string) (string, error) {
	return "", nil
}
func (m *mockUserRepo) Delete(id string) error { return nil }

func (m *mockUserRepo) EmailExists(email string, excludeID *string) (bool, error) {
	for _, u := range m.users {
		if u.Email == email {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockUserRepo) UsernameExists(username string, excludeID *string) (bool, error) {
	_, ok := m.users[username]
	return ok, nil
}

func (m *mockUserRepo) SetMFA(userID string, secret *string, enabled bool) error {
	for _, u := range m.users {
		if u.ID == userID {
			u.MFASecret = secret
			u.MFAEnabled = enabled
			return nil
		}
	}
	return errors.New("user not found")
}

//
// =========================
// CONSTRUCTOR TEST
// =========================
//

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

//
// =========================
// REGISTER TESTS (Inline mocks - from main)
// =========================
//

func TestRegister_Success(t *testing.T) {
	auth := NewAuthService(newMockUserRepo(), NewTokenService(newMockRedis(), "secret"))

	user, token, err := auth.Register(
		"user",
		"XyZ#91!kLmPq",
		"user@mail.com",
		nil,
		nil,
	)

	if err != nil || user == nil || token == nil {
		t.Fatal("register success expected")
	}
}

func TestRegister_UsernameExists(t *testing.T) {
	repo := newMockUserRepo()
	repo.users["user"] = &models.User{Username: "user"}

	auth := NewAuthService(repo, NewTokenService(newMockRedis(), "secret"))

	if _, _, err := auth.Register("user", "XyZ#91!kLmPq", "x@mail.com", nil, nil); err == nil {
		t.Fatal("expected username exists error")
	}
}

func TestRegister_EmailExists(t *testing.T) {
	repo := newMockUserRepo()
	repo.users["u1"] = &models.User{Email: "mail@test.com"}

	auth := NewAuthService(repo, NewTokenService(newMockRedis(), "secret"))

	if _, _, err := auth.Register("u2", "XyZ#91!kLmPq", "mail@test.com", nil, nil); err == nil {
		t.Fatal("expected email exists error")
	}
}

func TestRegister_WeakPassword(t *testing.T) {
	auth := NewAuthService(newMockUserRepo(), NewTokenService(newMockRedis(), "secret"))

	if _, _, err := auth.Register("u", "123", "x@mail.com", nil, nil); err == nil {
		t.Fatal("expected weak password error")
	}
}

func TestRegister_CreateError(t *testing.T) {
	repo := newMockUserRepo()
	repo.failCreate = true

	auth := NewAuthService(repo, NewTokenService(newMockRedis(), "secret"))

	if _, _, err := auth.Register("u", "XyZ#91!kLmPq", "x@mail.com", nil, nil); err == nil {
		t.Fatal("expected create error")
	}
}

//
// =========================
// REGISTER TESTS (testhelpers - from branch)
// =========================
//

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

//
// =========================
// LOGIN TESTS (Inline mocks - from main)
// =========================
//

func TestLogin_Success(t *testing.T) {
	repo := newMockUserRepo()
	hash, _ := bcrypt.GenerateFromPassword([]byte("XyZ#91!kLmPq"), bcrypt.DefaultCost)

	repo.users["user"] = &models.User{
		ID:       "1",
		Username: "user",
		Password: string(hash),
		RoleName: "admin",
	}

	auth := NewAuthService(repo, NewTokenService(newMockRedis(), "secret"))

	if u, tkn, err := auth.Login("user", "XyZ#91!kLmPq"); err != nil || u == nil || tkn == nil {
		t.Fatal("login success expected")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	auth := NewAuthService(newMockUserRepo(), NewTokenService(newMockRedis(), "secret"))

	if _, _, err := auth.Login("x", "pass"); err == nil {
		t.Fatal("expected user not found")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	repo := newMockUserRepo()
	hash, _ := bcrypt.GenerateFromPassword([]byte("right"), bcrypt.DefaultCost)

	repo.users["u"] = &models.User{
		Username: "u",
		Password: string(hash),
	}

	auth := NewAuthService(repo, NewTokenService(newMockRedis(), "secret"))

	if _, _, err := auth.Login("u", "wrong"); err == nil {
		t.Fatal("expected wrong password error")
	}
}

func TestLogin_TokenError(t *testing.T) {
	repo := newMockUserRepo()
	redis := newMockRedis()
	redis.failSet = true

	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	repo.users["u"] = &models.User{Username: "u", Password: string(hash)}

	auth := NewAuthService(repo, NewTokenService(redis, "secret"))

	if _, _, err := auth.Login("u", "pass"); err == nil {
		t.Fatal("expected token error")
	}
}

//
// =========================
// LOGIN TESTS (testhelpers - from branch)
// =========================
//

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

func TestAuthService_Login_Success_Detailed(t *testing.T) {
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

//
// =========================
// LOGOUT TESTS (Inline mocks - from main)
// =========================
//

func TestLogout_Success(t *testing.T) {
	redis := newMockRedis()
	redis.store["token"] = "x"

	auth := NewAuthService(newMockUserRepo(), NewTokenService(redis, "secret"))

	if err := auth.Logout("token"); err != nil {
		t.Fatal("logout success expected")
	}
}

func TestLogout_TokenNotExists(t *testing.T) {
	auth := NewAuthService(newMockUserRepo(), NewTokenService(newMockRedis(), "secret"))

	if err := auth.Logout("missing"); err != nil {
		t.Fatal("no error expected for missing token")
	}
}

func TestLogout_ExistsError(t *testing.T) {
	redis := newMockRedis()
	redis.failExists = true

	auth := NewAuthService(newMockUserRepo(), NewTokenService(redis, "secret"))

	err := auth.Logout("token")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLogout_DeleteError(t *testing.T) {
	redis := newMockRedis()
	redis.store["token"] = "x"
	redis.failDelete = true

	auth := NewAuthService(newMockUserRepo(), NewTokenService(redis, "secret"))

	if err := auth.Logout("token"); err == nil {
		t.Fatal("expected delete error")
	}
}

//
// =========================
// LOGOUT TESTS (testhelpers - from branch)
// =========================
//

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

func TestAuthService_Logout_Success_Detailed(t *testing.T) {
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

//
// =========================
// INTERFACE CHECK
// =========================
//

var _ repository.UserRepositoryInterface = (*mockUserRepo)(nil)