package services

import (
	"errors"
	"ikas/internal/models"
	"ikas/internal/repository"
)

type GulihService struct {
	repo repository.GulihRepositoryInterface
}

func NewGulihService(repo repository.GulihRepositoryInterface) *GulihService {
	return &GulihService{repo: repo}
}

func (s *GulihService) GetAll() ([]models.Gulih, error) {
	return s.repo.GetAll()
}

func (s *GulihService) GetByPerusahaan(perusahaanID string) ([]models.Gulih, error) {
	return s.repo.GetByPerusahaan(perusahaanID)
}

func (s *GulihService) GetByID(id string, userRole string, userPerusahaanID string) (*models.Gulih, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if userRole != "admin" && data.PerusahaanID != userPerusahaanID {
		return nil, errors.New("anda tidak memiliki akses ke data ini")
	}
	return data, nil
}
