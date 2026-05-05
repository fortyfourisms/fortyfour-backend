package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"survey/internal/dto"
	"survey/internal/middleware"
	"survey/internal/repository"
)

// SERVICE INTERFACE
type RisikoServiceInterface interface {
	ProcessEligibility(userID string, req dto.EligibilityRequest) (map[string]interface{}, error)
	ProcessAlasan(userID string, req dto.AlasanRequest) (map[string]interface{}, error)
	ProcessDampak(userID string, req dto.DampakRequest) (map[string]interface{}, error)
	ProcessPengendalian(userID string, req dto.PengendalianRequest) (map[string]interface{}, error)

	GetByUserID(userID string) (map[string]interface{}, error)
	GetByRespondentID(id int64) (map[string]interface{}, error)

	GetProgress(userID string) (dto.ProgressResponse, error)

	Navigate(userID string, req dto.NavigateRequest) (dto.ProgressResponse, error)
	SaveProgress(userID string, req dto.NavigateRequest) (dto.ProgressResponse, error)

	FinishSurvey(userID string) error
}

// HANDLER STRUCT
type RisikoHandler struct {
	svc RisikoServiceInterface
}

func NewRisikoHandler(svc RisikoServiceInterface) *RisikoHandler {
	return &RisikoHandler{svc: svc}
}

// HANDLER METHODS

// @Summary Submit Eligibility
// @Description Submit eligibility for risk assessment
// @Tags Risiko
// @Accept json
// @Produce json
// @Param request body dto.EligibilityRequest true "Eligibility request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/survey/risiko/eligibility [post]
func (h *RisikoHandler) SubmitEligibility(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req dto.EligibilityRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, 400, "invalid body")
		return
	}

	res, err := h.svc.ProcessEligibility(userID, req)
	handleResult(w, res, err)
}

// @Summary Submit Alasan
// @Description Submit reason for risk
// @Tags Risiko
// @Accept json
// @Produce json
// @Param request body dto.AlasanRequest true "Alasan request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/survey/risiko/reason [post]
func (h *RisikoHandler) SubmitAlasan(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req dto.AlasanRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, 400, "invalid body")
		return
	}

	res, err := h.svc.ProcessAlasan(userID, req)
	handleResult(w, res, err)
}

// @Summary Submit Dampak
// @Description Submit impact assessment for risk
// @Tags Risiko
// @Accept json
// @Produce json
// @Param request body dto.DampakRequest true "Dampak request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/survey/risiko/dampak [post]
func (h *RisikoHandler) SubmitDampak(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req dto.DampakRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, 400, "invalid body")
		return
	}

	res, err := h.svc.ProcessDampak(userID, req)
	handleResult(w, res, err)
}

// @Summary Submit Pengendalian
// @Description Submit control measures for risk
// @Tags Risiko
// @Accept json
// @Produce json
// @Param request body dto.PengendalianRequest true "Pengendalian request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/survey/risiko/pengendalian [post]
func (h *RisikoHandler) SubmitPengendalian(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req dto.PengendalianRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, 400, "invalid body")
		return
	}

	res, err := h.svc.ProcessPengendalian(userID, req)
	handleResult(w, res, err)
}

// @Summary Get My Risiko Data
// @Description Get current user's risk data
// @Tags Risiko
// @Accept json
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Failure 404 {object} dto.ErrorResponse
func (h *RisikoHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	data, err := h.svc.GetByUserID(userID)
	if err != nil {
		writeError(w, 404, err.Error())
		return
	}

	writeSuccess(w, data)
}

// @Summary Get Risiko by Respondent ID
// @Description Get risk data by respondent ID (admin only)
// @Tags Risiko
// @Accept json
// @Produce json
// @Param id path int true "Respondent ID"
// @Success 200 {object} dto.APIResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/survey/risiko/{id} [get]
func (h *RisikoHandler) GetByRespondentID(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())

	if role != "admin" {
		writeError(w, 403, "forbidden")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/survey/risiko/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, 400, "invalid id")
		return
	}

	data, err := h.svc.GetByRespondentID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, 404, "data tidak ditemukan")
			return
		}
		writeError(w, 500, err.Error())
		return
	}

	writeSuccess(w, data)
}

// @Summary Navigate Survey
// @Description Navigate through survey steps
// @Tags Risiko
// @Accept json
// @Produce json
// @Param request body dto.NavigateRequest true "Navigate request"
// @Success 200 {object} dto.APIResponse{data=dto.ProgressResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/survey/navigate [post]
func (h *RisikoHandler) Navigate(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req dto.NavigateRequest
	_ = decodeJSON(r, &req)

	res, err := h.svc.Navigate(userID, req)
	handleResult(w, res, err)
}

// @Summary Save Progress
// @Description Save current survey progress
// @Tags Risiko
// @Accept json
// @Produce json
// @Param request body dto.NavigateRequest true "Save progress request"
// @Success 200 {object} dto.APIResponse{data=dto.ProgressResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/survey/save-progress [post]
func (h *RisikoHandler) SaveProgress(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req dto.NavigateRequest
	_ = decodeJSON(r, &req)

	res, err := h.svc.SaveProgress(userID, req)
	handleResult(w, res, err)
}

// @Summary Finish Survey
// @Description Complete the survey
// @Tags Risiko
// @Accept json
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/survey/finish [post]
func (h *RisikoHandler) FinishSurvey(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	if err := h.svc.FinishSurvey(userID); err != nil {
		writeError(w, 500, err.Error())
		return
	}

	writeSuccess(w, "survey selesai")
}

// HELPERS
func decodeJSON(r *http.Request, dst interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

func handleResult(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, 404, err.Error())
			return
		}
		writeError(w, 400, err.Error())
		return
	}
	writeSuccess(w, data)
}

func writeSuccess(w http.ResponseWriter, data interface{}) {
	writeJSON(w, 200, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]interface{}{
		"success": false,
		"message": message,
	})
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
