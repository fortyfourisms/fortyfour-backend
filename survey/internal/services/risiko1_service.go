package services

import (
	"errors"
	"survey/internal/dto"
	"survey/internal/models"
	"survey/internal/repository"
	"survey/validator"
	"database/sql"
)

type RisikoService struct {
	repo *repository.RisikoRepository
}

func NewRisikoService(repo *repository.RisikoRepository) *RisikoService {
	return &RisikoService{repo: repo}
}

// STEP 1 — ELIGIBILITY
func (s *RisikoService) ProcessEligibility(req dto.EligibilityRequest) (map[string]interface{}, error) {

	if req.RespondenID == 0 {
		return nil, validation.ErrMissingRespondentID
	}

	if err := s.validateForeignKey(req.RespondenID, req.RisikoID); err != nil {
	return nil, err
	}

	data := models.RisikoEligibility{
		RespondenID:   req.RespondenID,
		RisikoID:      req.RisikoID,
		PernahTerjadi: req.PernahTerjadi,
	}

	if err := s.repo.UpsertEligibility(data); err != nil {
		return nil, err
	}

	nextStep := "reason"
	if req.PernahTerjadi {
		nextStep = "dampak"
	}

	return map[string]interface{}{
		"message":   "eligibility tersimpan",
		"next_step": nextStep,
	}, nil
}

// STEP 2A — ALASAN (JIKA TIDAK)
func (s *RisikoService) ProcessAlasan(req dto.AlasanRequest) (map[string]interface{}, error) {

	if req.RespondenID == 0 {
		return nil, validation.ErrMissingRespondentID
	}

	if err := s.validateForeignKey(req.RespondenID, req.RisikoID); err != nil {
	return nil, err
	}

	if req.Alasan == "" {
		return nil, validation.ErrMissingReason
	}

	data := models.RisikoAlasan{
		RespondenID: req.RespondenID,
		RisikoID:    req.RisikoID,
		Alasan:      req.Alasan,
	}

	if err := s.repo.UpsertAlasan(data); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message":   "alasan tersimpan",
		"next_step": "finish",
	}, nil
}

// STEP 2B — DAMPAK (JIKA YA)
func (s *RisikoService) ProcessDampak(req dto.DampakRequest) (map[string]interface{}, error) {

	if req.RespondenID == 0 {
		return nil, validation.ErrMissingRespondentID
	}

	if err := s.validateForeignKey(req.RespondenID, req.RisikoID); err != nil {
	return nil, err
	}

	// validasi impact & frekuensi
	if !req.DampakReputasi.Valid() ||
		!req.DampakOperasional.Valid() ||
		!req.DampakFinansial.Valid() ||
		!req.DampakHukum.Valid() {
		return nil, validation.ErrInvalidImpact
	}

	if !req.Frekuensi.Valid() {
		return nil, validation.ErrInvalidFreq
	}

	data := models.RisikoDampak{
		RespondenID:       req.RespondenID,
		RisikoID:          req.RisikoID,
		DampakReputasi:    req.DampakReputasi,
		DampakOperasional: req.DampakOperasional,
		DampakFinansial:   req.DampakFinansial,
		DampakHukum:       req.DampakHukum,
		Frekuensi:         req.Frekuensi,
	}

	if err := s.repo.UpsertDampak(data); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message":   "dampak & frekuensi tersimpan",
		"next_step": "pengendalian",
	}, nil
}

// STEP 2C — PENGENDALIAN
func (s *RisikoService) ProcessPengendalian(req dto.PengendalianRequest) (map[string]interface{}, error) {

	if req.RespondenID == 0 {
		return nil, validation.ErrMissingRespondentID
	}

	if err := s.validateForeignKey(req.RespondenID, req.RisikoID); err != nil {
	return nil, err
	}

	if req.AdaPengendalian && req.DeskripsiPengendalian == "" {
		return nil, validation.ErrMissingControl
	}

	data := models.RisikoPengendalian{
		RespondenID:           req.RespondenID,
		RisikoID:              req.RisikoID,
		AdaPengendalian:       req.AdaPengendalian,
		DeskripsiPengendalian: req.DeskripsiPengendalian,
	}

	if err := s.repo.UpsertPengendalian(data); err != nil {
		return nil, err
	}

	msg := "pengendalian tersimpan"
	if !req.AdaPengendalian {
		msg = "tidak ada pengendalian"
	}

	return map[string]interface{}{
		"message":   msg,
		"next_step": "finish",
	}, nil
}

// GET FULL DATA
func (s *RisikoService) GetByRespondentID(id int) (map[string]interface{}, error) {
	return s.repo.FindByRespondentID(id)
}

// PROGRESS
func (s *RisikoService) GetProgress(id int) (*models.SurveyProgress, error) {
	return s.repo.GetProgress(id)
}

// NAVIGATE
func (s *RisikoService) Navigate(req dto.NavigateRequest) (*models.SurveyProgress, error) {

	progress, err := s.repo.GetProgress(req.RespondenID)
	if err != nil {
		return nil, err
	}

	// ambil nilai sekarang (handle NULL)
	current := progress.RisikoID.Int64
	if !progress.RisikoID.Valid {
		current = 1
	}

	switch req.Direction {

	case "next":
		current++

	case "previous":
		if current > 1 {
			current--
		}

	default:
		return nil, errors.New("direction tidak valid")
	}

	// set kembali ke struct
	progress.RisikoID = sql.NullInt64{
		Int64: current,
		Valid: true,
	}

	err = s.repo.UpsertProgress(*progress)
	if err != nil {
		return nil, err
	}

	return progress, nil
}

// VALIDATE FOREIGN KEY
func (s *RisikoService) validateForeignKey(respondenID, risikoID int) error {

	existResponden, err := s.repo.ExistsResponden(respondenID)
	if err != nil {
		return err
	}
	if !existResponden {
		return errors.New("responden tidak ditemukan")
	}

	existRisiko, err := s.repo.ExistsRisiko(risikoID)
	if err != nil {
		return err
	}
	if !existRisiko {
		return errors.New("risiko tidak ditemukan")
	}

	return nil
}