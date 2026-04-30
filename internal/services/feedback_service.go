package services

import (
	"errors"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"

	"github.com/google/uuid"
)

type FeedbackService struct {
	repo repository.FeedbackRepositoryInterface
}

func NewFeedbackService(repo repository.FeedbackRepositoryInterface) *FeedbackService {
	return &FeedbackService{repo: repo}
}

func (s *FeedbackService) Upsert(idMateri, idUser string, req dto.UpsertFeedbackRequest) (*dto.FeedbackResponse, error) {
	if req.Konten == "" {
		return nil, errors.New("konten tidak boleh kosong")
	}

	feedback := &models.Feedback{
		ID:       uuid.New().String(),
		IDMateri: idMateri,
		IDUser:   idUser,
		Konten:   req.Konten,
	}

	if err := s.repo.Upsert(feedback); err != nil {
		return nil, err
	}

	// Ambil data terbaru (karena upsert bisa update, timestamps berubah)
	saved, err := s.repo.FindByUserAndMateri(idUser, idMateri)
	if err != nil {
		return nil, err
	}

	return mapFeedbackToResponse(saved), nil
}

func (s *FeedbackService) GetByUserAndMateri(idUser, idMateri string) (*dto.FeedbackResponse, error) {
	feedback, err := s.repo.FindByUserAndMateri(idUser, idMateri)
	if err != nil {
		return nil, errors.New("feedback tidak ditemukan")
	}
	return mapFeedbackToResponse(feedback), nil
}

// GetAllByMateri mengembalikan semua feedback untuk materi tertentu (untuk admin/staff).
func (s *FeedbackService) GetAllByMateri(idMateri string) ([]dto.FeedbackListItem, error) {
	items, err := s.repo.FindByMateri(idMateri)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []dto.FeedbackListItem{}
	}
	return items, nil
}

func mapFeedbackToResponse(f *models.Feedback) *dto.FeedbackResponse {
	return &dto.FeedbackResponse{
		ID:        f.ID,
		IDMateri:  f.IDMateri,
		Konten:    f.Konten,
		CreatedAt: f.CreatedAt.Format(time.RFC3339),
		UpdatedAt: f.UpdatedAt.Format(time.RFC3339),
	}
}
