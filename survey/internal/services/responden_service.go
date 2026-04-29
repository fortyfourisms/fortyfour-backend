package services

import (
	"errors"
	"strings"
	"time"

	"survey/internal/dto"
	"survey/internal/models"
)

type RespondenRepositoryInterface interface {
	Create(m models.Responden) error
	GetDetailByUserID(userID string) (*models.RespondenDetail, error)
	GetAllDetail() ([]models.RespondenDetail, error)
	GetDetailByID(id int) (*models.RespondenDetail, error)
	GetByID(id int) (*models.Responden, error)
	Update(id int, m models.Responden) error
}

type Validator interface {
	ValidateCreate(dto.CreateRespondenRequest) error
	ValidateUpdate(dto.UpdateRespondenRequest) error
}

type RespondenService struct {
	repo      RespondenRepositoryInterface
	validator Validator
}

// constructor
func NewRespondenService(repo RespondenRepositoryInterface, v Validator) *RespondenService {
	return &RespondenService{
		repo:      repo,
		validator: v,
	}
}

// CREATE
func (s *RespondenService) Create(req dto.CreateRespondenRequest) (*dto.RespondenResponse, error) {

	if err := s.validator.ValidateCreate(req); err != nil {
		return nil, err
	}

	model := models.Responden{
		UserID:             strings.TrimSpace(req.UserID),
		NoTelepon:          strings.TrimSpace(req.NoTelepon),
		Sektor:             strings.TrimSpace(req.Sektor),
		SektorLainnya:      toStringPtr(strings.TrimSpace(req.SektorLainnya)),
		SertifikatTraining: strings.TrimSpace(req.SertifikatTraining),
	}

	if err := s.repo.Create(model); err != nil {
		return nil, err
	}

	data, err := s.repo.GetDetailByUserID(model.UserID)
	if err != nil {
		return nil, err
	}

	resp := s.toResponse(data)
	return &resp, nil
}

// GET ALL
func (s *RespondenService) GetAll() ([]dto.RespondenResponse, error) {

	data, err := s.repo.GetAllDetail()
	if err != nil {
		return nil, err
	}

	var result []dto.RespondenResponse
	for i := range data {
		result = append(result, s.toResponse(&data[i]))
	}

	return result, nil
}

// GET BY ID
func (s *RespondenService) GetByID(id int) (*dto.RespondenResponse, error) {

	if id <= 0 {
		return nil, errors.New("id tidak valid")
	}

	data, err := s.repo.GetDetailByID(id)
	if err != nil {
		return nil, err
	}

	resp := s.toResponse(data)
	return &resp, nil
}

// UPDATE
func (s *RespondenService) Update(id int, req dto.UpdateRespondenRequest) (*dto.RespondenResponse, error) {

	if id <= 0 {
		return nil, errors.New("id tidak valid")
	}

	_, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.validator.ValidateUpdate(req); err != nil {
		return nil, err
	}

	model := models.Responden{
		NoTelepon:          strings.TrimSpace(req.NoTelepon),
		Sektor:             strings.TrimSpace(req.Sektor),
		SektorLainnya:      toStringPtr(strings.TrimSpace(req.SektorLainnya)),
		SertifikatTraining: strings.TrimSpace(req.SertifikatTraining),
	}

	if err := s.repo.Update(id, model); err != nil {
		return nil, err
	}

	updated, err := s.repo.GetDetailByID(id)
	if err != nil {
		return nil, err
	}

	resp := s.toResponse(updated)
	return &resp, nil
}

// HELPER
func toStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// RESPONSE MAPPING
func (s *RespondenService) toResponse(m *models.RespondenDetail) dto.RespondenResponse {

	return dto.RespondenResponse{
		ID:     m.ID,
		UserID: m.UserID,

		NamaLengkap:    safeString(m.NamaLengkap),
		Email:          safeString(m.Email),
		Jabatan:        safeString(m.Jabatan),
		NamaPerusahaan: safeString(m.NamaPerusahaan),
		PerusahaanID:   safeString(m.PerusahaanID),

		NoTelepon:          m.NoTelepon,
		Sektor:             m.Sektor,
		SektorLainnya:      safeString(m.SektorLainnya),
		SertifikatTraining: m.SertifikatTraining,

		CreatedAt: m.CreatedAt.Format(time.RFC3339),
		UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
	}
}