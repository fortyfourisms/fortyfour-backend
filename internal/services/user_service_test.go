package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/testhelpers"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

/*
=====================================
 TEST CREATE USER
=====================================
*/

func TestUserService_Create_Success(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := NewUserService(repo, "./uploads")

	req := dto.CreateUserRequest{
		Username: "alice",
		Password: "password123",
		Email:    "alice@example.com",
	}

	resp, err := service.Create(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if resp.Username != "alice" {
		t.Errorf("expected username 'alice', got '%s'", resp.Username)
	}
	if resp.Email != "alice@example.com" {
		t.Errorf("expected email 'alice@example.com', got '%s'", resp.Email)
	}

	// verify password is hashed in repo
	user, err := repo.FindByUsername("alice")
	if err != nil {
		t.Fatalf("expected user in repo, got error %v", err)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123")) != nil {
		t.Error("expected stored password to match provided password after hashing")
	}
}

func TestUserService_Create_InvalidEmail(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := NewUserService(repo, "./uploads")

	req := dto.CreateUserRequest{
		Username: "bob",
		Password: "password123",
		Email:    "not-an-email",
	}

	_, err := service.Create(req)
	if err == nil {
		t.Fatal("expected error for invalid email")
	}
	if err.Error() != "email tidak valid" {
		t.Errorf("expected 'email tidak valid', got '%s'", err.Error())
	}
}

func TestUserService_Create_InvalidUsername(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := NewUserService(repo, "./uploads")

	// Username terlalu pendek (< 3 karakter)
	req := dto.CreateUserRequest{
		Username: "ab",
		Password: "password123",
		Email:    "test@example.com",
	}

	_, err := service.Create(req)
	if err == nil {
		t.Fatal("expected error for invalid username")
	}
	if err.Error() != "username harus 3-50 karakter" {
		t.Errorf("expected 'username harus 3-50 karakter', got '%s'", err.Error())
	}
}

func TestUserService_Create_InvalidPassword(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := NewUserService(repo, "./uploads")

	// Password terlalu pendek (< 6 karakter)
	req := dto.CreateUserRequest{
		Username: "testuser",
		Password: "12345",
		Email:    "test@example.com",
	}

	_, err := service.Create(req)
	if err == nil {
		t.Fatal("expected error for invalid password")
	}
	if err.Error() != "password minimal 6 karakter" {
		t.Errorf("expected 'password minimal 6 karakter', got '%s'", err.Error())
	}
}

func TestUserService_Create_DuplicateUsernameOrEmail(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := NewUserService(repo, "./uploads")

	req := dto.CreateUserRequest{Username: "carol", Password: "password123", Email: "carol@example.com"}
	_, err := service.Create(req)
	if err != nil {
		t.Fatalf("setup create failed: %v", err)
	}

	// duplicate username
	dup := dto.CreateUserRequest{Username: "carol", Password: "newpass123", Email: "other@example.com"}
	_, err = service.Create(dup)
	if err == nil {
		t.Fatal("expected error for duplicate username")
	}
	if err.Error() != "username sudah digunakan" {
		t.Errorf("expected 'username sudah digunakan', got '%s'", err.Error())
	}

	// duplicate email
	dupEmail := dto.CreateUserRequest{Username: "other", Password: "newpass123", Email: "carol@example.com"}
	_, err = service.Create(dupEmail)
	if err == nil {
		t.Fatal("expected error for duplicate email")
	}
	if err.Error() != "email sudah digunakan" {
		t.Errorf("expected 'email sudah digunakan', got '%s'", err.Error())
	}
}

/*
=====================================
 TEST GET ALL USERS
=====================================
*/

func TestUserService_GetAll_Success(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := NewUserService(repo, "./uploads")

	// Create beberapa user
	users := []dto.CreateUserRequest{
		{Username: "user1", Password: "password1", Email: "user1@example.com"},
		{Username: "user2", Password: "password2", Email: "user2@example.com"},
		{Username: "user3", Password: "password3", Email: "user3@example.com"},
	}

	for _, req := range users {
		_, err := service.Create(req)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}
	}

	// Get all users
	result, err := service.GetAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 3 {
		t.Errorf("expected 3 users, got %d", len(result))
	}
}

func TestUserService_GetAll_EmptyResult(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := NewUserService(repo, "./uploads")

	result, err := service.GetAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected 0 users, got %d", len(result))
	}
}

/*
=====================================
 TEST GET BY ID
=====================================
*/

func TestUserService_GetByID_Success(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := NewUserService(repo, "./uploads")

	// Create user
	req := dto.CreateUserRequest{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
	}
	created, err := service.Create(req)
	if err != nil {
		t.Fatalf("setup create failed: %v", err)
	}

	// Get by ID
	result, err := service.GetByID(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != created.ID {
		t.Errorf("expected ID '%s', got '%s'", created.ID, result.ID)
	}
	if result.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", result.Username)
	}
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := NewUserService(repo, "./uploads")

	_, err := service.GetByID("non-existent-id")
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

/*
=====================================
 TEST UPDATE USER
=====================================
*/

func TestUserService_Update_Success_And_Conflict(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := NewUserService(repo, "./uploads")

	// create two users (usernames must be >=3 chars)
	a := dto.CreateUserRequest{Username: "user1", Password: "password1", Email: "u1@example.com"}
	b := dto.CreateUserRequest{Username: "user2", Password: "password2", Email: "u2@example.com"}
	ra, err := service.Create(a)
	if err != nil {
		t.Fatalf("create a failed: %v", err)
	}
	rb, err := service.Create(b)
	if err != nil {
		t.Fatalf("create b failed: %v", err)
	}

	// update a's username to new value
	newUsername := "user1-new"
	updateReq := dto.UpdateUserRequest{Username: &newUsername}
	updated, err := service.Update(ra.ID, updateReq)
	if err != nil {
		t.Fatalf("expected update success, got %v", err)
	}
	if updated.Username != newUsername {
		t.Errorf("expected username '%s', got '%s'", newUsername, updated.Username)
	}

	// try to update a to use b's username -> conflict
	conflictReq := dto.UpdateUserRequest{Username: &rb.Username}
	_, err = service.Update(ra.ID, conflictReq)
	if err == nil {
		t.Fatal("expected error for username already used")
	}
	if err.Error() != "username sudah digunakan" {
		t.Errorf("expected 'username sudah digunakan', got '%s'", err.Error())
	}
}

func TestUserService_Update_InvalidUsername(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := NewUserService(repo, "./uploads")

	// Create user
	req := dto.CreateUserRequest{Username: "testuser", Password: "password123", Email: "test@example.com"}
	created, err := service.Create(req)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// Update dengan username invalid (terlalu pendek)
	invalidUsername := "ab"
	updateReq := dto.UpdateUserRequest{Username: &invalidUsername}
	_, err = service.Update(created.ID, updateReq)
	if err == nil {
		t.Fatal("expected error for invalid username")
	}
	if err.Error() != "username harus 3-50 karakter" {
		t.Errorf("expected 'username harus 3-50 karakter', got '%s'", err.Error())
	}
}

func TestUserService_Update_InvalidEmail(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := NewUserService(repo, "./uploads")

	// Create user
	req := dto.CreateUserRequest{Username: "testuser", Password: "password123", Email: "test@example.com"}
	created, err := service.Create(req)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// Update dengan email invalid
	invalidEmail := "not-an-email"
	updateReq := dto.UpdateUserRequest{Email: &invalidEmail}
	_, err = service.Update(created.ID, updateReq)
	if err == nil {
		t.Fatal("expected error for invalid email")
	}
	if err.Error() != "email tidak valid" {
		t.Errorf("expected 'email tidak valid', got '%s'", err.Error())
	}
}

/*
=====================================
 TEST UPDATE PASSWORD
=====================================
*/

func TestUserService_UpdatePassword_Success(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := NewUserService(repo, "./uploads")

	// Create user
	req := dto.CreateUserRequest{Username: "testuser", Password: "oldpassword", Email: "test@example.com"}
	created, err := service.Create(req)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// Update password
	updateReq := dto.UpdateUserPasswordRequest{
		OldPassword: "oldpassword",
		NewPassword: "newpassword123",
	}
	err = service.UpdatePassword(created.ID, updateReq)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify password changed
	password, err := repo.GetPasswordByID(created.ID)
	if err != nil {
		t.Fatalf("failed to get password: %v", err)
	}

	if bcrypt.CompareHashAndPassword([]byte(password), []byte("newpassword123")) != nil {
		t.Error("expected password to be updated")
	}
}

func TestUserService_UpdatePassword_WrongOldPassword(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := NewUserService(repo, "./uploads")

	// Create user
	req := dto.CreateUserRequest{Username: "testuser", Password: "oldpassword", Email: "test@example.com"}
	created, err := service.Create(req)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// Update password dengan old password salah
	updateReq := dto.UpdateUserPasswordRequest{
		OldPassword: "wrongpassword",
		NewPassword: "newpassword123",
	}
	err = service.UpdatePassword(created.ID, updateReq)
	if err == nil {
		t.Fatal("expected error for wrong old password")
	}
	if err.Error() != "password lama tidak sesuai" {
		t.Errorf("expected 'password lama tidak sesuai', got '%s'", err.Error())
	}
}

func TestUserService_UpdatePassword_TooShort(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := NewUserService(repo, "./uploads")

	// Create user
	req := dto.CreateUserRequest{Username: "testuser", Password: "oldpassword", Email: "test@example.com"}
	created, err := service.Create(req)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// Update password dengan new password terlalu pendek
	updateReq := dto.UpdateUserPasswordRequest{
		OldPassword: "oldpassword",
		NewPassword: "12345", // < 6 karakter
	}
	err = service.UpdatePassword(created.ID, updateReq)
	if err == nil {
		t.Fatal("expected error for password too short")
	}
	if err.Error() != "password baru minimal 6 karakter" {
		t.Errorf("expected 'password baru minimal 6 karakter', got '%s'", err.Error())
	}
}

/*
=====================================
 TEST DELETE USER
=====================================
*/

func TestUserService_Delete_Success_And_NotFound(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := NewUserService(repo, "./uploads")

	req := dto.CreateUserRequest{Username: "todelete", Password: "password", Email: "td@example.com"}
	created, err := service.Create(req)
	if err != nil {
		t.Fatalf("setup create failed: %v", err)
	}

	// delete
	if err := service.Delete(created.ID); err != nil {
		t.Fatalf("expected delete success, got %v", err)
	}

	// ensure gone
	_, err = repo.FindByID(created.ID)
	if err == nil {
		t.Fatal("expected user to be deleted from repo")
	}

	// delete non-existent
	if err := service.Delete("nope"); err == nil {
		t.Fatal("expected error when deleting non-existent user")
	}
}