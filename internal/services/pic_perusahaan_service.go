package services

import (
	"errors"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"
	"strings"

	"github.com/google/uuid"
)

type PICService struct {
	repo *repository.PICRepository
}

func NewPICService(repo *repository.PICRepository) *PICService {
	return &PICService{repo: repo}
}

func (s *PICService) Create(req dto.CreatePICRequest) (*dto.PICResponse, error) {

	if req.Nama == nil || strings.TrimSpace(*req.Nama) == "" {
		return nil, errors.New("nama wajib diisi")
	}

	if req.IDPerusahaan == nil || strings.TrimSpace(*req.IDPerusahaan) == "" {
		return nil, errors.New("id_perusahaan wajib diisi")
	}

	id := uuid.New().String()

	if err := s.repo.Create(req, id); err != nil {
		return nil, err
	}

	return s.repo.GetByID(id)
}

func (s *PICService) GetAll() ([]dto.PICResponse, error) {
	return s.repo.GetAll()
}

func (s *PICService) GetByID(id string) (*dto.PICResponse, error) {
	return s.repo.GetByID(id)
}

func (s *PICService) Update(id string, req dto.UpdatePICRequest) (*dto.PICResponse, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Nama != nil {
		existing.Nama = *req.Nama
	}
	if req.Telepon != nil {
		existing.Telepon = *req.Telepon
	}
	if req.IDPerusahaan != nil {
		existing.IDPerusahaan = *req.IDPerusahaan
	}

	if err := s.repo.Update(id, *existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *PICService) Delete(id string) error {
	return s.repo.Delete(id)
}
