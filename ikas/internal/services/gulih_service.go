package services

import (
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

func (s *GulihService) GetByID(id string) (*models.Gulih, error) {
	return s.repo.GetByID(id)
}
