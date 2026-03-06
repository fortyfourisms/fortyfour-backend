package services

import (
	"context"
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
	// Nested creation routines have been removed per recent refactor.
	// Nilai kematangan initialized to 0, it will be automatically calculated when domains are filled later.
	nilaiKematangan := 0.0

	// Create IKAS record
	err := s.repo.Create(req, id, nilaiKematangan)
	if err != nil {
		return err
	}

	go s.publishIkasCreatedEvent(id, req, nilaiKematangan)

	return nil
}

func (s *IkasService) publishIkasCreatedEvent(ikasID string, req dto.CreateIkasRequest, nilaiKematangan float64) {
	if s.producer == nil {
		return // Skip jika producer tidak ada
	}

	event := dto_event.IkasCreatedEvent{
		IkasID:          ikasID,
		IDPerusahaan:    req.IDPerusahaan,
		Tanggal:         req.Tanggal,
		Responden:       req.Responden,
		NilaiKematangan: nilaiKematangan,
		TargetNilai:     req.TargetNilai,
		CreatedAt:       time.Now(),
	}

	ctx := context.Background()
	if err := s.producer.PublishIkasCreated(ctx, event); err != nil {
		// Log error tapi jangan fail create operation
		// logger.Error(err, "operation failed")
	}
}

func (s *IkasService) GetAll() ([]dto.IkasResponse, error) {
	return s.repo.GetAll()
}

func (s *IkasService) GetByID(id string) (*dto.IkasResponse, error) {
	return s.repo.GetByID(id)
}

func (s *IkasService) Update(id string, req dto.UpdateIkasRequest) (*dto.IkasResponse, error) {
	// Ambil data existing
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Simpan old nilai untuk event
	oldNilaiKematangan := existing.NilaiKematangan

	// Update data IKAS utama (Nested domain checks removed)
	err = s.repo.Update(id, req)
	if err != nil {
		return nil, err
	}

	// Fetch updated data
	updated, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Publish event dengan nilai terbaru (if existing, use 0 if there are changes but kematangan calculation continues to be dynamic in get mode)
	go s.publishIkasUpdatedEvent(id, updated.NilaiKematangan, oldNilaiKematangan)

	return updated, nil
}

func (s *IkasService) publishIkasUpdatedEvent(ikasID string, newNilai, oldNilai float64) {
	if s.producer == nil {
		return
	}

	event := dto_event.IkasUpdatedEvent{
		IkasID:             ikasID,
		OldNilaiKematangan: oldNilai,
		NewNilaiKematangan: newNilai,
		UpdatedAt:          time.Now(),
	}

	ctx := context.Background()
	if err := s.producer.PublishIkasUpdated(ctx, event); err != nil {
		// Log error tapi jangan fail update operation
		// logger.Error(err, "operation failed")
	}
}

func (s *IkasService) Delete(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	go s.publishIkasDeletedEvent(id)

	return nil
}

func (s *IkasService) publishIkasDeletedEvent(ikasID string) {
	if s.producer == nil {
		return
	}

	event := dto_event.IkasDeletedEvent{
		IkasID:    ikasID,
		DeletedAt: time.Now(),
	}

	ctx := context.Background()
	if err := s.producer.PublishIkasDeleted(ctx, event); err != nil {
		// Log error tapi jangan fail delete operation
		// logger.Error(err, "operation failed")
	}
}

func (s *IkasService) ImportFromExcel(fileData []byte) (*dto.IkasResponse, error) {
	// Parse Excel
	excelData, err := s.repo.ParseExcelForImport(fileData)
	if err != nil {
		return nil, err
	}

	// Generate ID baru
	newID := uuid.New().String()

	// Create menggunakan service Create yang sudah ada
	if err := s.Create(*excelData, newID); err != nil {
		return nil, err
	}

	// Ambil data yang baru dibuat
	resp, err := s.GetByID(newID)
	if err != nil {
		return nil, err
	}

	go s.publishIkasImportedEvent(newID, excelData, resp.NilaiKematangan)

	return resp, nil
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
		// Log error tapi jangan fail import operation
		// logger.Error(err, "operation failed")
	}
}
