package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/handlers"
	"fortyfour-backend/internal/middleware"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/* =========================
   MOCK SERVICE
========================= */

type mockPerusahaanService struct {
	mock.Mock
}

func (m *mockPerusahaanService) GetAll() ([]dto.PerusahaanResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.PerusahaanResponse), args.Error(1)
}

func (m *mockPerusahaanService) GetByID(id string) (*dto.PerusahaanResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PerusahaanResponse), args.Error(1)
}

func (m *mockPerusahaanService) Create(req dto.CreatePerusahaanRequest) (*dto.PerusahaanResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PerusahaanResponse), args.Error(1)
}

func (m *mockPerusahaanService) Update(id string, req dto.UpdatePerusahaanRequest) (*dto.PerusahaanResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PerusahaanResponse), args.Error(1)
}

func (m *mockPerusahaanService) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

/* =========================
   MOCK SSE SERVICE
========================= */

type mockSSEService struct {
	mock.Mock
}

func (m *mockSSEService) NotifyCreate(resource string, data interface{}, userID string) {
	m.Called(resource, data, userID)
}

func (m *mockSSEService) NotifyUpdate(resource string, data interface{}, userID string) {
	m.Called(resource, data, userID)
}

func (m *mockSSEService) NotifyDelete(resource string, id interface{}, userID string) {
	m.Called(resource, id, userID)
}

/* =========================
   HELPER FUNCTIONS
========================= */

func createMultipartRequest(method, url string, fields map[string]string, files map[string][]byte) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add text fields
	for k, v := range fields {
		_ = writer.WriteField(k, v)
	}

	// Add file fields with proper Content-Type
	for fieldName, fileData := range files {
		// Create form file with proper MIME header
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="test.jpg"`, fieldName))
		h.Set("Content-Type", "image/jpeg")

		part, err := writer.CreatePart(h)
		if err != nil {
			return nil, err
		}

		_, err = part.Write(fileData)
		if err != nil {
			return nil, err
		}
	}

	writer.Close()

	req := httptest.NewRequest(method, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func createTestImage(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, img, &jpeg.Options{Quality: 80})
	return buf.Bytes()
}

func createTestImagePNG(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	buf := new(bytes.Buffer)
	png.Encode(buf, img)
	return buf.Bytes()
}

func setupHandler(t *testing.T) (*handlers.PerusahaanHandler, *mockPerusahaanService, *mockSSEService, string) {
	mockSvc := new(mockPerusahaanService)
	mockSSE := new(mockSSEService)
	tmpDir := t.TempDir()
	handler := handlers.NewPerusahaanHandler(mockSvc, tmpDir, mockSSE)
	return handler, mockSvc, mockSSE, tmpDir
}

/* =========================
   TEST GET ALL
========================= */

func TestPerusahaanHandler_GetAll_Success(t *testing.T) {
	handler, mockSvc, mockSSE, _ := setupHandler(t)

	expectedData := []dto.PerusahaanResponse{
		{
			ID:             "1",
			NamaPerusahaan: "PT Test 1",
			Alamat:         "Jl. Test 1",
		},
		{
			ID:             "2",
			NamaPerusahaan: "PT Test 2",
			Alamat:         "Jl. Test 2",
		},
	}

	mockSvc.On("GetAll").Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/perusahaan", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []dto.PerusahaanResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "PT Test 1", response[0].NamaPerusahaan)
	assert.Equal(t, "PT Test 2", response[1].NamaPerusahaan)

	mockSvc.AssertExpectations(t)
	mockSSE.AssertNotCalled(t, "NotifyCreate")
	mockSSE.AssertNotCalled(t, "NotifyUpdate")
	mockSSE.AssertNotCalled(t, "NotifyDelete")
}

func TestPerusahaanHandler_GetAll_EmptyResult(t *testing.T) {
	handler, mockSvc, _, _ := setupHandler(t)

	mockSvc.On("GetAll").Return([]dto.PerusahaanResponse{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/perusahaan", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []dto.PerusahaanResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 0)

	mockSvc.AssertExpectations(t)
}

func TestPerusahaanHandler_GetAll_ServiceError(t *testing.T) {
	handler, mockSvc, _, _ := setupHandler(t)

	mockSvc.On("GetAll").Return(nil, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/api/perusahaan", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "database error")

	mockSvc.AssertExpectations(t)
}

/* =========================
   TEST GET BY ID
========================= */

func TestPerusahaanHandler_GetByID_Success(t *testing.T) {
	handler, mockSvc, _, _ := setupHandler(t)

	expectedData := &dto.PerusahaanResponse{
		ID:             "test-id-123",
		NamaPerusahaan: "PT Test",
		Alamat:         "Jl. Test No. 123",
		Email:          "test@test.com",
		Telepon:        "08123456789",
	}

	mockSvc.On("GetByID", "test-id-123").Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/perusahaan/test-id-123", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response dto.PerusahaanResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "test-id-123", response.ID)
	assert.Equal(t, "PT Test", response.NamaPerusahaan)
	assert.Equal(t, "Jl. Test No. 123", response.Alamat)

	mockSvc.AssertExpectations(t)
}

func TestPerusahaanHandler_GetByID_NotFound(t *testing.T) {
	handler, mockSvc, _, _ := setupHandler(t)

	mockSvc.On("GetByID", "nonexistent-id").Return(nil, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/api/perusahaan/nonexistent-id", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Data tidak ditemukan")

	mockSvc.AssertExpectations(t)
}

func TestPerusahaanHandler_GetByID_WithSpecialCharacters(t *testing.T) {
	handler, mockSvc, _, _ := setupHandler(t)

	specialIDs := []string{
		"id-with-dashes",
		"id_with_underscores",
		"123-456-789",
		"uuid-1234-5678-90ab-cdef",
	}

	for _, id := range specialIDs {
		t.Run(id, func(t *testing.T) {
			expectedData := &dto.PerusahaanResponse{
				ID:             id,
				NamaPerusahaan: "PT Test",
			}

			mockSvc.On("GetByID", id).Return(expectedData, nil).Once()

			req := httptest.NewRequest(http.MethodGet, "/api/perusahaan/"+id, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
			mockSvc.AssertExpectations(t)
		})
	}
}

/* =========================
   TEST CREATE
========================= */

func TestPerusahaanHandler_Create_Success_MinimalData(t *testing.T) {
	handler, mockSvc, mockSSE, _ := setupHandler(t)

	expectedResp := &dto.PerusahaanResponse{
		ID:             "new-id",
		NamaPerusahaan: "PT Baru",
	}

	mockSvc.On("Create", mock.AnythingOfType("dto.CreatePerusahaanRequest")).
		Return(expectedResp, nil)
	mockSSE.On("NotifyCreate", "perusahaan", expectedResp, "").Return()

	req, _ := createMultipartRequest(
		http.MethodPost,
		"/api/perusahaan",
		map[string]string{
			"nama_perusahaan": "PT Baru",
		},
		nil,
	)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response dto.PerusahaanResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "new-id", response.ID)
	assert.Equal(t, "PT Baru", response.NamaPerusahaan)

	mockSvc.AssertExpectations(t)
	mockSSE.AssertExpectations(t)
}

func TestPerusahaanHandler_Create_Success_CompleteData(t *testing.T) {
	handler, mockSvc, mockSSE, _ := setupHandler(t)

	expectedResp := &dto.PerusahaanResponse{
		ID:             "new-id",
		NamaPerusahaan: "PT Lengkap",
		Alamat:         "Jl. Complete No. 1",
		Email:          "complete@test.com",
		Telepon:        "08123456789",
	}

	mockSvc.On("Create", mock.AnythingOfType("dto.CreatePerusahaanRequest")).
		Return(expectedResp, nil)
	mockSSE.On("NotifyCreate", "perusahaan", expectedResp, "test-user-id").Return()

	req, _ := createMultipartRequest(
		http.MethodPost,
		"/api/perusahaan",
		map[string]string{
			"nama_perusahaan": "PT Lengkap",
			"alamat":          "Jl. Complete No. 1",
			"email":           "complete@test.com",
			"telepon":         "08123456789",
		},
		nil,
	)

	// Add user context
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "test-user-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	mockSvc.AssertExpectations(t)
	mockSSE.AssertExpectations(t)
}

func TestPerusahaanHandler_Create_Success_WithPhoto(t *testing.T) {
	handler, mockSvc, mockSSE, tmpDir := setupHandler(t)

	testImage := createTestImage(800, 600)

	expectedResp := &dto.PerusahaanResponse{
		ID:             "new-id",
		NamaPerusahaan: "PT Photo",
		Photo:          "photo.jpg",
	}

	mockSvc.On("Create", mock.MatchedBy(func(req dto.CreatePerusahaanRequest) bool {
		// Verify that photo field is set (non-nil and non-empty)
		return req.Photo != nil && *req.Photo != ""
	})).Return(expectedResp, nil)
	mockSSE.On("NotifyCreate", "perusahaan", expectedResp, "").Return()

	req, _ := createMultipartRequest(
		http.MethodPost,
		"/api/perusahaan",
		map[string]string{
			"nama_perusahaan": "PT Photo",
		},
		map[string][]byte{
			"photo": testImage,
		},
	)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	mockSvc.AssertExpectations(t)
	mockSSE.AssertExpectations(t)

	// Verify uploaded file exists
	files, _ := os.ReadDir(tmpDir)
	assert.Greater(t, len(files), 0, "File should be uploaded")
}

func TestPerusahaanHandler_Create_WithIDInURL_ShouldFail(t *testing.T) {
	handler, mockSvc, mockSSE, _ := setupHandler(t)

	req, _ := createMultipartRequest(
		http.MethodPost,
		"/api/perusahaan/some-id",
		map[string]string{
			"nama_perusahaan": "PT Test",
		},
		nil,
	)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "ID tidak diperlukan")

	mockSvc.AssertNotCalled(t, "Create")
	mockSSE.AssertNotCalled(t, "NotifyCreate")
}

func TestPerusahaanHandler_Create_InvalidFormData(t *testing.T) {
	handler, mockSvc, mockSSE, _ := setupHandler(t)

	// Create a request with invalid content type
	req := httptest.NewRequest(http.MethodPost, "/api/perusahaan", strings.NewReader("invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Gagal membaca form data")

	mockSvc.AssertNotCalled(t, "Create")
	mockSSE.AssertNotCalled(t, "NotifyCreate")
}

func TestPerusahaanHandler_Create_ServiceError(t *testing.T) {
	handler, mockSvc, mockSSE, _ := setupHandler(t)

	mockSvc.On("Create", mock.AnythingOfType("dto.CreatePerusahaanRequest")).
		Return(nil, errors.New("validation failed: nama_perusahaan is required"))

	req, _ := createMultipartRequest(
		http.MethodPost,
		"/api/perusahaan",
		map[string]string{
			"nama_perusahaan": "",
		},
		nil,
	)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "validation failed")

	mockSvc.AssertExpectations(t)
	mockSSE.AssertNotCalled(t, "NotifyCreate")
}

func TestPerusahaanHandler_Create_PhotoTooLarge(t *testing.T) {
	handler, mockSvc, mockSSE, _ := setupHandler(t)

	// Create a file larger than 10MB
	largeFile := make([]byte, 11<<20) // 11 MB

	req, _ := createMultipartRequest(
		http.MethodPost,
		"/api/perusahaan",
		map[string]string{
			"nama_perusahaan": "PT Large",
		},
		map[string][]byte{
			"photo": largeFile,
		},
	)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	mockSvc.AssertNotCalled(t, "Create")
	mockSSE.AssertNotCalled(t, "NotifyCreate")
}

func TestPerusahaanHandler_Create_InvalidPhotoFormat(t *testing.T) {
	handler, mockSvc, mockSSE, _ := setupHandler(t)

	invalidImage := []byte("not an image")

	// Create multipart with non-image content type
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("nama_perusahaan", "PT Test")

	// Create file part with text/plain content type (not image/*)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="photo"; filename="test.txt"`)
	h.Set("Content-Type", "text/plain")

	part, _ := writer.CreatePart(h)
	part.Write(invalidImage)

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/perusahaan", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Should fail due to invalid content type
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	mockSvc.AssertNotCalled(t, "Create")
	mockSSE.AssertNotCalled(t, "NotifyCreate")
}

/* =========================
   TEST UPDATE
========================= */

func TestPerusahaanHandler_Update_Success_MinimalData(t *testing.T) {
	handler, mockSvc, mockSSE, _ := setupHandler(t)

	expectedResp := &dto.PerusahaanResponse{
		ID:             "update-id",
		NamaPerusahaan: "PT Updated",
	}

	mockSvc.On("Update", "update-id", mock.AnythingOfType("dto.UpdatePerusahaanRequest")).
		Return(expectedResp, nil)
	mockSSE.On("NotifyUpdate", "perusahaan", expectedResp, "").Return()

	req, _ := createMultipartRequest(
		http.MethodPut,
		"/api/perusahaan/update-id",
		map[string]string{
			"nama_perusahaan": "PT Updated",
		},
		nil,
	)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response dto.PerusahaanResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "PT Updated", response.NamaPerusahaan)

	mockSvc.AssertExpectations(t)
	mockSSE.AssertExpectations(t)
}

func TestPerusahaanHandler_Update_Success_WithUserContext(t *testing.T) {
	handler, mockSvc, mockSSE, _ := setupHandler(t)

	expectedResp := &dto.PerusahaanResponse{
		ID:             "update-id",
		NamaPerusahaan: "PT Updated",
	}

	mockSvc.On("Update", "update-id", mock.AnythingOfType("dto.UpdatePerusahaanRequest")).
		Return(expectedResp, nil)
	mockSSE.On("NotifyUpdate", "perusahaan", expectedResp, "user-123").Return()

	req, _ := createMultipartRequest(
		http.MethodPut,
		"/api/perusahaan/update-id",
		map[string]string{
			"nama_perusahaan": "PT Updated",
		},
		nil,
	)

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-123")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	mockSvc.AssertExpectations(t)
	mockSSE.AssertExpectations(t)
}

func TestPerusahaanHandler_Update_WithNewPhoto(t *testing.T) {
	handler, mockSvc, mockSSE, tmpDir := setupHandler(t)

	testImage := createTestImage(800, 600)

	// Create old photo file to simulate existing photo
	oldPhotoPath := filepath.Join(tmpDir, "old-photo.jpg")
	os.WriteFile(oldPhotoPath, testImage, 0644)

	existingData := &dto.PerusahaanResponse{
		ID:             "update-id",
		NamaPerusahaan: "PT Test",
		Photo:          "old-photo.jpg",
	}

	mockSvc.On("GetByID", "update-id").Return(existingData, nil)

	expectedResp := &dto.PerusahaanResponse{
		ID:             "update-id",
		NamaPerusahaan: "PT Test",
		Photo:          "new-photo.jpg",
	}

	mockSvc.On("Update", "update-id", mock.MatchedBy(func(req dto.UpdatePerusahaanRequest) bool {
		return req.Photo != nil && *req.Photo != ""
	})).Return(expectedResp, nil)
	mockSSE.On("NotifyUpdate", "perusahaan", expectedResp, "").Return()

	req, _ := createMultipartRequest(
		http.MethodPut,
		"/api/perusahaan/update-id",
		map[string]string{
			"nama_perusahaan": "PT Test",
		},
		map[string][]byte{
			"photo": testImage,
		},
	)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	mockSvc.AssertExpectations(t)
	mockSSE.AssertExpectations(t)
}

func TestPerusahaanHandler_Update_WithoutID_ShouldFail(t *testing.T) {
	handler, mockSvc, mockSSE, _ := setupHandler(t)

	req, _ := createMultipartRequest(
		http.MethodPut,
		"/api/perusahaan",
		map[string]string{
			"nama_perusahaan": "PT Test",
		},
		nil,
	)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "ID wajib")

	mockSvc.AssertNotCalled(t, "Update")
	mockSSE.AssertNotCalled(t, "NotifyUpdate")
}

func TestPerusahaanHandler_Update_InvalidFormData(t *testing.T) {
	handler, mockSvc, mockSSE, _ := setupHandler(t)

	req := httptest.NewRequest(http.MethodPut, "/api/perusahaan/test-id", strings.NewReader("invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	mockSvc.AssertNotCalled(t, "Update")
	mockSSE.AssertNotCalled(t, "NotifyUpdate")
}

func TestPerusahaanHandler_Update_ServiceError(t *testing.T) {
	handler, mockSvc, mockSSE, _ := setupHandler(t)

	mockSvc.On("Update", "test-id", mock.AnythingOfType("dto.UpdatePerusahaanRequest")).
		Return(nil, errors.New("perusahaan not found"))

	req, _ := createMultipartRequest(
		http.MethodPut,
		"/api/perusahaan/test-id",
		map[string]string{
			"nama_perusahaan": "PT Test",
		},
		nil,
	)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "not found")

	mockSvc.AssertExpectations(t)
	mockSSE.AssertNotCalled(t, "NotifyUpdate")
}

/* =========================
   TEST DELETE
========================= */

func TestPerusahaanHandler_Delete_Success(t *testing.T) {
	handler, mockSvc, mockSSE, _ := setupHandler(t)

	mockSvc.On("GetByID", "delete-id").Return(&dto.PerusahaanResponse{
		ID:    "delete-id",
		Photo: "",
	}, nil)
	mockSvc.On("Delete", "delete-id").Return(nil)
	mockSSE.On("NotifyDelete", "perusahaan", "delete-id", "").Return()

	req := httptest.NewRequest(http.MethodDelete, "/api/perusahaan/delete-id", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Delete success", response["message"])

	mockSvc.AssertExpectations(t)
	mockSSE.AssertExpectations(t)
}

func TestPerusahaanHandler_Delete_Success_WithUserContext(t *testing.T) {
	handler, mockSvc, mockSSE, _ := setupHandler(t)

	mockSvc.On("GetByID", "delete-id").Return(&dto.PerusahaanResponse{
		ID:    "delete-id",
		Photo: "",
	}, nil)
	mockSvc.On("Delete", "delete-id").Return(nil)
	mockSSE.On("NotifyDelete", "perusahaan", "delete-id", "user-456").Return()

	req := httptest.NewRequest(http.MethodDelete, "/api/perusahaan/delete-id", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-456")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	mockSvc.AssertExpectations(t)
	mockSSE.AssertExpectations(t)
}

func TestPerusahaanHandler_Delete_WithPhoto(t *testing.T) {
	handler, mockSvc, mockSSE, tmpDir := setupHandler(t)

	// Create a test photo file
	photoPath := filepath.Join(tmpDir, "test-photo.jpg")
	testImage := createTestImage(100, 100)
	os.WriteFile(photoPath, testImage, 0644)

	mockSvc.On("GetByID", "delete-id").Return(&dto.PerusahaanResponse{
		ID:    "delete-id",
		Photo: "test-photo.jpg",
	}, nil)
	mockSvc.On("Delete", "delete-id").Return(nil)
	mockSSE.On("NotifyDelete", "perusahaan", "delete-id", "").Return()

	req := httptest.NewRequest(http.MethodDelete, "/api/perusahaan/delete-id", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	// Verify photo file is deleted
	_, err := os.Stat(photoPath)
	assert.True(t, os.IsNotExist(err), "Photo file should be deleted")

	mockSvc.AssertExpectations(t)
	mockSSE.AssertExpectations(t)
}

func TestPerusahaanHandler_Delete_WithoutID_ShouldFail(t *testing.T) {
	handler, mockSvc, mockSSE, _ := setupHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/perusahaan", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "ID wajib")

	mockSvc.AssertNotCalled(t, "Delete")
	mockSSE.AssertNotCalled(t, "NotifyDelete")
}

func TestPerusahaanHandler_Delete_ServiceError(t *testing.T) {
	handler, mockSvc, mockSSE, _ := setupHandler(t)

	mockSvc.On("GetByID", "test-id").Return(&dto.PerusahaanResponse{
		ID:    "test-id",
		Photo: "",
	}, nil)
	mockSvc.On("Delete", "test-id").Return(errors.New("foreign key constraint"))

	req := httptest.NewRequest(http.MethodDelete, "/api/perusahaan/test-id", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "foreign key")

	mockSvc.AssertExpectations(t)
	mockSSE.AssertNotCalled(t, "NotifyDelete")
}

/* =========================
   TEST METHOD NOT ALLOWED
========================= */

func TestPerusahaanHandler_MethodNotAllowed(t *testing.T) {
	handler, mockSvc, mockSSE, _ := setupHandler(t)

	methods := []string{http.MethodPatch, http.MethodOptions, http.MethodHead}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/perusahaan", nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

			mockSvc.AssertNotCalled(t, "GetAll")
			mockSvc.AssertNotCalled(t, "GetByID")
			mockSvc.AssertNotCalled(t, "Create")
			mockSvc.AssertNotCalled(t, "Update")
			mockSvc.AssertNotCalled(t, "Delete")
			mockSSE.AssertNotCalled(t, "NotifyCreate")
			mockSSE.AssertNotCalled(t, "NotifyUpdate")
			mockSSE.AssertNotCalled(t, "NotifyDelete")
		})
	}
}

/* =========================
   TEST ROUTING
========================= */

func TestPerusahaanHandler_Routing(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		setupMock      func(*mockPerusahaanService, *mockSSEService)
		expectedStatus int
	}{
		{
			name:   "GET all",
			method: http.MethodGet,
			path:   "/api/perusahaan",
			setupMock: func(ms *mockPerusahaanService, msse *mockSSEService) {
				ms.On("GetAll").Return([]dto.PerusahaanResponse{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "GET by ID",
			method: http.MethodGet,
			path:   "/api/perusahaan/123",
			setupMock: func(ms *mockPerusahaanService, msse *mockSSEService) {
				ms.On("GetByID", "123").Return(&dto.PerusahaanResponse{ID: "123"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST without ID",
			method:         http.MethodPost,
			path:           "/api/perusahaan",
			setupMock:      func(ms *mockPerusahaanService, msse *mockSSEService) {},
			expectedStatus: http.StatusBadRequest, // Will fail due to form parsing
		},
		{
			name:           "POST with ID - should fail",
			method:         http.MethodPost,
			path:           "/api/perusahaan/123",
			setupMock:      func(ms *mockPerusahaanService, msse *mockSSEService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "PUT with ID",
			method:         http.MethodPut,
			path:           "/api/perusahaan/123",
			setupMock:      func(ms *mockPerusahaanService, msse *mockSSEService) {},
			expectedStatus: http.StatusBadRequest, // Will fail due to form parsing
		},
		{
			name:           "PUT without ID - should fail",
			method:         http.MethodPut,
			path:           "/api/perusahaan",
			setupMock:      func(ms *mockPerusahaanService, msse *mockSSEService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "DELETE with ID",
			method: http.MethodDelete,
			path:   "/api/perusahaan/123",
			setupMock: func(ms *mockPerusahaanService, msse *mockSSEService) {
				ms.On("GetByID", "123").Return(&dto.PerusahaanResponse{ID: "123"}, nil)
				ms.On("Delete", "123").Return(nil)
				msse.On("NotifyDelete", "perusahaan", "123", "").Return()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "DELETE without ID - should fail",
			method:         http.MethodDelete,
			path:           "/api/perusahaan",
			setupMock:      func(ms *mockPerusahaanService, msse *mockSSEService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockSvc, mockSSE, _ := setupHandler(t)
			tt.setupMock(mockSvc, mockSSE)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockSvc.AssertExpectations(t)
			mockSSE.AssertExpectations(t)
		})
	}
}
