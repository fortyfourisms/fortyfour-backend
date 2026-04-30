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

	"github.com/rollbar/rollbar-go"
)

type PertanyaanProteksiProducerInterface interface {
	PublishPertanyaanProteksiCreated(ctx context.Context, event interface{}) error
	PublishPertanyaanProteksiUpdated(ctx context.Context, event interface{}) error
	PublishPertanyaanProteksiDeleted(ctx context.Context, event interface{}) error
}

type PertanyaanProteksiService struct {
	repo     repository.PertanyaanProteksiRepositoryInterface
	producer PertanyaanProteksiProducerInterface
	cache    cache.RedisInterface
}

func NewPertanyaanProteksiService(
	repo repository.PertanyaanProteksiRepositoryInterface,
	producer PertanyaanProteksiProducerInterface,
	cache cache.RedisInterface,
) *PertanyaanProteksiService {
	return &PertanyaanProteksiService{
		repo:     repo,
		producer: producer,
		cache:    cache,
	}
}

func validateProteksiIndexField(value *string, fieldName string) error {
	if value == nil {
		return nil
	}
	normalized := utils.NormalizeInput(*value)
	*value = normalized
	if utils.ContainsSQLInjectionPattern(*value) {
		return errors.New(fieldName + " mengandung karakter yang tidak diizinkan")
	}
	if !utils.IsValidInput(*value) {
		return errors.New(fieldName + " hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&")
	}
	return nil
}

func (s *PertanyaanProteksiService) validateCreate(req *dto.CreatePertanyaanProteksiRequest) error {
	if req.SubKategoriID <= 0 {
		return errors.New("sub_kategori_id tidak valid")
	}

	if req.RuangLingkupID <= 0 {
		return errors.New("ruang_lingkup_id tidak valid")
	}

	req.PertanyaanProteksi = utils.NormalizeInput(req.PertanyaanProteksi)
	if req.PertanyaanProteksi == "" {
		return errors.New("pertanyaan_proteksi tidak boleh kosong")
	}
	if len(req.PertanyaanProteksi) < 3 {
		return errors.New("pertanyaan_proteksi minimal 3 karakter")
	}
	if utils.ContainsSQLInjectionPattern(req.PertanyaanProteksi) {
		return errors.New("pertanyaan_proteksi mengandung karakter yang tidak diizinkan")
	}
	if !utils.IsValidInput(req.PertanyaanProteksi) {
		return errors.New("pertanyaan_proteksi hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&")
	}

	if err := validateProteksiIndexField(req.Index0, "index0"); err != nil {
		return err
	}
	if err := validateProteksiIndexField(req.Index1, "index1"); err != nil {
		return err
	}
	if err := validateProteksiIndexField(req.Index2, "index2"); err != nil {
		return err
	}
	if err := validateProteksiIndexField(req.Index3, "index3"); err != nil {
		return err
	}
	if err := validateProteksiIndexField(req.Index4, "index4"); err != nil {
		return err
	}
	if err := validateProteksiIndexField(req.Index5, "index5"); err != nil {
		return err
	}

	return nil
}

func (s *PertanyaanProteksiService) validateUpdate(req *dto.UpdatePertanyaanProteksiRequest) error {
	if req.SubKategoriID != nil {
		if *req.SubKategoriID <= 0 {
			return errors.New("sub_kategori_id tidak valid")
		}
	}

	if req.RuangLingkupID != nil {
		if *req.RuangLingkupID <= 0 {
			return errors.New("ruang_lingkup_id tidak valid")
		}
	}

	if req.PertanyaanProteksi != nil {
		normalized := utils.NormalizeInput(*req.PertanyaanProteksi)
		req.PertanyaanProteksi = &normalized
		if *req.PertanyaanProteksi == "" {
			return errors.New("pertanyaan_proteksi tidak boleh kosong")
		}
		if len(*req.PertanyaanProteksi) < 3 {
			return errors.New("pertanyaan_proteksi minimal 3 karakter")
		}
		if utils.ContainsSQLInjectionPattern(*req.PertanyaanProteksi) {
			return errors.New("pertanyaan_proteksi mengandung karakter yang tidak diizinkan")
		}
		if !utils.IsValidInput(*req.PertanyaanProteksi) {
			return errors.New("pertanyaan_proteksi hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&")
		}
	}

	if err := validateProteksiIndexField(req.Index0, "index0"); err != nil {
		return err
	}
	if err := validateProteksiIndexField(req.Index1, "index1"); err != nil {
		return err
	}
	if err := validateProteksiIndexField(req.Index2, "index2"); err != nil {
		return err
	}
	if err := validateProteksiIndexField(req.Index3, "index3"); err != nil {
		return err
	}
	if err := validateProteksiIndexField(req.Index4, "index4"); err != nil {
		return err
	}
	if err := validateProteksiIndexField(req.Index5, "index5"); err != nil {
		return err
	}

	return nil
}

func (s *PertanyaanProteksiService) Create(req dto.CreatePertanyaanProteksiRequest) (*dto.PertanyaanProteksiResponse, error) {
	if err := s.validateCreate(&req); err != nil {
		return nil, err
	}

	subKategoriExists, err := s.repo.CheckSubKategoriExists(req.SubKategoriID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	if !subKategoriExists {
		return nil, errors.New("sub_kategori_id tidak ditemukan")
	}

	ruangLingkupExists, err := s.repo.CheckRuangLingkupExists(req.RuangLingkupID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	if !ruangLingkupExists {
		return nil, errors.New("ruang_lingkup_id tidak ditemukan")
	}

	if err := s.producer.PublishPertanyaanProteksiCreated(context.Background(), dto_event.PertanyaanProteksiCreatedEvent{
		Request:   req,
		CreatedAt: time.Now(),
	}); err != nil {
		rollbar.Error(err)
		return nil, err
	}

	// Invalidate cache
	_ = s.cache.Delete(cache.CacheKeyPertanyaanProteksi)

	return nil, nil
}

func (s *PertanyaanProteksiService) GetAll() ([]dto.PertanyaanProteksiResponse, error) {
	// Try to get from cache
	cachedData, err := s.cache.Get(cache.CacheKeyPertanyaanProteksi)
	if err == nil && cachedData != "" {
		var questions []dto.PertanyaanProteksiResponse
		if err := json.Unmarshal([]byte(cachedData), &questions); err == nil {
			return questions, nil
		}
	}

	// If not in cache or error, get from DB
	questions, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	// Store in cache (fire and forget for performance)
	go func() {
		jsonData, err := json.Marshal(questions)
		if err == nil {
			_ = s.cache.Set(cache.CacheKeyPertanyaanProteksi, string(jsonData), cache.DefaultCacheExpiration)
		}
	}()

	return questions, nil
}

func (s *PertanyaanProteksiService) GetByID(id int) (*dto.PertanyaanProteksiResponse, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}

	return data, nil
}

func (s *PertanyaanProteksiService) Update(id int, req dto.UpdatePertanyaanProteksiRequest) (*dto.PertanyaanProteksiResponse, error) {
	if err := s.validateUpdate(&req); err != nil {
		return nil, err
	}

	_, err := s.repo.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		if err == sql.ErrNoRows {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}

	if req.SubKategoriID != nil {
		subKategoriExists, err := s.repo.CheckSubKategoriExists(*req.SubKategoriID)
		if err != nil {
			rollbar.Error(err)
			return nil, err
		}
		if !subKategoriExists {
			return nil, errors.New("sub_kategori_id tidak ditemukan")
		}
	}

	if req.RuangLingkupID != nil {
		ruangLingkupExists, err := s.repo.CheckRuangLingkupExists(*req.RuangLingkupID)
		if err != nil {
			rollbar.Error(err)
			return nil, err
		}
		if !ruangLingkupExists {
			return nil, errors.New("ruang_lingkup_id tidak ditemukan")
		}
	}

	if err := s.producer.PublishPertanyaanProteksiUpdated(context.Background(), dto_event.PertanyaanProteksiUpdatedEvent{
		ID:        id,
		Request:   req,
		UpdatedAt: time.Now(),
	}); err != nil {
		rollbar.Error(err)
		return nil, err
	}

	// Invalidate cache
	_ = s.cache.Delete(cache.CacheKeyPertanyaanProteksi)

	return nil, nil
}

func (s *PertanyaanProteksiService) Delete(id int) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("data tidak ditemukan")
		}
		return err
	}

	err = s.producer.PublishPertanyaanProteksiDeleted(context.Background(), dto_event.PertanyaanProteksiDeletedEvent{
		ID:        id,
		DeletedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	// Invalidate cache
	_ = s.cache.Delete(cache.CacheKeyPertanyaanProteksi)

	return nil
}
