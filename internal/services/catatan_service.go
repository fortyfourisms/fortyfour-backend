package services

import (
	"errors"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"

	"github.com/google/uuid"
)

type CatatanService struct {
	repo repository.CatatanRepositoryInterface
}

func NewCatatanService(repo repository.CatatanRepositoryInterface) *CatatanService {
	return &CatatanService{repo: repo}
}

func (s *CatatanService) Upsert(idMateri, idUser string, req dto.UpsertCatatanRequest) (*dto.CatatanPribadiResponse, error) {
	if req.Konten == "" {
		return nil, errors.New("konten tidak boleh kosong")
	}

	catatan := &models.CatatanPribadi{
		ID:       uuid.New().String(),
		IDMateri: idMateri,
		IDUser:   idUser,
		Konten:   req.Konten,
	}

	if err := s.repo.Upsert(catatan); err != nil {
		return nil, err
	}

	// Ambil data terbaru (karena upsert bisa update, timestamps berubah)
	saved, err := s.repo.FindByUserAndMateri(idUser, idMateri)
	if err != nil {
		return nil, err
	}

	return mapCatatanToResponse(saved), nil
}

func (s *CatatanService) GetByUserAndMateri(idUser, idMateri string) (*dto.CatatanPribadiResponse, error) {
	catatan, err := s.repo.FindByUserAndMateri(idUser, idMateri)
	if err != nil {
		return nil, errors.New("catatan tidak ditemukan")
	}
	return mapCatatanToResponse(catatan), nil
}

func mapCatatanToResponse(c *models.CatatanPribadi) *dto.CatatanPribadiResponse {
	return &dto.CatatanPribadiResponse{
		ID:        c.ID,
		IDMateri:  c.IDMateri,
		Konten:    c.Konten,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
	}
}
