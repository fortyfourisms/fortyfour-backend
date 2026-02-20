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

type PertanyaanIdentifikasiService struct {
	repo repository.PertanyaanIdentifikasiRepositoryInterface
}

func NewPertanyaanIdentifikasiService(repo repository.PertanyaanIdentifikasiRepositoryInterface) *PertanyaanIdentifikasiService {
	return &PertanyaanIdentifikasiService{repo: repo}
}

func (s *PertanyaanIdentifikasiService) validateCreate(req *dto.CreatePertanyaanIdentifikasiRequest) error {
	req.SubKategoriID = utils.NormalizeInput(req.SubKategoriID)
	if req.SubKategoriID == "" {
		return errors.New("sub_kategori_id tidak boleh kosong")
	}
	if !utils.IsValidUUID(req.SubKategoriID) {
		return errors.New("format sub_kategori_id tidak valid")
	}

	req.RuangLingkupID = utils.NormalizeInput(req.RuangLingkupID)
	if req.RuangLingkupID == "" {
		return errors.New("ruang_lingkup_id tidak boleh kosong")
	}
	if !utils.IsValidUUID(req.RuangLingkupID) {
		return errors.New("format ruang_lingkup_id tidak valid")
	}

	req.PertanyaanIdentifikasi = utils.NormalizeInput(req.PertanyaanIdentifikasi)
	if req.PertanyaanIdentifikasi == "" {
		return errors.New("pertanyaan_identifikasi tidak boleh kosong")
	}
	if len(req.PertanyaanIdentifikasi) < 3 {
		return errors.New("pertanyaan_identifikasi minimal 3 karakter")
	}

	return nil
}

func (s *PertanyaanIdentifikasiService) validateUpdate(req *dto.UpdatePertanyaanIdentifikasiRequest) error {
	if req.SubKategoriID != nil {
		normalized := utils.NormalizeInput(*req.SubKategoriID)
		req.SubKategoriID = &normalized
		if *req.SubKategoriID == "" {
			return errors.New("sub_kategori_id tidak boleh kosong")
		}
		if !utils.IsValidUUID(*req.SubKategoriID) {
			return errors.New("format sub_kategori_id tidak valid")
		}
	}

	if req.RuangLingkupID != nil {
		normalized := utils.NormalizeInput(*req.RuangLingkupID)
		req.RuangLingkupID = &normalized
		if *req.RuangLingkupID == "" {
			return errors.New("ruang_lingkup_id tidak boleh kosong")
		}
		if !utils.IsValidUUID(*req.RuangLingkupID) {
			return errors.New("format ruang_lingkup_id tidak valid")
		}
	}

	if req.PertanyaanIdentifikasi != nil {
		normalized := utils.NormalizeInput(*req.PertanyaanIdentifikasi)
		req.PertanyaanIdentifikasi = &normalized
		if *req.PertanyaanIdentifikasi == "" {
			return errors.New("pertanyaan_identifikasi tidak boleh kosong")
		}
		if len(*req.PertanyaanIdentifikasi) < 3 {
			return errors.New("pertanyaan_identifikasi minimal 3 karakter")
		}
	}

	return nil
}

func (s *PertanyaanIdentifikasiService) Create(req dto.CreatePertanyaanIdentifikasiRequest) (*dto.PertanyaanIdentifikasiResponse, error) {
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

	newID := uuid.New().String()

	if err := s.repo.Create(req, newID); err != nil {
		rollbar.Error(err)
		return nil, err
	}

	resp, err := s.repo.GetByID(newID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return resp, nil
}

func (s *PertanyaanIdentifikasiService) GetAll() ([]dto.PertanyaanIdentifikasiResponse, error) {
	return s.repo.GetAll()
}

func (s *PertanyaanIdentifikasiService) GetByID(id string) (*dto.PertanyaanIdentifikasiResponse, error) {
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

func (s *PertanyaanIdentifikasiService) Update(id string, req dto.UpdatePertanyaanIdentifikasiRequest) (*dto.PertanyaanIdentifikasiResponse, error) {
	if !utils.IsValidUUID(id) {
		return nil, errors.New("format ID tidak valid")
	}

	_, err := s.repo.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		if err == sql.ErrNoRows {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}

	if err := s.validateUpdate(&req); err != nil {
		return nil, err
	}

	if req.SubKategoriID != nil {
		exists, err := s.repo.CheckSubKategoriExists(*req.SubKategoriID)
		if err != nil {
			rollbar.Error(err)
			return nil, err
		}
		if !exists {
			return nil, errors.New("sub_kategori_id tidak ditemukan")
		}
	}

	if req.RuangLingkupID != nil {
		exists, err := s.repo.CheckRuangLingkupExists(*req.RuangLingkupID)
		if err != nil {
			rollbar.Error(err)
			return nil, err
		}
		if !exists {
			return nil, errors.New("ruang_lingkup_id tidak ditemukan")
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

func (s *PertanyaanIdentifikasiService) Delete(id string) error {
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
