package services

import (
	"database/sql"
	"errors"
	"ikas/internal/dto"
	"ikas/internal/repository"
	"ikas/internal/utils"

	"github.com/rollbar/rollbar-go"
)

type PertanyaanProteksiService struct {
	repo repository.PertanyaanProteksiRepositoryInterface
}

func NewPertanyaanProteksiService(repo repository.PertanyaanProteksiRepositoryInterface) *PertanyaanProteksiService {
	return &PertanyaanProteksiService{repo: repo}
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

func (s *PertanyaanProteksiService) GetAll() ([]dto.PertanyaanProteksiResponse, error) {
	return s.repo.GetAll()
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

func (s *PertanyaanProteksiService) Delete(id int) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("data tidak ditemukan")
		}
		return err
	}

	return s.repo.Delete(id)
}
