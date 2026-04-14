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
	GetByIkasIDFn        func(ikasID string) ([]models.Deteksi, error)
	GetByPerusahaanIDFn  func(perusahaanID string) ([]models.Deteksi, error)
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

func (m *mockDeteksiRepository) GetByPerusahaanID(perusahaanID string) ([]models.Deteksi, error) {
	if m.GetByPerusahaanIDFn != nil {
		return m.GetByPerusahaanIDFn(perusahaanID)
	}
	return nil, nil
}

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
				{ID: "1"},
				{ID: "2"},
			}, nil
		},
	}
	ikasRepo := new(mockIkasRepository)

	service := NewDeteksiService(repo, ikasRepo)

	// Admin can see all
	data, err := service.GetAll("admin")

	assert.NoError(t, err)
	assert.Len(t, data, 2)

	// Non-admin should fail
	_, err = service.GetAll("user")
	assert.Error(t, err)
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

	ikasRepo := new(mockIkasRepository)
	service := NewDeteksiService(repo, ikasRepo)

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

	ikasRepo := &mockIkasRepository{}
	service := NewDeteksiService(repo, ikasRepo)

	result, err := service.GetByID("invalid-id", "admin", "")

	assert.Error(t, err)
	assert.Nil(t, result)
}
