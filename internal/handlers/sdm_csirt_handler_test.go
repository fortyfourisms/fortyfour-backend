package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/services"
)

//
// MOCK SDM CSIRT SERVICE
//

type mockSdmCsirtService struct {
	CreateFn     func(dto.CreateSdmCsirtRequest) (string, error)
	GetAllFn     func() ([]dto.SdmCsirtResponse, error)
	GetByIDFn    func(string) (*dto.SdmCsirtResponse, error)
	UpdateFn     func(string, dto.UpdateSdmCsirtRequest) error
	DeleteFn     func(string) error
	GetByCsirtFn func(string) ([]dto.SdmCsirtResponse, error)
}

func (m *mockSdmCsirtService) Create(req dto.CreateSdmCsirtRequest) (string, error) {
	if m.CreateFn != nil {
		return m.CreateFn(req)
	}
	return "", nil
}

func (m *mockSdmCsirtService) GetAll() ([]dto.SdmCsirtResponse, error) {
	if m.GetAllFn != nil {
		return m.GetAllFn()
	}
	return []dto.SdmCsirtResponse{}, nil
}

func (m *mockSdmCsirtService) GetByID(id string) (*dto.SdmCsirtResponse, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(id)
	}
	return &dto.SdmCsirtResponse{ID: id}, nil
}

func (m *mockSdmCsirtService) Update(id string, req dto.UpdateSdmCsirtRequest) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(id, req)
	}
	return nil
}

func (m *mockSdmCsirtService) Delete(id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id)
	}
	return nil
}

func (m *mockSdmCsirtService) GetByCsirt(idCsirt string) ([]dto.SdmCsirtResponse, error) {
	if m.GetByCsirtFn != nil {
		return m.GetByCsirtFn(idCsirt)
	}
	return []dto.SdmCsirtResponse{}, nil
}

//
// MOCK CSIRT SERVICE
//

type mockCsirtServiceForSdm struct{}

func (m *mockCsirtServiceForSdm) GetAll() ([]dto.CsirtResponse, error) {
	return []dto.CsirtResponse{}, nil
}
func (m *mockCsirtServiceForSdm) GetByID(id string) (*dto.CsirtResponse, error) {
	return nil, errors.New("not found")
}
func (m *mockCsirtServiceForSdm) GetByPerusahaan(idPerusahaan string) ([]dto.CsirtResponse, error) {
	return []dto.CsirtResponse{}, nil
}
func (m *mockCsirtServiceForSdm) Create(req dto.CreateCsirtRequest) (*models.Csirt, error) {
	return nil, nil
}
func (m *mockCsirtServiceForSdm) Update(id string, req dto.UpdateCsirtRequest) (*models.Csirt, error) {
	return nil, nil
}
func (m *mockCsirtServiceForSdm) Delete(id string) error { return nil }

//
// HELPER
//

func newSdmHandler(mockSvc *mockSdmCsirtService) *SdmCsirtHandler {
	sseService := services.NewSSEService()
	return NewSdmCsirtHandler(mockSvc, &mockCsirtServiceForSdm{}, sseService)
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

	handler := newSdmHandler(mockSvc)

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

	handler := newSdmHandler(mockSvc)

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

	handler := newSdmHandler(mockSvc)

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
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return &dto.SdmCsirtResponse{ID: id, NamaPersonel: "Charlie"}, nil
		},
	}

	body, _ := json.Marshal(dto.CreateSdmCsirtRequest{
		NamaPersonel: stringPtr("Charlie"),
	})

	handler := newSdmHandler(mockSvc)

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

	handler := newSdmHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/api/sdm_csirt", bytes.NewBuffer([]byte(`{}`)))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestSdmCsirtHandler_Update_Success(t *testing.T) {
	mockSvc := &mockSdmCsirtService{
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return &dto.SdmCsirtResponse{ID: id, NamaPersonel: "Updated"}, nil
		},
		UpdateFn: func(id string, req dto.UpdateSdmCsirtRequest) error {
			return nil
		},
	}

	body, _ := json.Marshal(dto.UpdateSdmCsirtRequest{
		NamaPersonel: stringPtr("Updated"),
	})

	handler := newSdmHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPut, "/api/sdm_csirt/1", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestSdmCsirtHandler_Update_Failed(t *testing.T) {
	mockSvc := &mockSdmCsirtService{
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return &dto.SdmCsirtResponse{ID: id}, nil
		},
		UpdateFn: func(id string, req dto.UpdateSdmCsirtRequest) error {
			return errors.New("update failed")
		},
	}

	handler := newSdmHandler(mockSvc)

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

	handler := newSdmHandler(mockSvc)

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

	handler := newSdmHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/api/sdm_csirt/1", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestSdmCsirtHandler_MethodNotAllowed(t *testing.T) {
	mockSvc := &mockSdmCsirtService{}
	handler := newSdmHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPatch, "/api/sdm_csirt", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}

func TestSdmCsirtHandler_GetAll_Error(t *testing.T) {
	mockSvc := &mockSdmCsirtService{
		GetAllFn: func() ([]dto.SdmCsirtResponse, error) {
			return nil, errors.New("db error")
		},
	}
	handler := newSdmHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/sdm_csirt", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestSdmCsirtHandler_GetAll_ResponseBody(t *testing.T) {
	mockSvc := &mockSdmCsirtService{
		GetAllFn: func() ([]dto.SdmCsirtResponse, error) {
			return []dto.SdmCsirtResponse{
				{ID: "1", NamaPersonel: "Andi"},
				{ID: "2", NamaPersonel: "Budi"},
			}, nil
		},
	}
	handler := newSdmHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/sdm_csirt", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var result []dto.SdmCsirtResponse
	json.NewDecoder(rr.Body).Decode(&result)
	if len(result) != 2 {
		t.Errorf("expected 2 items, got %d", len(result))
	}
}

func TestSdmCsirtHandler_Update_InvalidBody(t *testing.T) {
	mockSvc := &mockSdmCsirtService{
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return &dto.SdmCsirtResponse{ID: id}, nil
		},
		UpdateFn: func(id string, req dto.UpdateSdmCsirtRequest) error {
			return errors.New("update failed")
		},
	}
	handler := newSdmHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPut, "/api/sdm_csirt/1", bytes.NewBuffer([]byte("bad json")))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusOK {
		t.Error("tidak boleh 200 jika update gagal")
	}
}

func TestSdmCsirtHandler_Delete_WithoutID(t *testing.T) {
	mockSvc := &mockSdmCsirtService{
		DeleteFn: func(id string) error {
			return errors.New("tidak ditemukan")
		},
	}
	handler := newSdmHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/api/sdm_csirt/non-existent", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for delete error, got %d", rr.Code)
	}
}