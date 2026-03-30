package services

import (
	"errors"
	"survey/internal/dto"
	"survey/internal/repository"
)

type RisikoService struct {
	repo *repository.RisikoRepository
}

func NewRisikoService(r *repository.RisikoRepository) *RisikoService {
	return &RisikoService{repo: r}
}

// VALIDASI
func (s *RisikoService) validate(req dto.CreateRisikoJawabanRequest) error {

	if req.RespondenID == 0 {
		return errors.New("responden_id wajib")
	}

	if req.RisikoID == 0 {
		return errors.New("risiko_id wajib")
	}

	if req.PernahTerjadi == "" {
		return errors.New("pernah_terjadi wajib")
	}

	return nil
}

// CREATE JAWABAN
func (s *RisikoService) CreateJawaban(req dto.CreateRisikoJawabanRequest) error {

	if err := s.validate(req); err != nil {
		return err
	}

	data := map[string]interface{}{
		"responden_id": req.RespondenID,
		"risiko_id": req.RisikoID,
		"pernah_terjadi": req.PernahTerjadi,
		"dampak_reputasi": req.DampakReputasi,
		"dampak_operasional": req.DampakOperasional,
		"dampak_finansial": req.DampakFinansial,
		"dampak_hukum": req.DampakHukum,
		"frekuensi": req.Frekuensi,
		"ada_pengendalian": req.AdaPengendalian,
		"deskripsi_pengendalian": req.DeskripsiPengendalian,
	}

	return s.repo.CreateJawaban(data)
}