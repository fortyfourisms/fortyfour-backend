package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"

	"github.com/google/uuid"
)

type SeCsirtService struct {
    repo *repository.SeCsirtRepository
}

func NewSeCsirtService(repo *repository.SeCsirtRepository) *SeCsirtService {
    return &SeCsirtService{repo: repo}
}

func (s *SeCsirtService) Create(req dto.CreateSeCsirtRequest) (string, error) {
    id := uuid.New().String()
    if err := s.repo.Create(req, id); err != nil {
        return "", err
    }
    return id, nil
}

func (s *SeCsirtService) GetAll() ([]dto.SeCsirtResponse, error) {
    return s.repo.GetAll()
}

func (s *SeCsirtService) GetByID(id string) (*dto.SeCsirtResponse, error) {
    return s.repo.GetByID(id)
}

func (s *SeCsirtService) Update(id string, req dto.SeCsirtResponse) error {
    return s.repo.Update(id, req)
}

func (s *SeCsirtService) Delete(id string) error {
    return s.repo.Delete(id)
}
