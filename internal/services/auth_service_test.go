package services

import (
	"errors"
	"testing"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/testhelpers"

	"golang.org/x/crypto/bcrypt"
)

//
// =========================
// MOCK REDIS
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

func (m *mockRedis) Scan(pattern string) ([]string, error) {
	return []string{}, nil
}

//
// =========================
// MOCK USER REPOSITORY
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
	// Set PasswordChangedAt agar tidak langsung expired saat Login
	if user.PasswordChangedAt.IsZero() {
		user.PasswordChangedAt = time.Now()
	}
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

func (m *mockUserRepo) UpdateStatus(userID string, status models.UserStatus) error {
	if u, ok := m.users[userID]; ok {
		u.Status = status
	}
	return nil
}
func (m *mockUserRepo) IncrementLoginAttempts(userID string) (int, error) {
	if u, ok := m.users[userID]; ok {
		u.LoginAttempts++
		return u.LoginAttempts, nil
	}
	return 0, nil
}
func (m *mockUserRepo) ResetLoginAttempts(userID string) error {
	if u, ok := m.users[userID]; ok {
		u.LoginAttempts = 0
	}
	return nil
}
func (m *mockUserRepo) UpdatePasswordChangedAt(userID string) error { return nil }

func (m *mockUserRepo) ExistsByPerusahaan(idPerusahaan string) (bool, error) {
	for _, u := range m.users {
		if u.IDPerusahaan != nil && *u.IDPerusahaan == idPerusahaan {
			return true, nil
		}
	}
	return false, nil
}

// newNotifSvc membuat NotificationService dengan mock redis untuk test
func newNotifSvc() *NotificationService {
	return NewNotificationService(newMockRedis())
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
			tokenService: NewTokenService(testhelpers.NewMockRedisClient(), "test-secret", false, "localhost"),
			want:         nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAuthService(tt.userRepo, tt.tokenService, newNotifSvc())
			if got == nil {
				t.Errorf("NewAuthService() returned nil")
			}
		})
	}
}

//
// =========================
// REGISTER TESTS
// =========================
//

func TestRegister_Success(t *testing.T) {
	auth := NewAuthService(newMockUserRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

	user, token, err := auth.Register(
		dto.RegisterRequest{
			Username: "user",
			Password: "XyZ#91!kLmPq",
			Email:    "user@mail.com",
		},
		testhelpers.NewMockPerusahaanService(),
	)

	if err != nil || user == nil || token == nil {
		t.Fatal("register success expected")
	}
}

func TestRegister_UsernameExists(t *testing.T) {
	repo := newMockUserRepo()
	repo.users["user"] = &models.User{Username: "user"}

	auth := NewAuthService(repo, NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

	if _, _, err := auth.Register(
		dto.RegisterRequest{Username: "user", Password: "XyZ#91!kLmPq", Email: "x@mail.com"},
		testhelpers.NewMockPerusahaanService(),
	); err == nil {
		t.Fatal("expected username exists error")
	}
}

func TestRegister_EmailExists(t *testing.T) {
	repo := newMockUserRepo()
	repo.users["u1"] = &models.User{Email: "mail@test.com"}

	auth := NewAuthService(repo, NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

	if _, _, err := auth.Register(
		dto.RegisterRequest{Username: "u2", Password: "XyZ#91!kLmPq", Email: "mail@test.com"},
		testhelpers.NewMockPerusahaanService(),
	); err == nil {
		t.Fatal("expected email exists error")
	}
}

func TestRegister_WeakPassword(t *testing.T) {
	auth := NewAuthService(newMockUserRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

	if _, _, err := auth.Register(
		dto.RegisterRequest{Username: "u", Password: "123", Email: "x@mail.com"},
		testhelpers.NewMockPerusahaanService(),
	); err == nil {
		t.Fatal("expected weak password error")
	}
}

func TestRegister_CreateError(t *testing.T) {
	repo := newMockUserRepo()
	repo.failCreate = true

	auth := NewAuthService(repo, NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

	if _, _, err := auth.Register(
		dto.RegisterRequest{Username: "u", Password: "XyZ#91!kLmPq", Email: "x@mail.com"},
		testhelpers.NewMockPerusahaanService(),
	); err == nil {
		t.Fatal("expected create error")
	}
}

func TestAuthService_Register_EmptyFields(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		email    string
		wantErr  bool
		errMsg   string
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
			tokenService := NewTokenService(redis, "test-secret", false, "localhost")
			s := NewAuthService(userRepo, tokenService, newNotifSvc())

			_, _, gotErr := s.Register(
				dto.RegisterRequest{Username: tt.username, Password: tt.password, Email: tt.email},
				testhelpers.NewMockPerusahaanService(),
			)
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

func TestAuthService_Register_WithRoleAndJabatan(t *testing.T) {
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret", false, "localhost")
	authService := NewAuthService(userRepo, tokenService, newNotifSvc())

	roleID := "role-123"
	idJabatan := "jabatan-456"

	user, tokens, err := authService.Register(
		dto.RegisterRequest{
			Username:  "testuser",
			Password:  "MySecureP@ssw0rd2024!",
			Email:     "test@example.com",
			RoleID:    &roleID,
			IDJabatan: &idJabatan,
		},
		testhelpers.NewMockPerusahaanService(),
	)

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
	tokenService := NewTokenService(redis, "test-secret", false, "localhost")
	authService := NewAuthService(userRepo, tokenService, newNotifSvc())

	authService.Register(
		dto.RegisterRequest{Username: "testuser", Password: "MySecureP@ssw0rd2024!", Email: "test@example.com"},
		testhelpers.NewMockPerusahaanService(),
	)

	_, _, err := authService.Register(
		dto.RegisterRequest{Username: "testuser", Password: "DifferentSecureP@ss2024!", Email: "another@example.com"},
		testhelpers.NewMockPerusahaanService(),
	)

	if err == nil {
		t.Fatal("expected error for duplicate username")
	}

	if err.Error() != "username sudah digunakan" {
		t.Errorf("expected 'username sudah digunakan' error, got '%s'", err.Error())
	}
}

//
// =========================
// LOGIN TESTS
// =========================
//

func TestLogin_Success(t *testing.T) {
	repo := newMockUserRepo()
	hash, _ := bcrypt.GenerateFromPassword([]byte("XyZ#91!kLmPq"), bcrypt.DefaultCost)

	repo.users["user"] = &models.User{
		ID:                "1",
		Username:          "user",
		Password:          string(hash),
		RoleName:          "admin",
		PasswordChangedAt: time.Now(), // pastikan tidak expired
	}

	auth := NewAuthService(repo, NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

	// User without MFA enabled should return tokens
	if u, tkn, err := auth.Login("user", "XyZ#91!kLmPq"); err != nil || u == nil || tkn == nil {
		t.Fatal("login success expected")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	auth := NewAuthService(newMockUserRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

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

	auth := NewAuthService(repo, NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

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

	auth := NewAuthService(repo, NewTokenService(redis, "secret", false, "localhost"), newNotifSvc())

	if _, _, err := auth.Login("u", "pass"); err == nil {
		t.Fatal("expected token error")
	}
}

func TestAuthService_Login_Success_Detailed(t *testing.T) {
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret", false, "localhost")
	authService := NewAuthService(userRepo, tokenService, newNotifSvc())

	authService.Register(
		dto.RegisterRequest{Username: "testuser", Password: "MySecureP@ssw0rd2024!", Email: "test@example.com"},
		testhelpers.NewMockPerusahaanService(),
	)

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

	// Note: tokens might be nil if MFA is required
	if tokens == nil && !user.MFAEnabled {
		t.Fatal("expected tokens when MFA is not enabled")
	}

	if tokens != nil {
		if tokens.AccessToken == "" {
			t.Error("expected access token to be generated")
		}

		if tokens.RefreshToken == "" {
			t.Error("expected refresh token to be generated")
		}
	}
}

func TestAuthService_Login_MFAEnabled_ReturnsNilTokens(t *testing.T) {
	repo := newMockUserRepo()
	hash, _ := bcrypt.GenerateFromPassword([]byte("XyZ#91!kLmPq"), bcrypt.DefaultCost)

	secret := "MFASECRET123"
	repo.users["user"] = &models.User{
		ID:                "1",
		Username:          "user",
		Password:          string(hash),
		RoleName:          "admin",
		MFAEnabled:        true,
		MFASecret:         &secret,
		PasswordChangedAt: time.Now(), // pastikan tidak expired
	}

	auth := NewAuthService(repo, NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

	user, tokens, err := auth.Login("user", "XyZ#91!kLmPq")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("expected user to be returned")
	}

	// When MFA is enabled, tokens should be nil
	if tokens != nil {
		t.Error("expected nil tokens when MFA is enabled")
	}
}

func TestAuthService_Login_InvalidUsername(t *testing.T) {
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret", false, "localhost")
	authService := NewAuthService(userRepo, tokenService, newNotifSvc())

	authService.Register(
		dto.RegisterRequest{Username: "testuser", Password: "MySecureP@ssw0rd2024!", Email: "test@example.com"},
		testhelpers.NewMockPerusahaanService(),
	)

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
	tokenService := NewTokenService(redis, "test-secret", false, "localhost")
	authService := NewAuthService(userRepo, tokenService, newNotifSvc())

	authService.Register(
		dto.RegisterRequest{Username: "testuser", Password: "MySecureP@ssw0rd2024!", Email: "test@example.com"},
		testhelpers.NewMockPerusahaanService(),
	)

	_, _, err := authService.Login("testuser", "WrongP@ssword!")

	if err == nil {
		t.Fatal("expected error for invalid password")
	}

	if err.Error() != "username atau password salah" {
		t.Errorf("expected 'username atau password salah', got '%s'", err.Error())
	}
}

//
// =========================
// LOGOUT TESTS
// =========================
//

func TestLogout_Success(t *testing.T) {
	redis := newMockRedis()
	redis.store["refresh_token:token"] = "x"

	auth := NewAuthService(newMockUserRepo(), NewTokenService(redis, "secret", false, "localhost"), newNotifSvc())

	if err := auth.Logout("token"); err != nil {
		t.Fatal("logout success expected")
	}
}

func TestLogout_TokenNotExists(t *testing.T) {
	auth := NewAuthService(newMockUserRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

	if err := auth.Logout("missing"); err != nil {
		t.Fatal("no error expected for missing token")
	}
}

func TestLogout_DeleteError(t *testing.T) {
	redis := newMockRedis()
	redis.store["refresh_token:token"] = "x"
	redis.failDelete = true

	auth := NewAuthService(newMockUserRepo(), NewTokenService(redis, "secret", false, "localhost"), newNotifSvc())

	if err := auth.Logout("token"); err == nil {
		t.Fatal("expected delete error")
	}
}

func TestAuthService_Logout_Success_Detailed(t *testing.T) {
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret", false, "localhost")
	authService := NewAuthService(userRepo, tokenService, newNotifSvc())

	// Register and get tokens
	_, tokens, err := authService.Register(
		dto.RegisterRequest{Username: "testuser", Password: "MySecureP@ssw0rd2024!", Email: "test@example.com"},
		testhelpers.NewMockPerusahaanService(),
	)
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

//
// =========================
// INTERFACE CHECK
// =========================
//

var _ repository.UserRepositoryInterface = (*mockUserRepo)(nil)
