package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/handlers"
)

// ═══════════════════════════════════════════════════════════════════════════
// MOCK: KonversiService
// ═══════════════════════════════════════════════════════════════════════════

type MockKonversiService struct {
	GetKonversiFn func(perusahaanID string) ([]dto.KonversiResponse, error)
}

func (m *MockKonversiService) GetKonversi(perusahaanID string) ([]dto.KonversiResponse, error) {
	if m.GetKonversiFn != nil {
		return m.GetKonversiFn(perusahaanID)
	}
	return nil, nil
}

// ═══════════════════════════════════════════════════════════════════════════
// TESTS
// ═══════════════════════════════════════════════════════════════════════════

func TestKonversiHandler_GetAll_Success(t *testing.T) {
	mockSvc := &MockKonversiService{
		GetKonversiFn: func(perusahaanID string) ([]dto.KonversiResponse, error) {
			return []dto.KonversiResponse{
				{PerusahaanID: "uuid-1", NamaPerusahaan: "PT A", TotalPoin: 3, Persentase: 75.0},
			}, nil
		},
	}
	handler := handlers.NewKonversiHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/konversi", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", rec.Code)
	}
}

func TestKonversiHandler_GetByID_QueryParam(t *testing.T) {
	var receivedID string
	mockSvc := &MockKonversiService{
		GetKonversiFn: func(perusahaanID string) ([]dto.KonversiResponse, error) {
			receivedID = perusahaanID
			return []dto.KonversiResponse{
				{PerusahaanID: perusahaanID, NamaPerusahaan: "PT A"},
			}, nil
		},
	}
	handler := handlers.NewKonversiHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/konversi?perusahaan_id=uuid-123", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", rec.Code)
	}
	if receivedID != "uuid-123" {
		t.Errorf("expected perusahaan_id 'uuid-123', got '%s'", receivedID)
	}
}

func TestKonversiHandler_GetByID_PathParam(t *testing.T) {
	var receivedID string
	mockSvc := &MockKonversiService{
		GetKonversiFn: func(perusahaanID string) ([]dto.KonversiResponse, error) {
			receivedID = perusahaanID
			return []dto.KonversiResponse{
				{PerusahaanID: perusahaanID, NamaPerusahaan: "PT A"},
			}, nil
		},
	}
	handler := handlers.NewKonversiHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/konversi/uuid-456", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", rec.Code)
	}
	if receivedID != "uuid-456" {
		t.Errorf("expected perusahaan_id 'uuid-456', got '%s'", receivedID)
	}
}

func TestKonversiHandler_ServiceError(t *testing.T) {
	mockSvc := &MockKonversiService{
		GetKonversiFn: func(perusahaanID string) ([]dto.KonversiResponse, error) {
			return nil, errors.New("db error")
		},
	}
	handler := handlers.NewKonversiHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/konversi", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status InternalServerError, got %v", rec.Code)
	}
}

func TestKonversiHandler_MethodNotAllowed(t *testing.T) {
	handler := handlers.NewKonversiHandler(&MockKonversiService{})

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}
	for _, method := range methods {
		req := httptest.NewRequest(method, "/api/konversi", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("[%s] expected status MethodNotAllowed, got %v", method, rec.Code)
		}
	}
}

func TestKonversiHandler_QueryParamTakesPrecedence(t *testing.T) {
	var receivedID string
	mockSvc := &MockKonversiService{
		GetKonversiFn: func(perusahaanID string) ([]dto.KonversiResponse, error) {
			receivedID = perusahaanID
			return []dto.KonversiResponse{}, nil
		},
	}
	handler := handlers.NewKonversiHandler(mockSvc)

	// Query param should take precedence over path param
	req := httptest.NewRequest(http.MethodGet, "/api/konversi/path-id?perusahaan_id=query-id", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", rec.Code)
	}
	if receivedID != "query-id" {
		t.Errorf("expected query param 'query-id' to take precedence, got '%s'", receivedID)
	}
}

func TestKonversiHandler_EmptyResult(t *testing.T) {
	mockSvc := &MockKonversiService{
		GetKonversiFn: func(perusahaanID string) ([]dto.KonversiResponse, error) {
			return []dto.KonversiResponse{}, nil
		},
	}
	handler := handlers.NewKonversiHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/konversi", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", rec.Code)
	}
}
