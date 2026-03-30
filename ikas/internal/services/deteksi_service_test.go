package services

import (
	"errors"
	"ikas/internal/dto"
	"ikas/internal/models"
	"ikas/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

//
// ===============================
// MOCK DETEKSI REPOSITORY
// ===============================
//

type mockDeteksiRepository struct {
	CreateFn  func(req dto.CreateDeteksiRequest, id string) error
	GetAllFn  func() ([]models.Deteksi, error)
	GetByIDFn func(id string) (*models.Deteksi, error)
	UpdateFn  func(id string, deteksi models.Deteksi) error
	DeleteFn  func(id string) error
}

func (m *mockDeteksiRepository) Create(req dto.CreateDeteksiRequest, id string) error {
	return m.CreateFn(req, id)
}

func (m *mockDeteksiRepository) GetAll() ([]models.Deteksi, error) {
	return m.GetAllFn()
}

func (m *mockDeteksiRepository) GetByID(id string) (*models.Deteksi, error) {
	return m.GetByIDFn(id)
}

func (m *mockDeteksiRepository) Update(id string, deteksi models.Deteksi) error {
	return m.UpdateFn(id, deteksi)
}

func (m *mockDeteksiRepository) Delete(id string) error {
	return m.DeleteFn(id)
}

// compile-time safety check
var _ repository.DeteksiRepositoryInterface = (*mockDeteksiRepository)(nil)

//
// ===============================
// TEST CREATE
// ===============================
//

func TestDeteksiService_Create_Success(t *testing.T) {
	repo := &mockDeteksiRepository{
		CreateFn: func(req dto.CreateDeteksiRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*models.Deteksi, error) {
			return &models.Deteksi{
				ID:           id,
				NilaiDeteksi: 85,
			}, nil
		},
	}

	service := NewDeteksiService(repo)

	req := dto.CreateDeteksiRequest{
		NilaiDeteksi: 85,
	}

	result, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 85.0, result.NilaiDeteksi)
}

func TestDeteksiService_Create_Error(t *testing.T) {
	repo := &mockDeteksiRepository{
		CreateFn: func(req dto.CreateDeteksiRequest, id string) error {
			return errors.New("gagal create")
		},
	}

	service := NewDeteksiService(repo)

	result, err := service.Create(dto.CreateDeteksiRequest{})

	assert.Error(t, err)
	assert.Nil(t, result)
}

//
// ===============================
// TEST GET ALL
// ===============================
//

func TestDeteksiService_GetAll_Success(t *testing.T) {
	repo := &mockDeteksiRepository{
		GetAllFn: func() ([]models.Deteksi, error) {
			return []models.Deteksi{
				{ID: "1", NilaiDeteksi: 70},
				{ID: "2", NilaiDeteksi: 80},
			}, nil
		},
	}

	service := NewDeteksiService(repo)

	result, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

//
// ===============================
// TEST GET BY ID
// ===============================
//

func TestDeteksiService_GetByID_Success(t *testing.T) {
	repo := &mockDeteksiRepository{
		GetByIDFn: func(id string) (*models.Deteksi, error) {
			return &models.Deteksi{
				ID:           id,
				NilaiDeteksi: 90,
			}, nil
		},
	}

	service := NewDeteksiService(repo)

	result, err := service.GetByID("uuid-test")

	assert.NoError(t, err)
	assert.Equal(t, 90.0, result.NilaiDeteksi)
}

func TestDeteksiService_GetByID_NotFound(t *testing.T) {
	repo := &mockDeteksiRepository{
		GetByIDFn: func(id string) (*models.Deteksi, error) {
			return nil, errors.New("data tidak ditemukan")
		},
	}

	service := NewDeteksiService(repo)

	result, err := service.GetByID("invalid-id")

	assert.Error(t, err)
	assert.Nil(t, result)
}

//
// ===============================
// TEST UPDATE
// ===============================
//

func TestDeteksiService_Update_Success(t *testing.T) {
	nilaiBaru := 95.0

	repo := &mockDeteksiRepository{
		GetByIDFn: func(id string) (*models.Deteksi, error) {
			return &models.Deteksi{
				ID:           id,
				NilaiDeteksi: 80,
			}, nil
		},
		UpdateFn: func(id string, deteksi models.Deteksi) error {
			return nil
		},
	}

	service := NewDeteksiService(repo)

	req := dto.UpdateDeteksiRequest{
		NilaiDeteksi: &nilaiBaru,
	}

	result, err := service.Update("uuid-test", req)

	assert.NoError(t, err)
	assert.Equal(t, nilaiBaru, result.NilaiDeteksi)
}

//
// ===============================
// TEST DELETE
// ===============================
//

func TestDeteksiService_Delete_Success(t *testing.T) {
	repo := &mockDeteksiRepository{
		DeleteFn: func(id string) error {
			return nil
		},
	}

	service := NewDeteksiService(repo)

	err := service.Delete("uuid-test")

	assert.NoError(t, err)
}
