package services

import (
	"errors"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
	"time"
)

type BeritaServiceInterface interface {
	Create(authorID string, req dto.CreateBeritaRequest) (*dto.BeritaResponse, error)
	GetAll() ([]dto.BeritaResponse, error)
	GetByID(id int64) (*dto.BeritaResponse, error)
	Update(id int64, req dto.UpdateBeritaRequest) (*dto.BeritaResponse, error)
	Delete(id int64) error
}

type BeritaService struct {
	repo repository.BeritaRepositoryInterface
}

func NewBeritaService(repo repository.BeritaRepositoryInterface) *BeritaService {
	return &BeritaService{repo: repo}
}

var _ BeritaServiceInterface = (*BeritaService)(nil)

func (s *BeritaService) Create(authorID string, req dto.CreateBeritaRequest) (*dto.BeritaResponse, error) {
	berita := &models.Berita{
		Judul:     req.Judul,
		Deskripsi: req.Deskripsi,
		AuthorID:  authorID,
	}

	if err := s.repo.Create(berita); err != nil {
		return nil, err
	}

	// Fetch again to get author details
	saved, err := s.repo.FindByID(berita.ID)
	if err != nil {
		return nil, err
	}

	return mapBeritaToResponse(saved), nil
}

func (s *BeritaService) GetAll() ([]dto.BeritaResponse, error) {
	list, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var res []dto.BeritaResponse
	for _, b := range list {
		res = append(res, *mapBeritaToResponse(&b))
	}
	return res, nil
}

func (s *BeritaService) GetByID(id int64) (*dto.BeritaResponse, error) {
	b, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, errors.New("berita tidak ditemukan")
	}
	return mapBeritaToResponse(b), nil
}

func (s *BeritaService) Update(id int64, req dto.UpdateBeritaRequest) (*dto.BeritaResponse, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("berita tidak ditemukan")
	}

	if req.Judul != nil {
		existing.Judul = *req.Judul
	}
	if req.Deskripsi != nil {
		existing.Deskripsi = *req.Deskripsi
	}

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return mapBeritaToResponse(existing), nil
}

func (s *BeritaService) Delete(id int64) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("berita tidak ditemukan")
	}

	return s.repo.Delete(id)
}

func mapBeritaToResponse(b *models.Berita) *dto.BeritaResponse {
	res := &dto.BeritaResponse{
		ID:        b.ID,
		Judul:     b.Judul,
		Deskripsi: b.Deskripsi,
		AuthorID:  b.AuthorID,
		CreatedAt: b.CreatedAt.Format(time.RFC3339),
		UpdatedAt: b.UpdatedAt.Format(time.RFC3339),
	}

	if b.Author != nil {
		res.Author = &dto.BeritaAuthor{
			ID:          b.Author.ID,
			Username:    b.Author.Username,
			DisplayName: b.Author.DisplayName,
		}
	}

	return res
}
