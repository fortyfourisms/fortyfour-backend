package services

import (
	"database/sql"
	"errors"
	"ikas/internal/dto"
	"ikas/internal/repository"
	"ikas/internal/utils"

	"github.com/rollbar/rollbar-go"
)

type JawabanProteksiService struct {
	repo repository.JawabanProteksiRepositoryInterface
}

func NewJawabanProteksiService(repo repository.JawabanProteksiRepositoryInterface) *JawabanProteksiService {
	return &JawabanProteksiService{repo: repo}
}

var validValidasiProteksi = map[string]bool{"yes": true, "no": true}

func (s *JawabanProteksiService) validateCreate(req *dto.CreateJawabanProteksiRequest) error {
	if req.PertanyaanProteksiID <= 0 {
		return errors.New("pertanyaan_proteksi_id tidak valid")
	}

	req.PerusahaanID = utils.NormalizeInput(req.PerusahaanID)
	if req.PerusahaanID == "" {
		return errors.New("perusahaan_id tidak boleh kosong")
	}
	if !utils.IsValidUUID(req.PerusahaanID) {
		return errors.New("format perusahaan_id tidak valid")
	}

	if req.JawabanProteksi != nil && (*req.JawabanProteksi < 0 || *req.JawabanProteksi > 5) {
		return errors.New("jawaban_proteksi harus bernilai antara 0 sampai 5, atau null untuk N/A")
	}

	if req.Validasi != nil {
		if req.Evidence == nil || utils.NormalizeInput(*req.Evidence) == "" {
			return errors.New("validasi hanya boleh diisi jika evidence ada")
		}
		if !validValidasiProteksi[*req.Validasi] {
			return errors.New("validasi hanya boleh berisi 'yes' atau 'no'")
		}
	}

	return nil
}

func (s *JawabanProteksiService) validateUpdate(req *dto.UpdateJawabanProteksiRequest, existingEvidence *string) error {
	if req.JawabanProteksi != nil && (*req.JawabanProteksi < 0 || *req.JawabanProteksi > 5) {
		return errors.New("jawaban_proteksi harus bernilai antara 0 sampai 5, atau null untuk N/A")
	}

	if req.Validasi != nil {
		if !validValidasiProteksi[*req.Validasi] {
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

func (s *JawabanProteksiService) Create(req dto.CreateJawabanProteksiRequest) (*dto.JawabanProteksiResponse, error) {
	if err := s.validateCreate(&req); err != nil {
		return nil, err
	}

	pertanyaanExists, err := s.repo.CheckPertanyaanExists(req.PertanyaanProteksiID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	if !pertanyaanExists {
		return nil, errors.New("pertanyaan_proteksi_id tidak ditemukan")
	}

	perusahaanExists, err := s.repo.CheckPerusahaanExists(req.PerusahaanID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	if !perusahaanExists {
		return nil, errors.New("perusahaan_id tidak ditemukan")
	}

	isDuplicate, err := s.repo.CheckDuplicate(req.PerusahaanID, req.PertanyaanProteksiID, 0)
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

	if err := s.repo.RecalculateProteksi(req.PerusahaanID); err != nil {
		rollbar.Error(err)
	}

	resp, err := s.repo.GetByID(int(lastID))
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return resp, nil
}

func (s *JawabanProteksiService) GetAll() ([]dto.JawabanProteksiResponse, error) {
	return s.repo.GetAll()
}

func (s *JawabanProteksiService) GetByID(id int) (*dto.JawabanProteksiResponse, error) {
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

func (s *JawabanProteksiService) GetByPerusahaan(perusahaanID string) ([]dto.JawabanProteksiResponse, error) {
	if !utils.IsValidUUID(perusahaanID) {
		return nil, errors.New("format perusahaan_id tidak valid")
	}
	return s.repo.GetByPerusahaan(perusahaanID)
}

func (s *JawabanProteksiService) GetByPertanyaan(pertanyaanID int) ([]dto.JawabanProteksiResponse, error) {
	if pertanyaanID <= 0 {
		return nil, errors.New("pertanyaan_proteksi_id tidak valid")
	}
	return s.repo.GetByPertanyaan(pertanyaanID)
}

func (s *JawabanProteksiService) Update(id int, req dto.UpdateJawabanProteksiRequest) (*dto.JawabanProteksiResponse, error) {
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

	if err := s.repo.RecalculateProteksi(existing.PerusahaanID); err != nil {
		rollbar.Error(err)
	}

	updated, err := s.repo.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return updated, nil
}

func (s *JawabanProteksiService) Delete(id int) error {
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

	if err := s.repo.RecalculateProteksi(existing.PerusahaanID); err != nil {
		rollbar.Error(err)
	}

	return nil
}
