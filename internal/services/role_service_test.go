package services

import (
	"context"
	"encoding/json"
	"errors"
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

// ============================================================
// TestRoleService_Create — tambahan
// ============================================================

func TestRoleService_Create_InvalidatesCache(t *testing.T) {
	rc := newRoleTestRedis()
	mockRepo := testhelpers.NewMockRoleRepository()
	service := NewRoleService(mockRepo, rc)

	// Pre-populate list cache
	setRoleCache(rc, keyList("role"), []*dto.RoleResponse{{ID: "lama", Name: "Lama"}})

	req := dto.CreateRoleRequest{Name: "role-baru", Description: "Deskripsi"}
	_, err := service.Create(req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Cache list harus sudah dihapus
	exists, _ := rc.Exists(keyList("role"))
	if exists {
		t.Error("cache list harus dihapus setelah Create")
	}
}

// ============================================================
// TestRoleService_GetByID — cache hit/miss
// ============================================================

func TestRoleService_GetByID_CacheHit_SkipRepo(t *testing.T) {
	rc := newRoleTestRedis()
	mockRepo := testhelpers.NewMockRoleRepository()
	service := NewRoleService(mockRepo, rc)

	cached := dto.RoleResponse{ID: "role-1", Name: "Dari Cache"}
	setRoleCache(rc, keyDetail("role", "role-1"), cached)

	result, err := service.GetByID("role-1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Name != "Dari Cache" {
		t.Errorf("expected 'Dari Cache', got '%s'", result.Name)
	}
}

func TestRoleService_GetByID_CacheMiss_SetsCache(t *testing.T) {
	rc := newRoleTestRedis()
	mockRepo := testhelpers.NewMockRoleRepository()
	service := NewRoleService(mockRepo, rc)

	role := &models.Role{ID: "role-db", Name: "Dari DB", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	mockRepo.Create(context.Background(), role)

	_, err := service.GetByID("role-db")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	exists, _ := rc.Exists(keyDetail("role", "role-db"))
	if !exists {
		t.Error("data harus di-cache setelah GetByID")
	}
}

// ============================================================
// TestRoleService_GetAll — tambahan
// ============================================================

func TestRoleService_GetAll_Empty(t *testing.T) {
	service, _ := setupRoleService()

	result, err := service.GetAll()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty list, got %d", len(result))
	}
}

func TestRoleService_GetAll_CacheHit_SkipRepo(t *testing.T) {
	rc := newRoleTestRedis()
	mockRepo := testhelpers.NewMockRoleRepository()
	service := NewRoleService(mockRepo, rc)

	cached := []*dto.RoleResponse{{ID: "c1", Name: "Cache Role"}}
	setRoleCache(rc, keyList("role"), cached)

	// Tambahkan role ke repo — tidak boleh diakses kalau cache hit
	mockRepo.Create(context.Background(), &models.Role{
		ID: "db-1", Name: "DB Role", CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})

	result, err := service.GetAll()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 || result[0].Name != "Cache Role" {
		t.Errorf("expected data dari cache, got %v", result)
	}
}

func TestRoleService_GetAll_CacheMiss_SetsCache(t *testing.T) {
	rc := newRoleTestRedis()
	mockRepo := testhelpers.NewMockRoleRepository()
	service := NewRoleService(mockRepo, rc)

	mockRepo.Create(context.Background(), &models.Role{
		ID: "r1", Name: "Role 1", CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})

	_, err := service.GetAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	exists, _ := rc.Exists(keyList("role"))
	if !exists {
		t.Error("data harus di-cache setelah GetAll")
	}
}

// ============================================================
// TestRoleService_Update — tambahan
// ============================================================

func TestRoleService_Update_DuplicateName_Error(t *testing.T) {
	service, mockRepo := setupRoleService()

	// Buat dua role
	mockRepo.Create(context.Background(), &models.Role{
		ID: "role-a", Name: "role-a", CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})
	mockRepo.Create(context.Background(), &models.Role{
		ID: "role-b", Name: "role-b", CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})

	// Coba ubah nama role-a menjadi nama role-b yang sudah ada
	_, err := service.Update("role-a", dto.UpdateRoleRequest{Name: "role-b"})

	if err == nil {
		t.Fatal("expected error untuk nama duplikat")
	}
	if err.Error() != "role name already exists" {
		t.Errorf("expected 'role name already exists', got '%s'", err.Error())
	}
}

func TestRoleService_Update_SameName_NoError(t *testing.T) {
	service, mockRepo := setupRoleService()

	mockRepo.Create(context.Background(), &models.Role{
		ID: "role-1", Name: "existing-name", CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})

	// Update dengan nama yang sama → tidak dianggap duplikat
	result, err := service.Update("role-1", dto.UpdateRoleRequest{
		Name:        "existing-name",
		Description: "Deskripsi baru",
	})

	if err != nil {
		t.Fatalf("expected no error untuk nama yang sama, got %v", err)
	}
	if result.Description != "Deskripsi baru" {
		t.Errorf("expected description 'Deskripsi baru', got '%s'", result.Description)
	}
}

func TestRoleService_Update_HanyaDescription(t *testing.T) {
	service, mockRepo := setupRoleService()

	mockRepo.Create(context.Background(), &models.Role{
		ID: "role-1", Name: "nama-lama", Description: "desc lama",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})

	result, err := service.Update("role-1", dto.UpdateRoleRequest{Description: "desc baru"})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Nama tidak berubah
	if result.Name != "nama-lama" {
		t.Errorf("expected nama tetap 'nama-lama', got '%s'", result.Name)
	}
	if result.Description != "desc baru" {
		t.Errorf("expected description 'desc baru', got '%s'", result.Description)
	}
}

func TestRoleService_Update_InvalidatesCache(t *testing.T) {
	rc := newRoleTestRedis()
	mockRepo := testhelpers.NewMockRoleRepository()
	service := NewRoleService(mockRepo, rc)

	mockRepo.Create(context.Background(), &models.Role{
		ID: "role-1", Name: "lama", CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})

	// Pre-populate cache
	setRoleCache(rc, keyDetail("role", "role-1"), dto.RoleResponse{ID: "role-1", Name: "lama"})
	setRoleCache(rc, keyList("role"), []*dto.RoleResponse{{ID: "role-1", Name: "lama"}})

	_, err := service.Update("role-1", dto.UpdateRoleRequest{Name: "baru"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	existsDetail, _ := rc.Exists(keyDetail("role", "role-1"))
	existsList, _ := rc.Exists(keyList("role"))
	if existsDetail {
		t.Error("cache detail harus dihapus setelah update")
	}
	if existsList {
		t.Error("cache list harus dihapus setelah update")
	}
}

// ============================================================
// TestRoleService_Delete — tambahan
// ============================================================

func TestRoleService_Delete_InvalidatesCache(t *testing.T) {
	rc := newRoleTestRedis()
	mockRepo := testhelpers.NewMockRoleRepository()
	service := NewRoleService(mockRepo, rc)

	mockRepo.Create(context.Background(), &models.Role{
		ID: "role-1", Name: "hapus-aku", CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})

	setRoleCache(rc, keyDetail("role", "role-1"), dto.RoleResponse{ID: "role-1"})
	setRoleCache(rc, keyList("role"), []*dto.RoleResponse{{ID: "role-1"}})

	err := service.Delete("role-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	existsDetail, _ := rc.Exists(keyDetail("role", "role-1"))
	existsList, _ := rc.Exists(keyList("role"))
	if existsDetail {
		t.Error("cache detail harus dihapus setelah delete")
	}
	if existsList {
		t.Error("cache list harus dihapus setelah delete")
	}
}

// ============================================================
// Helpers untuk test Redis pada role_service_test.go
// ============================================================

func newRoleTestRedis() *roleTestRedis {
	return &roleTestRedis{data: make(map[string]string)}
}

func setRoleCache(rc *roleTestRedis, key string, value interface{}) {
	b, _ := json.Marshal(value)
	rc.data[key] = string(b)
}

type roleTestRedis struct {
	data map[string]string
}

func (r *roleTestRedis) Set(key string, value interface{}, ttl time.Duration) error {
	if v, ok := value.(string); ok {
		r.data[key] = v
	}
	return nil
}
func (r *roleTestRedis) Get(key string) (string, error) {
	v, ok := r.data[key]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}
func (r *roleTestRedis) Delete(key string) error {
	delete(r.data, key)
	return nil
}
func (r *roleTestRedis) Exists(key string) (bool, error) {
	_, ok := r.data[key]
	return ok, nil
}
func (r *roleTestRedis) Scan(pattern string) ([]string, error) { return nil, nil }
func (r *roleTestRedis) Close() error                          { return nil }
