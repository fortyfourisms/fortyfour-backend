package services

import (
	"fortyfour-backend/internal/dto"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockSdmCsirtRepo struct {
	CreateFn  func(req dto.CreateSdmCsirtRequest, id string) error
	GetAllFn  func() ([]dto.SdmCsirtResponse, error)
	GetByIDFn func(id string) (*dto.SdmCsirtResponse, error)
	UpdateFn  func(id string, req dto.SdmCsirtResponse) error
	DeleteFn  func(id string) error
}

func (m *mockSdmCsirtRepo) Create(req dto.CreateSdmCsirtRequest, id string) error {
	return m.CreateFn(req, id)
}
func (m *mockSdmCsirtRepo) GetAll() ([]dto.SdmCsirtResponse, error) {
	return m.GetAllFn()
}
func (m *mockSdmCsirtRepo) GetByID(id string) (*dto.SdmCsirtResponse, error) {
	return m.GetByIDFn(id)
}
func (m *mockSdmCsirtRepo) Update(id string, req dto.SdmCsirtResponse) error {
	return m.UpdateFn(id, req)
}
func (m *mockSdmCsirtRepo) Delete(id string) error {
	return m.DeleteFn(id)
}

func strPtr(s string) *string {
	return &s
}

func TestSdmCsirtService_Create(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		CreateFn: func(req dto.CreateSdmCsirtRequest, id string) error {
			return nil
		},
	}

	service := NewSdmCsirtService(repo)

	id, err := service.Create(dto.CreateSdmCsirtRequest{
		NamaPersonel: strPtr("Andi"),
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestSdmCsirtService_GetAll(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		GetAllFn: func() ([]dto.SdmCsirtResponse, error) {
			return []dto.SdmCsirtResponse{
				{NamaPersonel: "Budi"},
			}, nil
		},
	}

	service := NewSdmCsirtService(repo)

	res, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestSdmCsirtService_Update(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return &dto.SdmCsirtResponse{
				NamaPersonel: "Old",
			}, nil
		},
		UpdateFn: func(id string, req dto.SdmCsirtResponse) error {
			return nil
		},
	}

	service := NewSdmCsirtService(repo)

	newName := "New"
	err := service.Update("1", dto.UpdateSdmCsirtRequest{
		NamaPersonel: &newName,
	})

	assert.NoError(t, err)
}

func TestSdmCsirtService_Delete(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		DeleteFn: func(id string) error {
			return nil
		},
	}

	service := NewSdmCsirtService(repo)

	err := service.Delete("1")

	assert.NoError(t, err)
}
