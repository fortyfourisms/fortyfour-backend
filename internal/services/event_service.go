package services

import (
	"errors"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
	"time"
)

type EventServiceInterface interface {
	Create(req dto.CreateEventRequest) (*dto.EventResponse, error)
	GetAll() ([]dto.EventResponse, error)
	GetByID(id int64) (*dto.EventResponse, error)
	Update(id int64, req dto.UpdateEventRequest) (*dto.EventResponse, error)
	Delete(id int64) error
}

type EventService struct {
	repo repository.EventRepositoryInterface
}

func NewEventService(repo repository.EventRepositoryInterface) *EventService {
	return &EventService{repo: repo}
}

var _ EventServiceInterface = (*EventService)(nil)

func (s *EventService) Create(req dto.CreateEventRequest) (*dto.EventResponse, error) {
	tanggal, err := time.Parse(time.RFC3339, req.Tanggal)
	if err != nil {
		return nil, errors.New("format tanggal tidak valid (gunakan RFC3339, contoh: 2024-12-31T15:00:00Z)")
	}

	event := &models.Event{
		Judul:     req.Judul,
		Deskripsi: req.Deskripsi,
		Tanggal:   tanggal,
		Lokasi:    req.Lokasi,
	}

	if err := s.repo.Create(event); err != nil {
		return nil, err
	}

	saved, err := s.repo.FindByID(event.ID)
	if err != nil {
		return nil, err
	}

	return mapEventToResponse(saved), nil
}

func (s *EventService) GetAll() ([]dto.EventResponse, error) {
	list, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var res []dto.EventResponse
	for _, e := range list {
		res = append(res, *mapEventToResponse(&e))
	}
	return res, nil
}

func (s *EventService) GetByID(id int64) (*dto.EventResponse, error) {
	e, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, errors.New("event tidak ditemukan")
	}
	return mapEventToResponse(e), nil
}

func (s *EventService) Update(id int64, req dto.UpdateEventRequest) (*dto.EventResponse, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("event tidak ditemukan")
	}

	if req.Judul != nil {
		existing.Judul = *req.Judul
	}
	if req.Deskripsi != nil {
		existing.Deskripsi = *req.Deskripsi
	}
	if req.Tanggal != nil {
		tanggal, err := time.Parse(time.RFC3339, *req.Tanggal)
		if err != nil {
			return nil, errors.New("format tanggal tidak valid")
		}
		existing.Tanggal = tanggal
	}
	if req.Lokasi != nil {
		existing.Lokasi = *req.Lokasi
	}

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return mapEventToResponse(existing), nil
}

func (s *EventService) Delete(id int64) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("event tidak ditemukan")
	}

	return s.repo.Delete(id)
}

func mapEventToResponse(e *models.Event) *dto.EventResponse {
	status := "upcoming"
	if e.Tanggal.Before(time.Now()) {
		status = "past"
	}

	res := &dto.EventResponse{
		ID:        e.ID,
		Judul:     e.Judul,
		Deskripsi: e.Deskripsi,
		Tanggal:   e.Tanggal.Format(time.RFC3339),
		Lokasi:    e.Lokasi,
		Status:    status,
		CreatedAt: e.CreatedAt.Format(time.RFC3339),
		UpdatedAt: e.UpdatedAt.Format(time.RFC3339),
	}

	return res
}
