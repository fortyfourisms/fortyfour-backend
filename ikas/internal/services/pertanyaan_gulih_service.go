package services

import (
	"database/sql"
	"errors"
	"ikas/internal/dto"
	"ikas/internal/repository"
	"ikas/internal/utils"

	"github.com/rollbar/rollbar-go"
)

type PertanyaanGulihService struct {
	repo repository.PertanyaanGulihRepositoryInterface
}

func NewPertanyaanGulihService(repo repository.PertanyaanGulihRepositoryInterface) *PertanyaanGulihService {
	return &PertanyaanGulihService{repo: repo}
}

func validateGulihIndexField(value *string, fieldName string) error {
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

func (s *PertanyaanGulihService) validateCreate(req *dto.CreatePertanyaanGulihRequest) error {
	if req.SubKategoriID <= 0 {
		return errors.New("sub_kategori_id tidak valid")
	}

	if req.RuangLingkupID <= 0 {
		return errors.New("ruang_lingkup_id tidak valid")
	}

	req.PertanyaanGulih = utils.NormalizeInput(req.PertanyaanGulih)
	if req.PertanyaanGulih == "" {
		return errors.New("pertanyaan_gulih tidak boleh kosong")
	}
	if len(req.PertanyaanGulih) < 3 {
		return errors.New("pertanyaan_gulih minimal 3 karakter")
	}
	if utils.ContainsSQLInjectionPattern(req.PertanyaanGulih) {
		return errors.New("pertanyaan_gulih mengandung karakter yang tidak diizinkan")
	}
	if !utils.IsValidInput(req.PertanyaanGulih) {
		return errors.New("pertanyaan_gulih hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&")
	}

	if err := validateGulihIndexField(req.Index0, "index0"); err != nil {
		return err
	}
	if err := validateGulihIndexField(req.Index1, "index1"); err != nil {
		return err
	}
	if err := validateGulihIndexField(req.Index2, "index2"); err != nil {
		return err
	}
	if err := validateGulihIndexField(req.Index3, "index3"); err != nil {
		return err
	}
	if err := validateGulihIndexField(req.Index4, "index4"); err != nil {
		return err
	}
	if err := validateGulihIndexField(req.Index5, "index5"); err != nil {
		return err
	}

	return nil
}

func (s *PertanyaanGulihService) validateUpdate(req *dto.UpdatePertanyaanGulihRequest) error {
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

	if req.PertanyaanGulih != nil {
		normalized := utils.NormalizeInput(*req.PertanyaanGulih)
		req.PertanyaanGulih = &normalized
		if *req.PertanyaanGulih == "" {
			return errors.New("pertanyaan_gulih tidak boleh kosong")
		}
		if len(*req.PertanyaanGulih) < 3 {
			return errors.New("pertanyaan_gulih minimal 3 karakter")
		}
		if utils.ContainsSQLInjectionPattern(*req.PertanyaanGulih) {
			return errors.New("pertanyaan_gulih mengandung karakter yang tidak diizinkan")
		}
		if !utils.IsValidInput(*req.PertanyaanGulih) {
			return errors.New("pertanyaan_gulih hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&")
		}
	}

	if err := validateGulihIndexField(req.Index0, "index0"); err != nil {
		return err
	}
	if err := validateGulihIndexField(req.Index1, "index1"); err != nil {
		return err
	}
	if err := validateGulihIndexField(req.Index2, "index2"); err != nil {
		return err
	}
	if err := validateGulihIndexField(req.Index3, "index3"); err != nil {
		return err
	}
	if err := validateGulihIndexField(req.Index4, "index4"); err != nil {
		return err
	}
	if err := validateGulihIndexField(req.Index5, "index5"); err != nil {
		return err
	}

	return nil
}

func (s *PertanyaanGulihService) Create(req dto.CreatePertanyaanGulihRequest) (*dto.PertanyaanGulihResponse, error) {
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

	lastID, err := s.repo.Create(req)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	resp, err := s.repo.GetByID(int(lastID))
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return resp, nil
}

func (s *PertanyaanGulihService) GetAll() ([]dto.PertanyaanGulihResponse, error) {
	return s.repo.GetAll()
}

func (s *PertanyaanGulihService) GetByID(id int) (*dto.PertanyaanGulihResponse, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}

	return data, nil
}

func (s *PertanyaanGulihService) Update(id int, req dto.UpdatePertanyaanGulihRequest) (*dto.PertanyaanGulihResponse, error) {
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

func (s *PertanyaanGulihService) Delete(id int) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("data tidak ditemukan")
		}
		return err
	}

	return s.repo.Delete(id)
}
