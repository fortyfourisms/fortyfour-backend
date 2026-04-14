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

type JawabanIdentifikasiProducerInterface interface {
	PublishJawabanIdentifikasiCreated(ctx context.Context, event interface{}) error
	PublishJawabanIdentifikasiUpdated(ctx context.Context, event interface{}) error
	PublishJawabanIdentifikasiDeleted(ctx context.Context, event interface{}) error
	PublishIkasAuditLog(ctx context.Context, event interface{}) error
}

type JawabanIdentifikasiService struct {
	repo     repository.JawabanIdentifikasiRepositoryInterface
	ikasRepo repository.IkasRepositoryInterface
	producer JawabanIdentifikasiProducerInterface
}

func NewJawabanIdentifikasiService(
	repo repository.JawabanIdentifikasiRepositoryInterface,
	ikasRepo repository.IkasRepositoryInterface,
	producer JawabanIdentifikasiProducerInterface,
) *JawabanIdentifikasiService {
	return &JawabanIdentifikasiService{
		repo:     repo,
		ikasRepo: ikasRepo,
		producer: producer,
	}
}

var validValidasi = map[string]bool{"yes": true, "no": true}

func (s *JawabanIdentifikasiService) validateCreate(req *dto.CreateJawabanIdentifikasiRequest, userRole string) error {
	if req.PertanyaanIdentifikasiID <= 0 {
		return errors.New("pertanyaan_identifikasi_id tidak valid")
	}

	req.IkasID = utils.NormalizeInput(req.IkasID)
	if req.IkasID == "" {
		return errors.New("ikas_id tidak boleh kosong")
	}
	if !utils.IsValidUUID(req.IkasID) {
		return errors.New("format ikas_id tidak valid")
	}

	// Jawaban must be provided and between 0.00 - 5.00
	if req.JawabanIdentifikasi == nil {
		return errors.New("jawaban_identifikasi tidak boleh kosong")
	}
	if *req.JawabanIdentifikasi < 0 || *req.JawabanIdentifikasi > 5 {
		return errors.New("jawaban_identifikasi harus bernilai antara 0 sampai 5")
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
		if !validValidasi[*req.Validasi] {
			return errors.New("validasi hanya boleh berisi 'yes' atau 'no'")
		}
	}

	return nil
}

func (s *JawabanIdentifikasiService) validateUpdate(req *dto.UpdateJawabanIdentifikasiRequest, existingEvidence *string, userRole string) error {
	// null = N/A (diperbolehkan), tapi jika diisi harus 0.00 - 5.00
	if req.JawabanIdentifikasi != nil && (*req.JawabanIdentifikasi < 0 || *req.JawabanIdentifikasi > 5) {
		return errors.New("jawaban_identifikasi harus bernilai antara 0 sampai 5, atau null untuk N/A")
	}

	// Restricted fields for non-admins
	if userRole != "admin" {
		if req.Validasi != nil || (req.Keterangan != nil && utils.NormalizeInput(*req.Keterangan) != "") {
			return errors.New("hanya admin yang dapat mengubah field validasi dan keterangan")
		}
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

func (s *JawabanIdentifikasiService) Create(req dto.CreateJawabanIdentifikasiRequest, userRole string, userPerusahaanID string) (string, error) {
	if err := s.validateCreate(&req, userRole); err != nil {
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

	ikasExists, err := s.repo.CheckIkasExists(req.IkasID)
	if err != nil {
		rollbar.Error(err)
		return "", err
	}
	if !ikasExists {
		return "", errors.New("ikas_id tidak ditemukan")
	}

	// VALIDASI KEPEMILIKAN: Hanya pemegang IKAS yang bisa mengisi jawaban (atau admin)
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
	// Check if already exists in the MAIN table
	isDuplicate, err := s.repo.CheckDuplicate(req.IkasID, req.PertanyaanIdentifikasiID, 0)
	if err != nil {
		rollbar.Error(err)
		return "", err
	}
	if isDuplicate {
		return "", errors.New("pertanyaan ini sudah pernah diisi untuk asesmen ini")
	}

	// Publish to RabbitMQ for Pola 2
	if err := s.producer.PublishJawabanIdentifikasiCreated(context.Background(), req); err != nil {
		rollbar.Error(err)
		return "", err
	}

	return "Berhasil menyimpan data", nil
}

func (s *JawabanIdentifikasiService) GetAll(userRole string) ([]dto.JawabanIdentifikasiResponse, error) {
	if userRole != "admin" {
		return nil, errors.New("anda tidak memiliki akses untuk melihat semua data")
	}
	return s.repo.GetAll()
}


func (s *JawabanIdentifikasiService) GetByID(id int, userRole string, userPerusahaanID string) (*dto.JawabanIdentifikasiResponse, error) {
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

	// Note: checking access now by comparing ikas session rather than perusahaan_id directly
	// but since we don't have the link here easily without fetching ikas, 
	// we might need a more robust check if needed.
	// For now, let's assume if they have the ikas_id they might have access, or fetch ikas to check company.
	
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

func (s *JawabanIdentifikasiService) GetByIkasID(ikasID string, userRole string, userPerusahaanID string) ([]dto.JawabanIdentifikasiResponse, error) {
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

func (s *JawabanIdentifikasiService) GetByPerusahaanID(perusahaanID string, userRole string, userPerusahaanID string) ([]dto.JawabanIdentifikasiResponse, error) {
	if userRole != "admin" {
		if perusahaanID != userPerusahaanID {
			return nil, errors.New("anda tidak memiliki akses ke data perusahaan ini")
		}
	}
	return s.repo.GetByPerusahaanID(perusahaanID)
}

func (s *JawabanIdentifikasiService) GetByPertanyaan(pertanyaanID int) ([]dto.JawabanIdentifikasiResponse, error) {
	if pertanyaanID <= 0 {
		return nil, errors.New("pertanyaan_identifikasi_id tidak valid")
	}
	return s.repo.GetByPertanyaan(pertanyaanID)
}

func (s *JawabanIdentifikasiService) Update(id int, req dto.UpdateJawabanIdentifikasiRequest, userID string, userRole string, userPerusahaanID string) error {
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
		auditEvent := dto_event.IkasAuditLogEvent{
			IkasID:    existing.IkasID,
			UserID:    userID,
			Action:    "UPDATE_IDENTIFIKASI",
			Changes:   changes,
			Timestamp: time.Now(),
		}
		_ = s.producer.PublishIkasAuditLog(context.Background(), auditEvent)
	}

	if err := s.producer.PublishJawabanIdentifikasiUpdated(context.Background(), event); err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}

func (s *JawabanIdentifikasiService) Delete(id int, userID string, userRole string, userPerusahaanID string) error {
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
	event := dto_event.JawabanIdentifikasiDeletedEvent{
		ID:        id,
		IkasID:    existing.IkasID,
		DeletedAt: time.Now(),
	}

	if s.producer != nil {
		auditEvent := dto_event.IkasAuditLogEvent{
			IkasID:    existing.IkasID,
			UserID:    userID,
			Action:    "DELETE_IDENTIFIKASI",
			Changes:   map[string]interface{}{"pertanyaan_id": existing.PertanyaanIdentifikasi.ID, "status": "deleted"},
			Timestamp: time.Now(),
		}
		_ = s.producer.PublishIkasAuditLog(context.Background(), auditEvent)
	}

	if err := s.producer.PublishJawabanIdentifikasiDeleted(context.Background(), event); err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}
