package services

import (
	"errors"
	"ikas/internal/models"
	"ikas/internal/repository"
)

type DeteksiService struct {
	repo     repository.DeteksiRepositoryInterface
	ikasRepo repository.IkasRepositoryInterface
}

func NewDeteksiService(
	repo repository.DeteksiRepositoryInterface,
	ikasRepo repository.IkasRepositoryInterface,
) *DeteksiService {
	return &DeteksiService{
		repo:     repo,
		ikasRepo: ikasRepo,
	}
}

func (s *DeteksiService) GetAll(userRole string) ([]models.Deteksi, error) {
	if userRole != "admin" {
		return nil, errors.New("anda tidak memiliki akses untuk melihat semua data")
	}
	return s.repo.GetAll()
}

func (s *DeteksiService) GetByIkasID(ikasID string, userRole string, userPerusahaanID string) ([]models.Deteksi, error) {
	if userRole != "admin" {
		owned, err := s.ikasRepo.CheckOwnership(ikasID, userPerusahaanID)
		if err != nil {
			return nil, err
		}
		if !owned {
			return nil, errors.New("anda tidak memiliki akses ke data asesmen ini")
		}
	}
	return s.repo.GetByIkasID(ikasID)
}

func (s *DeteksiService) GetByID(id string, userRole string, userPerusahaanID string) (*models.Deteksi, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if userRole != "admin" {
		owned, err := s.ikasRepo.CheckOwnership(data.IkasID, userPerusahaanID)
		if err != nil {
			return nil, err
		}
		if !owned {
			return nil, errors.New("anda tidak memiliki akses ke data ini")
		}
	}

	return data, nil
}

func (s *DeteksiService) GetByPerusahaanID(perusahaanID string, userRole string, userPerusahaanID string) ([]models.Deteksi, error) {
	if userRole != "admin" {
		if perusahaanID != userPerusahaanID {
			return nil, errors.New("anda tidak memiliki akses ke data perusahaan ini")
		}
	}
	return s.repo.GetByPerusahaanID(perusahaanID)
}
