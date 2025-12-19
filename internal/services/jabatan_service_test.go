package services

import (
	"errors"
	"testing"

	"fortyfour-backend/internal/dto"

	"github.com/stretchr/testify/assert"
)

/*
=====================================
 MOCK JABATAN REPOSITORY
=====================================
*/

type mockJabatanRepository struct {
	CreateFn  func(req dto.CreateJabatanRequest, id string) error
	GetByIDFn func(id string) (*dto.JabatanResponse, error)
	GetAllFn  func() ([]dto.JabatanResponse, error)
	UpdateFn  func(id string, jabatan dto.JabatanResponse) error
	DeleteFn  func(id string) error
}

func (m *mockJabatanRepository) Create(req dto.CreateJabatanRequest, id string) error {
	return m.CreateFn(req, id)
}

func (m *mockJabatanRepository) GetByID(id string) (*dto.JabatanResponse, error) {
	return m.GetByIDFn(id)
}

func (m *mockJabatanRepository) GetAll() ([]dto.JabatanResponse, error) {
	return m.GetAllFn()
}

func (m *mockJabatanRepository) Update(id string, jabatan dto.JabatanResponse) error {
	return m.UpdateFn(id, jabatan)
}

func (m *mockJabatanRepository) Delete(id string) error {
	return m.DeleteFn(id)
}

/*
=====================================
 TEST CREATE
=====================================
*/

func TestCreateJabatan_Success(t *testing.T) {
	nama := "Manager IT"

	repo := &mockJabatanRepository{
		CreateFn: func(req dto.CreateJabatanRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*dto.JabatanResponse, error) {
			return &dto.JabatanResponse{
				ID:          id,
				NamaJabatan: nama,
			}, nil
		},
	}

	service := NewJabatanService(repo)

	req := dto.CreateJabatanRequest{
		NamaJabatan: &nama,
	}

	result, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, nama, result.NamaJabatan)
}

func TestCreateJabatan_ValidationFailed(t *testing.T) {
	repo := &mockJabatanRepository{}

	service := NewJabatanService(repo)

	req := dto.CreateJabatanRequest{
		NamaJabatan: nil,
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCreateJabatan_RepoFailed(t *testing.T) {
	nama := "Manager"

	repo := &mockJabatanRepository{
		CreateFn: func(req dto.CreateJabatanRequest, id string) error {
			return errors.New("db error")
		},
	}

	service := NewJabatanService(repo)

	req := dto.CreateJabatanRequest{
		NamaJabatan: &nama,
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

/*
=====================================
 TEST GET ALL
=====================================
*/

func TestGetAllJabatan_Success(t *testing.T) {
	repo := &mockJabatanRepository{
		GetAllFn: func() ([]dto.JabatanResponse, error) {
			return []dto.JabatanResponse{
				{ID: "1", NamaJabatan: "Manager"},
				{ID: "2", NamaJabatan: "Staff"},
			}, nil
		},
	}

	service := NewJabatanService(repo)

	data, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, data, 2)
}

/*
=====================================
 TEST UPDATE
=====================================
*/

func TestUpdateJabatan_Success(t *testing.T) {
	newName := "Updated Jabatan"

	repo := &mockJabatanRepository{
		GetByIDFn: func(id string) (*dto.JabatanResponse, error) {
			return &dto.JabatanResponse{
				ID:          id,
				NamaJabatan: "Old Jabatan",
			}, nil
		},
		UpdateFn: func(id string, jabatan dto.JabatanResponse) error {
			return nil
		},
	}

	service := NewJabatanService(repo)

	req := dto.UpdateJabatanRequest{
		NamaJabatan: &newName,
	}

	result, err := service.Update("123", req)

	assert.NoError(t, err)
	assert.Equal(t, newName, result.NamaJabatan)
}

/*
=====================================
 TEST DELETE
=====================================
*/

func TestDeleteJabatan_Success(t *testing.T) {
	repo := &mockJabatanRepository{
		DeleteFn: func(id string) error {
			return nil
		},
	}

	service := NewJabatanService(repo)

	err := service.Delete("123")

	assert.NoError(t, err)
}
