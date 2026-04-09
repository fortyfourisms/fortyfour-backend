package services

import (
	"context"
	"database/sql"
	"errors"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/dto/dto_event"
	"fortyfour-backend/internal/rabbitmq"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/pkg/cache"
	"strings"
	"time"

	"github.com/google/uuid"
)

type PICService struct {
	repo     repository.PICRepositoryInterface
	rc       cache.RedisInterface
	producer *rabbitmq.Producer
}

func NewPICService(repo repository.PICRepositoryInterface, rc cache.RedisInterface, producer *rabbitmq.Producer) *PICService {
	return &PICService{
		repo:     repo,
		rc:       rc,
		producer: producer,
	}
}

func (s *PICService) Create(req dto.CreatePICRequest) (*dto.PICResponse, error) {
	if req.Nama == nil || strings.TrimSpace(*req.Nama) == "" {
		return nil, errors.New("nama wajib diisi")
	}

	if req.IDPerusahaan == nil || strings.TrimSpace(*req.IDPerusahaan) == "" {
		return nil, errors.New("id_perusahaan wajib diisi")
	}

	id := uuid.New().String()

	if err := s.repo.Create(req, id); err != nil {
		return nil, err
	}

	result, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	cacheSet(s.rc, keyDetail("pic", id), result, TTLDetail)
	cacheDelete(s.rc, keyList("pic"))
	if result.Perusahaan != nil {
		cacheDelete(s.rc, "pic_perusahaan_"+result.Perusahaan.ID)
	}

	if s.producer != nil {
		go func() {
			var perusahaanID string
			if result.Perusahaan != nil {
				perusahaanID = result.Perusahaan.ID
			}

			event := dto_event.PicCreatedEvent{
				ID:           result.ID,
				Nama:         result.Nama,
				Telepon:      result.Telepon,
				IDPerusahaan: perusahaanID,
				CreatedAt:    time.Now(),
			}
			s.producer.PublishPicCreated(context.Background(), event)
		}()
	}

	return result, nil
}

func (s *PICService) GetAll() ([]dto.PICResponse, error) {
	key := keyList("pic")
	var result []dto.PICResponse
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

func (s *PICService) GetByID(id string) (*dto.PICResponse, error) {
	key := keyDetail("pic", id)
	var result dto.PICResponse
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

func (s *PICService) GetByPerusahaan(idPerusahaan string) ([]dto.PICResponse, error) {
	key := "pic_perusahaan_" + idPerusahaan
	var result []dto.PICResponse
	if cacheGet(s.rc, key, &result) {
		return result, nil
	}

	result, err := s.repo.GetByPerusahaan(idPerusahaan)
	if err != nil {
		return nil, err
	}

	cacheSet(s.rc, key, result, TTLList)
	return result, nil
}

func (s *PICService) Update(id string, req dto.UpdatePICRequest) (*dto.PICResponse, error) {
	if err := s.repo.Update(id, req); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("pic not found")
		}
		return nil, err
	}

	updated, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	cacheDelete(s.rc, keyDetail("pic", id))
	cacheDelete(s.rc, keyList("pic"))
	if updated.Perusahaan != nil {
		cacheDelete(s.rc, "pic_perusahaan_"+updated.Perusahaan.ID)
	}

	if s.producer != nil {
		go func() {
			var perusahaanID string
			if updated.Perusahaan != nil {
				perusahaanID = updated.Perusahaan.ID
			}

			event := dto_event.PicUpdatedEvent{
				ID:           updated.ID,
				Nama:         updated.Nama,
				Telepon:      updated.Telepon,
				IDPerusahaan: perusahaanID,
				UpdatedAt:    time.Now(),
			}
			s.producer.PublishPicUpdated(context.Background(), event)
		}()
	}
	return updated, nil
}

func (s *PICService) Delete(id string) error {
	// Ambil data dulu untuk invalidate cache per perusahaan
	existing, _ := s.repo.GetByID(id)

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	cacheDelete(s.rc, keyDetail("pic", id))
	cacheDelete(s.rc, keyList("pic"))
	if existing != nil && existing.Perusahaan != nil {
		cacheDelete(s.rc, "pic_perusahaan_"+existing.Perusahaan.ID)
	}

	if s.producer != nil {
		go func() {
			event := dto_event.PicDeletedEvent{
				ID:        id,
				DeletedAt: time.Now(),
			}
			s.producer.PublishPicDeleted(context.Background(), event)
		}()
	}

	return nil
}
