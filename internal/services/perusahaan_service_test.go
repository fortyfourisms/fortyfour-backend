package services

import (
	"errors"
	"fortyfour-backend/internal/dto"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
=====================================
 MOCK PERUSAHAAN REPOSITORY
=====================================
*/

type mockPerusahaanRepository struct {
	CreateFn  func(req dto.CreatePerusahaanRequest, id string) error
	GetByIDFn func(id string) (*dto.PerusahaanResponse, error)
	GetAllFn  func() ([]dto.PerusahaanResponse, error)
	UpdateFn  func(id string, perusahaan dto.PerusahaanResponse) error
	DeleteFn  func(id string) error
}

func (m *mockPerusahaanRepository) Create(req dto.CreatePerusahaanRequest, id string) error {
	return m.CreateFn(req, id)
}

func (m *mockPerusahaanRepository) GetByID(id string) (*dto.PerusahaanResponse, error) {
	return m.GetByIDFn(id)
}

func (m *mockPerusahaanRepository) GetAll() ([]dto.PerusahaanResponse, error) {
	return m.GetAllFn()
}

func (m *mockPerusahaanRepository) Update(id string, perusahaan dto.PerusahaanResponse) error {
	return m.UpdateFn(id, perusahaan)
}

func (m *mockPerusahaanRepository) Delete(id string) error {
	return m.DeleteFn(id)
}

/*
=====================================
 MOCK SUB SEKTOR REPOSITORY
=====================================
*/

type mockSubSektorRepository struct {
	GetByIDFn       func(id string) (*dto.SubSektorResponse, error)
	GetAllFn        func() ([]dto.SubSektorResponse, error)
	GetBySektorIDFn func(sektorID string) ([]dto.SubSektorResponse, error)
}

func (m *mockSubSektorRepository) GetByID(id string) (*dto.SubSektorResponse, error) {
	return m.GetByIDFn(id)
}

func (m *mockSubSektorRepository) GetAll() ([]dto.SubSektorResponse, error) {
	return m.GetAllFn()
}

func (m *mockSubSektorRepository) GetBySektorID(sektorID string) ([]dto.SubSektorResponse, error) {
	return m.GetBySektorIDFn(sektorID)
}

/*
=====================================
 TEST CREATE PERUSAHAAN
=====================================
*/

func TestCreatePerusahaan_Success(t *testing.T) {
	nama := "PT Teknologi Maju"
	idSubSektor := "sub-sektor-id-123"

	perusahaanRepo := &mockPerusahaanRepository{
		CreateFn: func(req dto.CreatePerusahaanRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{
				ID:             id,
				NamaPerusahaan: nama,
				SubSektor: &dto.SubSektorResponse{
					ID:            idSubSektor,
					NamaSubSektor: "Elektronik",
					IDSektor:      "sektor-id-123",
					NamaSektor:    "ILMATE",
				},
			}, nil
		},
	}

	subSektorRepo := &mockSubSektorRepository{
		GetByIDFn: func(id string) (*dto.SubSektorResponse, error) {
			return &dto.SubSektorResponse{
				ID:            idSubSektor,
				NamaSubSektor: "Elektronik",
				IDSektor:      "sektor-id-123",
				NamaSektor:    "ILMATE",
			}, nil
		},
	}

	service := NewPerusahaanService(perusahaanRepo, subSektorRepo)

	req := dto.CreatePerusahaanRequest{
		NamaPerusahaan: &nama,
		IDSubSektor:    &idSubSektor,
	}

	result, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, nama, result.NamaPerusahaan)
	assert.NotNil(t, result.SubSektor)
	assert.Equal(t, "Elektronik", result.SubSektor.NamaSubSektor)
}

func TestCreatePerusahaan_ValidationFailed_NamaKosong(t *testing.T) {
	perusahaanRepo := &mockPerusahaanRepository{}
	subSektorRepo := &mockSubSektorRepository{}
	service := NewPerusahaanService(perusahaanRepo, subSektorRepo)

	idSubSektor := "sub-sektor-id-123"
	req := dto.CreatePerusahaanRequest{
		NamaPerusahaan: nil,
		IDSubSektor:    &idSubSektor,
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "nama_perusahaan wajib diisi", err.Error())
}

func TestCreatePerusahaan_ValidationFailed_SubSektorKosong(t *testing.T) {
	nama := "PT ABC"

	perusahaanRepo := &mockPerusahaanRepository{
		CreateFn: func(req dto.CreatePerusahaanRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{
				ID:             id,
				NamaPerusahaan: nama,
				SubSektor:      nil,
			}, nil
		},
	}

	subSektorRepo := &mockSubSektorRepository{}
	service := NewPerusahaanService(perusahaanRepo, subSektorRepo)

	req := dto.CreatePerusahaanRequest{
		NamaPerusahaan: &nama,
		IDSubSektor:    nil,
	}

	result, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, nama, result.NamaPerusahaan)
	assert.Nil(t, result.SubSektor)
}

func TestCreatePerusahaan_SubSektorNotFound(t *testing.T) {
	perusahaanRepo := &mockPerusahaanRepository{}

	subSektorRepo := &mockSubSektorRepository{
		GetByIDFn: func(id string) (*dto.SubSektorResponse, error) {
			return nil, errors.New("not found")
		},
	}

	service := NewPerusahaanService(perusahaanRepo, subSektorRepo)

	nama := "PT ABC"
	idSubSektor := "invalid-sub-sektor-id"

	req := dto.CreatePerusahaanRequest{
		NamaPerusahaan: &nama,
		IDSubSektor:    &idSubSektor,
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "sub sektor tidak ditemukan", err.Error())
}

func TestCreatePerusahaan_RepoFailed(t *testing.T) {
	nama := "PT ABC"
	idSubSektor := "sub-sektor-id-123"

	perusahaanRepo := &mockPerusahaanRepository{
		CreateFn: func(req dto.CreatePerusahaanRequest, id string) error {
			return errors.New("db error")
		},
	}

	subSektorRepo := &mockSubSektorRepository{
		GetByIDFn: func(id string) (*dto.SubSektorResponse, error) {
			return &dto.SubSektorResponse{
				ID:            idSubSektor,
				NamaSubSektor: "Elektronik",
			}, nil
		},
	}

	service := NewPerusahaanService(perusahaanRepo, subSektorRepo)

	req := dto.CreatePerusahaanRequest{
		NamaPerusahaan: &nama,
		IDSubSektor:    &idSubSektor,
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

/*
=====================================
 TEST GET
=====================================
*/

func TestGetAllPerusahaan_Success(t *testing.T) {
	perusahaanRepo := &mockPerusahaanRepository{
		GetAllFn: func() ([]dto.PerusahaanResponse, error) {
			return []dto.PerusahaanResponse{
				{
					ID:             "1",
					NamaPerusahaan: "PT A",
					SubSektor: &dto.SubSektorResponse{
						ID:            "sub-1",
						NamaSubSektor: "Elektronik",
						NamaSektor:    "ILMATE",
					},
				},
				{
					ID:             "2",
					NamaPerusahaan: "PT B",
					SubSektor: &dto.SubSektorResponse{
						ID:            "sub-2",
						NamaSubSektor: "Tekstil",
						NamaSektor:    "IKFT",
					},
				},
			}, nil
		},
	}

	subSektorRepo := &mockSubSektorRepository{}
	service := NewPerusahaanService(perusahaanRepo, subSektorRepo)

	data, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, data, 2)
	assert.Equal(t, "Elektronik", data[0].SubSektor.NamaSubSektor)
	assert.Equal(t, "Tekstil", data[1].SubSektor.NamaSubSektor)
}

func TestGetPerusahaanByID_Success(t *testing.T) {
	perusahaanRepo := &mockPerusahaanRepository{
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{
				ID:             id,
				NamaPerusahaan: "PT A",
				SubSektor: &dto.SubSektorResponse{
					ID:            "sub-1",
					NamaSubSektor: "Elektronik",
					IDSektor:      "sektor-1",
					NamaSektor:    "ILMATE",
				},
			}, nil
		},
	}

	subSektorRepo := &mockSubSektorRepository{}
	service := NewPerusahaanService(perusahaanRepo, subSektorRepo)

	result, err := service.GetByID("123")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "PT A", result.NamaPerusahaan)
	assert.Equal(t, "Elektronik", result.SubSektor.NamaSubSektor)
}

/*
=====================================
 TEST UPDATE PERUSAHAAN
=====================================
*/

func TestUpdatePerusahaan_Success(t *testing.T) {
	oldNama := "PT Lama"
	newNama := "PT Baru"

	perusahaanRepo := &mockPerusahaanRepository{
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{
				ID:             id,
				NamaPerusahaan: oldNama,
				SubSektor: &dto.SubSektorResponse{
					ID:            "sub-1",
					NamaSubSektor: "Elektronik",
					IDSektor:      "sektor-1",
					NamaSektor:    "ILMATE",
				},
			}, nil
		},
		UpdateFn: func(id string, perusahaan dto.PerusahaanResponse) error {
			return nil
		},
	}

	subSektorRepo := &mockSubSektorRepository{}
	service := NewPerusahaanService(perusahaanRepo, subSektorRepo)

	req := dto.UpdatePerusahaanRequest{
		NamaPerusahaan: &newNama,
	}

	result, err := service.Update("123", req)

	assert.NoError(t, err)
	assert.Equal(t, newNama, result.NamaPerusahaan)
	assert.NotNil(t, result.SubSektor)
	assert.Equal(t, "Elektronik", result.SubSektor.NamaSubSektor)
}

func TestUpdatePerusahaan_UpdateSubSektor_Success(t *testing.T) {
	oldSubSektorID := "sub-1"
	newSubSektorID := "sub-2"

	perusahaanRepo := &mockPerusahaanRepository{
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{
				ID:             id,
				NamaPerusahaan: "PT A",
				SubSektor: &dto.SubSektorResponse{
					ID:            oldSubSektorID,
					NamaSubSektor: "Elektronik",
					IDSektor:      "sektor-1",
					NamaSektor:    "ILMATE",
				},
			}, nil
		},
		UpdateFn: func(id string, perusahaan dto.PerusahaanResponse) error {
			return nil
		},
	}

	subSektorRepo := &mockSubSektorRepository{
		GetByIDFn: func(id string) (*dto.SubSektorResponse, error) {
			return &dto.SubSektorResponse{
				ID:            newSubSektorID,
				NamaSubSektor: "Otomotif",
				IDSektor:      "sektor-1",
				NamaSektor:    "ILMATE",
			}, nil
		},
	}

	service := NewPerusahaanService(perusahaanRepo, subSektorRepo)

	req := dto.UpdatePerusahaanRequest{
		IDSubSektor: &newSubSektorID,
	}

	result, err := service.Update("123", req)

	assert.NoError(t, err)
	assert.NotNil(t, result.SubSektor)
	assert.Equal(t, newSubSektorID, result.SubSektor.ID)
	assert.Equal(t, "Otomotif", result.SubSektor.NamaSubSektor)
}

func TestUpdatePerusahaan_InvalidSubSektor(t *testing.T) {
	invalidSubSektorID := "invalid-sub-sektor"

	perusahaanRepo := &mockPerusahaanRepository{
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{
				ID:             id,
				NamaPerusahaan: "PT A",
				SubSektor: &dto.SubSektorResponse{
					ID:            "sub-1",
					NamaSubSektor: "Elektronik",
				},
			}, nil
		},
	}

	subSektorRepo := &mockSubSektorRepository{
		GetByIDFn: func(id string) (*dto.SubSektorResponse, error) {
			return nil, errors.New("not found")
		},
	}

	service := NewPerusahaanService(perusahaanRepo, subSektorRepo)

	req := dto.UpdatePerusahaanRequest{
		IDSubSektor: &invalidSubSektorID,
	}

	result, err := service.Update("123", req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "sub sektor tidak ditemukan", err.Error())
}

/*
=====================================
 TEST DELETE PERUSAHAAN
=====================================
*/

func TestDeletePerusahaan_Success(t *testing.T) {
	perusahaanRepo := &mockPerusahaanRepository{
		DeleteFn: func(id string) error {
			return nil
		},
	}

	subSektorRepo := &mockSubSektorRepository{}
	service := NewPerusahaanService(perusahaanRepo, subSektorRepo)

	err := service.Delete("123")

	assert.NoError(t, err)
}

func TestDeletePerusahaan_Failed(t *testing.T) {
	perusahaanRepo := &mockPerusahaanRepository{
		DeleteFn: func(id string) error {
			return errors.New("db error")
		},
	}

	subSektorRepo := &mockSubSektorRepository{}
	service := NewPerusahaanService(perusahaanRepo, subSektorRepo)

	err := service.Delete("123")

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}
