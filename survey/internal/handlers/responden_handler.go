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
	role := strings.ToLower(strings.TrimSpace(middleware.GetRole(r.Context())))
	path := strings.TrimPrefix(r.URL.Path, "/api/survey/responden")
	path = strings.Trim(path, "/")

	switch r.Method {
	case http.MethodGet:
		if path == "me" {
			if role != "user" && role != "user_pic" {
				utils.RespondError(w, http.StatusForbidden, "Forbidden")
				return
			}
			h.GetMe(w, r)
			return
		}

		if path == "" {
			if role != "admin" && role != "staff" {
				utils.RespondError(w, http.StatusForbidden, "Forbidden")
				return
			}
			h.handleGetAll(w)
			return
		}

		if role != "admin" && role != "staff" {
			utils.RespondError(w, http.StatusForbidden, "Forbidden")
			return
		}
		h.handleGetByID(w, path)

	case http.MethodPost:
		if path != "me" {
			utils.RespondError(w, http.StatusForbidden, "Only /me allowed")
			return
		}

		if role != "user" && role != "user_pic" {
			utils.RespondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		h.UpsertMe(w, r)

	default:
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handleGetAll godoc
// @Summary Get all responden (Admin/Staff)
// @Description Retrieve a list of all survey respondents.
// @Tags Survey Responden
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Berhasil mengambil data responden"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/survey/responden [get]
func (h *RespondenHandler) handleGetAll(w http.ResponseWriter) {
	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Berhasil mengambil data responden", data)
}

// handleGetByID godoc
// @Summary Get responden by ID (Admin/Staff)
// @Description Retrieve specific respondent details by their UUID.
// @Tags Survey Responden
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Responden UUID"
// @Success 200 {object} map[string]interface{} "Berhasil mengambil detail responden"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Data tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/survey/responden/{id} [get]
func (h *RespondenHandler) handleGetByID(w http.ResponseWriter, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			utils.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Berhasil mengambil detail responden", data)
}

// GetMe godoc
// @Summary Get my respondent profile (User/User PIC)
// @Description Retrieve the respondent profile of the currently authenticated user.
// @Tags Survey Responden
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Berhasil mengambil data profil survey"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Data responden belum tersedia"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/survey/responden/me [get]
func (h *RespondenHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	data, err := h.service.GetByUserID(userID)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			utils.RespondError(w, http.StatusNotFound, "Data responden belum tersedia, silakan lengkapi profil")
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Berhasil mengambil data profil survey", data)
}

// UpsertMe godoc
// @Summary Update or create my respondent profile (User/User PIC)
// @Description Creates or updates the respondent profile of the currently authenticated user.
// @Tags Survey Responden
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateRespondenRequest true "Responden Data"
// @Success 200 {object} map[string]interface{} "Berhasil memperbarui data profil survey"
// @Failure 400 {object} map[string]interface{} "Invalid body or Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
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
