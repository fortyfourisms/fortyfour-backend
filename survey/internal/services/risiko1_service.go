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

const (
	SurveyStatusDraft         = "draft"
	SurveyStatusSubmitted     = "submitted"
	SurveyStatusEditRequested = "edit_requested"
	SurveyStatusEditApproved  = "edit_approved"
	SurveyStatusEditRejected  = "edit_rejected"
)

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
	if strings.TrimSpace(userID) == "" {
		return 0, errors.New("user_id wajib diisi")
	}
	id, err := s.repo.GetRespondentIDByUserID(userID)
	if err != nil {
		return 0, errors.New("responden tidak ditemukan")
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

func (s *RisikoService) ensureEditable(respondenID int64) (*models.SurveyProgress, error) {
	progress, err := s.ensureEditable(respondenID)
	if err != nil {
		return nil, err
	}
	if !canEditProgress(progress) {
		return nil, errors.New("survey sudah selesai, ajukan request edit ke admin")
	}
	return progress, nil
}

func markDraft(progress *models.SurveyProgress, step string, risikoID int) {
	progress.Selesai = false
	progress.Status = SurveyStatusDraft
	progress.LangkahSaatIni = sql.NullString{String: step, Valid: true}
	if risikoID > 0 {
		progress.RisikoID = sqlInt64(risikoID)
	}
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

	progress, err := s.ensureEditable(respondenID)
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

	risikoID, err := toInt(req.RisikoID)
	if err != nil {
		return nil, err
	}

	progress, err := s.ensureEditable(respondenID)
	if err != nil {
		return nil, err
	}

	data := models.RisikoAlasan{
		RespondenID: int64(respondenID),
		RisikoID:    toInt64Ptr(risikoID),
		Alasan:      req.Alasan,
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

	risikoID, err := toInt(req.RisikoID)
	if err != nil {
		return nil, err
	}

	progress, err := s.ensureEditable(respondenID)
	if err != nil {
		return nil, err
	}

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

	risikoID, err := toInt(req.RisikoID)
	if err != nil {
		return nil, err
	}

	progress, err := s.ensureEditable(respondenID)
	if err != nil {
		return nil, err
	}

	data := models.RisikoPengendalian{
		RespondenID:           int64(respondenID),
		RisikoID:              toInt64Ptr(risikoID),
		AdaPengendalian:       req.AdaPengendalian,
		DeskripsiPengendalian: toStringPtr(req.DeskripsiPengendalian),
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

	return progressToResponse(progress), nil
}

func progressToResponse(progress *models.SurveyProgress) dto.ProgressResponse {
	var risikoID *int
	if progress.RisikoID.Valid {
		val := int(progress.RisikoID.Int64)
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

// HELPER SQL
func sqlInt64(v int) sql.NullInt64 {
	return sql.NullInt64{
		Int64: int64(v),
		Valid: true,
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
	progress.RespondenID = int64(respondenID)
	progress.RisikoID = sqlInt64(req.CurrentRisk)
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
	progress.EditReviewedAt = sql.NullTime{Valid: false}
	progress.EditReviewedBy = sql.NullString{Valid: false}

	if err := s.repo.UpsertProgress(*progress); err != nil {
		return dto.ProgressResponse{}, err
	}

	return progressToResponse(progress), nil
}

func (s *RisikoService) ReviewEditRequest(adminID string, respondenID int64, req dto.ReviewEditRequest) (dto.ProgressResponse, error) {
	if strings.TrimSpace(adminID) == "" {
		return dto.ProgressResponse{}, errors.New("admin_id wajib diisi")
	}
	if respondenID <= 0 {
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
	progress.EditReviewedAt = sql.NullTime{Time: time.Now(), Valid: true}
	progress.EditReviewedBy = sql.NullString{String: adminID, Valid: true}

	if req.Approved {
		progress.Status = SurveyStatusEditApproved
		progress.Selesai = false
		progress.LangkahSaatIni = sql.NullString{String: "edit-approved", Valid: true}
	} else {
		progress.Status = SurveyStatusEditRejected
		progress.Selesai = true
		progress.LangkahSaatIni = sql.NullString{String: "edit-rejected", Valid: true}
	}

	if err := s.repo.UpsertProgress(*progress); err != nil {
		return dto.ProgressResponse{}, err
	}

	return progressToResponse(progress), nil
}
