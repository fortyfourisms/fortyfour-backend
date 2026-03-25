package services

import (
	"survey/internal/dto"
	"survey/internal/models"
	"survey/internal/repository"
)

type RisikoService struct {
	repo *repository.RisikoRepository
}

func NewRisikoService(r *repository.RisikoRepository) *RisikoService {
	return &RisikoService{repo: r}
}

func (s *RisikoService) Create(req dto.CreateRisikoRequest) error {

	data := models.RisikoSurvey{
		RespondenID: req.RespondenID,
		RisikoIP: req.RisikoIP,
		DampakReputasi: req.DampakReputasi,
		DampakOperasional: req.DampakOperasional,
		DampakFinansial: req.DampakFinansial,
		DampakHukum: req.DampakHukum,
		Frekuensi: req.Frekuensi,
		AdaPengendalian: req.AdaPengendalian,
		TindakanPengendalian: req.TindakanPengendalian,
	}

	return s.repo.Create(data)
}