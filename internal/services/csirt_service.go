package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"

	"github.com/google/uuid"
)

type CsirtService struct {
	repo *repository.CsirtRepository
}

func NewCsirtService(repo *repository.CsirtRepository) *CsirtService {
	return &CsirtService{repo: repo}
}

func (s *CsirtService) Create(req dto.CreateCsirtRequest) (*models.Csirt, error) {
	id := uuid.New().String()
	if err := s.repo.Create(req, id); err != nil {
		return nil, err
	}
	return s.repo.GetByID(id)
}

func (s *CsirtService) GetAll() ([]dto.CsirtResponse, error) {
	return s.repo.GetAllWithPerusahaan()
}

func (s *CsirtService) GetByID(id string) (*dto.CsirtResponse, error) {
	return s.repo.GetByIDWithPerusahaan(id)
}

func (s *CsirtService) Update(id string, req dto.UpdateCsirtRequest) (*models.Csirt, error) {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.NamaCsirt != nil {
		c.NamaCsirt = *req.NamaCsirt
	}
	if req.WebCsirt != nil {
		c.WebCsirt = *req.WebCsirt
	}
	if req.TeleponCsirt != nil {
		c.TeleponCsirt = req.TeleponCsirt
	}
	if req.PhotoCsirt != nil {
		c.PhotoCsirt = *req.PhotoCsirt
	}
	if req.FileRFC2350 != nil {
		c.FileRFC2350 = *req.FileRFC2350
	}
	if req.FilePublicKeyPGP != nil {
		c.FilePublicKeyPGP = *req.FilePublicKeyPGP
	}

	if err := s.repo.Update(id, *c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CsirtService) Delete(id string) error {
	return s.repo.Delete(id)
}
