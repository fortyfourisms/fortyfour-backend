package services

import (
	"errors"
	"ikas/internal/models"
	"ikas/internal/repository"
)

type GulihService struct {
	repo     repository.GulihRepositoryInterface
	ikasRepo repository.IkasRepositoryInterface
}

func NewGulihService(
	repo repository.GulihRepositoryInterface,
	ikasRepo repository.IkasRepositoryInterface,
) *GulihService {
	return &GulihService{
		repo:     repo,
		ikasRepo: ikasRepo,
	}
}

func (s *GulihService) GetAll(userRole string) ([]models.Gulih, error) {
	if userRole != "admin" {
		return nil, errors.New("anda tidak memiliki akses untuk melihat semua data")
	}
	return s.repo.GetAll()
}

func (s *GulihService) GetByIkasID(ikasID string, userRole string, userPerusahaanID string) ([]models.Gulih, error) {
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

func (s *GulihService) GetByID(id string, userRole string, userPerusahaanID string) (*models.Gulih, error) {
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

func (s *GulihService) GetByPerusahaanID(perusahaanID string, userRole string, userPerusahaanID string) ([]models.Gulih, error) {
	if userRole != "admin" {
		if perusahaanID != userPerusahaanID {
			return nil, errors.New("anda tidak memiliki akses ke data perusahaan ini")
		}
	}
	return s.repo.GetByPerusahaanID(perusahaanID)
}
