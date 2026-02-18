package services

import (
	"context"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/testhelpers"
	"testing"
	"time"

	"github.com/google/uuid"
)

func setupRoleService() (*RoleService, repository.RoleRepository) {
	mockRepo := testhelpers.NewMockRoleRepository()
	service := NewRoleService(mockRepo, nil)
	return service, mockRepo
}

func TestRoleService_Create_Success(t *testing.T) {
	service, mockRepo := setupRoleService()

	req := dto.CreateRoleRequest{
		Name:        "test-role",
		Description: "Test Description",
	}

	result, err := service.Create(req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.Name != "test-role" {
		t.Errorf("expected name 'test-role', got '%s'", result.Name)
	}

	// Verify it was created
	created, _ := mockRepo.GetByID(context.Background(), result.ID)
	if created == nil {
		t.Error("role should be created in repository")
	}
}

func TestRoleService_Create_DuplicateName(t *testing.T) {
	service, mockRepo := setupRoleService()

	// Create first role
	role1 := &models.Role{
		ID:          uuid.New().String(),
		Name:        "existing-role",
		Description: "Existing",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockRepo.Create(context.Background(), role1)

	// Try to create duplicate
	req := dto.CreateRoleRequest{
		Name:        "existing-role",
		Description: "Duplicate",
	}

	_, err := service.Create(req)

	if err == nil {
		t.Error("expected error for duplicate name")
	}
	if err.Error() != "role name already exists" {
		t.Errorf("expected 'role name already exists', got '%s'", err.Error())
	}
}

func TestRoleService_GetByID_Success(t *testing.T) {
	service, mockRepo := setupRoleService()

	role := &models.Role{
		ID:          "test-id",
		Name:        "test-role",
		Description: "Test Description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockRepo.Create(context.Background(), role)

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

func TestRoleService_GetByID_NotFound(t *testing.T) {
	service, _ := setupRoleService()

	_, err := service.GetByID("nonexistent")

	if err == nil {
		t.Error("expected error for nonexistent role")
	}
}

func TestRoleService_GetAll_Success(t *testing.T) {
	service, mockRepo := setupRoleService()

	role1 := &models.Role{
		ID:          uuid.New().String(),
		Name:        "role1",
		Description: "Role 1",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	role2 := &models.Role{
		ID:          uuid.New().String(),
		Name:        "role2",
		Description: "Role 2",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockRepo.Create(context.Background(), role1)
	mockRepo.Create(context.Background(), role2)

	result, err := service.GetAll()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 roles, got %d", len(result))
	}
}

func TestRoleService_Update_Success(t *testing.T) {
	service, mockRepo := setupRoleService()

	role := &models.Role{
		ID:          "test-id",
		Name:        "old-name",
		Description: "Old Description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockRepo.Create(context.Background(), role)

	req := dto.UpdateRoleRequest{
		Name:        "new-name",
		Description: "New Description",
	}

	result, err := service.Update("test-id", req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Name != "new-name" {
		t.Errorf("expected name 'new-name', got '%s'", result.Name)
	}
}

func TestRoleService_Update_NotFound(t *testing.T) {
	service, _ := setupRoleService()

	req := dto.UpdateRoleRequest{
		Name: "new-name",
	}

	_, err := service.Update("nonexistent", req)

	if err == nil {
		t.Error("expected error for nonexistent role")
	}
}

func TestRoleService_Delete_Success(t *testing.T) {
	service, mockRepo := setupRoleService()

	role := &models.Role{
		ID:          "test-id",
		Name:        "test-role",
		Description: "Test",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockRepo.Create(context.Background(), role)

	err := service.Delete("test-id")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify it was deleted
	_, err = mockRepo.GetByID(context.Background(), "test-id")
	if err == nil {
		t.Error("role should be deleted")
	}
}

func TestRoleService_Delete_NotFound(t *testing.T) {
	service, _ := setupRoleService()

	err := service.Delete("nonexistent")

	if err == nil {
		t.Error("expected error for nonexistent role")
	}
}
