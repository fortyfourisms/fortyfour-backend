package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"

	"github.com/google/uuid"
)

type GulihService struct {
	repo repository.GulihRepositoryInterface
}

func NewGulihService(repo repository.GulihRepositoryInterface) *GulihService {
	return &GulihService{repo: repo}
}

func (s *GulihService) Create(req dto.CreateGulihRequest) (*models.Gulih, error) {
	id := uuid.New().String()
	if err := s.repo.Create(req, id); err != nil {
		return nil, err
	}
	return s.repo.GetByID(id)
}

func (s *GulihService) GetAll() ([]models.Gulih, error) {
	return s.repo.GetAll()
}

func (s *GulihService) GetByID(id string) (*models.Gulih, error) {
	return s.repo.GetByID(id)
}

func (s *GulihService) Update(id string, req dto.UpdateGulihRequest) (*models.Gulih, error) {
	gulih, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.NilaiGulih != nil {
		gulih.NilaiGulih = *req.NilaiGulih
	}
	if req.NilaiSubdomain1 != nil {
		gulih.NilaiSubdomain1 = *req.NilaiSubdomain1
	}
	if req.NilaiSubdomain2 != nil {
		gulih.NilaiSubdomain2 = *req.NilaiSubdomain2
	}
	if req.NilaiSubdomain3 != nil {
		gulih.NilaiSubdomain3 = *req.NilaiSubdomain3
	}
	if req.NilaiSubdomain4 != nil {
		gulih.NilaiSubdomain4 = *req.NilaiSubdomain4
	}

	if err := s.repo.Update(id, *gulih); err != nil {
		return nil, err
	}

	return gulih, nil
}

func (s *GulihService) Delete(id string) error {
	return s.repo.Delete(id)
}
