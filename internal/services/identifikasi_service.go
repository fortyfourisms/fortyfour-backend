package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"

	"github.com/google/uuid"
)

type IdentifikasiService struct {
	repo *repository.IdentifikasiRepository
}

func NewIdentifikasiService(repo *repository.IdentifikasiRepository) *IdentifikasiService {
	return &IdentifikasiService{repo: repo}
}

func (s *IdentifikasiService) Create(req dto.CreateIdentifikasiRequest) (*models.Identifikasi, error) {
	id := uuid.New().String()
	if err := s.repo.Create(req, id); err != nil {
		return nil, err
	}

	return s.repo.GetByID(id)
}

func (s *IdentifikasiService) GetAll() ([]models.Identifikasi, error) {
	return s.repo.GetAll()
}

func (s *IdentifikasiService) GetByID(id string) (*models.Identifikasi, error) {
	return s.repo.GetByID(id)
}

func (s *IdentifikasiService) Update(id string, req dto.UpdateIdentifikasiRequest) (*models.Identifikasi, error) {
	identifikasi, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.NilaiIdentifikasi != nil {
		identifikasi.NilaiIdentifikasi = *req.NilaiIdentifikasi
	}
	if req.NilaiSubdomain1 != nil {
		identifikasi.NilaiSubdomain1 = *req.NilaiSubdomain1
	}
	if req.NilaiSubdomain2 != nil {
		identifikasi.NilaiSubdomain2 = *req.NilaiSubdomain2
	}
	if req.NilaiSubdomain3 != nil {
		identifikasi.NilaiSubdomain3 = *req.NilaiSubdomain3
	}
	if req.NilaiSubdomain4 != nil {
		identifikasi.NilaiSubdomain4 = *req.NilaiSubdomain4
	}
	if req.NilaiSubdomain5 != nil {
		identifikasi.NilaiSubdomain5 = *req.NilaiSubdomain5
	}

	if err := s.repo.Update(id, *identifikasi); err != nil {
		return nil, err
	}

	return identifikasi, nil
}

func (s *IdentifikasiService) Delete(id string) error {
	return s.repo.Delete(id)
}
