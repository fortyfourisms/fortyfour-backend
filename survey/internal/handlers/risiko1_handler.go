package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
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
	GetByRespondentID(id string) (map[string]interface{}, error)

	GetProgress(userID string) (dto.ProgressResponse, error)

	Navigate(userID string, req dto.NavigateRequest) (dto.ProgressResponse, error)
	SaveProgress(userID string, req dto.NavigateRequest) (dto.ProgressResponse, error)

	FinishSurvey(userID string) error
	RequestEdit(userID string, req dto.RequestEditRequest) (dto.ProgressResponse, error)
	ReviewEditRequest(adminID string, respondenID string, req dto.ReviewEditRequest) (dto.ProgressResponse, error)
}

// HANDLER STRUCT
type RisikoHandler struct {
	svc RisikoServiceInterface
}

func NewRisikoHandler(svc RisikoServiceInterface) *RisikoHandler {
	return &RisikoHandler{svc: svc}
}

// GetAllRisiko godoc
// @Summary Get all reference risks
// @Description Retrieve a list of all master risks available for the survey.
// @Tags Survey Risiko
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Berhasil mengambil data risiko"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/survey/risiko [get]
func (h *RisikoHandler) GetAllRisiko(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetAllRisiko()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Berhasil mengambil data risiko", data)
}

// SubmitEligibility godoc
// @Summary Submit eligibility answer
// @Description Submit whether a risk has occurred.
// @Tags Survey Risiko
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.EligibilityRequest true "Eligibility Data"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/survey/risiko/eligibility [post]
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

// SubmitAlasan godoc
// @Summary Submit reason for risk
// @Description Submit the reason why a risk occurred.
// @Tags Survey Risiko
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.AlasanRequest true "Reason Data"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/survey/risiko/reason [post]
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

// SubmitDampak godoc
// @Summary Submit risk impacts
// @Description Submit the impacts of the risk.
// @Tags Survey Risiko
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.DampakRequest true "Impact Data"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/survey/risiko/dampak [post]
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

// SubmitPengendalian godoc
// @Summary Submit risk control
// @Description Submit risk mitigation controls.
// @Tags Survey Risiko
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.PengendalianRequest true "Control Data"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/survey/risiko/pengendalian [post]
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

// GetMe godoc
// @Summary Get my risks
// @Description Retrieve the risk items for the currently authenticated user.
// @Tags Survey Risiko
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Berhasil mengambil data risiko"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Data risiko belum ditemukan"
// @Router /api/survey/risiko/me [get]
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

// GetByRespondentID godoc
// @Summary Get risks by respondent ID (Admin/Staff)
// @Description Retrieve the risk items for a specific respondent.
// @Tags Survey Risiko
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Respondent UUID"
// @Success 200 {object} map[string]interface{} "Berhasil mengambil data risiko responden"
// @Failure 400 {object} map[string]interface{} "ID tidak valid"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Data tidak ditemukan"
// @Router /api/survey/risiko/{id} [get]
func (h *RisikoHandler) GetByRespondentID(w http.ResponseWriter, r *http.Request) {
	role := strings.ToLower(strings.TrimSpace(middleware.GetRole(r.Context())))

	if role != "admin" && role != "staff" {
		utils.RespondError(w, http.StatusForbidden, "Forbidden")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/survey/risiko/")
	id := strings.Trim(idStr, "/")
	if id == "" {
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

// GetProgress godoc
// @Summary Get current survey progress
// @Description Retrieve the survey progress of the currently authenticated user.
// @Tags Survey Progress
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/survey/progress [get]
func (h *RisikoHandler) GetProgress(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	res, err := h.svc.GetProgress(userID)
	handleResult(w, res, err)
}

// Navigate godoc
// @Summary Navigate survey
// @Description Navigate between steps (prev/next) in the survey.
// @Tags Survey Progress
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.NavigateRequest false "Navigation Data"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/survey/navigate [post]
func (h *RisikoHandler) Navigate(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req dto.NavigateRequest
	if err := decodeOptionalJSON(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	res, err := h.svc.Navigate(userID, req)
	handleResult(w, res, err)
}

// SaveProgress godoc
// @Summary Save survey progress
// @Description Save current state and progress.
// @Tags Survey Progress
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.NavigateRequest false "Progress Data"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/survey/save-progress [post]
func (h *RisikoHandler) SaveProgress(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req dto.NavigateRequest
	if err := decodeOptionalJSON(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	res, err := h.svc.SaveProgress(userID, req)
	handleResult(w, res, err)
}

// FinishSurvey godoc
// @Summary Finish survey
// @Description Submit and finish the survey.
// @Tags Survey Progress
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Survey berhasil diselesaikan"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/survey/finish [post]
func (h *RisikoHandler) FinishSurvey(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	if err := h.svc.FinishSurvey(userID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Survey berhasil diselesaikan", nil)
}

// RequestEdit godoc
// @Summary Request to edit a submitted survey
// @Description Request permission to edit a survey that has already been submitted.
// @Tags Survey Progress
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.RequestEditRequest true "Edit Request Data"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/survey/edit-requests [post]
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

// ReviewEditRequest godoc
// @Summary Review survey edit request (Admin)
// @Description Approve or reject a user's request to edit their survey.
// @Tags Survey Progress
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Respondent UUID"
// @Param request body dto.ReviewEditRequest true "Review Data"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Router /api/survey/edit-requests/{id} [put]
func (h *RisikoHandler) ReviewEditRequest(w http.ResponseWriter, r *http.Request) {
	role := strings.ToLower(strings.TrimSpace(middleware.GetRole(r.Context())))
	if role != "admin" {
		utils.RespondError(w, http.StatusForbidden, "Forbidden")
		return
	}

	adminID := middleware.GetUserID(r.Context())
	idStr := strings.TrimPrefix(r.URL.Path, "/api/survey/edit-requests/")
	respondenID := strings.Trim(idStr, "/")
	if respondenID == "" {
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

func decodeOptionalJSON(r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return nil
	}
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(dst); err != nil {
		if errors.Is(err, http.ErrBodyReadAfterClose) {
			return err
		}
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	return nil
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
