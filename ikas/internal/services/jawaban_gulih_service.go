package services

import (
	"context"
	"database/sql"
	"errors"
	"ikas/internal/dto"
	"ikas/internal/dto/dto_event"
	"ikas/internal/rabbitmq"
	"ikas/internal/repository"
	"ikas/internal/utils"
	"time"

	"github.com/rollbar/rollbar-go"
)

type JawabanGulihService struct {
	repo     repository.JawabanGulihRepositoryInterface
	producer *rabbitmq.Producer
}

func NewJawabanGulihService(repo repository.JawabanGulihRepositoryInterface, producer *rabbitmq.Producer) *JawabanGulihService {
	return &JawabanGulihService{
		repo:     repo,
		producer: producer,
	}
}

var validValidasiGulih = map[string]bool{"yes": true, "no": true}

func (s *JawabanGulihService) validateCreate(req *dto.CreateJawabanGulihRequest) error {
	if req.PertanyaanGulihID <= 0 {
		return errors.New("pertanyaan_gulih_id tidak valid")
	}

	req.PerusahaanID = utils.NormalizeInput(req.PerusahaanID)
	if req.PerusahaanID == "" {
		return errors.New("perusahaan_id tidak boleh kosong")
	}
	if !utils.IsValidUUID(req.PerusahaanID) {
		return errors.New("format perusahaan_id tidak valid")
	}

	if req.JawabanGulih == nil {
		return errors.New("jawaban_gulih tidak boleh kosong")
	}
	if *req.JawabanGulih < 0 || *req.JawabanGulih > 5 {
		return errors.New("jawaban_gulih harus bernilai antara 0 sampai 5")
	}

	if req.Validasi != nil {
		if req.Evidence == nil || utils.NormalizeInput(*req.Evidence) == "" {
			return errors.New("validasi hanya boleh diisi jika evidence ada")
		}
		if !validValidasiGulih[*req.Validasi] {
			return errors.New("validasi hanya boleh berisi 'yes' atau 'no'")
		}
	}

	return nil
}

func (s *JawabanGulihService) validateUpdate(req *dto.UpdateJawabanGulihRequest, existingEvidence *string) error {
	if req.JawabanGulih != nil && (*req.JawabanGulih < 0 || *req.JawabanGulih > 5) {
		return errors.New("jawaban_gulih harus bernilai antara 0 sampai 5, atau null untuk N/A")
	}

	if req.Validasi != nil {
		if !validValidasiGulih[*req.Validasi] {
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

func (s *JawabanGulihService) Create(req dto.CreateJawabanGulihRequest) (string, error) {
	if err := s.validateCreate(&req); err != nil {
		return "", err
	}

	pertanyaanExists, err := s.repo.CheckPertanyaanExists(req.PertanyaanGulihID)
	if err != nil {
		rollbar.Error(err)
		return "", err
	}
	if !pertanyaanExists {
		return "", errors.New("pertanyaan_gulih_id tidak ditemukan")
	}

	perusahaanExists, err := s.repo.CheckPerusahaanExists(req.PerusahaanID)
	if err != nil {
		rollbar.Error(err)
		return "", err
	}
	if !perusahaanExists {
		return "", errors.New("perusahaan_id tidak ditemukan")
	}

	// Synchronous Duplicate Check (Pola 2 Refinement)
	isDuplicate, err := s.repo.CheckDuplicate(req.PerusahaanID, req.PertanyaanGulihID, 0)
	if err != nil {
		rollbar.Error(err)
		return "", err
	}
	if isDuplicate {
		return "", errors.New("pertanyaan ini sudah pernah diisi oleh perusahaan Anda")
	}

	// Publish to RabbitMQ for Pola 2
	if err := s.producer.PublishJawabanGulihCreated(context.Background(), req); err != nil {
		rollbar.Error(err)
		return "", err
	}

	return "Berhasil menyimpan data", nil
}

func (s *JawabanGulihService) GetAll() ([]dto.JawabanGulihResponse, error) {
	return s.repo.GetAll()
}

func (s *JawabanGulihService) GetByID(id int) (*dto.JawabanGulihResponse, error) {
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

func (s *JawabanGulihService) GetByPerusahaan(perusahaanID string) ([]dto.JawabanGulihResponse, error) {
	if !utils.IsValidUUID(perusahaanID) {
		return nil, errors.New("format perusahaan_id tidak valid")
	}
	return s.repo.GetByPerusahaan(perusahaanID)
}

func (s *JawabanGulihService) GetByPertanyaan(pertanyaanID int) ([]dto.JawabanGulihResponse, error) {
	if pertanyaanID <= 0 {
		return nil, errors.New("pertanyaan_gulih_id tidak valid")
	}
	return s.repo.GetByPertanyaan(pertanyaanID)
}

func (s *JawabanGulihService) Update(id int, req dto.UpdateJawabanGulihRequest) error {
	if id <= 0 {
		return errors.New("format ID tidak valid")
	}

	// Existence Check
	existing, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("data tidak ditemukan")
		}
		return err
	}

	if err := s.validateUpdate(&req, existing.Evidence); err != nil {
		return err
	}

	// Publish Update Event (Pola 2)
	event := dto_event.JawabanGulihUpdatedEvent{
		ID:        id,
		Request:   req,
		UpdatedAt: time.Now(),
	}

	if err := s.producer.PublishJawabanGulihUpdated(context.Background(), event); err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}

func (s *JawabanGulihService) Delete(id int) error {
	if id <= 0 {
		return errors.New("format ID tidak valid")
	}

	// Existence Check
	existing, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("data tidak ditemukan")
		}
		return err
	}

	// Publish Delete Event (Pola 2)
	event := dto_event.JawabanGulihDeletedEvent{
		ID:           id,
		PerusahaanID: existing.PerusahaanID,
		DeletedAt:    time.Now(),
	}

	if err := s.producer.PublishJawabanGulihDeleted(context.Background(), event); err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}
