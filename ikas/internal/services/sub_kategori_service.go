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

type SubKategoriProducerInterface interface {
	PublishSubKategoriCreated(ctx context.Context, event interface{}) error
	PublishSubKategoriUpdated(ctx context.Context, event interface{}) error
	PublishSubKategoriDeleted(ctx context.Context, event interface{}) error
}

type SubKategoriService struct {
	repo     repository.SubKategoriRepositoryInterface
	producer SubKategoriProducerInterface
}

func NewSubKategoriService(repo repository.SubKategoriRepositoryInterface, producer SubKategoriProducerInterface) *SubKategoriService {
	return &SubKategoriService{
		repo:     repo,
		producer: producer,
	}
}

func (s *SubKategoriService) validateCreate(req *dto.CreateSubKategoriRequest) error {
	if req.KategoriID <= 0 {
		return errors.New("kategori_id tidak valid")
	}

	req.NamaSubKategori = utils.NormalizeInput(req.NamaSubKategori)

	// NOT NULL: tidak boleh kosong
	if req.NamaSubKategori == "" {
		return errors.New("nama_sub_kategori tidak boleh kosong")
	}

	// Min karakter
	if len(req.NamaSubKategori) < 3 {
		return errors.New("nama_sub_kategori minimal 3 karakter")
	}

	// Max karakter
	if len(req.NamaSubKategori) > 500 {
		return errors.New("nama_sub_kategori maksimal 500 karakter")
	}

	// Validasi SQL Injection pattern (blacklist)
	if utils.ContainsSQLInjectionPattern(req.NamaSubKategori) {
		return errors.New("nama_sub_kategori mengandung karakter yang tidak diizinkan")
	}

	// Validasi karakter yang diizinkan
	if !utils.IsValidInput(req.NamaSubKategori) {
		return errors.New("nama_sub_kategori hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&")
	}

	return nil
}

// Validasi untuk Update
func (s *SubKategoriService) validateUpdate(req *dto.UpdateSubKategoriRequest) error {
	// Validasi kategori_id jika dikirim
	if req.KategoriID != nil {
		if *req.KategoriID <= 0 {
			return errors.New("kategori_id tidak valid")
		}
	}

	// Jika field nama_sub_kategori dikirim (bukan nil), lakukan validasi
	if req.NamaSubKategori != nil {
		// Normalisasi: trim whitespace + hilangkan multiple spaces
		normalized := utils.NormalizeInput(*req.NamaSubKategori)
		req.NamaSubKategori = &normalized

		// NOT NULL: tidak boleh string kosong
		if *req.NamaSubKategori == "" {
			return errors.New("nama_sub_kategori tidak boleh kosong")
		}

		// Min karakter
		if len(*req.NamaSubKategori) < 3 {
			return errors.New("nama_sub_kategori minimal 3 karakter")
		}

		// Max karakter
		if len(*req.NamaSubKategori) > 500 {
			return errors.New("nama_sub_kategori maksimal 500 karakter")
		}

		// Validasi SQL Injection pattern (blacklist)
		if utils.ContainsSQLInjectionPattern(*req.NamaSubKategori) {
			return errors.New("nama_sub_kategori mengandung karakter yang tidak diizinkan")
		}

		// Validasi karakter yang diizinkan (whitelist - lebih ketat)
		if !utils.IsValidInput(*req.NamaSubKategori) {
			return errors.New("nama_sub_kategori hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&")
		}
	}

	return nil
}

func (s *SubKategoriService) Create(req dto.CreateSubKategoriRequest) (*dto.SubKategoriResponse, error) {
	// Validasi input
	if err := s.validateCreate(&req); err != nil {
		return nil, err
	}

	// Cek apakah kategori_id ada di tabel kategori (foreign key validation)
	kategoriExists, err := s.repo.CheckKategoriExists(req.KategoriID)
	if err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}
	if !kategoriExists {
		return nil, errors.New("kategori_id tidak ditemukan")
	}

	// Cek duplikasi data (case-insensitive, whitespace-trimmed) dalam kategori yang sama
	isDuplicate, err := s.repo.CheckDuplicateName(req.KategoriID, req.NamaSubKategori, 0)
	if err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}
	if isDuplicate {
		return nil, errors.New("nama_sub_kategori sudah ada dalam kategori ini")
	}

	// Opsi 1: tulis ke DB dulu
	newID, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	// Publish event ke RabbitMQ untuk notifikasi SSE (fire-and-forget)
	go func() {
		_ = s.producer.PublishSubKategoriCreated(context.Background(), dto_event.SubKategoriCreatedEvent{
			Request:   req,
			CreatedAt: time.Now(),
		})
	}()

	return s.repo.GetByID(int(newID))
}

func (s *SubKategoriService) GetAll() ([]dto.SubKategoriResponse, error) {
	return s.repo.GetAll()
}

func (s *SubKategoriService) GetByID(id int) (*dto.SubKategoriResponse, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}
	return data, nil
}

func (s *SubKategoriService) Update(id int, req dto.UpdateSubKategoriRequest) (*dto.SubKategoriResponse, error) {

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

	// Jika kategori_id diubah, cek apakah kategori baru exists
	if req.KategoriID != nil {
		kategoriExists, err := s.repo.CheckKategoriExists(*req.KategoriID)
		if err != nil {
			logger.Error(err, "operation failed")
			return nil, err
		}
		if !kategoriExists {
			return nil, errors.New("kategori_id tidak ditemukan")
		}
	}

	// Tentukan kategori_id yang akan digunakan untuk pengecekan duplikasi
	checkKategoriID := existing.KategoriID
	if req.KategoriID != nil {
		checkKategoriID = *req.KategoriID
	}

	// Cek duplikasi nama sub_kategori dalam kategori yang sama
	if req.NamaSubKategori != nil {
		isDuplicate, err := s.repo.CheckDuplicateName(checkKategoriID, *req.NamaSubKategori, id)
		if err != nil {
			logger.Error(err, "operation failed")
			return nil, err
		}
		if isDuplicate {
			return nil, errors.New("nama_sub_kategori sudah ada dalam kategori ini")
		}
	}

	// Opsi 1: tulis ke DB dulu
	if err := s.repo.Update(id, req); err != nil {
		return nil, err
	}

	// Publish event ke RabbitMQ untuk notifikasi SSE (fire-and-forget)
	go func() {
		_ = s.producer.PublishSubKategoriUpdated(context.Background(), dto_event.SubKategoriUpdatedEvent{
			ID:        id,
			Request:   req,
			UpdatedAt: time.Now(),
		})
	}()

	return nil, nil
}

func (s *SubKategoriService) Delete(id int) error {

	_, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("data tidak ditemukan")
		}
		return err
	}

	// Opsi 1: hapus dari DB dulu
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// Publish event ke RabbitMQ untuk notifikasi SSE (fire-and-forget)
	go func() {
		_ = s.producer.PublishSubKategoriDeleted(context.Background(), dto_event.SubKategoriDeletedEvent{
			ID:        id,
			DeletedAt: time.Now(),
		})
	}()

	return nil
}
