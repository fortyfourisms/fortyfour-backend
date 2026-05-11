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

type JawabanGulihProducerInterface interface {
	PublishJawabanGulihCreated(ctx context.Context, event interface{}) error
	PublishJawabanGulihUpdated(ctx context.Context, event interface{}) error
	PublishJawabanGulihDeleted(ctx context.Context, event interface{}) error
	PublishIkasAuditLog(ctx context.Context, event interface{}) error
}

type JawabanGulihService struct {
	repo           repository.JawabanGulihRepositoryInterface
	ikasRepo       repository.IkasRepositoryInterface
	pertanyaanRepo repository.PertanyaanGulihRepositoryInterface
	producer       JawabanGulihProducerInterface
	ikasSvc        *IkasService
	cache          cache.RedisInterface
}

func NewJawabanGulihService(
	repo repository.JawabanGulihRepositoryInterface,
	ikasRepo repository.IkasRepositoryInterface,
	pertanyaanRepo repository.PertanyaanGulihRepositoryInterface,
	producer JawabanGulihProducerInterface,
	ikasSvc *IkasService,
	cache cache.RedisInterface,
) *JawabanGulihService {
	return &JawabanGulihService{
		repo:           repo,
		ikasRepo:       ikasRepo,
		pertanyaanRepo: pertanyaanRepo,
		producer:       producer,
		ikasSvc:        ikasSvc,
		cache:          cache,
	}
}

var validValidasiGulih = map[string]bool{"yes": true, "no": true}

func (s *JawabanGulihService) validateCreate(req *dto.CreateJawabanGulihRequest, userRole string) error {
	if req.PertanyaanGulihID <= 0 {
		return errors.New("pertanyaan_gulih_id tidak valid")
	}

	req.IkasID = utils.NormalizeInput(req.IkasID)
	if req.IkasID == "" {
		return errors.New("ikas_id tidak boleh kosong")
	}
	if !utils.IsValidUUID(req.IkasID) {
		return errors.New("format ikas_id tidak valid")
	}

	if req.JawabanGulih == nil {
		return errors.New("jawaban_gulih tidak boleh kosong")
	}
	if *req.JawabanGulih < 0 || *req.JawabanGulih > 5 {
		return errors.New("jawaban_gulih harus bernilai antara 0 sampai 5")
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
		if !validValidasiGulih[*req.Validasi] {
			return errors.New("validasi hanya boleh berisi 'yes' atau 'no'")
		}
	}

	return nil
}

func (s *JawabanGulihService) validateUpdate(req *dto.UpdateJawabanGulihRequest, existingEvidence *string, userRole string) error {
	if req.JawabanGulih != nil && (*req.JawabanGulih < 0 || *req.JawabanGulih > 5) {
		return errors.New("jawaban_gulih harus bernilai antara 0 sampai 5, atau null for N/A")
	}

	// Restricted fields for non-admins
	if userRole != "admin" && userRole != "staff" {
		if req.Validasi != nil || (req.Keterangan != nil && utils.NormalizeInput(*req.Keterangan) != "") {
			return errors.New("hanya admin yang dapat mengubah field validasi dan keterangan")
		}
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

func (s *JawabanGulihService) Create(req dto.CreateJawabanGulihRequest, userRole string, userPerusahaanID string) (string, error) {
	if err := s.validateCreate(&req, userRole); err != nil {
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

	ikasExists, err := s.repo.CheckIkasExists(req.IkasID)
	if err != nil {
		rollbar.Error(err)
		return "", err
	}
	if !ikasExists {
		return "", errors.New("ikas_id tidak ditemukan")
	}

	// VALIDASI KEPEMILIKAN
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
	isDuplicate, err := s.repo.CheckDuplicate(req.IkasID, req.PertanyaanGulihID, 0)
	if err != nil {
		rollbar.Error(err)
		return "", err
	}
	if isDuplicate {
		return "", errors.New("pertanyaan ini sudah pernah diisi untuk asesmen ini")
	}

	// Publish to RabbitMQ for Pola 2
	if err := s.producer.PublishJawabanGulihCreated(context.Background(), req); err != nil {
		rollbar.Error(err)
		return "", err
	}

	if s.cache != nil {
		s.cache.Delete(fmt.Sprintf("%s%s", cache.CacheKeyPrefixJawabanGulih, req.IkasID))
		s.cache.Delete(fmt.Sprintf("%s%s", cache.CacheKeyPrefixGulih, req.IkasID))
	}

	return "Berhasil menyimpan data", nil
}

func (s *JawabanGulihService) GetAll(userRole string) ([]dto.JawabanGulihResponse, error) {
	if userRole != "admin" && userRole != "staff" {
		return nil, errors.New("anda tidak memiliki akses untuk melihat semua data")
	}
	return s.repo.GetAll()
}

func (s *JawabanGulihService) GetByID(id int, userRole string, userPerusahaanID string) (*dto.JawabanGulihResponse, error) {
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

	if userRole != "admin" && userRole != "staff" && ikasData.Perusahaan.ID != userPerusahaanID {
		return nil, errors.New("anda tidak memiliki akses ke data ini")
	}

	return data, nil
}

func (s *JawabanGulihService) GetByIkasID(ikasID string, userRole string, userPerusahaanID string) (*dto.UnifiedJawabanGulihResponse, error) {
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

	cacheKey := fmt.Sprintf("%s%s", cache.CacheKeyPrefixJawabanGulih, ikasID)
	if s.cache != nil {
		cachedData, err := s.cache.Get(cacheKey)
		if err == nil && cachedData != "" {
			var data dto.UnifiedJawabanGulihResponse
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
			// Combine data
			exists := make(map[int]bool)
			for _, d := range data {
				exists[d.PertanyaanGulih.ID] = true
			}
			for _, b := range bufferData {
				if !exists[b.PertanyaanGulih.ID] {
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

	response := &dto.UnifiedJawabanGulihResponse{
		Data:                 data,
		Count:                len(data),
		IsDraft:              isDraft,
		CompletionPercentage: completionPercentage,
	}

	if s.cache != nil {
		go func(key string, dataToCache *dto.UnifiedJawabanGulihResponse) {
			jsonData, err := json.Marshal(dataToCache)
			if err == nil {
				_ = s.cache.Set(key, string(jsonData), cache.DefaultCacheExpiration)
			}
		}(cacheKey, response)
	}

	return response, nil
}

func (s *JawabanGulihService) GetByPerusahaanID(perusahaanID string, userRole string, userPerusahaanID string) ([]dto.JawabanGulihResponse, error) {
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
func (s *JawabanGulihService) GetByPertanyaan(pertanyaanID int, userRole string, userPerusahaanID string) ([]dto.JawabanGulihResponse, error) {
	if pertanyaanID <= 0 {
		return nil, errors.New("pertanyaan_gulih_id tidak valid")
	}
	if userRole != "admin" && userRole != "staff" {
		return s.repo.GetByPertanyaanAndPerusahaan(pertanyaanID, userPerusahaanID)
	}
	return s.repo.GetByPertanyaan(pertanyaanID)
}

func (s *JawabanGulihService) Update(id int, req dto.UpdateJawabanGulihRequest, userID string, userRole string, userPerusahaanID string) (int, string, error) {
	if id <= 0 {
		return 0, "", errors.New("format ID tidak valid")
	}

	existing, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", errors.New("data tidak ditemukan")
		}
		return 0, "", err
	}

	ikasData, err := s.ikasRepo.GetByID(existing.IkasID)
	if err != nil {
		return 0, "", errors.New("gagal memverifikasi kepemilikan asesmen")
	}

	if ikasData.IsValidated {
		return 0, "", errors.New("data asesmen ini sudah divalidasi dan tidak dapat diubah")
	}

	msg := "Berhasil menyimpan data"

	// ---------------- AUTO CLONE/CARRY OVER ----------------
	if s.ikasSvc != nil {
		newIkasID, carryErr := s.ikasSvc.TriggerCarryOverIfNeeded(context.Background(), existing.IkasID, true)
		if carryErr == nil && newIkasID != existing.IkasID {
			newJawabanID, errId := s.repo.GetIDByIkasAndPertanyaan(newIkasID, existing.PertanyaanGulih.ID)
			if errId == nil && newJawabanID > 0 {
				id = newJawabanID
				msg = "Data tahun lalu tidak dapat diubah. Update telah dialihkan otomatis ke record tahun berjalan (carry-over)."

				existing, err = s.repo.GetByID(id)
				if err != nil {
					return 0, "", err
				}

				ikasData, err = s.ikasRepo.GetByID(existing.IkasID)
				if err != nil {
					return 0, "", errors.New("gagal memverifikasi kepemilikan asesmen setelah clone")
				}
			}
		}
	}
	// -------------------------------------------------------

	if userRole != "admin" && userRole != "staff" && ikasData.Perusahaan.ID != userPerusahaanID {
		return 0, "", errors.New("anda tidak memiliki akses untuk mengubah data ini")
	}

	if err := s.validateUpdate(&req, existing.Evidence, userRole); err != nil {
		return 0, "", err
	}

	// Publish Update Event (Pola 2)
	event := dto_event.JawabanGulihUpdatedEvent{
		ID:        id,
		Request:   req,
		UpdatedAt: time.Now(),
	}

	// Change detection for audit log
	changes := make(map[string]interface{})
	if req.JawabanGulih != nil && (existing.JawabanGulih == nil || *req.JawabanGulih != *existing.JawabanGulih) {
		oldVal := interface{}(nil)
		if existing.JawabanGulih != nil {
			oldVal = *existing.JawabanGulih
		}
		changes["jawaban_gulih"] = map[string]interface{}{"old": oldVal, "new": *req.JawabanGulih}
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
				"id":   existing.PertanyaanGulih.ID,
				"teks": existing.PertanyaanGulih.PertanyaanGulih,
			},
			Diff: changes,
		}
		changesJSON, _ := json.Marshal(payload)
		auditEvent := dto_event.IkasAuditLogEvent{
			IkasID:    existing.IkasID,
			UserID:    userID,
			Action:    "UPDATE_GULIH",
			Changes:   changesJSON,
			Timestamp: time.Now(),
		}
		_ = s.producer.PublishIkasAuditLog(context.Background(), auditEvent)
	}

	if err := s.producer.PublishJawabanGulihUpdated(context.Background(), event); err != nil {
		rollbar.Error(err)
		return 0, "", err
	}

	if s.cache != nil {
		s.cache.Delete(fmt.Sprintf("%s%s", cache.CacheKeyPrefixJawabanGulih, existing.IkasID))
		s.cache.Delete(fmt.Sprintf("%s%s", cache.CacheKeyPrefixGulih, existing.IkasID))
	}

	return id, msg, nil
}

func (s *JawabanGulihService) Delete(id int, userID string, userRole string, userPerusahaanID string) error {
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
	event := dto_event.JawabanGulihDeletedEvent{
		ID:        id,
		IkasID:    existing.IkasID,
		DeletedAt: time.Now(),
	}

	if s.producer != nil {
		payload := struct {
			Pertanyaan interface{} `json:"pertanyaan"`
			Status     string      `json:"status"`
		}{
			Pertanyaan: map[string]interface{}{
				"id":   existing.PertanyaanGulih.ID,
				"teks": existing.PertanyaanGulih.PertanyaanGulih,
			},
			Status: "deleted",
		}
		changesJSON, _ := json.Marshal(payload)
		auditEvent := dto_event.IkasAuditLogEvent{
			IkasID:    existing.IkasID,
			UserID:    userID,
			Action:    "DELETE_GULIH",
			Changes:   changesJSON,
			Timestamp: time.Now(),
		}
		_ = s.producer.PublishIkasAuditLog(context.Background(), auditEvent)
	}

	if err := s.producer.PublishJawabanGulihDeleted(context.Background(), event); err != nil {
		rollbar.Error(err)
		return err
	}

	if s.cache != nil {
		s.cache.Delete(fmt.Sprintf("%s%s", cache.CacheKeyPrefixJawabanGulih, existing.IkasID))
		s.cache.Delete(fmt.Sprintf("%s%s", cache.CacheKeyPrefixGulih, existing.IkasID))
	}

	return nil
}
