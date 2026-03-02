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

type JawabanIdentifikasiService struct {
	repo repository.JawabanIdentifikasiRepositoryInterface
}

func NewJawabanIdentifikasiService(repo repository.JawabanIdentifikasiRepositoryInterface) *JawabanIdentifikasiService {
	return &JawabanIdentifikasiService{repo: repo}
}

var validValidasi = map[string]bool{"yes": true, "no": true}

func (s *JawabanIdentifikasiService) validateCreate(req *dto.CreateJawabanIdentifikasiRequest) error {
	req.PertanyaanIdentifikasiID = utils.NormalizeInput(req.PertanyaanIdentifikasiID)
	if req.PertanyaanIdentifikasiID == "" {
		return errors.New("pertanyaan_identifikasi_id tidak boleh kosong")
	}
	if !utils.IsValidUUID(req.PertanyaanIdentifikasiID) {
		return errors.New("format pertanyaan_identifikasi_id tidak valid")
	}

	req.PerusahaanID = utils.NormalizeInput(req.PerusahaanID)
	if req.PerusahaanID == "" {
		return errors.New("perusahaan_id tidak boleh kosong")
	}
	if !utils.IsValidUUID(req.PerusahaanID) {
		return errors.New("format perusahaan_id tidak valid")
	}

	// null = N/A (diperbolehkan), tapi jika diisi harus 0.00 - 5.00
	if req.JawabanIdentifikasi != nil && (*req.JawabanIdentifikasi < 0 || *req.JawabanIdentifikasi > 5) {
		return errors.New("jawaban_identifikasi harus bernilai antara 0 sampai 5, atau null untuk N/A")
	}

	if req.Validasi != nil {
		if req.Evidence == nil || utils.NormalizeInput(*req.Evidence) == "" {
			return errors.New("validasi hanya boleh diisi jika evidence ada")
		}
		if !validValidasi[*req.Validasi] {
			return errors.New("validasi hanya boleh berisi 'yes' atau 'no'")
		}
	}

	return nil
}

func (s *JawabanIdentifikasiService) validateUpdate(req *dto.UpdateJawabanIdentifikasiRequest, existingEvidence *string) error {
	// null = N/A (diperbolehkan), tapi jika diisi harus 0.00 - 5.00
	if req.JawabanIdentifikasi != nil && (*req.JawabanIdentifikasi < 0 || *req.JawabanIdentifikasi > 5) {
		return errors.New("jawaban_identifikasi harus bernilai antara 0 sampai 5, atau null untuk N/A")
	}

	if req.Validasi != nil {
		if !validValidasi[*req.Validasi] {
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

func (s *JawabanIdentifikasiService) Create(req dto.CreateJawabanIdentifikasiRequest) (*dto.JawabanIdentifikasiResponse, error) {
	if err := s.validateCreate(&req); err != nil {
		return nil, err
	}

	pertanyaanExists, err := s.repo.CheckPertanyaanExists(req.PertanyaanIdentifikasiID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	if !pertanyaanExists {
		return nil, errors.New("pertanyaan_identifikasi_id tidak ditemukan")
	}

	perusahaanExists, err := s.repo.CheckPerusahaanExists(req.PerusahaanID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	if !perusahaanExists {
		return nil, errors.New("perusahaan_id tidak ditemukan")
	}

	isDuplicate, err := s.repo.CheckDuplicate(req.PerusahaanID, req.PertanyaanIdentifikasiID, "")
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	if isDuplicate {
		return nil, errors.New("jawaban untuk pertanyaan ini sudah ada untuk perusahaan tersebut")
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

func (s *JawabanIdentifikasiService) GetAll() ([]dto.JawabanIdentifikasiResponse, error) {
	return s.repo.GetAll()
}

func (s *JawabanIdentifikasiService) GetByID(id string) (*dto.JawabanIdentifikasiResponse, error) {
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

func (s *JawabanIdentifikasiService) GetByPerusahaan(perusahaanID string) ([]dto.JawabanIdentifikasiResponse, error) {
	if !utils.IsValidUUID(perusahaanID) {
		return nil, errors.New("format perusahaan_id tidak valid")
	}
	return s.repo.GetByPerusahaan(perusahaanID)
}

func (s *JawabanIdentifikasiService) GetByPertanyaan(pertanyaanID string) ([]dto.JawabanIdentifikasiResponse, error) {
	if !utils.IsValidUUID(pertanyaanID) {
		return nil, errors.New("format pertanyaan_identifikasi_id tidak valid")
	}
	return s.repo.GetByPertanyaan(pertanyaanID)
}

func (s *JawabanIdentifikasiService) Update(id string, req dto.UpdateJawabanIdentifikasiRequest) (*dto.JawabanIdentifikasiResponse, error) {
	if !utils.IsValidUUID(id) {
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

	updated, err := s.repo.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return updated, nil
}

func (s *JawabanIdentifikasiService) Delete(id string) error {
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
