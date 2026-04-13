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
// MOCK DETEKSI REPOSITORY
// ===============================
//

type mockDeteksiRepository struct {
	GetAllFn      func() ([]models.Deteksi, error)
	GetByIDFn     func(id string) (*models.Deteksi, error)
	GetByIkasIDFn func(ikasID string) ([]models.Deteksi, error)
}

func (m *mockDeteksiRepository) GetAll() ([]models.Deteksi, error) {
	return m.GetAllFn()
}

func (m *mockDeteksiRepository) GetByID(id string) (*models.Deteksi, error) {
	return m.GetByIDFn(id)
}

func (m *mockDeteksiRepository) GetByIkasID(ikasID string) ([]models.Deteksi, error) {
	if m.GetByIkasIDFn != nil {
		return m.GetByIkasIDFn(ikasID)
	}
	return nil, nil
}

// compile-time safety check
var _ repository.DeteksiRepositoryInterface = (*mockDeteksiRepository)(nil)

//
// ===============================
// TEST GET ALL
// ===============================
//

func TestDeteksiService_GetAll_Success(t *testing.T) {
	repo := &mockDeteksiRepository{
		GetAllFn: func() ([]models.Deteksi, error) {
			return []models.Deteksi{
				{ID: "1", NilaiDeteksi: 70},
				{ID: "2", NilaiDeteksi: 80},
			}, nil
		},
	}

	service := NewDeteksiService(repo)

	result, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

//
// ===============================
// TEST GET BY ID
// ===============================
//

func TestDeteksiService_GetByID_Success(t *testing.T) {
	repo := &mockDeteksiRepository{
		GetByIDFn: func(id string) (*models.Deteksi, error) {
			return &models.Deteksi{
				ID:           id,
				NilaiDeteksi: 90,
			}, nil
		},
	}

	service := NewDeteksiService(repo)

	result, err := service.GetByID("uuid-test", "admin", "")

	assert.NoError(t, err)
	assert.Equal(t, 90.0, result.NilaiDeteksi)
}

func TestDeteksiService_GetByID_NotFound(t *testing.T) {
	repo := &mockDeteksiRepository{
		GetByIDFn: func(id string) (*models.Deteksi, error) {
			return nil, errors.New("data tidak ditemukan")
		},
	}

	service := NewDeteksiService(repo)

	result, err := service.GetByID("invalid-id", "admin", "")

	assert.Error(t, err)
	assert.Nil(t, result)
}
