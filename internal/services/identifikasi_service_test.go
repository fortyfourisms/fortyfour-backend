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
// MOCK IDENTIFIKASI REPOSITORY
// ===============================
//

type mockIdentifikasiRepository struct {
	CreateFn  func(req dto.CreateIdentifikasiRequest, id string) error
	GetAllFn  func() ([]models.Identifikasi, error)
	GetByIDFn func(id string) (*models.Identifikasi, error)
	UpdateFn  func(id string, identifikasi models.Identifikasi) error
	DeleteFn  func(id string) error
}

func (m *mockIdentifikasiRepository) Create(req dto.CreateIdentifikasiRequest, id string) error {
	return m.CreateFn(req, id)
}

func (m *mockIdentifikasiRepository) GetAll() ([]models.Identifikasi, error) {
	return m.GetAllFn()
}

func (m *mockIdentifikasiRepository) GetByID(id string) (*models.Identifikasi, error) {
	return m.GetByIDFn(id)
}

func (m *mockIdentifikasiRepository) Update(id string, identifikasi models.Identifikasi) error {
	return m.UpdateFn(id, identifikasi)
}

func (m *mockIdentifikasiRepository) Delete(id string) error {
	return m.DeleteFn(id)
}

// Compile-time check
var _ repository.IdentifikasiRepositoryInterface = (*mockIdentifikasiRepository)(nil)

//
// ===============================
// TEST CREATE
// ===============================
//

func TestIdentifikasiService_Create_Success(t *testing.T) {
	repo := &mockIdentifikasiRepository{
		CreateFn: func(req dto.CreateIdentifikasiRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*models.Identifikasi, error) {
			return &models.Identifikasi{
				ID:                id,
				NilaiIdentifikasi: 4.2,
				NilaiSubdomain1:   4,
				NilaiSubdomain2:   4.5,
				NilaiSubdomain3:   4.1,
				NilaiSubdomain4:   3.9,
				NilaiSubdomain5:   4.0,
			}, nil
		},
	}

	service := NewIdentifikasiService(repo)

	req := dto.CreateIdentifikasiRequest{
		NilaiIdentifikasi: 4.2,
		NilaiSubdomain1:   4,
		NilaiSubdomain2:   4.5,
		NilaiSubdomain3:   4.1,
		NilaiSubdomain4:   3.9,
		NilaiSubdomain5:   4.0,
	}

	result, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 4.2, result.NilaiIdentifikasi)
}

//
// ===============================
// TEST GET ALL
// ===============================
//

func TestIdentifikasiService_GetAll_Success(t *testing.T) {
	repo := &mockIdentifikasiRepository{
		GetAllFn: func() ([]models.Identifikasi, error) {
			return []models.Identifikasi{
				{ID: "1"},
				{ID: "2"},
			}, nil
		},
	}

	service := NewIdentifikasiService(repo)

	data, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, data, 2)
}

//
// ===============================
// TEST GET BY ID
// ===============================
//

func TestIdentifikasiService_GetByID_Success(t *testing.T) {
	repo := &mockIdentifikasiRepository{
		GetByIDFn: func(id string) (*models.Identifikasi, error) {
			return &models.Identifikasi{
				ID:                id,
				NilaiIdentifikasi: 3.8,
			}, nil
		},
	}

	service := NewIdentifikasiService(repo)

	result, err := service.GetByID("uuid-test")

	assert.NoError(t, err)
	assert.Equal(t, "uuid-test", result.ID)
}

func TestIdentifikasiService_GetByID_NotFound(t *testing.T) {
	repo := &mockIdentifikasiRepository{
		GetByIDFn: func(id string) (*models.Identifikasi, error) {
			return nil, errors.New("data tidak ditemukan")
		},
	}

	service := NewIdentifikasiService(repo)

	result, err := service.GetByID("invalid-id")

	assert.Error(t, err)
	assert.Nil(t, result)
}

//
// ===============================
// TEST UPDATE
// ===============================
//

func TestIdentifikasiService_Update_Success(t *testing.T) {
	newValue := 4.8

	repo := &mockIdentifikasiRepository{
		GetByIDFn: func(id string) (*models.Identifikasi, error) {
			return &models.Identifikasi{
				ID:                id,
				NilaiIdentifikasi: 3.5,
			}, nil
		},
		UpdateFn: func(id string, identifikasi models.Identifikasi) error {
			return nil
		},
	}

	service := NewIdentifikasiService(repo)

	req := dto.UpdateIdentifikasiRequest{
		NilaiIdentifikasi: &newValue,
	}

	result, err := service.Update("uuid-test", req)

	assert.NoError(t, err)
	assert.Equal(t, newValue, result.NilaiIdentifikasi)
}

//
// ===============================
// TEST DELETE
// ===============================
//

func TestIdentifikasiService_Delete_Success(t *testing.T) {
	repo := &mockIdentifikasiRepository{
		DeleteFn: func(id string) error {
			return nil
		},
	}

	service := NewIdentifikasiService(repo)

	err := service.Delete("uuid-test")

	assert.NoError(t, err)
}
