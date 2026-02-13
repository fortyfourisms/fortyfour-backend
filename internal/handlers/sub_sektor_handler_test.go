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

type mockSubSektorService struct {
	GetAllFn        func() ([]dto.SubSektorResponse, error)
	GetByIDFn       func(string) (*dto.SubSektorResponse, error)
	GetBySektorIDFn func(string) ([]dto.SubSektorResponse, error)
}

func (m *mockSubSektorService) GetAll() ([]dto.SubSektorResponse, error) {
	return m.GetAllFn()
}

func (m *mockSubSektorService) GetByID(id string) (*dto.SubSektorResponse, error) {
	return m.GetByIDFn(id)
}

func (m *mockSubSektorService) GetBySektorID(sektorID string) ([]dto.SubSektorResponse, error) {
	return m.GetBySektorIDFn(sektorID)
}

//
// TESTS – SUCCESS PATH
//

func TestSubSektorHandler_GetAll_Success(t *testing.T) {
	mockSvc := &mockSubSektorService{
		GetAllFn: func() ([]dto.SubSektorResponse, error) {
			return []dto.SubSektorResponse{
				{ID: "1", NamaSubSektor: "Sub A"},
			}, nil
		},
	}

	handler := NewSubSektorHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/sub_sektor", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestSubSektorHandler_GetByID_Success(t *testing.T) {
	mockSvc := &mockSubSektorService{
		GetByIDFn: func(id string) (*dto.SubSektorResponse, error) {
			return &dto.SubSektorResponse{
				ID:            id,
				NamaSubSektor: "Sub A",
			}, nil
		},
	}

	handler := NewSubSektorHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/sub_sektor/1", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestSubSektorHandler_GetBySektorID_Success(t *testing.T) {
	mockSvc := &mockSubSektorService{
		GetBySektorIDFn: func(sektorID string) ([]dto.SubSektorResponse, error) {
			return []dto.SubSektorResponse{
				{ID: "1", NamaSubSektor: "Sub A"},
				{ID: "2", NamaSubSektor: "Sub B"},
			}, nil
		},
	}

	handler := NewSubSektorHandler(mockSvc)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/sub_sektor/by_sektor/10",
		nil,
	)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

//
// TESTS – ERROR PATH
//

func TestSubSektorHandler_GetByID_NotFound(t *testing.T) {
	mockSvc := &mockSubSektorService{
		GetByIDFn: func(id string) (*dto.SubSektorResponse, error) {
			return nil, errors.New("not found")
		},
	}

	handler := NewSubSektorHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/sub_sektor/999", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestSubSektorHandler_GetBySektorID_ServiceError(t *testing.T) {
	mockSvc := &mockSubSektorService{
		GetBySektorIDFn: func(sektorID string) ([]dto.SubSektorResponse, error) {
			return nil, errors.New("db error")
		},
	}

	handler := NewSubSektorHandler(mockSvc)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/sub_sektor/by_sektor/10",
		nil,
	)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}

//
// TESTS – EDGE / ROUTING
//

func TestSubSektorHandler_MethodNotAllowed(t *testing.T) {
	mockSvc := &mockSubSektorService{}

	handler := NewSubSektorHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/api/sub_sektor", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}

func TestSubSektorHandler_BySektor_WrongMethod(t *testing.T) {
	mockSvc := &mockSubSektorService{}

	handler := NewSubSektorHandler(mockSvc)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/sub_sektor/by_sektor/1",
		nil,
	)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}
