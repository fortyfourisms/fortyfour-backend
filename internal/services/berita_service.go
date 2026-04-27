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

type BeritaServiceInterface interface {
	Create(authorID string, req dto.CreateBeritaRequest) error
	GetAll() ([]dto.BeritaResponse, error)
	GetByID(id int64) (*dto.BeritaResponse, error)
	Update(id int64, req dto.UpdateBeritaRequest) error
	Delete(id int64) error
}

type BeritaService struct {
	repo     repository.BeritaRepositoryInterface
	producer *internalRmq.Producer
}

func NewBeritaService(repo repository.BeritaRepositoryInterface, producer *internalRmq.Producer) *BeritaService {
	return &BeritaService{
		repo:     repo,
		producer: producer,
	}
}

var _ BeritaServiceInterface = (*BeritaService)(nil)

func (s *BeritaService) Create(authorID string, req dto.CreateBeritaRequest) error {
	event := dto_event.BeritaCreatedEvent{
		AuthorID:  authorID,
		Request:   req,
		CreatedAt: time.Now(),
	}

	return s.producer.PublishBeritaCreated(context.Background(), event)
}

func (s *BeritaService) GetAll() ([]dto.BeritaResponse, error) {
	list, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var res []dto.BeritaResponse
	for _, b := range list {
		res = append(res, *mapBeritaToResponse(&b))
	}
	return res, nil
}

func (s *BeritaService) GetByID(id int64) (*dto.BeritaResponse, error) {
	b, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, errors.New("berita tidak ditemukan")
	}
	return mapBeritaToResponse(b), nil
}

func (s *BeritaService) Update(id int64, req dto.UpdateBeritaRequest) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("berita tidak ditemukan")
	}

	event := dto_event.BeritaUpdatedEvent{
		ID:        id,
		Request:   req,
		UpdatedAt: time.Now(),
	}

	return s.producer.PublishBeritaUpdated(context.Background(), event)
}

func (s *BeritaService) Delete(id int64) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("berita tidak ditemukan")
	}

	event := dto_event.BeritaDeletedEvent{
		ID:        id,
		DeletedAt: time.Now(),
	}

	return s.producer.PublishBeritaDeleted(context.Background(), event)
}


func mapBeritaToResponse(b *models.Berita) *dto.BeritaResponse {
	res := &dto.BeritaResponse{
		ID:        b.ID,
		Judul:     b.Judul,
		Deskripsi: b.Deskripsi,
		AuthorID:  b.AuthorID,
		CreatedAt: b.CreatedAt.Format(time.RFC3339),
		UpdatedAt: b.UpdatedAt.Format(time.RFC3339),
	}

	if b.Author != nil {
		res.Author = &dto.BeritaAuthor{
			ID:          b.Author.ID,
			Username:    b.Author.Username,
			DisplayName: b.Author.DisplayName,
		}
	}

	return res
}
