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

	GetAllEditRequests() ([]dto.EditRequestItemResponse, error)
	GetMyEditRequest(userID string) (*dto.EditRequestItemResponse, error)
}

// HANDLER STRUCT
type RisikoHandler struct {
	svc RisikoServiceInterface
}

func NewRisikoHandler(svc RisikoServiceInterface) *RisikoHandler {
	return &RisikoHandler{svc: svc}
}

// @Summary Get All Risiko Aktif
// @Description Get all active risks for survey
// @Tags Risiko
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/survey/risiko [get]
func (h *RisikoHandler) GetAllRisiko(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetAllRisiko()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Berhasil mengambil data risiko", data)
}

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
	if !middleware.HasRole(r.Context(), "user_pic") {
		utils.RespondError(w, http.StatusForbidden, "Hanya user_pic yang dapat mengisi survey")
		return
	}
	userID := middleware.GetUserID(r.Context())

	var req dto.EligibilityRequest
	if err := decodeJSON(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	res, err := h.svc.ProcessEligibility(userID, req)
	handleResult(w, res, err)
}

// @Summary Submit Alasan
// @Description Submit reason for not having risk
// @Tags Risiko
// @Accept json
// @Produce json
// @Param request body dto.AlasanRequest true "Alasan request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/survey/risiko/reason [post]
func (h *RisikoHandler) SubmitAlasan(w http.ResponseWriter, r *http.Request) {
	if !middleware.HasRole(r.Context(), "user_pic") {
		utils.RespondError(w, http.StatusForbidden, "Hanya user_pic yang dapat mengisi survey")
		return
	}
	userID := middleware.GetUserID(r.Context())

	var req dto.AlasanRequest
	if err := decodeJSON(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	res, err := h.svc.ProcessAlasan(userID, req)
	handleResult(w, res, err)
}

// @Summary Submit Dampak
// @Description Submit impact of risk
// @Tags Risiko
// @Accept json
// @Produce json
// @Param request body dto.DampakRequest true "Dampak request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/survey/risiko/dampak [post]
func (h *RisikoHandler) SubmitDampak(w http.ResponseWriter, r *http.Request) {
	if !middleware.HasRole(r.Context(), "user_pic") {
		utils.RespondError(w, http.StatusForbidden, "Hanya user_pic yang dapat mengisi survey")
		return
	}
	userID := middleware.GetUserID(r.Context())

	var req dto.DampakRequest
	if err := decodeJSON(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	res, err := h.svc.ProcessDampak(userID, req)
	handleResult(w, res, err)
}

// @Summary Submit Pengendalian
// @Description Submit mitigation/control for risk
// @Tags Risiko
// @Accept json
// @Produce json
// @Param request body dto.PengendalianRequest true "Pengendalian request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/survey/risiko/pengendalian [post]
func (h *RisikoHandler) SubmitPengendalian(w http.ResponseWriter, r *http.Request) {
	if !middleware.HasRole(r.Context(), "user_pic") {
		utils.RespondError(w, http.StatusForbidden, "Hanya user_pic yang dapat mengisi survey")
		return
	}
	userID := middleware.GetUserID(r.Context())

	var req dto.PengendalianRequest
	if err := decodeJSON(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	res, err := h.svc.ProcessPengendalian(userID, req)
	handleResult(w, res, err)
}

// @Summary Get My Risiko Data
// @Description Get current user's risk assessment data
// @Tags Risiko
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Failure 404 {object} dto.ErrorResponse
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

// @Summary Get Risiko Data by Respondent ID
// @Description Get risk assessment data for a specific respondent (Admin/Staff only)
// @Tags Risiko
// @Produce json
// @Param id path string true "Respondent ID"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/survey/risiko/{id} [get]
func (h *RisikoHandler) GetByRespondentID(w http.ResponseWriter, r *http.Request) {
	if !middleware.HasRole(r.Context(), "admin", "staff") {
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

// @Summary Get Survey Progress
// @Description Get current survey progress for the user
// @Tags Survey Progress
// @Produce json
// @Success 200 {object} dto.APIResponse{data=dto.ProgressResponse}
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/survey/progress [get]
func (h *RisikoHandler) GetProgress(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	res, err := h.svc.GetProgress(userID)
	handleResult(w, res, err)
}

// @Summary Navigate Survey
// @Description Navigate to next or previous risk
// @Tags Survey Progress
// @Accept json
// @Produce json
// @Param request body dto.NavigateRequest true "Navigation request"
// @Success 200 {object} dto.APIResponse{data=dto.ProgressResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/survey/navigate [post]
func (h *RisikoHandler) Navigate(w http.ResponseWriter, r *http.Request) {
	if !middleware.HasRole(r.Context(), "user_pic") {
		utils.RespondError(w, http.StatusForbidden, "Hanya user_pic yang dapat mengisi survey")
		return
	}
	userID := middleware.GetUserID(r.Context())

	var req dto.NavigateRequest
	if err := decodeOptionalJSON(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	res, err := h.svc.Navigate(userID, req)
	handleResult(w, res, err)
}

// @Summary Save Survey Progress
// @Description Save current risk progress
// @Tags Survey Progress
// @Accept json
// @Produce json
// @Param request body dto.NavigateRequest true "Save progress request"
// @Success 200 {object} dto.APIResponse{data=dto.ProgressResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/survey/save-progress [post]
func (h *RisikoHandler) SaveProgress(w http.ResponseWriter, r *http.Request) {
	if !middleware.HasRole(r.Context(), "user_pic") {
		utils.RespondError(w, http.StatusForbidden, "Hanya user_pic yang dapat mengisi survey")
		return
	}
	userID := middleware.GetUserID(r.Context())

	var req dto.NavigateRequest
	if err := decodeOptionalJSON(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	res, err := h.svc.SaveProgress(userID, req)
	handleResult(w, res, err)
}

// @Summary Finish Survey
// @Description Mark survey as completed
// @Tags Survey Progress
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/survey/finish [post]
func (h *RisikoHandler) FinishSurvey(w http.ResponseWriter, r *http.Request) {
	if !middleware.HasRole(r.Context(), "user_pic") {
		utils.RespondError(w, http.StatusForbidden, "Hanya user_pic yang dapat mengisi survey")
		return
	}
	userID := middleware.GetUserID(r.Context())

	if err := h.svc.FinishSurvey(userID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Survey berhasil diselesaikan", nil)
}

// @Summary Request Edit
// @Description Request to edit survey data
// @Tags Survey Edit Request
// @Accept json
// @Produce json
// @Param request body dto.RequestEditRequest true "Edit request reason"
// @Success 200 {object} dto.APIResponse{data=dto.ProgressResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/survey/request-edit [post]
func (h *RisikoHandler) RequestEdit(w http.ResponseWriter, r *http.Request) {
	if !middleware.HasRole(r.Context(), "user_pic") {
		utils.RespondError(w, http.StatusForbidden, "Hanya user_pic yang dapat mengajukan edit request")
		return
	}
	userID := middleware.GetUserID(r.Context())

	var req dto.RequestEditRequest
	if err := decodeJSON(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	res, err := h.svc.RequestEdit(userID, req)
	handleResult(w, res, err)
}

// @Summary Review Edit Request
// @Description Approve or reject a request to edit survey data. Admin specifies action ("approve" or "reject").
// @Tags Survey Edit Request
// @Accept json
// @Produce json
// @Param id path string true "Respondent ID"
// @Param request body dto.ReviewEditRequest true "Provide 'action' (approve/reject) and optional 'response' string"
// @Success 200 {object} dto.APIResponse{data=dto.ProgressResponse}
// @Failure 400 {object} dto.ErrorResponse "Invalid body or respondent ID"
// @Failure 403 {object} dto.ErrorResponse "Forbidden - Admin or Staff only"
// @Router /api/survey/edit-requests/{id} [post]
func (h *RisikoHandler) ReviewEditRequest(w http.ResponseWriter, r *http.Request) {
	if !middleware.HasRole(r.Context(), "admin") {
		utils.RespondError(w, http.StatusForbidden, "Hanya Admin yang dapat memproses edit request")
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

// @Summary Get Edit Requests
// @Description Get survey edit requests. Admin/Staff see all, User sees their own.
// @Tags Survey Edit Request
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=[]dto.EditRequestItemResponse}
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/survey/edit-requests [get]
func (h *RisikoHandler) GetEditRequests(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	if middleware.HasRole(r.Context(), "admin", "staff") {
		res, err := h.svc.GetAllEditRequests()
		handleResult(w, res, err)
		return
	}

	// For User/UserPIC, get only their own
	res, err := h.svc.GetMyEditRequest(userID)
	if err != nil {
		handleResult(w, nil, err)
		return
	}

	// Wrap single item in array for consistency if needed,
	// but the DTO return for MyEditRequest is a single object.
	// Let's keep it as is or wrap it.
	utils.RespondSuccess(w, http.StatusOK, "Success", []interface{}{res})
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
