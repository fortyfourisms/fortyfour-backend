package services

import (
	"errors"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"

	"github.com/google/uuid"
	"github.com/rollbar/rollbar-go"
)

type JabatanService struct {
	repo repository.JabatanRepositoryInterface
}

func NewJabatanService(repo repository.JabatanRepositoryInterface) *JabatanService {
	return &JabatanService{repo: repo}
}

func (s *JabatanService) Create(req dto.CreateJabatanRequest) (*dto.JabatanResponse, error) {
	if req.NamaJabatan == nil || strings.TrimSpace(*req.NamaJabatan) == "" {
		return nil, errors.New("nama_jabatan wajib diisi")
	}

	id := uuid.New().String()

	if err := s.repo.Create(req, id); err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return s.repo.GetByID(id)
}

func (s *JabatanService) GetAll() ([]dto.JabatanResponse, error) {
	return s.repo.GetAll()
}

func (s *JabatanService) GetByID(id string) (*dto.JabatanResponse, error) {
	return s.repo.GetByID(id)
}

func (s *JabatanService) Update(id string, req dto.UpdateJabatanRequest) (*dto.JabatanResponse, error) {
	jabatan, err := s.repo.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	if req.NamaJabatan != nil {
		jabatan.NamaJabatan = *req.NamaJabatan
	}

	if err := s.repo.Update(id, *jabatan); err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return jabatan, nil
}

func (s *JabatanService) Delete(id string) error {
	return s.repo.Delete(id)
}
