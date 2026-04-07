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
	repo         repository.KelasRepositoryInterface
	materiRepo   repository.MateriRepositoryInterface
	progressRepo repository.ProgressRepositoryInterface
	rc           cache.RedisInterface
}

func NewKelasService(
	repo repository.KelasRepositoryInterface,
	materiRepo repository.MateriRepositoryInterface,
	progressRepo repository.ProgressRepositoryInterface,
	rc cache.RedisInterface,
) *KelasService {
	return &KelasService{
		repo:         repo,
		materiRepo:   materiRepo,
		progressRepo: progressRepo,
		rc:           rc,
	}
}

// ── Helpers ──────────────────────────────────────────────────────────────────

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
		ID:          m.ID,
		IDKelas:     m.IDKelas,
		Judul:       m.Judul,
		Tipe:        m.Tipe,
		Urutan:      m.Urutan,
		YoutubeID:   m.YoutubeID,
		PDFPath:     m.PDFPath,
		DurasiDetik: m.DurasiDetik,
		CreatedAt:   m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   m.UpdatedAt.Format(time.RFC3339),
	}
}

// ── Admin: CRUD Kelas ─────────────────────────────────────────────────────────

func (s *KelasService) Create(req dto.CreateKelasRequest, createdBy string) (*dto.KelasResponse, error) {
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
		CreatedBy: createdBy,
	}

	if err := s.repo.Create(kelas); err != nil {
		return nil, err
	}

	cacheDelete(s.rc, keyList("kelas"))
	return mapKelasToResponse(kelas), nil
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

	cacheDelete(s.rc, keyDetail("kelas", id))
	cacheDelete(s.rc, keyList("kelas"))
	return mapKelasToResponse(kelas), nil
}

func (s *KelasService) Delete(id string) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return errors.New("kelas tidak ditemukan")
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	cacheDelete(s.rc, keyDetail("kelas", id))
	cacheDelete(s.rc, keyList("kelas"))
	return nil
}

// ── User & Admin: Read Kelas ──────────────────────────────────────────────────

// GetAll mengembalikan list kelas.
// Jika onlyPublished=true (untuk user biasa), hanya tampilkan yang published.
func (s *KelasService) GetAll(onlyPublished bool) ([]dto.KelasResponse, error) {
	key := keyList("kelas")
	var cached []dto.KelasResponse
	if cacheGet(s.rc, key, &cached) {
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

	cacheSet(s.rc, key, result, TTLList)
	return result, nil
}

// GetDetail mengembalikan detail kelas beserta daftar materi.
// Jika userID tidak kosong, progress user juga disertakan.
func (s *KelasService) GetDetail(id, userID string) (*dto.KelasResponse, error) {
	kelas, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("kelas tidak ditemukan")
	}

	resp := mapKelasToResponse(kelas)

	// Ambil semua materi dalam kelas
	materiList, err := s.materiRepo.FindByKelas(id)
	if err != nil {
		return nil, err
	}

	// Ambil progress user jika ada
	var progressMap map[string]*models.UserMateriProgress
	if userID != "" {
		progressList, err := s.progressRepo.FindByUserAndKelas(userID, id)
		if err == nil {
			progressMap = make(map[string]*models.UserMateriProgress, len(progressList))
			for i := range progressList {
				progressMap[progressList[i].IDMateri] = &progressList[i]
			}
		}
	}

	materiSelesai := 0
	kuisSelesai := false

	materiResponses := make([]dto.MateriResponse, 0, len(materiList))
	for _, m := range materiList {
		mr := mapMateriToResponse(&m)

		// Inject progress ke response materi
		if progressMap != nil {
			if p, ok := progressMap[m.ID]; ok {
				mr.IsCompleted = p.IsCompleted
				mr.LastWatchedSeconds = p.LastWatchedSeconds
			}
		}

		// Hitung ringkasan progress
		if mr.IsCompleted {
			if m.Tipe == models.MateriTipeVideo || m.Tipe == models.MateriTipePDF {
				materiSelesai++
			}
			if m.Tipe == models.MateriTipeKuis {
				kuisSelesai = true
			}
		}

		materiResponses = append(materiResponses, mr)
	}

	resp.Materi = materiResponses

	// Ringkasan progress keseluruhan kelas
	if userID != "" {
		totalMedia := 0
		for _, m := range materiList {
			if m.Tipe != models.MateriTipeKuis {
				totalMedia++
			}
		}
		isSelesai := materiSelesai == totalMedia && totalMedia > 0 && kuisSelesai
		resp.Progress = &dto.KelasProgress{
			TotalMateri:    totalMedia,
			MateriSelesai:  materiSelesai,
			KuisSelesai:    kuisSelesai,
			IsKelasSelesai: isSelesai,
		}
	}

	return resp, nil
}
