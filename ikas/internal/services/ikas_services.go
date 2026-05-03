package services

import (
	"context"
	"encoding/json"
	"fmt"
	"ikas/internal/dto"
	"ikas/internal/dto/dto_event"
	"ikas/internal/repository"
	"ikas/internal/utils"
	"ikas/pkg/cache"
	"time"

	"github.com/google/uuid"
	"github.com/rollbar/rollbar-go"
)

type IkasProducerInterface interface {
	PublishIkasCreated(ctx context.Context, event interface{}) error
	PublishIkasUpdated(ctx context.Context, event interface{}) error
	PublishIkasDeleted(ctx context.Context, event interface{}) error
	PublishIkasAuditLog(ctx context.Context, event interface{}) error
	PublishIkasImported(ctx context.Context, event interface{}) error
	PublishJawabanIdentifikasiCreated(ctx context.Context, event interface{}) error
	PublishJawabanProteksiCreated(ctx context.Context, event interface{}) error
	PublishJawabanDeteksiCreated(ctx context.Context, event interface{}) error
	PublishJawabanGulihCreated(ctx context.Context, event interface{}) error
	PublishIkasEditRequested(ctx context.Context, event interface{}) error
	PublishIkasEditActioned(ctx context.Context, event interface{}) error
}

type IkasService struct {
	repo             repository.IkasRepositoryInterface
	identifikasiRepo repository.IdentifikasiRepositoryInterface
	proteksiRepo     repository.ProteksiRepositoryInterface
	deteksiRepo      repository.DeteksiRepositoryInterface
	gulihRepo        repository.GulihRepositoryInterface
	jawabanIdenRepo  repository.JawabanIdentifikasiRepositoryInterface
	jawabanProtRepo  repository.JawabanProteksiRepositoryInterface
	jawabanDetRepo   repository.JawabanDeteksiRepositoryInterface
	jawabanGulihRepo repository.JawabanGulihRepositoryInterface
	producer         IkasProducerInterface
	cache            cache.RedisInterface
}

func NewIkasService(
	repo repository.IkasRepositoryInterface,
	identifikasiRepo repository.IdentifikasiRepositoryInterface,
	proteksiRepo repository.ProteksiRepositoryInterface,
	deteksiRepo repository.DeteksiRepositoryInterface,
	gulihRepo repository.GulihRepositoryInterface,
	jawabanIdenRepo repository.JawabanIdentifikasiRepositoryInterface,
	jawabanProtRepo repository.JawabanProteksiRepositoryInterface,
	jawabanDetRepo repository.JawabanDeteksiRepositoryInterface,
	jawabanGulihRepo repository.JawabanGulihRepositoryInterface,
	producer IkasProducerInterface,
	cache cache.RedisInterface,
) *IkasService {
	return &IkasService{
		repo:             repo,
		identifikasiRepo: identifikasiRepo,
		proteksiRepo:     proteksiRepo,
		deteksiRepo:      deteksiRepo,
		gulihRepo:        gulihRepo,
		jawabanIdenRepo:  jawabanIdenRepo,
		jawabanProtRepo:  jawabanProtRepo,
		jawabanDetRepo:   jawabanDetRepo,
		jawabanGulihRepo: jawabanGulihRepo,
		producer:         producer,
		cache:            cache,
	}
}

func (s *IkasService) validateCreate(req *dto.CreateIkasRequest) error {
	req.IDPerusahaan = utils.NormalizeInput(req.IDPerusahaan)
	if req.IDPerusahaan == "" {
		return fmt.Errorf("id_perusahaan tidak boleh kosong")
	}
	if !utils.IsValidUUID(req.IDPerusahaan) {
		return fmt.Errorf("format id_perusahaan tidak valid")
	}

	req.Responden = utils.NormalizeInput(req.Responden)
	if req.Responden == "" {
		return fmt.Errorf("responden tidak boleh kosong")
	}

	req.Telepon = utils.NormalizeInput(req.Telepon)
	if req.Telepon == "" {
		return fmt.Errorf("telepon tidak boleh kosong")
	}

	req.Jabatan = utils.NormalizeInput(req.Jabatan)
	if req.Jabatan == "" {
		return fmt.Errorf("jabatan tidak boleh kosong")
	}

	if req.TargetNilai < 0 || req.TargetNilai > 5 {
		return fmt.Errorf("target_nilai harus bernilai antara 0 sampai 5")
	}

	if req.Tanggal != "" {
		if _, err := time.Parse("2006-01-02", req.Tanggal); err != nil {
			return fmt.Errorf("format tanggal tidak valid (harus YYYY-MM-DD)")
		}
	}

	return nil
}

func (s *IkasService) validateUpdate(req *dto.UpdateIkasRequest) error {
	if req.IDPerusahaan != nil {
		normalized := utils.NormalizeInput(*req.IDPerusahaan)
		req.IDPerusahaan = &normalized
		if *req.IDPerusahaan == "" {
			return fmt.Errorf("id_perusahaan tidak boleh kosong")
		}
		if !utils.IsValidUUID(*req.IDPerusahaan) {
			return fmt.Errorf("format id_perusahaan tidak valid")
		}
	}

	if req.Responden != nil {
		normalized := utils.NormalizeInput(*req.Responden)
		req.Responden = &normalized
		if *req.Responden == "" {
			return fmt.Errorf("responden tidak boleh kosong")
		}
	}

	if req.Telepon != nil {
		normalized := utils.NormalizeInput(*req.Telepon)
		req.Telepon = &normalized
		if *req.Telepon == "" {
			return fmt.Errorf("telepon tidak boleh kosong")
		}
	}

	if req.Jabatan != nil {
		normalized := utils.NormalizeInput(*req.Jabatan)
		req.Jabatan = &normalized
		if *req.Jabatan == "" {
			return fmt.Errorf("jabatan tidak boleh kosong")
		}
	}

	if req.TargetNilai != nil {
		if *req.TargetNilai < 0 || *req.TargetNilai > 5 {
			return fmt.Errorf("target_nilai harus bernilai antara 0 sampai 5")
		}
	}

	if req.Tanggal != nil && *req.Tanggal != "" {
		if _, err := time.Parse("2006-01-02", *req.Tanggal); err != nil {
			return fmt.Errorf("format tanggal tidak valid (harus YYYY-MM-DD)")
		}
	}

	return nil
}

func (s *IkasService) Create(ctx context.Context, req dto.CreateIkasRequest, id string, userID string, userRole string, userPerusahaanID string) error {
	// 1. Enforce perusahaan ownership for non-admin users
	if userRole != "admin" && userRole != "staff" {
		// User must be affiliated with a company before they can create IKAS
		if userPerusahaanID == "" || userPerusahaanID == "null" {
			return fmt.Errorf("akun Anda belum terasosiasi dengan perusahaan. Silakan lengkapi profil perusahaan terlebih dahulu")
		}
		// User can only create IKAS for their own company
		if req.IDPerusahaan != userPerusahaanID {
			return fmt.Errorf("anda tidak memiliki akses untuk membuat data IKAS atas nama perusahaan lain")
		}
	}

	// 2. Synchronous Validation
	if err := s.validateCreate(&req); err != nil {
		return err
	}

	// 3. Synchronous Existence Check (Perusahaan)
	perusahaanExists, err := s.repo.CheckExistsByPerusahaanID(req.IDPerusahaan)
	if err != nil {
		return err
	}
	if !perusahaanExists {
		return fmt.Errorf("perusahaan dengan ID %s tidak ditemukan", req.IDPerusahaan)
	}

	// Extract year from Tanggal for yearly uniqueness check
	var tahun int
	if t, err := time.Parse("2006-01-02", req.Tanggal); err == nil {
		tahun = t.Year()
	} else {
		// Fallback to current year if format is invalid (should be validated by now)
		tahun = time.Now().Year()
	}

	exists, err := s.repo.CheckExistsByPerusahaanIDAndYear(req.IDPerusahaan, tahun)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("Data IKAS untuk perusahaan ini di tahun %d sudah ada", tahun)
	}

	event := dto_event.IkasCreatedEvent{
		IkasID:          id,
		IDPerusahaan:    req.IDPerusahaan,
		Tanggal:         req.Tanggal,
		Responden:       req.Responden,
		Telepon:         req.Telepon,
		Jabatan:         req.Jabatan,
		TargetNilai:     req.TargetNilai,
		NilaiKematangan: 0.0,
		UserID:          userID,
		CreatedAt:       time.Now(),
	}

	if s.producer == nil {
		return nil
	}

	if err := s.producer.PublishIkasCreated(ctx, event); err != nil {
		return err
	}

	// Audit Log for Create
	auditEvent := dto_event.IkasAuditLogEvent{
		IkasID:    id,
		UserID:    userID,
		Action:    "CREATE_IKAS",
		Changes:   map[string]interface{}{"perusahaan_id": req.IDPerusahaan, "tanggal": req.Tanggal, "responden": req.Responden},
		Timestamp: time.Now(),
	}
	_ = s.producer.PublishIkasAuditLog(ctx, auditEvent)

	// Invalidate Cache
	if s.cache != nil {
		s.cache.Delete(cache.CacheKeyIkasRecords)
		s.cache.Delete(fmt.Sprintf("%s:perusahaan:%s", cache.CacheKeyIkasRecords, req.IDPerusahaan))
	}

	return nil
}

func (s *IkasService) GetAll(userRole string) ([]dto.IkasResponse, error) {
	if userRole != "admin" && userRole != "staff" {
		return nil, fmt.Errorf("anda tidak memiliki akses untuk melihat semua data")
	}

	if s.cache != nil {
		cachedData, err := s.cache.Get(cache.CacheKeyIkasRecords)
		if err == nil && cachedData != "" {
			var result []dto.IkasResponse
			if err := json.Unmarshal([]byte(cachedData), &result); err == nil {
				return result, nil
			}
		}
	}

	data, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	if s.cache != nil && len(data) > 0 {
		go func(dataToCache []dto.IkasResponse) {
			jsonData, err := json.Marshal(dataToCache)
			if err == nil {
				err = s.cache.Set(cache.CacheKeyIkasRecords, string(jsonData), cache.DefaultCacheExpiration)
				if err != nil {
					rollbar.Error(fmt.Errorf("redis caching error: %v", err))
				}
			}
		}(data)
	}

	return data, nil
}

func (s *IkasService) GetByPerusahaan(perusahaanID string) ([]dto.IkasResponse, error) {
	cacheKey := fmt.Sprintf("%s:perusahaan:%s", cache.CacheKeyIkasRecords, perusahaanID)
	if s.cache != nil {
		cachedData, err := s.cache.Get(cacheKey)
		if err == nil && cachedData != "" {
			var result []dto.IkasResponse
			if err := json.Unmarshal([]byte(cachedData), &result); err == nil {
				return result, nil
			}
		}
	}

	data, err := s.repo.GetByPerusahaan(perusahaanID)
	if err != nil {
		return nil, err
	}

	if s.cache != nil && len(data) > 0 {
		go func(dataToCache []dto.IkasResponse, key string) {
			jsonData, err := json.Marshal(dataToCache)
			if err == nil {
				err = s.cache.Set(key, string(jsonData), cache.DefaultCacheExpiration)
				if err != nil {
					rollbar.Error(fmt.Errorf("redis caching error: %v", err))
				}
			}
		}(data, cacheKey)
	}

	return data, nil
}

func (s *IkasService) GetByID(id string, userRole string, userPerusahaanID string) (*dto.IkasResponse, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if userRole != "admin" && userRole != "staff" && data.Perusahaan != nil && data.Perusahaan.ID != userPerusahaanID {
		return nil, fmt.Errorf("anda tidak memiliki akses ke data ini")
	}
	return data, nil
}

func (s *IkasService) Update(ctx context.Context, id string, req dto.UpdateIkasRequest, userID string, userRole string, userPerusahaanID string) (string, error) {
	// 1. Synchronous Validation
	if err := s.validateUpdate(&req); err != nil {
		return id, err
	}

	// 2. Get existing IKAS meta
	current, err := s.repo.GetByID(id)
	if err != nil {
		return id, err
	}

	// 3. Cross-Perusahaan Check (Security)
	if userRole != "admin" && userRole != "staff" && current.Perusahaan != nil && current.Perusahaan.ID != userPerusahaanID {
		return id, fmt.Errorf("anda tidak memiliki akses untuk mengubah data ini")
	}

	// 4. Check if new IDPerusahaan exists (if provided)
	if req.IDPerusahaan != nil && *req.IDPerusahaan != "" && *req.IDPerusahaan != current.Perusahaan.ID {
		exists, err := s.repo.CheckExistsByPerusahaanID(*req.IDPerusahaan)
		if err != nil {
			return id, err
		}
		if !exists {
			return id, fmt.Errorf("perusahaan dengan ID %s tidak ditemukan", *req.IDPerusahaan)
		}
	}

	// 5. Validation Check (LOCKED)
	if current.IsValidated {
		return id, fmt.Errorf("data asesmen ini sudah divalidasi dan tidak dapat diubah")
	}

	newID, err := s.TriggerCarryOverIfNeeded(ctx, id, false)
	if err != nil {
		return id, err
	}
	id = newID

	// 5. Update process on the resolved 'id'
	changes := make(map[string]interface{})
	if req.IDPerusahaan != nil && *req.IDPerusahaan != current.Perusahaan.ID {
		changes["id_perusahaan"] = map[string]interface{}{"old": current.Perusahaan.ID, "new": *req.IDPerusahaan}
	}
	if req.Tanggal != nil && *req.Tanggal != current.Tanggal {
		changes["tanggal"] = map[string]interface{}{"old": current.Tanggal, "new": *req.Tanggal}
	}
	if req.Responden != nil && *req.Responden != current.Responden {
		changes["responden"] = map[string]interface{}{"old": current.Responden, "new": *req.Responden}
	}
	if req.Telepon != nil && *req.Telepon != current.Telepon {
		changes["telepon"] = map[string]interface{}{"old": current.Telepon, "new": *req.Telepon}
	}
	if req.Jabatan != nil && *req.Jabatan != current.Jabatan {
		changes["jabatan"] = map[string]interface{}{"old": current.Jabatan, "new": *req.Jabatan}
	}
	if req.TargetNilai != nil && *req.TargetNilai != current.TargetNilai {
		changes["target_nilai"] = map[string]interface{}{"old": current.TargetNilai, "new": *req.TargetNilai}
	}

	if s.producer == nil {
		return id, nil
	}

	// Publish audit log if there are changes
	if len(changes) > 0 {
		auditEvent := dto_event.IkasAuditLogEvent{
			IkasID:    id,
			UserID:    userID,
			Action:    "UPDATE",
			Changes:   changes,
			Timestamp: time.Now(),
		}
		_ = s.producer.PublishIkasAuditLog(ctx, auditEvent)
	}

	// Publish update event
	event := dto_event.IkasUpdatedEvent{
		IkasID:       id,
		IDPerusahaan: req.IDPerusahaan,
		Tanggal:      req.Tanggal,
		Responden:    req.Responden,
		Telepon:      req.Telepon,
		Jabatan:      req.Jabatan,
		TargetNilai:  req.TargetNilai,
		UserID:       userID,
		UpdatedAt:    time.Now(),
	}

	if err := s.producer.PublishIkasUpdated(ctx, event); err != nil {
		return id, err
	}

	// Invalidate Cache
	if s.cache != nil {
		s.cache.Delete(cache.CacheKeyIkasRecords)
		if current.Perusahaan != nil {
			s.cache.Delete(fmt.Sprintf("%s:perusahaan:%s", cache.CacheKeyIkasRecords, current.Perusahaan.ID))
		}
		if req.IDPerusahaan != nil {
			s.cache.Delete(fmt.Sprintf("%s:perusahaan:%s", cache.CacheKeyIkasRecords, *req.IDPerusahaan))
		}
	}

	return id, nil
}

// TriggerCarryOverIfNeeded detects cross-year edits and creates a new year clone if needed.
// It returns the newly created ID if a clone was made, or the existing ID.
// If isAnswerUpdate is true and a new year record already exists, it forwards to the latest ID so answers update the new year automatically.
func (s *IkasService) TriggerCarryOverIfNeeded(ctx context.Context, id string, isAnswerUpdate bool) (string, error) {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return id, err
	}

	var existingYear int
	dateStr := current.Tanggal
	if len(dateStr) > 10 {
		dateStr = dateStr[:10]
	}
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		existingYear = t.Year()
	}

	if existingYear == 0 && current.CreatedAt != "" {
		caStr := current.CreatedAt
		if len(caStr) > 10 {
			caStr = caStr[:10]
		}
		if t, err := time.Parse("2006-01-02", caStr); err == nil {
			existingYear = t.Year()
		}
	}

	targetYear := time.Now().Year()

	if existingYear < targetYear {
		latest, err := s.repo.GetLatestByPerusahaan(current.Perusahaan.ID)
		if err != nil {
			return id, err
		}

		var latestYear int
		if latest != nil {
			lDate := latest.Tanggal
			if len(lDate) > 10 {
				lDate = lDate[:10]
			}
			if t, err := time.Parse("2006-01-02", lDate); err == nil {
				latestYear = t.Year()
			}
		}

		if latestYear < targetYear {
			newID, err := s.handleCarryOver(ctx, id, targetYear)
			if err != nil {
				return id, fmt.Errorf("gagal melakukan carry-over data: %v", err)
			}
			return newID, nil
		} else if latest != nil && latestYear == targetYear && isAnswerUpdate {
			// A target year record already exists, we return latest.ID!
			// Because for isolated answer updates, we WANT them to automatically route to the newest record.
			return latest.ID, nil
		}
	}

	return id, nil
}

func (s *IkasService) handleCarryOver(ctx context.Context, sourceID string, targetYear int) (string, error) {
	newID := uuid.New().String()
	targetDate := fmt.Sprintf("%d-01-01", targetYear)

	// 1. Phase 1: Create Initial Ikas Record (parents must exist before children)
	if err := s.repo.CreateInitial(sourceID, newID, targetDate); err != nil {
		return "", err
	}

	// 2. Phase 2: Clone Score Summaries (these now reference the existing newID)
	newIden, err := s.identifikasiRepo.CloneByIkasID(sourceID, newID)
	if err != nil {
		return "", err
	}
	newProt, err := s.proteksiRepo.CloneByIkasID(sourceID, newID)
	if err != nil {
		return "", err
	}
	newDet, err := s.deteksiRepo.CloneByIkasID(sourceID, newID)
	if err != nil {
		return "", err
	}
	newGulih, err := s.gulihRepo.CloneByIkasID(sourceID, newID)
	if err != nil {
		return "", err
	}

	// 3. Phase 3: Update Ikas Record with domain score links
	if err := s.repo.UpdateDomainLinks(newID, newIden, newProt, newDet, newGulih); err != nil {
		return "", err
	}

	// 4. Phase 4: Clone all Answers (mass copy)
	if err := s.jawabanIdenRepo.CloneByIkasID(sourceID, newID); err != nil {
		return "", err
	}
	if err := s.jawabanProtRepo.CloneByIkasID(sourceID, newID); err != nil {
		return "", err
	}
	if err := s.jawabanDetRepo.CloneByIkasID(sourceID, newID); err != nil {
		return "", err
	}
	if err := s.jawabanGulihRepo.CloneByIkasID(sourceID, newID); err != nil {
		return "", err
	}

	return newID, nil
}

func (s *IkasService) Delete(ctx context.Context, id string, userID string, userRole string, userPerusahaanID string) error {
	// Check existence
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if userRole != "admin" && existing.Perusahaan != nil && existing.Perusahaan.ID != userPerusahaanID {
		return fmt.Errorf("anda tidak memiliki akses untuk menghapus data ini")
	}

	if existing.IsValidated && userRole != "admin" {
		return fmt.Errorf("data asesmen ini sudah divalidasi dan tidak dapat dihapus")
	}

	// Publish delete event
	event := dto_event.IkasDeletedEvent{
		IkasID:    id,
		UserID:    userID,
		DeletedAt: time.Now(),
	}

	if s.producer == nil {
		return nil
	}

	if err := s.producer.PublishIkasDeleted(ctx, event); err != nil {
		return err
	}

	// Audit Log for Delete
	auditEvent := dto_event.IkasAuditLogEvent{
		IkasID:    id,
		UserID:    userID,
		Action:    "DELETE_IKAS",
		Changes:   map[string]interface{}{"status": "deleted"},
		Timestamp: time.Now(),
	}
	_ = s.producer.PublishIkasAuditLog(ctx, auditEvent)

	// Invalidate Cache
	if s.cache != nil {
		s.cache.Delete(cache.CacheKeyIkasRecords)
		if existing.Perusahaan != nil {
			s.cache.Delete(fmt.Sprintf("%s:perusahaan:%s", cache.CacheKeyIkasRecords, existing.Perusahaan.ID))
		}
	}

	return nil
}

func (s *IkasService) ImportFromExcel(ctx context.Context, fileData []byte, userID string, userRole string, userPerusahaanID string) (string, error) {
	excelData, err := s.repo.ParseExcelForImport(fileData)
	if err != nil {
		return "", err
	}

	newID := uuid.New().String()

	// 1. Create main IKAS record (ownership enforcement delegated to Create)
	if err := s.Create(ctx, excelData.IkasRequest, newID, userID, userRole, userPerusahaanID); err != nil {
		return "", err
	}

	// 2. Publish events for each subdomain to trigger automatic processing
	// pass NEW ID (ikasID) instead of perusahaanID

	// Identifikasi
	for _, ans := range excelData.JawabanIdentifikasi {
		event := dto.CreateJawabanIdentifikasiRequest{
			PertanyaanIdentifikasiID: ans.PertanyaanID,
			IkasID:                   newID,
			JawabanIdentifikasi:      &ans.Jawaban,
		}
		s.producer.PublishJawabanIdentifikasiCreated(context.Background(), event)
	}

	// Proteksi
	for _, ans := range excelData.JawabanProteksi {
		event := dto.CreateJawabanProteksiRequest{
			PertanyaanProteksiID: ans.PertanyaanID,
			IkasID:               newID,
			JawabanProteksi:      &ans.Jawaban,
		}
		s.producer.PublishJawabanProteksiCreated(context.Background(), event)
	}

	// Deteksi
	for _, ans := range excelData.JawabanDeteksi {
		event := dto.CreateJawabanDeteksiRequest{
			PertanyaanDeteksiID: ans.PertanyaanID,
			IkasID:              newID,
			JawabanDeteksi:      &ans.Jawaban,
		}
		s.producer.PublishJawabanDeteksiCreated(context.Background(), event)
	}

	// Gulih
	for _, ans := range excelData.JawabanGulih {
		event := dto.CreateJawabanGulihRequest{
			PertanyaanGulihID: ans.PertanyaanID,
			IkasID:            newID,
			JawabanGulih:      &ans.Jawaban,
		}
		s.producer.PublishJawabanGulihCreated(context.Background(), event)
	}

	return newID, nil
}
func (s *IkasService) ValidateIkas(ctx context.Context, id string, status bool, userID string) error {
	// Existence Check
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if existing.IsValidated == status {
		return nil // Already in target state
	}

	// Update status
	if err := s.repo.UpdateValidationStatus(id, status); err != nil {
		return err
	}

	// Audit Log
	action := "VALIDATE_IKAS"
	if !status {
		action = "UNVALIDATE_IKAS"
	}

	auditEvent := dto_event.IkasAuditLogEvent{
		IkasID:    id,
		UserID:    userID,
		Action:    action,
		Changes:   map[string]interface{}{"is_validated": status},
		Timestamp: time.Now(),
	}
	if s.producer != nil {
		_ = s.producer.PublishIkasAuditLog(ctx, auditEvent)
	}

	// Invalidate Cache
	if s.cache != nil {
		s.cache.Delete(cache.CacheKeyIkasRecords)
		if existing.Perusahaan != nil {
			s.cache.Delete(fmt.Sprintf("%s:perusahaan:%s", cache.CacheKeyIkasRecords, existing.Perusahaan.ID))
		}
	}

	return nil
}
func (s *IkasService) ExportByIDPDF(ctx context.Context, id string, userRole string, userPerusahaanID string) (*dto.IkasResponse, []byte, error) {
	// 1. Get complete data with ownership check
	ikas, err := s.GetByID(id, userRole, userPerusahaanID)
	if err != nil {
		return nil, nil, err
	}

	// 2. Generate PDF using ikas/internal/utils
	pdfBytes, err := utils.GenerateIkasPDF(ikas)
	if err != nil {
		return nil, nil, err
	}

	return ikas, pdfBytes, nil
}

func (s *IkasService) RequestEdit(ctx context.Context, id string, reason string, userID string, userRole string, userPerusahaanID string) error {
	// 1. Get current data
	ikas, err := s.GetByID(id, userRole, userPerusahaanID)
	if err != nil {
		return err
	}

	// 2. Validate state
	if !ikas.IsValidated {
		return fmt.Errorf("data IKAS belum divalidasi, tidak perlu meminta izin edit")
	}
	if ikas.EditRequestStatus == "pending" {
		return fmt.Errorf("permintaan edit sebelumnya masih menunggu persetujuan admin")
	}

	// 3. Update status in DB
	err = s.repo.UpdateRequestEditStatus(id, "pending", reason)
	if err != nil {
		return fmt.Errorf("gagal mengajukan permintaan edit: %v", err)
	}

	// 4. Publish Event
	if s.producer != nil {
		event := dto_event.IkasEditRequestedEvent{
			IkasID:         ikas.ID,
			NamaPerusahaan: ikas.Perusahaan.NamaPerusahaan,
			Responden:      ikas.Responden,
			Reason:         reason,
			RequestedAt:    time.Now(),
		}
		_ = s.producer.PublishIkasEditRequested(ctx, event)
	}

	// Audit Log
	auditEvent := dto_event.IkasAuditLogEvent{
		IkasID:    id,
		UserID:    userID,
		Action:    "REQUEST_EDIT_IKAS",
		Changes:   map[string]interface{}{"reason": reason},
		Timestamp: time.Now(),
	}
	if s.producer != nil {
		_ = s.producer.PublishIkasAuditLog(ctx, auditEvent)
	}

	return nil
}

func (s *IkasService) ApproveEdit(ctx context.Context, id string, userID string) error {
	// 1. Get current data (admin check is handled by handler/middleware)
	ikas, err := s.GetByID(id, "admin", "")
	if err != nil {
		return err
	}

	if ikas.EditRequestStatus != "pending" {
		return fmt.Errorf("tidak ada permintaan edit aktif untuk data ini")
	}

	// 2. Unlock data and reset status
	err = s.repo.UpdateValidationStatus(id, false)
	if err != nil {
		return err
	}

	err = s.repo.UpdateRequestEditStatus(id, "none", "")
	if err != nil {
		return err
	}

	// 3. Publish Event
	if s.producer != nil {
		event := dto_event.IkasEditActionedEvent{
			IkasID:         ikas.ID,
			NamaPerusahaan: ikas.Perusahaan.NamaPerusahaan,
			Status:         "approved",
			ActionedAt:     time.Now(),
		}
		_ = s.producer.PublishIkasEditActioned(ctx, event)
	}

	// Audit Log
	auditEvent := dto_event.IkasAuditLogEvent{
		IkasID:    id,
		UserID:    userID,
		Action:    "APPROVE_EDIT_IKAS",
		Changes:   map[string]interface{}{"status": "approved"},
		Timestamp: time.Now(),
	}
	if s.producer != nil {
		_ = s.producer.PublishIkasAuditLog(ctx, auditEvent)
	}

	// Invalidate Cache
	if s.cache != nil {
		s.cache.Delete(cache.CacheKeyIkasRecords)
		if ikas.Perusahaan != nil {
			s.cache.Delete(fmt.Sprintf("%s:perusahaan:%s", cache.CacheKeyIkasRecords, ikas.Perusahaan.ID))
		}
	}

	return nil
}

func (s *IkasService) RejectEdit(ctx context.Context, id string, adminReason string, userID string) error {
	// 1. Get current data
	ikas, err := s.GetByID(id, "admin", "")
	if err != nil {
		return err
	}

	if ikas.EditRequestStatus != "pending" {
		return fmt.Errorf("tidak ada permintaan edit aktif untuk data ini")
	}

	// 2. Update status to rejected
	err = s.repo.UpdateRequestEditStatus(id, "rejected", adminReason)
	if err != nil {
		return err
	}

	// 3. Publish Event
	if s.producer != nil {
		event := dto_event.IkasEditActionedEvent{
			IkasID:         ikas.ID,
			NamaPerusahaan: ikas.Perusahaan.NamaPerusahaan,
			Status:         "rejected",
			AdminReason:    adminReason,
			ActionedAt:     time.Now(),
		}
		_ = s.producer.PublishIkasEditActioned(ctx, event)
	}

	// Audit Log
	auditEvent := dto_event.IkasAuditLogEvent{
		IkasID:    id,
		UserID:    userID,
		Action:    "REJECT_EDIT_IKAS",
		Changes:   map[string]interface{}{"status": "rejected", "reason": adminReason},
		Timestamp: time.Now(),
	}
	if s.producer != nil {
		_ = s.producer.PublishIkasAuditLog(ctx, auditEvent)
	}

	return nil
}
