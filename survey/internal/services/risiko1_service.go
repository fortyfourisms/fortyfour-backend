package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"survey/internal/dto"
	"survey/internal/models"
	validation "survey/validator"
)

// CONFIG
const risikoCacheTTL = 300 * time.Second

// REPOSITORY
type RisikoRepositoryInterface interface {
	ExistsResponden(int64) (bool, error)
	ExistsRisiko(int64) (bool, error)

	UpsertEligibility(models.RisikoEligibility) error
	UpsertAlasan(models.RisikoAlasan) error
	UpsertDampak(models.RisikoDampak) error
	UpsertPengendalian(models.RisikoPengendalian) error

	FindByRespondentID(int64) (map[string]interface{}, error)

	GetProgress(int64) (*models.SurveyProgress, error)
	UpsertProgress(models.SurveyProgress) error
	InsertCustomRisiko(int64, string) (int, error)

	GetRespondentIDByUserID(userID string) (int64, error)
}

// SERVICE
type RisikoService struct {
	repo  RisikoRepositoryInterface
	cache CacheRepository
}

func NewRisikoService(repo RisikoRepositoryInterface, cache CacheRepository) *RisikoService {
	return &RisikoService{repo: repo, cache: cache}
}

// CACHE
func progressKey(id int64) string {
	return fmt.Sprintf("risiko:progress:%d", id)
}

func respondentKey(id int64) string {
	return fmt.Sprintf("risiko:data:%d", id)
}

func (s *RisikoService) setCache(key string, val any) {
	if s.cache == nil {
		return
	}
	b, _ := json.Marshal(val)
	_ = s.cache.Set(context.Background(), key, string(b), int(risikoCacheTTL.Seconds()))
}

func (s *RisikoService) invalidate(id int64) {
	if s.cache == nil {
		return
	}
	_ = s.cache.Del(context.Background(), progressKey(id))
	_ = s.cache.Del(context.Background(), respondentKey(id))
}

// HELPER
func (s *RisikoService) getRespondenID(userID string) (int64, error) {
	id, err := s.repo.GetRespondentIDByUserID(userID)
	if err != nil {
		return 0, errors.New("responden tidak ditemukan")
	}
	return id, nil
}

func toInt(ptr *int) (int, error) {
	if ptr == nil {
		return 0, validation.ErrMissingRisikoID
	}
	return *ptr, nil
}

func toInt64Ptr(v int) *int64 {
	val := int64(v)
	return &val
}

func toStringPtr(s string) *string {
	return &s
}

// STEP 1
func (s *RisikoService) ProcessEligibility(userID string, req dto.EligibilityRequest) (map[string]interface{}, error) {

	if err := validation.ValidateEligibilityRequest(req); err != nil {
		return nil, err
	}

	respondenID, err := s.getRespondenID(userID)
	if err != nil {
		return nil, err
	}

	risikoID, err := toInt(req.RisikoID)
	if err != nil {
		return nil, err
	}

	data := models.RisikoEligibility{
		RespondenID:   int64(respondenID),
		RisikoID:      toInt64Ptr(risikoID),
		PernahTerjadi: req.PernahTerjadi,
	}

	if err := s.repo.UpsertEligibility(data); err != nil {
		return nil, err
	}

	s.invalidate(respondenID)

	next := "reason"
	if req.PernahTerjadi {
		next = "dampak"
	}

	return map[string]interface{}{
		"message":   "eligibility tersimpan",
		"next_step": next,
	}, nil
}

// STEP 2A
func (s *RisikoService) ProcessAlasan(userID string, req dto.AlasanRequest) (map[string]interface{}, error) {

	if err := validation.ValidateAlasanRequest(req); err != nil {
		return nil, err
	}

	respondenID, err := s.getRespondenID(userID)
	if err != nil {
		return nil, err
	}

	risikoID, _ := toInt(req.RisikoID)

	data := models.RisikoAlasan{
		RespondenID: int64(respondenID),
		RisikoID:    toInt64Ptr(risikoID),
		Alasan:      req.Alasan,
	}

	if err := s.repo.UpsertAlasan(data); err != nil {
		return nil, err
	}

	s.invalidate(respondenID)

	return map[string]interface{}{
		"message":   "alasan tersimpan",
		"next_step": "finish",
	}, nil
}

// STEP 2B
func (s *RisikoService) ProcessDampak(userID string, req dto.DampakRequest) (map[string]interface{}, error) {

	if err := validation.ValidateDampakRequest(req); err != nil {
		return nil, err
	}

	respondenID, err := s.getRespondenID(userID)
	if err != nil {
		return nil, err
	}

	risikoID, _ := toInt(req.RisikoID)

	data := models.RisikoDampak{
		RespondenID:       int64(respondenID),
		RisikoID:          toInt64Ptr(risikoID),
		DampakReputasi:    strconv.Itoa(int(req.DampakReputasi)),
		DampakOperasional: strconv.Itoa(int(req.DampakOperasional)),
		DampakFinansial:   strconv.Itoa(int(req.DampakFinansial)),
		DampakHukum:       strconv.Itoa(int(req.DampakHukum)),
		Frekuensi:         strconv.Itoa(int(req.Frekuensi)),
	}

	if err := s.repo.UpsertDampak(data); err != nil {
		return nil, err
	}

	s.invalidate(respondenID)

	return map[string]interface{}{
		"message":   "dampak tersimpan",
		"next_step": "pengendalian",
	}, nil
}

// STEP 2C
func (s *RisikoService) ProcessPengendalian(userID string, req dto.PengendalianRequest) (map[string]interface{}, error) {

	if err := validation.ValidatePengendalianRequest(req); err != nil {
		return nil, err
	}

	respondenID, err := s.getRespondenID(userID)
	if err != nil {
		return nil, err
	}

	risikoID, _ := toInt(req.RisikoID)

	data := models.RisikoPengendalian{
		RespondenID:           int64(respondenID),
		RisikoID:              toInt64Ptr(risikoID),
		AdaPengendalian:       req.AdaPengendalian,
		DeskripsiPengendalian: toStringPtr(req.DeskripsiPengendalian),
	}

	if err := s.repo.UpsertPengendalian(data); err != nil {
		return nil, err
	}

	s.invalidate(respondenID)

	return map[string]interface{}{
		"message":   "pengendalian tersimpan",
		"next_step": "finish",
	}, nil
}

func (s *RisikoService) GetByUserID(userID string) (map[string]interface{}, error) {
	respondenID, err := s.getRespondenID(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByRespondentID(respondenID)
}

func (s *RisikoService) GetByRespondentID(id int64) (map[string]interface{}, error) {
	return s.repo.FindByRespondentID(id)
}

func (s *RisikoService) GetProgress(userID string) (dto.ProgressResponse, error) {
	respondenID, err := s.getRespondenID(userID)
	if err != nil {
		return dto.ProgressResponse{}, err
	}

	progress, err := s.repo.GetProgress(respondenID)
	if err != nil {
		return dto.ProgressResponse{}, err
	}

	var risikoID *int
	if progress.RisikoID.Valid {
		val := int(progress.RisikoID.Int64)
		risikoID = &val
	}

	var langkahSaatIni *string
	if progress.LangkahSaatIni.Valid {
		langkahSaatIni = &progress.LangkahSaatIni.String
	}

	return dto.ProgressResponse{
		RespondenID:    progress.RespondenID,
		RisikoID:       risikoID,
		LangkahSaatIni: langkahSaatIni,
		Selesai:        progress.Selesai,
	}, nil
}

// HELPER SQL
func sqlInt64(v int) sql.NullInt64 {
	return sql.NullInt64{
		Int64: int64(v),
		Valid: true,
	}
}

func intPtr(v int) *int {
	return &v
}

func (s *RisikoService) Navigate(userID string, req dto.NavigateRequest) (dto.ProgressResponse, error) {
	respondenID, err := s.getRespondenID(userID)
	if err != nil {
		return dto.ProgressResponse{}, err
	}

	progress, err := s.repo.GetProgress(respondenID)
	if err != nil {
		return dto.ProgressResponse{}, err
	}

	currentRisk := 0
	if progress.RisikoID.Valid {
		currentRisk = int(progress.RisikoID.Int64)
	}

	switch strings.ToLower(req.Direction) {
	case "next":
		currentRisk++
	case "prev":
		if currentRisk > 1 {
			currentRisk--
		}
	default:
		return dto.ProgressResponse{}, errors.New("invalid direction")
	}

	if currentRisk > 0 {
		progress.RisikoID = sqlInt64(currentRisk)
	} else {
		progress.RisikoID = sql.NullInt64{Valid: false}
	}
	progress.LangkahSaatIni = sql.NullString{String: "navigate", Valid: true}

	if err := s.repo.UpsertProgress(*progress); err != nil {
		return dto.ProgressResponse{}, err
	}

	var risikoID *int
	if currentRisk > 0 {
		risikoID = intPtr(currentRisk)
	}

	return dto.ProgressResponse{
		RespondenID:    int64(respondenID),
		RisikoID:       risikoID,
		LangkahSaatIni: toStringPtr("navigate"),
		Selesai:        progress.Selesai,
	}, nil
}

func (s *RisikoService) SaveProgress(userID string, req dto.NavigateRequest) (dto.ProgressResponse, error) {
	respondenID, err := s.getRespondenID(userID)
	if err != nil {
		return dto.ProgressResponse{}, err
	}

	progress := models.SurveyProgress{
		RespondenID:    int64(respondenID),
		RisikoID:       sqlInt64(req.CurrentRisk),
		LangkahSaatIni: sql.NullString{String: "save-progress", Valid: true},
		Selesai:        false,
	}

	if err := s.repo.UpsertProgress(progress); err != nil {
		return dto.ProgressResponse{}, err
	}

	return dto.ProgressResponse{
		RespondenID:    int64(respondenID),
		RisikoID:       intPtr(req.CurrentRisk),
		LangkahSaatIni: toStringPtr("save-progress"),
		Selesai:        false,
	}, nil
}

func (s *RisikoService) CreateCustomRisiko(req dto.CustomRisikoRequest) (int, error) {
	nama := strings.TrimSpace(req.NamaRisiko)
	if nama == "" {
		return 0, errors.New("nama risiko wajib diisi")
	}

	return s.repo.InsertCustomRisiko(req.RespondenID, nama)
}

func (s *RisikoService) FinishSurvey(userID string) error {
	respondenID, err := s.getRespondenID(userID)
	if err != nil {
		return err
	}

	progress, err := s.repo.GetProgress(respondenID)
	if err != nil {
		return err
	}

	progress.Selesai = true
	progress.LangkahSaatIni = sql.NullString{String: "finish", Valid: true}

	return s.repo.UpsertProgress(*progress)
}
