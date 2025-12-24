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
 TEST CREATE PERUSAHAAN
=====================================
*/

func TestCreatePerusahaan_Success(t *testing.T) {
	nama := "PT Teknologi Maju"
	sektor := "Teknologi"

	repo := &mockPerusahaanRepository{
		CreateFn: func(req dto.CreatePerusahaanRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{
				ID:             id,
				NamaPerusahaan: nama,
				Sektor:         sektor,
			}, nil
		},
	}

	service := NewPerusahaanService(repo)

	req := dto.CreatePerusahaanRequest{
		NamaPerusahaan: &nama,
		Sektor:         &sektor,
	}

	result, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, nama, result.NamaPerusahaan)
	assert.Equal(t, sektor, result.Sektor)
}

func TestCreatePerusahaan_ValidationFailed_NamaKosong(t *testing.T) {
	repo := &mockPerusahaanRepository{}
	service := NewPerusahaanService(repo)

	sektor := "Teknologi"
	req := dto.CreatePerusahaanRequest{
		NamaPerusahaan: nil,
		Sektor:         &sektor,
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "nama_perusahaan wajib diisi", err.Error())
}

func TestCreatePerusahaan_ValidationFailed_SektorKosong(t *testing.T) {
	repo := &mockPerusahaanRepository{}
	service := NewPerusahaanService(repo)

	nama := "PT ABC"
	req := dto.CreatePerusahaanRequest{
		NamaPerusahaan: &nama,
		Sektor:         nil,
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "sektor wajib diisi", err.Error())
}

func TestCreatePerusahaan_InvalidSektor(t *testing.T) {
	repo := &mockPerusahaanRepository{}
	service := NewPerusahaanService(repo)

	nama := "PT ABC"
	sektor := "Pertanian"

	req := dto.CreatePerusahaanRequest{
		NamaPerusahaan: &nama,
		Sektor:         &sektor,
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "sektor tidak valid", err.Error())
}

func TestCreatePerusahaan_RepoFailed(t *testing.T) {
	nama := "PT ABC"
	sektor := "Teknologi"

	repo := &mockPerusahaanRepository{
		CreateFn: func(req dto.CreatePerusahaanRequest, id string) error {
			return errors.New("db error")
		},
	}

	service := NewPerusahaanService(repo)

	req := dto.CreatePerusahaanRequest{
		NamaPerusahaan: &nama,
		Sektor:         &sektor,
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
	repo := &mockPerusahaanRepository{
		GetAllFn: func() ([]dto.PerusahaanResponse, error) {
			return []dto.PerusahaanResponse{
				{ID: "1", NamaPerusahaan: "PT A", Sektor: "Teknologi"},
				{ID: "2", NamaPerusahaan: "PT B", Sektor: "Keuangan"},
			}, nil
		},
	}

	service := NewPerusahaanService(repo)

	data, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, data, 2)
}

func TestGetPerusahaanByID_Success(t *testing.T) {
	repo := &mockPerusahaanRepository{
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{
				ID:             id,
				NamaPerusahaan: "PT A",
				Sektor:         "Teknologi",
			}, nil
		},
	}

	service := NewPerusahaanService(repo)

	result, err := service.GetByID("123")

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

/*
=====================================
 TEST UPDATE PERUSAHAAN
=====================================
*/

func TestUpdatePerusahaan_Success(t *testing.T) {
	oldNama := "PT Lama"
	newNama := "PT Baru"
	sektor := "Teknologi"

	repo := &mockPerusahaanRepository{
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{
				ID:             id,
				NamaPerusahaan: oldNama,
				Sektor:         sektor,
			}, nil
		},
		UpdateFn: func(id string, perusahaan dto.PerusahaanResponse) error {
			return nil
		},
	}

	service := NewPerusahaanService(repo)

	req := dto.UpdatePerusahaanRequest{
		NamaPerusahaan: &newNama,
	}

	result, err := service.Update("123", req)

	assert.NoError(t, err)
	assert.Equal(t, newNama, result.NamaPerusahaan)
	assert.Equal(t, sektor, result.Sektor)
}

func TestUpdatePerusahaan_InvalidSektor(t *testing.T) {
	sektorInvalid := "Pertanian"

	repo := &mockPerusahaanRepository{
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{
				ID:             id,
				NamaPerusahaan: "PT A",
				Sektor:         "Teknologi",
			}, nil
		},
	}

	service := NewPerusahaanService(repo)

	req := dto.UpdatePerusahaanRequest{
		Sektor: &sektorInvalid,
	}

	result, err := service.Update("123", req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "sektor tidak valid", err.Error())
}

/*
=====================================
 TEST DELETE PERUSAHAAN
=====================================
*/

func TestDeletePerusahaan_Success(t *testing.T) {
	repo := &mockPerusahaanRepository{
		DeleteFn: func(id string) error {
			return nil
		},
	}

	service := NewPerusahaanService(repo)

	err := service.Delete("123")

	assert.NoError(t, err)
}
