package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"

	"github.com/google/uuid"
)

type IkasService struct {
	repo *repository.IkasRepository
}

func NewIkasService(repo *repository.IkasRepository) *IkasService {
	return &IkasService{repo: repo}
}

func (s *IkasService) Create(req dto.CreateIkasRequest) (*models.Ikas, error) {
	id := uuid.New().String()

	if err := s.repo.Create(req, id); err != nil {
		return nil, err
	}

	return s.repo.GetByID(id)
}

func (s *IkasService) GetAll() ([]models.Ikas, error) {
	return s.repo.GetAll()
}

func (s *IkasService) GetByID(id string) (*models.Ikas, error) {
	return s.repo.GetByID(id)
}

func (s *IkasService) Update(id string, req dto.UpdateIkasRequest) (*models.Ikas, error) {
	ikas, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update hanya field yang dikirim (non-nil)
	if req.IDStakeholder != nil {
		ikas.IDStakeholder = *req.IDStakeholder
	}
	if req.Tanggal != nil {
		ikas.Tanggal = *req.Tanggal
	}
	if req.Responden != nil {
		ikas.Responden = *req.Responden
	}
	if req.Telepon != nil {
		ikas.Telepon = *req.Telepon
	}
	if req.Jabatan != nil {
		ikas.Jabatan = *req.Jabatan
	}
	if req.NilaiKematangan != nil {
		ikas.NilaiKematangan = *req.NilaiKematangan
	}
	if req.TargetNilai != nil {
		ikas.TargetNilai = *req.TargetNilai
	}
	if req.IDIdentifikasi != nil {
		ikas.IDIdentifikasi = *req.IDIdentifikasi
	}
	if req.IDProteksi != nil {
		ikas.IDProteksi = *req.IDProteksi
	}
	if req.IDDeteksi != nil {
		ikas.IDDeteksi = *req.IDDeteksi
	}
	if req.IDGulih != nil {
		ikas.IDGulih = *req.IDGulih
	}

	if err := s.repo.Update(id, *ikas); err != nil {
		return nil, err
	}

	return ikas, nil
}

func (s *IkasService) Delete(id string) error {
	return s.repo.Delete(id)
}