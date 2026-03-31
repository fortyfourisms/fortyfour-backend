package services

import (
	"errors"
	"testing"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ════════════════════════════════════════════════════════════════════════════
// MOCK SE SERVICE
// ════════════════════════════════════════════════════════════════════════════

type mockSEServiceForExport struct {
	mock.Mock
}

func (m *mockSEServiceForExport) GetAll() ([]dto.SEResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.SEResponse), args.Error(1)
}

func (m *mockSEServiceForExport) GetByID(id string) (*dto.SEResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SEResponse), args.Error(1)
}

func (m *mockSEServiceForExport) GetByPerusahaan(idPerusahaan string) ([]dto.SEResponse, error) {
	args := m.Called(idPerusahaan)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.SEResponse), args.Error(1)
}

func (m *mockSEServiceForExport) Create(req dto.CreateSERequest) (*dto.SEResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SEResponse), args.Error(1)
}

func (m *mockSEServiceForExport) Update(id string, req dto.UpdateSERequest) (*dto.SEResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SEResponse), args.Error(1)
}

func (m *mockSEServiceForExport) Delete(id string) error {
	return m.Called(id).Error(0)
}

// ════════════════════════════════════════════════════════════════════════════
// MOCK CSIRT SERVICE
// ════════════════════════════════════════════════════════════════════════════

type mockCsirtServiceForExport struct {
	mock.Mock
}

func (m *mockCsirtServiceForExport) GetAll() ([]dto.CsirtResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.CsirtResponse), args.Error(1)
}

func (m *mockCsirtServiceForExport) GetByID(id string) (*dto.CsirtResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CsirtResponse), args.Error(1)
}

func (m *mockCsirtServiceForExport) GetByPerusahaan(idPerusahaan string) ([]dto.CsirtResponse, error) {
	args := m.Called(idPerusahaan)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.CsirtResponse), args.Error(1)
}

func (m *mockCsirtServiceForExport) Create(req dto.CreateCsirtRequest) (*models.Csirt, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Csirt), args.Error(1)
}

func (m *mockCsirtServiceForExport) Update(id string, req dto.UpdateCsirtRequest) (*models.Csirt, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Csirt), args.Error(1)
}

func (m *mockCsirtServiceForExport) Delete(id string) error {
	return m.Called(id).Error(0)
}

// ════════════════════════════════════════════════════════════════════════════
// HELPER — dummy SE/CSIRT data
// ════════════════════════════════════════════════════════════════════════════

func dummySEResponse(id, idPerusahaan, namaPerusahaan string) dto.SEResponse {
	return dto.SEResponse{
		ID:                      id,
		IDPerusahaan:            idPerusahaan,
		NamaSE:                  "Sistem Informasi",
		IpSE:                    "192.168.1.1",
		AsNumberSE:              "AS12345",
		PengelolaSE:             "IT Dept",
		NilaiInvestasi:          "A",
		AnggaranOperasional:     "A",
		KepatuhanPeraturan:      "A",
		TeknikKriptografi:       "A",
		JumlahPengguna:          "A",
		DataPribadi:             "A",
		KlasifikasiData:         "A",
		KekritisanProses:        "A",
		DampakKegagalan:         "A",
		PotensiKerugiandanDampakNegatif: "A",
		TotalBobot:              50,
		KategoriSE:              "Strategis",
		Perusahaan: &dto.PerusahaanMiniResponse{
			ID:             idPerusahaan,
			NamaPerusahaan: namaPerusahaan,
		},
	}
}

func dummyCsirtResponse(id, idPerusahaan, namaPerusahaan string) dto.CsirtResponse {
	return dto.CsirtResponse{
		ID:        id,
		NamaCsirt: "CSIRT Test",
		WebCsirt:  "https://csirt.test",
		Perusahaan: dto.PerusahaanResponse{
			ID:             idPerusahaan,
			NamaPerusahaan: namaPerusahaan,
		},
	}
}

// ════════════════════════════════════════════════════════════════════════════
// SE EXPORT SERVICE TESTS
// ════════════════════════════════════════════════════════════════════════════

func TestSEExportService_ExportAllPDF_Success(t *testing.T) {
	mockSvc := new(mockSEServiceForExport)
	mockSvc.On("GetAll").Return([]dto.SEResponse{
		dummySEResponse("se-1", "p-1", "Perusahaan A"),
		dummySEResponse("se-2", "p-2", "Perusahaan B"),
	}, nil)

	svc := NewSEExportService(mockSvc)
	pdfBytes, err := svc.ExportAllPDF()

	assert.NoError(t, err)
	assert.NotNil(t, pdfBytes)
	assert.Greater(t, len(pdfBytes), 0)
	mockSvc.AssertExpectations(t)
}

func TestSEExportService_ExportAllPDF_ServiceError(t *testing.T) {
	mockSvc := new(mockSEServiceForExport)
	mockSvc.On("GetAll").Return(nil, errors.New("db error"))

	svc := NewSEExportService(mockSvc)
	pdfBytes, err := svc.ExportAllPDF()

	assert.Error(t, err)
	assert.Nil(t, pdfBytes)
	assert.Equal(t, "db error", err.Error())
}

func TestSEExportService_ExportAllPDF_EmptyData(t *testing.T) {
	mockSvc := new(mockSEServiceForExport)
	mockSvc.On("GetAll").Return([]dto.SEResponse{}, nil)

	svc := NewSEExportService(mockSvc)
	pdfBytes, err := svc.ExportAllPDF()

	assert.Error(t, err)
	assert.Nil(t, pdfBytes)
	assert.Equal(t, "tidak ada data SE untuk diexport", err.Error())
}

func TestSEExportService_ExportByPerusahaanPDF_Success(t *testing.T) {
	mockSvc := new(mockSEServiceForExport)
	mockSvc.On("GetByPerusahaan", "p-abc").Return([]dto.SEResponse{
		dummySEResponse("se-1", "p-abc", "Perusahaan ABC"),
	}, nil)

	svc := NewSEExportService(mockSvc)
	pdfBytes, err := svc.ExportByPerusahaanPDF("p-abc")

	assert.NoError(t, err)
	assert.NotNil(t, pdfBytes)
	assert.Greater(t, len(pdfBytes), 0)
	mockSvc.AssertExpectations(t)
}

func TestSEExportService_ExportByPerusahaanPDF_EmptyID(t *testing.T) {
	svc := NewSEExportService(new(mockSEServiceForExport))
	pdfBytes, err := svc.ExportByPerusahaanPDF("   ")

	assert.Error(t, err)
	assert.Nil(t, pdfBytes)
	assert.Equal(t, "id_perusahaan wajib diisi", err.Error())
}

func TestSEExportService_ExportByPerusahaanPDF_ServiceError(t *testing.T) {
	mockSvc := new(mockSEServiceForExport)
	mockSvc.On("GetByPerusahaan", "p-abc").Return(nil, errors.New("db error"))

	svc := NewSEExportService(mockSvc)
	pdfBytes, err := svc.ExportByPerusahaanPDF("p-abc")

	assert.Error(t, err)
	assert.Nil(t, pdfBytes)
}

func TestSEExportService_ExportByPerusahaanPDF_EmptyData(t *testing.T) {
	mockSvc := new(mockSEServiceForExport)
	mockSvc.On("GetByPerusahaan", "p-abc").Return([]dto.SEResponse{}, nil)

	svc := NewSEExportService(mockSvc)
	pdfBytes, err := svc.ExportByPerusahaanPDF("p-abc")

	assert.Error(t, err)
	assert.Nil(t, pdfBytes)
	assert.Equal(t, "tidak ada data SE untuk perusahaan ini", err.Error())
}

func TestSEExportService_ExportByIDPDF_Success(t *testing.T) {
	expected := dummySEResponse("se-1", "p-abc", "Perusahaan ABC")
	mockSvc := new(mockSEServiceForExport)
	mockSvc.On("GetByID", "se-1").Return(&expected, nil)

	svc := NewSEExportService(mockSvc)
	se, pdfBytes, err := svc.ExportByIDPDF("se-1")

	assert.NoError(t, err)
	assert.NotNil(t, se)
	assert.Equal(t, "se-1", se.ID)
	assert.NotNil(t, pdfBytes)
	assert.Greater(t, len(pdfBytes), 0)
	mockSvc.AssertExpectations(t)
}

func TestSEExportService_ExportByIDPDF_EmptyID(t *testing.T) {
	svc := NewSEExportService(new(mockSEServiceForExport))
	se, pdfBytes, err := svc.ExportByIDPDF("  ")

	assert.Error(t, err)
	assert.Nil(t, se)
	assert.Nil(t, pdfBytes)
	assert.Equal(t, "id wajib diisi", err.Error())
}

func TestSEExportService_ExportByIDPDF_NotFound(t *testing.T) {
	mockSvc := new(mockSEServiceForExport)
	mockSvc.On("GetByID", "not-exist").Return(nil, errors.New("not found"))

	svc := NewSEExportService(mockSvc)
	se, pdfBytes, err := svc.ExportByIDPDF("not-exist")

	assert.Error(t, err)
	assert.Nil(t, se)
	assert.Nil(t, pdfBytes)
	assert.Equal(t, "data tidak ditemukan", err.Error())
}

// ════════════════════════════════════════════════════════════════════════════
// CSIRT EXPORT SERVICE TESTS
// ════════════════════════════════════════════════════════════════════════════

func TestCsirtExportService_ExportAllPDF_Success(t *testing.T) {
	mockSvc := new(mockCsirtServiceForExport)
	mockSvc.On("GetAll").Return([]dto.CsirtResponse{
		dummyCsirtResponse("csirt-1", "p-1", "Perusahaan A"),
		dummyCsirtResponse("csirt-2", "p-2", "Perusahaan B"),
	}, nil)

	svc := NewCsirtExportService(mockSvc)
	pdfBytes, err := svc.ExportAllPDF()

	assert.NoError(t, err)
	assert.NotNil(t, pdfBytes)
	assert.Greater(t, len(pdfBytes), 0)
	mockSvc.AssertExpectations(t)
}

func TestCsirtExportService_ExportAllPDF_ServiceError(t *testing.T) {
	mockSvc := new(mockCsirtServiceForExport)
	mockSvc.On("GetAll").Return(nil, errors.New("db error"))

	svc := NewCsirtExportService(mockSvc)
	pdfBytes, err := svc.ExportAllPDF()

	assert.Error(t, err)
	assert.Nil(t, pdfBytes)
	assert.Equal(t, "db error", err.Error())
}

func TestCsirtExportService_ExportAllPDF_EmptyData(t *testing.T) {
	mockSvc := new(mockCsirtServiceForExport)
	mockSvc.On("GetAll").Return([]dto.CsirtResponse{}, nil)

	svc := NewCsirtExportService(mockSvc)
	pdfBytes, err := svc.ExportAllPDF()

	assert.Error(t, err)
	assert.Nil(t, pdfBytes)
	assert.Equal(t, "tidak ada data CSIRT untuk diexport", err.Error())
}

func TestCsirtExportService_ExportByPerusahaanPDF_Success(t *testing.T) {
	mockSvc := new(mockCsirtServiceForExport)
	mockSvc.On("GetByPerusahaan", "p-abc").Return([]dto.CsirtResponse{
		dummyCsirtResponse("csirt-1", "p-abc", "Perusahaan ABC"),
	}, nil)

	svc := NewCsirtExportService(mockSvc)
	pdfBytes, err := svc.ExportByPerusahaanPDF("p-abc")

	assert.NoError(t, err)
	assert.NotNil(t, pdfBytes)
	assert.Greater(t, len(pdfBytes), 0)
	mockSvc.AssertExpectations(t)
}

func TestCsirtExportService_ExportByPerusahaanPDF_EmptyID(t *testing.T) {
	svc := NewCsirtExportService(new(mockCsirtServiceForExport))
	pdfBytes, err := svc.ExportByPerusahaanPDF("")

	assert.Error(t, err)
	assert.Nil(t, pdfBytes)
	assert.Equal(t, "id_perusahaan wajib diisi", err.Error())
}

func TestCsirtExportService_ExportByPerusahaanPDF_ServiceError(t *testing.T) {
	mockSvc := new(mockCsirtServiceForExport)
	mockSvc.On("GetByPerusahaan", "p-abc").Return(nil, errors.New("db error"))

	svc := NewCsirtExportService(mockSvc)
	pdfBytes, err := svc.ExportByPerusahaanPDF("p-abc")

	assert.Error(t, err)
	assert.Nil(t, pdfBytes)
}

func TestCsirtExportService_ExportByPerusahaanPDF_EmptyData(t *testing.T) {
	mockSvc := new(mockCsirtServiceForExport)
	mockSvc.On("GetByPerusahaan", "p-abc").Return([]dto.CsirtResponse{}, nil)

	svc := NewCsirtExportService(mockSvc)
	pdfBytes, err := svc.ExportByPerusahaanPDF("p-abc")

	assert.Error(t, err)
	assert.Nil(t, pdfBytes)
	assert.Equal(t, "tidak ada data CSIRT untuk perusahaan ini", err.Error())
}

func TestCsirtExportService_ExportByIDPDF_Success(t *testing.T) {
	expected := dummyCsirtResponse("csirt-1", "p-abc", "Perusahaan ABC")
	mockSvc := new(mockCsirtServiceForExport)
	mockSvc.On("GetByID", "csirt-1").Return(&expected, nil)

	svc := NewCsirtExportService(mockSvc)
	csirt, pdfBytes, err := svc.ExportByIDPDF("csirt-1")

	assert.NoError(t, err)
	assert.NotNil(t, csirt)
	assert.Equal(t, "csirt-1", csirt.ID)
	assert.NotNil(t, pdfBytes)
	assert.Greater(t, len(pdfBytes), 0)
	mockSvc.AssertExpectations(t)
}

func TestCsirtExportService_ExportByIDPDF_EmptyID(t *testing.T) {
	svc := NewCsirtExportService(new(mockCsirtServiceForExport))
	csirt, pdfBytes, err := svc.ExportByIDPDF("")

	assert.Error(t, err)
	assert.Nil(t, csirt)
	assert.Nil(t, pdfBytes)
	assert.Equal(t, "id wajib diisi", err.Error())
}

func TestCsirtExportService_ExportByIDPDF_NotFound(t *testing.T) {
	mockSvc := new(mockCsirtServiceForExport)
	mockSvc.On("GetByID", "not-exist").Return(nil, errors.New("not found"))

	svc := NewCsirtExportService(mockSvc)
	csirt, pdfBytes, err := svc.ExportByIDPDF("not-exist")

	assert.Error(t, err)
	assert.Nil(t, csirt)
	assert.Nil(t, pdfBytes)
	assert.Equal(t, "data tidak ditemukan", err.Error())
}