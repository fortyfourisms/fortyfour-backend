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
	return &RisikoService{r}
}

func (s *RisikoService) Create(req dto.CreateRisikoRequest) (dto.RisikoResponse, error) {

	if req.NamaRisiko == "" {
		return dto.RisikoResponse{}, errors.New("nama risiko wajib diisi")
	}

	data, err := s.repo.Create(repository.Risiko(req))
	if err != nil {
		return dto.RisikoResponse{}, err
	}

	return dto.RisikoResponse(data), nil
}

func (s *RisikoService) GetAll() ([]dto.RisikoResponse, error) {

	list, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var res []dto.RisikoResponse
	for _, v := range list {
		res = append(res, dto.RisikoResponse(v))
	}

	return res, nil
}