package handlers

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
)

//
// MOCK SERVICE
//

type mockCsirtService struct {
	GetAllFn  func() ([]dto.CsirtResponse, error)
	GetByIDFn func(string) (*dto.CsirtResponse, error)
	CreateFn  func(dto.CreateCsirtRequest) (*models.Csirt, error)
	UpdateFn  func(string, dto.UpdateCsirtRequest) (*models.Csirt, error)
	DeleteFn  func(string) error
}

func (m *mockCsirtService) GetAll() ([]dto.CsirtResponse, error) {
	return m.GetAllFn()
}
func (m *mockCsirtService) GetByID(id string) (*dto.CsirtResponse, error) {
	return m.GetByIDFn(id)
}
func (m *mockCsirtService) Create(req dto.CreateCsirtRequest) (*models.Csirt, error) {
	return m.CreateFn(req)
}
func (m *mockCsirtService) Update(id string, req dto.UpdateCsirtRequest) (*models.Csirt, error) {
	return m.UpdateFn(id, req)
}
func (m *mockCsirtService) Delete(id string) error {
	return m.DeleteFn(id)
}

//
// HELPERS
//

func createMultipartRequest(
	t *testing.T,
	method, url string,
	fields map[string]string,
	files map[string]string,
) (*http.Request, func()) {

	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for k, v := range fields {
		_ = writer.WriteField(k, v)
	}

	for field, filename := range files {
		part, err := writer.CreateFormFile(field, filename)
		if err != nil {
			t.Fatalf("failed create form file: %v", err)
		}
		_, _ = io.WriteString(part, "dummy content")
	}

	writer.Close()

	req := httptest.NewRequest(method, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	cleanup := func() {
		_ = os.RemoveAll("uploads")
	}

	return req, cleanup
}

//
// TESTS — GET ALL
//

func TestCsirtHandler_GetAll_Success(t *testing.T) {
	mockSvc := &mockCsirtService{
		GetAllFn: func() ([]dto.CsirtResponse, error) {
			return []dto.CsirtResponse{
				{ID: "1", NamaCsirt: "CSIRT A"},
			}, nil
		},
	}

	handler := NewCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/csirt", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestCsirtHandler_GetAll_ServiceError(t *testing.T) {
	mockSvc := &mockCsirtService{
		GetAllFn: func() ([]dto.CsirtResponse, error) {
			return nil, errors.New("db error")
		},
	}

	handler := NewCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/csirt", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}

//
// TESTS — GET BY ID
//

func TestCsirtHandler_GetByID_Success(t *testing.T) {
	mockSvc := &mockCsirtService{
		GetByIDFn: func(id string) (*dto.CsirtResponse, error) {
			return &dto.CsirtResponse{
				ID:        id,
				NamaCsirt: "CSIRT A",
			}, nil
		},
	}

	handler := NewCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/csirt/1", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestCsirtHandler_GetByID_NotFound(t *testing.T) {
	mockSvc := &mockCsirtService{
		GetByIDFn: func(id string) (*dto.CsirtResponse, error) {
			return nil, errors.New("not found")
		},
	}

	handler := NewCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/csirt/123", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

//
// TESTS — CREATE
//

func TestCsirtHandler_Create_InvalidMultipart(t *testing.T) {
	mockSvc := &mockCsirtService{}

	handler := NewCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/api/csirt", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestCsirtHandler_Create_ServiceError(t *testing.T) {
	mockSvc := &mockCsirtService{
		CreateFn: func(req dto.CreateCsirtRequest) (*models.Csirt, error) {
			return nil, errors.New("create failed")
		},
	}

	handler := NewCsirtHandler(mockSvc)

	req, cleanup := createMultipartRequest(
		t,
		http.MethodPost,
		"/api/csirt",
		map[string]string{
			"id_perusahaan": "1",
			"nama_csirt":    "CSIRT",
		},
		nil,
	)
	t.Cleanup(cleanup)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

//
// TESTS — DELETE
//

func TestCsirtHandler_Delete_ServiceError(t *testing.T) {
	mockSvc := &mockCsirtService{
		DeleteFn: func(id string) error {
			return errors.New("delete failed")
		},
	}

	handler := NewCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/api/csirt/1", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

//
// TESTS — METHOD NOT ALLOWED
//

func TestCsirtHandler_MethodNotAllowed(t *testing.T) {
	mockSvc := &mockCsirtService{}

	handler := NewCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPatch, "/api/csirt", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}
