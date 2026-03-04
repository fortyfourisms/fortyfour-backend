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

func (m *mockIkasRepo) CreateIdentifikasi(d *dto.CreateIdentifikasiData) (int64, float64, error) {
	return 1, 80, nil
}

func (m *mockIkasRepo) FindPerusahaanByName(namaPerusahaan string) (string, error) {
	return "", nil
}

func (m *mockIkasRepo) CreateProteksi(d *dto.CreateProteksiData) (int64, float64, error) {
	return 2, 70, nil
}

func (m *mockIkasRepo) CreateDeteksi(d *dto.CreateDeteksiData) (int64, float64, error) {
	return 3, 90, nil
}

func (m *mockIkasRepo) CreateGulih(d *dto.CreateGulihData) (int64, float64, error) {
	return 4, 60, nil
}

func (m *mockIkasRepo) Create(
	req dto.CreateIkasRequest,
	id string,
	nilai float64,
	idI, idP, idD, idG int,
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
			ID:                1,
			NilaiIdentifikasi: 80,
		},
		Proteksi: &dto.ProteksiInIkas{
			ID:            2,
			NilaiProteksi: 70,
		},
		Deteksi: &dto.DeteksiInIkas{
			ID:           3,
			NilaiDeteksi: 90,
		},
		Gulih: &dto.GulihInIkas{
			ID:         4,
			NilaiGulih: 60,
		},
	}, nil
}

func (m *mockIkasRepo) UpdateIdentifikasi(id int, d *dto.UpdateIdentifikasiData) (float64, error) {
	return 85, nil
}

func (m *mockIkasRepo) UpdateProteksi(id int, d *dto.UpdateProteksiData) (float64, error) {
	return 75, nil
}

func (m *mockIkasRepo) UpdateDeteksi(id int, d *dto.UpdateDeteksiData) (float64, error) {
	return 95, nil
}

func (m *mockIkasRepo) UpdateGulih(id int, d *dto.UpdateGulihData) (float64, error) {
	return 65, nil
}

func (m *mockIkasRepo) Update(id string, req dto.UpdateIkasRequest) error {
	return nil
}

func (m *mockIkasRepo) Delete(id string) error {
	return nil
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
	service := NewIkasService(repo, nil)

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
	service := NewIkasService(repo, nil)

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
	service := NewIkasService(repo, nil)

	resp, err := service.ImportFromExcel([]byte("fake excel"))

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
