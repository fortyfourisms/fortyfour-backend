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
// MOCK PROTEKSI REPOSITORY
// ===============================
//

type mockProteksiRepository struct {
	CreateFn  func(req dto.CreateProteksiRequest, id string) error
	GetAllFn  func() ([]models.Proteksi, error)
	GetByIDFn func(id string) (*models.Proteksi, error)
	UpdateFn  func(id string, proteksi models.Proteksi) error
	DeleteFn  func(id string) error
}

func (m *mockProteksiRepository) Create(req dto.CreateProteksiRequest, id string) error {
	return m.CreateFn(req, id)
}

func (m *mockProteksiRepository) GetAll() ([]models.Proteksi, error) {
	return m.GetAllFn()
}

func (m *mockProteksiRepository) GetByID(id string) (*models.Proteksi, error) {
	return m.GetByIDFn(id)
}

func (m *mockProteksiRepository) Update(id string, proteksi models.Proteksi) error {
	return m.UpdateFn(id, proteksi)
}

func (m *mockProteksiRepository) Delete(id string) error {
	return m.DeleteFn(id)
}

// compile-time safety check
var _ repository.ProteksiRepositoryInterface = (*mockProteksiRepository)(nil)

//
// ===============================
// TEST CREATE
// ===============================
//

func TestProteksiService_Create_Success(t *testing.T) {
	repo := &mockProteksiRepository{
		CreateFn: func(req dto.CreateProteksiRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*models.Proteksi, error) {
			return &models.Proteksi{
				ID:            id,
				NilaiProteksi: 80,
			}, nil
		},
	}

	service := NewProteksiService(repo)

	req := dto.CreateProteksiRequest{
		NilaiProteksi: 80,
	}

	result, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 80.0, result.NilaiProteksi)
}

func TestProteksiService_Create_Error(t *testing.T) {
	repo := &mockProteksiRepository{
		CreateFn: func(req dto.CreateProteksiRequest, id string) error {
			return errors.New("gagal create")
		},
	}

	service := NewProteksiService(repo)

	result, err := service.Create(dto.CreateProteksiRequest{})

	assert.Error(t, err)
	assert.Nil(t, result)
}

//
// ===============================
// TEST GET ALL
// ===============================
//

func TestProteksiService_GetAll_Success(t *testing.T) {
	repo := &mockProteksiRepository{
		GetAllFn: func() ([]models.Proteksi, error) {
			return []models.Proteksi{
				{ID: "1", NilaiProteksi: 70},
				{ID: "2", NilaiProteksi: 80},
			}, nil
		},
	}

	service := NewProteksiService(repo)

	result, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

//
// ===============================
// TEST GET BY ID
// ===============================
//

func TestProteksiService_GetByID_Success(t *testing.T) {
	repo := &mockProteksiRepository{
		GetByIDFn: func(id string) (*models.Proteksi, error) {
			return &models.Proteksi{
				ID:            id,
				NilaiProteksi: 90,
			}, nil
		},
	}

	service := NewProteksiService(repo)

	result, err := service.GetByID("uuid-test")

	assert.NoError(t, err)
	assert.Equal(t, 90.0, result.NilaiProteksi)
}

func TestProteksiService_GetByID_NotFound(t *testing.T) {
	repo := &mockProteksiRepository{
		GetByIDFn: func(id string) (*models.Proteksi, error) {
			return nil, errors.New("data tidak ditemukan")
		},
	}

	service := NewProteksiService(repo)

	result, err := service.GetByID("invalid-id")

	assert.Error(t, err)
	assert.Nil(t, result)
}

//
// ===============================
// TEST UPDATE
// ===============================
//

func TestProteksiService_Update_Success(t *testing.T) {
	nilaiBaru := 95.0

	repo := &mockProteksiRepository{
		GetByIDFn: func(id string) (*models.Proteksi, error) {
			return &models.Proteksi{
				ID:            id,
				NilaiProteksi: 80,
			}, nil
		},
		UpdateFn: func(id string, proteksi models.Proteksi) error {
			return nil
		},
	}

	service := NewProteksiService(repo)

	req := dto.UpdateProteksiRequest{
		NilaiProteksi: &nilaiBaru,
	}

	result, err := service.Update("uuid-test", req)

	assert.NoError(t, err)
	assert.Equal(t, nilaiBaru, result.NilaiProteksi)
}

//
// ===============================
// TEST DELETE
// ===============================
//

func TestProteksiService_Delete_Success(t *testing.T) {
	repo := &mockProteksiRepository{
		DeleteFn: func(id string) error {
			return nil
		},
	}

	service := NewProteksiService(repo)

	err := service.Delete("uuid-test")

	assert.NoError(t, err)
}
