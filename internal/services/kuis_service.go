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

type KuisService struct {
	repo         repository.KuisAttemptRepositoryInterface
	soalRepo     repository.SoalRepositoryInterface
	materiRepo   repository.MateriRepositoryInterface
	progressRepo repository.ProgressRepositoryInterface
	rc           cache.RedisInterface
}

func NewKuisService(
	repo repository.KuisAttemptRepositoryInterface,
	soalRepo repository.SoalRepositoryInterface,
	materiRepo repository.MateriRepositoryInterface,
	progressRepo repository.ProgressRepositoryInterface,
	rc cache.RedisInterface,
) *KuisService {
	return &KuisService{
		repo:         repo,
		soalRepo:     soalRepo,
		materiRepo:   materiRepo,
		progressRepo: progressRepo,
		rc:           rc,
	}
}

// ── Start Kuis ────────────────────────────────────────────────────────────────

// Start memvalidasi prerequisite lalu membuat attempt baru.
// User harus sudah menyelesaikan minimal 1 video ATAU 1 pdf dalam kelas yang sama.
func (s *KuisService) Start(userID, materiID string) (*dto.StartKuisResponse, error) {
	// 1. Pastikan materi ada dan bertipe kuis
	materi, err := s.materiRepo.FindByID(materiID)
	if err != nil {
		return nil, errors.New("materi tidak ditemukan")
	}
	if materi.Tipe != models.MateriTipeKuis {
		return nil, errors.New("materi ini bukan bertipe kuis")
	}

	// 2. Cek prerequisite: user harus sudah selesai minimal 1 video/pdf dalam kelas
	hasCompleted, err := s.progressRepo.HasCompletedAnyMedia(userID, materi.IDKelas)
	if err != nil {
		return nil, err
	}
	if !hasCompleted {
		return nil, errors.New("selesaikan setidaknya satu materi (video atau pdf) sebelum mengerjakan kuis")
	}

	// 3. Cek apakah ada attempt yang belum selesai (belum di-submit)
	latest, err := s.repo.FindLatestByUserAndMateri(userID, materiID)
	if err == nil && latest != nil && latest.FinishedAt == nil {
		// Kembalikan attempt yang sudah ada + soal-soalnya
		return s.buildStartResponse(latest.ID, materiID)
	}

	// 4. Buat attempt baru
	attempt := &models.KuisAttempt{
		ID:       uuid.New().String(),
		IDUser:   userID,
		IDMateri: materiID,
	}
	if err := s.repo.Create(attempt); err != nil {
		return nil, err
	}

	return s.buildStartResponse(attempt.ID, materiID)
}

func (s *KuisService) buildStartResponse(attemptID, materiID string) (*dto.StartKuisResponse, error) {
	soalList, err := s.soalRepo.FindByMateri(materiID)
	if err != nil {
		return nil, err
	}
	if len(soalList) == 0 {
		return nil, errors.New("kuis ini belum memiliki soal")
	}

	// Map ke SoalUserResponse — is_correct disembunyikan
	soalResp := make([]dto.SoalUserResponse, 0, len(soalList))
	for _, soal := range soalList {
		pilihan := make([]dto.PilihanUserResponse, 0, len(soal.Pilihan))
		for _, p := range soal.Pilihan {
			pilihan = append(pilihan, dto.PilihanUserResponse{
				ID:     p.ID,
				Teks:   p.Teks,
				Urutan: p.Urutan,
			})
		}
		soalResp = append(soalResp, dto.SoalUserResponse{
			ID:         soal.ID,
			Pertanyaan: soal.Pertanyaan,
			Urutan:     soal.Urutan,
			Pilihan:    pilihan,
		})
	}

	return &dto.StartKuisResponse{
		AttemptID: attemptID,
		IDMateri:  materiID,
		Soal:      soalResp,
	}, nil
}

// ── Submit Kuis ───────────────────────────────────────────────────────────────

// Submit menerima jawaban user, menghitung skor, dan menyimpan hasilnya.
func (s *KuisService) Submit(userID, attemptID string, req dto.SubmitKuisRequest) (*dto.KuisResultResponse, error) {
	// 1. Validasi attempt milik user dan belum selesai
	attempt, err := s.repo.FindByID(attemptID)
	if err != nil {
		return nil, errors.New("attempt tidak ditemukan")
	}
	if attempt.IDUser != userID {
		return nil, errors.New("attempt bukan milik user ini")
	}
	if attempt.FinishedAt != nil {
		return nil, errors.New("kuis ini sudah pernah dikerjakan")
	}

	// 2. Ambil semua soal kuis ini
	soalList, err := s.soalRepo.FindByMateri(attempt.IDMateri)
	if err != nil {
		return nil, err
	}

	// Buat map idSoal -> soal untuk lookup cepat
	soalMap := make(map[string]*models.Soal, len(soalList))
	for i := range soalList {
		soalMap[soalList[i].ID] = &soalList[i]
	}

	// 3. Validasi: semua soal harus dijawab, tidak boleh lebih
	if len(req.Jawaban) != len(soalList) {
		return nil, errors.New("jumlah jawaban tidak sesuai dengan jumlah soal")
	}

	// Pastikan tidak ada soal yang dijawab dua kali
	answeredSoal := make(map[string]bool, len(req.Jawaban))
	for _, j := range req.Jawaban {
		if answeredSoal[j.IDSoal] {
			return nil, errors.New("terdapat soal yang dijawab lebih dari satu kali")
		}
		answeredSoal[j.IDSoal] = true
	}

	// 4. Hitung skor
	jawabanModels := make([]models.KuisJawaban, 0, len(req.Jawaban))
	totalBenar := 0
	detailHasil := make([]dto.HasilSoalResponse, 0, len(soalList))

	for _, j := range req.Jawaban {
		soal, ok := soalMap[j.IDSoal]
		if !ok {
			return nil, errors.New("soal dengan id " + j.IDSoal + " tidak ditemukan dalam kuis ini")
		}

		// Ambil pilihan yang dipilih user
		pilihanUser, err := s.soalRepo.FindPilihanByID(j.IDPilihan)
		if err != nil {
			return nil, errors.New("pilihan jawaban tidak ditemukan")
		}
		// Pastikan pilihan ini memang milik soal yang bersangkutan
		if pilihanUser.IDSoal != j.IDSoal {
			return nil, errors.New("pilihan jawaban tidak sesuai dengan soal")
		}

		// Ambil pilihan benar untuk feedback
		pilihanBenar, err := s.soalRepo.FindCorrectPilihan(j.IDSoal)
		if err != nil {
			return nil, err
		}

		isCorrect := pilihanUser.IsCorrect
		if isCorrect {
			totalBenar++
		}

		jawabanModels = append(jawabanModels, models.KuisJawaban{
			ID:        uuid.New().String(),
			IDAttempt: attemptID,
			IDSoal:    j.IDSoal,
			IDPilihan: j.IDPilihan,
			IsCorrect: isCorrect,
		})

		detailHasil = append(detailHasil, dto.HasilSoalResponse{
			IDSoal:         soal.ID,
			Pertanyaan:     soal.Pertanyaan,
			IDPilihanUser:  j.IDPilihan,
			IDPilihanBenar: pilihanBenar.ID,
			IsCorrect:      isCorrect,
		})
	}

	skor := float64(totalBenar) / float64(len(soalList)) * 100

	// 5. Simpan hasil
	if err := s.repo.Finish(attemptID, skor, totalBenar, jawabanModels); err != nil {
		return nil, err
	}

	now := time.Now()
	return &dto.KuisResultResponse{
		AttemptID:  attemptID,
		Skor:       skor,
		TotalSoal:  len(soalList),
		TotalBenar: totalBenar,
		FinishedAt: now.Format(time.RFC3339),
		Detail:     detailHasil,
	}, nil
}

// ── Get Result ────────────────────────────────────────────────────────────────

// GetResult mengembalikan hasil attempt yang sudah selesai.
func (s *KuisService) GetResult(userID, attemptID string) (*dto.KuisResultResponse, error) {
	attempt, err := s.repo.FindByID(attemptID)
	if err != nil {
		return nil, errors.New("attempt tidak ditemukan")
	}
	if attempt.IDUser != userID {
		return nil, errors.New("attempt bukan milik user ini")
	}
	if attempt.FinishedAt == nil {
		return nil, errors.New("kuis belum selesai dikerjakan")
	}

	// Ambil jawaban user
	jawabanList, err := s.repo.FindJawabanByAttempt(attemptID)
	if err != nil {
		return nil, err
	}

	// Ambil soal untuk tampilkan pertanyaan dan pilihan benar
	soalList, err := s.soalRepo.FindByMateri(attempt.IDMateri)
	if err != nil {
		return nil, err
	}
	soalMap := make(map[string]*models.Soal, len(soalList))
	for i := range soalList {
		soalMap[soalList[i].ID] = &soalList[i]
	}

	detail := make([]dto.HasilSoalResponse, 0, len(jawabanList))
	for _, j := range jawabanList {
		soal := soalMap[j.IDSoal]
		pilihanBenar, err := s.soalRepo.FindCorrectPilihan(j.IDSoal)
		if err != nil {
			continue
		}

		pertanyaan := ""
		if soal != nil {
			pertanyaan = soal.Pertanyaan
		}

		detail = append(detail, dto.HasilSoalResponse{
			IDSoal:         j.IDSoal,
			Pertanyaan:     pertanyaan,
			IDPilihanUser:  j.IDPilihan,
			IDPilihanBenar: pilihanBenar.ID,
			IsCorrect:      j.IsCorrect,
		})
	}

	return &dto.KuisResultResponse{
		AttemptID:  attempt.ID,
		Skor:       attempt.Skor,
		TotalSoal:  attempt.TotalSoal,
		TotalBenar: attempt.TotalBenar,
		FinishedAt: attempt.FinishedAt.Format(time.RFC3339),
		Detail:     detail,
	}, nil
}

// ── Riwayat Attempt ───────────────────────────────────────────────────────────

// GetAttemptsByUser mengembalikan semua attempt user untuk satu kuis.
// Berguna jika kuis diizinkan untuk diulang.
func (s *KuisService) GetAttemptsByUser(userID, materiID string) ([]models.KuisAttempt, error) {
	materi, err := s.materiRepo.FindByID(materiID)
	if err != nil {
		return nil, errors.New("materi tidak ditemukan")
	}
	if materi.Tipe != models.MateriTipeKuis {
		return nil, errors.New("materi ini bukan bertipe kuis")
	}
	return s.repo.FindByUserAndMateri(userID, materiID)
}