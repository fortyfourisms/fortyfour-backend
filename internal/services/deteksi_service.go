package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"

	"github.com/google/uuid"
	"github.com/rollbar/rollbar-go"
)

type DeteksiService struct {
	repo repository.DeteksiRepositoryInterface
}

func NewDeteksiService(repo repository.DeteksiRepositoryInterface) *DeteksiService {
	return &DeteksiService{repo: repo}
}

func (s *DeteksiService) Create(req dto.CreateDeteksiRequest) (*models.Deteksi, error) {
	id := uuid.New().String()
	if err := s.repo.Create(req, id); err != nil {
		rollbar.Error(err)
		return nil, err
	}
	return s.repo.GetByID(id)
}

func (s *DeteksiService) GetAll() ([]models.Deteksi, error) {
	return s.repo.GetAll()
}

func (s *DeteksiService) GetByID(id string) (*models.Deteksi, error) {
	return s.repo.GetByID(id)
}

func (s *DeteksiService) Update(id string, req dto.UpdateDeteksiRequest) (*models.Deteksi, error) {
	d, err := s.repo.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	if req.NilaiDeteksi != nil {
		d.NilaiDeteksi = *req.NilaiDeteksi
	}
	if req.NilaiSubdomain1 != nil {
		d.NilaiSubdomain1 = *req.NilaiSubdomain1
	}
	if req.NilaiSubdomain2 != nil {
		d.NilaiSubdomain2 = *req.NilaiSubdomain2
	}
	if req.NilaiSubdomain3 != nil {
		d.NilaiSubdomain3 = *req.NilaiSubdomain3
	}

	if err := s.repo.Update(id, *d); err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return d, nil
}

func (s *DeteksiService) Delete(id string) error {
	return s.repo.Delete(id)
}
