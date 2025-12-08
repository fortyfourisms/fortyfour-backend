package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"

	"github.com/google/uuid"
)

type ProteksiService struct {
	repo *repository.ProteksiRepository
}

func NewProteksiService(repo *repository.ProteksiRepository) *ProteksiService {
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

	sum := 0
	totalNilai := 0.0

	if proteksi.NilaiSubdomain1 > 0 {
		sum++
		totalNilai += proteksi.NilaiSubdomain1
	}
	if proteksi.NilaiSubdomain2 > 0 {
		sum++
		totalNilai += proteksi.NilaiSubdomain2
	}
	if proteksi.NilaiSubdomain3 > 0 {
		sum++
		totalNilai += proteksi.NilaiSubdomain3
	}
	if proteksi.NilaiSubdomain4 > 0 {
		sum++
		totalNilai += proteksi.NilaiSubdomain4
	}
	if proteksi.NilaiSubdomain5 > 0 {
		sum++
		totalNilai += proteksi.NilaiSubdomain5
	}
	if proteksi.NilaiSubdomain6 > 0 {
		sum++
		totalNilai += proteksi.NilaiSubdomain6
	}

	if sum > 0 {
		proteksi.NilaiProteksi = totalNilai / float64(sum)
	} else {
		proteksi.NilaiProteksi = 0
	}

	if err := s.repo.Update(id, *proteksi); err != nil {
		return nil, err
	}

	return proteksi, nil
}

func (s *ProteksiService) Delete(id string) error {
	return s.repo.Delete(id)
}