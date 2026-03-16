package services

import (
	"context"
	"fmt"
	"ikas/internal/dto"
	"ikas/internal/dto/dto_event"
	"ikas/internal/rabbitmq"
	"ikas/internal/repository"
	"time"

	"github.com/google/uuid"
)

type IkasService struct {
	repo     repository.IkasRepositoryInterface
	producer *rabbitmq.Producer
}

func NewIkasService(repo repository.IkasRepositoryInterface, producer *rabbitmq.Producer) *IkasService {
	return &IkasService{
		repo:     repo,
		producer: producer,
	}
}

func (s *IkasService) Create(req dto.CreateIkasRequest, id string) error {
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
		CreatedAt:       time.Now(),
	}

	if s.producer == nil {
		return nil
	}

	ctx := context.Background()
	if err := s.producer.PublishIkasCreated(ctx, event); err != nil {
		return err
	}

	return nil
}

func (s *IkasService) GetAll() ([]dto.IkasResponse, error) {
	return s.repo.GetAll()
}

func (s *IkasService) GetByID(id string) (*dto.IkasResponse, error) {
	return s.repo.GetByID(id)
}

func (s *IkasService) Update(id string, req dto.UpdateIkasRequest) error {
	// Check existence
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
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
		UpdatedAt:    time.Now(),
	}

	if s.producer == nil {
		return nil
	}

	ctx := context.Background()
	if err := s.producer.PublishIkasUpdated(ctx, event); err != nil {
		return err
	}

	return nil
}

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func getFloatValue(f *float64) float64 {
	if f == nil {
		return 0.0
	}
	return *f
}

func (s *IkasService) Delete(id string) error {
	// Check existence
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Publish delete event
	event := dto_event.IkasDeletedEvent{
		IkasID:    id,
		DeletedAt: time.Now(),
	}

	if s.producer == nil {
		return nil
	}

	ctx := context.Background()
	if err := s.producer.PublishIkasDeleted(ctx, event); err != nil {
		return err
	}

	return nil
}

func (s *IkasService) ImportFromExcel(fileData []byte) (string, error) {
	excelData, err := s.repo.ParseExcelForImport(fileData)
	if err != nil {
		return "", err
	}

	newID := uuid.New().String()

	// 1. Create main IKAS record
	if err := s.Create(excelData.IkasRequest, newID); err != nil {
		return "", err
	}

	// 2. Publish events for each subdomain to trigger automatic processing
	perusahaanID := excelData.IkasRequest.IDPerusahaan

	// Identifikasi
	for _, ans := range excelData.JawabanIdentifikasi {
		event := dto.CreateJawabanIdentifikasiRequest{
			PertanyaanIdentifikasiID: ans.PertanyaanID,
			PerusahaanID:           perusahaanID,
			JawabanIdentifikasi:    &ans.Jawaban,
		}
		s.producer.PublishJawabanIdentifikasiCreated(context.Background(), event)
	}

	// Proteksi
	for _, ans := range excelData.JawabanProteksi {
		event := dto.CreateJawabanProteksiRequest{
			PertanyaanProteksiID: ans.PertanyaanID,
			PerusahaanID:       perusahaanID,
			JawabanProteksi:    &ans.Jawaban,
		}
		s.producer.PublishJawabanProteksiCreated(context.Background(), event)
	}

	// Deteksi
	for _, ans := range excelData.JawabanDeteksi {
		event := dto.CreateJawabanDeteksiRequest{
			PertanyaanDeteksiID: ans.PertanyaanID,
			PerusahaanID:      perusahaanID,
			JawabanDeteksi:    &ans.Jawaban,
		}
		s.producer.PublishJawabanDeteksiCreated(context.Background(), event)
	}

	// Gulih
	for _, ans := range excelData.JawabanGulih {
		event := dto.CreateJawabanGulihRequest{
			PertanyaanGulihID: ans.PertanyaanID,
			PerusahaanID:     perusahaanID,
			JawabanGulih:    &ans.Jawaban,
		}
		s.producer.PublishJawabanGulihCreated(context.Background(), event)
	}

	return newID, nil
}

func (s *IkasService) publishIkasImportedEvent(ikasID string, req *dto.CreateIkasRequest, nilaiKematangan float64) {
	if s.producer == nil {
		return
	}

	event := dto_event.IkasImportedEvent{
		IkasID:          ikasID,
		IDPerusahaan:    req.IDPerusahaan,
		NilaiKematangan: nilaiKematangan,
		ImportedAt:      time.Now(),
	}

	ctx := context.Background()
	if err := s.producer.PublishIkasImported(ctx, event); err != nil {
		// Log error but don't fail import
	}
}
