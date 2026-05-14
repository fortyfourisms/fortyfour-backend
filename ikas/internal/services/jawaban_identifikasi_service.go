package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"ikas/internal/dto"
	"ikas/internal/dto/dto_event"
	"ikas/internal/repository"
	"ikas/internal/utils"
	"ikas/pkg/cache"
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
	repo           repository.JawabanIdentifikasiRepositoryInterface
	ikasRepo       repository.IkasRepositoryInterface
	pertanyaanRepo repository.PertanyaanIdentifikasiRepositoryInterface
	producer       JawabanIdentifikasiProducerInterface
	ikasSvc        *IkasService
	cache          cache.RedisInterface
}

func NewJawabanIdentifikasiService(
	repo repository.JawabanIdentifikasiRepositoryInterface,
	ikasRepo repository.IkasRepositoryInterface,
	pertanyaanRepo repository.PertanyaanIdentifikasiRepositoryInterface,
	producer JawabanIdentifikasiProducerInterface,
	ikasSvc *IkasService,
	cache cache.RedisInterface,
) *JawabanIdentifikasiService {
	return &JawabanIdentifikasiService{
		repo:           repo,
		ikasRepo:       ikasRepo,
		pertanyaanRepo: pertanyaanRepo,
		producer:       producer,
		ikasSvc:        ikasSvc,
		cache:          cache,
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
	if userRole != "admin" && userRole != "staff" {
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
	if userRole != "admin" && userRole != "staff" {
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
	if userRole != "admin" && userRole != "staff" {
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

	if s.cache != nil {
		s.cache.Delete(fmt.Sprintf("%s%s", cache.CacheKeyPrefixJawabanIdentifikasi, req.IkasID))
		s.cache.Delete(fmt.Sprintf("%s%s", cache.CacheKeyPrefixIdentifikasi, req.IkasID))
	}

	return "Berhasil menyimpan data", nil
}

func (s *JawabanIdentifikasiService) GetAll(userRole string) ([]dto.JawabanIdentifikasiResponse, error) {
	if userRole != "admin" && userRole != "staff" {
		return nil, errors.New("anda tidak memiliki akses untuk melihat semua data")
	}
	return s.repo.GetAll()
}

func (s *JawabanIdentifikasiService) GetByUUID(uuid string, userRole string, userPerusahaanID string) (*dto.JawabanIdentifikasiResponse, error) {
	if uuid == "" {
		return nil, errors.New("format ID tidak valid")
	}

	data, err := s.repo.GetByUUID(uuid)
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

	if userRole != "admin" && userRole != "staff" && ikasData.Perusahaan.ID != userPerusahaanID {
		return nil, errors.New("anda tidak memiliki akses ke data ini")
	}

	return data, nil
}

func (s *JawabanIdentifikasiService) GetByIkasID(ikasID string, userRole string, userPerusahaanID string) (*dto.UnifiedJawabanIdentifikasiResponse, error) {
	if !utils.IsValidUUID(ikasID) {
		return nil, errors.New("format ikas_id tidak valid")
	}

	if userRole != "admin" && userRole != "staff" {
		owned, err := s.ikasRepo.CheckOwnership(ikasID, userPerusahaanID)
		if err != nil {
			return nil, err
		}
		if !owned {
			return nil, errors.New("anda tidak memiliki akses ke data asesmen ini")
		}
	}

	cacheKey := fmt.Sprintf("%s%s", cache.CacheKeyPrefixJawabanIdentifikasi, ikasID)
	if s.cache != nil {
		cachedData, err := s.cache.Get(cacheKey)
		if err == nil && cachedData != "" {
			var data dto.UnifiedJawabanIdentifikasiResponse
			if err := json.Unmarshal([]byte(cachedData), &data); err == nil {
				return &data, nil
			}
		}
	}

	// 1. Get total questions
	totalPertanyaan, err := s.pertanyaanRepo.GetTotalCount()
	if err != nil {
		return nil, err
	}

	// 2. Get from main table
	data, err := s.repo.GetByIkasID(ikasID)
	if err != nil {
		return nil, err
	}

	isDraft := false
	if len(data) < totalPertanyaan {
		// 3. If incomplete, check buffer
		bufferData, err := s.repo.GetByIkasIDFromBuffer(ikasID)
		if err == nil && len(bufferData) > 0 {
			// Combine data (prefer main table if duplicate, though shouldn't happen)
			exists := make(map[int]bool)
			for _, d := range data {
				exists[d.PertanyaanIdentifikasi.ID] = true
			}
			for _, b := range bufferData {
				if !exists[b.PertanyaanIdentifikasi.ID] {
					data = append(data, b)
				}
			}
			isDraft = true
		}
	}

	// Calculate completion percentage
	completionPercentage := 0.0
	if totalPertanyaan > 0 {
		completionPercentage = utils.RoundToTwo((float64(len(data)) / float64(totalPertanyaan)) * 100.0)
	}

	response := &dto.UnifiedJawabanIdentifikasiResponse{
		Data:                 data,
		Count:                len(data),
		IsDraft:              isDraft,
		CompletionPercentage: completionPercentage,
	}

	if s.cache != nil {
		go func(key string, dataToCache *dto.UnifiedJawabanIdentifikasiResponse) {
			jsonData, err := json.Marshal(dataToCache)
			if err == nil {
				_ = s.cache.Set(key, string(jsonData), cache.DefaultCacheExpiration)
			}
		}(cacheKey, response)
	}

	return response, nil
}

func (s *JawabanIdentifikasiService) GetByPerusahaanID(perusahaanID string, userRole string, userPerusahaanID string) ([]dto.JawabanIdentifikasiResponse, error) {
	if userRole != "admin" && userRole != "staff" {
		if perusahaanID != userPerusahaanID {
			return nil, errors.New("anda tidak memiliki akses ke data perusahaan ini")
		}
	}
	return s.repo.GetByPerusahaanID(perusahaanID)
}

// GetByPertanyaan retrieves answers for a given question.
// For admin/staff, returns all answers across companies.
// For regular users, enforces company-scoped filtering.
func (s *JawabanIdentifikasiService) GetByPertanyaan(pertanyaanID int, userRole string, userPerusahaanID string) ([]dto.JawabanIdentifikasiResponse, error) {
	if pertanyaanID <= 0 {
		return nil, errors.New("pertanyaan_identifikasi_id tidak valid")
	}
	if userRole != "admin" && userRole != "staff" {
		return s.repo.GetByPertanyaanAndPerusahaan(pertanyaanID, userPerusahaanID)
	}
	return s.repo.GetByPertanyaan(pertanyaanID)
}

func (s *JawabanIdentifikasiService) Update(uuid string, req dto.UpdateJawabanIdentifikasiRequest, userID string, userRole string, userPerusahaanID string) (string, string, error) {
	if uuid == "" {
		return "", "", errors.New("format ID tidak valid")
	}

	existing, err := s.repo.GetByUUID(uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", errors.New("data tidak ditemukan")
		}
		return "", "", err
	}

	ikasData, err := s.ikasRepo.GetByID(existing.IkasID)
	if err != nil {
		return "", "", errors.New("gagal memverifikasi kepemilikan asesmen")
	}

	if ikasData.IsValidated {
		return "", "", errors.New("data asesmen ini sudah divalidasi dan tidak dapat diubah")
	}

	msg := "Berhasil menyimpan data"

	// ---------------- AUTO CLONE/CARRY OVER ----------------
	if s.ikasSvc != nil {
		newIkasID, carryErr := s.ikasSvc.TriggerCarryOverIfNeeded(context.Background(), existing.IkasID, true)
		if carryErr == nil && newIkasID != existing.IkasID {
			// Find the UUID of the cloned record
			clonedData, errCloned := s.repo.GetByIkasID(newIkasID)
			if errCloned == nil {
				for _, d := range clonedData {
					if d.PertanyaanIdentifikasi.ID == existing.PertanyaanIdentifikasi.ID {
						uuid = d.UUID
						msg = "Data tahun lalu tidak dapat diubah. Update telah dialihkan otomatis ke record tahun berjalan (carry-over)."

						existing, err = s.repo.GetByUUID(uuid)
						if err != nil {
							return "", "", err
						}

						ikasData, err = s.ikasRepo.GetByID(existing.IkasID)
						if err != nil {
							return "", "", errors.New("gagal memverifikasi kepemilikan asesmen setelah clone")
						}
						break
					}
				}
			}
		}
	}
	// -------------------------------------------------------

	if userRole != "admin" && userRole != "staff" && ikasData.Perusahaan.ID != userPerusahaanID {
		return "", "", errors.New("anda tidak memiliki akses untuk mengubah data ini")
	}

	if err := s.validateUpdate(&req, existing.Evidence, userRole); err != nil {
		return "", "", err
	}

	// Publish Update Event (Pola 2)
	event := dto_event.JawabanIdentifikasiUpdatedEvent{
		UUID:      uuid,
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
		payload := struct {
			Pertanyaan interface{} `json:"pertanyaan"`
			Diff       interface{} `json:"diff"`
		}{
			Pertanyaan: map[string]interface{}{
				"id":   existing.PertanyaanIdentifikasi.ID,
				"teks": existing.PertanyaanIdentifikasi.PertanyaanIdentifikasi,
			},
			Diff: changes,
		}
		changesJSON, _ := json.Marshal(payload)
		auditEvent := dto_event.IkasAuditLogEvent{
			IkasID:    existing.IkasID,
			UserID:    userID,
			Action:    "UPDATE_IDENTIFIKASI",
			Changes:   changesJSON,
			Timestamp: time.Now(),
		}
		_ = s.producer.PublishIkasAuditLog(context.Background(), auditEvent)
	}

	if err := s.producer.PublishJawabanIdentifikasiUpdated(context.Background(), event); err != nil {
		rollbar.Error(err)
		return "", "", err
	}

	if s.cache != nil {
		s.cache.Delete(fmt.Sprintf("%s%s", cache.CacheKeyPrefixJawabanIdentifikasi, existing.IkasID))
		s.cache.Delete(fmt.Sprintf("%s%s", cache.CacheKeyPrefixIdentifikasi, existing.IkasID))
	}

	return uuid, msg, nil
}

func (s *JawabanIdentifikasiService) Delete(uuid string, userID string, userRole string, userPerusahaanID string) error {
	if uuid == "" {
		return errors.New("format ID tidak valid")
	}

	// Existence Check
	existing, err := s.repo.GetByUUID(uuid)
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
		UUID:      uuid,
		IkasID:    existing.IkasID,
		DeletedAt: time.Now(),
	}

	if s.producer != nil {
		payload := struct {
			Pertanyaan interface{} `json:"pertanyaan"`
			Status     string      `json:"status"`
		}{
			Pertanyaan: map[string]interface{}{
				"id":   existing.PertanyaanIdentifikasi.ID,
				"teks": existing.PertanyaanIdentifikasi.PertanyaanIdentifikasi,
			},
			Status: "deleted",
		}
		changesJSON, _ := json.Marshal(payload)
		auditEvent := dto_event.IkasAuditLogEvent{
			IkasID:    existing.IkasID,
			UserID:    userID,
			Action:    "DELETE_IDENTIFIKASI",
			Changes:   changesJSON,
			Timestamp: time.Now(),
		}
		_ = s.producer.PublishIkasAuditLog(context.Background(), auditEvent)
	}

	if err := s.producer.PublishJawabanIdentifikasiDeleted(context.Background(), event); err != nil {
		rollbar.Error(err)
		return err
	}

	if s.cache != nil {
		s.cache.Delete(fmt.Sprintf("%s%s", cache.CacheKeyPrefixJawabanIdentifikasi, existing.IkasID))
		s.cache.Delete(fmt.Sprintf("%s%s", cache.CacheKeyPrefixIdentifikasi, existing.IkasID))
	}

	return nil
}
