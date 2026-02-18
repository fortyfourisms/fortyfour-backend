package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"fortyfour-backend/internal/dto"
)

//
// MOCK SERVICE
//

type mockSdmCsirtService struct {
	CreateFn  func(dto.CreateSdmCsirtRequest) (string, error)
	GetAllFn  func() ([]dto.SdmCsirtResponse, error)
	GetByIDFn func(string) (*dto.SdmCsirtResponse, error)
	UpdateFn  func(string, dto.UpdateSdmCsirtRequest) error
	DeleteFn  func(string) error
}

func (m *mockSdmCsirtService) Create(req dto.CreateSdmCsirtRequest) (string, error) {
	return m.CreateFn(req)
}

func (m *mockSdmCsirtService) GetAll() ([]dto.SdmCsirtResponse, error) {
	return m.GetAllFn()
}

func (m *mockSdmCsirtService) GetByID(id string) (*dto.SdmCsirtResponse, error) {
	return m.GetByIDFn(id)
}

func (m *mockSdmCsirtService) Update(id string, req dto.UpdateSdmCsirtRequest) error {
	return m.UpdateFn(id, req)
}

func (m *mockSdmCsirtService) Delete(id string) error {
	return m.DeleteFn(id)
}

//
// TESTS
//

func TestSdmCsirtHandler_GetAll_Success(t *testing.T) {
	mockSvc := &mockSdmCsirtService{
		GetAllFn: func() ([]dto.SdmCsirtResponse, error) {
			return []dto.SdmCsirtResponse{
				{ID: "1", NamaPersonel: "Andi"},
			}, nil
		},
	}

	handler := NewSdmCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/sdm_csirt", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestSdmCsirtHandler_GetByID_Success(t *testing.T) {
	mockSvc := &mockSdmCsirtService{
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return &dto.SdmCsirtResponse{
				ID:           id,
				NamaPersonel: "Budi",
			}, nil
		},
	}

	handler := NewSdmCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/sdm_csirt/123", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestSdmCsirtHandler_GetByID_NotFound(t *testing.T) {
	mockSvc := &mockSdmCsirtService{
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return nil, errors.New("not found")
		},
	}

	handler := NewSdmCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/sdm_csirt/999", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestSdmCsirtHandler_Create_Success(t *testing.T) {
	mockSvc := &mockSdmCsirtService{
		CreateFn: func(req dto.CreateSdmCsirtRequest) (string, error) {
			return "uuid-123", nil
		},
	}

	body, _ := json.Marshal(dto.CreateSdmCsirtRequest{
		NamaPersonel: stringPtr("Charlie"),
	})

	handler := NewSdmCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/api/sdm_csirt", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
}

func TestSdmCsirtHandler_Create_Failed(t *testing.T) {
	mockSvc := &mockSdmCsirtService{
		CreateFn: func(req dto.CreateSdmCsirtRequest) (string, error) {
			return "", errors.New("create failed")
		},
	}

	handler := NewSdmCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/api/sdm_csirt", bytes.NewBuffer([]byte(`{}`)))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestSdmCsirtHandler_Update_Success(t *testing.T) {
	mockSvc := &mockSdmCsirtService{
		UpdateFn: func(id string, req dto.UpdateSdmCsirtRequest) error {
			return nil
		},
	}

	body, _ := json.Marshal(dto.UpdateSdmCsirtRequest{
		NamaPersonel: stringPtr("Updated"),
	})

	handler := NewSdmCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPut, "/api/sdm_csirt/1", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestSdmCsirtHandler_Update_Failed(t *testing.T) {
	mockSvc := &mockSdmCsirtService{
		UpdateFn: func(id string, req dto.UpdateSdmCsirtRequest) error {
			return errors.New("update failed")
		},
	}

	handler := NewSdmCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPut, "/api/sdm_csirt/1", bytes.NewBuffer([]byte(`{}`)))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestSdmCsirtHandler_Delete_Success(t *testing.T) {
	mockSvc := &mockSdmCsirtService{
		DeleteFn: func(id string) error {
			return nil
		},
	}

	handler := NewSdmCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/api/sdm_csirt/1", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestSdmCsirtHandler_Delete_Failed(t *testing.T) {
	mockSvc := &mockSdmCsirtService{
		DeleteFn: func(id string) error {
			return errors.New("delete failed")
		},
	}

	handler := NewSdmCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/api/sdm_csirt/1", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}
