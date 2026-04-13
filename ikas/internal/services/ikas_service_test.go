package services

import (
	"context"
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

func (m *mockIkasRepo) CheckExistsByPerusahaanID(idPerusahaan string) (bool, error) {
	if idPerusahaan == "perusahaan-ada" {
		return true, nil
	}
	return false, nil
}

func (m *mockIkasRepo) CheckExistsByPerusahaanIDAndYear(id string, year int) (bool, error) {
	if id == "perusahaan-ada" {
		return true, nil
	}
	return false, nil
}

func (m *mockIkasRepo) GetIDByPerusahaanID(idPerusahaan string) (string, error) {
	return "ikas-id", nil
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

func (m *mockIkasRepo) ParseExcelForImport(b []byte) (*dto.ParsedExcelData, error) {
	return &dto.ParsedExcelData{}, nil
}

func (m *mockIkasRepo) GetByPerusahaan(perusahaanID string) ([]dto.IkasResponse, error) {
	return []dto.IkasResponse{}, nil
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
	err := service.Create(context.Background(), req, "ikas-id", "test-user")
	assert.NoError(t, err)
}

func TestIkasService_Create_Duplicate(t *testing.T) {
	repo := &mockIkasRepo{}
	service := NewIkasService(repo, nil)

	req := dto.CreateIkasRequest{
		IDPerusahaan: "perusahaan-ada",
	}

	err := service.Create(context.Background(), req, "ikas-id", "test-user")
	assert.Error(t, err)
	// Update expected error message to match yearly validation
	assert.Contains(t, err.Error(), "sudah ada")
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
	err := service.Update(context.Background(), "ikas-id", req, "test-user", "admin", "")
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

	id, err := service.ImportFromExcel(context.Background(), []byte("fake excel"), "test-user")
	assert.NoError(t, err) // Should work and return generated ID
	assert.NotEmpty(t, id)
}
