package services

import (
	"context"
	"database/sql"
	"errors"
	"ikas/internal/dto"
	"ikas/internal/dto/dto_event"
	"ikas/internal/repository"
	"ikas/internal/utils"
	"time"

	"github.com/rollbar/rollbar-go"
)

type JawabanProteksiProducerInterface interface {
	PublishJawabanProteksiCreated(ctx context.Context, event interface{}) error
	PublishJawabanProteksiUpdated(ctx context.Context, event interface{}) error
	PublishJawabanProteksiDeleted(ctx context.Context, event interface{}) error
	PublishIkasAuditLog(ctx context.Context, event interface{}) error
}

type JawabanProteksiService struct {
	repo     repository.JawabanProteksiRepositoryInterface
	ikasRepo repository.IkasRepositoryInterface
	producer JawabanProteksiProducerInterface
}

func NewJawabanProteksiService(
	repo repository.JawabanProteksiRepositoryInterface,
	ikasRepo repository.IkasRepositoryInterface,
	producer JawabanProteksiProducerInterface,
) *JawabanProteksiService {
	return &JawabanProteksiService{
		repo:     repo,
		ikasRepo: ikasRepo,
		producer: producer,
	}
}

var validValidasiProteksi = map[string]bool{"yes": true, "no": true}

func (s *JawabanProteksiService) validateCreate(req *dto.CreateJawabanProteksiRequest, userRole string) error {
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

	if req.JawabanProteksi == nil {
		return errors.New("jawaban_proteksi tidak boleh kosong")
	}
	if *req.JawabanProteksi < 0 || *req.JawabanProteksi > 5 {
		return errors.New("jawaban_proteksi harus bernilai antara 0 sampai 5")
	}

	// Restricted fields for non-admins
	if userRole != "admin" {
		if req.Validasi != nil || (req.Keterangan != nil && utils.NormalizeInput(*req.Keterangan) != "") {
			return errors.New("hanya admin yang dapat mengisi field validasi dan keterangan")
		}
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

func (s *JawabanProteksiService) validateUpdate(req *dto.UpdateJawabanProteksiRequest, existingEvidence *string, userRole string) error {
	if req.JawabanProteksi != nil && (*req.JawabanProteksi < 0 || *req.JawabanProteksi > 5) {
		return errors.New("jawaban_proteksi harus bernilai antara 0 sampai 5, atau null for N/A")
	}

	// Restricted fields for non-admins
	if userRole != "admin" {
		if req.Validasi != nil || (req.Keterangan != nil && utils.NormalizeInput(*req.Keterangan) != "") {
			return errors.New("hanya admin yang dapat mengubah field validasi dan keterangan")
		}
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

func (s *JawabanProteksiService) Create(req dto.CreateJawabanProteksiRequest, userRole string) (string, error) {
	if err := s.validateCreate(&req, userRole); err != nil {
		return "", err
	}

	pertanyaanExists, err := s.repo.CheckPertanyaanExists(req.PertanyaanProteksiID)
	if err != nil {
		rollbar.Error(err)
		return "", err
	}
	if !pertanyaanExists {
		return "", errors.New("pertanyaan_proteksi_id tidak ditemukan")
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
	isDuplicate, err := s.repo.CheckDuplicate(req.PerusahaanID, req.PertanyaanProteksiID, 0)
	if err != nil {
		rollbar.Error(err)
		return "", err
	}
	if isDuplicate {
		return "", errors.New("pertanyaan ini sudah pernah diisi oleh perusahaan Anda")
	}

	// Publish to RabbitMQ for Pola 2
	if err := s.producer.PublishJawabanProteksiCreated(context.Background(), req); err != nil {
		rollbar.Error(err)
		return "", err
	}

	return "Berhasil menyimpan data", nil
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

func (s *JawabanProteksiService) Update(id int, req dto.UpdateJawabanProteksiRequest, userID string, userRole string) error {
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

	if err := s.validateUpdate(&req, existing.Evidence, userRole); err != nil {
		return err
	}

	// Publish Update Event (Pola 2)
	event := dto_event.JawabanProteksiUpdatedEvent{
		ID:        id,
		Request:   req,
		UpdatedAt: time.Now(),
	}

	// Change detection for audit log
	changes := make(map[string]interface{})
	if req.JawabanProteksi != nil && (existing.JawabanProteksi == nil || *req.JawabanProteksi != *existing.JawabanProteksi) {
		oldVal := interface{}(nil)
		if existing.JawabanProteksi != nil {
			oldVal = *existing.JawabanProteksi
		}
		changes["jawaban_proteksi"] = map[string]interface{}{"old": oldVal, "new": *req.JawabanProteksi}
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
				Action:    "UPDATE_PROTEKSI",
				Changes:   changes,
				Timestamp: time.Now(),
			}
			_ = s.producer.PublishIkasAuditLog(context.Background(), auditEvent)
		}
	}

	if err := s.producer.PublishJawabanProteksiUpdated(context.Background(), event); err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}

func (s *JawabanProteksiService) Delete(id int, userID string) error {
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
	event := dto_event.JawabanProteksiDeletedEvent{
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
				Action:    "DELETE_PROTEKSI",
				Changes:   map[string]interface{}{"pertanyaan_id": existing.PertanyaanProteksi.ID, "status": "deleted"},
				Timestamp: time.Now(),
			}
			_ = s.producer.PublishIkasAuditLog(context.Background(), auditEvent)
		}
	}

	if err := s.producer.PublishJawabanProteksiDeleted(context.Background(), event); err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}
