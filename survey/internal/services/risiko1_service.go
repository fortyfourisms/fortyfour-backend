package services

import (
	"database/sql"
	"errors"
	"survey/internal/dto"
	"survey/internal/models"
	"survey/validator"
)

type RisikoRepositoryInterface interface {
	ExistsResponden(int) (bool, error)
	ExistsRisiko(int) (bool, error)
	ExistsCustomRisiko(int) (bool, error)
	UpsertEligibility(models.RisikoEligibility) error
	UpsertAlasan(models.RisikoAlasan) error
	UpsertDampak(models.RisikoDampak) error
	UpsertPengendalian(models.RisikoPengendalian) error
	FindByRespondentID(int) (map[string]interface{}, error)
	GetProgress(int) (*models.SurveyProgress, error)
	UpsertProgress(models.SurveyProgress) error
	InsertCustomRisiko(int, string) (int, error)
}

const CustomRiskIndex = 14

type RisikoService struct {
	repo RisikoRepositoryInterface
}

func NewRisikoService(repo RisikoRepositoryInterface) *RisikoService {
	return &RisikoService{repo: repo}
}

// HELPER
func toIntPtr(v int) *int {
	if v == 0 {
		return nil
	}
	return &v
}

func (s *RisikoService) validateFK(respondenID, risikoID, customID int) error {

	existResponden, err := s.repo.ExistsResponden(respondenID)
	if err != nil {
		return err
	}
	if !existResponden {
		return errors.New("responden tidak ditemukan")
	}

	// custom risk
	if customID != 0 {
		exist, err := s.repo.ExistsCustomRisiko(customID)
		if err != nil {
			return err
		}
		if !exist {
			return errors.New("custom risiko tidak ditemukan")
		}
		return nil
	}

	// master risk
	exist, err := s.repo.ExistsRisiko(risikoID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("risiko tidak ditemukan")
	}

	return nil
}

func assignRisk(risikoID int, customID int) (int, *int) {
	if customID != 0 {
		return 0, toIntPtr(customID)
	}
	return risikoID, nil
}

// STEP 1 — ELIGIBILITY
func (s *RisikoService) ProcessEligibility(req dto.EligibilityRequest) (map[string]interface{}, error) {

	if req.RespondenID == 0 {
		return nil, validation.ErrMissingRespondentID
	}

	if err := s.validateFK(req.RespondenID, req.RisikoID, req.CustomRisikoID); err != nil {
		return nil, err
	}

	risikoID, customPtr := assignRisk(req.RisikoID, req.CustomRisikoID)

	data := models.RisikoEligibility{
		RespondenID:    req.RespondenID,
		RisikoID:       risikoID,
		CustomRisikoID: customPtr,
		PernahTerjadi:  req.PernahTerjadi,
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

// STEP 2A — ALASAN
func (s *RisikoService) ProcessAlasan(req dto.AlasanRequest) (map[string]interface{}, error) {

	if req.RespondenID == 0 {
		return nil, validation.ErrMissingRespondentID
	}

	if req.Alasan == "" {
		return nil, validation.ErrMissingReason
	}

	if err := s.validateFK(req.RespondenID, req.RisikoID, req.CustomRisikoID); err != nil {
		return nil, err
	}

	risikoID, customPtr := assignRisk(req.RisikoID, req.CustomRisikoID)

	data := models.RisikoAlasan{
		RespondenID:    req.RespondenID,
		RisikoID:       risikoID,
		CustomRisikoID: customPtr,
		Alasan:         req.Alasan,
	}

	if err := s.repo.UpsertAlasan(data); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message":   "alasan tersimpan",
		"next_step": "finish",
	}, nil
}

// STEP 2B — DAMPAK
func (s *RisikoService) ProcessDampak(req dto.DampakRequest) (map[string]interface{}, error) {

	if req.RespondenID == 0 {
		return nil, validation.ErrMissingRespondentID
	}

	if err := s.validateFK(req.RespondenID, req.RisikoID, req.CustomRisikoID); err != nil {
		return nil, err
	}

	if !req.DampakReputasi.Valid() ||
		!req.DampakOperasional.Valid() ||
		!req.DampakFinansial.Valid() ||
		!req.DampakHukum.Valid() {
		return nil, validation.ErrInvalidImpact
	}

	if !req.Frekuensi.Valid() {
		return nil, validation.ErrInvalidFreq
	}

	risikoID, customPtr := assignRisk(req.RisikoID, req.CustomRisikoID)

	data := models.RisikoDampak{
		RespondenID:       req.RespondenID,
		RisikoID:          risikoID,
		CustomRisikoID:    customPtr,
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

	if req.AdaPengendalian && req.DeskripsiPengendalian == "" {
		return nil, validation.ErrMissingControl
	}

	if err := s.validateFK(req.RespondenID, req.RisikoID, req.CustomRisikoID); err != nil {
		return nil, err
	}

	risikoID, customPtr := assignRisk(req.RisikoID, req.CustomRisikoID)

	data := models.RisikoPengendalian{
		RespondenID:           req.RespondenID,
		RisikoID:              risikoID,
		CustomRisikoID:        customPtr,
		AdaPengendalian:       req.AdaPengendalian,
		DeskripsiPengendalian: req.DeskripsiPengendalian,
	}

	if err := s.repo.UpsertPengendalian(data); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message":   "pengendalian tersimpan",
		"next_step": "finish",
	}, nil
}

func (s *RisikoService) GetByRespondentID(id int) (map[string]interface{}, error) {
	return s.repo.FindByRespondentID(id)
}

func (s *RisikoService) GetProgress(id int) (dto.ProgressResponse, error) {

	progress, err := s.repo.GetProgress(id)
	if err != nil {
		return dto.ProgressResponse{}, err
	}

	return mapProgressToResponse(progress), nil
}

// NAVIGATE
func (s *RisikoService) Navigate(req dto.NavigateRequest) (dto.ProgressResponse, error) {

	progress, err := s.repo.GetProgress(req.RespondenID)
	if err != nil {
		return dto.ProgressResponse{}, err
	}

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
		return dto.ProgressResponse{}, errors.New("direction tidak valid")
	}

	step := "normal"
	if current == CustomRiskIndex {
		step = "custom_risk"
	}

	progress.RisikoID = sql.NullInt64{Int64: current, Valid: true}
	progress.LangkahSaatIni = sql.NullString{String: step, Valid: true}

	if err := s.repo.UpsertProgress(*progress); err != nil {
		return dto.ProgressResponse{}, err
	}

	return mapProgressToResponse(progress), nil
}

// SAVE PROGRESS
func (s *RisikoService) SaveProgress(req dto.NavigateRequest) (dto.ProgressResponse, error) {

	progress, err := s.repo.GetProgress(req.RespondenID)
	if err != nil {
		return dto.ProgressResponse{}, err
	}

	progress.RisikoID = sql.NullInt64{
		Int64: int64(req.CurrentRisk),
		Valid: true,
	}

	progress.LangkahSaatIni = sql.NullString{
		String: "paused",
		Valid:  true,
	}

	progress.Selesai = false

	if err := s.repo.UpsertProgress(*progress); err != nil {
		return dto.ProgressResponse{}, err
	}

	return mapProgressToResponse(progress), nil
}

// CUSTOM RISIKO
func (s *RisikoService) CreateCustomRisiko(req dto.CustomRisikoRequest) (int, error) {

	if req.NamaRisiko == "" {
		return 0, errors.New("nama risiko wajib diisi")
	}

	return s.repo.InsertCustomRisiko(req.RespondenID, req.NamaRisiko)
}

// FINISH
func (s *RisikoService) FinishSurvey(respondenID int) error {

	progress, err := s.repo.GetProgress(respondenID)
	if err != nil {
		return err
	}

	progress.Selesai = true
	progress.LangkahSaatIni = sql.NullString{
		String: "finish",
		Valid:  true,
	}

	return s.repo.UpsertProgress(*progress)
}

// MAPPING
func mapProgressToResponse(p *models.SurveyProgress) dto.ProgressResponse {

	var risikoID *int
	if p.RisikoID.Valid {
		val := int(p.RisikoID.Int64)
		risikoID = &val
	}

	var langkah *string
	if p.LangkahSaatIni.Valid {
		langkah = &p.LangkahSaatIni.String
	}

	return dto.ProgressResponse{
		RespondenID:    p.RespondenID,
		RisikoID:       risikoID,
		LangkahSaatIni: langkah,
		Selesai:        p.Selesai,
	}
}
