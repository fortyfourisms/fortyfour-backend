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
	GetByIkasIDFn        func(ikasID string) ([]models.Identifikasi, error)
	GetByPerusahaanIDFn  func(perusahaanID string) ([]models.Identifikasi, error)
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

func (m *mockIdentifikasiRepository) GetByPerusahaanID(perusahaanID string) ([]models.Identifikasi, error) {
	if m.GetByPerusahaanIDFn != nil {
		return m.GetByPerusahaanIDFn(perusahaanID)
	}
	return nil, nil
}

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
	ikasRepo := new(mockIkasRepository)

	service := NewIdentifikasiService(repo, ikasRepo)

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

func TestIdentifikasiService_GetByID_Success(t *testing.T) {
	repo := &mockIdentifikasiRepository{
		GetByIDFn: func(id string) (*models.Identifikasi, error) {
			return &models.Identifikasi{
				ID:                id,
				NilaiIdentifikasi: 3.8,
			}, nil
		},
	}

	ikasRepo := new(mockIkasRepository)
	service := NewIdentifikasiService(repo, ikasRepo)

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

	ikasRepo := &mockIkasRepository{}
	service := NewIdentifikasiService(repo, ikasRepo)

	result, err := service.GetByID("invalid-id", "admin", "")

	assert.Error(t, err)
	assert.Nil(t, result)
}
