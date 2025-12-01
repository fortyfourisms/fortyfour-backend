package services

import (
	"fortyfour-backend/internal/testhelpers"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register_Success(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := NewAuthService(repo, "test-secret")

	user, token, err := service.Register("testuser", "password123", "test@example.com")

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

	if token == "" {
		t.Error("expected token to be generated")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
	if err != nil {
		t.Error("password should be hashed correctly")
	}
}

func TestAuthService_Register_DuplicateUsername(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := NewAuthService(repo, "test-secret")

	service.Register("testuser", "password123", "test@example.com")

	_, _, err := service.Register("testuser", "newpassword", "another@example.com")

	if err == nil {
		t.Fatal("expected error for duplicate username")
	}

	if err.Error() != "username already exists" {
		t.Errorf("expected 'username already exists' error, got '%s'", err.Error())
	}
}

func TestAuthService_Register_EmptyFields(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := NewAuthService(repo, "test-secret")

	testCases := []struct {
		name     string
		username string
		password string
		email    string
	}{
		{"empty username", "", "password", "email@test.com"},
		{"empty password", "user", "", "email@test.com"},
		{"empty email", "user", "password", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := service.Register(tc.username, tc.password, tc.email)

			if err == nil {
				t.Fatal("expected error for empty field")
			}

			if err.Error() != "all fields are required" {
				t.Errorf("expected 'all fields are required', got '%s'", err.Error())
			}
		})
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := NewAuthService(repo, "test-secret")

	service.Register("testuser", "password123", "test@example.com")

	user, token, err := service.Login("testuser", "password123")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("expected user to be returned")
	}

	if user.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", user.Username)
	}

	if token == "" {
		t.Error("expected token to be generated")
	}
}

func TestAuthService_Login_InvalidUsername(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := NewAuthService(repo, "test-secret")

	_, _, err := service.Login("nonexistent", "password123")

	if err == nil {
		t.Fatal("expected error for invalid username")
	}

	if err.Error() != "invalid credentials" {
		t.Errorf("expected 'invalid credentials', got '%s'", err.Error())
	}
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := NewAuthService(repo, "test-secret")

	service.Register("testuser", "password123", "test@example.com")

	_, _, err := service.Login("testuser", "wrongpassword")

	if err == nil {
		t.Fatal("expected error for invalid password")
	}

	if err.Error() != "invalid credentials" {
		t.Errorf("expected 'invalid credentials', got '%s'", err.Error())
	}
}
