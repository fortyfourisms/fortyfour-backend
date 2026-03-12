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
		IDPerusahaan: getStringValue(req.IDPerusahaan),
		Tanggal:      getStringValue(req.Tanggal),
		Responden:    getStringValue(req.Responden),
		Telepon:      getStringValue(req.Telepon),
		Jabatan:      getStringValue(req.Jabatan),
		TargetNilai:  getFloatValue(req.TargetNilai),
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

	if err := s.Create(*excelData, newID); err != nil {
		return "", err
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
