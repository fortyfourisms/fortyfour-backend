package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"ikas/internal/dto"
	"ikas/internal/dto/dto_event"
	"ikas/internal/repository"
	"ikas/internal/utils"
	"ikas/pkg/cache"
	"time"

	"fortyfour-backend/pkg/logger"
)

type DomainProducerInterface interface {
	PublishDomainCreated(ctx context.Context, event interface{}) error
	PublishDomainUpdated(ctx context.Context, event interface{}) error
	PublishDomainDeleted(ctx context.Context, event interface{}) error
}

type DomainService struct {
	repo     repository.DomainRepositoryInterface
	producer DomainProducerInterface
	cache    cache.RedisInterface
}

func NewDomainService(repo repository.DomainRepositoryInterface, producer DomainProducerInterface, cache cache.RedisInterface) *DomainService {
	return &DomainService{repo: repo, producer: producer, cache: cache}
}

func (s *DomainService) validateCreate(req *dto.CreateDomainRequest) error {
	req.NamaDomain = utils.NormalizeInput(req.NamaDomain)

	if req.NamaDomain == "" {
		return errors.New("nama_domain tidak boleh kosong")
	}
	if len(req.NamaDomain) < 3 {
		return errors.New("nama_domain minimal 3 karakter")
	}
	if len(req.NamaDomain) > 50 {
		return errors.New("nama_domain maksimal 50 karakter")
	}
	if utils.ContainsSQLInjectionPattern(req.NamaDomain) {
		return errors.New("nama_domain mengandung karakter tidak diizinkan")
	}
	if !utils.IsValidInput(req.NamaDomain) {
		return errors.New("nama_domain hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&")
	}
	return nil
}

func (s *DomainService) validateUpdate(req *dto.UpdateDomainRequest) error {
	if req.NamaDomain != nil {
		normalized := utils.NormalizeInput(*req.NamaDomain)
		req.NamaDomain = &normalized

		if *req.NamaDomain == "" {
			return errors.New("nama_domain tidak boleh kosong")
		}
		if len(*req.NamaDomain) < 3 {
			return errors.New("nama_domain minimal 3 karakter")
		}
		if len(*req.NamaDomain) > 50 {
			return errors.New("nama_domain maksimal 50 karakter")
		}
		if utils.ContainsSQLInjectionPattern(*req.NamaDomain) {
			return errors.New("nama_domain mengandung karakter tidak diizinkan")
		}
		if !utils.IsValidInput(*req.NamaDomain) {
			return errors.New("nama_domain hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&")
		}
	}
	return nil
}

func (s *DomainService) Create(req dto.CreateDomainRequest) (*dto.DomainResponse, error) {
	if err := s.validateCreate(&req); err != nil {
		return nil, err
	}

	dup, err := s.repo.CheckDuplicateName(req.NamaDomain, 0)
	if err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}
	if dup {
		return nil, errors.New("nama_domain sudah ada")
	}

	if err := s.producer.PublishDomainCreated(context.Background(), dto_event.DomainCreatedEvent{
		Request:   req,
		CreatedAt: time.Now(),
	}); err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}

	if s.cache != nil {
		s.cache.Delete(cache.CacheKeyDomain)
	}

	return nil, nil
}

func (s *DomainService) GetAll() ([]dto.DomainResponse, error) {
	if s.cache != nil {
		cachedData, err := s.cache.Get(cache.CacheKeyDomain)
		if err == nil && cachedData != "" {
			var data []dto.DomainResponse
			if err := json.Unmarshal([]byte(cachedData), &data); err == nil {
				return data, nil
			}
		}
	}

	data, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		go func(dataToCache []dto.DomainResponse) {
			jsonData, err := json.Marshal(dataToCache)
			if err == nil {
				_ = s.cache.Set(cache.CacheKeyDomain, string(jsonData), cache.DefaultCacheExpiration)
			}
		}(data)
	}

	return data, nil
}

func (s *DomainService) GetByID(id int) (*dto.DomainResponse, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}
	return data, nil
}

func (s *DomainService) Update(id int, req dto.UpdateDomainRequest) (*dto.DomainResponse, error) {
	if _, err := s.repo.GetByID(id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}

	if err := s.validateUpdate(&req); err != nil {
		return nil, err
	}

	if req.NamaDomain != nil {
		dup, err := s.repo.CheckDuplicateName(*req.NamaDomain, id)
		if err != nil {
			logger.Error(err, "operation failed")
			return nil, err
		}
		if dup {
			return nil, errors.New("nama_domain sudah ada")
		}
	}

	if err := s.producer.PublishDomainUpdated(context.Background(), dto_event.DomainUpdatedEvent{
		ID:        id,
		Request:   req,
		UpdatedAt: time.Now(),
	}); err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}

	if s.cache != nil {
		s.cache.Delete(cache.CacheKeyDomain)
	}

	return nil, nil
}

func (s *DomainService) Delete(id int) error {

	if _, err := s.repo.GetByID(id); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("data tidak ditemukan")
		}
		return err
	}

	hasKategori, err := s.repo.CheckHasKategori(id)
	if err != nil {
		return err
	}
	if hasKategori {
		return errors.New("tidak dapat menghapus domain karena masih memiliki kategori")
	}

	if err := s.producer.PublishDomainDeleted(context.Background(), dto_event.DomainDeletedEvent{
		ID:        id,
		DeletedAt: time.Now(),
	}); err != nil {
		logger.Error(err, "operation failed")
		return err
	}

	if s.cache != nil {
		s.cache.Delete(cache.CacheKeyDomain)
	}

	return nil
}
