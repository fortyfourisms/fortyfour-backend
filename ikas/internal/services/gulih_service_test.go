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
// MOCK GULIH REPOSITORY
// ===============================
//

type mockGulihRepository struct {
	GetAllFn      func() ([]models.Gulih, error)
	GetByIDFn     func(id string) (*models.Gulih, error)
	GetByIkasIDFn        func(ikasID string) ([]models.Gulih, error)
	GetByPerusahaanIDFn  func(perusahaanID string) ([]models.Gulih, error)
	CloneByIkasIDFn      func(sourceID, targetID string) (string, error)
}

func (m *mockGulihRepository) GetAll() ([]models.Gulih, error) {
	return m.GetAllFn()
}

func (m *mockGulihRepository) GetByID(id string) (*models.Gulih, error) {
	return m.GetByIDFn(id)
}

func (m *mockGulihRepository) GetByIkasID(ikasID string) ([]models.Gulih, error) {
	if m.GetByIkasIDFn != nil {
		return m.GetByIkasIDFn(ikasID)
	}
	return nil, nil
}

func (m *mockGulihRepository) GetByPerusahaanID(perusahaanID string) ([]models.Gulih, error) {
	if m.GetByPerusahaanIDFn != nil {
		return m.GetByPerusahaanIDFn(perusahaanID)
	}
	return nil, nil
}

func (m *mockGulihRepository) CloneByIkasID(oldIkasID string, newIkasID string) (string, error) {
	if m.CloneByIkasIDFn != nil {
		return m.CloneByIkasIDFn(oldIkasID, newIkasID)
	}
	return "", nil
}

var _ repository.GulihRepositoryInterface = (*mockGulihRepository)(nil)

//
// ===============================
// TEST GET ALL
// ===============================
//

func TestGulihService_GetAll_Success(t *testing.T) {
	repo := &mockGulihRepository{
		GetAllFn: func() ([]models.Gulih, error) {
			return []models.Gulih{
				{ID: "1"},
				{ID: "2"},
			}, nil
		},
	}
	ikasRepo := new(mockIkasRepository)

	service := NewGulihService(repo, ikasRepo)

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

func TestGulihService_GetByID_Success(t *testing.T) {
	repo := &mockGulihRepository{
		GetByIDFn: func(id string) (*models.Gulih, error) {
			return &models.Gulih{
				ID:         id,
				NilaiGulih: 90,
			}, nil
		},
	}

	ikasRepo := new(mockIkasRepository)
	service := NewGulihService(repo, ikasRepo)

	result, err := service.GetByID("uuid-test", "admin", "")

	assert.NoError(t, err)
	assert.Equal(t, 90.0, result.NilaiGulih)
}

func TestGulihService_GetByID_NotFound(t *testing.T) {
	repo := &mockGulihRepository{
		GetByIDFn: func(id string) (*models.Gulih, error) {
			return nil, errors.New("data tidak ditemukan")
		},
	}

	ikasRepo := &mockIkasRepository{}
	service := NewGulihService(repo, ikasRepo)

	result, err := service.GetByID("invalid-id", "admin", "")

	assert.Error(t, err)
	assert.Nil(t, result)
}
