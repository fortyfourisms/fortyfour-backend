package services

import (
	"context"
	"errors"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/dto/dto_event"
	"fortyfour-backend/internal/models"
	internalRmq "fortyfour-backend/internal/rabbitmq"
	"fortyfour-backend/internal/repository"
	"time"
)

type EventServiceInterface interface {
	Create(req dto.CreateEventRequest) error
	GetAll() ([]dto.EventResponse, error)
	GetByID(id int64) (*dto.EventResponse, error)
	Update(id int64, req dto.UpdateEventRequest) error
	Delete(id int64) error
}

type EventService struct {
	repo     repository.EventRepositoryInterface
	producer *internalRmq.Producer
}

func NewEventService(repo repository.EventRepositoryInterface, producer *internalRmq.Producer) *EventService {
	return &EventService{
		repo:     repo,
		producer: producer,
	}
}

var _ EventServiceInterface = (*EventService)(nil)

func (s *EventService) Create(req dto.CreateEventRequest) error {
	if _, err := time.Parse(time.RFC3339, req.Tanggal); err != nil {
		return errors.New("format tanggal tidak valid (gunakan RFC3339, contoh: 2024-12-31T15:00:00Z)")
	}

	event := dto_event.EventCreatedEvent{
		Request:   req,
		CreatedAt: time.Now(),
	}

	return s.producer.PublishEventCreated(context.Background(), event)
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

func (s *EventService) Update(id int64, req dto.UpdateEventRequest) error {
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

	event := dto_event.EventUpdatedEvent{
		ID:        id,
		Request:   req,
		UpdatedAt: time.Now(),
	}

	return s.producer.PublishEventUpdated(context.Background(), event)
}

func (s *EventService) Delete(id int64) error {
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

	return s.producer.PublishEventDeleted(context.Background(), event)
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
		Lokasi:    e.Lokasi,
		Tanggal:   e.Tanggal.Format(time.RFC3339),
		Status:    status,
		CreatedAt: e.CreatedAt.Format(time.RFC3339),
		UpdatedAt: e.UpdatedAt.Format(time.RFC3339),
	}

	return res
}

