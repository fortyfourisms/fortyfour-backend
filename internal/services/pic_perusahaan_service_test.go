package services

import (
	"errors"
	"testing"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"

	"github.com/stretchr/testify/assert"
)

//
// ===============================
// MOCK PIC REPOSITORY
// ===============================
//

type mockPICRepository struct {
	CreateFn  func(req dto.CreatePICRequest, id string) error
	GetByIDFn func(id string) (*dto.PICResponse, error)
	GetAllFn  func() ([]dto.PICResponse, error)
	UpdateFn  func(id string, req dto.UpdatePICRequest) error
	DeleteFn  func(id string) error
}

func (m *mockPICRepository) Create(req dto.CreatePICRequest, id string) error {
	return m.CreateFn(req, id)
}

func (m *mockPICRepository) GetByID(id string) (*dto.PICResponse, error) {
	return m.GetByIDFn(id)
}

func (m *mockPICRepository) GetAll() ([]dto.PICResponse, error) {
	return m.GetAllFn()
}

func (m *mockPICRepository) Update(id string, req dto.UpdatePICRequest) error {
	return m.UpdateFn(id, req)
}

func (m *mockPICRepository) Delete(id string) error {
	return m.DeleteFn(id)
}

// Compile-time check
var _ repository.PICRepositoryInterface = (*mockPICRepository)(nil)

//
// ===============================
// TEST CREATE
// ===============================
//

func TestPICService_Create_Success(t *testing.T) {
	nama := "John Doe"
	idPerusahaan := "uuid-perusahaan"

	repo := &mockPICRepository{
		CreateFn: func(req dto.CreatePICRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*dto.PICResponse, error) {
			return &dto.PICResponse{
				ID:   id,
				Nama: nama,
				Perusahaan: &dto.PerusahaanInPIC{
					ID:             idPerusahaan,
					NamaPerusahaan: "PT Contoh",
				},
			}, nil
		},
	}

	service := NewPICService(repo, nil)

	req := dto.CreatePICRequest{
		Nama:         &nama,
		IDPerusahaan: &idPerusahaan,
	}

	result, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, nama, result.Nama)
	assert.NotNil(t, result.Perusahaan)
	assert.Equal(t, idPerusahaan, result.Perusahaan.ID)
}

func TestPICService_Create_ValidationError(t *testing.T) {
	repo := &mockPICRepository{}
	service := NewPICService(repo, nil)

	req := dto.CreatePICRequest{}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

//
// ===============================
// TEST GET ALL
// ===============================
//

func TestPICService_GetAll_Success(t *testing.T) {
	repo := &mockPICRepository{
		GetAllFn: func() ([]dto.PICResponse, error) {
			return []dto.PICResponse{
				{ID: "1", Nama: "PIC 1"},
				{ID: "2", Nama: "PIC 2"},
			}, nil
		},
	}

	service := NewPICService(repo, nil)

	result, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

//
// ===============================
// TEST GET BY ID
// ===============================
//

func TestPICService_GetByID_Success(t *testing.T) {
	repo := &mockPICRepository{
		GetByIDFn: func(id string) (*dto.PICResponse, error) {
			return &dto.PICResponse{
				ID:   id,
				Nama: "PIC Test",
				Perusahaan: &dto.PerusahaanInPIC{
					ID:             "uuid-perusahaan",
					NamaPerusahaan: "PT Test",
				},
			}, nil
		},
	}

	service := NewPICService(repo, nil)

	result, err := service.GetByID("uuid-test")

	assert.NoError(t, err)
	assert.Equal(t, "PIC Test", result.Nama)
	assert.NotNil(t, result.Perusahaan)
}

func TestPICService_GetByID_NotFound(t *testing.T) {
	repo := &mockPICRepository{
		GetByIDFn: func(id string) (*dto.PICResponse, error) {
			return nil, errors.New("data tidak ditemukan")
		},
	}

	service := NewPICService(repo, nil)

	result, err := service.GetByID("invalid-id")

	assert.Error(t, err)
	assert.Nil(t, result)
}

//
// ===============================
// TEST UPDATE
// ===============================
//

func TestPICService_Update_Success(t *testing.T) {
	namaBaru := "Nama Baru"

	repo := &mockPICRepository{
		UpdateFn: func(id string, req dto.UpdatePICRequest) error {
			return nil
		},
		GetByIDFn: func(id string) (*dto.PICResponse, error) {
			return &dto.PICResponse{
				ID:   id,
				Nama: namaBaru,
			}, nil
		},
	}

	service := NewPICService(repo, nil)

	req := dto.UpdatePICRequest{
		Nama: &namaBaru,
	}

	result, err := service.Update("uuid-test", req)

	assert.NoError(t, err)
	assert.Equal(t, namaBaru, result.Nama)
}

//
// ===============================
// TEST DELETE
// ===============================
//

func TestPICService_Delete_Success(t *testing.T) {
	repo := &mockPICRepository{
		DeleteFn: func(id string) error {
			return nil
		},
	}

	service := NewPICService(repo, nil)

	err := service.Delete("uuid-test")

	assert.NoError(t, err)
}
