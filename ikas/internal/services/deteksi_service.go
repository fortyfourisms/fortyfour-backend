package services

import (
	"ikas/internal/models"
	"ikas/internal/repository"
)

type DeteksiService struct {
	repo repository.DeteksiRepositoryInterface
}

func NewDeteksiService(repo repository.DeteksiRepositoryInterface) *DeteksiService {
	return &DeteksiService{repo: repo}
}

func (s *DeteksiService) GetAll() ([]models.Deteksi, error) {
	return s.repo.GetAll()
}

func (s *DeteksiService) GetByID(id string) (*models.Deteksi, error) {
	return s.repo.GetByID(id)
}
