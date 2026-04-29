package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/handlers"
	"fortyfour-backend/internal/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockBeritaService
type MockBeritaService struct {
	CreateFunc  func(authorID string, req dto.CreateBeritaRequest) error
	GetAllFunc  func() ([]dto.BeritaResponse, error)
	GetByIDFunc func(id int64) (*dto.BeritaResponse, error)
	UpdateFunc  func(id int64, req dto.UpdateBeritaRequest) error
	DeleteFunc  func(id int64) error
}

func (m *MockBeritaService) Create(authorID string, req dto.CreateBeritaRequest) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(authorID, req)
	}
	return nil
}

func (m *MockBeritaService) GetAll() ([]dto.BeritaResponse, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc()
	}
	return nil, nil
}

func (m *MockBeritaService) GetByID(id int64) (*dto.BeritaResponse, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *MockBeritaService) Update(id int64, req dto.UpdateBeritaRequest) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(id, req)
	}
	return nil
}

func (m *MockBeritaService) Delete(id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func TestBeritaHandler_ServeHTTP(t *testing.T) {}

func TestBeritaHandler_GetAll(t *testing.T) {
	mockSvc := &MockBeritaService{
		GetAllFunc: func() ([]dto.BeritaResponse, error) {
			return []dto.BeritaResponse{{ID: 1, Judul: "Test Berita"}}, nil
		},
	}
	handler := handlers.NewBeritaHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/berita", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", rec.Code)
	}
}

func TestBeritaHandler_GetAll_Error(t *testing.T) {
	mockSvc := &MockBeritaService{
		GetAllFunc: func() ([]dto.BeritaResponse, error) {
			return nil, errors.New("db error")
		},
	}
	handler := handlers.NewBeritaHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/berita", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status InternalServerError, got %v", rec.Code)
	}
}

func TestBeritaHandler_GetByID(t *testing.T) {
	mockSvc := &MockBeritaService{
		GetByIDFunc: func(id int64) (*dto.BeritaResponse, error) {
			return &dto.BeritaResponse{ID: 1, Judul: "Test Berita"}, nil
		},
	}
	handler := handlers.NewBeritaHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/berita/1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", rec.Code)
	}
}

func TestBeritaHandler_GetByID_InvalidID(t *testing.T) {
	handler := handlers.NewBeritaHandler(&MockBeritaService{})

	req := httptest.NewRequest(http.MethodGet, "/api/berita/abc", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest, got %v", rec.Code)
	}
}

func TestBeritaHandler_GetByID_Error(t *testing.T) {
	mockSvc := &MockBeritaService{
		GetByIDFunc: func(id int64) (*dto.BeritaResponse, error) {
			return nil, errors.New("not found")
		},
	}
	handler := handlers.NewBeritaHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/berita/1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status NotFound, got %v", rec.Code)
	}
}

func TestBeritaHandler_Create(t *testing.T) {
	mockSvc := &MockBeritaService{
		CreateFunc: func(authorID string, req dto.CreateBeritaRequest) error {
			return nil
		},
	}
	handler := handlers.NewBeritaHandler(mockSvc)

	body := dto.CreateBeritaRequest{Judul: "Test Berita", Deskripsi: "Deskripsi Berita"}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/berita", bytes.NewReader(bodyBytes))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "author-123")
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status Created, got %v", rec.Code)
	}
}

func TestBeritaHandler_Create_InvalidBody(t *testing.T) {
	handler := handlers.NewBeritaHandler(&MockBeritaService{})

	req := httptest.NewRequest(http.MethodPost, "/api/berita", bytes.NewReader([]byte("{invalid-json}")))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest, got %v", rec.Code)
	}
}

func TestBeritaHandler_Create_ValidationError(t *testing.T) {
	handler := handlers.NewBeritaHandler(&MockBeritaService{})

	body := dto.CreateBeritaRequest{Judul: "T"} // Judul terlalu pendek
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/berita", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest, got %v", rec.Code)
	}
}

func TestBeritaHandler_Create_Unauthorized(t *testing.T) {
	handler := handlers.NewBeritaHandler(&MockBeritaService{})

	body := dto.CreateBeritaRequest{Judul: "Test Berita", Deskripsi: "Deskripsi Berita"}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/berita", bytes.NewReader(bodyBytes))
	// No user context
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status Unauthorized, got %v", rec.Code)
	}
}

func TestBeritaHandler_Create_ServiceError(t *testing.T) {
	mockSvc := &MockBeritaService{
		CreateFunc: func(authorID string, req dto.CreateBeritaRequest) error {
			return errors.New("service error")
		},
	}
	handler := handlers.NewBeritaHandler(mockSvc)

	body := dto.CreateBeritaRequest{Judul: "Test Berita", Deskripsi: "Deskripsi Berita"}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/berita", bytes.NewReader(bodyBytes))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "author-123")
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status InternalServerError, got %v", rec.Code)
	}
}

func TestBeritaHandler_Update(t *testing.T) {
	mockSvc := &MockBeritaService{
		UpdateFunc: func(id int64, req dto.UpdateBeritaRequest) error {
			return nil
		},
	}
	handler := handlers.NewBeritaHandler(mockSvc)

	judul := "Updated Berita"
	deskripsi := "<script>alert('xss')</script>Deskripsi"
	body := dto.UpdateBeritaRequest{Judul: &judul, Deskripsi: &deskripsi}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/berita/1", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", rec.Code)
	}
}

func TestBeritaHandler_Update_InvalidID(t *testing.T) {
	handler := handlers.NewBeritaHandler(&MockBeritaService{})

	req := httptest.NewRequest(http.MethodPut, "/api/berita/abc", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest, got %v", rec.Code)
	}
}

func TestBeritaHandler_Update_InvalidBody(t *testing.T) {
	handler := handlers.NewBeritaHandler(&MockBeritaService{})

	req := httptest.NewRequest(http.MethodPut, "/api/berita/1", bytes.NewReader([]byte("{invalid}")))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest, got %v", rec.Code)
	}
}

func TestBeritaHandler_Update_ValidationError(t *testing.T) {
	handler := handlers.NewBeritaHandler(&MockBeritaService{})

	judul := "T" // Terlalu pendek
	body := dto.UpdateBeritaRequest{Judul: &judul}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/berita/1", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest, got %v", rec.Code)
	}
}

func TestBeritaHandler_Update_ServiceError(t *testing.T) {
	mockSvc := &MockBeritaService{
		UpdateFunc: func(id int64, req dto.UpdateBeritaRequest) error {
			return errors.New("service error")
		},
	}
	handler := handlers.NewBeritaHandler(mockSvc)

	judul := "Updated Berita"
	body := dto.UpdateBeritaRequest{Judul: &judul}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/berita/1", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest, got %v", rec.Code)
	}
}

func TestBeritaHandler_Delete(t *testing.T) {
	mockSvc := &MockBeritaService{
		DeleteFunc: func(id int64) error {
			return nil
		},
	}
	handler := handlers.NewBeritaHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/api/berita/1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", rec.Code)
	}
}

func TestBeritaHandler_Delete_InvalidID(t *testing.T) {
	handler := handlers.NewBeritaHandler(&MockBeritaService{})

	req := httptest.NewRequest(http.MethodDelete, "/api/berita/abc", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest, got %v", rec.Code)
	}
}

func TestBeritaHandler_Delete_ServiceError(t *testing.T) {
	mockSvc := &MockBeritaService{
		DeleteFunc: func(id int64) error {
			return errors.New("service error")
		},
	}
	handler := handlers.NewBeritaHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/api/berita/1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest, got %v", rec.Code)
	}
}

func TestBeritaHandler_ServeHTTP_MethodNotAllowed(t *testing.T) {
	handler := handlers.NewBeritaHandler(&MockBeritaService{})

	req := httptest.NewRequest(http.MethodPatch, "/api/berita/1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status MethodNotAllowed, got %v", rec.Code)
	}
}
