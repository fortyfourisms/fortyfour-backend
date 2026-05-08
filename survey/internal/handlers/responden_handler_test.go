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
	GetAllFunc         func() ([]dto.RespondenResponse, error)
	GetByIDFunc        func(int) (*dto.RespondenResponse, error)
	GetByUserIDFunc    func(string) (*dto.RespondenResponse, error)
	UpsertByUserIDFunc func(string, dto.CreateRespondenRequest) (*dto.RespondenResponse, error)
}

func (m *mockService) GetAll() ([]dto.RespondenResponse, error) {
	return m.GetAllFunc()
}

func (m *mockService) GetByID(id int) (*dto.RespondenResponse, error) {
	return m.GetByIDFunc(id)
}

func (m *mockService) GetByUserID(userID string) (*dto.RespondenResponse, error) {
	return m.GetByUserIDFunc(userID)
}

func (m *mockService) UpsertByUserID(userID string, req dto.CreateRespondenRequest) (*dto.RespondenResponse, error) {
	return m.UpsertByUserIDFunc(userID, req)
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
	req.Header.Set("X-User-Role", "admin")
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
	req.Header.Set("X-User-Role", "admin")
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
	req.Header.Set("X-User-Role", "admin")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// GET ME SUCCESS
func TestGetMe_Success(t *testing.T) {
	mock := &mockService{
		GetByUserIDFunc: func(userID string) (*dto.RespondenResponse, error) {
			return &dto.RespondenResponse{}, nil
		},
	}

	handler := NewRespondenHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/survey/responden/me", nil)
	req.Header.Set("X-User-ID", "user1")
	req.Header.Set("X-User-Role", "user")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// UPSERT ME SUCCESS
func TestUpsertMe_Success(t *testing.T) {
	mock := &mockService{
		UpsertByUserIDFunc: func(userID string, req dto.CreateRespondenRequest) (*dto.RespondenResponse, error) {
			return &dto.RespondenResponse{}, nil
		},
	}

	handler := NewRespondenHandler(mock)

	body, _ := json.Marshal(dto.CreateRespondenRequest{})
	req := httptest.NewRequest(http.MethodPost, "/api/survey/responden/me", bytes.NewBuffer(body))
	req.Header.Set("X-User-ID", "user1")
	req.Header.Set("X-User-Role", "user")
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
