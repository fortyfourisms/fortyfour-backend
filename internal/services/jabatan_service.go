package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/dto/dto_event"
	"fortyfour-backend/internal/rabbitmq"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/pkg/cache"

	"github.com/google/uuid"
)

type JabatanService struct {
	repo     repository.JabatanRepositoryInterface
	rc       cache.RedisInterface
	producer *rabbitmq.Producer
}

func NewJabatanService(repo repository.JabatanRepositoryInterface, rc cache.RedisInterface, producer *rabbitmq.Producer) *JabatanService {
	return &JabatanService{
		repo:     repo,
		rc:       rc,
		producer: producer,
	}
}

func (s *JabatanService) Create(req dto.CreateJabatanRequest) (*dto.JabatanResponse, error) {
	if req.NamaJabatan == nil || strings.TrimSpace(*req.NamaJabatan) == "" {
		return nil, errors.New("nama_jabatan wajib diisi")
	}

	id := uuid.New().String()
	if err := s.repo.Create(req, id); err != nil {
		return nil, err
	}

	result, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	cacheSet(s.rc, keyDetail("jabatan", id), result, TTLDetail)
	cacheDelete(s.rc, keyList("jabatan"))

	if s.producer != nil {
		go func() {
			event := dto_event.JabatanCreatedEvent{
				ID:          result.ID,
				NamaJabatan: result.NamaJabatan,
				CreatedAt:   time.Now(),
			}
			s.producer.PublishJabatanCreated(context.Background(), event)
		}()
	}

	return result, nil
}

func (s *JabatanService) GetAll() ([]dto.JabatanResponse, error) {
	key := keyList("jabatan")
	var result []dto.JabatanResponse
	if cacheGet(s.rc, key, &result) {
		return result, nil
	}

	result, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	cacheSet(s.rc, key, result, TTLList)
	return result, nil
}

func (s *JabatanService) GetByID(id string) (*dto.JabatanResponse, error) {
	key := keyDetail("jabatan", id)
	var result dto.JabatanResponse
	if cacheGet(s.rc, key, &result) {
		return &result, nil
	}

	data, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	cacheSet(s.rc, key, data, TTLDetail)
	return data, nil
}

func (s *JabatanService) Update(id string, req dto.UpdateJabatanRequest) (*dto.JabatanResponse, error) {
	jabatan, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.NamaJabatan != nil {
		jabatan.NamaJabatan = *req.NamaJabatan
	}

	if err := s.repo.Update(id, *jabatan); err != nil {
		return nil, err
	}

	cacheDelete(s.rc, keyDetail("jabatan", id))
	cacheDelete(s.rc, keyList("jabatan"))

	if s.producer != nil {
		go func() {
			event := dto_event.JabatanUpdatedEvent{
				ID:          jabatan.ID,
				NamaJabatan: jabatan.NamaJabatan,
				UpdatedAt:   time.Now(),
			}
			s.producer.PublishJabatanUpdated(context.Background(), event)
		}()
	}

	return jabatan, nil
}

func (s *JabatanService) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	cacheDelete(s.rc, keyDetail("jabatan", id))
	cacheDelete(s.rc, keyList("jabatan"))

	if s.producer != nil {
		go func() {
			event := dto_event.JabatanDeletedEvent{
				ID:        id,
				DeletedAt: time.Now(),
			}
			s.producer.PublishJabatanDeleted(context.Background(), event)
		}()
	}

	return nil
}
