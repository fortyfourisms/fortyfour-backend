package services

import (
	"errors"
	"ikas/internal/models"
	"ikas/internal/repository"
)

type IdentifikasiService struct {
	repo     repository.IdentifikasiRepositoryInterface
	ikasRepo repository.IkasRepositoryInterface
}

func NewIdentifikasiService(
	repo repository.IdentifikasiRepositoryInterface,
	ikasRepo repository.IkasRepositoryInterface,
) *IdentifikasiService {
	return &IdentifikasiService{
		repo:     repo,
		ikasRepo: ikasRepo,
	}
}

func (s *IdentifikasiService) GetAll(userRole string) ([]models.Identifikasi, error) {
	if userRole != "admin" {
		return nil, errors.New("anda tidak memiliki akses untuk melihat semua data")
	}
	return s.repo.GetAll()
}

func (s *IdentifikasiService) GetByIkasID(ikasID string, userRole string, userPerusahaanID string) ([]models.Identifikasi, error) {
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

func (s *IdentifikasiService) GetByID(id string, userRole string, userPerusahaanID string) (*models.Identifikasi, error) {
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

func (s *IdentifikasiService) GetByPerusahaanID(perusahaanID string, userRole string, userPerusahaanID string) ([]models.Identifikasi, error) {
	if userRole != "admin" {
		if perusahaanID != userPerusahaanID {
			return nil, errors.New("anda tidak memiliki akses ke data perusahaan ini")
		}
	}
	return s.repo.GetByPerusahaanID(perusahaanID)
}
