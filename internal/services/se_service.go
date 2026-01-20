package services

import (
	"errors"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"

	"github.com/google/uuid"
)

type SEService interface {
	Create(req dto.CreateSERequest) (*dto.SEResponse, error)
	GetAll() ([]dto.SEResponse, error)
	GetByID(id string) (*dto.SEResponse, error)
	Update(id string, req dto.UpdateSERequest) (*dto.SEResponse, error)
	Delete(id string) error
}

type seService struct {
	repo repository.SERepositoryInterface
}

func NewSEService(repo repository.SERepositoryInterface) SEService {
	return &seService{repo: repo}
}

/* =======================
   CREATE
======================= */

func (s *seService) Create(req dto.CreateSERequest) (*dto.SEResponse, error) {
	if req.IDPerusahaan == nil || strings.TrimSpace(*req.IDPerusahaan) == "" {
		return nil, errors.New("id_perusahaan wajib diisi")
	}
	if req.IDSubSektor == nil || strings.TrimSpace(*req.IDSubSektor) == "" {
		return nil, errors.New("id_sub_sektor wajib diisi")
	}

	totalBobot, err := hitungTotalBobotCreate(req)
	if err != nil {
		return nil, err
	}

	kategori := hitungKategoriSE(totalBobot)
	id := uuid.NewString()

	if err := s.repo.Create(req, id, totalBobot, kategori); err != nil {
		return nil, err
	}

	return s.repo.GetByID(id)
}

/* =======================
   READ
======================= */

func (s *seService) GetAll() ([]dto.SEResponse, error) {
	return s.repo.GetAll()
}

func (s *seService) GetByID(id string) (*dto.SEResponse, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("id wajib diisi")
	}
	return s.repo.GetByID(id)
}

/* =======================
   UPDATE
======================= */

func (s *seService) Update(id string, req dto.UpdateSERequest) (*dto.SEResponse, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("id wajib diisi")
	}

	totalBobot, err := hitungTotalBobotUpdate(req)
	if err != nil {
		return nil, err
	}

	kategori := hitungKategoriSE(totalBobot)

	if err := s.repo.Update(id, req, totalBobot, kategori); err != nil {
		return nil, err
	}

	return s.repo.GetByID(id)
}

/* =======================
   DELETE
======================= */

func (s *seService) Delete(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id wajib diisi")
	}
	return s.repo.Delete(id)
}

/* ======================================================
   HELPER FUNCTIONS (A/B/C → BOBOT)
====================================================== */

func jawabanKeBobot(jawaban string) (int, error) {
	switch strings.ToUpper(strings.TrimSpace(jawaban)) {
	case "A":
		return 5, nil
	case "B":
		return 2, nil
	case "C":
		return 1, nil
	default:
		return 0, errors.New("jawaban harus A, B, atau C")
	}
}

func hitungTotalBobotCreate(req dto.CreateSERequest) (int, error) {
	qs := []string{
		req.Q1, req.Q2, req.Q3, req.Q4, req.Q5,
		req.Q6, req.Q7, req.Q8, req.Q9, req.Q10,
	}

	total := 0
	for _, q := range qs {
		bobot, err := jawabanKeBobot(q)
		if err != nil {
			return 0, err
		}
		total += bobot
	}
	return total, nil
}

func hitungTotalBobotUpdate(req dto.UpdateSERequest) (int, error) {
	qs := []string{
		req.Q1, req.Q2, req.Q3, req.Q4, req.Q5,
		req.Q6, req.Q7, req.Q8, req.Q9, req.Q10,
	}

	total := 0
	for _, q := range qs {
		bobot, err := jawabanKeBobot(q)
		if err != nil {
			return 0, err
		}
		total += bobot
	}
	return total, nil
}

/* =======================
   KATEGORISASI SE
======================= */

func hitungKategoriSE(totalBobot int) string {
	switch {
	case totalBobot >= 35 && totalBobot <= 50:
		return "Strategis"
	case totalBobot >= 16 && totalBobot <= 34:
		return "Tinggi"
	case totalBobot >= 10 && totalBobot <= 15:
		return "Rendah"
	default:
		return ""
	}
}
