package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/*
=====================================
 MOCK SE SERVICE
=====================================
*/

type MockSEService struct {
	mock.Mock
}

func (m *MockSEService) Create(req dto.CreateSERequest) (*dto.SEResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SEResponse), args.Error(1)
}

func (m *MockSEService) GetAll() ([]dto.SEResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.SEResponse), args.Error(1)
}

func (m *MockSEService) GetByID(id string) (*dto.SEResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SEResponse), args.Error(1)
}

func (m *MockSEService) GetByPerusahaan(idPerusahaan string) ([]dto.SEResponse, error) {
	args := m.Called(idPerusahaan)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.SEResponse), args.Error(1)
}

func (m *MockSEService) Update(id string, req dto.UpdateSERequest) (*dto.SEResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SEResponse), args.Error(1)
}

func (m *MockSEService) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

/*
=====================================
 MOCK SSE SERVICE
=====================================
*/

type MockSSEServiceForSE struct {
	mock.Mock
}

func (m *MockSSEServiceForSE) NotifyCreate(resource string, data interface{}, userID string) {
	m.Called(resource, data, userID)
}

func (m *MockSSEServiceForSE) NotifyUpdate(resource string, data interface{}, userID string) {
	m.Called(resource, data, userID)
}

func (m *MockSSEServiceForSE) NotifyDelete(resource string, id interface{}, userID string) {
	m.Called(resource, id, userID)
}

/*
=====================================
 HELPER FUNCTIONS
=====================================
*/

func setupSEHandler() (*SEHandler, *MockSEService, *MockSSEServiceForSE) {
	mockService := new(MockSEService)
	mockSSE := new(MockSSEServiceForSE)
	handler := NewSEHandler(mockService, mockSSE)
	return handler, mockService, mockSSE
}

// withAdminContext set role admin di context request
func withAdminContext(req *http.Request) *http.Request {
	ctx := context.WithValue(req.Context(), middleware.RoleKey, "admin")
	return req.WithContext(ctx)
}

// withUserContext set role user + id_perusahaan di context request
func withUserContext(req *http.Request, idPerusahaan string) *http.Request {
	ctx := context.WithValue(req.Context(), middleware.RoleKey, "user")
	ctx = context.WithValue(ctx, middleware.IDPerusahaanKey, idPerusahaan)
	return req.WithContext(ctx)
}

func createValidSERequestForHandler() dto.CreateSERequest {
	return dto.CreateSERequest{
		IDPerusahaan:                    "perusahaan-123",
		NamaSE:                          "Sistem Informasi",
		IpSE:                            "192.168.1.1",
		AsNumberSE:                      "AS12345",
		PengelolaSE:                     "IT Dept",
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
 TEST GET ALL
=====================================
*/

func TestSEHandler_GetAll_Success(t *testing.T) {
	handler, mockService, _ := setupSEHandler()

	expectedData := []dto.SEResponse{
		{
			ID:         "se-1",
			NamaSE:     "SE 1",
			KategoriSE: "Strategis",
			TotalBobot: 50,
		},
		{
			ID:         "se-2",
			NamaSE:     "SE 2",
			KategoriSE: "Tinggi",
			TotalBobot: 30,
		},
	}

	mockService.On("GetAll").Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/se", nil)
	req = withAdminContext(req)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.SEResponse
	json.NewDecoder(w.Body).Decode(&response)
	assert.Len(t, response, 2)
	assert.Equal(t, "se-1", response[0].ID)
	assert.Equal(t, "Strategis", response[0].KategoriSE)

	mockService.AssertExpectations(t)
}

func TestSEHandler_GetAll_ServiceError(t *testing.T) {
	handler, mockService, _ := setupSEHandler()

	mockService.On("GetAll").Return(nil, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/api/se", nil)
	req = withAdminContext(req)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	assert.Contains(t, response["error"], "database error")

	mockService.AssertExpectations(t)
}

/*
=====================================
 TEST GET BY ID
=====================================
*/

func TestSEHandler_GetByID_Success(t *testing.T) {
	handler, mockService, _ := setupSEHandler()

	expectedData := &dto.SEResponse{
		ID:         "se-123",
		NamaSE:     "Sistem Informasi",
		KategoriSE: "Strategis",
		TotalBobot: 45,
	}

	mockService.On("GetByID", "se-123").Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/se/se-123", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.SEResponse
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, "se-123", response.ID)
	assert.Equal(t, "Strategis", response.KategoriSE)

	mockService.AssertExpectations(t)
}

func TestSEHandler_GetByID_NotFound(t *testing.T) {
	handler, mockService, _ := setupSEHandler()

	mockService.On("GetByID", "invalid-id").Return(nil, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/api/se/invalid-id", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	assert.Contains(t, response["error"], "tidak ditemukan")

	mockService.AssertExpectations(t)
}

/*
=====================================
 TEST CREATE
=====================================
*/

func TestSEHandler_Create_Success(t *testing.T) {
	handler, mockService, mockSSE := setupSEHandler()

	reqBody := createValidSERequestForHandler()

	expectedResponse := &dto.SEResponse{
		ID:           "new-se-id",
		IDPerusahaan: reqBody.IDPerusahaan,
		NamaSE:       reqBody.NamaSE,
		IpSE:         reqBody.IpSE,
		TotalBobot:   50,
		KategoriSE:   "Strategis",
	}

	mockService.On("Create", mock.AnythingOfType("dto.CreateSERequest")).Return(expectedResponse, nil)
	mockSSE.On("NotifyCreate", "se", expectedResponse, "")

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/se", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = withAdminContext(req)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response dto.SEResponse
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, reqBody.NamaSE, response.NamaSE)
	assert.Equal(t, 50, response.TotalBobot)
	assert.Equal(t, "Strategis", response.KategoriSE)

	mockService.AssertExpectations(t)
	mockSSE.AssertExpectations(t)
}

func TestSEHandler_Create_InvalidBody(t *testing.T) {
	handler, _, _ := setupSEHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/se", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	req = withAdminContext(req)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid request body")
}

func TestSEHandler_Create_ServiceError(t *testing.T) {
	handler, mockService, _ := setupSEHandler()

	reqBody := createValidSERequestForHandler()

	mockService.On("Create", mock.AnythingOfType("dto.CreateSERequest")).
		Return(nil, errors.New("validation error: nama_se wajib diisi"))

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/se", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = withAdminContext(req)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	assert.Contains(t, response["error"], "validation error")

	mockService.AssertExpectations(t)
}

func TestSEHandler_Create_KategoriStrategis(t *testing.T) {
	handler, mockService, mockSSE := setupSEHandler()

	reqBody := createValidSERequestForHandler()

	expectedResponse := &dto.SEResponse{
		ID:         "new-se-id",
		NamaSE:     reqBody.NamaSE,
		TotalBobot: 50,
		KategoriSE: "Strategis",
	}

	mockService.On("Create", mock.AnythingOfType("dto.CreateSERequest")).Return(expectedResponse, nil)
	mockSSE.On("NotifyCreate", "se", expectedResponse, "")

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/se", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = withAdminContext(req)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response dto.SEResponse
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, "Strategis", response.KategoriSE)

	mockService.AssertExpectations(t)
}

func TestSEHandler_Create_KategoriTinggi(t *testing.T) {
	handler, mockService, mockSSE := setupSEHandler()

	reqBody := createValidSERequestForHandler()
	// Adjust untuk kategori Tinggi (bobot 16-34)
	reqBody.DataPribadi = "B"
	reqBody.KlasifikasiData = "B"
	reqBody.KekritisanProses = "B"
	reqBody.DampakKegagalan = "C"
	reqBody.PotensiKerugiandanDampakNegatif = "C"

	expectedResponse := &dto.SEResponse{
		ID:         "new-se-id",
		TotalBobot: 33,
		KategoriSE: "Tinggi",
	}

	mockService.On("Create", mock.AnythingOfType("dto.CreateSERequest")).Return(expectedResponse, nil)
	mockSSE.On("NotifyCreate", "se", expectedResponse, "")

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/se", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = withAdminContext(req)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response dto.SEResponse
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, "Tinggi", response.KategoriSE)

	mockService.AssertExpectations(t)
}

/*
=====================================
 TEST UPDATE
=====================================
*/

func TestSEHandler_Update_Success(t *testing.T) {
	handler, mockService, mockSSE := setupSEHandler()

	newNama := "Updated SE Name"
	reqBody := dto.UpdateSERequest{
		NamaSE: &newNama,
	}

	expectedResponse := &dto.SEResponse{
		ID:         "se-123",
		NamaSE:     "Updated SE Name",
		TotalBobot: 50,
		KategoriSE: "Strategis",
	}

	mockService.On("Update", "se-123", mock.AnythingOfType("dto.UpdateSERequest")).
		Return(expectedResponse, nil)
	mockSSE.On("NotifyUpdate", "se", expectedResponse, "")

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/se/se-123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = withAdminContext(req)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.SEResponse
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, "Updated SE Name", response.NamaSE)

	mockService.AssertExpectations(t)
	mockSSE.AssertExpectations(t)
}

func TestSEHandler_Update_Recategorize(t *testing.T) {
	handler, mockService, mockSSE := setupSEHandler()

	// Update karakteristik untuk mengubah kategori
	newNilai := "C"
	reqBody := dto.UpdateSERequest{
		NilaiInvestasi:      &newNilai,
		AnggaranOperasional: &newNilai,
		KepatuhanPeraturan:  &newNilai,
	}

	expectedResponse := &dto.SEResponse{
		ID:         "se-123",
		TotalBobot: 28,
		KategoriSE: "Tinggi", // Re-categorized from Strategis to Tinggi
	}

	mockService.On("Update", "se-123", mock.AnythingOfType("dto.UpdateSERequest")).
		Return(expectedResponse, nil)
	mockSSE.On("NotifyUpdate", "se", expectedResponse, "")

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/se/se-123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = withAdminContext(req)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.SEResponse
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, "Tinggi", response.KategoriSE)
	assert.Equal(t, 28, response.TotalBobot)

	mockService.AssertExpectations(t)
}

func TestSEHandler_Update_InvalidBody(t *testing.T) {
	handler, _, _ := setupSEHandler()

	req := httptest.NewRequest(http.MethodPut, "/api/se/se-123", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	req = withAdminContext(req)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSEHandler_Update_ServiceError(t *testing.T) {
	handler, mockService, _ := setupSEHandler()

	reqBody := dto.UpdateSERequest{}

	mockService.On("Update", "se-123", mock.AnythingOfType("dto.UpdateSERequest")).
		Return(nil, errors.New("update failed"))

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/se/se-123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = withAdminContext(req)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockService.AssertExpectations(t)
}

/*
=====================================
 TEST DELETE
=====================================
*/

func TestSEHandler_Delete_Success(t *testing.T) {
	handler, mockService, mockSSE := setupSEHandler()

	mockService.On("Delete", "se-123").Return(nil)
	mockSSE.On("NotifyDelete", "se", "se-123", "")

	req := httptest.NewRequest(http.MethodDelete, "/api/se/se-123", nil)
	req = withAdminContext(req)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, "Delete success", response["message"])

	mockService.AssertExpectations(t)
	mockSSE.AssertExpectations(t)
}

func TestSEHandler_Delete_ServiceError(t *testing.T) {
	handler, mockService, _ := setupSEHandler()

	mockService.On("Delete", "se-123").Return(errors.New("delete failed"))

	req := httptest.NewRequest(http.MethodDelete, "/api/se/se-123", nil)
	req = withAdminContext(req)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockService.AssertExpectations(t)
}

/*
=====================================
 TEST METHOD NOT ALLOWED
=====================================
*/

func TestSEHandler_MethodNotAllowed(t *testing.T) {
	handler, _, _ := setupSEHandler()

	req := httptest.NewRequest(http.MethodPatch, "/api/se", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

/*
=====================================
 TEST SERVE HTTP ROUTING
=====================================
*/

func TestSEHandler_ServeHTTP_Routes(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		setupMock      func(*MockSEService, *MockSSEServiceForSE)
	}{
		{
			name:           "GET all",
			method:         http.MethodGet,
			path:           "/api/se",
			expectedStatus: http.StatusOK,
			setupMock: func(ms *MockSEService, msse *MockSSEServiceForSE) {
				ms.On("GetAll").Return([]dto.SEResponse{}, nil)
			},
		},
		{
			name:           "GET by ID",
			method:         http.MethodGet,
			path:           "/api/se/123",
			expectedStatus: http.StatusOK,
			setupMock: func(ms *MockSEService, msse *MockSSEServiceForSE) {
				ms.On("GetByID", "123").Return(&dto.SEResponse{ID: "123"}, nil)
			},
		},
		{
			name:           "POST create",
			method:         http.MethodPost,
			path:           "/api/se",
			expectedStatus: http.StatusBadRequest, // Invalid body
			setupMock:      func(ms *MockSEService, msse *MockSSEServiceForSE) {},
		},
		{
			name:           "PUT update",
			method:         http.MethodPut,
			path:           "/api/se/123",
			expectedStatus: http.StatusBadRequest, // Invalid body
			setupMock:      func(ms *MockSEService, msse *MockSSEServiceForSE) {},
		},
		{
			name:           "DELETE",
			method:         http.MethodDelete,
			path:           "/api/se/123",
			expectedStatus: http.StatusOK,
			setupMock: func(ms *MockSEService, msse *MockSSEServiceForSE) {
				ms.On("Delete", "123").Return(nil)
				msse.On("NotifyDelete", "se", "123", "")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockService, mockSSE := setupSEHandler()
			tt.setupMock(mockService, mockSSE)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			req = withAdminContext(req)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

/*
=====================================
 TEST OWNERSHIP — GET ALL AS USER
=====================================
*/

func TestSEHandler_GetAll_AsUser_FilterByPerusahaan(t *testing.T) {
	handler, mockService, _ := setupSEHandler()

	expectedData := []dto.SEResponse{
		{ID: "se-1", NamaSE: "SE Milik Perusahaan", IDPerusahaan: "perusahaan-abc"},
	}
	mockService.On("GetByPerusahaan", "perusahaan-abc").Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/se", nil)
	req = withUserContext(req, "perusahaan-abc")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []dto.SEResponse
	json.NewDecoder(w.Body).Decode(&response)
	assert.Len(t, response, 1)
	assert.Equal(t, "perusahaan-abc", response[0].IDPerusahaan)
	mockService.AssertExpectations(t)
}

func TestSEHandler_GetAll_AsUser_NoPerusahaan_Forbidden(t *testing.T) {
	handler, _, _ := setupSEHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/se", nil)
	ctx := context.WithValue(req.Context(), middleware.RoleKey, "user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestSEHandler_GetAll_AsUser_ServiceError(t *testing.T) {
	handler, mockService, _ := setupSEHandler()

	mockService.On("GetByPerusahaan", "perusahaan-abc").Return(nil, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/se", nil)
	req = withUserContext(req, "perusahaan-abc")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

/*
=====================================
 TEST OWNERSHIP — GET BY ID AS USER
=====================================
*/

func TestSEHandler_GetByID_AsUser_OwnData_Success(t *testing.T) {
	handler, mockService, _ := setupSEHandler()

	mockService.On("GetByID", "se-123").Return(&dto.SEResponse{
		ID: "se-123", IDPerusahaan: "perusahaan-abc",
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/se/se-123", nil)
	req = withUserContext(req, "perusahaan-abc")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestSEHandler_GetByID_AsUser_OtherPerusahaan_Forbidden(t *testing.T) {
	handler, mockService, _ := setupSEHandler()

	mockService.On("GetByID", "se-123").Return(&dto.SEResponse{
		ID: "se-123", IDPerusahaan: "perusahaan-lain",
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/se/se-123", nil)
	req = withUserContext(req, "perusahaan-abc")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	mockService.AssertExpectations(t)
}

/*
=====================================
 TEST OWNERSHIP — CREATE AS USER
=====================================
*/

func TestSEHandler_Create_AsUser_IDPerusahaanForcedFromJWT(t *testing.T) {
	handler, mockService, mockSSE := setupSEHandler()

	reqBody := createValidSERequestForHandler()
	reqBody.IDPerusahaan = "perusahaan-lain" // harus di-override oleh JWT

	expectedResponse := &dto.SEResponse{
		ID: "new-se", IDPerusahaan: "perusahaan-abc", NamaSE: reqBody.NamaSE, KategoriSE: "Strategis",
	}
	mockService.On("Create", mock.MatchedBy(func(r dto.CreateSERequest) bool {
		return r.IDPerusahaan == "perusahaan-abc"
	})).Return(expectedResponse, nil)
	mockSSE.On("NotifyCreate", "se", expectedResponse, "")

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/se", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = withUserContext(req, "perusahaan-abc")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestSEHandler_Create_AsUser_NoPerusahaan_Forbidden(t *testing.T) {
	handler, _, _ := setupSEHandler()

	body, _ := json.Marshal(createValidSERequestForHandler())
	req := httptest.NewRequest(http.MethodPost, "/api/se", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RoleKey, "user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

/*
=====================================
 TEST OWNERSHIP — UPDATE AS USER
=====================================
*/

func TestSEHandler_Update_AsUser_OwnData_Success(t *testing.T) {
	handler, mockService, mockSSE := setupSEHandler()

	newNama := "Updated"
	mockService.On("GetByID", "se-123").Return(&dto.SEResponse{
		ID: "se-123", IDPerusahaan: "perusahaan-abc",
	}, nil)
	mockService.On("Update", "se-123", mock.AnythingOfType("dto.UpdateSERequest")).Return(&dto.SEResponse{
		ID: "se-123", NamaSE: "Updated",
	}, nil)
	mockSSE.On("NotifyUpdate", "se", mock.Anything, "")

	body, _ := json.Marshal(dto.UpdateSERequest{NamaSE: &newNama})
	req := httptest.NewRequest(http.MethodPut, "/api/se/se-123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = withUserContext(req, "perusahaan-abc")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestSEHandler_Update_AsUser_OtherPerusahaan_Forbidden(t *testing.T) {
	handler, mockService, _ := setupSEHandler()

	mockService.On("GetByID", "se-123").Return(&dto.SEResponse{
		ID: "se-123", IDPerusahaan: "perusahaan-lain",
	}, nil)

	body, _ := json.Marshal(dto.UpdateSERequest{})
	req := httptest.NewRequest(http.MethodPut, "/api/se/se-123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = withUserContext(req, "perusahaan-abc")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	mockService.AssertExpectations(t)
}

func TestSEHandler_Update_AsUser_NotFound(t *testing.T) {
	handler, mockService, _ := setupSEHandler()

	mockService.On("GetByID", "se-nonexistent").Return(nil, errors.New("not found"))

	body, _ := json.Marshal(dto.UpdateSERequest{})
	req := httptest.NewRequest(http.MethodPut, "/api/se/se-nonexistent", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = withUserContext(req, "perusahaan-abc")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

/*
=====================================
 TEST OWNERSHIP — DELETE AS USER
=====================================
*/

func TestSEHandler_Delete_AsUser_OwnData_Success(t *testing.T) {
	handler, mockService, mockSSE := setupSEHandler()

	mockService.On("GetByID", "se-123").Return(&dto.SEResponse{
		ID: "se-123", IDPerusahaan: "perusahaan-abc",
	}, nil)
	mockService.On("Delete", "se-123").Return(nil)
	mockSSE.On("NotifyDelete", "se", "se-123", "")

	req := httptest.NewRequest(http.MethodDelete, "/api/se/se-123", nil)
	req = withUserContext(req, "perusahaan-abc")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestSEHandler_Delete_AsUser_OtherPerusahaan_Forbidden(t *testing.T) {
	handler, mockService, _ := setupSEHandler()

	mockService.On("GetByID", "se-123").Return(&dto.SEResponse{
		ID: "se-123", IDPerusahaan: "perusahaan-lain",
	}, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/se/se-123", nil)
	req = withUserContext(req, "perusahaan-abc")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	mockService.AssertExpectations(t)
}

func TestSEHandler_Delete_AsUser_NotFound(t *testing.T) {
	handler, mockService, _ := setupSEHandler()

	mockService.On("GetByID", "se-nonexistent").Return(nil, errors.New("not found"))

	req := httptest.NewRequest(http.MethodDelete, "/api/se/se-nonexistent", nil)
	req = withUserContext(req, "perusahaan-abc")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}
