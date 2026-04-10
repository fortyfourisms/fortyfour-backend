package services

import (
	"errors"
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

func (s *ProteksiService) GetByPerusahaan(perusahaanID string) ([]models.Proteksi, error) {
	return s.repo.GetByPerusahaan(perusahaanID)
}

func (s *ProteksiService) GetByID(id string, userRole string, userPerusahaanID string) (*models.Proteksi, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if userRole != "admin" && data.PerusahaanID != userPerusahaanID {
		return nil, errors.New("anda tidak memiliki akses ke data ini")
	}
	return data, nil
}
