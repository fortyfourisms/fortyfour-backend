package services

import (
	"errors"
	"ikas/internal/models"
	"ikas/internal/repository"
)

type IdentifikasiService struct {
	repo repository.IdentifikasiRepositoryInterface
}

func NewIdentifikasiService(repo repository.IdentifikasiRepositoryInterface) *IdentifikasiService {
	return &IdentifikasiService{repo: repo}
}

func (s *IdentifikasiService) GetAll() ([]models.Identifikasi, error) {
	return s.repo.GetAll()
}

func (s *IdentifikasiService) GetByPerusahaan(perusahaanID string) ([]models.Identifikasi, error) {
	return s.repo.GetByPerusahaan(perusahaanID)
}

func (s *IdentifikasiService) GetByID(id string, userRole string, userPerusahaanID string) (*models.Identifikasi, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if userRole != "admin" && data.PerusahaanID != userPerusahaanID {
		return nil, errors.New("anda tidak memiliki akses ke data ini")
	}
	return data, nil
}
