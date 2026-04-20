package services

import (
	"context"
	"database/sql"
	"errors"
	"ikas/internal/dto"
	"ikas/internal/dto/dto_event"
	"ikas/internal/repository"
	"ikas/internal/utils"
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
}

func NewRuangLingkupService(repo repository.RuangLingkupRepositoryInterface, producer RuangLingkupProducerInterface) *RuangLingkupService {
	return &RuangLingkupService{
		repo:     repo,
		producer: producer,
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

	newID, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	err = s.producer.PublishRuangLingkupCreated(context.Background(), dto_event.RuangLingkupCreatedEvent{
		Request:   req,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return s.repo.GetByID(int(newID))
}

func (s *RuangLingkupService) GetAll() ([]dto.RuangLingkupResponse, error) {
	return s.repo.GetAll()
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

	if err := s.repo.Update(id, req); err != nil {
		return nil, err
	}

	err = s.producer.PublishRuangLingkupUpdated(context.Background(), dto_event.RuangLingkupUpdatedEvent{
		ID:        id,
		Request:   req,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, err
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

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	err = s.producer.PublishRuangLingkupDeleted(context.Background(), dto_event.RuangLingkupDeletedEvent{
		ID:        id,
		DeletedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}
