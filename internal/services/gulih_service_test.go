package services

import (
	"errors"
	"testing"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"

	"github.com/stretchr/testify/assert"
)

//
// ===============================
// MOCK GULIH REPOSITORY
// ===============================
//

type mockGulihRepository struct {
	CreateFn  func(req dto.CreateGulihRequest, id string) error
	GetAllFn  func() ([]models.Gulih, error)
	GetByIDFn func(id string) (*models.Gulih, error)
	UpdateFn  func(id string, gulih models.Gulih) error
	DeleteFn  func(id string) error
}

func (m *mockGulihRepository) Create(req dto.CreateGulihRequest, id string) error {
	return m.CreateFn(req, id)
}

func (m *mockGulihRepository) GetAll() ([]models.Gulih, error) {
	return m.GetAllFn()
}

func (m *mockGulihRepository) GetByID(id string) (*models.Gulih, error) {
	return m.GetByIDFn(id)
}

func (m *mockGulihRepository) Update(id string, gulih models.Gulih) error {
	return m.UpdateFn(id, gulih)
}

func (m *mockGulihRepository) Delete(id string) error {
	return m.DeleteFn(id)
}

// compile-time safety check
var _ repository.GulihRepositoryInterface = (*mockGulihRepository)(nil)

//
// ===============================
// TEST CREATE
// ===============================
//

func TestGulihService_Create_Success(t *testing.T) {
	repo := &mockGulihRepository{
		CreateFn: func(req dto.CreateGulihRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*models.Gulih, error) {
			return &models.Gulih{
				ID:         id,
				NilaiGulih: 88,
			}, nil
		},
	}

	service := NewGulihService(repo)

	req := dto.CreateGulihRequest{
		NilaiGulih: 88,
	}

	result, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 88.0, result.NilaiGulih)
}

func TestGulihService_Create_Error(t *testing.T) {
	repo := &mockGulihRepository{
		CreateFn: func(req dto.CreateGulihRequest, id string) error {
			return errors.New("gagal create")
		},
	}

	service := NewGulihService(repo)

	result, err := service.Create(dto.CreateGulihRequest{})

	assert.Error(t, err)
	assert.Nil(t, result)
}

//
// ===============================
// TEST GET ALL
// ===============================
//

func TestGulihService_GetAll_Success(t *testing.T) {
	repo := &mockGulihRepository{
		GetAllFn: func() ([]models.Gulih, error) {
			return []models.Gulih{
				{ID: "1", NilaiGulih: 70},
				{ID: "2", NilaiGulih: 80},
			}, nil
		},
	}

	service := NewGulihService(repo)

	result, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

//
// ===============================
// TEST GET BY ID
// ===============================
//

func TestGulihService_GetByID_Success(t *testing.T) {
	repo := &mockGulihRepository{
		GetByIDFn: func(id string) (*models.Gulih, error) {
			return &models.Gulih{
				ID:         id,
				NilaiGulih: 90,
			}, nil
		},
	}

	service := NewGulihService(repo)

	result, err := service.GetByID("uuid-test")

	assert.NoError(t, err)
	assert.Equal(t, 90.0, result.NilaiGulih)
}

func TestGulihService_GetByID_NotFound(t *testing.T) {
	repo := &mockGulihRepository{
		GetByIDFn: func(id string) (*models.Gulih, error) {
			return nil, errors.New("data tidak ditemukan")
		},
	}

	service := NewGulihService(repo)

	result, err := service.GetByID("invalid-id")

	assert.Error(t, err)
	assert.Nil(t, result)
}

//
// ===============================
// TEST UPDATE
// ===============================
//

func TestGulihService_Update_Success(t *testing.T) {
	nilaiBaru := 95.0

	repo := &mockGulihRepository{
		GetByIDFn: func(id string) (*models.Gulih, error) {
			return &models.Gulih{
				ID:         id,
				NilaiGulih: 80,
			}, nil
		},
		UpdateFn: func(id string, gulih models.Gulih) error {
			return nil
		},
	}

	service := NewGulihService(repo)

	req := dto.UpdateGulihRequest{
		NilaiGulih: &nilaiBaru,
	}

	result, err := service.Update("uuid-test", req)

	assert.NoError(t, err)
	assert.Equal(t, nilaiBaru, result.NilaiGulih)
}

//
// ===============================
// TEST DELETE
// ===============================
//

func TestGulihService_Delete_Success(t *testing.T) {
	repo := &mockGulihRepository{
		DeleteFn: func(id string) error {
			return nil
		},
	}

	service := NewGulihService(repo)

	err := service.Delete("uuid-test")

	assert.NoError(t, err)
}
