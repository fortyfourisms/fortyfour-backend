package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"
)

type PerusahaanService struct {
	perusahaanRepo *repository.PerusahaanRepository
}

func NewPerusahaanService(perusahaanRepo *repository.PerusahaanRepository) *PerusahaanService {
	return &PerusahaanService{perusahaanRepo: perusahaanRepo}
}

func (s *PerusahaanService) Create(req dto.PerusahaanRequest) (*dto.PerusahaanResponse, error) {
	id, err := s.perusahaanRepo.Create(req)
	if err != nil {
		return nil, err
	}

	// Fetch data lengkap setelah insert
	p, err := s.perusahaanRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *PerusahaanService) GetAll() ([]dto.PerusahaanResponse, error) {
	return s.perusahaanRepo.GetAll()
}

func (s *PerusahaanService) GetByID(id int) (*dto.PerusahaanResponse, error) {
	return s.perusahaanRepo.GetByID(id)
}

func (s *PerusahaanService) Update(id int, req dto.PerusahaanRequest) error {
	return s.perusahaanRepo.Update(id, req)
}

func (s *PerusahaanService) Delete(id int) error {
	return s.perusahaanRepo.Delete(id)
}
