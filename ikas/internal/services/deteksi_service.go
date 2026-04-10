package services

import (
	"errors"
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

func (s *DeteksiService) GetByPerusahaan(perusahaanID string) ([]models.Deteksi, error) {
	return s.repo.GetByPerusahaan(perusahaanID)
}

func (s *DeteksiService) GetByID(id string, userRole string, userPerusahaanID string) (*models.Deteksi, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if userRole != "admin" && data.PerusahaanID != userPerusahaanID {
		return nil, errors.New("anda tidak memiliki akses ke data ini")
	}
	return data, nil
}
