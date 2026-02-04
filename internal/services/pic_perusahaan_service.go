package services

import (
	"errors"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"
	"strings"

	"github.com/google/uuid"
	"github.com/rollbar/rollbar-go"
)

type PICService struct {
	repo repository.PICRepositoryInterface
}

func NewPICService(repo repository.PICRepositoryInterface) *PICService {
	return &PICService{repo: repo}
}

func (s *PICService) Create(req dto.CreatePICRequest) (*dto.PICResponse, error) {

	if req.Nama == nil || strings.TrimSpace(*req.Nama) == "" {
		return nil, errors.New("nama wajib diisi")
	}

	if req.IDPerusahaan == nil || strings.TrimSpace(*req.IDPerusahaan) == "" {
		return nil, errors.New("id_perusahaan wajib diisi")
	}

	id := uuid.New().String()

	if err := s.repo.Create(req, id); err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return s.repo.GetByID(id)
}

func (s *PICService) GetAll() ([]dto.PICResponse, error) {
	return s.repo.GetAll()
}

func (s *PICService) GetByID(id string) (*dto.PICResponse, error) {
	return s.repo.GetByID(id)
}

func (s *PICService) Update(id string, req dto.UpdatePICRequest) (*dto.PICResponse, error) {
	if err := s.repo.Update(id, req); err != nil {
		rollbar.Error(err)
		return nil, err
	}

	updated, err := s.repo.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return updated, nil
}

func (s *PICService) Delete(id string) error {
	return s.repo.Delete(id)
}
