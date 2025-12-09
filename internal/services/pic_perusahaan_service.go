package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"

	"github.com/google/uuid"
)

type PICPerusahaanService struct {
	Repo *repository.PICPerusahaanRepository
}

func NewPICPerusahaanService(repo *repository.PICPerusahaanRepository) *PICPerusahaanService {
	return &PICPerusahaanService{Repo: repo}
}

func (s *PICPerusahaanService) Create(req dto.CreatePICPerusahaanRequest) (models.PICPerusahaan, error) {
	pic := models.PICPerusahaan{
		ID:           uuid.New().String(),
		Nama:         req.Nama,
		Telepon:      req.Telepon,
		IDPerusahaan: req.IDPerusahaan,
	}

	err := s.Repo.Create(pic)
	return pic, err
}

func (s *PICPerusahaanService) GetAll() ([]models.PICPerusahaan, error) {
	return s.Repo.GetAll()
}

func (s *PICPerusahaanService) GetByID(id string) (models.PICPerusahaan, error) {
	return s.Repo.GetByID(id)
}

func (s *PICPerusahaanService) Update(id string, req dto.UpdatePICPerusahaanRequest) (models.PICPerusahaan, error) {
	pic := models.PICPerusahaan{
		Nama:         req.Nama,
		Telepon:      req.Telepon,
		IDPerusahaan: req.IDPerusahaan,
	}

	err := s.Repo.Update(id, pic)
	return pic, err
}

func (s *PICPerusahaanService) Delete(id string) error {
	return s.Repo.Delete(id)
}
