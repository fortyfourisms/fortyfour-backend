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

type JawabanDeteksiService struct {
	repo     repository.JawabanDeteksiRepositoryInterface
	ikasRepo repository.IkasRepositoryInterface
	producer *rabbitmq.Producer
}

func NewJawabanDeteksiService(
	repo repository.JawabanDeteksiRepositoryInterface,
	ikasRepo repository.IkasRepositoryInterface,
	producer *rabbitmq.Producer,
) *JawabanDeteksiService {
	return &JawabanDeteksiService{
		repo:     repo,
		ikasRepo: ikasRepo,
		producer: producer,
	}
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

	if req.JawabanDeteksi == nil {
		return errors.New("jawaban_deteksi tidak boleh kosong")
	}
	if *req.JawabanDeteksi < 0 || *req.JawabanDeteksi > 5 {
		return errors.New("jawaban_deteksi harus bernilai antara 0 sampai 5")
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

func (s *JawabanDeteksiService) Create(req dto.CreateJawabanDeteksiRequest) (string, error) {
	if err := s.validateCreate(&req); err != nil {
		return "", err
	}

	pertanyaanExists, err := s.repo.CheckPertanyaanExists(req.PertanyaanDeteksiID)
	if err != nil {
		rollbar.Error(err)
		return "", err
	}
	if !pertanyaanExists {
		return "", errors.New("pertanyaan_deteksi_id tidak ditemukan")
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
	isDuplicate, err := s.repo.CheckDuplicate(req.PerusahaanID, req.PertanyaanDeteksiID, 0)
	if err != nil {
		rollbar.Error(err)
		return "", err
	}
	if isDuplicate {
		return "", errors.New("pertanyaan ini sudah pernah diisi oleh perusahaan Anda")
	}

	// Publish to RabbitMQ for Pola 2
	if err := s.producer.PublishJawabanDeteksiCreated(context.Background(), req); err != nil {
		rollbar.Error(err)
		return "", err
	}

	return "Berhasil menyimpan data", nil
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

func (s *JawabanDeteksiService) Update(id int, req dto.UpdateJawabanDeteksiRequest, userID string) error {
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
	event := dto_event.JawabanDeteksiUpdatedEvent{
		ID:        id,
		Request:   req,
		UpdatedAt: time.Now(),
	}

	// Change detection for audit log
	changes := make(map[string]interface{})
	if req.JawabanDeteksi != nil && (existing.JawabanDeteksi == nil || *req.JawabanDeteksi != *existing.JawabanDeteksi) {
		oldVal := interface{}(nil)
		if existing.JawabanDeteksi != nil {
			oldVal = *existing.JawabanDeteksi
		}
		changes["jawaban_deteksi"] = map[string]interface{}{"old": oldVal, "new": *req.JawabanDeteksi}
	}
	if req.Evidence != nil && (existing.Evidence == nil || *req.Evidence != *existing.Evidence) {
		oldVal := interface{}(nil)
		if existing.Evidence != nil {
			oldVal = *existing.Evidence
		}
		changes["evidence"] = map[string]interface{}{"old": oldVal, "new": *req.Evidence}
	}
	if req.Validasi != nil && (existing.Validasi == nil || *req.Validasi != *existing.Validasi) {
		oldVal := interface{}(nil)
		if existing.Validasi != nil {
			oldVal = *existing.Validasi
		}
		changes["validasi"] = map[string]interface{}{"old": oldVal, "new": *req.Validasi}
	}
	if req.Keterangan != nil && (existing.Keterangan == nil || *req.Keterangan != *existing.Keterangan) {
		oldVal := interface{}(nil)
		if existing.Keterangan != nil {
			oldVal = *existing.Keterangan
		}
		changes["keterangan"] = map[string]interface{}{"old": oldVal, "new": *req.Keterangan}
	}

	if s.producer != nil && len(changes) > 0 {
		ikasID, err := s.ikasRepo.GetIDByPerusahaanID(existing.PerusahaanID)
		if err == nil {
			auditEvent := dto_event.IkasAuditLogEvent{
				IkasID:    ikasID,
				UserID:    userID,
				Action:    "UPDATE_DETEKSI",
				Changes:   changes,
				Timestamp: time.Now(),
			}
			_ = s.producer.PublishIkasAuditLog(context.Background(), auditEvent)
		}
	}

	if err := s.producer.PublishJawabanDeteksiUpdated(context.Background(), event); err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}

func (s *JawabanDeteksiService) Delete(id int, userID string) error {
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
	event := dto_event.JawabanDeteksiDeletedEvent{
		ID:           id,
		PerusahaanID: existing.PerusahaanID,
		DeletedAt:    time.Now(),
	}

	if s.producer != nil {
		ikasID, err := s.ikasRepo.GetIDByPerusahaanID(existing.PerusahaanID)
		if err == nil {
			auditEvent := dto_event.IkasAuditLogEvent{
				IkasID:    ikasID,
				UserID:    userID,
				Action:    "DELETE_DETEKSI",
				Changes:   map[string]interface{}{"pertanyaan_id": existing.PertanyaanDeteksi.ID, "status": "deleted"},
				Timestamp: time.Now(),
			}
			_ = s.producer.PublishIkasAuditLog(context.Background(), auditEvent)
		}
	}

	if err := s.producer.PublishJawabanDeteksiDeleted(context.Background(), event); err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}
