package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"

	"github.com/google/uuid"
)

type SdmCsirtService struct {
	repo repository.SdmCsirtRepositoryInterface
}

func NewSdmCsirtService(repo repository.SdmCsirtRepositoryInterface) *SdmCsirtService {
	return &SdmCsirtService{repo: repo}
}

func (s *SdmCsirtService) Create(req dto.CreateSdmCsirtRequest) (string, error) {
	id := uuid.New().String()
	return id, s.repo.Create(req, id)
}

func (s *SdmCsirtService) GetAll() ([]dto.SdmCsirtResponse, error) {
	return s.repo.GetAll()
}

func (s *SdmCsirtService) GetByID(id string) (*dto.SdmCsirtResponse, error) {
	return s.repo.GetByID(id)
}

func (s *SdmCsirtService) Update(id string, req dto.UpdateSdmCsirtRequest) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if req.NamaPersonel != nil {
		existing.NamaPersonel = *req.NamaPersonel
	}
	if req.JabatanCsirt != nil {
		existing.JabatanCsirt = *req.JabatanCsirt
	}
	if req.JabatanPerusahaan != nil {
		existing.JabatanPerusahaan = *req.JabatanPerusahaan
	}
	if req.Skill != nil {
		existing.Skill = *req.Skill
	}
	if req.Sertifikasi != nil {
		existing.Sertifikasi = *req.Sertifikasi
	}

	return s.repo.Update(id, *existing)
}

func (s *SdmCsirtService) Delete(id string) error {
	return s.repo.Delete(id)
}
