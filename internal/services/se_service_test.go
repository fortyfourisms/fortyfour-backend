package services

import (
	"errors"
	"testing"

	"fortyfour-backend/internal/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/*
=====================================
 MOCK SE REPOSITORY
=====================================
*/

type MockSERepository struct {
	mock.Mock
}

func (m *MockSERepository) Create(req dto.CreateSERequest, id string, totalBobot int, kategori string) error {
	args := m.Called(req, id, totalBobot, kategori)
	return args.Error(0)
}

func (m *MockSERepository) GetAll() ([]dto.SEResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.SEResponse), args.Error(1)
}

func (m *MockSERepository) GetByID(id string) (*dto.SEResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SEResponse), args.Error(1)
}

func (m *MockSERepository) Update(id string, req dto.UpdateSERequest, totalBobot int, kategori string) error {
	args := m.Called(id, req, totalBobot, kategori)
	return args.Error(0)
}

func (m *MockSERepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSERepository) GetByPerusahaan(idPerusahaan string) ([]dto.SEResponse, error) {
	args := m.Called(idPerusahaan)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.SEResponse), args.Error(1)
}

/*
=====================================
 HELPER FUNCTIONS
=====================================
*/

func setupSEService() (SEService, *MockSERepository) {
	mockRepo := new(MockSERepository)
	service := NewSEService(mockRepo, nil, nil)
	return service, mockRepo
}

func createValidSERequest() dto.CreateSERequest {
	return dto.CreateSERequest{
		IDPerusahaan: "perusahaan-123",
		NamaSE:       "Sistem Informasi Keuangan",
		IpSE:         "192.168.1.1",
		AsNumberSE:   "AS12345",
		PengelolaSE:  "IT Department",
		// Semua A = 10 x 5 = 50 (Strategis)
		NilaiInvestasi:                  "A",
		AnggaranOperasional:             "A",
		KepatuhanPeraturan:              "A",
		TeknikKriptografi:               "A",
		JumlahPengguna:                  "A",
		DataPribadi:                     "A",
		KlasifikasiData:                 "A",
		KekritisanProses:                "A",
		DampakKegagalan:                 "A",
		PotensiKerugiandanDampakNegatif: "A",
	}
}

/*
=====================================
 TEST CREATE - VALIDATION
=====================================
*/

func TestSEService_Create_MissingIDPerusahaan(t *testing.T) {
	service, _ := setupSEService()

	req := createValidSERequest()
	req.IDPerusahaan = ""

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "id_perusahaan wajib diisi", err.Error())
}

func TestSEService_Create_MissingNamaSE(t *testing.T) {
	service, _ := setupSEService()

	req := createValidSERequest()
	req.NamaSE = ""

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "nama_se wajib diisi", err.Error())
}

func TestSEService_Create_MissingIpSE(t *testing.T) {
	service, _ := setupSEService()

	req := createValidSERequest()
	req.IpSE = ""

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "ip_se wajib diisi", err.Error())
}

func TestSEService_Create_MissingAsNumber(t *testing.T) {
	service, _ := setupSEService()

	req := createValidSERequest()
	req.AsNumberSE = ""

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "as_number_se wajib diisi", err.Error())
}

func TestSEService_Create_MissingPengelola(t *testing.T) {
	service, _ := setupSEService()

	req := createValidSERequest()
	req.PengelolaSE = ""

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "pengelola_se wajib diisi", err.Error())
}

func TestSEService_Create_InvalidKarakteristik(t *testing.T) {
	service, _ := setupSEService()

	req := createValidSERequest()
	req.NilaiInvestasi = "X"

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "jawaban harus A, B, atau C")
}

/*
=====================================
 TEST CREATE - KATEGORISASI LOGIC
=====================================
*/

func TestSEService_Create_KategoriStrategis(t *testing.T) {
	service, mockRepo := setupSEService()

	// Total bobot = 50 (10 x A = 10 x 5) → Strategis
	req := createValidSERequest()

	expectedResponse := &dto.SEResponse{
		ID:         "new-se-id",
		NamaSE:     req.NamaSE,
		TotalBobot: 50,
		KategoriSE: "Strategis",
	}

	mockRepo.On("Create", req, mock.AnythingOfType("string"), 50, "Strategis").Return(nil)
	mockRepo.On("GetByID", mock.AnythingOfType("string")).Return(expectedResponse, nil)

	result, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 50, result.TotalBobot)
	assert.Equal(t, "Strategis", result.KategoriSE)

	mockRepo.AssertExpectations(t)
}

func TestSEService_Create_KategoriTinggi(t *testing.T) {
	service, mockRepo := setupSEService()

	req := createValidSERequest()
	// Mix A, B, C untuk total bobot di range 16-34 (Tinggi)
	// 5A (25) + 3B (6) + 2C (2) = 33 → Tinggi
	req.NilaiInvestasi = "A"
	req.AnggaranOperasional = "A"
	req.KepatuhanPeraturan = "A"
	req.TeknikKriptografi = "A"
	req.JumlahPengguna = "A"
	req.DataPribadi = "B"
	req.KlasifikasiData = "B"
	req.KekritisanProses = "B"
	req.DampakKegagalan = "C"
	req.PotensiKerugiandanDampakNegatif = "C"

	expectedResponse := &dto.SEResponse{
		ID:         "new-se-id",
		TotalBobot: 33,
		KategoriSE: "Tinggi",
	}

	mockRepo.On("Create", req, mock.AnythingOfType("string"), 33, "Tinggi").Return(nil)
	mockRepo.On("GetByID", mock.AnythingOfType("string")).Return(expectedResponse, nil)

	result, err := service.Create(req)

	assert.NoError(t, err)
	assert.Equal(t, 33, result.TotalBobot)
	assert.Equal(t, "Tinggi", result.KategoriSE)

	mockRepo.AssertExpectations(t)
}

func TestSEService_Create_KategoriRendah(t *testing.T) {
	service, mockRepo := setupSEService()

	req := createValidSERequest()
	// Semua C = 10 x 1 = 10 → Rendah
	req.NilaiInvestasi = "C"
	req.AnggaranOperasional = "C"
	req.KepatuhanPeraturan = "C"
	req.TeknikKriptografi = "C"
	req.JumlahPengguna = "C"
	req.DataPribadi = "C"
	req.KlasifikasiData = "C"
	req.KekritisanProses = "C"
	req.DampakKegagalan = "C"
	req.PotensiKerugiandanDampakNegatif = "C"

	expectedResponse := &dto.SEResponse{
		ID:         "new-se-id",
		TotalBobot: 10,
		KategoriSE: "Rendah",
	}

	mockRepo.On("Create", req, mock.AnythingOfType("string"), 10, "Rendah").Return(nil)
	mockRepo.On("GetByID", mock.AnythingOfType("string")).Return(expectedResponse, nil)

	result, err := service.Create(req)

	assert.NoError(t, err)
	assert.Equal(t, 10, result.TotalBobot)
	assert.Equal(t, "Rendah", result.KategoriSE)

	mockRepo.AssertExpectations(t)
}

/*
=====================================
 TEST CREATE - SUCCESS CASES
=====================================
*/

func TestSEService_Create_Success(t *testing.T) {
	service, mockRepo := setupSEService()

	req := createValidSERequest()

	expectedResponse := &dto.SEResponse{
		ID:           "new-se-id",
		IDPerusahaan: req.IDPerusahaan,
		NamaSE:       req.NamaSE,
		IpSE:         req.IpSE,
		TotalBobot:   50,
		KategoriSE:   "Strategis",
	}

	mockRepo.On("Create", req, mock.AnythingOfType("string"), 50, "Strategis").Return(nil)
	mockRepo.On("GetByID", mock.AnythingOfType("string")).Return(expectedResponse, nil)

	result, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.NamaSE, result.NamaSE)
	assert.Equal(t, req.IpSE, result.IpSE)

	mockRepo.AssertExpectations(t)
}

func TestSEService_Create_RepositoryError(t *testing.T) {
	service, mockRepo := setupSEService()

	req := createValidSERequest()

	mockRepo.On("Create", req, mock.AnythingOfType("string"), 50, "Strategis").
		Return(errors.New("database error"))

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockRepo.AssertExpectations(t)
}

/*
=====================================
 TEST GET ALL
=====================================
*/

func TestSEService_GetAll_Success(t *testing.T) {
	service, mockRepo := setupSEService()

	expectedData := []dto.SEResponse{
		{ID: "se-1", NamaSE: "SE 1", KategoriSE: "Strategis", TotalBobot: 50},
		{ID: "se-2", NamaSE: "SE 2", KategoriSE: "Tinggi", TotalBobot: 30},
	}

	mockRepo.On("GetAll").Return(expectedData, nil)

	result, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "se-1", result[0].ID)
	assert.Equal(t, "Strategis", result[0].KategoriSE)

	mockRepo.AssertExpectations(t)
}

func TestSEService_GetAll_EmptyResult(t *testing.T) {
	service, mockRepo := setupSEService()

	mockRepo.On("GetAll").Return([]dto.SEResponse{}, nil)

	result, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, result, 0)

	mockRepo.AssertExpectations(t)
}

func TestSEService_GetAll_RepositoryError(t *testing.T) {
	service, mockRepo := setupSEService()

	mockRepo.On("GetAll").Return(nil, errors.New("database error"))

	result, err := service.GetAll()

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

/*
=====================================
 TEST GET BY ID
=====================================
*/

func TestSEService_GetByID_Success(t *testing.T) {
	service, mockRepo := setupSEService()

	expectedData := &dto.SEResponse{
		ID:         "se-123",
		NamaSE:     "Sistem Informasi",
		KategoriSE: "Strategis",
		TotalBobot: 45,
	}

	mockRepo.On("GetByID", "se-123").Return(expectedData, nil)

	result, err := service.GetByID("se-123")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "se-123", result.ID)
	assert.Equal(t, "Strategis", result.KategoriSE)

	mockRepo.AssertExpectations(t)
}

func TestSEService_GetByID_EmptyID(t *testing.T) {
	service, _ := setupSEService()

	result, err := service.GetByID("")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "id wajib diisi", err.Error())
}

func TestSEService_GetByID_NotFound(t *testing.T) {
	service, mockRepo := setupSEService()

	mockRepo.On("GetByID", "invalid-id").Return(nil, errors.New("not found"))

	result, err := service.GetByID("invalid-id")

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

/*
=====================================
 TEST GET BY PERUSAHAAN
=====================================
*/

func TestSEService_GetByPerusahaan_Success(t *testing.T) {
	service, mockRepo := setupSEService()

	expectedData := []dto.SEResponse{
		{ID: "se-1", NamaSE: "SE 1", IDPerusahaan: "perusahaan-abc", KategoriSE: "Strategis", TotalBobot: 50},
		{ID: "se-2", NamaSE: "SE 2", IDPerusahaan: "perusahaan-abc", KategoriSE: "Tinggi", TotalBobot: 30},
	}

	mockRepo.On("GetByPerusahaan", "perusahaan-abc").Return(expectedData, nil)

	result, err := service.GetByPerusahaan("perusahaan-abc")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "perusahaan-abc", result[0].IDPerusahaan)

	mockRepo.AssertExpectations(t)
}

func TestSEService_GetByPerusahaan_EmptyID(t *testing.T) {
	service, _ := setupSEService()

	result, err := service.GetByPerusahaan("")

	assert.Error(t, err)
	assert.Equal(t, "id_perusahaan wajib diisi", err.Error())
	assert.Nil(t, result)
}

func TestSEService_GetByPerusahaan_RepositoryError(t *testing.T) {
	service, mockRepo := setupSEService()

	mockRepo.On("GetByPerusahaan", "perusahaan-abc").Return(nil, errors.New("database error"))

	result, err := service.GetByPerusahaan("perusahaan-abc")

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

/*
=====================================
 TEST UPDATE
=====================================
*/

func TestSEService_Update_Success_Recategorize(t *testing.T) {
	service, mockRepo := setupSEService()

	// Existing data: Strategis (bobot 50)
	existing := &dto.SEResponse{
		ID:                              "se-123",
		NamaSE:                          "Old Name",
		NilaiInvestasi:                  "A",
		AnggaranOperasional:             "A",
		KepatuhanPeraturan:              "A",
		TeknikKriptografi:               "A",
		JumlahPengguna:                  "A",
		DataPribadi:                     "A",
		KlasifikasiData:                 "A",
		KekritisanProses:                "A",
		DampakKegagalan:                 "A",
		PotensiKerugiandanDampakNegatif: "A",
		TotalBobot:                      50,
		KategoriSE:                      "Strategis",
	}

	// Update some fields to C → lower total bobot
	newNilaiInvestasi := "C"
	newAnggaranOperasional := "C"
	newKepatuhanPeraturan := "C"
	newTeknikKriptografi := "C"
	newJumlahPengguna := "C"

	req := dto.UpdateSERequest{
		NilaiInvestasi:      &newNilaiInvestasi,
		AnggaranOperasional: &newAnggaranOperasional,
		KepatuhanPeraturan:  &newKepatuhanPeraturan,
		TeknikKriptografi:   &newTeknikKriptografi,
		JumlahPengguna:      &newJumlahPengguna,
		// 5C (5) + 5A (25) = 30 → Tinggi
	}

	// First GetByID call returns existing data
	mockRepo.On("GetByID", "se-123").Return(existing, nil).Once()
	// Update is called
	mockRepo.On("Update", "se-123", req, 30, "Tinggi").Return(nil).Once()

	result, err := service.Update("se-123", req)

	assert.NoError(t, err)
	assert.Equal(t, 30, result.TotalBobot)
	assert.Equal(t, "Tinggi", result.KategoriSE)

	mockRepo.AssertExpectations(t)
}

func TestSEService_Update_EmptyID(t *testing.T) {
	service, _ := setupSEService()

	req := dto.UpdateSERequest{}

	result, err := service.Update("", req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "id wajib diisi", err.Error())
}

func TestSEService_Update_NotFound(t *testing.T) {
	service, mockRepo := setupSEService()

	req := dto.UpdateSERequest{}

	mockRepo.On("GetByID", "invalid-id").Return(nil, errors.New("not found"))

	result, err := service.Update("invalid-id", req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "data tidak ditemukan", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestSEService_Update_RepositoryError(t *testing.T) {
	service, mockRepo := setupSEService()

	existing := &dto.SEResponse{
		ID:                              "se-123",
		NilaiInvestasi:                  "A",
		AnggaranOperasional:             "A",
		KepatuhanPeraturan:              "A",
		TeknikKriptografi:               "A",
		JumlahPengguna:                  "A",
		DataPribadi:                     "A",
		KlasifikasiData:                 "A",
		KekritisanProses:                "A",
		DampakKegagalan:                 "A",
		PotensiKerugiandanDampakNegatif: "A",
	}

	req := dto.UpdateSERequest{}

	mockRepo.On("GetByID", "se-123").Return(existing, nil)
	mockRepo.On("Update", "se-123", req, 50, "Strategis").Return(errors.New("update failed"))

	result, err := service.Update("se-123", req)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

/*
=====================================
 TEST DELETE
=====================================
*/

func TestSEService_Delete_Success(t *testing.T) {
	service, mockRepo := setupSEService()

	existingSE := &dto.SEResponse{ID: "se-123", IDPerusahaan: "perusahaan-1"}
	mockRepo.On("GetByID", "se-123").Return(existingSE, nil)
	mockRepo.On("Delete", "se-123").Return(nil)

	err := service.Delete("se-123")

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestSEService_Delete_EmptyID(t *testing.T) {
	service, _ := setupSEService()

	err := service.Delete("")

	assert.Error(t, err)
	assert.Equal(t, "id wajib diisi", err.Error())
}

func TestSEService_Delete_RepositoryError(t *testing.T) {
	service, mockRepo := setupSEService()

	mockRepo.On("GetByID", "se-123").Return((*dto.SEResponse)(nil), nil)
	mockRepo.On("Delete", "se-123").Return(errors.New("delete failed"))

	err := service.Delete("se-123")

	assert.Error(t, err)
	assert.Equal(t, "delete failed", err.Error())

	mockRepo.AssertExpectations(t)
}
