package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"survey/internal/dto"
	"survey/internal/middleware"
	"survey/internal/utils"
)

// SERVICE
type RespondenServiceInterface interface {
	GetAll() ([]dto.RespondenResponse, error)
	GetByID(id string) (*dto.RespondenResponse, error)
	GetByUserID(userID string) (*dto.RespondenResponse, error)
	UpsertByUserID(userID string, req dto.CreateRespondenRequest) (*dto.RespondenResponse, error)
}

// HANDLER
type RespondenHandler struct {
	service RespondenServiceInterface
}

func NewRespondenHandler(service RespondenServiceInterface) *RespondenHandler {
	return &RespondenHandler{service: service}
}

// ROUTER
func (h *RespondenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/survey/responden")
	path = strings.Trim(path, "/")

	switch r.Method {
	case http.MethodGet:
		if path == "me" {
			if !middleware.HasRole(r.Context(), "user", "user_pic", "admin", "staff") {
				utils.RespondError(w, http.StatusForbidden, "Forbidden")
				return
			}
			h.GetMe(w, r)
			return
		}

		if path == "" {
			if !middleware.HasRole(r.Context(), "admin", "staff") {
				utils.RespondError(w, http.StatusForbidden, "Forbidden")
				return
			}
			h.handleGetAll(w)
			return
		}

		if !middleware.HasRole(r.Context(), "admin", "staff") {
			utils.RespondError(w, http.StatusForbidden, "Forbidden")
			return
		}
		h.handleGetByID(w, path)

	case http.MethodPost:
		if path != "me" && path != "" {
			utils.RespondError(w, http.StatusForbidden, "Hanya /me atau base path yang diizinkan")
			return
		}

		if !middleware.HasRole(r.Context(), "user_pic") {
			utils.RespondError(w, http.StatusForbidden, "Hanya user_pic yang dapat mengisi data responden")
			return
		}

		h.UpsertMe(w, r)

	default:
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// @Summary Get All Responden
// @Description Get all survey respondents (Admin/Staff only)
// @Tags Responden
// @Produce json
// @Success 200 {object} dto.APIResponse{data=[]dto.RespondenResponse}
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/survey/responden [get]
func (h *RespondenHandler) handleGetAll(w http.ResponseWriter) {
	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Berhasil mengambil data responden", data)
}

// @Summary Get Responden by ID
// @Description Get details of a specific respondent (Admin/Staff only)
// @Tags Responden
// @Produce json
// @Param id path string true "Responden ID"
// @Success 200 {object} dto.APIResponse{data=dto.RespondenResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/survey/responden/{id} [get]
func (h *RespondenHandler) handleGetByID(w http.ResponseWriter, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Berhasil mengambil detail responden", data)
}

// @Summary Get My Responden Profile
// @Description Get current user's respondent profile
// @Tags Responden
// @Produce json
// @Success 200 {object} dto.APIResponse{data=dto.RespondenResponse}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/survey/responden/me [get]
func (h *RespondenHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	data, err := h.service.GetByUserID(userID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Berhasil mengambil data profil survey", data)
}

// @Summary Upsert My Responden Profile
// @Description Create or update current user's respondent profile
// @Tags Responden
// @Accept json
// @Produce json
// @Param request body dto.CreateRespondenRequest true "Respondent data"
// @Success 200 {object} dto.APIResponse{data=dto.RespondenResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/survey/responden/me [post]
func (h *RespondenHandler) UpsertMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.CreateRespondenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid body")
		return
	}
	if strings.TrimSpace(req.IdPerusahaan) == "" {
		req.IdPerusahaan = middleware.GetPerusahaanID(r.Context())
	}

	resp, err := h.service.UpsertByUserID(userID, req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Berhasil memperbarui data profil survey", resp)
}
