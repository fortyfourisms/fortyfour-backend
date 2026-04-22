package services

import (
	"errors"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/pkg/cache"

	"github.com/google/uuid"
)

type FilePendukungService struct {
	repo       repository.FilePendukungRepositoryInterface
	materiRepo repository.MateriRepositoryInterface
	rc         cache.RedisInterface
}

func NewFilePendukungService(
	repo repository.FilePendukungRepositoryInterface,
	materiRepo repository.MateriRepositoryInterface,
	rc cache.RedisInterface,
) *FilePendukungService {
	return &FilePendukungService{repo: repo, materiRepo: materiRepo, rc: rc}
}

func (s *FilePendukungService) Create(idMateri, namaFile, filePath string, ukuran int64) (*dto.FilePendukungResponse, error) {
	// Pastikan materi ada
	if _, err := s.materiRepo.FindByID(idMateri); err != nil {
		return nil, errors.New("materi tidak ditemukan")
	}

	fp := &models.FilePendukung{
		ID:       uuid.New().String(),
		IDMateri: idMateri,
		NamaFile: namaFile,
		FilePath: filePath,
		Ukuran:   ukuran,
	}

	if err := s.repo.Create(fp); err != nil {
		return nil, err
	}

	return mapFilePendukungToResponse(fp), nil
}

func (s *FilePendukungService) GetByMateri(idMateri string) ([]dto.FilePendukungResponse, error) {
	list, err := s.repo.FindByMateri(idMateri)
	if err != nil {
		return nil, err
	}

	result := make([]dto.FilePendukungResponse, 0, len(list))
	for _, fp := range list {
		fp := fp
		result = append(result, *mapFilePendukungToResponse(&fp))
	}
	return result, nil
}

func (s *FilePendukungService) Delete(id string) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return errors.New("file pendukung tidak ditemukan")
	}
	return s.repo.Delete(id)
}

func (s *FilePendukungService) FindByID(id string) (*models.FilePendukung, error) {
	return s.repo.FindByID(id)
}

func mapFilePendukungToResponse(fp *models.FilePendukung) *dto.FilePendukungResponse {
	return &dto.FilePendukungResponse{
		ID:        fp.ID,
		IDMateri:  fp.IDMateri,
		NamaFile:  fp.NamaFile,
		FilePath:  fp.FilePath,
		Ukuran:    fp.Ukuran,
		CreatedAt: fp.CreatedAt.Format(time.RFC3339),
	}
}
