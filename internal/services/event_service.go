package services

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/dto/dto_event"
	"fortyfour-backend/internal/models"
	internalRmq "fortyfour-backend/internal/rabbitmq"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/utils"
	"fortyfour-backend/pkg/cache"
	"strings"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type EventServiceInterface interface {
	Create(req dto.CreateEventRequest) error
	GetAll() ([]dto.EventResponse, error)
	GetByID(id string) (*dto.EventResponse, error)
	GetBySlug(slug string) (*dto.EventResponse, error)
	Update(id string, req dto.UpdateEventRequest) error
	Delete(id string) error
	Register(eventID string, req dto.CreateEventRegistrationRequest) (*dto.EventRegistrationResponse, error)
	DownloadRegistrationPDF(registrationID int64) ([]byte, string, error)
}

type EventService struct {
	repo     repository.EventRepositoryInterface
	producer *internalRmq.Producer
	rc       cache.RedisInterface
}

func NewEventService(
	repo repository.EventRepositoryInterface,
	producer *internalRmq.Producer,
	rc cache.RedisInterface,
) *EventService {
	return &EventService{
		repo:     repo,
		producer: producer,
		rc:       rc,
	}
}

var _ EventServiceInterface = (*EventService)(nil)

func (s *EventService) Create(req dto.CreateEventRequest) error {
	if _, err := time.Parse(time.RFC3339, req.Tanggal); err != nil {
		return errors.New("format tanggal tidak valid (gunakan RFC3339, contoh: 2024-12-31T15:00:00Z)")
	}

	// Check slug uniqueness
	slug := utils.Slugify(req.Judul)
	existing, err := s.repo.FindBySlug(slug)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("event dengan judul tersebut sudah ada")
	}

	event := dto_event.EventCreatedEvent{
		Request:   req,
		CreatedAt: time.Now(),
	}

	if s.producer != nil {
		err := s.producer.PublishEventCreated(context.Background(), event)
		if err != nil {
			return err
		}
	}

	// Invalidate list cache
	cacheDelete(s.rc, keyList("event"))
	return nil
}

func (s *EventService) GetAll() ([]dto.EventResponse, error) {
	key := keyList("event")
	var res []dto.EventResponse

	if cacheGet(s.rc, key, &res) {
		return res, nil
	}

	list, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	for _, e := range list {
		res = append(res, *mapEventToResponse(&e))
	}

	cacheSet(s.rc, key, res, TTLList)
	return res, nil
}

func (s *EventService) GetByID(id string) (*dto.EventResponse, error) {
	key := keyDetail("event", id)
	var res dto.EventResponse

	if cacheGet(s.rc, key, &res) {
		return &res, nil
	}

	e, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, errors.New("event tidak ditemukan")
	}

	resp := mapEventToResponse(e)
	cacheSet(s.rc, key, resp, TTLDetail)
	return resp, nil
}

func (s *EventService) GetBySlug(slug string) (*dto.EventResponse, error) {
	key := keyDetail("event", "slug:"+slug)
	var res dto.EventResponse

	if cacheGet(s.rc, key, &res) {
		return &res, nil
	}

	e, err := s.repo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, errors.New("event tidak ditemukan")
	}

	resp := mapEventToResponse(e)
	cacheSet(s.rc, key, resp, TTLDetail)
	return resp, nil
}

func (s *EventService) Update(id string, req dto.UpdateEventRequest) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("event tidak ditemukan")
	}

	if req.Tanggal != nil {
		if _, err := time.Parse(time.RFC3339, *req.Tanggal); err != nil {
			return errors.New("format tanggal tidak valid")
		}
	}

	// Check slug uniqueness if title is being updated
	if req.Judul != nil {
		slug := utils.Slugify(*req.Judul)
		existingWithSlug, err := s.repo.FindBySlug(slug)
		if err != nil {
			return err
		}
		if existingWithSlug != nil && existingWithSlug.ID != id {
			return errors.New("event dengan judul tersebut sudah ada")
		}
	}

	event := dto_event.EventUpdatedEvent{
		ID:        id,
		Request:   req,
		UpdatedAt: time.Now(),
	}

	if s.producer != nil {
		err = s.producer.PublishEventUpdated(context.Background(), event)
		if err != nil {
			return err
		}
	}

	// Invalidate caches
	cacheDelete(s.rc, keyList("event"))
	cacheDelete(s.rc, keyDetail("event", id))
	if existing.Slug != "" {
		cacheDelete(s.rc, keyDetail("event", "slug:"+existing.Slug))
	}
	return nil
}

func (s *EventService) Delete(id string) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("event tidak ditemukan")
	}

	event := dto_event.EventDeletedEvent{
		ID:        id,
		DeletedAt: time.Now(),
	}

	if s.producer != nil {
		err = s.producer.PublishEventDeleted(context.Background(), event)
		if err != nil {
			return err
		}
	}

	// Invalidate caches
	cacheDelete(s.rc, keyList("event"))
	cacheDelete(s.rc, keyDetail("event", id))
	if existing.Slug != "" {
		cacheDelete(s.rc, keyDetail("event", "slug:"+existing.Slug))
	}
	return nil
}

func (s *EventService) Register(eventID string, req dto.CreateEventRegistrationRequest) (*dto.EventRegistrationResponse, error) {
	event, err := s.repo.FindByID(eventID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, errors.New("event tidak ditemukan")
	}

	exists, err := s.repo.ExistsRegistrationByEventAndEmail(eventID, strings.TrimSpace(req.Email))
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email sudah terdaftar pada event ini")
	}

	payload := utils.EventRegistrationQRPayload{
		EventID:    event.ID,
		EventTitle: event.Judul,
		Nama:       req.Nama,
		Email:      req.Email,
		Perusahaan: req.Perusahaan,
		Jabatan:    req.Jabatan,
		NoHP:       req.NoHP,
		Sektor:     req.Sektor,
		QRToken:    uuid.NewString(),
	}

	reg := &models.EventRegistration{
		EventID:    eventID,
		Nama:       strings.TrimSpace(req.Nama),
		Email:      strings.TrimSpace(req.Email),
		Perusahaan: strings.TrimSpace(req.Perusahaan),
		Jabatan:    strings.TrimSpace(req.Jabatan),
		NoHP:       strings.TrimSpace(req.NoHP),
		Sektor:     strings.TrimSpace(req.Sektor),
		QRToken:    payload.QRToken,
	}

	rawPayload, err := payload.JSON()
	if err != nil {
		return nil, err
	}
	reg.QRPayload = rawPayload

	if err := s.repo.CreateRegistration(reg); err != nil {
		// Handle concurrent duplicate registration (MySQL UNIQUE constraint violation)
		var mysqlErr *mysqlDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, errors.New("email sudah terdaftar pada event ini")
		}
		return nil, err
	}

	payload.RegistrationID = reg.ID
	rawPayload, err = payload.JSON()
	if err != nil {
		return nil, err
	}
	reg.QRPayload = rawPayload

	if err := s.repo.UpdateRegistrationPayload(reg.ID, rawPayload); err != nil {
		return nil, err
	}

	// Invalidate caches
	cacheDelete(s.rc, keyList("event"))
	cacheDelete(s.rc, keyDetail("event", eventID))

	qrPNG, err := utils.GenerateQRCodePNG(rawPayload, 256)
	if err != nil {
		return nil, err
	}

	return &dto.EventRegistrationResponse{
		ID:           reg.ID,
		EventID:      reg.EventID,
		Nama:         reg.Nama,
		Email:        reg.Email,
		Perusahaan:   reg.Perusahaan,
		Jabatan:      reg.Jabatan,
		NoHP:         reg.NoHP,
		Sektor:       reg.Sektor,
		QRToken:      reg.QRToken,
		QRPayload:    rawPayload,
		QRCodeBase64: base64.StdEncoding.EncodeToString(qrPNG),
		DownloadURL:  fmt.Sprintf("/api/kegiatan/registrasi/%d/download", reg.ID),
		CreatedAt:    reg.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *EventService) DownloadRegistrationPDF(registrationID int64) ([]byte, string, error) {
	reg, err := s.repo.FindRegistrationByID(registrationID)
	if err != nil {
		return nil, "", err
	}
	if reg == nil {
		return nil, "", errors.New("registrasi event tidak ditemukan")
	}

	event, err := s.repo.FindByID(reg.EventID)
	if err != nil {
		return nil, "", err
	}
	if event == nil {
		return nil, "", errors.New("event tidak ditemukan")
	}

	qrPNG, err := utils.GenerateQRCodePNG(reg.QRPayload, 256)
	if err != nil {
		return nil, "", err
	}

	pdf, err := utils.GenerateEventRegistrationPDF(event.Judul, reg, qrPNG)
	if err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("registrasi-event-%s-%d.pdf", reg.EventID, reg.ID)
	return pdf, filename, nil
}

func mapEventToResponse(e *models.Event) *dto.EventResponse {
	status := "upcoming"
	if e.Tanggal.Before(time.Now()) {
		status = "past"
	}

	res := &dto.EventResponse{
		ID:        e.ID,
		Slug:      e.Slug,
		Judul:     e.Judul,
		Deskripsi: e.Deskripsi,
		Lokasi:    e.Lokasi,
		Tanggal:   e.Tanggal.Format(time.RFC3339),
		Status:    status,
		CreatedAt: e.CreatedAt.Format(time.RFC3339),
		UpdatedAt: e.UpdatedAt.Format(time.RFC3339),
	}

	return res
}

