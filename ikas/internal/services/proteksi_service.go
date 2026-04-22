package services

import (
	"errors"
	"ikas/internal/models"
	"ikas/internal/repository"
)

type ProteksiService struct {
	repo     repository.ProteksiRepositoryInterface
	ikasRepo repository.IkasRepositoryInterface
}

func NewProteksiService(
	repo repository.ProteksiRepositoryInterface,
	ikasRepo repository.IkasRepositoryInterface,
) *ProteksiService {
	return &ProteksiService{
		repo:     repo,
		ikasRepo: ikasRepo,
	}
}

func (s *ProteksiService) GetAll(userRole string) ([]models.Proteksi, error) {
	if userRole != "admin" {
		return nil, errors.New("anda tidak memiliki akses untuk melihat semua data")
	}
	return s.repo.GetAll()
}

func (s *ProteksiService) GetByIkasID(ikasID string, userRole string, userPerusahaanID string) ([]models.Proteksi, error) {
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

func (s *ProteksiService) GetByID(id string, userRole string, userPerusahaanID string) (*models.Proteksi, error) {
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

func (s *ProteksiService) GetByPerusahaanID(perusahaanID string, userRole string, userPerusahaanID string) ([]models.Proteksi, error) {
	if userRole != "admin" {
		if perusahaanID != userPerusahaanID {
			return nil, errors.New("anda tidak memiliki akses ke data perusahaan ini")
		}
	}
	return s.repo.GetByPerusahaanID(perusahaanID)
}
