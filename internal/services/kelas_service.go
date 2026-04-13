package services

import (
	"errors"
	"strings"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/pkg/cache"

	"github.com/google/uuid"
)

type KelasService struct {
	repo            repository.KelasRepositoryInterface
	materiRepo      repository.MateriRepositoryInterface
	progressRepo    repository.ProgressRepositoryInterface
	kuisRepo        repository.KuisRepositoryInterface
	attemptRepo     repository.KuisAttemptRepositoryInterface
	sertifikatRepo  repository.SertifikatRepositoryInterface
	fpRepo          repository.FilePendukungRepositoryInterface
	rc              cache.RedisInterface
}

func NewKelasService(
	repo repository.KelasRepositoryInterface,
	materiRepo repository.MateriRepositoryInterface,
	progressRepo repository.ProgressRepositoryInterface,
	kuisRepo repository.KuisRepositoryInterface,
	attemptRepo repository.KuisAttemptRepositoryInterface,
	sertifikatRepo repository.SertifikatRepositoryInterface,
	fpRepo repository.FilePendukungRepositoryInterface,
	rc cache.RedisInterface,
) *KelasService {
	return &KelasService{
		repo:           repo,
		materiRepo:     materiRepo,
		progressRepo:   progressRepo,
		kuisRepo:       kuisRepo,
		attemptRepo:    attemptRepo,
		sertifikatRepo: sertifikatRepo,
		fpRepo:         fpRepo,
		rc:             rc,
	}
}

// ── Admin: CRUD Kelas ─────────────────────────────────────────────────────────

func (s *KelasService) Create(req dto.CreateKelasRequest, userID string) (*dto.KelasResponse, error) {
	judul := strings.TrimSpace(req.Judul)
	if judul == "" {
		return nil, errors.New("judul wajib diisi")
	}

	kelas := &models.Kelas{
		ID:        uuid.New().String(),
		Judul:     judul,
		Deskripsi: req.Deskripsi,
		Thumbnail: req.Thumbnail,
		Status:    models.KelasStatusDraft,
		CreatedBy: userID,
	}

	if err := s.repo.Create(kelas); err != nil {
		return nil, err
	}

	cacheDelete(s.rc, keyList("kelas"))
	resp := mapKelasToResponse(kelas)
	return resp, nil
}

func (s *KelasService) Update(id string, req dto.UpdateKelasRequest) (*dto.KelasResponse, error) {
	kelas, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("kelas tidak ditemukan")
	}

	if req.Judul != nil {
		trimmed := strings.TrimSpace(*req.Judul)
		if trimmed == "" {
			return nil, errors.New("judul tidak boleh kosong")
		}
		kelas.Judul = trimmed
	}
	if req.Deskripsi != nil {
		kelas.Deskripsi = req.Deskripsi
	}
	if req.Thumbnail != nil {
		kelas.Thumbnail = req.Thumbnail
	}
	if req.Status != nil {
		kelas.Status = models.KelasStatus(*req.Status)
	}

	if err := s.repo.Update(kelas); err != nil {
		return nil, err
	}

	cacheDelete(s.rc, keyList("kelas"))
	cacheDelete(s.rc, keyDetail("kelas", id))
	resp := mapKelasToResponse(kelas)
	return resp, nil
}

func (s *KelasService) Delete(id string) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return errors.New("kelas tidak ditemukan")
	}
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	cacheDelete(s.rc, keyList("kelas"))
	cacheDelete(s.rc, keyDetail("kelas", id))
	return nil
}

// ── GetAll ────────────────────────────────────────────────────────────────────

func (s *KelasService) GetAll(onlyPublished bool) ([]dto.KelasResponse, error) {
	cacheKey := keyList("kelas")
	if onlyPublished {
		cacheKey = keyList("kelas:published")
	}

	var cached []dto.KelasResponse
	if cacheGet(s.rc, cacheKey, &cached) {
		return cached, nil
	}

	kelasList, err := s.repo.FindAll(onlyPublished)
	if err != nil {
		return nil, err
	}

	result := make([]dto.KelasResponse, 0, len(kelasList))
	for _, k := range kelasList {
		k := k
		result = append(result, *mapKelasToResponse(&k))
	}

	cacheSet(s.rc, cacheKey, result, TTLList)
	return result, nil
}

// ── GetDetail (dengan materi, kuis, progress, sertifikat) ─────────────────────

func (s *KelasService) GetDetail(id, userID string) (*dto.KelasResponse, error) {
	kelas, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("kelas tidak ditemukan")
	}

	resp := mapKelasToResponse(kelas)

	// Load materi
	materiList, err := s.materiRepo.FindByKelas(id)
	if err == nil {
		// Load progress user
		progressMap := make(map[string]*models.UserMateriProgress)
		if userID != "" {
			progressList, err := s.progressRepo.FindByUserAndKelas(userID, id)
			if err == nil {
				for i := range progressList {
					progressMap[progressList[i].IDMateri] = &progressList[i]
				}
			}
		}

		materiResponses := make([]dto.MateriResponse, 0, len(materiList))
		for _, m := range materiList {
			m := m
			mr := mapMateriToResponse(&m)

			// Inject progress
			if p, ok := progressMap[m.ID]; ok {
				mr.IsCompleted = p.IsCompleted
				mr.LastWatchedSeconds = p.LastWatchedSeconds
			}

			// Load file pendukung
			fps, err := s.fpRepo.FindByMateri(m.ID)
			if err == nil && len(fps) > 0 {
				fpResponses := make([]dto.FilePendukungResponse, 0, len(fps))
				for _, fp := range fps {
					fpResponses = append(fpResponses, dto.FilePendukungResponse{
						ID:        fp.ID,
						IDMateri:  fp.IDMateri,
						NamaFile:  fp.NamaFile,
						FilePath:  fp.FilePath,
						Ukuran:    fp.Ukuran,
						CreatedAt: fp.CreatedAt.Format(time.RFC3339),
					})
				}
				mr.FilePendukung = fpResponses
			}

			// Load kuis per-materi
			kuis, err := s.kuisRepo.FindByMateri(m.ID)
			if err == nil && kuis != nil {
				kr := mapKuisToResponse(kuis)
				mr.Kuis = kr
			}

			materiResponses = append(materiResponses, mr)
		}
		resp.Materi = materiResponses
	}

	// Load kuis list
	kuisList, err := s.kuisRepo.FindByKelas(id)
	if err == nil {
		kuisResponses := make([]dto.KuisResponse, 0, len(kuisList))
		for _, k := range kuisList {
			k := k
			kuisResponses = append(kuisResponses, *mapKuisToResponse(&k))
		}
		resp.KuisList = kuisResponses
	}

	// Hitung progress
	if userID != "" {
		progress := s.calculateProgress(userID, id, materiList, kuisList)
		resp.Progress = progress

		// Cek sertifikat
		cert, err := s.sertifikatRepo.FindByUserAndKelas(userID, id)
		if err == nil && cert != nil {
			resp.Sertifikat = &dto.SertifikatResponse{
				ID:              cert.ID,
				NomorSertifikat: cert.NomorSertifikat,
				IDKelas:         cert.IDKelas,
				IDUser:          cert.IDUser,
				NamaPeserta:     cert.NamaPeserta,
				NamaKelas:       cert.NamaKelas,
				TanggalTerbit:   cert.TanggalTerbit.Format("2006-01-02"),
				PDFPath:         cert.PDFPath,
				CreatedAt:       cert.CreatedAt.Format(time.RFC3339),
			}
		}
	}

	return resp, nil
}

// ── Progress Calculation ──────────────────────────────────────────────────────

func (s *KelasService) calculateProgress(userID, kelasID string, materiList []models.Materi, kuisList []models.Kuis) *dto.KelasProgress {
	progress := &dto.KelasProgress{
		TotalMateri: len(materiList),
	}

	// Materi selesai
	progressList, err := s.progressRepo.FindByUserAndKelas(userID, kelasID)
	if err == nil {
		for _, p := range progressList {
			if p.IsCompleted {
				progress.MateriSelesai++
			}
		}
	}

	// Kuis
	totalKuis := 0
	kuisLulus := 0
	kuisAkhirLulus := false

	for _, k := range kuisList {
		if k.IsFinal {
			// Cek kuis akhir
			attempts, err := s.attemptRepo.FindByUserAndKuis(userID, k.ID)
			if err == nil {
				for _, a := range attempts {
					if a.IsPassed {
						kuisAkhirLulus = true
						break
					}
				}
			}
		} else {
			totalKuis++
			// Cek apakah sudah lulus
			attempts, err := s.attemptRepo.FindByUserAndKuis(userID, k.ID)
			if err == nil {
				for _, a := range attempts {
					if a.IsPassed {
						kuisLulus++
						break
					}
				}
			}
		}
	}

	progress.TotalKuis = totalKuis
	progress.KuisLulus = kuisLulus
	progress.KuisAkhirLulus = kuisAkhirLulus

	// Persentase
	totalItems := progress.TotalMateri + progress.TotalKuis
	if totalItems > 0 {
		completedItems := progress.MateriSelesai + progress.KuisLulus
		progress.PersentaseProgress = float64(completedItems) / float64(totalItems) * 100
	}

	// Kelas selesai = semua materi + semua kuis + kuis akhir lulus
	progress.IsKelasSelesai = progress.MateriSelesai >= progress.TotalMateri &&
		progress.KuisLulus >= progress.TotalKuis &&
		kuisAkhirLulus

	return progress
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func mapKelasToResponse(k *models.Kelas) *dto.KelasResponse {
	return &dto.KelasResponse{
		ID:        k.ID,
		Judul:     k.Judul,
		Deskripsi: k.Deskripsi,
		Thumbnail: k.Thumbnail,
		Status:    k.Status,
		CreatedBy: k.CreatedBy,
		CreatedAt: k.CreatedAt.Format(time.RFC3339),
		UpdatedAt: k.UpdatedAt.Format(time.RFC3339),
	}
}

func mapMateriToResponse(m *models.Materi) dto.MateriResponse {
	return dto.MateriResponse{
		ID:               m.ID,
		IDKelas:          m.IDKelas,
		Judul:            m.Judul,
		Tipe:             m.Tipe,
		Urutan:           m.Urutan,
		YoutubeID:        m.YoutubeID,
		DurasiDetik:      m.DurasiDetik,
		KontenHTML:       m.KontenHTML,
		DeskripsiSingkat: m.DeskripsiSingkat,
		Kategori:         m.Kategori,
		CreatedAt:        m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        m.UpdatedAt.Format(time.RFC3339),
	}
}
