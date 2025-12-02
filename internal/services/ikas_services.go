package services

import (
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
)

type IkasService interface {
	CreateIkas(ikas *models.Ikas) error
	GetAllIkas() ([]models.Ikas, error)
	GetIkasByID(id int) (*models.Ikas, error)
	UpdateIkas(id int, ikas *models.Ikas) error
	DeleteIkas(id int) error
}

type ikasService struct {
	repo repository.IkasRepository
}

func NewIkasService(repo repository.IkasRepository) IkasService {
	return &ikasService{repo: repo}
}

func (s *ikasService) CreateIkas(ikas *models.Ikas) error {
	return s.repo.Create(ikas)
}

func (s *ikasService) GetAllIkas() ([]models.Ikas, error) {
	return s.repo.GetAll()
}

func (s *ikasService) GetIkasByID(id int) (*models.Ikas, error) {
	return s.repo.GetByID(id)
}

func (s *ikasService) UpdateIkas(id int, ikas *models.Ikas) error {
	return s.repo.Update(id, ikas)
}

func (s *ikasService) DeleteIkas(id int) error {
	return s.repo.Delete(id)
}