package services

import (
	"errors"
	"testing"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type eventServiceRepoMock struct {
	eventByID                    *models.Event
	eventByIDErr                 error
	registrationByID             *models.EventRegistration
	registrationByIDErr          error
	existsRegistration           bool
	existsRegistrationErr        error
	createRegistrationErr        error
	updateRegistrationPayloadErr error
	createdRegistration          *models.EventRegistration
	updatedPayload               string
}

func (m *eventServiceRepoMock) Create(event *models.Event) error { return nil }
func (m *eventServiceRepoMock) FindAll() ([]models.Event, error) { return nil, nil }
func (m *eventServiceRepoMock) FindByID(id int64) (*models.Event, error) {
	return m.eventByID, m.eventByIDErr
}
func (m *eventServiceRepoMock) Update(event *models.Event) error { return nil }
func (m *eventServiceRepoMock) Delete(id int64) error            { return nil }
func (m *eventServiceRepoMock) CreateRegistration(reg *models.EventRegistration) error {
	if m.createRegistrationErr != nil {
		return m.createRegistrationErr
	}
	reg.ID = 101
	reg.CreatedAt = time.Date(2026, 4, 27, 10, 0, 0, 0, time.UTC)
	reg.UpdatedAt = reg.CreatedAt
	cp := *reg
	m.createdRegistration = &cp
	return nil
}
func (m *eventServiceRepoMock) FindRegistrationByID(id int64) (*models.EventRegistration, error) {
	return m.registrationByID, m.registrationByIDErr
}
func (m *eventServiceRepoMock) ExistsRegistrationByEventAndEmail(eventID int64, email string) (bool, error) {
	return m.existsRegistration, m.existsRegistrationErr
}
func (m *eventServiceRepoMock) UpdateRegistrationPayload(id int64, payload string) error {
	m.updatedPayload = payload
	return m.updateRegistrationPayloadErr
}

func TestEventService_Register_Success(t *testing.T) {
	repo := &eventServiceRepoMock{
		eventByID: &models.Event{
			ID:      9,
			Judul:   "Cyber Drill 2026",
			Tanggal: time.Now().Add(24 * time.Hour),
		},
	}
	svc := NewEventService(repo)

	resp, err := svc.Register(9, dto.CreateEventRegistrationRequest{
		Nama:       "Budi Santoso",
		Email:      "budi@example.com",
		Perusahaan: "PT ABC",
		Jabatan:    "IT Manager",
		NoHP:       "08123456789",
		Sektor:     "Energi",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int64(101), resp.ID)
	assert.Equal(t, "/api/kegiatan/registrasi/101/download", resp.DownloadURL)
	assert.NotEmpty(t, resp.QRCodeBase64)
	assert.Contains(t, resp.QRPayload, "\"registration_id\":101")
	assert.Equal(t, repo.updatedPayload, resp.QRPayload)
	assert.Equal(t, "Cyber Drill 2026", repo.eventByID.Judul)
}

func TestEventService_Register_DuplicateEmail(t *testing.T) {
	repo := &eventServiceRepoMock{
		eventByID:          &models.Event{ID: 9, Judul: "Cyber Drill 2026"},
		existsRegistration: true,
	}
	svc := NewEventService(repo)

	resp, err := svc.Register(9, dto.CreateEventRegistrationRequest{
		Nama:       "Budi Santoso",
		Email:      "budi@example.com",
		Perusahaan: "PT ABC",
		Jabatan:    "IT Manager",
		NoHP:       "08123456789",
		Sektor:     "Energi",
	})

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "email sudah terdaftar pada event ini", err.Error())
}

func TestEventService_DownloadRegistrationPDF_Success(t *testing.T) {
	repo := &eventServiceRepoMock{
		eventByID: &models.Event{ID: 9, Judul: "Cyber Drill 2026"},
		registrationByID: &models.EventRegistration{
			ID:         101,
			EventID:    9,
			Nama:       "Budi Santoso",
			Email:      "budi@example.com",
			Perusahaan: "PT ABC",
			Jabatan:    "IT Manager",
			NoHP:       "08123456789",
			Sektor:     "Energi",
			QRToken:    "qr-123",
			QRPayload:  `{"registration_id":101,"nama":"Budi Santoso"}`,
		},
	}
	svc := NewEventService(repo)

	pdf, filename, err := svc.DownloadRegistrationPDF(101)

	require.NoError(t, err)
	assert.NotEmpty(t, pdf)
	assert.Equal(t, "registrasi-event-9-101.pdf", filename)
	assert.Equal(t, "%PDF", string(pdf[:4]))
}

func TestEventService_DownloadRegistrationPDF_NotFound(t *testing.T) {
	repo := &eventServiceRepoMock{
		registrationByIDErr: errors.New("db error"),
	}
	svc := NewEventService(repo)

	pdf, filename, err := svc.DownloadRegistrationPDF(101)

	require.Error(t, err)
	assert.Nil(t, pdf)
	assert.Empty(t, filename)
}
