package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"fortyfour-backend/internal/dto"
)

//
// MOCK SERVICE
//

type mockSektorService struct {
	GetAllFn  func() ([]dto.SektorResponse, error)
	GetByIDFn func(string) (*dto.SektorResponse, error)
}

func (m *mockSektorService) GetAll() ([]dto.SektorResponse, error) {
	return m.GetAllFn()
}

func (m *mockSektorService) GetByID(id string) (*dto.SektorResponse, error) {
	return m.GetByIDFn(id)
}

//
// TESTS – SUCCESS PATH
//

func TestSektorHandler_GetAll_Success(t *testing.T) {
	mockSvc := &mockSektorService{
		GetAllFn: func() ([]dto.SektorResponse, error) {
			return []dto.SektorResponse{
				{
					ID:        "1",
					NamaSektor: "Sektor A",
				},
			}, nil
		},
	}

	handler := NewSektorHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/sektor", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestSektorHandler_GetByID_Success(t *testing.T) {
	mockSvc := &mockSektorService{
		GetByIDFn: func(id string) (*dto.SektorResponse, error) {
			return &dto.SektorResponse{
				ID:        id,
				NamaSektor: "Sektor A",
			}, nil
		},
	}

	handler := NewSektorHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/sektor/1", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

//
// TESTS – ERROR PATH
//

func TestSektorHandler_GetAll_InternalError(t *testing.T) {
	mockSvc := &mockSektorService{
		GetAllFn: func() ([]dto.SektorResponse, error) {
			return nil, errors.New("db error")
		},
	}

	handler := NewSektorHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/sektor", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}

func TestSektorHandler_GetByID_NotFound(t *testing.T) {
	mockSvc := &mockSektorService{
		GetByIDFn: func(id string) (*dto.SektorResponse, error) {
			return nil, errors.New("not found")
		},
	}

	handler := NewSektorHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/sektor/999", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

//
// TESTS – METHOD NOT ALLOWED
//

func TestSektorHandler_MethodNotAllowed(t *testing.T) {
	mockSvc := &mockSektorService{}

	handler := NewSektorHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/api/sektor", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}
