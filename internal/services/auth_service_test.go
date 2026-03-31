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

// FindByEmail mencari user berdasarkan email
func (m *mockUserRepo) FindByEmail(email string) (*models.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("not found")
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

/*
=====================================
 TEST MFA - SetupMFA
=====================================
*/

func setupAuthService() (*AuthService, *mockUserRepo, *mockRedis) {
	repo := newMockUserRepo()
	redis := newMockRedis()
	tokenSvc := NewTokenService(redis, "test-secret", false, "localhost")
	svc := NewAuthService(repo, tokenSvc, newNotifSvc())
	return svc, repo, redis
}

func TestAuthService_SetupMFA_Success(t *testing.T) {
	svc, repo, _ := setupAuthService()

	user := &models.User{ID: "user-1", Username: "testuser", Email: "test@example.com"}
	repo.users["user-1"] = user

	uri, secret, err := svc.SetupMFA("user-1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if uri == "" {
		t.Error("expected provisioning URI, got empty string")
	}
	if secret == "" {
		t.Error("expected secret, got empty string")
	}
}

func TestAuthService_SetupMFA_UserNotFound(t *testing.T) {
	svc, _, _ := setupAuthService()

	uri, secret, err := svc.SetupMFA("tidak-ada")

	if err == nil {
		t.Error("expected error for non-existent user")
	}
	if uri != "" || secret != "" {
		t.Error("expected empty uri and secret on error")
	}
}

func TestAuthService_SetupMFA_RedisError(t *testing.T) {
	svc, repo, redis := setupAuthService()

	user := &models.User{ID: "user-1", Username: "testuser", Email: "test@example.com"}
	repo.users["user-1"] = user
	redis.failSet = true

	_, _, err := svc.SetupMFA("user-1")

	if err == nil {
		t.Error("expected error when redis fails")
	}
}

/*
=====================================
 TEST MFA - CreateMFASetupToken
=====================================
*/

func TestAuthService_CreateMFASetupToken_Success(t *testing.T) {
	svc, _, _ := setupAuthService()

	token, err := svc.CreateMFASetupToken("user-1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Error("expected non-empty setup token")
	}
}

func TestAuthService_CreateMFASetupToken_RedisError(t *testing.T) {
	svc, _, redis := setupAuthService()
	redis.failSet = true

	token, err := svc.CreateMFASetupToken("user-1")

	if err == nil {
		t.Error("expected error when redis fails")
	}
	if token != "" {
		t.Error("expected empty token on error")
	}
}

func TestAuthService_CreateMFASetupToken_UniquePerCall(t *testing.T) {
	svc, _, _ := setupAuthService()

	token1, _ := svc.CreateMFASetupToken("user-1")
	token2, _ := svc.CreateMFASetupToken("user-1")

	if token1 == token2 {
		t.Error("expected unique tokens per call")
	}
}

/*
=====================================
 TEST MFA - ValidateMFASetupToken
=====================================
*/

func TestAuthService_ValidateMFASetupToken_Success(t *testing.T) {
	svc, _, _ := setupAuthService()

	token, _ := svc.CreateMFASetupToken("user-1")

	userID, err := svc.ValidateMFASetupToken(token)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if userID != "user-1" {
		t.Errorf("expected userID 'user-1', got '%s'", userID)
	}
}

func TestAuthService_ValidateMFASetupToken_InvalidToken(t *testing.T) {
	svc, _, _ := setupAuthService()

	userID, err := svc.ValidateMFASetupToken("token-tidak-ada")

	if err == nil {
		t.Error("expected error for invalid token")
	}
	if userID != "" {
		t.Error("expected empty userID on error")
	}
}

/*
=====================================
 TEST MFA - CreateMFAPending
=====================================
*/

func TestAuthService_CreateMFAPending_Success(t *testing.T) {
	svc, _, _ := setupAuthService()

	token, err := svc.CreateMFAPending("user-1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Error("expected non-empty mfa pending token")
	}
}

func TestAuthService_CreateMFAPending_RedisError(t *testing.T) {
	svc, _, redis := setupAuthService()
	redis.failSet = true

	token, err := svc.CreateMFAPending("user-1")

	if err == nil {
		t.Error("expected error when redis fails")
	}
	if token != "" {
		t.Error("expected empty token on error")
	}
}

func TestAuthService_CreateMFAPending_UniquePerCall(t *testing.T) {
	svc, _, _ := setupAuthService()

	token1, _ := svc.CreateMFAPending("user-1")
	token2, _ := svc.CreateMFAPending("user-1")

	if token1 == token2 {
		t.Error("expected unique pending tokens per call")
	}
}

/*
=====================================
 TEST MFA - EnableMFA
=====================================
*/

func TestAuthService_EnableMFA_SetupNotFound(t *testing.T) {
	svc, _, _ := setupAuthService()

	// Tidak ada mfa_setup di redis → harus error
	err := svc.EnableMFA("user-1", "123456")

	if err == nil {
		t.Error("expected error when mfa setup not found in redis")
	}
	if err.Error() != "mfa setup expired or not found" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestAuthService_EnableMFA_InvalidCode(t *testing.T) {
	svc, repo, redis := setupAuthService()

	user := &models.User{ID: "user-1", Username: "testuser", Email: "test@example.com"}
	repo.users["user-1"] = user

	// Setup MFA dulu agar secret tersimpan di redis
	_, _, err := svc.SetupMFA("user-1")
	if err != nil {
		t.Fatalf("setup MFA failed: %v", err)
	}
	_ = redis

	// Coba enable dengan kode salah
	err = svc.EnableMFA("user-1", "000000")

	if err == nil {
		t.Error("expected error for invalid TOTP code")
	}
	if err.Error() != "invalid mfa code" {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

/*
=====================================
 TEST MFA - VerifyMFA
=====================================
*/

func TestAuthService_VerifyMFA_InvalidToken(t *testing.T) {
	svc, _, _ := setupAuthService()

	user, tokens, err := svc.VerifyMFA("token-tidak-ada", "123456")

	if err == nil {
		t.Error("expected error for invalid mfa token")
	}
	if user != nil || tokens != nil {
		t.Error("expected nil user and tokens on error")
	}
	if err.Error() != "invalid or expired mfa token" {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestAuthService_VerifyMFA_UserNotFound(t *testing.T) {
	svc, _, redis := setupAuthService()

	// Set pending token yang menunjuk ke user yang tidak ada
	redis.store["mfa_pending:test-token"] = "user-tidak-ada"

	user, tokens, err := svc.VerifyMFA("test-token", "123456")

	if err == nil {
		t.Error("expected error when user not found")
	}
	if user != nil || tokens != nil {
		t.Error("expected nil user and tokens on error")
	}
}

func TestAuthService_VerifyMFA_MFANotConfigured(t *testing.T) {
	svc, repo, redis := setupAuthService()

	// User tanpa MFASecret
	user := &models.User{ID: "user-1", Username: "testuser", Email: "test@example.com", MFASecret: nil}
	repo.users["user-1"] = user
	redis.store["mfa_pending:test-token"] = "user-1"

	result, tokens, err := svc.VerifyMFA("test-token", "123456")

	if err == nil {
		t.Error("expected error when mfa not configured")
	}
	if result != nil || tokens != nil {
		t.Error("expected nil result and tokens")
	}
	if err.Error() != "mfa not configured" {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestAuthService_VerifyMFA_InvalidCode(t *testing.T) {
	svc, repo, redis := setupAuthService()

	secret := "JBSWY3DPEHPK3PXP" // valid base32 TOTP secret
	user := &models.User{ID: "user-1", Username: "testuser", Email: "test@example.com", MFASecret: &secret, MFAEnabled: true}
	repo.users["user-1"] = user
	redis.store["mfa_pending:test-token"] = "user-1"

	result, tokens, err := svc.VerifyMFA("test-token", "000000")

	if err == nil {
		t.Error("expected error for invalid TOTP code")
	}
	if result != nil || tokens != nil {
		t.Error("expected nil result and tokens on invalid code")
	}
	if err.Error() != "invalid mfa code" {
		t.Errorf("unexpected error: %s", err.Error())
	}
}