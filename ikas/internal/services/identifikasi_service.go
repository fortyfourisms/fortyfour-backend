package services

import (
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

func (s *IdentifikasiService) GetByID(id string) (*models.Identifikasi, error) {
	return s.repo.GetByID(id)
}
