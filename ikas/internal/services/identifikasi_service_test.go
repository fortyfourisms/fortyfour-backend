package services

import (
	"errors"
	"ikas/internal/models"
	"ikas/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

//
// ===============================
// MOCK IDENTIFIKASI REPOSITORY
// ===============================
//

type mockIdentifikasiRepository struct {
	GetAllFn      func() ([]models.Identifikasi, error)
	GetByIDFn     func(id string) (*models.Identifikasi, error)
	GetByIkasIDFn func(ikasID string) ([]models.Identifikasi, error)
}

func (m *mockIdentifikasiRepository) GetAll() ([]models.Identifikasi, error) {
	return m.GetAllFn()
}

func (m *mockIdentifikasiRepository) GetByID(id string) (*models.Identifikasi, error) {
	return m.GetByIDFn(id)
}

func (m *mockIdentifikasiRepository) GetByIkasID(ikasID string) ([]models.Identifikasi, error) {
	if m.GetByIkasIDFn != nil {
		return m.GetByIkasIDFn(ikasID)
	}
	return nil, nil
}

// Compile-time check
var _ repository.IdentifikasiRepositoryInterface = (*mockIdentifikasiRepository)(nil)

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

	result, err := service.GetByID("uuid-test", "admin", "")

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

	result, err := service.GetByID("invalid-id", "admin", "")

	assert.Error(t, err)
	assert.Nil(t, result)
}
