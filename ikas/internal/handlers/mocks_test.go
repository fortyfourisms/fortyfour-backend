package handlers

import (
	"github.com/stretchr/testify/mock"
	"ikas/internal/dto"
	"ikas/internal/repository"
)

type mockIkasRepository struct {
	mock.Mock
}

func (m *mockIkasRepository) Create(req dto.CreateIkasRequest, id string, nilai float64) error {
	args := m.Called(req, id, nilai)
	return args.Error(0)
}

func (m *mockIkasRepository) GetAll() ([]dto.IkasResponse, error) {
	args := m.Called()
	return args.Get(0).([]dto.IkasResponse), args.Error(1)
}

func (m *mockIkasRepository) GetByID(id string) (*dto.IkasResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.IkasResponse), args.Error(1)
}

func (m *mockIkasRepository) Update(id string, req dto.UpdateIkasRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *mockIkasRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockIkasRepository) CheckExistsByPerusahaanID(perusahaanID string) (bool, error) {
	args := m.Called(perusahaanID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *mockIkasRepository) CheckExistsByPerusahaanIDAndYear(id string, year int) (bool, error) {
	args := m.Called(id, year)
	return args.Get(0).(bool), args.Error(1)
}

func (m *mockIkasRepository) CheckOwnership(ikasID string, perusahaanID string) (bool, error) {
	args := m.Called(ikasID, perusahaanID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *mockIkasRepository) FindPerusahaanByName(namaPerusahaan string) (string, error) {
	args := m.Called(namaPerusahaan)
	return args.Get(0).(string), args.Error(1)
}

func (m *mockIkasRepository) GetIDByPerusahaanID(param string) (string, error) {
	args := m.Called(param)
	return args.Get(0).(string), args.Error(1)
}

func (m *mockIkasRepository) GetByPerusahaan(perusahaanID string) ([]dto.IkasResponse, error) {
	args := m.Called(perusahaanID)
	return args.Get(0).([]dto.IkasResponse), args.Error(1)
}

func (m *mockIkasRepository) ParseExcelForImport(fileData []byte) (*dto.ParsedExcelData, error) {
	args := m.Called(fileData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ParsedExcelData), args.Error(1)
}

func (m *mockIkasRepository) UpdateValidationStatus(id string, status bool) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *mockIkasRepository) IsLocked(id string) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

var _ repository.IkasRepositoryInterface = (*mockIkasRepository)(nil)
