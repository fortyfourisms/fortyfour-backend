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

type JawabanDeteksiProducerInterface interface {
	PublishJawabanDeteksiCreated(ctx context.Context, event interface{}) error
	PublishJawabanDeteksiUpdated(ctx context.Context, event interface{}) error
	PublishJawabanDeteksiDeleted(ctx context.Context, event interface{}) error
	PublishIkasAuditLog(ctx context.Context, event interface{}) error
}

type JawabanDeteksiService struct {
	repo     repository.JawabanDeteksiRepositoryInterface
	ikasRepo repository.IkasRepositoryInterface
	producer JawabanDeteksiProducerInterface
}

func NewJawabanDeteksiService(
	repo repository.JawabanDeteksiRepositoryInterface,
	ikasRepo repository.IkasRepositoryInterface,
	producer JawabanDeteksiProducerInterface,
) *JawabanDeteksiService {
	return &JawabanDeteksiService{
		repo:     repo,
		ikasRepo: ikasRepo,
		producer: producer,
	}
}

var validValidasiDeteksi = map[string]bool{"yes": true, "no": true}

func (s *JawabanDeteksiService) validateCreate(req *dto.CreateJawabanDeteksiRequest, userRole string) error {
	if req.PertanyaanDeteksiID <= 0 {
		return errors.New("pertanyaan_deteksi_id tidak valid")
	}

	req.IkasID = utils.NormalizeInput(req.IkasID)
	if req.IkasID == "" {
		return errors.New("ikas_id tidak boleh kosong")
	}
	if !utils.IsValidUUID(req.IkasID) {
		return errors.New("format ikas_id tidak valid")
	}

	if req.JawabanDeteksi == nil {
		return errors.New("jawaban_deteksi tidak boleh kosong")
	}
	if *req.JawabanDeteksi < 0 || *req.JawabanDeteksi > 5 {
		return errors.New("jawaban_deteksi harus bernilai antara 0 sampai 5")
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
		if !validValidasiDeteksi[*req.Validasi] {
			return errors.New("validasi hanya boleh berisi 'yes' atau 'no'")
		}
	}

	return nil
}

func (s *JawabanDeteksiService) validateUpdate(req *dto.UpdateJawabanDeteksiRequest, existingEvidence *string, userRole string) error {
	if req.JawabanDeteksi != nil && (*req.JawabanDeteksi < 0 || *req.JawabanDeteksi > 5) {
		return errors.New("jawaban_deteksi harus bernilai antara 0 sampai 5, atau null for N/A")
	}

	// Restricted fields for non-admins
	if userRole != "admin" {
		if req.Validasi != nil || (req.Keterangan != nil && utils.NormalizeInput(*req.Keterangan) != "") {
			return errors.New("hanya admin yang dapat mengubah field validasi dan keterangan")
		}
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

func (s *JawabanDeteksiService) Create(req dto.CreateJawabanDeteksiRequest, userRole string, userPerusahaanID string) (string, error) {
	if err := s.validateCreate(&req, userRole); err != nil {
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

	ikasExists, err := s.repo.CheckIkasExists(req.IkasID)
	if err != nil {
		rollbar.Error(err)
		return "", err
	}
	if !ikasExists {
		return "", errors.New("ikas_id tidak ditemukan")
	}

	// VALIDASI KEPEMILIKAN
	if userRole != "admin" {
		owned, err := s.ikasRepo.CheckOwnership(req.IkasID, userPerusahaanID)
		if err != nil {
			rollbar.Error(err)
			return "", err
		}
		if !owned {
			return "", errors.New("anda tidak memiliki akses ke data asesmen ini")
		}
	}

	// CHECK LOCK
	locked, err := s.ikasRepo.IsLocked(req.IkasID)
	if err != nil {
		rollbar.Error(err)
		return "", err
	}
	if locked {
		return "", errors.New("data asesmen ini sudah divalidasi dan tidak dapat diubah")
	}

	// Synchronous Duplicate Check (Pola 2 Refinement)
	isDuplicate, err := s.repo.CheckDuplicate(req.IkasID, req.PertanyaanDeteksiID, 0)
	if err != nil {
		rollbar.Error(err)
		return "", err
	}
	if isDuplicate {
		return "", errors.New("pertanyaan ini sudah pernah diisi untuk asesmen ini")
	}

	// Publish to RabbitMQ for Pola 2
	if err := s.producer.PublishJawabanDeteksiCreated(context.Background(), req); err != nil {
		rollbar.Error(err)
		return "", err
	}

	return "Berhasil menyimpan data", nil
}

func (s *JawabanDeteksiService) GetAll(userRole string) ([]dto.JawabanDeteksiResponse, error) {
	if userRole != "admin" {
		return nil, errors.New("anda tidak memiliki akses untuk melihat semua data")
	}
	return s.repo.GetAll()
}

func (s *JawabanDeteksiService) GetByID(id int, userRole string, userPerusahaanID string) (*dto.JawabanDeteksiResponse, error) {
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

	// Fetch ikas to check ownership
	ikasData, err := s.ikasRepo.GetByID(data.IkasID)
	if err != nil {
		return nil, errors.New("gagal memverifikasi kepemilikan asesmen")
	}

	if userRole != "admin" && ikasData.Perusahaan.ID != userPerusahaanID {
		return nil, errors.New("anda tidak memiliki akses ke data ini")
	}

	return data, nil
}

func (s *JawabanDeteksiService) GetByIkasID(ikasID string, userRole string, userPerusahaanID string) ([]dto.JawabanDeteksiResponse, error) {
	if !utils.IsValidUUID(ikasID) {
		return nil, errors.New("format ikas_id tidak valid")
	}

	if userRole != "admin" {
		owned, err := s.ikasRepo.CheckOwnership(ikasID, userPerusahaanID)
		if err != nil {
			return nil, err
		}
		if !owned {
			return nil, errors.New("anda tidak memiliki akses ke data asesmen ini")
		}
	}

	return s.repo.GetByIkasID(ikasID)
}

func (s *JawabanDeteksiService) GetByPerusahaanID(perusahaanID string, userRole string, userPerusahaanID string) ([]dto.JawabanDeteksiResponse, error) {
	if userRole != "admin" {
		if perusahaanID != userPerusahaanID {
			return nil, errors.New("anda tidak memiliki akses ke data perusahaan ini")
		}
	}
	return s.repo.GetByPerusahaanID(perusahaanID)
}

func (s *JawabanDeteksiService) GetByPertanyaan(pertanyaanID int) ([]dto.JawabanDeteksiResponse, error) {
	if pertanyaanID <= 0 {
		return nil, errors.New("pertanyaan_deteksi_id tidak valid")
	}
	return s.repo.GetByPertanyaan(pertanyaanID)
}

func (s *JawabanDeteksiService) Update(id int, req dto.UpdateJawabanDeteksiRequest, userID string, userRole string, userPerusahaanID string) error {
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

	// Fetch ikas to check ownership
	ikasData, err := s.ikasRepo.GetByID(existing.IkasID)
	if err != nil {
		return errors.New("gagal memverifikasi kepemilikan asesmen")
	}

	if userRole != "admin" && ikasData.Perusahaan.ID != userPerusahaanID {
		return errors.New("anda tidak memiliki akses untuk mengubah data ini")
	}

	if ikasData.IsValidated {
		return errors.New("data asesmen ini sudah divalidasi dan tidak dapat diubah")
	}

	if err := s.validateUpdate(&req, existing.Evidence, userRole); err != nil {
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
		auditEvent := dto_event.IkasAuditLogEvent{
			IkasID:    existing.IkasID,
			UserID:    userID,
			Action:    "UPDATE_DETEKSI",
			Changes:   changes,
			Timestamp: time.Now(),
		}
		_ = s.producer.PublishIkasAuditLog(context.Background(), auditEvent)
	}

	if err := s.producer.PublishJawabanDeteksiUpdated(context.Background(), event); err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}

func (s *JawabanDeteksiService) Delete(id int, userID string, userRole string, userPerusahaanID string) error {
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

	// Fetch ikas to check ownership
	ikasData, err := s.ikasRepo.GetByID(existing.IkasID)
	if err != nil {
		return errors.New("gagal memverifikasi kepemilikan asesmen")
	}

	if userRole != "admin" && ikasData.Perusahaan.ID != userPerusahaanID {
		return errors.New("anda tidak memiliki akses untuk menghapus data ini")
	}

	if ikasData.IsValidated {
		return errors.New("data asesmen ini sudah divalidasi dan tidak dapat dihapus")
	}

	// Publish Delete Event (Pola 2)
	event := dto_event.JawabanDeteksiDeletedEvent{
		ID:        id,
		IkasID:    existing.IkasID,
		DeletedAt: time.Now(),
	}

	if s.producer != nil {
		auditEvent := dto_event.IkasAuditLogEvent{
			IkasID:    existing.IkasID,
			UserID:    userID,
			Action:    "DELETE_DETEKSI",
			Changes:   map[string]interface{}{"pertanyaan_id": existing.PertanyaanDeteksi.ID, "status": "deleted"},
			Timestamp: time.Now(),
		}
		_ = s.producer.PublishIkasAuditLog(context.Background(), auditEvent)
	}

	if err := s.producer.PublishJawabanDeteksiDeleted(context.Background(), event); err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}
