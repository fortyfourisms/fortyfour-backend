package services

import (
	"database/sql"
	"errors"
	"ikas/internal/dto"
	"ikas/internal/repository"
	"ikas/internal/utils"

	"github.com/rollbar/rollbar-go"
)

type PertanyaanIdentifikasiService struct {
	repo repository.PertanyaanIdentifikasiRepositoryInterface
}

func NewPertanyaanIdentifikasiService(repo repository.PertanyaanIdentifikasiRepositoryInterface) *PertanyaanIdentifikasiService {
	return &PertanyaanIdentifikasiService{repo: repo}
}

func validateIndexField(value *string, fieldName string) error {
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

func (s *PertanyaanIdentifikasiService) validateCreate(req *dto.CreatePertanyaanIdentifikasiRequest) error {
	if req.SubKategoriID <= 0 {
		return errors.New("sub_kategori_id tidak valid")
	}

	if req.RuangLingkupID <= 0 {
		return errors.New("ruang_lingkup_id tidak valid")
	}

	req.PertanyaanIdentifikasi = utils.NormalizeInput(req.PertanyaanIdentifikasi)
	if req.PertanyaanIdentifikasi == "" {
		return errors.New("pertanyaan_identifikasi tidak boleh kosong")
	}

	if err := validateIndexField(req.Index0, "index0"); err != nil {
		return err
	}
	if err := validateIndexField(req.Index1, "index1"); err != nil {
		return err
	}
	if err := validateIndexField(req.Index2, "index2"); err != nil {
		return err
	}
	if err := validateIndexField(req.Index3, "index3"); err != nil {
		return err
	}
	if err := validateIndexField(req.Index4, "index4"); err != nil {
		return err
	}
	if err := validateIndexField(req.Index5, "index5"); err != nil {
		return err
	}

	return nil
}

func (s *PertanyaanIdentifikasiService) validateUpdate(req *dto.UpdatePertanyaanIdentifikasiRequest) error {
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

	if req.PertanyaanIdentifikasi != nil {
		*req.PertanyaanIdentifikasi = utils.NormalizeInput(*req.PertanyaanIdentifikasi)
		if *req.PertanyaanIdentifikasi == "" {
			return errors.New("pertanyaan_identifikasi tidak boleh kosong")
		}
	}

	if err := validateIndexField(req.Index0, "index0"); err != nil {
		return err
	}
	if err := validateIndexField(req.Index1, "index1"); err != nil {
		return err
	}
	if err := validateIndexField(req.Index2, "index2"); err != nil {
		return err
	}
	if err := validateIndexField(req.Index3, "index3"); err != nil {
		return err
	}
	if err := validateIndexField(req.Index4, "index4"); err != nil {
		return err
	}
	if err := validateIndexField(req.Index5, "index5"); err != nil {
		return err
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

func (s *PertanyaanIdentifikasiService) GetAll() ([]dto.PertanyaanIdentifikasiResponse, error) {
	return s.repo.GetAll()
}

func (s *PertanyaanIdentifikasiService) GetByID(id int) (*dto.PertanyaanIdentifikasiResponse, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}

	return data, nil
}

func (s *PertanyaanIdentifikasiService) Update(id int, req dto.UpdatePertanyaanIdentifikasiRequest) (*dto.PertanyaanIdentifikasiResponse, error) {
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

func (s *PertanyaanIdentifikasiService) Delete(id int) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("data tidak ditemukan")
		}
		return err
	}

	return s.repo.Delete(id)
}
