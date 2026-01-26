package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/testhelpers"
	"os"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func setupUserService() (*UserService, *testhelpers.MockUserRepository) {
	mockRepo := testhelpers.NewMockUserRepository()
	uploadPath := "./test_uploads"
	_ = os.MkdirAll(uploadPath, os.ModePerm)
	service := NewUserService(mockRepo, uploadPath)
	return service, mockRepo
}

/*
=====================================
 TEST GET ALL USERS
=====================================
*/

func TestUserService_GetAll_Success(t *testing.T) {
	service, mockRepo := setupUserService()

	user1 := testhelpers.CreateTestUser("id1", "user1", "user1@test.com")
	user2 := testhelpers.CreateTestUser("id2", "user2", "user2@test.com")
	_ = mockRepo.Create(user1)
	_ = mockRepo.Create(user2)

	result, err := service.GetAll()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 users, got %d", len(result))
	}
}

func TestUserService_GetAll_EmptyResult(t *testing.T) {
	service, _ := setupUserService()

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
	service, mockRepo := setupUserService()

	user := testhelpers.CreateTestUser("test-id", "testuser", "test@test.com")
	_ = mockRepo.Create(user)

	result, err := service.GetByID("test-id")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.ID != "test-id" {
		t.Errorf("expected ID 'test-id', got '%s'", result.ID)
	}
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	service, _ := setupUserService()

	_, err := service.GetByID("nonexistent")

	if err == nil {
		t.Error("expected error for nonexistent user")
	}
}

/*
=====================================
 TEST CREATE USER
=====================================
*/

func TestUserService_Create_Success(t *testing.T) {
	service, _ := setupUserService()

	req := dto.CreateUserRequest{
		Username: "newuser",
		Password: "MySecureP@ssw0rd2024!",
		Email:    "newuser@test.com",
	}

	result, err := service.Create(req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.Username != "newuser" {
		t.Errorf("expected username 'newuser', got '%s'", result.Username)
	}
	if result.Email != "newuser@test.com" {
		t.Errorf("expected email 'newuser@test.com', got '%s'", result.Email)
	}

	// verify password is hashed
	user, _ := service.repo.FindByUsername("newuser")
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("MySecureP@ssw0rd2024!")) != nil {
		t.Error("expected stored password to match provided password after hashing")
	}
}

func TestUserService_Create_InvalidEmail(t *testing.T) {
	service, _ := setupUserService()

	req := dto.CreateUserRequest{
		Username: "testuser",
		Password: "MySecureP@ssw0rd2024!",
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
	service, _ := setupUserService()

	req := dto.CreateUserRequest{
		Username: "ab", // < 3 karakter
		Password: "MySecureP@ssw0rd2024!",
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
	service, _ := setupUserService()

	req := dto.CreateUserRequest{
		Username: "testuser",
		Password: "12345", // terlalu pendek
		Email:    "test@example.com",
	}

	_, err := service.Create(req)
	if err == nil {
		t.Fatal("expected error for invalid password")
	}
	// Sesuaikan dengan error message yang sebenarnya
	if err.Error() != "password minimal 8 karakter" {
		t.Errorf("expected 'password minimal 8 karakter', got '%s'", err.Error())
	}
}

func TestUserService_Create_DuplicateUsername(t *testing.T) {
	service, mockRepo := setupUserService()

	user := testhelpers.CreateTestUser("id1", "existinguser", "existing@test.com")
	_ = mockRepo.Create(user)

	req := dto.CreateUserRequest{
		Username: "existinguser",
		Password: "MySecureP@ssw0rd2024!",
		Email:    "new@test.com",
	}

	_, err := service.Create(req)

	if err == nil {
		t.Error("expected error for duplicate username")
	}
	if err.Error() != "username sudah digunakan" {
		t.Errorf("expected 'username sudah digunakan', got '%s'", err.Error())
	}
}

func TestUserService_Create_DuplicateEmail(t *testing.T) {
	service, mockRepo := setupUserService()

	user := testhelpers.CreateTestUser("id1", "user1", "existing@test.com")
	_ = mockRepo.Create(user)

	req := dto.CreateUserRequest{
		Username: "newuser",
		Password: "MySecureP@ssw0rd2024!",
		Email:    "existing@test.com",
	}

	_, err := service.Create(req)

	if err == nil {
		t.Error("expected error for duplicate email")
	}
	if err.Error() != "email sudah digunakan" {
		t.Errorf("expected 'email sudah digunakan', got '%s'", err.Error())
	}
}

/*
=====================================
 TEST UPDATE USER
=====================================
*/

func TestUserService_Update_Success(t *testing.T) {
	service, mockRepo := setupUserService()

	user := testhelpers.CreateTestUser("test-id", "olduser", "old@test.com")
	err := mockRepo.Create(user)
	if err != nil {
		t.Fatalf("failed setup user: %v", err)
	}

	newUsername := "newuser"
	newEmail := "new@test.com"
	req := dto.UpdateUserRequest{
		Username: &newUsername,
		Email:    &newEmail,
	}

	result, err := service.Update("test-id", req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Username != "newuser" {
		t.Errorf("expected username 'newuser', got '%s'", result.Username)
	}
	if result.Email != "new@test.com" {
		t.Errorf("expected email 'new@test.com', got '%s'", result.Email)
	}
}

func TestUserService_Update_Conflict(t *testing.T) {
	service, mockRepo := setupUserService()

	// create two users
	a := testhelpers.CreateTestUser("id-a", "user1", "u1@example.com")
	b := testhelpers.CreateTestUser("id-b", "user2", "u2@example.com")
	_ = mockRepo.Create(a)
	_ = mockRepo.Create(b)

	// try to update a to use b's username -> conflict
	conflictReq := dto.UpdateUserRequest{Username: &b.Username}
	_, err := service.Update(a.ID, conflictReq)
	if err == nil {
		t.Fatal("expected error for username already used")
	}
	if err.Error() != "username sudah digunakan" {
		t.Errorf("expected 'username sudah digunakan', got '%s'", err.Error())
	}
}

func TestUserService_Update_InvalidUsername(t *testing.T) {
	service, mockRepo := setupUserService()

	user := testhelpers.CreateTestUser("test-id", "testuser", "test@example.com")
	_ = mockRepo.Create(user)

	invalidUsername := "ab"
	updateReq := dto.UpdateUserRequest{Username: &invalidUsername}
	_, err := service.Update("test-id", updateReq)
	if err == nil {
		t.Fatal("expected error for invalid username")
	}
	if err.Error() != "username harus 3-50 karakter" {
		t.Errorf("expected 'username harus 3-50 karakter', got '%s'", err.Error())
	}
}

func TestUserService_Update_InvalidEmail(t *testing.T) {
	service, mockRepo := setupUserService()

	user := testhelpers.CreateTestUser("test-id", "testuser", "test@example.com")
	_ = mockRepo.Create(user)

	invalidEmail := "not-an-email"
	updateReq := dto.UpdateUserRequest{Email: &invalidEmail}
	_, err := service.Update("test-id", updateReq)
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
	service, mockRepo := setupUserService()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("OldP@ssword123!"), bcrypt.DefaultCost)
	user := testhelpers.CreateTestUser("test-id", "testuser", "test@test.com")
	user.Password = string(hashedPassword)
	_ = mockRepo.Create(user)

	req := dto.UpdateUserPasswordRequest{
		OldPassword:        "OldP@ssword123!",
		NewPassword:        "NewP@ssword456!",
		ConfirmNewPassword: "NewP@ssword456!",
	}

	err := service.UpdatePassword("test-id", req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify password changed
	password, err := mockRepo.GetPasswordByID("test-id")
	if err != nil {
		t.Fatalf("failed to get password: %v", err)
	}

	if bcrypt.CompareHashAndPassword([]byte(password), []byte("NewP@ssword456!")) != nil {
		t.Error("expected password to be updated")
	}
}

func TestUserService_UpdatePassword_WrongOldPassword(t *testing.T) {
	service, mockRepo := setupUserService()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("CorrectP@ssword123!"), bcrypt.DefaultCost)
	user := testhelpers.CreateTestUser("test-id", "testuser", "test@test.com")
	user.Password = string(hashedPassword)
	_ = mockRepo.Create(user)

	req := dto.UpdateUserPasswordRequest{
		OldPassword:        "WrongP@ssword123!",
		NewPassword:        "NewP@ssword456!",
		ConfirmNewPassword: "NewP@ssword456!",
	}

	err := service.UpdatePassword("test-id", req)

	if err == nil {
		t.Error("expected error for wrong old password")
	}
	if err.Error() != "password lama tidak sesuai" {
		t.Errorf("expected 'password lama tidak sesuai', got '%s'", err.Error())
	}
}

func TestUserService_UpdatePassword_TooShort(t *testing.T) {
	service, mockRepo := setupUserService()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("OldP@ssword123!"), bcrypt.DefaultCost)
	user := testhelpers.CreateTestUser("test-id", "testuser", "test@example.com")
	user.Password = string(hashedPassword)
	_ = mockRepo.Create(user)

	req := dto.UpdateUserPasswordRequest{
		OldPassword:        "OldP@ssword123!",
		NewPassword:        "Short1!", // kurang dari 8 karakter
		ConfirmNewPassword: "Short1!",
	}
	err := service.UpdatePassword("test-id", req)
	if err == nil {
		t.Fatal("expected error for password too short")
	}
	// Sesuaikan dengan validasi yang ada
	if err.Error() != "password minimal 8 karakter" {
		t.Errorf("expected 'password minimal 8 karakter', got '%s'", err.Error())
	}
}

/*
=====================================
 TEST DELETE USER
=====================================
*/

func TestUserService_Delete_Success(t *testing.T) {
	service, mockRepo := setupUserService()

	user := testhelpers.CreateTestUser("test-id", "testuser", "test@test.com")
	_ = mockRepo.Create(user)

	err := service.Delete("test-id")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify it was deleted
	_, err = mockRepo.FindByID("test-id")
	if err == nil {
		t.Error("user should be deleted")
	}
}

func TestUserService_Delete_NotFound(t *testing.T) {
	service, _ := setupUserService()

	err := service.Delete("nonexistent-id")

	if err == nil {
		t.Error("expected error for nonexistent user")
	}
}