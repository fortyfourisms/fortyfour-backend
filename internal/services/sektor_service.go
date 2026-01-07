package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"
)

type SektorServiceInterface interface {
	GetAll() ([]dto.SektorResponse, error)
	GetByID(id string) (*dto.SektorResponse, error)
}

type SektorService struct {
	repo repository.SektorRepositoryInterface
}

func NewSektorService(repo repository.SektorRepositoryInterface) *SektorService {
	return &SektorService{repo: repo}
}

func (s *SektorService) GetAll() ([]dto.SektorResponse, error) {
	return s.repo.GetAll()
}

func (s *SektorService) GetByID(id string) (*dto.SektorResponse, error) {
	return s.repo.GetByID(id)
}