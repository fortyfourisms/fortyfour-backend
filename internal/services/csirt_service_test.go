package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockCsirtRepo struct {
	CreateFn                func(req dto.CreateCsirtRequest, id string) error
	GetByIDFn               func(id string) (*models.Csirt, error)
	GetAllWithPerusahaanFn  func() ([]dto.CsirtResponse, error)
	GetByIDWithPerusahaanFn func(id string) (*dto.CsirtResponse, error)
	UpdateFn                func(id string, csirt models.Csirt) error
	DeleteFn                func(id string) error
}

func (m *mockCsirtRepo) Create(req dto.CreateCsirtRequest, id string) error {
	return m.CreateFn(req, id)
}
func (m *mockCsirtRepo) GetByID(id string) (*models.Csirt, error) {
	return m.GetByIDFn(id)
}
func (m *mockCsirtRepo) GetAllWithPerusahaan() ([]dto.CsirtResponse, error) {
	return m.GetAllWithPerusahaanFn()
}
func (m *mockCsirtRepo) GetByIDWithPerusahaan(id string) (*dto.CsirtResponse, error) {
	return m.GetByIDWithPerusahaanFn(id)
}
func (m *mockCsirtRepo) Update(id string, csirt models.Csirt) error {
	return m.UpdateFn(id, csirt)
}
func (m *mockCsirtRepo) Delete(id string) error {
	return m.DeleteFn(id)
}

func TestCsirtService_Create_Success(t *testing.T) {
	repo := &mockCsirtRepo{
		CreateFn: func(req dto.CreateCsirtRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return &models.Csirt{ID: id, NamaCsirt: "CSIRT A"}, nil
		},
	}

	service := NewCsirtService(repo)

	res, err := service.Create(dto.CreateCsirtRequest{
		NamaCsirt: "CSIRT A",
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "CSIRT A", res.NamaCsirt)
}

func TestCsirtService_GetAll(t *testing.T) {
	repo := &mockCsirtRepo{
		GetAllWithPerusahaanFn: func() ([]dto.CsirtResponse, error) {
			return []dto.CsirtResponse{
				{NamaCsirt: "CSIRT A"},
			}, nil
		},
	}

	service := NewCsirtService(repo)

	res, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestCsirtService_Update(t *testing.T) {
	repo := &mockCsirtRepo{
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return &models.Csirt{ID: id, NamaCsirt: "Old"}, nil
		},
		UpdateFn: func(id string, csirt models.Csirt) error {
			return nil
		},
	}

	service := NewCsirtService(repo)

	newName := "New CSIRT"
	res, err := service.Update("1", dto.UpdateCsirtRequest{
		NamaCsirt: &newName,
	})

	assert.NoError(t, err)
	assert.Equal(t, "New CSIRT", res.NamaCsirt)
}

func TestCsirtService_Delete(t *testing.T) {
	repo := &mockCsirtRepo{
		DeleteFn: func(id string) error {
			return nil
		},
	}

	service := NewCsirtService(repo)

	err := service.Delete("1")

	assert.NoError(t, err)
}
