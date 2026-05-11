package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"survey/internal/dto"
	"survey/internal/middleware"
	"survey/internal/models"
	"survey/internal/repository"
	"survey/internal/utils"
)

// SERVICE INTERFACE
type RisikoServiceInterface interface {
	GetAllRisiko() ([]models.RisikoResponse, error)

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
	RequestEdit(userID string, req dto.RequestEditRequest) (dto.ProgressResponse, error)
	ReviewEditRequest(adminID string, respondenID int64, req dto.ReviewEditRequest) (dto.ProgressResponse, error)
}

// HANDLER STRUCT
type RisikoHandler struct {
	svc RisikoServiceInterface
}

func NewRisikoHandler(svc RisikoServiceInterface) *RisikoHandler {
	return &RisikoHandler{svc: svc}
}

// GetAllRisiko handles GET /api/survey/risiko
func (h *RisikoHandler) GetAllRisiko(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetAllRisiko()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Berhasil mengambil data risiko", data)
}

// SubmitEligibility handles POST /api/survey/risiko/eligibility
func (h *RisikoHandler) SubmitEligibility(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req dto.EligibilityRequest
	if err := decodeJSON(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	res, err := h.svc.ProcessEligibility(userID, req)
	handleResult(w, res, err)
}

// SubmitAlasan handles POST /api/survey/risiko/reason
func (h *RisikoHandler) SubmitAlasan(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req dto.AlasanRequest
	if err := decodeJSON(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	res, err := h.svc.ProcessAlasan(userID, req)
	handleResult(w, res, err)
}

// SubmitDampak handles POST /api/survey/risiko/dampak
func (h *RisikoHandler) SubmitDampak(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req dto.DampakRequest
	if err := decodeJSON(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	res, err := h.svc.ProcessDampak(userID, req)
	handleResult(w, res, err)
}

// SubmitPengendalian handles POST /api/survey/risiko/pengendalian
func (h *RisikoHandler) SubmitPengendalian(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req dto.PengendalianRequest
	if err := decodeJSON(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	res, err := h.svc.ProcessPengendalian(userID, req)
	handleResult(w, res, err)
}

// GetMe handles GET /api/survey/risiko/me
func (h *RisikoHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	data, err := h.svc.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.RespondError(w, http.StatusNotFound, "Data risiko belum ditemukan")
			return
		}
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Berhasil mengambil data risiko", data)
}

// GetByRespondentID handles GET /api/survey/risiko/{id}
func (h *RisikoHandler) GetByRespondentID(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())

	if role != "admin" && role != "staff" {
		utils.RespondError(w, http.StatusForbidden, "Forbidden")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/survey/risiko/")
	id, err := strconv.ParseInt(strings.Trim(idStr, "/"), 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	data, err := h.svc.GetByRespondentID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.RespondError(w, http.StatusNotFound, "Data tidak ditemukan")
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Berhasil mengambil data risiko responden", data)
}

// GetProgress handles GET /api/survey/progress
func (h *RisikoHandler) GetProgress(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	res, err := h.svc.GetProgress(userID)
	handleResult(w, res, err)
}

// Navigate handles POST /api/survey/navigate
func (h *RisikoHandler) Navigate(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req dto.NavigateRequest
	_ = decodeJSON(r, &req)

	res, err := h.svc.Navigate(userID, req)
	handleResult(w, res, err)
}

// SaveProgress handles POST /api/survey/save-progress
func (h *RisikoHandler) SaveProgress(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req dto.NavigateRequest
	_ = decodeJSON(r, &req)

	res, err := h.svc.SaveProgress(userID, req)
	handleResult(w, res, err)
}

// FinishSurvey handles POST /api/survey/finish
func (h *RisikoHandler) FinishSurvey(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	if err := h.svc.FinishSurvey(userID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Survey berhasil diselesaikan", nil)
}

func (h *RisikoHandler) RequestEdit(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req dto.RequestEditRequest
	if err := decodeJSON(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	res, err := h.svc.RequestEdit(userID, req)
	handleResult(w, res, err)
}

func (h *RisikoHandler) ReviewEditRequest(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "admin" && role != "staff" {
		utils.RespondError(w, http.StatusForbidden, "Forbidden")
		return
	}

	adminID := middleware.GetUserID(r.Context())
	idStr := strings.TrimPrefix(r.URL.Path, "/api/survey/edit-requests/")
	respondenID, err := strconv.ParseInt(strings.Trim(idStr, "/"), 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID responden tidak valid")
		return
	}

	var req dto.ReviewEditRequest
	if err := decodeJSON(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	res, err := h.svc.ReviewEditRequest(adminID, respondenID, req)
	handleResult(w, res, err)
}

// HELPERS
func decodeJSON(r *http.Request, dst interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

func handleResult(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(w, http.StatusOK, "Success", data)
}
