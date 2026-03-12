package services

import (
	"ikas/internal/dto"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
========================================
 MOCK REPOSITORY
========================================
*/

type mockIkasRepo struct{}

func (m *mockIkasRepo) FindPerusahaanByName(namaPerusahaan string) (string, error) {
	return "perusahaan-id-1", nil
}

// CREATE IKAS
func (m *mockIkasRepo) Create(
	req dto.CreateIkasRequest,
	id string,
	nilai float64,
) error {
	return nil
}

func (m *mockIkasRepo) GetAll() ([]dto.IkasResponse, error) {
	return []dto.IkasResponse{}, nil
}

func (m *mockIkasRepo) GetByID(id string) (*dto.IkasResponse, error) {
	return &dto.IkasResponse{
		ID:              id,
		NilaiKematangan: 75,
	}, nil
}

func (m *mockIkasRepo) Update(id string, req dto.UpdateIkasRequest) error {
	return nil
}

func (m *mockIkasRepo) Delete(id string) error {
	return nil
}

func (m *mockIkasRepo) ParseExcelForImport(b []byte) (*dto.CreateIkasRequest, error) {
	return &dto.CreateIkasRequest{}, nil
}

// Mock methods for Excel import process in repo
func (m *mockIkasRepo) SaveImportedData(id string, data *dto.CreateIkasRequest) error {
	return nil
}

/*
========================================
 TEST CREATE
========================================
*/

func TestIkasService_Create_Success(t *testing.T) {
	repo := &mockIkasRepo{}
	// Nil producer for test
	service := NewIkasService(repo, nil)

	req := dto.CreateIkasRequest{
		IDPerusahaan: "1",
	}

	// This should return nil now because producer check skips it
	err := service.Create(req, "ikas-id")
	assert.NoError(t, err)
}

/*
========================================
 TEST UPDATE
========================================
*/

func TestIkasService_Update_Async(t *testing.T) {
	repo := &mockIkasRepo{}
	service := NewIkasService(repo, nil)

	val := 10.0
	req := dto.UpdateIkasRequest{
		TargetNilai: &val,
	}

	// Update now returns nil error on nil producer
	err := service.Update("ikas-id", req)
	assert.NoError(t, err)
}

/*
========================================
 TEST IMPORT
========================================
*/

func TestIkasService_ImportFromExcel_Async(t *testing.T) {
	repo := &mockIkasRepo{}
	service := NewIkasService(repo, nil)

	id, err := service.ImportFromExcel([]byte("fake excel"))
	assert.NoError(t, err) // Should work and return generated ID
	assert.NotEmpty(t, id)
}
