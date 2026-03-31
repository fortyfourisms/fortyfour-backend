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
// MOCK PROTEKSI REPOSITORY
// ===============================
//

type mockProteksiRepository struct {
	GetAllFn  func() ([]models.Proteksi, error)
	GetByIDFn func(id string) (*models.Proteksi, error)
}

func (m *mockProteksiRepository) GetAll() ([]models.Proteksi, error) {
	return m.GetAllFn()
}

func (m *mockProteksiRepository) GetByID(id string) (*models.Proteksi, error) {
	return m.GetByIDFn(id)
}

// compile-time safety check
var _ repository.ProteksiRepositoryInterface = (*mockProteksiRepository)(nil)

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
