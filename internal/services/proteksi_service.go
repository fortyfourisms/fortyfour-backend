package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"

	"github.com/google/uuid"
)

type ProteksiService struct {
	repo repository.ProteksiRepositoryInterface
}

func NewProteksiService(repo repository.ProteksiRepositoryInterface) *ProteksiService {
	return &ProteksiService{repo: repo}
}

func (s *ProteksiService) Create(req dto.CreateProteksiRequest) (*models.Proteksi, error) {
	id := uuid.New().String()

	if err := s.repo.Create(req, id); err != nil {
		return nil, err
	}

	return s.repo.GetByID(id)
}

func (s *ProteksiService) GetAll() ([]models.Proteksi, error) {
	return s.repo.GetAll()
}

func (s *ProteksiService) GetByID(id string) (*models.Proteksi, error) {
	return s.repo.GetByID(id)
}

func (s *ProteksiService) Update(id string, req dto.UpdateProteksiRequest) (*models.Proteksi, error) {
	proteksi, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if req.NilaiProteksi != nil {
		proteksi.NilaiProteksi = *req.NilaiProteksi
	}
	if req.NilaiSubdomain1 != nil {
		proteksi.NilaiSubdomain1 = *req.NilaiSubdomain1
	}
	if req.NilaiSubdomain2 != nil {
		proteksi.NilaiSubdomain2 = *req.NilaiSubdomain2
	}
	if req.NilaiSubdomain3 != nil {
		proteksi.NilaiSubdomain3 = *req.NilaiSubdomain3
	}
	if req.NilaiSubdomain4 != nil {
		proteksi.NilaiSubdomain4 = *req.NilaiSubdomain4
	}
	if req.NilaiSubdomain5 != nil {
		proteksi.NilaiSubdomain5 = *req.NilaiSubdomain5
	}
	if req.NilaiSubdomain6 != nil {
		proteksi.NilaiSubdomain6 = *req.NilaiSubdomain6
	}

	if err := s.repo.Update(id, *proteksi); err != nil {
		return nil, err
	}

	return proteksi, nil
}

func (s *ProteksiService) Delete(id string) error {
	return s.repo.Delete(id)
}
