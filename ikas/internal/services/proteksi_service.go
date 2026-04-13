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

func (s *ProteksiService) GetByIkasID(ikasID string) ([]models.Proteksi, error) {
	return s.repo.GetByIkasID(ikasID)
}

func (s *ProteksiService) GetByID(id string, userRole string, userPerusahaanID string) (*models.Proteksi, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if userRole != "admin" {
		// temporary workaround
	}
	return data, nil
}
