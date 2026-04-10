package services

import (
	"context"
	"fmt"
	"ikas/internal/dto"
	"ikas/internal/dto/dto_event"
	"ikas/internal/repository"
	"time"

	"github.com/google/uuid"
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
}

type IkasService struct {
	repo     repository.IkasRepositoryInterface
	producer IkasProducerInterface
}

func NewIkasService(repo repository.IkasRepositoryInterface, producer IkasProducerInterface) *IkasService {
	return &IkasService{
		repo:     repo,
		producer: producer,
	}
}

func (s *IkasService) Create(ctx context.Context, req dto.CreateIkasRequest, id string, userID string) error {
	// Check if IKAS for this company already exists
	exists, err := s.repo.CheckExistsByPerusahaanID(req.IDPerusahaan)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("Data IKAS untuk perusahaan ini sudah ada")
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

	return nil
}

func (s *IkasService) GetAll() ([]dto.IkasResponse, error) {
	return s.repo.GetAll()
}

func (s *IkasService) GetByPerusahaan(perusahaanID string) ([]dto.IkasResponse, error) {
	return s.repo.GetByPerusahaan(perusahaanID)
}

func (s *IkasService) GetByID(id string, userRole string, userPerusahaanID string) (*dto.IkasResponse, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if userRole != "admin" && data.Perusahaan != nil && data.Perusahaan.ID != userPerusahaanID {
		return nil, fmt.Errorf("anda tidak memiliki akses ke data ini")
	}
	return data, nil
}

func (s *IkasService) Update(ctx context.Context, id string, req dto.UpdateIkasRequest, userID string, userRole string, userPerusahaanID string) error {
	// Check existence and get current state
	current, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if userRole != "admin" && current.Perusahaan != nil && current.Perusahaan.ID != userPerusahaanID {
		return fmt.Errorf("anda tidak memiliki akses untuk mengubah data ini")
	}

	// Change detection for audit log
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
		return nil
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
		return err
	}

	return nil
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

	return nil
}

func (s *IkasService) ImportFromExcel(ctx context.Context, fileData []byte, userID string) (string, error) {
	excelData, err := s.repo.ParseExcelForImport(fileData)
	if err != nil {
		return "", err
	}

	newID := uuid.New().String()

	// 1. Create main IKAS record
	if err := s.Create(ctx, excelData.IkasRequest, newID, userID); err != nil {
		return "", err
	}

	// 2. Publish events for each subdomain to trigger automatic processing
	perusahaanID := excelData.IkasRequest.IDPerusahaan

	// Identifikasi
	for _, ans := range excelData.JawabanIdentifikasi {
		event := dto.CreateJawabanIdentifikasiRequest{
			PertanyaanIdentifikasiID: ans.PertanyaanID,
			PerusahaanID:             perusahaanID,
			JawabanIdentifikasi:      &ans.Jawaban,
		}
		s.producer.PublishJawabanIdentifikasiCreated(context.Background(), event)
	}

	// Proteksi
	for _, ans := range excelData.JawabanProteksi {
		event := dto.CreateJawabanProteksiRequest{
			PertanyaanProteksiID: ans.PertanyaanID,
			PerusahaanID:         perusahaanID,
			JawabanProteksi:      &ans.Jawaban,
		}
		s.producer.PublishJawabanProteksiCreated(context.Background(), event)
	}

	// Deteksi
	for _, ans := range excelData.JawabanDeteksi {
		event := dto.CreateJawabanDeteksiRequest{
			PertanyaanDeteksiID: ans.PertanyaanID,
			PerusahaanID:        perusahaanID,
			JawabanDeteksi:      &ans.Jawaban,
		}
		s.producer.PublishJawabanDeteksiCreated(context.Background(), event)
	}

	// Gulih
	for _, ans := range excelData.JawabanGulih {
		event := dto.CreateJawabanGulihRequest{
			PertanyaanGulihID: ans.PertanyaanID,
			PerusahaanID:      perusahaanID,
			JawabanGulih:      &ans.Jawaban,
		}
		s.producer.PublishJawabanGulihCreated(context.Background(), event)
	}

	return newID, nil
}
