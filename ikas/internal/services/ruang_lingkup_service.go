package services

import (
	"database/sql"
	"errors"
	"ikas/internal/dto"
	"ikas/internal/repository"
	"ikas/internal/utils"

	"fortyfour-backend/pkg/logger"
	"github.com/google/uuid"
)

type RuangLingkupService struct {
	repo repository.RuangLingkupRepositoryInterface
}

func NewRuangLingkupService(repo repository.RuangLingkupRepositoryInterface) *RuangLingkupService {
	return &RuangLingkupService{repo: repo}
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
	isDuplicate, err := s.repo.CheckDuplicateName(req.NamaRuangLingkup, "")
	if err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}
	if isDuplicate {
		return nil, errors.New("nama_ruang_lingkup sudah ada")
	}

	// Generate UUID
	newID := uuid.New().String()

	if err := s.repo.Create(req, newID); err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}

	// Ambil data yang baru dibuat
	resp, err := s.repo.GetByID(newID)
	if err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}

	return resp, nil
}

func (s *RuangLingkupService) GetAll() ([]dto.RuangLingkupResponse, error) {
	return s.repo.GetAll()
}

func (s *RuangLingkupService) GetByID(id string) (*dto.RuangLingkupResponse, error) {
	// Validasi format UUID untuk mencegah SQL injection via ID
	if !utils.IsValidUUID(id) {
		return nil, errors.New("format ID tidak valid")
	}

	data, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}
	return data, nil
}

func (s *RuangLingkupService) Update(id string, req dto.UpdateRuangLingkupRequest) (*dto.RuangLingkupResponse, error) {
	// Validasi format UUID untuk mencegah SQL injection via ID
	if !utils.IsValidUUID(id) {
		return nil, errors.New("format ID tidak valid")
	}

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

	// Update
	if err := s.repo.Update(id, req); err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}

	// Ambil data terbaru
	updated, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}

	return updated, nil
}

func (s *RuangLingkupService) Delete(id string) error {
	// Validasi format UUID untuk mencegah SQL injection via ID
	if !utils.IsValidUUID(id) {
		return errors.New("format ID tidak valid")
	}

	_, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("data tidak ditemukan")
		}
		return err
	}

	return s.repo.Delete(id)
}
