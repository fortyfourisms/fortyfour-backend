package services

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"survey/internal/dto"
	"survey/internal/models"
	"survey/internal/repository"
)

type RespondenService struct {
	repo *repository.RespondenRepository
}

func NewRespondenService(repo *repository.RespondenRepository) *RespondenService {
	return &RespondenService{repo: repo}
}

//Validation
func (s *RespondenService) validate(req *dto.CreateRespondenRequest) error {

	req.NamaLengkap = strings.TrimSpace(req.NamaLengkap)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if req.NamaLengkap == "" {
		return errors.New("nama_lengkap tidak boleh kosong")
	}

	if len(req.NamaLengkap) < 3 {
		return errors.New("nama_lengkap minimal 3 karakter")
	}

	if req.Email == "" {
		return errors.New("email tidak boleh kosong")
	}

	if !strings.Contains(req.Email, "@") {
		return errors.New("format email tidak valid")
	}

	return nil
}

//Create
func (s *RespondenService) Create(req dto.CreateRespondenRequest) (*dto.RespondenResponse, error) {

	if err := s.validate(&req); err != nil {
		return nil, err
	}

	model := s.toModel(req)

	if err := s.repo.Create(model); err != nil {
		return nil, err
	}

	// Ambil data terakhir
	all, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	if len(all) == 0 {
		return nil, errors.New("gagal mengambil data setelah create")
	}

	created := all[len(all)-1]
	resp := s.toResponse(&created)

	return &resp, nil
}

//Get
func (s *RespondenService) GetAll() ([]dto.RespondenResponse, error) {

	data, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var result []dto.RespondenResponse
	for i := range data {
		result = append(result, s.toResponse(&data[i]))
	}

	return result, nil
}

func (s *RespondenService) GetByID(id int) (*dto.RespondenResponse, error) {

	if id <= 0 {
		return nil, errors.New("id tidak valid")
	}

	data, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}

	resp := s.toResponse(data)
	return &resp, nil
}

//Update
func (s *RespondenService) Update(id int, req dto.UpdateRespondenRequest) (*dto.RespondenResponse, error) {

	if id <= 0 {
		return nil, errors.New("id tidak valid")
	}

	_, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}

	//Convert update DTO ke model
	model := models.Responden{
		NamaLengkap:        strings.TrimSpace(req.NamaLengkap),
		Jabatan:            strings.TrimSpace(req.Jabatan),
		Perusahaan:         strings.TrimSpace(req.Perusahaan),
		Email:              strings.ToLower(strings.TrimSpace(req.Email)),
		NoTelepon:          strings.TrimSpace(req.NoTelepon),
		Sektor:             strings.TrimSpace(req.Sektor),
		SektorLainnya:      strings.TrimSpace(req.SektorLainnya),
		SertifikatTraining: strings.TrimSpace(req.SertifikatTraining),
		UpdatedAt:          time.Now(),
	}

	if err := s.repo.Update(id, model); err != nil {
		return nil, err
	}

	updated, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	resp := s.toResponse(updated)
	return &resp, nil
}

//Delete
func (s *RespondenService) Delete(id int) error {

	if id <= 0 {
		return errors.New("id tidak valid")
	}

	_, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("data tidak ditemukan")
		}
		return err
	}

	return s.repo.Delete(id)
}

//Helpers
func (s *RespondenService) toModel(req dto.CreateRespondenRequest) models.Responden {

	return models.Responden{
		NamaLengkap:        strings.TrimSpace(req.NamaLengkap),
		Jabatan:            strings.TrimSpace(req.Jabatan),
		Perusahaan:         strings.TrimSpace(req.Perusahaan),
		Email:              strings.ToLower(strings.TrimSpace(req.Email)),
		NoTelepon:          strings.TrimSpace(req.NoTelepon),
		Sektor:             strings.TrimSpace(req.Sektor),
		SektorLainnya:      strings.TrimSpace(req.SektorLainnya),
		SertifikatTraining: strings.TrimSpace(req.SertifikatTraining),
	}
}

func (s *RespondenService) toResponse(m *models.Responden) dto.RespondenResponse {

	return dto.RespondenResponse{
		ID:                 m.ID,
		NamaLengkap:        m.NamaLengkap,
		Jabatan:            m.Jabatan,
		Perusahaan:         m.Perusahaan,
		Email:              m.Email,
		NoTelepon:          m.NoTelepon,
		Sektor:             m.Sektor,
		SektorLainnya:      m.SektorLainnya,
		SertifikatTraining: m.SertifikatTraining,
		CreatedAt:          m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          m.UpdatedAt.Format(time.RFC3339),
	}
} 