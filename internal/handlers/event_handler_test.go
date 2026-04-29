package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"fortyfour-backend/internal/dto"

	"github.com/stretchr/testify/assert"
)

type mockEventService struct {
	registerFn func(eventID int64, req dto.CreateEventRegistrationRequest) (*dto.EventRegistrationResponse, error)
	pdfFn      func(registrationID int64) ([]byte, string, error)
}

func (m *mockEventService) Create(req dto.CreateEventRequest) (*dto.EventResponse, error) {
	return nil, nil
}
func (m *mockEventService) GetAll() ([]dto.EventResponse, error)         { return nil, nil }
func (m *mockEventService) GetByID(id int64) (*dto.EventResponse, error) { return nil, nil }
func (m *mockEventService) Update(id int64, req dto.UpdateEventRequest) (*dto.EventResponse, error) {
	return nil, nil
}
func (m *mockEventService) Delete(id int64) error { return nil }
func (m *mockEventService) Register(eventID int64, req dto.CreateEventRegistrationRequest) (*dto.EventRegistrationResponse, error) {
	return m.registerFn(eventID, req)
}
func (m *mockEventService) DownloadRegistrationPDF(registrationID int64) ([]byte, string, error) {
	return m.pdfFn(registrationID)
}

func TestEventHandler_Register(t *testing.T) {
	handler := NewEventHandler(&mockEventService{
		registerFn: func(eventID int64, req dto.CreateEventRegistrationRequest) (*dto.EventRegistrationResponse, error) {
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
		pdfFn: func(registrationID int64) ([]byte, string, error) { return nil, "", nil },
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

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Berhasil registrasi event")
	assert.Contains(t, w.Body.String(), "/api/kegiatan/registrasi/10/download")
}

func TestEventHandler_DownloadRegistrationPDF(t *testing.T) {
	handler := NewEventHandler(&mockEventService{
		registerFn: func(eventID int64, req dto.CreateEventRegistrationRequest) (*dto.EventRegistrationResponse, error) {
			return nil, nil
		},
		pdfFn: func(registrationID int64) ([]byte, string, error) {
			return []byte("%PDF-1.4 fake"), "registrasi-event-9-10.pdf", nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/kegiatan/registrasi/10/download", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "registrasi-event-9-10.pdf")
	assert.Equal(t, "%PDF-1.4 fake", w.Body.String())
}
