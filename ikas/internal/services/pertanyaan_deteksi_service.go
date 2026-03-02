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

type PertanyaanDeteksiService struct {
	repo repository.PertanyaanDeteksiRepositoryInterface
}

func NewPertanyaanDeteksiService(repo repository.PertanyaanDeteksiRepositoryInterface) *PertanyaanDeteksiService {
	return &PertanyaanDeteksiService{repo: repo}
}

func validateDeteksiIndexField(value *string, fieldName string) error {
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

func (s *PertanyaanDeteksiService) validateCreate(req *dto.CreatePertanyaanDeteksiRequest) error {
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

	req.PertanyaanDeteksi = utils.NormalizeInput(req.PertanyaanDeteksi)
	if req.PertanyaanDeteksi == "" {
		return errors.New("pertanyaan_deteksi tidak boleh kosong")
	}
	if len(req.PertanyaanDeteksi) < 3 {
		return errors.New("pertanyaan_deteksi minimal 3 karakter")
	}
	if utils.ContainsSQLInjectionPattern(req.PertanyaanDeteksi) {
		return errors.New("pertanyaan_deteksi mengandung karakter yang tidak diizinkan")
	}
	if !utils.IsValidInput(req.PertanyaanDeteksi) {
		return errors.New("pertanyaan_deteksi hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&")
	}

	if err := validateDeteksiIndexField(req.Index0, "index0"); err != nil {
		return err
	}
	if err := validateDeteksiIndexField(req.Index1, "index1"); err != nil {
		return err
	}
	if err := validateDeteksiIndexField(req.Index2, "index2"); err != nil {
		return err
	}
	if err := validateDeteksiIndexField(req.Index3, "index3"); err != nil {
		return err
	}
	if err := validateDeteksiIndexField(req.Index4, "index4"); err != nil {
		return err
	}
	if err := validateDeteksiIndexField(req.Index5, "index5"); err != nil {
		return err
	}

	return nil
}

func (s *PertanyaanDeteksiService) validateUpdate(req *dto.UpdatePertanyaanDeteksiRequest) error {
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

	if req.PertanyaanDeteksi != nil {
		normalized := utils.NormalizeInput(*req.PertanyaanDeteksi)
		req.PertanyaanDeteksi = &normalized
		if *req.PertanyaanDeteksi == "" {
			return errors.New("pertanyaan_deteksi tidak boleh kosong")
		}
		if len(*req.PertanyaanDeteksi) < 3 {
			return errors.New("pertanyaan_deteksi minimal 3 karakter")
		}
		if utils.ContainsSQLInjectionPattern(*req.PertanyaanDeteksi) {
			return errors.New("pertanyaan_deteksi mengandung karakter yang tidak diizinkan")
		}
		if !utils.IsValidInput(*req.PertanyaanDeteksi) {
			return errors.New("pertanyaan_deteksi hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&")
		}
	}

	if err := validateDeteksiIndexField(req.Index0, "index0"); err != nil {
		return err
	}
	if err := validateDeteksiIndexField(req.Index1, "index1"); err != nil {
		return err
	}
	if err := validateDeteksiIndexField(req.Index2, "index2"); err != nil {
		return err
	}
	if err := validateDeteksiIndexField(req.Index3, "index3"); err != nil {
		return err
	}
	if err := validateDeteksiIndexField(req.Index4, "index4"); err != nil {
		return err
	}
	if err := validateDeteksiIndexField(req.Index5, "index5"); err != nil {
		return err
	}

	return nil
}

func (s *PertanyaanDeteksiService) Create(req dto.CreatePertanyaanDeteksiRequest) (*dto.PertanyaanDeteksiResponse, error) {
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

func (s *PertanyaanDeteksiService) GetAll() ([]dto.PertanyaanDeteksiResponse, error) {
	return s.repo.GetAll()
}

func (s *PertanyaanDeteksiService) GetByID(id string) (*dto.PertanyaanDeteksiResponse, error) {
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

func (s *PertanyaanDeteksiService) Update(id string, req dto.UpdatePertanyaanDeteksiRequest) (*dto.PertanyaanDeteksiResponse, error) {
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

func (s *PertanyaanDeteksiService) Delete(id string) error {
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
