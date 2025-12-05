package services

import (
	"errors"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"

	"github.com/google/uuid"
)

type PerusahaanService struct {
	repo *repository.PerusahaanRepository
}

func NewPerusahaanService(repo *repository.PerusahaanRepository) *PerusahaanService {
	return &PerusahaanService{repo: repo}
}

func (s *PerusahaanService) Create(req dto.CreatePerusahaanRequest) (*dto.PerusahaanResponse, error) {
	if req.NamaPerusahaan == nil || req.JenisUsaha == nil {
		return nil, errors.New("nama_perusahaan dan jenis_usaha wajib diisi")
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
	if req.JenisUsaha != nil {
		perusahaan.JenisUsaha = *req.JenisUsaha
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
