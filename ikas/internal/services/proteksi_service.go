package services

import (
	"ikas/internal/models"
	"ikas/internal/repository"
)

type ProteksiService struct {
	repo repository.ProteksiRepositoryInterface
}

func NewProteksiService(repo repository.ProteksiRepositoryInterface) *ProteksiService {
	return &ProteksiService{repo: repo}
}

func (s *ProteksiService) GetAll() ([]models.Proteksi, error) {
	return s.repo.GetAll()
}

func (s *ProteksiService) GetByID(id string) (*models.Proteksi, error) {
	return s.repo.GetByID(id)
}
