package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"
)

type SubSektorServiceInterface interface {
	GetAll() ([]dto.SubSektorResponse, error)
	GetByID(id string) (*dto.SubSektorResponse, error)
	GetBySektorID(sektorID string) ([]dto.SubSektorResponse, error)
}

type SubSektorService struct {
	repo repository.SubSektorRepositoryInterface
}

func NewSubSektorService(repo repository.SubSektorRepositoryInterface) *SubSektorService {
	return &SubSektorService{repo: repo}
}

func (s *SubSektorService) GetAll() ([]dto.SubSektorResponse, error) {
	return s.repo.GetAll()
}

func (s *SubSektorService) GetByID(id string) (*dto.SubSektorResponse, error) {
	return s.repo.GetByID(id)
}

func (s *SubSektorService) GetBySektorID(sektorID string) ([]dto.SubSektorResponse, error) {
	return s.repo.GetBySektorID(sektorID)
}
