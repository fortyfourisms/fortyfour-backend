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
	return "e39d8349-6aad-4315-998c-4d3799e49a92", nil
}

func (m *mockIkasRepo) CheckExistsByPerusahaanID(idPerusahaan string) (bool, error) {
	// Always return true for valid-looking UUIDs in tests
	return true, nil
}

func (m *mockIkasRepo) CheckExistsByPerusahaanIDAndYear(id string, year int) (bool, error) {
	if id == "e39d8349-6aad-4315-998c-4d3799e49a92" {
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
		Tanggal:         "2026-01-01",
		Perusahaan:      &dto.PerusahaanInIkas{ID: "e39d8349-6aad-4315-998c-4d3799e49a92"},
	}, nil
}

func (m *mockIkasRepo) GetLatestByPerusahaan(perusahaanID string) (*dto.IkasResponse, error) {
	return nil, nil
}

func (m *mockIkasRepo) Update(id string, req dto.UpdateIkasRequest) error {
	return nil
}

func (m *mockIkasRepo) Delete(id string) error {
	return nil
}

func (m *mockIkasRepo) ParseExcelForImport(b []byte) (*dto.ParsedExcelData, error) {
	return &dto.ParsedExcelData{
		IkasRequest: dto.CreateIkasRequest{
			IDPerusahaan: "7d7ae6c3-eae1-4e66-bc3e-c75af9c9302c",
			Responden:    "Test Responden",
			Telepon:      "081234567890",
			Jabatan:      "CIO",
			TargetNilai:  3.0,
		},
	}, nil
}

func (m *mockIkasRepo) GetByPerusahaan(perusahaanID string) ([]dto.IkasResponse, error) {
	return []dto.IkasResponse{}, nil
}

func (m *mockIkasRepo) SaveImportedData(id string, data *dto.CreateIkasRequest) error {
	return nil
}

func (m *mockIkasRepo) CheckOwnership(ikasID string, perusahaanID string) (bool, error) {
	return true, nil
}

func (m *mockIkasRepo) UpdateValidationStatus(id string, status bool) error {
	return nil
}

func (m *mockIkasRepo) CreateInitial(sourceID, targetID, targetDate string) error {
	return nil
}

func (m *mockIkasRepo) UpdateDomainLinks(ikasID, identifikasiID, proteksiID, deteksiID, gulihID string) error {
	return nil
}

func (m *mockIkasRepo) IsLocked(id string) (bool, error) {
	return false, nil
}

/*
========================================
 TEST CREATE
========================================
*/

func TestIkasService_Create_Success(t *testing.T) {
	repo := &mockIkasRepo{}
	service := NewIkasService(repo, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	req := dto.CreateIkasRequest{
		IDPerusahaan: "7d7ae6c3-eae1-4e66-bc3e-c75af9c9302c",
		Responden:    "Test Responden",
		Telepon:      "081234567890",
		Jabatan:      "CIO",
		TargetNilai:  3.0,
	}

	err := service.Create(context.Background(), req, "ikas-id", "test-user")
	assert.NoError(t, err)
}

func TestIkasService_Create_Duplicate(t *testing.T) {
	repo := &mockIkasRepo{}
	service := NewIkasService(repo, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	req := dto.CreateIkasRequest{
		IDPerusahaan: "e39d8349-6aad-4315-998c-4d3799e49a92", // ID that returns true for duplicate check in mock
		Responden:    "Test Responden",
		Telepon:      "081234567890",
		Jabatan:      "CIO",
		TargetNilai:  3.0,
	}

	err := service.Create(context.Background(), req, "ikas-id", "test-user")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sudah ada")
}

/*
========================================
 TEST UPDATE
========================================
*/

func TestIkasService_Update_Async(t *testing.T) {
	repo := &mockIkasRepo{}
	service := NewIkasService(repo, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	val := 3.0
	req := dto.UpdateIkasRequest{
		TargetNilai: &val,
	}

	_, err := service.Update(context.Background(), "ikas-id", req, "test-user", "admin", "")
	assert.NoError(t, err)
}

func TestIkasService_GetAll_Admin(t *testing.T) {
	repo := &mockIkasRepo{}
	service := NewIkasService(repo, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	_, err := service.GetAll("admin")
	assert.NoError(t, err)
}

func TestIkasService_GetAll_NonAdmin(t *testing.T) {
	repo := &mockIkasRepo{}
	service := NewIkasService(repo, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	_, err := service.GetAll("user")
	assert.Error(t, err)
}

/*
========================================
 TEST IMPORT
========================================
*/

func TestIkasService_ImportFromExcel_Async(t *testing.T) {
	repo := &mockIkasRepo{}
	service := NewIkasService(repo, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	id, err := service.ImportFromExcel(context.Background(), []byte("fake excel"), "test-user")
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
}
