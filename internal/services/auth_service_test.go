package services

import (
	"errors"
	"testing"
	"time"

	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"

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

//
// =========================
// REGISTER TESTS
// =========================
//

func TestRegister_Success(t *testing.T) {
	auth := NewAuthService(newMockUserRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"))

	user, err := auth.Register(
		"user",
		"XyZ#91!kLmPq",
		"user@mail.com",
		nil,
		nil,
	)

	if err != nil || user == nil {
		t.Fatal("register success expected")
	}
}

func TestRegister_UsernameExists(t *testing.T) {
	repo := newMockUserRepo()
	repo.users["user"] = &models.User{Username: "user"}

	auth := NewAuthService(repo, NewTokenService(newMockRedis(), "secret", false, "localhost"))

	if _, err := auth.Register("user", "XyZ#91!kLmPq", "x@mail.com", nil, nil); err == nil {
		t.Fatal("expected username exists error")
	}
}

func TestRegister_EmailExists(t *testing.T) {
	repo := newMockUserRepo()
	repo.users["u1"] = &models.User{Email: "mail@test.com"}

	auth := NewAuthService(repo, NewTokenService(newMockRedis(), "secret", false, "localhost"))

	if _, err := auth.Register("u2", "XyZ#91!kLmPq", "mail@test.com", nil, nil); err == nil {
		t.Fatal("expected email exists error")
	}
}

func TestRegister_WeakPassword(t *testing.T) {
	auth := NewAuthService(newMockUserRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"))

	if _, err := auth.Register("u", "123", "x@mail.com", nil, nil); err == nil {
		t.Fatal("expected weak password error")
	}
}

func TestRegister_CreateError(t *testing.T) {
	repo := newMockUserRepo()
	repo.failCreate = true

	auth := NewAuthService(repo, NewTokenService(newMockRedis(), "secret", false, "localhost"))

	if _, err := auth.Register("u", "XyZ#91!kLmPq", "x@mail.com", nil, nil); err == nil {
		t.Fatal("expected create error")
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
		ID:       "1",
		Username: "user",
		Password: string(hash),
		RoleName: "admin",
	}

	auth := NewAuthService(repo, NewTokenService(newMockRedis(), "secret", false, "localhost"))

	if u, err := auth.Login("user", "XyZ#91!kLmPq"); err != nil || u == nil {
		t.Fatal("login success expected")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	auth := NewAuthService(newMockUserRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"))

	if _, err := auth.Login("x", "pass"); err == nil {
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

	auth := NewAuthService(repo, NewTokenService(newMockRedis(), "secret", false, "localhost"))

	if _, err := auth.Login("u", "wrong"); err == nil {
		t.Fatal("expected wrong password error")
	}
}

//
// =========================
// LOGOUT TESTS
// =========================
//

func TestLogout_Success(t *testing.T) {
	redis := newMockRedis()
	redis.store["token"] = "x"

	auth := NewAuthService(newMockUserRepo(), NewTokenService(redis, "secret", false, "localhost"))

	if err := auth.Logout("token"); err != nil {
		t.Fatal("logout success expected")
	}
}

func TestLogout_TokenNotExists(t *testing.T) {
	auth := NewAuthService(newMockUserRepo(), NewTokenService(newMockRedis(), "secret", false, "localhost"))

	if err := auth.Logout("missing"); err != nil {
		t.Fatal("no error expected for missing token")
	}
}

func TestLogout_ExistsError(t *testing.T) {
	redis := newMockRedis()
	redis.failExists = true

	auth := NewAuthService(newMockUserRepo(), NewTokenService(redis, "secret", false, "localhost"))

	err := auth.Logout("token")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLogout_DeleteError(t *testing.T) {
	redis := newMockRedis()
	redis.store["token"] = "x"
	redis.failDelete = true

	auth := NewAuthService(newMockUserRepo(), NewTokenService(redis, "secret", false, "localhost"))

	if err := auth.Logout("token"); err == nil {
		t.Fatal("expected delete error")
	}
}

//
// =========================
// INTERFACE CHECK
// =========================
//

var _ repository.UserRepositoryInterface = (*mockUserRepo)(nil)
