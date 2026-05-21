package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"survey/internal/dto"
	"survey/internal/models"
	validation "survey/validator"
)

// CONFIG
const risikoCacheTTL = 300 * time.Second

const (
	SurveyStatusDraft         = "draft"
	SurveyStatusSubmitted     = "submitted"
	SurveyStatusEditRequested = "edit_requested"
	SurveyStatusEditApproved  = "edit_approved"
	SurveyStatusEditRejected  = "edit_rejected"
)

// REPOSITORY
type RisikoRepositoryInterface interface {
	GetAllRisiko() ([]models.RisikoResponse, error)
	ExistsResponden(string) (bool, error)
	ExistsRisiko(string) (bool, error)
	GetRisikoIDByUrutan(int) (string, error)
	GetUrutanByRisikoID(string) (int, error)

	UpsertEligibility(models.RisikoEligibility) error
	UpsertAlasan(models.RisikoAlasan) error
	UpsertDampak(models.RisikoDampak) error
	UpsertPengendalian(models.RisikoPengendalian) error

	FindByRespondentID(string) (map[string]interface{}, error)

	GetProgress(string) (*models.SurveyProgress, error)
	UpsertProgress(models.SurveyProgress) error
	InsertCustomRisiko(string, string) (int, error)

	GetRespondentIDByUserID(userID string) (string, error)
}

// SERVICE
type RisikoService struct {
	repo  RisikoRepositoryInterface
	cache CacheRepository
}

func NewRisikoService(repo RisikoRepositoryInterface, cache CacheRepository) *RisikoService {
	return &RisikoService{repo: repo, cache: cache}
}

func (s *RisikoService) GetAllRisiko() ([]models.RisikoResponse, error) {
	return s.repo.GetAllRisiko()
}

// CACHE
func progressKey(id string) string {
	return fmt.Sprintf("risiko:progress:%s", id)
}

func respondentKey(id string) string {
	return fmt.Sprintf("risiko:data:%s", id)
}

func (s *RisikoService) setCache(key string, val any) {
	if s.cache == nil {
		return
	}
	b, _ := json.Marshal(val)
	_ = s.cache.Set(context.Background(), key, string(b), int(risikoCacheTTL.Seconds()))
}

func (s *RisikoService) invalidate(id string) {
	if s.cache == nil {
		return
	}
	_ = s.cache.Del(context.Background(), progressKey(id))
	_ = s.cache.Del(context.Background(), respondentKey(id))
}

// HELPER
func (s *RisikoService) getRespondenID(userID string) (string, error) {
	if strings.TrimSpace(userID) == "" {
		return "", errors.New("user_id wajib diisi")
	}
	id, err := s.repo.GetRespondentIDByUserID(userID)
	if err != nil {
		return "", errors.New("responden tidak ditemukan")
	}
	return id, nil
}

func canEditProgress(progress *models.SurveyProgress) bool {
	status := progress.Status
	if status == "" {
		status = SurveyStatusDraft
	}
	return !progress.Selesai || status == SurveyStatusDraft || status == SurveyStatusEditApproved
}

func (s *RisikoService) ensureEditable(respondenID string) (*models.SurveyProgress, error) {
	progress, err := s.repo.GetProgress(respondenID)
	if err != nil {
		return nil, err
	}
	if !canEditProgress(progress) {
		return nil, errors.New("survey sudah selesai, ajukan request edit ke admin")
	}
	return progress, nil
}

func markDraft(progress *models.SurveyProgress, step string, risikoID string) {
	progress.Selesai = false
	progress.Status = SurveyStatusDraft
	progress.LangkahSaatIni = sql.NullString{String: step, Valid: true}
	if risikoID != "" {
		progress.RisikoID = sql.NullString{String: risikoID, Valid: true}
	}
}

func toStr(ptr *string) (string, error) {
	if ptr == nil {
		return "", validation.ErrMissingRisikoID
	}
	return *ptr, nil
}

func toNullableTrimmedString(s string) *string {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func (s *RisikoService) resolveRisikoIDByUrutan(currentRisk int) sql.NullString {
	if currentRisk <= 0 {
		return sql.NullString{Valid: false}
	}

	risikoID, err := s.repo.GetRisikoIDByUrutan(currentRisk)
	if err != nil || risikoID == "" {
		return sql.NullString{Valid: false}
	}

	return sql.NullString{String: risikoID, Valid: true}
}

func (s *RisikoService) resolveCurrentUrutan(progress *models.SurveyProgress, fallback int) int {
	if progress != nil && progress.RisikoID.Valid {
		if urutan, err := s.repo.GetUrutanByRisikoID(progress.RisikoID.String); err == nil && urutan > 0 {
			return urutan
		}
	}

	if fallback > 0 {
		return fallback
	}

	return 0
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

	risikoID, err := toStr(req.RisikoID)
	if err != nil {
		return nil, err
	}

	progress, err := s.ensureEditable(respondenID)
	if err != nil {
		return nil, err
	}

	data := models.RisikoEligibility{
		RespondenID:   respondenID,
		RisikoID:      req.RisikoID,
		PernahTerjadi: req.PernahTerjadi,
	}

	if err := s.repo.UpsertEligibility(data); err != nil {
		return nil, err
	}
	markDraft(progress, "eligibility", risikoID)
	if err := s.repo.UpsertProgress(*progress); err != nil {
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

	risikoID, err := toStr(req.RisikoID)
	if err != nil {
		return nil, err
	}

	progress, err := s.ensureEditable(respondenID)
	if err != nil {
		return nil, err
	}

	data := models.RisikoAlasan{
		RespondenID: respondenID,
		RisikoID:    req.RisikoID,
		Alasan:      strings.TrimSpace(req.Alasan),
	}

	if err := s.repo.UpsertAlasan(data); err != nil {
		return nil, err
	}
	markDraft(progress, "reason", risikoID)
	if err := s.repo.UpsertProgress(*progress); err != nil {
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

	risikoID, err := toStr(req.RisikoID)
	if err != nil {
		return nil, err
	}

	progress, err := s.ensureEditable(respondenID)
	if err != nil {
		return nil, err
	}

	data := models.RisikoDampak{
		RespondenID:       respondenID,
		RisikoID:          req.RisikoID,
		DampakReputasi:    models.MapImpactIntToString(req.DampakReputasi),
		DampakOperasional: models.MapImpactIntToString(req.DampakOperasional),
		DampakFinansial:   models.MapImpactIntToString(req.DampakFinansial),
		DampakHukum:       models.MapImpactIntToString(req.DampakHukum),
		Frekuensi:         models.MapFrequencyIntToString(req.Frekuensi),
	}

	if err := s.repo.UpsertDampak(data); err != nil {
		return nil, err
	}
	markDraft(progress, "dampak", risikoID)
	if err := s.repo.UpsertProgress(*progress); err != nil {
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

	risikoID, err := toStr(req.RisikoID)
	if err != nil {
		return nil, err
	}

	progress, err := s.ensureEditable(respondenID)
	if err != nil {
		return nil, err
	}

	data := models.RisikoPengendalian{
		RespondenID:           respondenID,
		RisikoID:              req.RisikoID,
		AdaPengendalian:       req.AdaPengendalian,
		DeskripsiPengendalian: toNullableTrimmedString(req.DeskripsiPengendalian),
	}

	if err := s.repo.UpsertPengendalian(data); err != nil {
		return nil, err
	}
	markDraft(progress, "pengendalian", risikoID)
	if err := s.repo.UpsertProgress(*progress); err != nil {
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

func (s *RisikoService) GetByRespondentID(id string) (map[string]interface{}, error) {
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

	return progressToResponse(progress), nil
}

func progressToResponse(progress *models.SurveyProgress) dto.ProgressResponse {
	var risikoID *string
	if progress.RisikoID.Valid {
		val := progress.RisikoID.String
		risikoID = &val
	}

	var langkahSaatIni *string
	if progress.LangkahSaatIni.Valid {
		langkahSaatIni = &progress.LangkahSaatIni.String
	}

	var editReason *string
	if progress.EditReason.Valid {
		editReason = &progress.EditReason.String
	}

	var editResponse *string
	if progress.EditResponse.Valid {
		editResponse = &progress.EditResponse.String
	}

	status := progress.Status
	if status == "" {
		status = SurveyStatusDraft
	}

	return dto.ProgressResponse{
		RespondenID:    progress.RespondenID,
		RisikoID:       risikoID,
		LangkahSaatIni: langkahSaatIni,
		Selesai:        progress.Selesai,
		Status:         status,
		EditReason:     editReason,
		EditResponse:   editResponse,
	}
}

func (s *RisikoService) Navigate(userID string, req dto.NavigateRequest) (dto.ProgressResponse, error) {
	respondenID, err := s.getRespondenID(userID)
	if err != nil {
		return dto.ProgressResponse{}, err
	}

	progress, err := s.ensureEditable(respondenID)
	if err != nil {
		return dto.ProgressResponse{}, err
	}

	currentRisk := s.resolveCurrentUrutan(progress, req.CurrentRisk)

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

	progress.RisikoID = s.resolveRisikoIDByUrutan(currentRisk)
	progress.LangkahSaatIni = sql.NullString{String: "navigate", Valid: true}
	progress.Status = SurveyStatusDraft
	progress.Selesai = false

	if err := s.repo.UpsertProgress(*progress); err != nil {
		return dto.ProgressResponse{}, err
	}
	return progressToResponse(progress), nil
}

func (s *RisikoService) SaveProgress(userID string, req dto.NavigateRequest) (dto.ProgressResponse, error) {
	respondenID, err := s.getRespondenID(userID)
	if err != nil {
		return dto.ProgressResponse{}, err
	}

	progress, err := s.ensureEditable(respondenID)
	if err != nil {
		return dto.ProgressResponse{}, err
	}
	progress.RespondenID = respondenID
	progress.RisikoID = s.resolveRisikoIDByUrutan(req.CurrentRisk)
	progress.LangkahSaatIni = sql.NullString{String: "save-progress", Valid: true}
	progress.Selesai = false
	progress.Status = SurveyStatusDraft

	if err := s.repo.UpsertProgress(*progress); err != nil {
		return dto.ProgressResponse{}, err
	}

	return progressToResponse(progress), nil
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
	progress.Status = SurveyStatusSubmitted
	progress.LangkahSaatIni = sql.NullString{String: "finish", Valid: true}
	progress.SubmittedAt = sql.NullTime{Time: time.Now(), Valid: true}

	return s.repo.UpsertProgress(*progress)
}

func (s *RisikoService) RequestEdit(userID string, req dto.RequestEditRequest) (dto.ProgressResponse, error) {
	respondenID, err := s.getRespondenID(userID)
	if err != nil {
		return dto.ProgressResponse{}, err
	}

	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		return dto.ProgressResponse{}, errors.New("alasan request edit wajib diisi")
	}

	progress, err := s.repo.GetProgress(respondenID)
	if err != nil {
		return dto.ProgressResponse{}, err
	}

	status := progress.Status
	if status == "" {
		status = SurveyStatusDraft
	}
	if !progress.Selesai || status == SurveyStatusDraft || status == SurveyStatusEditApproved {
		return dto.ProgressResponse{}, errors.New("survey masih dapat diedit")
	}
	if status == SurveyStatusEditRequested {
		return dto.ProgressResponse{}, errors.New("request edit masih menunggu persetujuan admin")
	}

	progress.Status = SurveyStatusEditRequested
	progress.EditReason = sql.NullString{String: reason, Valid: true}
	progress.EditResponse = sql.NullString{Valid: false}
	progress.EditRequestedAt = sql.NullTime{Time: time.Now(), Valid: true}
	progress.EditApprovedAt = sql.NullTime{Valid: false}
	progress.EditApprovedBy = sql.NullString{Valid: false}
	progress.EditRejectedAt = sql.NullTime{Valid: false}
	progress.EditRejectedBy = sql.NullString{Valid: false}

	if err := s.repo.UpsertProgress(*progress); err != nil {
		return dto.ProgressResponse{}, err
	}

	return progressToResponse(progress), nil
}

func (s *RisikoService) ReviewEditRequest(adminID string, respondenID string, req dto.ReviewEditRequest) (dto.ProgressResponse, error) {
	if strings.TrimSpace(adminID) == "" {
		return dto.ProgressResponse{}, errors.New("admin_id wajib diisi")
	}
	if strings.TrimSpace(respondenID) == "" {
		return dto.ProgressResponse{}, errors.New("responden_id tidak valid")
	}

	progress, err := s.repo.GetProgress(respondenID)
	if err != nil {
		return dto.ProgressResponse{}, err
	}
	if progress.Status != SurveyStatusEditRequested {
		return dto.ProgressResponse{}, errors.New("tidak ada request edit yang menunggu persetujuan")
	}

	response := strings.TrimSpace(req.Response)
	progress.EditResponse = sql.NullString{String: response, Valid: response != ""}

	if req.Action == "approve" {
		progress.Status = SurveyStatusEditApproved
		progress.Selesai = false
		progress.LangkahSaatIni = sql.NullString{String: "edit-approved", Valid: true}
		progress.EditApprovedAt = sql.NullTime{Time: time.Now(), Valid: true}
		progress.EditApprovedBy = sql.NullString{String: adminID, Valid: true}
	} else {
		progress.Status = SurveyStatusEditRejected
		progress.Selesai = true
		progress.LangkahSaatIni = sql.NullString{String: "edit-rejected", Valid: true}
		progress.EditRejectedAt = sql.NullTime{Time: time.Now(), Valid: true}
		progress.EditRejectedBy = sql.NullString{String: adminID, Valid: true}
	}

	if err := s.repo.UpsertProgress(*progress); err != nil {
		return dto.ProgressResponse{}, err
	}

	return progressToResponse(progress), nil
}
