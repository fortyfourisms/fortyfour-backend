package services

import (
	"testing"

	"fortyfour-backend/internal/dto"

	"github.com/stretchr/testify/assert"
)

/*
========================================
 MOCK REPOSITORY
========================================
*/

type mockIkasRepo struct{}

func (m *mockIkasRepo) CreateIdentifikasi(id string, d *dto.CreateIdentifikasiData) (float64, error) {
	return 80, nil
}

func (m *mockIkasRepo) FindPerusahaanByName(namaPerusahaan string) (string, error) {
	return "", nil
}

func (m *mockIkasRepo) CreateProteksi(id string, d *dto.CreateProteksiData) (float64, error) {
	return 70, nil
}

func (m *mockIkasRepo) CreateDeteksi(id string, d *dto.CreateDeteksiData) (float64, error) {
	return 90, nil
}

func (m *mockIkasRepo) CreateGulih(id string, d *dto.CreateGulihData) (float64, error) {
	return 60, nil
}

func (m *mockIkasRepo) Create(
	req dto.CreateIkasRequest,
	id string,
	nilai float64,
	idI, idP, idD, idG string,
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
		Identifikasi: &dto.IdentifikasiInIkas{
			ID:                "iden-id",
			NilaiIdentifikasi: 80,
		},
		Proteksi: &dto.ProteksiInIkas{
			ID:            "prot-id",
			NilaiProteksi: 70,
		},
		Deteksi: &dto.DeteksiInIkas{
			ID:           "det-id",
			NilaiDeteksi: 90,
		},
		Gulih: &dto.GulihInIkas{
			ID:         "gul-id",
			NilaiGulih: 60,
		},
	}, nil
}

func (m *mockIkasRepo) UpdateIdentifikasi(id string, d *dto.UpdateIdentifikasiData) (float64, error) {
	return 85, nil
}

func (m *mockIkasRepo) UpdateProteksi(id string, d *dto.UpdateProteksiData) (float64, error) {
	return 75, nil
}

func (m *mockIkasRepo) UpdateDeteksi(id string, d *dto.UpdateDeteksiData) (float64, error) {
	return 95, nil
}

func (m *mockIkasRepo) UpdateGulih(id string, d *dto.UpdateGulihData) (float64, error) {
	return 65, nil
}

func (m *mockIkasRepo) Update(id string, req dto.UpdateIkasRequest) error {
	return nil
}

func (m *mockIkasRepo) Delete(id string) error {
	return nil
}

func (m *mockIkasRepo) ImportFromExcel(raw []byte) (*dto.IkasResponse, error) {
	return &dto.IkasResponse{ID: "imp-1", Responden: "imported"}, nil
}

func (m *mockIkasRepo) ParseExcelForImport(b []byte) (*dto.CreateIkasRequest, error) {
	return &dto.CreateIkasRequest{
		Identifikasi: &dto.CreateIdentifikasiData{},
		Proteksi:     &dto.CreateProteksiData{},
		Deteksi:      &dto.CreateDeteksiData{},
		Gulih:        &dto.CreateGulihData{},
	}, nil
}

/*
========================================
 TEST CREATE
========================================
*/

func TestIkasService_Create_Success(t *testing.T) {
	repo := &mockIkasRepo{}
	service := NewIkasService(repo)

	req := dto.CreateIkasRequest{
		IDPerusahaan: "1",
		Identifikasi: &dto.CreateIdentifikasiData{},
		Proteksi:     &dto.CreateProteksiData{},
		Deteksi:      &dto.CreateDeteksiData{},
		Gulih:        &dto.CreateGulihData{},
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
	service := NewIkasService(repo)

	val := 10.0
	req := dto.UpdateIkasRequest{
		Identifikasi: &dto.UpdateIdentifikasiData{
			NilaiSubdomain1: &val,
		},
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
	service := NewIkasService(repo)

	resp, err := service.ImportFromExcel([]byte("fake excel"))

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
