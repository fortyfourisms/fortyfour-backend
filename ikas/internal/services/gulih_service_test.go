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
	GetByIkasIDFn func(ikasID string) ([]models.Gulih, error)
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

// compile-time safety check
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

	service := NewGulihService(repo)

	result, err := service.GetByID("invalid-id", "admin", "")

	assert.Error(t, err)
	assert.Nil(t, result)
}
