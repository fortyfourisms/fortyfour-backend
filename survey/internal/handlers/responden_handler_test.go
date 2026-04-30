package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"survey/internal/dto"
)

// MOCK SERVICE
type mockService struct {
	GetAllFunc func() ([]dto.RespondenResponse, error)
	GetByIDFunc func(int) (*dto.RespondenResponse, error)
	CreateFunc func(dto.CreateRespondenRequest) (*dto.RespondenResponse, error)
	UpdateFunc func(int, dto.UpdateRespondenRequest) (*dto.RespondenResponse, error)
}

func (m *mockService) GetAll() ([]dto.RespondenResponse, error) {
	return m.GetAllFunc()
}

func (m *mockService) GetByID(id int) (*dto.RespondenResponse, error) {
	return m.GetByIDFunc(id)
}

func (m *mockService) Create(req dto.CreateRespondenRequest) (*dto.RespondenResponse, error) {
	return m.CreateFunc(req)
}

func (m *mockService) Update(id int, req dto.UpdateRespondenRequest) (*dto.RespondenResponse, error) {
	return m.UpdateFunc(id, req)
}

// GET ALL
func TestGetAllResponden(t *testing.T) {
	mock := &mockService{
		GetAllFunc: func() ([]dto.RespondenResponse, error) {
			return []dto.RespondenResponse{}, nil
		},
	}

	handler := NewRespondenHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/survey/responden", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// GET BY ID SUCCESS
func TestGetByID_Success(t *testing.T) {
	mock := &mockService{
		GetByIDFunc: func(id int) (*dto.RespondenResponse, error) {
			return &dto.RespondenResponse{}, nil
		},
	}

	handler := NewRespondenHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/survey/responden/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// GET BY ID INVALID
func TestGetByID_InvalidID(t *testing.T) {
	handler := NewRespondenHandler(&mockService{})

	req := httptest.NewRequest(http.MethodGet, "/api/survey/responden/abc", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// CREATE SUCCESS
func TestCreateResponden_Success(t *testing.T) {
	mock := &mockService{
		CreateFunc: func(req dto.CreateRespondenRequest) (*dto.RespondenResponse, error) {
			return &dto.RespondenResponse{}, nil
		},
	}

	handler := NewRespondenHandler(mock)

	body, _ := json.Marshal(dto.CreateRespondenRequest{})
	req := httptest.NewRequest(http.MethodPost, "/api/survey/responden", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

// CREATE INVALID BODY
func TestCreateResponden_InvalidBody(t *testing.T) {
	handler := NewRespondenHandler(&mockService{})

	req := httptest.NewRequest(http.MethodPost, "/api/survey/responden", bytes.NewBuffer([]byte("invalid")))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// UPDATE SUCCESS
func TestUpdateResponden_Success(t *testing.T) {
	mock := &mockService{
		UpdateFunc: func(id int, req dto.UpdateRespondenRequest) (*dto.RespondenResponse, error) {
			return &dto.RespondenResponse{}, nil
		},
	}

	handler := NewRespondenHandler(mock)

	body, _ := json.Marshal(dto.UpdateRespondenRequest{})
	req := httptest.NewRequest(http.MethodPut, "/api/survey/responden/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// METHOD NOT ALLOWED
func TestMethodNotAllowed(t *testing.T) {
	handler := NewRespondenHandler(&mockService{})

	req := httptest.NewRequest(http.MethodDelete, "/api/survey/responden", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}