package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"survey/internal/dto"
	"survey/internal/models"
	"survey/validator"
)

const CustomRiskIndex = 14

// =======================
// REPOSITORY INTERFACE
// =======================
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

// =======================
// REDIS CACHE INTERFACE
// =======================
// type CacheRepository interface {
// 	Get(ctx context.Context, key string) (string, error)
// 	Set(ctx context.Context, key string, value string, ttlSeconds int) error
// 	Del(ctx context.Context, key string) error
// }

// =======================
// SERVICE
// =======================
type RisikoService struct {
	repo  RisikoRepositoryInterface
	cache CacheRepository
}

func NewRisikoService(repo RisikoRepositoryInterface, cache CacheRepository) *RisikoService {
	return &RisikoService{
		repo:  repo,
		cache: cache,
	}
}

// =======================
// CACHE KEY
// =======================
func progressKey(id int) string {
	return fmt.Sprintf("progress:%d", id)
}

func respondentKey(id int) string {
	return fmt.Sprintf("respondent:%d", id)
}

// =======================
// CACHE HELPERS
// =======================
func (s *RisikoService) setCache(key string, value any, ttl int) {
	if s.cache == nil {
		return
	}

	b, err := json.Marshal(value)
	if err != nil {
		return
	}

	_ = s.cache.Set(context.Background(), key, string(b), ttl)
}

func (s *RisikoService) getCache(key string, target any) bool {
	if s.cache == nil {
		return false
	}

	data, ok, err := s.cache.Get(context.Background(), key)
	if err != nil || !ok {
		return false
	}

	return json.Unmarshal([]byte(data), target) == nil
}

func (s *RisikoService) invalidate(respondenID int) {
	if s.cache == nil {
		return
	}

	_ = s.cache.Del(context.Background(), progressKey(respondenID))
	_ = s.cache.Del(context.Background(), respondentKey(respondenID))
}

// =======================
// HELPER
// =======================
func toIntPtr(v int) *int {
	if v == 0 {
		return nil
	}
	return &v
}

func assignRisk(risikoID int, customID int) (int, *int) {
	if customID != 0 {
		return 0, toIntPtr(customID)
	}
	return risikoID, nil
}

// =======================
// VALIDATION FK
// =======================
func (s *RisikoService) validateFK(respondenID, risikoID, customID int) error {

	existResponden, err := s.repo.ExistsResponden(respondenID)
	if err != nil {
		return err
	}
	if !existResponden {
		return errors.New("responden tidak ditemukan")
	}

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

	exist, err := s.repo.ExistsRisiko(risikoID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("risiko tidak ditemukan")
	}

	return nil
}

// =======================
// STEP 1 ELIGIBILITY
// =======================
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

	s.invalidate(req.RespondenID)

	nextStep := "reason"
	if req.PernahTerjadi {
		nextStep = "dampak"
	}

	return map[string]interface{}{
		"message":   "eligibility tersimpan",
		"next_step": nextStep,
	}, nil
}

// =======================
// STEP 2A ALASAN
// =======================
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

	s.invalidate(req.RespondenID)

	return map[string]interface{}{
		"message":   "alasan tersimpan",
		"next_step": "finish",
	}, nil
}

// =======================
// STEP 2B DAMPAK
// =======================
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

	s.invalidate(req.RespondenID)

	return map[string]interface{}{
		"message":   "dampak tersimpan",
		"next_step": "pengendalian",
	}, nil
}

// =======================
// STEP 2C PENGENDALIAN
// =======================
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

	s.invalidate(req.RespondenID)

	return map[string]interface{}{
		"message":   "pengendalian tersimpan",
		"next_step": "finish",
	}, nil
}

// =======================
// GET PROGRESS (CACHE)
// =======================
func (s *RisikoService) GetProgress(id int) (dto.ProgressResponse, error) {

	var cached dto.ProgressResponse
	if s.getCache(progressKey(id), &cached) {
		return cached, nil
	}

	progress, err := s.repo.GetProgress(id)
	if err != nil {
		return dto.ProgressResponse{}, err
	}

	resp := mapProgressToResponse(progress)
	s.setCache(progressKey(id), resp, 300)

	return resp, nil
}

// =======================
// FIND RESPONDENT (CACHE)
// =======================
func (s *RisikoService) GetByRespondentID(id int) (map[string]interface{}, error) {

	var cached map[string]interface{}
	if s.getCache(respondentKey(id), &cached) {
		return cached, nil
	}

	data, err := s.repo.FindByRespondentID(id)
	if err != nil {
		return nil, err
	}

	s.setCache(respondentKey(id), data, 300)
	return data, nil
}

// =======================
// NAVIGATE
// =======================
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

	s.invalidate(req.RespondenID)

	return mapProgressToResponse(progress), nil
}

// =======================
// SAVE PROGRESS
// =======================
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

	s.invalidate(req.RespondenID)

	return mapProgressToResponse(progress), nil
}

// =======================
// CUSTOM RISIKO
// =======================
func (s *RisikoService) CreateCustomRisiko(req dto.CustomRisikoRequest) (int, error) {

	if req.NamaRisiko == "" {
		return 0, errors.New("nama risiko wajib diisi")
	}

	return s.repo.InsertCustomRisiko(req.RespondenID, req.NamaRisiko)
}

// =======================
// FINISH
// =======================
func (s *RisikoService) FinishSurvey(respondenID int) error {

	progress, err := s.repo.GetProgress(respondenID)
	if err != nil {
		return err
	}

	progress.Selesai = true
	progress.LangkahSaatIni = sql.NullString{String: "finish", Valid: true}

	s.invalidate(respondenID)

	return s.repo.UpsertProgress(*progress)
}

// =======================
// MAPPING
// =======================
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