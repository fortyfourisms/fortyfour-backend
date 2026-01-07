package handlers_test

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/handlers"

	"github.com/stretchr/testify/assert"
)

/* =========================
   MOCK SERVICE
========================= */

type mockPerusahaanService struct {
	getAllFn  func() ([]dto.PerusahaanResponse, error)
	getByIDFn func(string) (*dto.PerusahaanResponse, error)
	createFn  func(dto.CreatePerusahaanRequest) (*dto.PerusahaanResponse, error)
	updateFn  func(string, dto.UpdatePerusahaanRequest) (*dto.PerusahaanResponse, error)
	deleteFn  func(string) error
}

func (m *mockPerusahaanService) GetAll() ([]dto.PerusahaanResponse, error) {
	return m.getAllFn()
}

func (m *mockPerusahaanService) GetByID(id string) (*dto.PerusahaanResponse, error) {
	return m.getByIDFn(id)
}

func (m *mockPerusahaanService) Create(req dto.CreatePerusahaanRequest) (*dto.PerusahaanResponse, error) {
	return m.createFn(req)
}

func (m *mockPerusahaanService) Update(id string, req dto.UpdatePerusahaanRequest) (*dto.PerusahaanResponse, error) {
	return m.updateFn(id, req)
}

func (m *mockPerusahaanService) Delete(id string) error {
	return m.deleteFn(id)
}

/* =========================
   MOCK SSE
========================= */

type mockSSEService struct {
	createCalled bool
	updateCalled bool
	deleteCalled bool
}

func (m *mockSSEService) NotifyCreate(string, interface{}, string) {
	m.createCalled = true
}

func (m *mockSSEService) NotifyUpdate(string, interface{}, string) {
	m.updateCalled = true
}

func (m *mockSSEService) NotifyDelete(string, interface{}, string) {
	m.deleteCalled = true
}

/* =========================
   LOCAL MULTIPART HELPER
========================= */

func createMultipartRequest(method, url string, fields map[string]string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for k, v := range fields {
		_ = writer.WriteField(k, v)
	}

	writer.Close()

	req := httptest.NewRequest(method, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

/* =========================
   TESTS
========================= */

func TestPerusahaanHandler_GetAll_Success(t *testing.T) {
	mockSvc := &mockPerusahaanService{
		getAllFn: func() ([]dto.PerusahaanResponse, error) {
			return []dto.PerusahaanResponse{{ID: "1"}}, nil
		},
	}
	mockSSE := &mockSSEService{}

	handler := handlers.NewPerusahaanHandler(mockSvc, "/tmp", mockSSE)

	req := httptest.NewRequest(http.MethodGet, "/api/perusahaan", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestPerusahaanHandler_GetByID_NotFound(t *testing.T) {
	mockSvc := &mockPerusahaanService{
		getByIDFn: func(string) (*dto.PerusahaanResponse, error) {
			return nil, errors.New("not found")
		},
	}
	handler := handlers.NewPerusahaanHandler(mockSvc, "/tmp", &mockSSEService{})

	req := httptest.NewRequest(http.MethodGet, "/api/perusahaan/1", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestPerusahaanHandler_Create_Success(t *testing.T) {
	mockSvc := &mockPerusahaanService{
		createFn: func(dto.CreatePerusahaanRequest) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{ID: "1"}, nil
		},
	}
	mockSSE := &mockSSEService{}

	handler := handlers.NewPerusahaanHandler(mockSvc, "/tmp", mockSSE)

	req, _ := createMultipartRequest(
		http.MethodPost,
		"/api/perusahaan",
		map[string]string{
			"nama_perusahaan": "PT Test",
		},
	)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.True(t, mockSSE.createCalled)
}

func TestPerusahaanHandler_Update_Success(t *testing.T) {
	mockSvc := &mockPerusahaanService{
		updateFn: func(string, dto.UpdatePerusahaanRequest) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{ID: "1"}, nil
		},
		getByIDFn: func(string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{ID: "1"}, nil
		},
	}
	mockSSE := &mockSSEService{}
	handler := handlers.NewPerusahaanHandler(mockSvc, "/tmp", mockSSE)

	req, _ := createMultipartRequest(
		http.MethodPut,
		"/api/perusahaan/1",
		map[string]string{
			"nama_perusahaan": "Updated",
		},
	)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.True(t, mockSSE.updateCalled)
}

func TestPerusahaanHandler_Delete_Success(t *testing.T) {
	mockSvc := &mockPerusahaanService{
		deleteFn: func(string) error { return nil },
		getByIDFn: func(string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{ID: "1", Photo: ""}, nil
		},
	}
	mockSSE := &mockSSEService{}
	handler := handlers.NewPerusahaanHandler(mockSvc, "/tmp", mockSSE)

	req := httptest.NewRequest(http.MethodDelete, "/api/perusahaan/1", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.True(t, mockSSE.deleteCalled)
}
