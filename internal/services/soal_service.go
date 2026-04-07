package services

import (
	"errors"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/pkg/cache"

	"github.com/google/uuid"
)

type SoalService struct {
	repo       repository.SoalRepositoryInterface
	materiRepo repository.MateriRepositoryInterface
	rc         cache.RedisInterface
}

func NewSoalService(
	repo repository.SoalRepositoryInterface,
	materiRepo repository.MateriRepositoryInterface,
	rc cache.RedisInterface,
) *SoalService {
	return &SoalService{repo: repo, materiRepo: materiRepo, rc: rc}
}

// ── Admin: CRUD Soal ──────────────────────────────────────────────────────────

func (s *SoalService) Create(idMateri string, req dto.CreateSoalRequest) (*dto.SoalResponse, error) {
	materi, err := s.materiRepo.FindByID(idMateri)
	if err != nil {
		return nil, errors.New("materi tidak ditemukan")
	}
	if materi.Tipe != models.MateriTipeKuis {
		return nil, errors.New("materi ini bukan bertipe kuis")
	}

	pertanyaan := strings.TrimSpace(req.Pertanyaan)
	if pertanyaan == "" {
		return nil, errors.New("pertanyaan wajib diisi")
	}
	if err := validatePilihan(req.Pilihan); err != nil {
		return nil, err
	}

	soal := &models.Soal{
		ID:         uuid.New().String(),
		IDMateri:   idMateri,
		Pertanyaan: pertanyaan,
		Urutan:     req.Urutan,
	}

	pilihan := buildPilihan(soal.ID, req.Pilihan)

	if err := s.repo.Create(soal, pilihan); err != nil {
		return nil, err
	}

	soal.Pilihan = pilihan
	return toSoalResponse(soal), nil
}

func (s *SoalService) Update(id string, req dto.UpdateSoalRequest) (*dto.SoalResponse, error) {
	soal, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("soal tidak ditemukan")
	}

	if req.Pertanyaan != nil {
		trimmed := strings.TrimSpace(*req.Pertanyaan)
		if trimmed == "" {
			return nil, errors.New("pertanyaan tidak boleh kosong")
		}
		soal.Pertanyaan = trimmed
	}
	if req.Urutan != nil {
		soal.Urutan = *req.Urutan
	}

	var pilihan []models.PilihanJawaban
	if len(req.Pilihan) > 0 {
		if err := validatePilihan(req.Pilihan); err != nil {
			return nil, err
		}
		pilihan = buildPilihan(soal.ID, req.Pilihan)
	}

	if err := s.repo.Update(soal, pilihan); err != nil {
		return nil, err
	}

	soal.Pilihan = pilihan
	return toSoalResponse(soal), nil
}

func (s *SoalService) Delete(id string) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return errors.New("soal tidak ditemukan")
	}
	return s.repo.Delete(id)
}

// GetByMateri untuk admin (tampilkan is_correct)
func (s *SoalService) GetByMateri(idMateri string) ([]dto.SoalResponse, error) {
	soalList, err := s.repo.FindByMateri(idMateri)
	if err != nil {
		return nil, err
	}

	result := make([]dto.SoalResponse, 0, len(soalList))
	for _, soal := range soalList {
		soal := soal
		result = append(result, *toSoalResponse(&soal))
	}
	return result, nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func validatePilihan(pilihan []dto.CreatePilihanRequest) error {
	if len(pilihan) < 2 {
		return errors.New("soal harus memiliki minimal 2 pilihan jawaban")
	}
	if len(pilihan) > 5 {
		return errors.New("soal maksimal memiliki 5 pilihan jawaban")
	}

	correctCount := 0
	for _, p := range pilihan {
		if strings.TrimSpace(p.Teks) == "" {
			return errors.New("teks pilihan jawaban tidak boleh kosong")
		}
		if p.IsCorrect {
			correctCount++
		}
	}
	if correctCount == 0 {
		return errors.New("harus ada tepat 1 pilihan jawaban yang benar")
	}
	if correctCount > 1 {
		return errors.New("hanya boleh ada 1 pilihan jawaban yang benar")
	}
	return nil
}

func buildPilihan(idSoal string, reqs []dto.CreatePilihanRequest) []models.PilihanJawaban {
	pilihan := make([]models.PilihanJawaban, 0, len(reqs))
	for _, p := range reqs {
		pilihan = append(pilihan, models.PilihanJawaban{
			ID:        uuid.New().String(),
			IDSoal:    idSoal,
			Teks:      strings.TrimSpace(p.Teks),
			IsCorrect: p.IsCorrect,
			Urutan:    p.Urutan,
		})
	}
	return pilihan
}

func toSoalResponse(soal *models.Soal) *dto.SoalResponse {
	pilihan := make([]dto.PilihanResponse, 0, len(soal.Pilihan))
	for _, p := range soal.Pilihan {
		pilihan = append(pilihan, dto.PilihanResponse{
			ID:        p.ID,
			Teks:      p.Teks,
			IsCorrect: p.IsCorrect,
			Urutan:    p.Urutan,
		})
	}
	return &dto.SoalResponse{
		ID:         soal.ID,
		IDMateri:   soal.IDMateri,
		Pertanyaan: soal.Pertanyaan,
		Urutan:     soal.Urutan,
		Pilihan:    pilihan,
	}
}
