package services

import (
	"database/sql"
	"errors"
	"ikas/internal/dto"
	"ikas/internal/repository"
	"ikas/internal/utils"

	"github.com/rollbar/rollbar-go"
)

type JawabanDeteksiService struct {
	repo repository.JawabanDeteksiRepositoryInterface
}

func NewJawabanDeteksiService(repo repository.JawabanDeteksiRepositoryInterface) *JawabanDeteksiService {
	return &JawabanDeteksiService{repo: repo}
}

var validValidasiDeteksi = map[string]bool{"yes": true, "no": true}

func (s *JawabanDeteksiService) validateCreate(req *dto.CreateJawabanDeteksiRequest) error {
	if req.PertanyaanDeteksiID <= 0 {
		return errors.New("pertanyaan_deteksi_id tidak valid")
	}

	req.PerusahaanID = utils.NormalizeInput(req.PerusahaanID)
	if req.PerusahaanID == "" {
		return errors.New("perusahaan_id tidak boleh kosong")
	}
	if !utils.IsValidUUID(req.PerusahaanID) {
		return errors.New("format perusahaan_id tidak valid")
	}

	if req.JawabanDeteksi != nil && (*req.JawabanDeteksi < 0 || *req.JawabanDeteksi > 5) {
		return errors.New("jawaban_deteksi harus bernilai antara 0 sampai 5, atau null untuk N/A")
	}

	if req.Validasi != nil {
		if req.Evidence == nil || utils.NormalizeInput(*req.Evidence) == "" {
			return errors.New("validasi hanya boleh diisi jika evidence ada")
		}
		if !validValidasiDeteksi[*req.Validasi] {
			return errors.New("validasi hanya boleh berisi 'yes' atau 'no'")
		}
	}

	return nil
}

func (s *JawabanDeteksiService) validateUpdate(req *dto.UpdateJawabanDeteksiRequest, existingEvidence *string) error {
	if req.JawabanDeteksi != nil && (*req.JawabanDeteksi < 0 || *req.JawabanDeteksi > 5) {
		return errors.New("jawaban_deteksi harus bernilai antara 0 sampai 5, atau null untuk N/A")
	}

	if req.Validasi != nil {
		if !validValidasiDeteksi[*req.Validasi] {
			return errors.New("validasi hanya boleh berisi 'yes' atau 'no'")
		}
		effectiveEvidence := existingEvidence
		if req.Evidence != nil {
			effectiveEvidence = req.Evidence
		}
		if effectiveEvidence == nil || utils.NormalizeInput(*effectiveEvidence) == "" {
			return errors.New("validasi hanya boleh diisi jika evidence ada")
		}
	}

	return nil
}

func (s *JawabanDeteksiService) Create(req dto.CreateJawabanDeteksiRequest) (*dto.JawabanDeteksiResponse, error) {
	if err := s.validateCreate(&req); err != nil {
		return nil, err
	}

	pertanyaanExists, err := s.repo.CheckPertanyaanExists(req.PertanyaanDeteksiID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	if !pertanyaanExists {
		return nil, errors.New("pertanyaan_deteksi_id tidak ditemukan")
	}

	perusahaanExists, err := s.repo.CheckPerusahaanExists(req.PerusahaanID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	if !perusahaanExists {
		return nil, errors.New("perusahaan_id tidak ditemukan")
	}

	isDuplicate, err := s.repo.CheckDuplicate(req.PerusahaanID, req.PertanyaanDeteksiID, 0)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	if isDuplicate {
		return nil, errors.New("jawaban untuk pertanyaan ini sudah ada untuk perusahaan tersebut")
	}

	lastID, err := s.repo.Create(req)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	if err := s.repo.RecalculateDeteksi(req.PerusahaanID); err != nil {
		rollbar.Error(err)
	}

	resp, err := s.repo.GetByID(int(lastID))
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return resp, nil
}

func (s *JawabanDeteksiService) GetAll() ([]dto.JawabanDeteksiResponse, error) {
	return s.repo.GetAll()
}

func (s *JawabanDeteksiService) GetByID(id int) (*dto.JawabanDeteksiResponse, error) {
	if id <= 0 {
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

func (s *JawabanDeteksiService) GetByPerusahaan(perusahaanID string) ([]dto.JawabanDeteksiResponse, error) {
	if !utils.IsValidUUID(perusahaanID) {
		return nil, errors.New("format perusahaan_id tidak valid")
	}
	return s.repo.GetByPerusahaan(perusahaanID)
}

func (s *JawabanDeteksiService) GetByPertanyaan(pertanyaanID int) ([]dto.JawabanDeteksiResponse, error) {
	if pertanyaanID <= 0 {
		return nil, errors.New("pertanyaan_deteksi_id tidak valid")
	}
	return s.repo.GetByPertanyaan(pertanyaanID)
}

func (s *JawabanDeteksiService) Update(id int, req dto.UpdateJawabanDeteksiRequest) (*dto.JawabanDeteksiResponse, error) {
	if id <= 0 {
		return nil, errors.New("format ID tidak valid")
	}

	existing, err := s.repo.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		if err == sql.ErrNoRows {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}

	if err := s.validateUpdate(&req, existing.Evidence); err != nil {
		return nil, err
	}

	if err := s.repo.Update(id, req); err != nil {
		rollbar.Error(err)
		return nil, err
	}

	if err := s.repo.RecalculateDeteksi(existing.PerusahaanID); err != nil {
		rollbar.Error(err)
	}

	updated, err := s.repo.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return updated, nil
}

func (s *JawabanDeteksiService) Delete(id int) error {
	if id <= 0 {
		return errors.New("format ID tidak valid")
	}

	existing, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("data tidak ditemukan")
		}
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	if err := s.repo.RecalculateDeteksi(existing.PerusahaanID); err != nil {
		rollbar.Error(err)
	}

	return nil
}
