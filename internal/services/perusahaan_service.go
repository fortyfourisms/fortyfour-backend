package services

import (
	"errors"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"
	"strings"

	"github.com/google/uuid"
)

type PerusahaanService struct {
	repo *repository.PerusahaanRepository
}

func NewPerusahaanService(repo *repository.PerusahaanRepository) *PerusahaanService {
	return &PerusahaanService{repo: repo}
}

func (s *PerusahaanService) Create(req dto.CreatePerusahaanRequest) (*dto.PerusahaanResponse, error) {

	if req.NamaPerusahaan == nil || strings.TrimSpace(*req.NamaPerusahaan) == "" {
		return nil, errors.New("nama_perusahaan wajib diisi")
	}

	if req.Sektor == nil || strings.TrimSpace(*req.Sektor) == "" {
		return nil, errors.New("sektor wajib diisi")
	}

	id := uuid.New().String()
	if err := s.repo.Create(req, id); err != nil {
		return nil, err
	}
	return s.repo.GetByID(id)
}

func (s *PerusahaanService) GetAll() ([]dto.PerusahaanResponse, error) {
	return s.repo.GetAll()
}

func (s *PerusahaanService) GetByID(id string) (*dto.PerusahaanResponse, error) {
	return s.repo.GetByID(id)
}

func (s *PerusahaanService) Update(id string, req dto.UpdatePerusahaanRequest) (*dto.PerusahaanResponse, error) {
	perusahaan, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.NamaPerusahaan != nil {
		perusahaan.NamaPerusahaan = *req.NamaPerusahaan
	}
	if req.Sektor != nil {
		perusahaan.Sektor = *req.Sektor
	}
	if req.Alamat != nil {
		perusahaan.Alamat = *req.Alamat
	}
	if req.Telepon != nil {
		perusahaan.Telepon = *req.Telepon
	}
	if req.Email != nil {
		perusahaan.Email = *req.Email
	}
	if req.Website != nil {
		perusahaan.Website = *req.Website
	}
	if req.Photo != nil {
		perusahaan.Photo = *req.Photo
	}

	if err := s.repo.Update(id, *perusahaan); err != nil {
		return nil, err
	}

	return perusahaan, nil
}

func (s *PerusahaanService) Delete(id string) error {
	return s.repo.Delete(id)
}
