package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/dto/dto_event"
	"fortyfour-backend/internal/models"
	internalRmq "fortyfour-backend/internal/rabbitmq"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/pkg/cache"
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
	rc       cache.RedisInterface
}

func NewBeritaService(
	repo repository.BeritaRepositoryInterface,
	producer *internalRmq.Producer,
	rc cache.RedisInterface,
) *BeritaService {
	return &BeritaService{
		repo:     repo,
		producer: producer,
		rc:       rc,
	}
}

var _ BeritaServiceInterface = (*BeritaService)(nil)

func (s *BeritaService) Create(authorID string, req dto.CreateBeritaRequest) error {
	event := dto_event.BeritaCreatedEvent{
		AuthorID:  authorID,
		Request:   req,
		CreatedAt: time.Now(),
	}

	if s.producer != nil {
		err := s.producer.PublishBeritaCreated(context.Background(), event)
		if err != nil {
			return err
		}
	}

	// Invalidate list cache
	cacheDelete(s.rc, keyList("berita"))
	return nil
}

func (s *BeritaService) GetAll() ([]dto.BeritaResponse, error) {
	key := keyList("berita")
	var res []dto.BeritaResponse

	if cacheGet(s.rc, key, &res) {
		return res, nil
	}

	list, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	for _, b := range list {
		res = append(res, *mapBeritaToResponse(&b))
	}

	cacheSet(s.rc, key, res, TTLList)
	return res, nil
}

func (s *BeritaService) GetByID(id int64) (*dto.BeritaResponse, error) {
	key := keyDetail("berita", fmt.Sprintf("%d", id))
	var res dto.BeritaResponse

	if cacheGet(s.rc, key, &res) {
		return &res, nil
	}

	b, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, errors.New("berita tidak ditemukan")
	}

	resp := mapBeritaToResponse(b)
	cacheSet(s.rc, key, resp, TTLDetail)
	return resp, nil
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

	if s.producer != nil {
		err = s.producer.PublishBeritaUpdated(context.Background(), event)
		if err != nil {
			return err
		}
	}

	// Invalidate caches
	cacheDelete(s.rc, keyList("berita"))
	cacheDelete(s.rc, keyDetail("berita", fmt.Sprintf("%d", id)))
	return nil
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

	if s.producer != nil {
		err = s.producer.PublishBeritaDeleted(context.Background(), event)
		if err != nil {
			return err
		}
	}

	// Invalidate caches
	cacheDelete(s.rc, keyList("berita"))
	cacheDelete(s.rc, keyDetail("berita", fmt.Sprintf("%d", id)))
	return nil
}

func mapBeritaToResponse(b *models.Berita) *dto.BeritaResponse {
	var tags []string
	if b.Tags != "" {
		_ = json.Unmarshal([]byte(b.Tags), &tags)
	}
	if tags == nil {
		tags = []string{}
	}

	res := &dto.BeritaResponse{
		ID:        b.ID,
		Judul:     b.Judul,
		Deskripsi: b.Deskripsi,
		Tags:      tags,
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
