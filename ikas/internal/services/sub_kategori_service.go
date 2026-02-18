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

type SubKategoriService struct {
	repo repository.SubKategoriRepositoryInterface
}

func NewSubKategoriService(repo repository.SubKategoriRepositoryInterface) *SubKategoriService {
	return &SubKategoriService{repo: repo}
}

func (s *SubKategoriService) validateCreate(req *dto.CreateSubKategoriRequest) error {
	req.KategoriID = utils.NormalizeInput(req.KategoriID)

	if req.KategoriID == "" {
		return errors.New("kategori_id tidak boleh kosong")
	}

	// Validasi format UUID
	if !utils.IsValidUUID(req.KategoriID) {
		return errors.New("format kategori_id tidak valid")
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
		normalized := utils.NormalizeInput(*req.KategoriID)
		req.KategoriID = &normalized

		// NOT NULL: tidak boleh string kosong
		if *req.KategoriID == "" {
			return errors.New("kategori_id tidak boleh kosong")
		}

		// Validasi format UUID
		if !utils.IsValidUUID(*req.KategoriID) {
			return errors.New("format kategori_id tidak valid")
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
		rollbar.Error(err)
		return nil, err
	}
	if !kategoriExists {
		return nil, errors.New("kategori_id tidak ditemukan")
	}

	// Cek duplikasi data (case-insensitive, whitespace-trimmed) dalam kategori yang sama
	isDuplicate, err := s.repo.CheckDuplicateName(req.KategoriID, req.NamaSubKategori, "")
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	if isDuplicate {
		return nil, errors.New("nama_sub_kategori sudah ada dalam kategori ini")
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

func (s *SubKategoriService) GetAll() ([]dto.SubKategoriResponse, error) {
	return s.repo.GetAll()
}

func (s *SubKategoriService) GetByID(id string) (*dto.SubKategoriResponse, error) {
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

func (s *SubKategoriService) Update(id string, req dto.UpdateSubKategoriRequest) (*dto.SubKategoriResponse, error) {
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

	// Jika kategori_id diubah, cek apakah kategori baru exists
	if req.KategoriID != nil {
		kategoriExists, err := s.repo.CheckKategoriExists(*req.KategoriID)
		if err != nil {
			rollbar.Error(err)
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
			rollbar.Error(err)
			return nil, err
		}
		if isDuplicate {
			return nil, errors.New("nama_sub_kategori sudah ada dalam kategori ini")
		}
	}

	// Update
	if err := s.repo.Update(id, req); err != nil {
		rollbar.Error(err)
		return nil, err
	}

	// Ambil data terbaru
	updated, err := s.repo.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return updated, nil
}

func (s *SubKategoriService) Delete(id string) error {
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
