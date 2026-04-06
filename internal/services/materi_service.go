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

type MateriService struct {
	repo         repository.MateriRepositoryInterface
	kelasRepo    repository.KelasRepositoryInterface
	progressRepo repository.ProgressRepositoryInterface
	rc           cache.RedisInterface
}

func NewMateriService(
	repo repository.MateriRepositoryInterface,
	kelasRepo repository.KelasRepositoryInterface,
	progressRepo repository.ProgressRepositoryInterface,
	rc cache.RedisInterface,
) *MateriService {
	return &MateriService{
		repo:         repo,
		kelasRepo:    kelasRepo,
		progressRepo: progressRepo,
		rc:           rc,
	}
}

// ── Admin: CRUD Materi ────────────────────────────────────────────────────────

func (s *MateriService) Create(idKelas string, req dto.CreateMateriRequest) (*dto.MateriResponse, error) {
	// Pastikan kelas ada
	if _, err := s.kelasRepo.FindByID(idKelas); err != nil {
		return nil, errors.New("kelas tidak ditemukan")
	}

	judul := strings.TrimSpace(req.Judul)
	if judul == "" {
		return nil, errors.New("judul wajib diisi")
	}

	tipe := models.MateriTipe(req.Tipe)

	// Validasi field wajib sesuai tipe
	switch tipe {
	case models.MateriTipeVideo:
		if req.YoutubeID == nil || strings.TrimSpace(*req.YoutubeID) == "" {
			return nil, errors.New("youtube_id wajib diisi untuk tipe video")
		}
	case models.MateriTipePDF:
		if req.PDFPath == nil || strings.TrimSpace(*req.PDFPath) == "" {
			return nil, errors.New("pdf_path wajib diisi untuk tipe pdf")
		}
	case models.MateriTipeKuis:
		// kuis tidak butuh field tambahan saat create materi,
		// soal ditambahkan via endpoint terpisah
	default:
		return nil, errors.New("tipe tidak valid, harus: video, pdf, atau kuis")
	}

	materi := &models.Materi{
		ID:          uuid.New().String(),
		IDKelas:     idKelas,
		Judul:       judul,
		Tipe:        tipe,
		Urutan:      req.Urutan,
		YoutubeID:   req.YoutubeID,
		PDFPath:     req.PDFPath,
		DurasiDetik: req.DurasiDetik,
	}

	if err := s.repo.Create(materi); err != nil {
		return nil, err
	}

	cacheDelete(s.rc, keyList("materi:"+idKelas))
	resp := mapMateriToResponse(materi)
	return &resp, nil
}

func (s *MateriService) Update(id string, req dto.UpdateMateriRequest) (*dto.MateriResponse, error) {
	materi, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("materi tidak ditemukan")
	}

	if req.Judul != nil {
		trimmed := strings.TrimSpace(*req.Judul)
		if trimmed == "" {
			return nil, errors.New("judul tidak boleh kosong")
		}
		materi.Judul = trimmed
	}
	if req.Urutan != nil {
		materi.Urutan = *req.Urutan
	}
	if req.YoutubeID != nil {
		materi.YoutubeID = req.YoutubeID
	}
	if req.PDFPath != nil {
		materi.PDFPath = req.PDFPath
	}
	if req.DurasiDetik != nil {
		materi.DurasiDetik = req.DurasiDetik
	}

	if err := s.repo.Update(materi); err != nil {
		return nil, err
	}

	cacheDelete(s.rc, keyList("materi:"+materi.IDKelas))
	resp := mapMateriToResponse(materi)
	return &resp, nil
}

func (s *MateriService) Delete(id string) error {
	materi, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("materi tidak ditemukan")
	}

	idKelas := materi.IDKelas
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// Rapikan urutan materi yang tersisa
	_ = s.repo.ReorderUrutan(idKelas)

	cacheDelete(s.rc, keyList("materi:"+idKelas))
	return nil
}

// ── User: Update Progress ─────────────────────────────────────────────────────

// UpdateProgress dipakai untuk:
//   - Video: update last_watched_seconds, tandai selesai jika >= 80% durasi
//   - PDF  : langsung tandai selesai (is_completed=true)
func (s *MateriService) UpdateProgress(userID, materiID string, req dto.UpdateProgressRequest) (*dto.ProgressResponse, error) {
	materi, err := s.repo.FindByID(materiID)
	if err != nil {
		return nil, errors.New("materi tidak ditemukan")
	}

	if materi.Tipe == models.MateriTipeKuis {
		return nil, errors.New("progress kuis dikelola melalui endpoint kuis")
	}

	// Ambil atau buat progress record
	progress, err := s.progressRepo.FindByUserAndMateri(userID, materiID)
	if err != nil {
		// Belum ada, buat baru
		progress = &models.UserMateriProgress{
			ID:       uuid.New().String(),
			IDUser:   userID,
			IDMateri: materiID,
		}
	}

	// Jangan mundurkan progress yang sudah selesai
	if progress.IsCompleted {
		return toProgressResponse(progress), nil
	}

	switch materi.Tipe {
	case models.MateriTipeVideo:
		if req.LastWatchedSeconds != nil {
			progress.LastWatchedSeconds = *req.LastWatchedSeconds
		}

		// Tandai selesai jika sudah nonton >= 80% durasi
		if materi.DurasiDetik != nil && *materi.DurasiDetik > 0 {
			pct := float64(progress.LastWatchedSeconds) / float64(*materi.DurasiDetik)
			if pct >= 0.8 {
				req.IsCompleted = true
			}
		}
		// Admin bisa force complete lewat is_completed=true
		if req.IsCompleted {
			progress.IsCompleted = true
			now := time.Now()
			progress.CompletedAt = &now
		}

	case models.MateriTipePDF:
		// PDF: selesai saat dibuka / dikonfirmasi client
		if req.IsCompleted {
			progress.IsCompleted = true
			now := time.Now()
			progress.CompletedAt = &now
		}
	}

	if err := s.progressRepo.Upsert(progress); err != nil {
		return nil, err
	}

	return toProgressResponse(progress), nil
}

func toProgressResponse(p *models.UserMateriProgress) *dto.ProgressResponse {
	resp := &dto.ProgressResponse{
		IDMateri:           p.IDMateri,
		IsCompleted:        p.IsCompleted,
		LastWatchedSeconds: p.LastWatchedSeconds,
	}
	if p.CompletedAt != nil {
		s := p.CompletedAt.Format(time.RFC3339)
		resp.CompletedAt = &s
	}
	return resp
}