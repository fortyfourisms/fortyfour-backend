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
	GetAllFn      func() ([]models.Proteksi, error)
	GetByIDFn     func(id string) (*models.Proteksi, error)
	GetByIkasIDFn        func(ikasID string) ([]models.Proteksi, error)
	GetByPerusahaanIDFn  func(perusahaanID string) ([]models.Proteksi, error)
}

func (m *mockProteksiRepository) GetAll() ([]models.Proteksi, error) {
	return m.GetAllFn()
}

func (m *mockProteksiRepository) GetByID(id string) (*models.Proteksi, error) {
	return m.GetByIDFn(id)
}

func (m *mockProteksiRepository) GetByIkasID(ikasID string) ([]models.Proteksi, error) {
	if m.GetByIkasIDFn != nil {
		return m.GetByIkasIDFn(ikasID)
	}
	return nil, nil
}

func (m *mockProteksiRepository) GetByPerusahaanID(perusahaanID string) ([]models.Proteksi, error) {
	if m.GetByPerusahaanIDFn != nil {
		return m.GetByPerusahaanIDFn(perusahaanID)
	}
	return nil, nil
}

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
				{ID: "1"},
				{ID: "2"},
			}, nil
		},
	}
	ikasRepo := new(mockIkasRepository)

	service := NewProteksiService(repo, ikasRepo)

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

func TestProteksiService_GetByID_Success(t *testing.T) {
	repo := &mockProteksiRepository{
		GetByIDFn: func(id string) (*models.Proteksi, error) {
			return &models.Proteksi{
				ID:            id,
				NilaiProteksi: 90,
			}, nil
		},
	}

	ikasRepo := new(mockIkasRepository)
	service := NewProteksiService(repo, ikasRepo)

	result, err := service.GetByID("uuid-test", "admin", "")

	assert.NoError(t, err)
	assert.Equal(t, 90.0, result.NilaiProteksi)
}

func TestProteksiService_GetByID_NotFound(t *testing.T) {
	repo := &mockProteksiRepository{
		GetByIDFn: func(id string) (*models.Proteksi, error) {
			return nil, errors.New("data tidak ditemukan")
		},
	}

	ikasRepo := &mockIkasRepository{}
	service := NewProteksiService(repo, ikasRepo)

	result, err := service.GetByID("invalid-id", "admin", "")

	assert.Error(t, err)
	assert.Nil(t, result)
}
