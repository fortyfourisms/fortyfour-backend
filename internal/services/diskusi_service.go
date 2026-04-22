package services

import (
	"errors"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"

	"github.com/google/uuid"
)

type DiskusiService struct {
	repo     repository.DiskusiRepositoryInterface
	userRepo repository.UserRepositoryInterface
}

func NewDiskusiService(
	repo repository.DiskusiRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
) *DiskusiService {
	return &DiskusiService{repo: repo, userRepo: userRepo}
}

func (s *DiskusiService) Create(idMateri, idUser string, req dto.CreateDiskusiRequest) (*dto.DiskusiResponse, error) {
	diskusi := &models.Diskusi{
		ID:       uuid.New().String(),
		IDMateri: idMateri,
		IDUser:   idUser,
		IDParent: req.IDParent,
		Konten:   req.Konten,
	}

	if err := s.repo.Create(diskusi); err != nil {
		return nil, err
	}

	return s.buildResponse(diskusi)
}

func (s *DiskusiService) GetByMateri(idMateri string) ([]dto.DiskusiResponse, error) {
	list, err := s.repo.FindByMateri(idMateri)
	if err != nil {
		return nil, err
	}

	result := make([]dto.DiskusiResponse, 0, len(list))
	for _, d := range list {
		d := d
		resp, err := s.buildResponse(&d)
		if err != nil {
			continue
		}
		// Load replies
		replies, err := s.repo.FindReplies(d.ID)
		if err == nil {
			replyResponses := make([]dto.DiskusiResponse, 0, len(replies))
			for _, r := range replies {
				r := r
				rResp, err := s.buildResponse(&r)
				if err != nil {
					continue
				}
				replyResponses = append(replyResponses, *rResp)
			}
			resp.Replies = replyResponses
		}
		result = append(result, *resp)
	}
	return result, nil
}

func (s *DiskusiService) Update(id, userID string, req dto.UpdateDiskusiRequest) (*dto.DiskusiResponse, error) {
	diskusi, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("diskusi tidak ditemukan")
	}
	if diskusi.IDUser != userID {
		return nil, errors.New("anda hanya bisa mengedit diskusi milik sendiri")
	}

	diskusi.Konten = req.Konten
	if err := s.repo.Update(diskusi); err != nil {
		return nil, err
	}

	return s.buildResponse(diskusi)
}

func (s *DiskusiService) Delete(id, userID, role string) error {
	diskusi, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("diskusi tidak ditemukan")
	}
	// Admin bisa hapus semua, user hanya miliknya
	if role != "admin" && diskusi.IDUser != userID {
		return errors.New("anda hanya bisa menghapus diskusi milik sendiri")
	}
	return s.repo.Delete(id)
}

func (s *DiskusiService) buildResponse(d *models.Diskusi) (*dto.DiskusiResponse, error) {
	namaUser := ""
	var fotoProfile *string
	user, err := s.userRepo.FindByID(d.IDUser)
	if err == nil {
		if user.DisplayName != nil {
			namaUser = *user.DisplayName
		} else {
			namaUser = user.Username
		}
		fotoProfile = user.FotoProfile
	}

	return &dto.DiskusiResponse{
		ID:          d.ID,
		IDMateri:    d.IDMateri,
		IDUser:      d.IDUser,
		NamaUser:    namaUser,
		FotoProfile: fotoProfile,
		IDParent:    d.IDParent,
		Konten:      d.Konten,
		CreatedAt:   d.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   d.UpdatedAt.Format(time.RFC3339),
	}, nil
}
