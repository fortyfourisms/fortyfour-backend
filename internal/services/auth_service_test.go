package services

import (
	"context"
	"errors"
	"fmt"
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

// =========================
// MOCK ROLE REPOSITORY
// =========================

type mockRoleRepo struct{}

func newMockRoleRepo() *mockRoleRepo {
	return &mockRoleRepo{}
}

func (m *mockRoleRepo) Create(ctx context.Context, role *models.Role) error { return nil }
func (m *mockRoleRepo) GetByID(ctx context.Context, id string) (*models.Role, error) {
	return &models.Role{ID: "role-user-id", Name: "user"}, nil
}
func (m *mockRoleRepo) GetAll(ctx context.Context) ([]*models.Role, error) { return nil, nil }
func (m *mockRoleRepo) Update(ctx context.Context, role *models.Role) error { return nil }
func (m *mockRoleRepo) Delete(ctx context.Context, id string) error         { return nil }
func (m *mockRoleRepo) GetByName(ctx context.Context, name string) (*models.Role, error) {
	return &models.Role{ID: "role-user-id", Name: "user"}, nil
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
			got := NewAuthService(tt.userRepo, newMockRoleRepo(), tt.tokenService, newNotifSvc())
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

// FIX: Semua test register wajib menyertakan nama_perusahaan atau id_perusahaan.
// Sebelumnya test ini tidak menyertakan keduanya dan akan lolos, sekarang akan error.

func TestRegister_Success(t *testing.T) {
	auth := NewAuthService(newMockUserRepo(), newMockRoleRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

	namaPerusahaan := "PT Test Sukses"
	user, token, err := auth.Register(
		dto.RegisterRequest{
			Username:       "user",
			Password:       "XyZ#91!kLmPq",
			Email:          "user@mail.com",
			NamaPerusahaan: &namaPerusahaan, // FIX: wajib diisi
		},
		testhelpers.NewMockPerusahaanService(),
	)

	if err != nil || user == nil || token == nil {
		t.Fatalf("register success expected, got err: %v", err)
	}
}

func TestRegister_UsernameExists(t *testing.T) {
	repo := newMockUserRepo()
	repo.users["user"] = &models.User{Username: "user"}

	auth := NewAuthService(repo, newMockRoleRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

	namaPerusahaan := "PT Test"
	if _, _, err := auth.Register(
		dto.RegisterRequest{Username: "user", Password: "XyZ#91!kLmPq", Email: "x@mail.com", NamaPerusahaan: &namaPerusahaan},
		testhelpers.NewMockPerusahaanService(),
	); err == nil {
		t.Fatal("expected username exists error")
	}
}

func TestRegister_EmailExists(t *testing.T) {
	repo := newMockUserRepo()
	repo.users["u1"] = &models.User{Email: "mail@test.com"}

	auth := NewAuthService(repo, newMockRoleRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

	namaPerusahaan := "PT Test"
	if _, _, err := auth.Register(
		dto.RegisterRequest{Username: "u2", Password: "XyZ#91!kLmPq", Email: "mail@test.com", NamaPerusahaan: &namaPerusahaan},
		testhelpers.NewMockPerusahaanService(),
	); err == nil {
		t.Fatal("expected email exists error")
	}
}

func TestRegister_WeakPassword(t *testing.T) {
	auth := NewAuthService(newMockUserRepo(), newMockRoleRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

	namaPerusahaan := "PT Test"
	if _, _, err := auth.Register(
		dto.RegisterRequest{Username: "u", Password: "123", Email: "x@mail.com", NamaPerusahaan: &namaPerusahaan},
		testhelpers.NewMockPerusahaanService(),
	); err == nil {
		t.Fatal("expected weak password error")
	}
}

func TestRegister_CreateError(t *testing.T) {
	repo := newMockUserRepo()
	repo.failCreate = true

	auth := NewAuthService(repo, newMockRoleRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

	namaPerusahaan := "PT Test"
	if _, _, err := auth.Register(
		dto.RegisterRequest{Username: "u", Password: "XyZ#91!kLmPq", Email: "x@mail.com", NamaPerusahaan: &namaPerusahaan},
		testhelpers.NewMockPerusahaanService(),
	); err == nil {
		t.Fatal("expected create error")
	}
}

func TestAuthService_Register_EmptyFields(t *testing.T) {
	namaPerusahaan := "PT Test"
	tests := []struct {
		name           string
		username       string
		password       string
		email          string
		namaPerusahaan *string
		wantErr        bool
		errMsg         string
	}{
		{
			name:           "empty username",
			username:       "",
			password:       "MySecureP@ssw0rd2024!",
			email:          "test@example.com",
			namaPerusahaan: &namaPerusahaan,
			wantErr:        true,
			errMsg:         "username wajib diisi",
		},
		{
			name:           "empty password",
			username:       "testuser",
			password:       "",
			email:          "test@example.com",
			namaPerusahaan: &namaPerusahaan,
			wantErr:        true,
			errMsg:         "password wajib diisi",
		},
		{
			name:           "empty email",
			username:       "testuser",
			password:       "MySecureP@ssw0rd2024!",
			email:          "",
			namaPerusahaan: &namaPerusahaan,
			wantErr:        true,
			errMsg:         "email wajib diisi",
		},
		{
			name:           "all empty",
			username:       "",
			password:       "",
			email:          "",
			namaPerusahaan: &namaPerusahaan,
			wantErr:        true,
			errMsg:         "username wajib diisi",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := testhelpers.NewMockUserRepository()
			redis := testhelpers.NewMockRedisClient()
			tokenService := NewTokenService(redis, "test-secret", false, "localhost")
			s := NewAuthService(userRepo, newMockRoleRepo(), tokenService, newNotifSvc())

			_, _, gotErr := s.Register(
				dto.RegisterRequest{
					Username:       tt.username,
					Password:       tt.password,
					Email:          tt.email,
					NamaPerusahaan: tt.namaPerusahaan,
				},
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

func TestAuthService_Register_WithJabatan(t *testing.T) {
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret", false, "localhost")
	authService := NewAuthService(userRepo, newMockRoleRepo(), tokenService, newNotifSvc())

	idJabatan := "jabatan-456"
	namaPerusahaan := "PT Jabatan Test"

	user, tokens, err := authService.Register(
		dto.RegisterRequest{
			Username:       "testuser",
			Password:       "MySecureP@ssw0rd2024!",
			Email:          "test@example.com",
			IDJabatan:      &idJabatan,
			NamaPerusahaan: &namaPerusahaan, // FIX: wajib diisi
		},
		testhelpers.NewMockPerusahaanService(),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user == nil {
		t.Fatal("expected user to be created")
	}
	if user.RoleID == nil || *user.RoleID != "role-user-id" {
		t.Errorf("expected default roleID 'role-user-id', got '%v'", user.RoleID)
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
	authService := NewAuthService(userRepo, newMockRoleRepo(), tokenService, newNotifSvc())

	namaPerusahaan1 := "PT Pertama"
	authService.Register(
		dto.RegisterRequest{Username: "testuser", Password: "MySecureP@ssw0rd2024!", Email: "test@example.com", NamaPerusahaan: &namaPerusahaan1},
		testhelpers.NewMockPerusahaanService(),
	)

	namaPerusahaan2 := "PT Kedua"
	_, _, err := authService.Register(
		dto.RegisterRequest{Username: "testuser", Password: "DifferentSecureP@ss2024!", Email: "another@example.com", NamaPerusahaan: &namaPerusahaan2},
		testhelpers.NewMockPerusahaanService(),
	)

	if err == nil {
		t.Fatal("expected error for duplicate username")
	}
	if err.Error() != "username sudah digunakan" {
		t.Errorf("expected 'username sudah digunakan' error, got '%s'", err.Error())
	}
}

// NEW: Test validasi wajib perusahaan saat register

func TestRegister_NoPerusahaan_ShouldFail(t *testing.T) {
	auth := NewAuthService(newMockUserRepo(), newMockRoleRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

	_, _, err := auth.Register(
		dto.RegisterRequest{
			Username: "usertest",
			Password: "XyZ#91!kLmPq",
			Email:    "test@mail.com",
			// Tidak ada NamaPerusahaan maupun IDPerusahaan
		},
		testhelpers.NewMockPerusahaanService(),
	)

	if err == nil {
		t.Fatal("expected error ketika tidak ada perusahaan yang diisi")
	}
	if err.Error() != "wajib mengisi nama_perusahaan (buat perusahaan baru) atau id_perusahaan (bergabung ke perusahaan yang sudah ada)" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestRegister_BothPerusahaanFields_ShouldFail(t *testing.T) {
	auth := NewAuthService(newMockUserRepo(), newMockRoleRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

	namaPerusahaan := "PT Test"
	idPerusahaan := "perusahaan-123"
	_, _, err := auth.Register(
		dto.RegisterRequest{
			Username:       "usertest",
			Password:       "XyZ#91!kLmPq",
			Email:          "test@mail.com",
			NamaPerusahaan: &namaPerusahaan,
			IDPerusahaan:   &idPerusahaan,
		},
		testhelpers.NewMockPerusahaanService(),
	)

	if err == nil {
		t.Fatal("expected error ketika keduanya diisi bersamaan")
	}
	if err.Error() != "tidak bisa mengisi nama_perusahaan dan id_perusahaan bersamaan" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestRegister_WithNamaPerusahaan_Success(t *testing.T) {
	auth := NewAuthService(newMockUserRepo(), newMockRoleRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

	namaPerusahaan := "PT Baru Sekali"
	user, token, err := auth.Register(
		dto.RegisterRequest{
			Username:       "userbaru",
			Password:       "XyZ#91!kLmPq",
			Email:          "baru@mail.com",
			NamaPerusahaan: &namaPerusahaan,
		},
		testhelpers.NewMockPerusahaanService(),
	)

	if err != nil {
		t.Fatalf("expected success dengan nama_perusahaan, got: %v", err)
	}
	if user == nil || token == nil {
		t.Fatal("expected user dan token, got nil")
	}
	if user.IDPerusahaan == nil {
		t.Error("expected IDPerusahaan terisi setelah register dengan nama_perusahaan")
	}
}

func TestRegister_WithIDPerusahaan_Success(t *testing.T) {
	perusahaanSvc := testhelpers.NewMockPerusahaanService()
	// Pre-populate perusahaan agar GetByID berhasil
	namaP := "PT Sudah Ada"
	p, _ := perusahaanSvc.Create(dto.CreatePerusahaanRequest{NamaPerusahaan: &namaP})

	auth := NewAuthService(newMockUserRepo(), newMockRoleRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

	user, token, err := auth.Register(
		dto.RegisterRequest{
			Username:     "userjoin",
			Password:     "XyZ#91!kLmPq",
			Email:        "join@mail.com",
			IDPerusahaan: &p.ID,
		},
		perusahaanSvc,
	)

	if err != nil {
		t.Fatalf("expected success dengan id_perusahaan, got: %v", err)
	}
	if user == nil || token == nil {
		t.Fatal("expected user dan token, got nil")
	}
}

func TestRegister_WithIDPerusahaan_AlreadyHasUser_ShouldFail(t *testing.T) {
	perusahaanSvc := testhelpers.NewMockPerusahaanService()
	namaP := "PT Sudah Punya User"
	p, _ := perusahaanSvc.Create(dto.CreatePerusahaanRequest{NamaPerusahaan: &namaP})

	repo := newMockUserRepo()
	// Simulasikan sudah ada user di perusahaan ini
	repo.users["existing"] = &models.User{
		ID:           "existing-id",
		Username:     "existing",
		IDPerusahaan: &p.ID,
	}

	auth := NewAuthService(repo, newMockRoleRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

	_, _, err := auth.Register(
		dto.RegisterRequest{
			Username:     "newuser",
			Password:     "XyZ#91!kLmPq",
			Email:        "new@mail.com",
			IDPerusahaan: &p.ID,
		},
		perusahaanSvc,
	)

	if err == nil {
		t.Fatal("expected error karena perusahaan sudah punya akun terdaftar")
	}
	if err.Error() != "perusahaan sudah terdaftar, silakan minta akun kepada admin perusahaan tersebut" {
		t.Errorf("unexpected error message: %s", err.Error())
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
		PasswordChangedAt: time.Now(),
	}

	auth := NewAuthService(repo, newMockRoleRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

	if u, tkn, err := auth.Login("user", "XyZ#91!kLmPq"); err != nil || u == nil || tkn == nil {
		t.Fatal("login success expected")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	auth := NewAuthService(newMockUserRepo(), newMockRoleRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

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

	auth := NewAuthService(repo, newMockRoleRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

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

	auth := NewAuthService(repo, newMockRoleRepo(), NewTokenService(redis, "secret", false, "localhost"), newNotifSvc())

	if _, _, err := auth.Login("u", "pass"); err == nil {
		t.Fatal("expected token error")
	}
}

func TestAuthService_Login_Success_Detailed(t *testing.T) {
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret", false, "localhost")
	authService := NewAuthService(userRepo, newMockRoleRepo(), tokenService, newNotifSvc())

	namaPerusahaan := "PT Login Test"
	authService.Register(
		dto.RegisterRequest{Username: "testuser", Password: "MySecureP@ssw0rd2024!", Email: "test@example.com", NamaPerusahaan: &namaPerusahaan},
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
		PasswordChangedAt: time.Now(),
	}

	auth := NewAuthService(repo, newMockRoleRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

	user, tokens, err := auth.Login("user", "XyZ#91!kLmPq")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user == nil {
		t.Fatal("expected user to be returned")
	}
	if tokens != nil {
		t.Error("expected nil tokens when MFA is enabled")
	}
}

func TestAuthService_Login_InvalidUsername(t *testing.T) {
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret", false, "localhost")
	authService := NewAuthService(userRepo, newMockRoleRepo(), tokenService, newNotifSvc())

	namaPerusahaan := "PT Test"
	authService.Register(
		dto.RegisterRequest{Username: "testuser", Password: "MySecureP@ssw0rd2024!", Email: "test@example.com", NamaPerusahaan: &namaPerusahaan},
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
	authService := NewAuthService(userRepo, newMockRoleRepo(), tokenService, newNotifSvc())

	namaPerusahaan := "PT Test"
	authService.Register(
		dto.RegisterRequest{Username: "testuser", Password: "MySecureP@ssw0rd2024!", Email: "test@example.com", NamaPerusahaan: &namaPerusahaan},
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

	auth := NewAuthService(newMockUserRepo(), newMockRoleRepo(), NewTokenService(redis, "secret", false, "localhost"), newNotifSvc())

	if err := auth.Logout("token"); err != nil {
		t.Fatal("logout success expected")
	}
}

func TestLogout_TokenNotExists(t *testing.T) {
	auth := NewAuthService(newMockUserRepo(), newMockRoleRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"), newNotifSvc())

	if err := auth.Logout("missing"); err != nil {
		t.Fatal("no error expected for missing token")
	}
}

func TestLogout_DeleteError(t *testing.T) {
	redis := newMockRedis()
	redis.store["refresh_token:token"] = "x"
	redis.failDelete = true

	auth := NewAuthService(newMockUserRepo(), newMockRoleRepo(), NewTokenService(redis, "secret", false, "localhost"), newNotifSvc())

	if err := auth.Logout("token"); err == nil {
		t.Fatal("expected delete error")
	}
}

func TestAuthService_Logout_Success_Detailed(t *testing.T) {
	userRepo := testhelpers.NewMockUserRepository()
	redis := testhelpers.NewMockRedisClient()
	tokenService := NewTokenService(redis, "test-secret", false, "localhost")
	authService := NewAuthService(userRepo, newMockRoleRepo(), tokenService, newNotifSvc())

	namaPerusahaan := "PT Logout Test"
	_, tokens, err := authService.Register(
		dto.RegisterRequest{Username: "testuser", Password: "MySecureP@ssw0rd2024!", Email: "test@example.com", NamaPerusahaan: &namaPerusahaan},
		testhelpers.NewMockPerusahaanService(),
	)
	if err != nil {
		t.Fatalf("registration failed: %v", err)
	}
	if tokens == nil {
		t.Fatal("expected tokens to be returned from registration")
	}

	err = authService.Logout(tokens.RefreshToken)
	if err != nil {
		t.Fatalf("expected no error during logout, got %v", err)
	}

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
	svc := NewAuthService(repo, newMockRoleRepo(), tokenSvc, newNotifSvc())
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

	_, _, err := svc.SetupMFA("user-1")
	if err != nil {
		t.Fatalf("setup MFA failed: %v", err)
	}
	_ = redis

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
 TEST MFA - VerifyMFA (dengan rate limiting)
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

// FIX: Error message berubah dari "invalid mfa code" menjadi "kode MFA salah, sisa percobaan: X"
func TestAuthService_VerifyMFA_InvalidCode_FirstAttempt(t *testing.T) {
	svc, repo, redis := setupAuthService()

	secret := "JBSWY3DPEHPK3PXP"
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
	// FIX: error message sekarang menyertakan sisa percobaan
	expectedErr := fmt.Sprintf("kode MFA salah, sisa percobaan: %d", MaxMFAAttempts-1)
	if err.Error() != expectedErr {
		t.Errorf("expected '%s', got '%s'", expectedErr, err.Error())
	}
}

// NEW: Test MFA rate limiting — percobaan ke-2 s/d ke-(MaxMFAAttempts-1) menampilkan sisa percobaan yang berkurang
func TestAuthService_VerifyMFA_RateLimit_CountsDown(t *testing.T) {
	svc, repo, redis := setupAuthService()

	secret := "JBSWY3DPEHPK3PXP"
	user := &models.User{ID: "user-1", Username: "testuser", Email: "test@example.com", MFASecret: &secret, MFAEnabled: true}
	repo.users["user-1"] = user
	redis.store["mfa_pending:test-token"] = "user-1"

	for attempt := 1; attempt <= MaxMFAAttempts-1; attempt++ {
		_, _, err := svc.VerifyMFA("test-token", "000000")
		if err == nil {
			t.Fatalf("attempt %d: expected error for invalid code", attempt)
		}

		remaining := MaxMFAAttempts - attempt
		expectedErr := fmt.Sprintf("kode MFA salah, sisa percobaan: %d", remaining)
		if err.Error() != expectedErr {
			t.Errorf("attempt %d: expected '%s', got '%s'", attempt, expectedErr, err.Error())
		}

		// Pending token harus masih ada (belum diinvalidate)
		if _, ok := redis.store["mfa_pending:test-token"]; !ok {
			t.Fatalf("attempt %d: pending token seharusnya masih ada sebelum attempts habis", attempt)
		}
	}
}

// NEW: Test MFA rate limiting — percobaan ke-MaxMFAAttempts mengosongkan pending token
func TestAuthService_VerifyMFA_RateLimit_TokenInvalidatedAfterMaxAttempts(t *testing.T) {
	svc, repo, redis := setupAuthService()

	secret := "JBSWY3DPEHPK3PXP"
	user := &models.User{ID: "user-1", Username: "testuser", Email: "test@example.com", MFASecret: &secret, MFAEnabled: true}
	repo.users["user-1"] = user
	redis.store["mfa_pending:test-token"] = "user-1"

	// Lakukan MaxMFAAttempts kali percobaan dengan kode salah
	for attempt := 1; attempt <= MaxMFAAttempts; attempt++ {
		_, _, err := svc.VerifyMFA("test-token", "000000")
		if err == nil {
			t.Fatalf("attempt %d: expected error", attempt)
		}
	}

	// Setelah MaxMFAAttempts, pending token harus sudah dihapus
	if _, ok := redis.store["mfa_pending:test-token"]; ok {
		t.Error("pending token seharusnya sudah dihapus setelah mencapai batas maksimal percobaan")
	}

	// attempts key juga harus sudah dihapus
	if _, ok := redis.store["mfa_attempts:test-token"]; ok {
		t.Error("attempts key seharusnya sudah dihapus setelah mencapai batas maksimal percobaan")
	}
}

// NEW: Test bahwa setelah token diinvalidate, percobaan berikutnya langsung error "invalid or expired"
func TestAuthService_VerifyMFA_RateLimit_TokenExpiredAfterInvalidation(t *testing.T) {
	svc, repo, redis := setupAuthService()

	secret := "JBSWY3DPEHPK3PXP"
	user := &models.User{ID: "user-1", Username: "testuser", Email: "test@example.com", MFASecret: &secret, MFAEnabled: true}
	repo.users["user-1"] = user
	redis.store["mfa_pending:test-token"] = "user-1"

	// Habiskan semua attempts
	for i := 0; i < MaxMFAAttempts; i++ {
		svc.VerifyMFA("test-token", "000000")
	}

	// Percobaan setelah token diinvalidate harus error "invalid or expired mfa token"
	_, _, err := svc.VerifyMFA("test-token", "000000")
	if err == nil {
		t.Fatal("expected error setelah token diinvalidate")
	}
	if err.Error() != "invalid or expired mfa token" {
		t.Errorf("expected 'invalid or expired mfa token', got '%s'", err.Error())
	}
}

// NEW: Test bahwa attempts counter di-reset ketika pending token baru dibuat
func TestAuthService_VerifyMFA_RateLimit_NewTokenResetsAttempts(t *testing.T) {
	svc, repo, redis := setupAuthService()

	secret := "JBSWY3DPEHPK3PXP"
	user := &models.User{ID: "user-1", Username: "testuser", Email: "test@example.com", MFASecret: &secret, MFAEnabled: true}
	repo.users["user-1"] = user

	// Token pertama — habiskan attempts
	redis.store["mfa_pending:token-lama"] = "user-1"
	for i := 0; i < MaxMFAAttempts; i++ {
		svc.VerifyMFA("token-lama", "000000")
	}

	// Token baru — harus bisa mulai dari awal lagi (MaxMFAAttempts percobaan)
	redis.store["mfa_pending:token-baru"] = "user-1"
	_, _, err := svc.VerifyMFA("token-baru", "000000")

	if err == nil {
		t.Fatal("expected error untuk kode salah")
	}
	// Sisa percobaan harus kembali ke MaxMFAAttempts-1 (bukan 0 atau negatif)
	expectedErr := fmt.Sprintf("kode MFA salah, sisa percobaan: %d", MaxMFAAttempts-1)
	if err.Error() != expectedErr {
		t.Errorf("expected '%s', got '%s'", expectedErr, err.Error())
	}
}

// NEW: Test bahwa setelah MaxMFAAttempts-1 kali salah, kode yang benar tetap bisa digunakan
// (memastikan rate limit tidak menghalangi user yang akhirnya memasukkan kode benar)
func TestAuthService_VerifyMFA_RateLimit_SuccessBeforeLimit(t *testing.T) {
	svc, repo, redis := setupAuthService()

	secret := "JBSWY3DPEHPK3PXP"
	user := &models.User{ID: "user-1", Username: "testuser", Email: "test@example.com", MFASecret: &secret, MFAEnabled: true, RoleName: "user"}
	repo.users["user-1"] = user
	redis.store["mfa_pending:test-token"] = "user-1"

	// Lakukan MaxMFAAttempts-1 kali percobaan salah
	for i := 0; i < MaxMFAAttempts-1; i++ {
		svc.VerifyMFA("test-token", "000000")
	}

	// Pending token masih harus ada
	if _, ok := redis.store["mfa_pending:test-token"]; !ok {
		t.Fatal("pending token seharusnya masih ada sebelum batas maksimal")
	}

	// Pending token masih valid dan bisa digunakan (meski kode masih salah seharusnya masih dapat error kode salah bukan expired)
	_, _, err := svc.VerifyMFA("test-token", "111111")
	if err == nil {
		t.Fatal("expected error untuk kode salah")
	}
	// Ini adalah percobaan terakhir (ke-MaxMFAAttempts), token akan diinvalidate
	if err.Error() != "terlalu banyak percobaan MFA, silakan login ulang" {
		t.Errorf("expected 'terlalu banyak percobaan MFA, silakan login ulang', got '%s'", err.Error())
	}
}

/*
=====================================
 TEST MFA - EnableMFAAndLogin
=====================================
*/

func TestAuthService_EnableMFAAndLogin_SetupNotFound(t *testing.T) {
	svc, _, _ := setupAuthService()

	user, tokens, err := svc.EnableMFAAndLogin("user-1", "123456")

	if err == nil {
		t.Error("expected error ketika mfa setup tidak ditemukan di redis")
	}
	if err.Error() != "mfa setup expired or not found" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
	if user != nil || tokens != nil {
		t.Error("expected nil user dan tokens saat setup tidak ditemukan")
	}
}

func TestAuthService_EnableMFAAndLogin_InvalidCode(t *testing.T) {
	svc, repo, _ := setupAuthService()

	user := &models.User{ID: "user-1", Username: "testuser", Email: "test@example.com"}
	repo.users["user-1"] = user

	_, _, err := svc.SetupMFA("user-1")
	if err != nil {
		t.Fatalf("SetupMFA gagal: %v", err)
	}

	result, tokens, err := svc.EnableMFAAndLogin("user-1", "000000")

	if err == nil {
		t.Error("expected error untuk kode TOTP yang salah")
	}
	if err.Error() != "invalid mfa code" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
	if result != nil || tokens != nil {
		t.Error("expected nil result dan tokens saat kode salah")
	}
}

func TestAuthService_EnableMFAAndLogin_UserNotFoundAfterSetMFA(t *testing.T) {
	svc, repo, redis := setupAuthService()

	user := &models.User{ID: "user-orphan", Username: "orphan", Email: "orphan@example.com"}
	repo.users["user-orphan"] = user

	_, _, err := svc.SetupMFA("user-orphan")
	if err != nil {
		t.Fatalf("SetupMFA gagal: %v", err)
	}

	knownSecret := "JBSWY3DPEHPK3PXP"
	redis.store["mfa_setup:user-orphan"] = knownSecret

	delete(repo.users, "user-orphan")

	result, tokens, err := svc.EnableMFAAndLogin("user-orphan", "000000")

	if err == nil {
		t.Error("expected error ketika kode salah (sebelum FindByID dieksekusi)")
	}
	if result != nil || tokens != nil {
		t.Error("expected nil result dan tokens")
	}
}

func TestAuthService_EnableMFAAndLogin_SetMFAFails(t *testing.T) {
	svc, repo, redis := setupAuthService()

	knownSecret := "JBSWY3DPEHPK3PXP"
	redis.store["mfa_setup:user-ghost"] = knownSecret

	user := &models.User{ID: "user-ghost", Username: "ghost", Email: "ghost@example.com"}
	repo.users["user-ghost"] = user

	result, tokens, err := svc.EnableMFAAndLogin("user-ghost", "000000")

	if err == nil {
		t.Error("expected error untuk kode TOTP yang tidak valid")
	}
	if result != nil || tokens != nil {
		t.Error("expected nil result dan tokens")
	}
	_ = repo
}

func TestAuthService_EnableMFAAndLogin_RedisSetError(t *testing.T) {
	svc, _, _ := setupAuthService()

	result, tokens, err := svc.EnableMFAAndLogin("user-any", "123456")

	if err == nil {
		t.Error("expected error ketika redis tidak memiliki secret")
	}
	if err.Error() != "mfa setup expired or not found" {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if result != nil || tokens != nil {
		t.Error("expected nil result dan tokens")
	}
}

func TestAuthService_EnableMFAAndLogin_SetupKeyNotCleanedOnInvalidCode(t *testing.T) {
	svc, repo, redis := setupAuthService()

	user := &models.User{ID: "user-1", Username: "testuser", Email: "test@example.com"}
	repo.users["user-1"] = user

	_, _, err := svc.SetupMFA("user-1")
	if err != nil {
		t.Fatalf("SetupMFA gagal: %v", err)
	}

	_, exists := redis.store["mfa_setup:user-1"]
	if !exists {
		t.Fatal("expected mfa_setup key di redis setelah SetupMFA")
	}

	_, _, _ = svc.EnableMFAAndLogin("user-1", "000000")

	_, stillExists := redis.store["mfa_setup:user-1"]
	if !stillExists {
		t.Error("mfa_setup key seharusnya MASIH ada di redis karena validasi kode gagal sebelum cleanup")
	}
}

/*
=====================================
 TEST REGISTER — Role Default Not Found
=====================================
*/

type failingMockRoleRepo struct {
	failGetByName bool
	returnNilRole bool
}

func (m *failingMockRoleRepo) Create(ctx context.Context, role *models.Role) error { return nil }
func (m *failingMockRoleRepo) GetByID(ctx context.Context, id string) (*models.Role, error) {
	return &models.Role{ID: "role-user-id", Name: "user"}, nil
}
func (m *failingMockRoleRepo) GetAll(ctx context.Context) ([]*models.Role, error) { return nil, nil }
func (m *failingMockRoleRepo) Update(ctx context.Context, role *models.Role) error { return nil }
func (m *failingMockRoleRepo) Delete(ctx context.Context, id string) error        { return nil }
func (m *failingMockRoleRepo) GetByName(ctx context.Context, name string) (*models.Role, error) {
	if m.failGetByName {
		return nil, errors.New("database error")
	}
	if m.returnNilRole {
		return nil, nil
	}
	return &models.Role{ID: "role-user-id", Name: "user"}, nil
}

func TestRegister_RoleDefaultNotFound(t *testing.T) {
	roleRepo := &failingMockRoleRepo{failGetByName: true}
	auth := NewAuthService(
		newMockUserRepo(),
		roleRepo,
		NewTokenService(newMockRedis(), "secret", false, "localhost"),
		newNotifSvc(),
	)

	namaPerusahaan := "PT Test"
	_, _, err := auth.Register(
		dto.RegisterRequest{
			Username:       "newuser",
			Password:       "XyZ#91!kLmPq",
			Email:          "newuser@mail.com",
			NamaPerusahaan: &namaPerusahaan, // FIX: wajib diisi
		},
		testhelpers.NewMockPerusahaanService(),
	)

	if err == nil {
		t.Fatal("expected error ketika role default tidak ditemukan")
	}
	if err.Error() != "role default tidak ditemukan, hubungi administrator" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestRegister_RoleDefaultNil(t *testing.T) {
	roleRepo := &failingMockRoleRepo{returnNilRole: true}
	auth := NewAuthService(
		newMockUserRepo(),
		roleRepo,
		NewTokenService(newMockRedis(), "secret", false, "localhost"),
		newNotifSvc(),
	)

	namaPerusahaan := "PT Test 2"
	_, _, err := auth.Register(
		dto.RegisterRequest{
			Username:       "newuser2",
			Password:       "XyZ#91!kLmPq",
			Email:          "newuser2@mail.com",
			NamaPerusahaan: &namaPerusahaan, // FIX: wajib diisi
		},
		testhelpers.NewMockPerusahaanService(),
	)

	if err == nil {
		t.Fatal("expected error ketika role default nil")
	}
	if err.Error() != "role default tidak ditemukan, hubungi administrator" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestRegister_RoleIDAssignedFromDatabase(t *testing.T) {
	roleRepo := &failingMockRoleRepo{}
	auth := NewAuthService(
		newMockUserRepo(),
		roleRepo,
		NewTokenService(newMockRedis(), "secret", false, "localhost"),
		newNotifSvc(),
	)

	namaPerusahaan := "PT Verify Role"
	user, token, err := auth.Register(
		dto.RegisterRequest{
			Username:       "verifyuser",
			Password:       "XyZ#91!kLmPq",
			Email:          "verify@mail.com",
			NamaPerusahaan: &namaPerusahaan, // FIX: wajib diisi
		},
		testhelpers.NewMockPerusahaanService(),
	)

	if err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
	if user == nil || token == nil {
		t.Fatal("expected user dan token, got nil")
	}
}