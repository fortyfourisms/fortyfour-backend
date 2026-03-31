package services

import (
	"context"
	"database/sql"
	"errors"
	"ikas/internal/dto"
	"ikas/internal/dto/dto_event"
	"ikas/internal/rabbitmq"
	"ikas/internal/repository"
	"ikas/internal/utils"
	"time"

	"fortyfour-backend/pkg/logger"
)

type KategoriService struct {
	repo     repository.KategoriRepositoryInterface
	producer *rabbitmq.Producer
}

func NewKategoriService(repo repository.KategoriRepositoryInterface, producer *rabbitmq.Producer) *KategoriService {
	return &KategoriService{
		repo:     repo,
		producer: producer,
	}
}

// Validasi untuk Create
func (s *KategoriService) validateCreate(req *dto.CreateKategoriRequest) error {
	if req.DomainID <= 0 {
		return errors.New("domain_id tidak valid")
	}
	req.NamaKategori = utils.NormalizeInput(req.NamaKategori)

	if req.NamaKategori == "" {
		return errors.New("nama_kategori tidak boleh kosong")
	}

	if len(req.NamaKategori) < 3 {
		return errors.New("nama_kategori minimal 3 karakter")
	}

	if len(req.NamaKategori) > 500 {
		return errors.New("nama_kategori maksimal 500 karakter")
	}

	// Validasi SQL Injection pattern
	if utils.ContainsSQLInjectionPattern(req.NamaKategori) {
		return errors.New("nama_kategori mengandung karakter yang tidak diizinkan")
	}

	// Validasi karakter yang diizinkan
	if !utils.IsValidInput(req.NamaKategori) {
		return errors.New("nama_kategori hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&")
	}

	return nil
}

// Validasi untuk Update
func (s *KategoriService) validateUpdate(req *dto.UpdateKategoriRequest) error {
	if req.DomainID != nil {
		if *req.DomainID <= 0 {
			return errors.New("domain_id tidak valid")
		}
	}

	if req.NamaKategori != nil {
		normalized := utils.NormalizeInput(*req.NamaKategori)
		req.NamaKategori = &normalized

		if *req.NamaKategori == "" {
			return errors.New("nama_kategori tidak boleh kosong")
		}

		if len(*req.NamaKategori) < 3 {
			return errors.New("nama_kategori minimal 3 karakter")
		}

		if len(*req.NamaKategori) > 500 {
			return errors.New("nama_kategori maksimal 500 karakter")
		}

		// Validasi SQL Injection pattern
		if utils.ContainsSQLInjectionPattern(*req.NamaKategori) {
			return errors.New("nama_kategori mengandung karakter yang tidak diizinkan")
		}

		// Validasi karakter yang diizinkan
		if !utils.IsValidInput(*req.NamaKategori) {
			return errors.New("nama_kategori hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&")
		}
	}

	return nil
}

func (s *KategoriService) Create(req dto.CreateKategoriRequest) (*dto.KategoriResponse, error) {
	// Validasi input
	if err := s.validateCreate(&req); err != nil {
		return nil, err
	}

	// Cek apakah domain_id ada di tabel domain (foreign key validation)
	domainExists, err := s.repo.CheckDomainExists(req.DomainID)
	if err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}
	if !domainExists {
		return nil, errors.New("domain_id tidak ditemukan")
	}

	// Cek duplikasi data (case-insensitive, whitespace-trimmed) dalam domain yang sama
	isDuplicate, err := s.repo.CheckDuplicateName(req.DomainID, req.NamaKategori, 0)
	if err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}
	if isDuplicate {
		return nil, errors.New("nama_kategori sudah ada dalam domain ini")
	}

	if err := s.producer.PublishKategoriCreated(context.Background(), dto_event.KategoriCreatedEvent{
		Request:   req,
		CreatedAt: time.Now(),
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *KategoriService) GetAll() ([]dto.KategoriResponse, error) {
	return s.repo.GetAll()
}

func (s *KategoriService) GetByID(id int) (*dto.KategoriResponse, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}
	return data, nil
}

func (s *KategoriService) Update(id int, req dto.UpdateKategoriRequest) (*dto.KategoriResponse, error) {

	// Cek apakah data ada
	existing, err := s.repo.GetByID(id)
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

	// Jika domain_id diubah, cek apakah domain baru exists
	if req.DomainID != nil {
		domainExists, err := s.repo.CheckDomainExists(*req.DomainID)
		if err != nil {
			logger.Error(err, "operation failed")
			return nil, err
		}
		if !domainExists {
			return nil, errors.New("domain_id tidak ditemukan")
		}
	}

	checkDomainID := existing.DomainID
	if req.DomainID != nil {
		checkDomainID = *req.DomainID
	}

	if req.NamaKategori != nil {
		isDuplicate, err := s.repo.CheckDuplicateName(checkDomainID, *req.NamaKategori, id)
		if err != nil {
			logger.Error(err, "operation failed")
			return nil, err
		}
		if isDuplicate {
			return nil, errors.New("nama_kategori sudah ada dalam domain ini")
		}
	}

	if err := s.producer.PublishKategoriUpdated(context.Background(), dto_event.KategoriUpdatedEvent{
		ID:        id,
		Request:   req,
		UpdatedAt: time.Now(),
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *KategoriService) Delete(id int) error {

	_, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("data tidak ditemukan")
		}
		return err
	}

	return s.producer.PublishKategoriDeleted(context.Background(), dto_event.KategoriDeletedEvent{
		ID:        id,
		DeletedAt: time.Now(),
	})
}
