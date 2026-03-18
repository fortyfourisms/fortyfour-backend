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

type JawabanIdentifikasiService struct {
	repo     repository.JawabanIdentifikasiRepositoryInterface
	ikasRepo repository.IkasRepositoryInterface
	producer *rabbitmq.Producer
}

func NewJawabanIdentifikasiService(
	repo repository.JawabanIdentifikasiRepositoryInterface,
	ikasRepo repository.IkasRepositoryInterface,
	producer *rabbitmq.Producer,
) *JawabanIdentifikasiService {
	return &JawabanIdentifikasiService{
		repo:     repo,
		ikasRepo: ikasRepo,
		producer: producer,
	}
}

var validValidasi = map[string]bool{"yes": true, "no": true}

func (s *JawabanIdentifikasiService) validateCreate(req *dto.CreateJawabanIdentifikasiRequest) error {
	if req.PertanyaanIdentifikasiID <= 0 {
		return errors.New("pertanyaan_identifikasi_id tidak valid")
	}

	req.PerusahaanID = utils.NormalizeInput(req.PerusahaanID)
	if req.PerusahaanID == "" {
		return errors.New("perusahaan_id tidak boleh kosong")
	}
	if !utils.IsValidUUID(req.PerusahaanID) {
		return errors.New("format perusahaan_id tidak valid")
	}

	// Jawaban must be provided and between 0.00 - 5.00
	if req.JawabanIdentifikasi == nil {
		return errors.New("jawaban_identifikasi tidak boleh kosong")
	}
	if *req.JawabanIdentifikasi < 0 || *req.JawabanIdentifikasi > 5 {
		return errors.New("jawaban_identifikasi harus bernilai antara 0 sampai 5")
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

func (s *JawabanIdentifikasiService) Create(req dto.CreateJawabanIdentifikasiRequest) (string, error) {
	if err := s.validateCreate(&req); err != nil {
		return "", err
	}

	pertanyaanExists, err := s.repo.CheckPertanyaanExists(req.PertanyaanIdentifikasiID)
	if err != nil {
		rollbar.Error(err)
		return "", err
	}
	if !pertanyaanExists {
		return "", errors.New("pertanyaan_identifikasi_id tidak ditemukan")
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
	// Check if already exists in the MAIN table
	isDuplicate, err := s.repo.CheckDuplicate(req.PerusahaanID, req.PertanyaanIdentifikasiID, 0)
	if err != nil {
		rollbar.Error(err)
		return "", err
	}
	if isDuplicate {
		return "", errors.New("pertanyaan ini sudah pernah diisi oleh perusahaan Anda")
	}

	// Publish to RabbitMQ for Pola 2
	if err := s.producer.PublishJawabanIdentifikasiCreated(context.Background(), req); err != nil {
		rollbar.Error(err)
		return "", err
	}

	return "Berhasil menyimpan data", nil
}

func (s *JawabanIdentifikasiService) GetAll() ([]dto.JawabanIdentifikasiResponse, error) {
	return s.repo.GetAll()
}

func (s *JawabanIdentifikasiService) GetByID(id int) (*dto.JawabanIdentifikasiResponse, error) {
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

func (s *JawabanIdentifikasiService) GetByPerusahaan(perusahaanID string) ([]dto.JawabanIdentifikasiResponse, error) {
	if !utils.IsValidUUID(perusahaanID) {
		return nil, errors.New("format perusahaan_id tidak valid")
	}
	return s.repo.GetByPerusahaan(perusahaanID)
}

func (s *JawabanIdentifikasiService) GetByPertanyaan(pertanyaanID int) ([]dto.JawabanIdentifikasiResponse, error) {
	if pertanyaanID <= 0 {
		return nil, errors.New("pertanyaan_identifikasi_id tidak valid")
	}
	return s.repo.GetByPertanyaan(pertanyaanID)
}

func (s *JawabanIdentifikasiService) Update(id int, req dto.UpdateJawabanIdentifikasiRequest, userID string) error {
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
	event := dto_event.JawabanIdentifikasiUpdatedEvent{
		ID:        id,
		Request:   req,
		UpdatedAt: time.Now(),
	}

	// Change detection for audit log
	changes := make(map[string]interface{})
	if req.JawabanIdentifikasi != nil && (existing.JawabanIdentifikasi == nil || *req.JawabanIdentifikasi != *existing.JawabanIdentifikasi) {
		oldVal := interface{}(nil)
		if existing.JawabanIdentifikasi != nil {
			oldVal = *existing.JawabanIdentifikasi
		}
		changes["jawaban_identifikasi"] = map[string]interface{}{"old": oldVal, "new": *req.JawabanIdentifikasi}
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
				Action:    "UPDATE_IDENTIFIKASI",
				Changes:   changes,
				Timestamp: time.Now(),
			}
			_ = s.producer.PublishIkasAuditLog(context.Background(), auditEvent)
		}
	}

	if err := s.producer.PublishJawabanIdentifikasiUpdated(context.Background(), event); err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}

func (s *JawabanIdentifikasiService) Delete(id int, userID string) error {
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
	event := dto_event.JawabanIdentifikasiDeletedEvent{
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
				Action:    "DELETE_IDENTIFIKASI",
				Changes:   map[string]interface{}{"pertanyaan_id": existing.PertanyaanIdentifikasi.ID, "status": "deleted"},
				Timestamp: time.Now(),
			}
			_ = s.producer.PublishIkasAuditLog(context.Background(), auditEvent)
		}
	}

	if err := s.producer.PublishJawabanIdentifikasiDeleted(context.Background(), event); err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}
