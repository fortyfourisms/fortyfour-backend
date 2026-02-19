package services

import (
	"database/sql"
	"errors"
	"ikas/internal/dto"
	"ikas/internal/repository"
	"ikas/internal/utils"

	"github.com/google/uuid"
	"github.com/rollbar/rollbar-go"
)

type KategoriService struct {
	repo repository.KategoriRepositoryInterface
}

func NewKategoriService(repo repository.KategoriRepositoryInterface) *KategoriService {
	return &KategoriService{repo: repo}
}

// Validasi untuk Create
func (s *KategoriService) validateCreate(req *dto.CreateKategoriRequest) error {
	// Validasi domain_id
	req.DomainID = utils.NormalizeInput(req.DomainID)

	if req.DomainID == "" {
		return errors.New("domain_id tidak boleh kosong")
	}

	// Validasi format UUID
	if !utils.IsValidUUID(req.DomainID) {
		return errors.New("format domain_id tidak valid")
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
		normalized := utils.NormalizeInput(*req.DomainID)
		req.DomainID = &normalized

		if *req.DomainID == "" {
			return errors.New("domain_id tidak boleh kosong")
		}

		// Validasi format UUID
		if !utils.IsValidUUID(*req.DomainID) {
			return errors.New("format domain_id tidak valid")
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
		rollbar.Error(err)
		return nil, err
	}
	if !domainExists {
		return nil, errors.New("domain_id tidak ditemukan")
	}

	// Cek duplikasi data (case-insensitive, whitespace-trimmed) dalam domain yang sama
	isDuplicate, err := s.repo.CheckDuplicateName(req.DomainID, req.NamaKategori, "")
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	if isDuplicate {
		return nil, errors.New("nama_kategori sudah ada dalam domain ini")
	}

	// Generate UUID
	newID := uuid.New().String()

	if err := s.repo.Create(req, newID); err != nil {
		rollbar.Error(err)
		return nil, err
	}

	// Ambil data yang baru dibuat
	resp, err := s.repo.GetByID(newID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return resp, nil
}

func (s *KategoriService) GetAll() ([]dto.KategoriResponse, error) {
	return s.repo.GetAll()
}

func (s *KategoriService) GetByID(id string) (*dto.KategoriResponse, error) {
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

func (s *KategoriService) Update(id string, req dto.UpdateKategoriRequest) (*dto.KategoriResponse, error) {
	// Validasi format UUID untuk mencegah SQL injection via ID
	if !utils.IsValidUUID(id) {
		return nil, errors.New("format ID tidak valid")
	}

	// Cek apakah data ada
	existing, err := s.repo.GetByID(id)
	if err != nil {
		rollbar.Error(err)
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
			rollbar.Error(err)
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
			rollbar.Error(err)
			return nil, err
		}
		if isDuplicate {
			return nil, errors.New("nama_kategori sudah ada dalam domain ini")
		}
	}

	if err := s.repo.Update(id, req); err != nil {
		rollbar.Error(err)
		return nil, err
	}

	updated, err := s.repo.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return updated, nil
}

func (s *KategoriService) Delete(id string) error {
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
