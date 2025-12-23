package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"
)

type IkasService struct {
	repo *repository.IkasRepository
}

func NewIkasService(repo *repository.IkasRepository) *IkasService {
	return &IkasService{repo: repo}
}

func (s *IkasService) Create(req dto.CreateIkasRequest, id string) error {
	return s.repo.Create(req, id)
}

func (s *IkasService) GetAll() ([]dto.IkasResponse, error) {
	return s.repo.GetAll()
}

func (s *IkasService) GetByID(id string) (*dto.IkasResponse, error) {
	return s.repo.GetByID(id)
}

func (s *IkasService) Update(id string, req dto.UpdateIkasRequest) (*dto.IkasResponse, error) {
	// Update data
	if err := s.repo.Update(id, req); err != nil {
		return nil, err
	}

	// Ambil data terbaru dengan JOIN
	updated, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *IkasService) Delete(id string) error {
	return s.repo.Delete(id)
}
