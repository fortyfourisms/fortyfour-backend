package services

import (
	"errors"
	"fmt"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/utils"

	"github.com/google/uuid"
)

type SertifikatService struct {
	repo         repository.SertifikatRepositoryInterface
	kelasRepo    repository.KelasRepositoryInterface
	progressRepo repository.ProgressRepositoryInterface
	attemptRepo  repository.KuisAttemptRepositoryInterface
	kuisRepo     repository.KuisRepositoryInterface
	userRepo     repository.UserRepositoryInterface
}

func NewSertifikatService(
	repo repository.SertifikatRepositoryInterface,
	kelasRepo repository.KelasRepositoryInterface,
	progressRepo repository.ProgressRepositoryInterface,
	attemptRepo repository.KuisAttemptRepositoryInterface,
	kuisRepo repository.KuisRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
) *SertifikatService {
	return &SertifikatService{
		repo:         repo,
		kelasRepo:    kelasRepo,
		progressRepo: progressRepo,
		attemptRepo:  attemptRepo,
		kuisRepo:     kuisRepo,
		userRepo:     userRepo,
	}
}

// Generate mengecek eligibility dan meng-generate sertifikat.
func (s *SertifikatService) Generate(userID, kelasID string) (*dto.SertifikatResponse, error) {
	// Cek apakah sudah ada sertifikat
	existing, err := s.repo.FindByUserAndKelas(userID, kelasID)
	if err == nil && existing != nil {
		return mapSertifikatToResponse(existing), nil
	}

	// Cek kelas ada
	kelas, err := s.kelasRepo.FindByID(kelasID)
	if err != nil {
		return nil, errors.New("kelas tidak ditemukan")
	}

	// Cek user ada
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	// 1. Cek semua materi selesai
	allMateri, err := s.progressRepo.HasCompletedAllMateri(userID, kelasID)
	if err != nil {
		return nil, err
	}
	if !allMateri {
		return nil, errors.New("selesaikan semua materi terlebih dahulu")
	}

	// 2. Cek semua kuis non-final lulus
	allKuis, err := s.attemptRepo.HasPassedAllKuisInKelas(userID, kelasID)
	if err != nil {
		return nil, err
	}
	if !allKuis {
		return nil, errors.New("lulus semua kuis per-materi terlebih dahulu")
	}

	// 3. Cek kuis akhir lulus (jika ada)
	finalKuis, err := s.kuisRepo.FindFinalByKelas(kelasID)
	if err == nil && finalKuis != nil {
		attempts, err := s.attemptRepo.FindByUserAndKuis(userID, finalKuis.ID)
		if err != nil || len(attempts) == 0 {
			return nil, errors.New("selesaikan kuis akhir terlebih dahulu")
		}
		hasPassed := false
		for _, a := range attempts {
			if a.IsPassed {
				hasPassed = true
				break
			}
		}
		if !hasPassed {
			return nil, errors.New("anda belum lulus kuis akhir")
		}
	}

	// Generate sertifikat
	namaPeserta := user.Username
	if user.DisplayName != nil && *user.DisplayName != "" {
		namaPeserta = *user.DisplayName
	}

	now := time.Now()
	nomorSertifikat := generateNomorSertifikat(now)

	sertifikat := &models.Sertifikat{
		ID:              uuid.New().String(),
		NomorSertifikat: nomorSertifikat,
		IDKelas:         kelasID,
		IDUser:          userID,
		NamaPeserta:     namaPeserta,
		NamaKelas:       kelas.Judul,
		TanggalTerbit:   now,
	}

	// Generate PDF
	pdfPath, err := utils.GenerateSertifikatPDF(sertifikat)
	if err == nil && pdfPath != "" {
		sertifikat.PDFPath = &pdfPath
	}

	if err := s.repo.Create(sertifikat); err != nil {
		return nil, err
	}

	return mapSertifikatToResponse(sertifikat), nil
}

func (s *SertifikatService) GetByUserAndKelas(userID, kelasID string) (*dto.SertifikatResponse, error) {
	cert, err := s.repo.FindByUserAndKelas(userID, kelasID)
	if err != nil {
		return nil, errors.New("sertifikat belum tersedia")
	}
	return mapSertifikatToResponse(cert), nil
}

func (s *SertifikatService) GetByID(id string) (*dto.SertifikatResponse, error) {
	cert, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("sertifikat tidak ditemukan")
	}
	return mapSertifikatToResponse(cert), nil
}

func (s *SertifikatService) GetByUser(userID string) ([]dto.SertifikatResponse, error) {
	list, err := s.repo.FindByUser(userID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.SertifikatResponse, 0, len(list))
	for _, c := range list {
		c := c
		result = append(result, *mapSertifikatToResponse(&c))
	}
	return result, nil
}

func mapSertifikatToResponse(s *models.Sertifikat) *dto.SertifikatResponse {
	resp := &dto.SertifikatResponse{
		ID:              s.ID,
		NomorSertifikat: s.NomorSertifikat,
		IDKelas:         s.IDKelas,
		IDUser:          s.IDUser,
		NamaPeserta:     s.NamaPeserta,
		NamaKelas:       s.NamaKelas,
		TanggalTerbit:   s.TanggalTerbit.Format("2006-01-02"),
		PDFPath:         s.PDFPath,
		CreatedAt:       s.CreatedAt.Format(time.RFC3339),
	}
	return resp
}

func generateNomorSertifikat(t time.Time) string {
	return fmt.Sprintf("CERT/%s/%s", t.Format("20060102"), uuid.New().String()[:8])
}
