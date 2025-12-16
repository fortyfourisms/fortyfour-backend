package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"

	"github.com/google/uuid"
)

type SdmCsirtService struct {
    repo *repository.SdmCsirtRepository
}

func NewSdmCsirtService(repo *repository.SdmCsirtRepository) *SdmCsirtService {
    return &SdmCsirtService{repo: repo}
}

func (s *SdmCsirtService) Create(req dto.CreateSdmCsirtRequest) (string, error) {
    id := uuid.New().String()
    if err := s.repo.Create(req, id); err != nil {
        return "", err
    }
    return id, nil
}

func (s *SdmCsirtService) GetAll() ([]dto.SdmCsirtResponse, error) {
    return s.repo.GetAll()
}

func (s *SdmCsirtService) GetByID(id string) (*dto.SdmCsirtResponse, error) {
    return s.repo.GetByID(id)
}

func (s *SdmCsirtService) Update(id string, req dto.SdmCsirtResponse) error {
    return s.repo.Update(id, req)
}

func (s *SdmCsirtService) Delete(id string) error {
    return s.repo.Delete(id)
}
