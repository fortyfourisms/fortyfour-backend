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

type RuangLingkupProducerInterface interface {
	PublishRuangLingkupCreated(ctx context.Context, event interface{}) error
	PublishRuangLingkupUpdated(ctx context.Context, event interface{}) error
	PublishRuangLingkupDeleted(ctx context.Context, event interface{}) error
}

type RuangLingkupService struct {
	repo     repository.RuangLingkupRepositoryInterface
	producer RuangLingkupProducerInterface
	cache    cache.RedisInterface
}

func NewRuangLingkupService(
	repo repository.RuangLingkupRepositoryInterface,
	producer RuangLingkupProducerInterface,
	cache cache.RedisInterface,
) *RuangLingkupService {
	return &RuangLingkupService{
		repo:     repo,
		producer: producer,
		cache:    cache,
	}
}

// Validasi untuk Create
func (s *RuangLingkupService) validateCreate(req *dto.CreateRuangLingkupRequest) error {
	// Normalisasi: trim whitespace + hilangkan multiple spaces
	req.NamaRuangLingkup = utils.NormalizeInput(req.NamaRuangLingkup)

	// NOT NULL: tidak boleh kosong
	if req.NamaRuangLingkup == "" {
		return errors.New("nama_ruang_lingkup tidak boleh kosong")
	}

	// Min karakter
	if len(req.NamaRuangLingkup) < 3 {
		return errors.New("nama_ruang_lingkup minimal 3 karakter")
	}

	// Max karakter
	if len(req.NamaRuangLingkup) > 50 {
		return errors.New("nama_ruang_lingkup maksimal 50 karakter")
	}

	// Validasi SQL Injection pattern (blacklist)
	if utils.ContainsSQLInjectionPattern(req.NamaRuangLingkup) {
		return errors.New("nama_ruang_lingkup mengandung karakter yang tidak diizinkan")
	}

	// Validasi karakter yang diizinkan
	if !utils.IsValidInput(req.NamaRuangLingkup) {
		return errors.New("nama_ruang_lingkup hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&")
	}

	return nil
}

// Validasi untuk Update
func (s *RuangLingkupService) validateUpdate(req *dto.UpdateRuangLingkupRequest) error {
	// Jika field dikirim (bukan nil), lakukan validasi
	if req.NamaRuangLingkup != nil {
		// Normalisasi: trim whitespace + hilangkan multiple spaces
		normalized := utils.NormalizeInput(*req.NamaRuangLingkup)
		req.NamaRuangLingkup = &normalized

		// NOT NULL: tidak boleh string kosong
		if *req.NamaRuangLingkup == "" {
			return errors.New("nama_ruang_lingkup tidak boleh kosong")
		}

		// Min karakter
		if len(*req.NamaRuangLingkup) < 3 {
			return errors.New("nama_ruang_lingkup minimal 3 karakter")
		}

		// Max karakter
		if len(*req.NamaRuangLingkup) > 50 {
			return errors.New("nama_ruang_lingkup maksimal 50 karakter")
		}

		// Validasi SQL Injection pattern (blacklist)
		if utils.ContainsSQLInjectionPattern(*req.NamaRuangLingkup) {
			return errors.New("nama_ruang_lingkup mengandung karakter yang tidak diizinkan")
		}

		// Validasi karakter yang diizinkan (whitelist - lebih ketat)
		if !utils.IsValidInput(*req.NamaRuangLingkup) {
			return errors.New("nama_ruang_lingkup hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&")
		}
	}

	return nil
}

func (s *RuangLingkupService) Create(req dto.CreateRuangLingkupRequest) (*dto.RuangLingkupResponse, error) {
	// Validasi input
	if err := s.validateCreate(&req); err != nil {
		return nil, err
	}

	// Cek duplikasi data (case-insensitive, whitespace-trimmed)
	isDuplicate, err := s.repo.CheckDuplicateName(req.NamaRuangLingkup, 0)
	if err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}
	if isDuplicate {
		return nil, errors.New("nama_ruang_lingkup sudah ada")
	}

	err = s.producer.PublishRuangLingkupCreated(context.Background(), dto_event.RuangLingkupCreatedEvent{
		Request:   req,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		s.cache.Delete(cache.CacheKeyRuangLingkup)
	}

	return nil, nil
}

func (s *RuangLingkupService) GetAll() ([]dto.RuangLingkupResponse, error) {
	if s.cache != nil {
		cachedData, err := s.cache.Get(cache.CacheKeyRuangLingkup)
		if err == nil && cachedData != "" {
			var data []dto.RuangLingkupResponse
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
		go func(dataToCache []dto.RuangLingkupResponse) {
			jsonData, err := json.Marshal(dataToCache)
			if err == nil {
				_ = s.cache.Set(cache.CacheKeyRuangLingkup, string(jsonData), cache.DefaultCacheExpiration)
			}
		}(data)
	}

	return data, nil
}

func (s *RuangLingkupService) GetByID(id int) (*dto.RuangLingkupResponse, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}
	return data, nil
}

func (s *RuangLingkupService) Update(id int, req dto.UpdateRuangLingkupRequest) (*dto.RuangLingkupResponse, error) {

	_, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error(err, "operation failed")
		if err == sql.ErrNoRows {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}

	// Validasi input
	if err := s.validateUpdate(&req); err != nil {
		return nil, err
	}

	// Cek duplikasi nama
	if req.NamaRuangLingkup != nil {
		isDuplicate, err := s.repo.CheckDuplicateName(*req.NamaRuangLingkup, id)
		if err != nil {
			logger.Error(err, "operation failed")
			return nil, err
		}
		if isDuplicate {
			return nil, errors.New("nama_ruang_lingkup sudah ada")
		}
	}

	err = s.producer.PublishRuangLingkupUpdated(context.Background(), dto_event.RuangLingkupUpdatedEvent{
		ID:        id,
		Request:   req,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		s.cache.Delete(cache.CacheKeyRuangLingkup)
	}

	return nil, nil
}

func (s *RuangLingkupService) Delete(id int) error {

	_, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("data tidak ditemukan")
		}
		return err
	}

	hasPertanyaan, err := s.repo.CheckHasPertanyaan(id)
	if err != nil {
		return err
	}
	if hasPertanyaan {
		return errors.New("tidak dapat menghapus ruang lingkup karena masih digunakan oleh pertanyaan")
	}

	err = s.producer.PublishRuangLingkupDeleted(context.Background(), dto_event.RuangLingkupDeletedEvent{
		ID:        id,
		DeletedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	if s.cache != nil {
		s.cache.Delete(cache.CacheKeyRuangLingkup)
	}

	return nil
}
