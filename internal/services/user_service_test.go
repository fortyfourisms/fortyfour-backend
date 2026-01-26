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
func TestUserService_GetByID_Success(t *testing.T) {
	service, mockRepo := setupUserService()

	user := testhelpers.CreateTestUser("test-id", "testuser", "test@test.com")
	_ = mockRepo.Create(user)

	result, err := service.GetByID("test-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil || result.ID != "test-id" {
		t.Fatal("expected valid user")
	}
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	service, _ := setupUserService()

	_, err := service.GetByID("not-found")
	if err == nil {
		t.Error("expected error for nonexistent user")
	}
}

func TestUserService_Create_Success(t *testing.T) {
	service, _ := setupUserService()

	req := dto.CreateUserRequest{
		Username: "newuser",
		Password: "P@sjord13!",
		Email:    "newuser@test.com",
	}

	result, err := service.Create(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Username != "newuser" {
		t.Errorf("expected username 'newuser', got '%s'", result.Username)
	}
}

func TestUserService_Create_DuplicateUsername(t *testing.T) {
	service, mockRepo := setupUserService()

	user := testhelpers.CreateTestUser("id1", "existinguser", "exist@test.com")
	_ = mockRepo.Create(user)

	req := dto.CreateUserRequest{
		Username: "existinguser",
		Password: "An0ther#Strong_2026!",
		Email:    "new@test.com",
	}

	_, err := service.Create(req)
	if err == nil {
		t.Error("expected error for duplicate username")
	}
}

func TestUserService_Create_DuplicateEmail(t *testing.T) {
	service, mockRepo := setupUserService()

	user := testhelpers.CreateTestUser("id1", "user1", "exist@test.com")
	_ = mockRepo.Create(user)

	req := dto.CreateUserRequest{
		Username: "newuser",
		Password: "An0ther#Strong_2026!",
		Email:    "exist@test.com",
	}

	_, err := service.Create(req)
	if err == nil {
		t.Error("expected error for duplicate email")
	}
}
func TestUserService_Update_Success(t *testing.T) {
	service, mockRepo := setupUserService()

	// ✅ pastikan ID konsisten
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
		t.Errorf("expected username updated")
	}
}

func TestUserService_UpdatePassword_Success(t *testing.T) {
	service, mockRepo := setupUserService()

	// Create user with hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Corr3tB4e!"), bcrypt.DefaultCost)
	user := testhelpers.CreateTestUser("test-id", "testuser", "test@test.com")
	user.Password = string(hashedPassword)
	mockRepo.Create(user)

	req := dto.UpdateUserPasswordRequest{
		OldPassword:        "Corr3tB4e!",
		NewPassword:        "Corr3tB4e!!",
		ConfirmNewPassword: "Corr3tB4e!!",
	}

	err := service.UpdatePassword("test-id", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
func TestUserService_UpdatePassword_WrongOldPassword(t *testing.T) {
	service, mockRepo := setupUserService()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Corr3tB4e!"), bcrypt.DefaultCost)
	user := testhelpers.CreateTestUser("test-id", "testuser", "test@test.com")
	user.Password = string(hashedPassword)
	mockRepo.Create(user)

	req := dto.UpdateUserPasswordRequest{
		OldPassword:        "Corr3tB4e!!",
		NewPassword:        "Corr3tB4e#!",
		ConfirmNewPassword: "Corr3tB4e#!",
	}

	err := service.UpdatePassword("test-id", req)
	if err == nil {
		t.Error("expected error for wrong old password")
	}
}

func TestUserService_Delete_Success(t *testing.T) {
	service, mockRepo := setupUserService()

	user := testhelpers.CreateTestUser("test-id", "user", "user@test.com")
	_ = mockRepo.Create(user)

	err := service.Delete("test-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = mockRepo.FindByID("test-id")
	if err == nil {
		t.Error("user should be deleted")
	}
}
