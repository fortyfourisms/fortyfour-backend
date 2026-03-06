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

/*
========================================
 TEST CREATE
========================================
*/

func TestIkasService_Create_Success(t *testing.T) {
	repo := &mockIkasRepo{}
	service := NewIkasService(repo, nil)

	req := dto.CreateIkasRequest{
		IDPerusahaan: "1",
	}

	err := service.Create(req, "ikas-id")

	assert.NoError(t, err)
}

/*
========================================
 TEST UPDATE
========================================
*/

func TestIkasService_Update_Recalculate(t *testing.T) {
	repo := &mockIkasRepo{}
	service := NewIkasService(repo, nil)

	val := 10.0
	req := dto.UpdateIkasRequest{
		TargetNilai: &val,
	}

	resp, err := service.Update("ikas-id", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

/*
========================================
 TEST IMPORT
========================================
*/

func TestIkasService_ImportFromExcel(t *testing.T) {
	repo := &mockIkasRepo{}
	service := NewIkasService(repo, nil)

	resp, err := service.ImportFromExcel([]byte("fake excel"))

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
