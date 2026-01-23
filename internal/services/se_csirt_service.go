package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"

	"github.com/google/uuid"
	"github.com/rollbar/rollbar-go"
)

type SeCsirtService struct {
	repo repository.SeCsirtRepositoryInterface
}

func NewSeCsirtService(repo repository.SeCsirtRepositoryInterface) *SeCsirtService {
	return &SeCsirtService{repo: repo}
}

func (s *SeCsirtService) Create(req dto.CreateSeCsirtRequest) (string, error) {
	id := uuid.New().String()
	return id, s.repo.Create(req, id)
}

func (s *SeCsirtService) GetAll() ([]dto.SeCsirtResponse, error) {
	return s.repo.GetAll()
}

func (s *SeCsirtService) GetByID(id string) (*dto.SeCsirtResponse, error) {
	return s.repo.GetByID(id)
}

func (s *SeCsirtService) Update(id string, req dto.UpdateSeCsirtRequest) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	if req.NamaSe != nil {
		existing.NamaSe = *req.NamaSe
	}
	if req.IpSe != nil {
		existing.IpSe = *req.IpSe
	}
	if req.AsNumberSe != nil {
		existing.AsNumberSe = *req.AsNumberSe
	}
	if req.PengelolaSe != nil {
		existing.Pengelola = *req.PengelolaSe
	}
	if req.FiturSe != nil {
		existing.FiturSe = *req.FiturSe
	}
	if req.KategoriSe != nil {
		existing.KategoriSe = *req.KategoriSe
	}

	return s.repo.Update(id, *existing)
}

func (s *SeCsirtService) Delete(id string) error {
	return s.repo.Delete(id)
}
