package services

import (
	"errors"
	"ikas/internal/dto"
	"ikas/internal/repository"

	"github.com/google/uuid"
	"github.com/rollbar/rollbar-go"
)

type RuangLingkupService struct {
	repo repository.RuangLingkupRepositoryInterface
}

func NewRuangLingkupService(repo repository.RuangLingkupRepositoryInterface) *RuangLingkupService {
	return &RuangLingkupService{repo: repo}
}

func (s *RuangLingkupService) Create(req dto.CreateRuangLingkupRequest) (*dto.RuangLingkupResponse, error) {
	// Cek duplikat
	isDuplicate, err := s.repo.CheckDuplicateName(req.NamaRuangLingkup, "")
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	if isDuplicate {
		return nil, errors.New("nama_ruang_lingkup sudah ada")
	}

	// Generate UUID di service layer
	newID := uuid.New().String()

	if err := s.repo.Create(req, newID); err != nil {
		rollbar.Error(err)
		return nil, err
	}

	// Ambil data yang baru dibuat
	resp, err := s.repo.GetByID(newID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return resp, nil
}

func (s *RuangLingkupService) GetAll() ([]dto.RuangLingkupResponse, error) {
	return s.repo.GetAll()
}

func (s *RuangLingkupService) GetByID(id string) (*dto.RuangLingkupResponse, error) {
	return s.repo.GetByID(id)
}

func (s *RuangLingkupService) Update(id string, req dto.UpdateRuangLingkupRequest) (*dto.RuangLingkupResponse, error) {
	// Cek existing dulu
	_, err := s.repo.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	// Cek duplikat nama
	if req.NamaRuangLingkup != nil {
		isDuplicate, err := s.repo.CheckDuplicateName(*req.NamaRuangLingkup, id)
		if err != nil {
			rollbar.Error(err)
			return nil, err
		}
		if isDuplicate {
			return nil, errors.New("nama_ruang_lingkup sudah ada")
		}
	}

	// Update
	if err := s.repo.Update(id, req); err != nil {
		rollbar.Error(err)
		return nil, err
	}

	// Ambil data terbaru
	updated, err := s.repo.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return updated, nil
}

func (s *RuangLingkupService) Delete(id string) error {
	return s.repo.Delete(id)
}
