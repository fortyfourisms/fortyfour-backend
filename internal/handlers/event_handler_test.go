package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockEventService implements services.EventServiceInterface for tests
type MockEventService struct {
	CreateFunc  func(req dto.CreateEventRequest) error
	GetAllFunc  func() ([]dto.EventResponse, error)
	GetByIDFunc func(id int64) (*dto.EventResponse, error)
	UpdateFunc  func(id int64, req dto.UpdateEventRequest) error
	DeleteFunc  func(id int64) error
	RegisterFunc func(eventID int64, req dto.CreateEventRegistrationRequest) (*dto.EventRegistrationResponse, error)
	PDFFunc      func(registrationID int64) ([]byte, string, error)
}

func (m *MockEventService) Create(req dto.CreateEventRequest) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(req)
	}
	return nil
}

func (m *MockEventService) GetAll() ([]dto.EventResponse, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc()
	}
	return nil, nil
}

func (m *MockEventService) GetByID(id int64) (*dto.EventResponse, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *MockEventService) Update(id int64, req dto.UpdateEventRequest) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(id, req)
	}
	return nil
}

func (m *MockEventService) Delete(id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func (m *MockEventService) Register(eventID int64, req dto.CreateEventRegistrationRequest) (*dto.EventRegistrationResponse, error) {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(eventID, req)
	}
	return nil, nil
}

func (m *MockEventService) DownloadRegistrationPDF(registrationID int64) ([]byte, string, error) {
	if m.PDFFunc != nil {
		return m.PDFFunc(registrationID)
	}
	return nil, "", nil
}

// ── CRUD Tests (from main) ──────────────────────────────────────────────────

func TestEventHandler_GetAll(t *testing.T) {
	mockSvc := &MockEventService{
		GetAllFunc: func() ([]dto.EventResponse, error) {
			return []dto.EventResponse{{ID: 1, Judul: "Test Event"}}, nil
		},
	}
	handler := handlers.NewEventHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/kegiatan", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", rec.Code)
	}
}

func TestEventHandler_GetAll_Error(t *testing.T) {
	mockSvc := &MockEventService{
		GetAllFunc: func() ([]dto.EventResponse, error) {
			return nil, errors.New("db error")
		},
	}
	handler := handlers.NewEventHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/kegiatan", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status InternalServerError, got %v", rec.Code)
	}
}

func TestEventHandler_GetByID(t *testing.T) {
	mockSvc := &MockEventService{
		GetByIDFunc: func(id int64) (*dto.EventResponse, error) {
			return &dto.EventResponse{ID: 1, Judul: "Test Event"}, nil
		},
	}
	handler := handlers.NewEventHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/kegiatan/1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", rec.Code)
	}
}

func TestEventHandler_GetByID_InvalidID(t *testing.T) {
	handler := handlers.NewEventHandler(&MockEventService{})

	req := httptest.NewRequest(http.MethodGet, "/api/kegiatan/abc", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest, got %v", rec.Code)
	}
}

func TestEventHandler_GetByID_Error(t *testing.T) {
	mockSvc := &MockEventService{
		GetByIDFunc: func(id int64) (*dto.EventResponse, error) {
			return nil, errors.New("not found")
		},
	}
	handler := handlers.NewEventHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/kegiatan/1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status NotFound, got %v", rec.Code)
	}
}

func TestEventHandler_Create(t *testing.T) {
	mockSvc := &MockEventService{
		CreateFunc: func(req dto.CreateEventRequest) error {
			return nil
		},
	}
	handler := handlers.NewEventHandler(mockSvc)

	body := dto.CreateEventRequest{
		Judul:     "Test Event",
		Deskripsi: "Deskripsi",
		Lokasi:    "Lokasi",
		Tanggal:   "2026-12-31T10:00:00Z",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/kegiatan", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status Created, got %v", rec.Code)
	}
}

func TestEventHandler_Create_InvalidBody(t *testing.T) {
	handler := handlers.NewEventHandler(&MockEventService{})

	req := httptest.NewRequest(http.MethodPost, "/api/kegiatan", bytes.NewReader([]byte("{invalid-json}")))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest, got %v", rec.Code)
	}
}

func TestEventHandler_Create_ValidationError(t *testing.T) {
	handler := handlers.NewEventHandler(&MockEventService{})

	body := dto.CreateEventRequest{Judul: "T"} // Judul terlalu pendek
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/kegiatan", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest, got %v", rec.Code)
	}
}

func TestEventHandler_Create_ServiceError(t *testing.T) {
	mockSvc := &MockEventService{
		CreateFunc: func(req dto.CreateEventRequest) error {
			return errors.New("service error")
		},
	}
	handler := handlers.NewEventHandler(mockSvc)

	body := dto.CreateEventRequest{
		Judul:     "Test Event",
		Deskripsi: "Deskripsi",
		Lokasi:    "Lokasi",
		Tanggal:   "2026-12-31T10:00:00Z",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/kegiatan", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status InternalServerError, got %v", rec.Code)
	}
}

func TestEventHandler_Update(t *testing.T) {
	mockSvc := &MockEventService{
		UpdateFunc: func(id int64, req dto.UpdateEventRequest) error {
			return nil
		},
	}
	handler := handlers.NewEventHandler(mockSvc)

	judul := "Updated Event"
	deskripsi := "<script>alert('xss')</script>Deskripsi"
	body := dto.UpdateEventRequest{Judul: &judul, Deskripsi: &deskripsi}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/kegiatan/1", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", rec.Code)
	}
}

func TestEventHandler_Update_InvalidID(t *testing.T) {
	handler := handlers.NewEventHandler(&MockEventService{})

	req := httptest.NewRequest(http.MethodPut, "/api/kegiatan/abc", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest, got %v", rec.Code)
	}
}

func TestEventHandler_Update_InvalidBody(t *testing.T) {
	handler := handlers.NewEventHandler(&MockEventService{})

	req := httptest.NewRequest(http.MethodPut, "/api/kegiatan/1", bytes.NewReader([]byte("{invalid}")))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest, got %v", rec.Code)
	}
}

func TestEventHandler_Update_ValidationError(t *testing.T) {
	handler := handlers.NewEventHandler(&MockEventService{})

	judul := "T" // Terlalu pendek
	body := dto.UpdateEventRequest{Judul: &judul}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/kegiatan/1", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest, got %v", rec.Code)
	}
}

func TestEventHandler_Update_ServiceError(t *testing.T) {
	mockSvc := &MockEventService{
		UpdateFunc: func(id int64, req dto.UpdateEventRequest) error {
			return errors.New("service error")
		},
	}
	handler := handlers.NewEventHandler(mockSvc)

	judul := "Updated Event"
	body := dto.UpdateEventRequest{Judul: &judul}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/kegiatan/1", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest, got %v", rec.Code)
	}
}

func TestEventHandler_Delete(t *testing.T) {
	mockSvc := &MockEventService{
		DeleteFunc: func(id int64) error {
			return nil
		},
	}
	handler := handlers.NewEventHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/api/kegiatan/1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", rec.Code)
	}
}

func TestEventHandler_Delete_InvalidID(t *testing.T) {
	handler := handlers.NewEventHandler(&MockEventService{})

	req := httptest.NewRequest(http.MethodDelete, "/api/kegiatan/abc", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest, got %v", rec.Code)
	}
}

func TestEventHandler_Delete_ServiceError(t *testing.T) {
	mockSvc := &MockEventService{
		DeleteFunc: func(id int64) error {
			return errors.New("service error")
		},
	}
	handler := handlers.NewEventHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/api/kegiatan/1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest, got %v", rec.Code)
	}
}

func TestEventHandler_ServeHTTP_MethodNotAllowed(t *testing.T) {
	handler := handlers.NewEventHandler(&MockEventService{})

	req := httptest.NewRequest(http.MethodPatch, "/api/kegiatan/1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status MethodNotAllowed, got %v", rec.Code)
	}
}

// ── Registration Tests (from branch) ────────────────────────────────────────

func TestEventHandler_Register(t *testing.T) {
	handler := handlers.NewEventHandler(&MockEventService{
		RegisterFunc: func(eventID int64, req dto.CreateEventRegistrationRequest) (*dto.EventRegistrationResponse, error) {
			return &dto.EventRegistrationResponse{
				ID:           10,
				EventID:      eventID,
				Nama:         req.Nama,
				Email:        req.Email,
				Perusahaan:   req.Perusahaan,
				Jabatan:      req.Jabatan,
				NoHP:         req.NoHP,
				Sektor:       req.Sektor,
				QRCodeBase64: "ZmFrZQ==",
				DownloadURL:  "/api/kegiatan/registrasi/10/download",
			}, nil
		},
	})

	body, _ := json.Marshal(dto.CreateEventRegistrationRequest{
		Nama:       "Budi Santoso",
		Email:      "budi@example.com",
		Perusahaan: "PT ABC",
		Jabatan:    "IT Manager",
		NoHP:       "08123456789",
		Sektor:     "Energi",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/kegiatan/9/registrasi", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status Created, got %v", w.Code)
	}
}

func TestEventHandler_DownloadRegistrationPDF(t *testing.T) {
	handler := handlers.NewEventHandler(&MockEventService{
		PDFFunc: func(registrationID int64) ([]byte, string, error) {
			return []byte("%PDF-1.4 fake"), "registrasi-event-9-10.pdf", nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/kegiatan/registrasi/10/download", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", w.Code)
	}
	if w.Header().Get("Content-Type") != "application/pdf" {
		t.Errorf("expected Content-Type application/pdf, got %v", w.Header().Get("Content-Type"))
	}
}
