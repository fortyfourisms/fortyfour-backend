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
	attemptRepo  repository.KuisAttemptRepositoryInterface
	soalRepo     repository.SoalRepositoryInterface
	kuisRepo     repository.KuisRepositoryInterface
	progressRepo repository.ProgressRepositoryInterface
	rc           cache.RedisInterface
}

func NewKuisService(
	attemptRepo repository.KuisAttemptRepositoryInterface,
	soalRepo repository.SoalRepositoryInterface,
	kuisRepo repository.KuisRepositoryInterface,
	progressRepo repository.ProgressRepositoryInterface,
	rc cache.RedisInterface,
) *KuisService {
	return &KuisService{
		attemptRepo:  attemptRepo,
		soalRepo:     soalRepo,
		kuisRepo:     kuisRepo,
		progressRepo: progressRepo,
		rc:           rc,
	}
}

// ── Admin: CRUD Kuis ──────────────────────────────────────────────────────────

func (s *KuisService) CreateKuis(idKelas string, req dto.CreateKuisRequest) (*dto.KuisResponse, error) {
	kuis := &models.Kuis{
		ID:           uuid.New().String(),
		IDKelas:      idKelas,
		IDMateri:     req.IDMateri,
		Judul:        req.Judul,
		Deskripsi:    req.Deskripsi,
		DurasiMenit:  req.DurasiMenit,
		PassingGrade: req.PassingGrade,
		IsFinal:      req.IsFinal,
		Urutan:       req.Urutan,
	}

	// Validasi: hanya boleh ada 1 kuis final per kelas
	if kuis.IsFinal {
		existing, err := s.kuisRepo.FindFinalByKelas(idKelas)
		if err == nil && existing != nil {
			return nil, errors.New("sudah ada kuis akhir untuk kelas ini")
		}
		kuis.IDMateri = nil // kuis final tidak terikat ke materi
	}

	if err := s.kuisRepo.Create(kuis); err != nil {
		return nil, err
	}

	return mapKuisToResponse(kuis), nil
}

func (s *KuisService) UpdateKuis(id string, req dto.UpdateKuisRequest) (*dto.KuisResponse, error) {
	kuis, err := s.kuisRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("kuis tidak ditemukan")
	}

	if req.Judul != nil {
		kuis.Judul = *req.Judul
	}
	if req.Deskripsi != nil {
		kuis.Deskripsi = req.Deskripsi
	}
	if req.DurasiMenit != nil {
		kuis.DurasiMenit = req.DurasiMenit
	}
	if req.PassingGrade != nil {
		kuis.PassingGrade = *req.PassingGrade
	}
	if req.IsFinal != nil {
		kuis.IsFinal = *req.IsFinal
	}
	if req.Urutan != nil {
		kuis.Urutan = *req.Urutan
	}

	if err := s.kuisRepo.Update(kuis); err != nil {
		return nil, err
	}

	return mapKuisToResponse(kuis), nil
}

func (s *KuisService) DeleteKuis(id string) error {
	if _, err := s.kuisRepo.FindByID(id); err != nil {
		return errors.New("kuis tidak ditemukan")
	}
	return s.kuisRepo.Delete(id)
}

func (s *KuisService) GetKuisByKelas(idKelas string) ([]dto.KuisResponse, error) {
	kuisList, err := s.kuisRepo.FindByKelas(idKelas)
	if err != nil {
		return nil, err
	}

	result := make([]dto.KuisResponse, 0, len(kuisList))
	for _, k := range kuisList {
		k := k
		result = append(result, *mapKuisToResponse(&k))
	}
	return result, nil
}

// ── User: Start Kuis ──────────────────────────────────────────────────────────

// Start memvalidasi prerequisite lalu membuat attempt baru.
func (s *KuisService) Start(userID, kuisID string) (*dto.StartKuisResponse, error) {
	// 1. Pastikan kuis ada
	kuis, err := s.kuisRepo.FindByID(kuisID)
	if err != nil {
		return nil, errors.New("kuis tidak ditemukan")
	}

	// 2. Cek prerequisite untuk kuis akhir
	if kuis.IsFinal {
		// Semua materi harus selesai
		allDone, err := s.progressRepo.HasCompletedAllMateri(userID, kuis.IDKelas)
		if err != nil {
			return nil, err
		}
		if !allDone {
			return nil, errors.New("selesaikan semua materi sebelum mengerjakan kuis akhir")
		}
		// Semua kuis non-final harus lulus
		allPassed, err := s.attemptRepo.HasPassedAllKuisInKelas(userID, kuis.IDKelas)
		if err != nil {
			return nil, err
		}
		if !allPassed {
			return nil, errors.New("lulus semua kuis per-materi sebelum mengerjakan kuis akhir")
		}
	}

	// 3. Cek apakah ada attempt yang belum selesai
	latest, err := s.attemptRepo.FindLatestByUserAndKuis(userID, kuisID)
	if err == nil && latest != nil && latest.FinishedAt == nil {
		return s.buildStartResponse(latest.ID, kuisID)
	}

	// 4. Buat attempt baru
	attempt := &models.KuisAttempt{
		ID:     uuid.New().String(),
		IDUser: userID,
		IDKuis: kuisID,
	}
	if err := s.attemptRepo.Create(attempt); err != nil {
		return nil, err
	}

	return s.buildStartResponse(attempt.ID, kuisID)
}

func (s *KuisService) buildStartResponse(attemptID, kuisID string) (*dto.StartKuisResponse, error) {
	soalList, err := s.soalRepo.FindByKuis(kuisID)
	if err != nil {
		return nil, err
	}
	if len(soalList) == 0 {
		return nil, errors.New("kuis ini belum memiliki soal")
	}

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
		IDKuis:    kuisID,
		Soal:      soalResp,
	}, nil
}

// ── User: Submit Kuis ─────────────────────────────────────────────────────────

func (s *KuisService) Submit(userID, attemptID string, req dto.SubmitKuisRequest) (*dto.KuisResultResponse, error) {
	// 1. Validasi attempt
	attempt, err := s.attemptRepo.FindByID(attemptID)
	if err != nil {
		return nil, errors.New("attempt tidak ditemukan")
	}
	if attempt.IDUser != userID {
		return nil, errors.New("attempt bukan milik user ini")
	}
	if attempt.FinishedAt != nil {
		return nil, errors.New("kuis ini sudah pernah dikerjakan")
	}

	// 2. Ambil kuis & soal
	kuis, err := s.kuisRepo.FindByID(attempt.IDKuis)
	if err != nil {
		return nil, errors.New("kuis tidak ditemukan")
	}

	soalList, err := s.soalRepo.FindByKuis(attempt.IDKuis)
	if err != nil {
		return nil, err
	}

	soalMap := make(map[string]*models.Soal, len(soalList))
	for i := range soalList {
		soalMap[soalList[i].ID] = &soalList[i]
	}

	// 3. Validasi jawaban
	if len(req.Jawaban) != len(soalList) {
		return nil, errors.New("jumlah jawaban tidak sesuai dengan jumlah soal")
	}

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

		pilihanUser, err := s.soalRepo.FindPilihanByID(j.IDPilihan)
		if err != nil {
			return nil, errors.New("pilihan jawaban tidak ditemukan")
		}
		if pilihanUser.IDSoal != j.IDSoal {
			return nil, errors.New("pilihan jawaban tidak sesuai dengan soal")
		}

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
	isPassed := skor >= kuis.PassingGrade

	// 5. Simpan hasil
	if err := s.attemptRepo.Finish(attemptID, skor, totalBenar, isPassed, jawabanModels); err != nil {
		return nil, err
	}

	now := time.Now()
	return &dto.KuisResultResponse{
		AttemptID:  attemptID,
		Skor:       skor,
		TotalSoal:  len(soalList),
		TotalBenar: totalBenar,
		IsPassed:   isPassed,
		FinishedAt: now.Format(time.RFC3339),
		Detail:     detailHasil,
	}, nil
}

// ── Get Result ────────────────────────────────────────────────────────────────

func (s *KuisService) GetResult(userID, attemptID string) (*dto.KuisResultResponse, error) {
	attempt, err := s.attemptRepo.FindByID(attemptID)
	if err != nil {
		return nil, errors.New("attempt tidak ditemukan")
	}
	if attempt.IDUser != userID {
		return nil, errors.New("attempt bukan milik user ini")
	}
	if attempt.FinishedAt == nil {
		return nil, errors.New("kuis belum selesai dikerjakan")
	}

	jawabanList, err := s.attemptRepo.FindJawabanByAttempt(attemptID)
	if err != nil {
		return nil, err
	}

	soalList, err := s.soalRepo.FindByKuis(attempt.IDKuis)
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
		IsPassed:   attempt.IsPassed,
		FinishedAt: attempt.FinishedAt.Format(time.RFC3339),
		Detail:     detail,
	}, nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func mapKuisToResponse(k *models.Kuis) *dto.KuisResponse {
	return &dto.KuisResponse{
		ID:           k.ID,
		IDKelas:      k.IDKelas,
		IDMateri:     k.IDMateri,
		Judul:        k.Judul,
		Deskripsi:    k.Deskripsi,
		DurasiMenit:  k.DurasiMenit,
		PassingGrade: k.PassingGrade,
		IsFinal:      k.IsFinal,
		Urutan:       k.Urutan,
		CreatedAt:    k.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    k.UpdatedAt.Format(time.RFC3339),
	}
}
