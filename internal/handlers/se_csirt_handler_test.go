package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/handlers"
)

/*
|--------------------------------------------------------------------------
| MOCK SERVICE
|--------------------------------------------------------------------------
*/

type mockSeCsirtService struct {
	CreateFn  func(req dto.CreateSeCsirtRequest) (string, error)
	GetAllFn  func() ([]dto.SeCsirtResponse, error)
	GetByIDFn func(id string) (*dto.SeCsirtResponse, error)
	UpdateFn  func(id string, req dto.UpdateSeCsirtRequest) error
	DeleteFn  func(id string) error
}

func (m *mockSeCsirtService) Create(req dto.CreateSeCsirtRequest) (string, error) {
	return m.CreateFn(req)
}

func (m *mockSeCsirtService) GetAll() ([]dto.SeCsirtResponse, error) {
	return m.GetAllFn()
}

func (m *mockSeCsirtService) GetByID(id string) (*dto.SeCsirtResponse, error) {
	return m.GetByIDFn(id)
}

func (m *mockSeCsirtService) Update(id string, req dto.UpdateSeCsirtRequest) error {
	return m.UpdateFn(id, req)
}

func (m *mockSeCsirtService) Delete(id string) error {
	return m.DeleteFn(id)
}

/*
|--------------------------------------------------------------------------
| TEST: GET ALL
|--------------------------------------------------------------------------
*/

func TestSeCsirtHandler_GetAll_Success(t *testing.T) {
	mockSvc := &mockSeCsirtService{
		GetAllFn: func() ([]dto.SeCsirtResponse, error) {
			return []dto.SeCsirtResponse{
				{ID: "1", NamaSe: "SE A"},
				{ID: "2", NamaSe: "SE B"},
			}, nil
		},
	}

	handler := handlers.NewSeCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/se_csirt", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestSeCsirtHandler_GetAll_Error(t *testing.T) {
	mockSvc := &mockSeCsirtService{
		GetAllFn: func() ([]dto.SeCsirtResponse, error) {
			return nil, errors.New("db error")
		},
	}

	handler := handlers.NewSeCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/se_csirt", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

/*
|--------------------------------------------------------------------------
| TEST: GET BY ID
|--------------------------------------------------------------------------
*/

func TestSeCsirtHandler_GetByID_Success(t *testing.T) {
	mockSvc := &mockSeCsirtService{
		GetByIDFn: func(id string) (*dto.SeCsirtResponse, error) {
			return &dto.SeCsirtResponse{ID: id, NamaSe: "SE Test"}, nil
		},
	}

	handler := handlers.NewSeCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/se_csirt/123", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestSeCsirtHandler_GetByID_NotFound(t *testing.T) {
	mockSvc := &mockSeCsirtService{
		GetByIDFn: func(id string) (*dto.SeCsirtResponse, error) {
			return nil, errors.New("not found")
		},
	}

	handler := handlers.NewSeCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/se_csirt/404", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

/*
|--------------------------------------------------------------------------
| TEST: CREATE
|--------------------------------------------------------------------------
*/

func TestSeCsirtHandler_Create_Success(t *testing.T) {
	mockSvc := &mockSeCsirtService{
		CreateFn: func(req dto.CreateSeCsirtRequest) (string, error) {
			return "new-id", nil
		},
	}

	handler := handlers.NewSeCsirtHandler(mockSvc)

	body, _ := json.Marshal(dto.CreateSeCsirtRequest{
		NamaSe: stringPtr("SE Baru"),
	})

	req := httptest.NewRequest(http.MethodPost, "/api/se_csirt", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

func TestSeCsirtHandler_Create_Error(t *testing.T) {
	mockSvc := &mockSeCsirtService{
		CreateFn: func(req dto.CreateSeCsirtRequest) (string, error) {
			return "", errors.New("create failed")
		},
	}

	handler := handlers.NewSeCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/api/se_csirt", bytes.NewBuffer([]byte(`{}`)))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

/*
|--------------------------------------------------------------------------
| TEST: UPDATE
|--------------------------------------------------------------------------
*/

func TestSeCsirtHandler_Update_Success(t *testing.T) {
	mockSvc := &mockSeCsirtService{
		UpdateFn: func(id string, req dto.UpdateSeCsirtRequest) error {
			return nil
		},
	}

	handler := handlers.NewSeCsirtHandler(mockSvc)

	body, _ := json.Marshal(dto.UpdateSeCsirtRequest{
		NamaSe: stringPtr("Updated"),
	})

	req := httptest.NewRequest(http.MethodPut, "/api/se_csirt/1", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestSeCsirtHandler_Update_Error(t *testing.T) {
	mockSvc := &mockSeCsirtService{
		UpdateFn: func(id string, req dto.UpdateSeCsirtRequest) error {
			return errors.New("update error")
		},
	}

	handler := handlers.NewSeCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPut, "/api/se_csirt/1", bytes.NewBuffer([]byte(`{}`)))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

/*
|--------------------------------------------------------------------------
| TEST: DELETE
|--------------------------------------------------------------------------
*/

func TestSeCsirtHandler_Delete_Success(t *testing.T) {
	mockSvc := &mockSeCsirtService{
		DeleteFn: func(id string) error {
			return nil
		},
	}

	handler := handlers.NewSeCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/api/se_csirt/1", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestSeCsirtHandler_Delete_Error(t *testing.T) {
	mockSvc := &mockSeCsirtService{
		DeleteFn: func(id string) error {
			return errors.New("delete error")
		},
	}

	handler := handlers.NewSeCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/api/se_csirt/1", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

/*
|--------------------------------------------------------------------------
| HELPER
|--------------------------------------------------------------------------
*/

func stringPtr(s string) *string {
	return &s
}
