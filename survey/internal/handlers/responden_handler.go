package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"survey/internal/dto"
	"survey/internal/middleware"
	"survey/internal/utils"
)

// SERVICE
type RespondenServiceInterface interface {
	GetAll() ([]dto.RespondenResponse, error)
	GetByID(id int) (*dto.RespondenResponse, error)

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

	role := middleware.GetRole(r.Context())

	path := strings.TrimPrefix(r.URL.Path, "/api/survey/responden")
	path = strings.Trim(path, "/")

	switch r.Method {

	// GET
	case http.MethodGet:

		// USER: /me
		if path == "me" {
			if role != "user" {
				utils.RespondError(w, http.StatusForbidden, "forbidden")
				return
			}
			h.GetMe(w, r)
			return
		}

		// ADMIN: GET ALL
		if path == "" {
			if role != "admin" {
				utils.RespondError(w, http.StatusForbidden, "forbidden")
				return
			}
			h.handleGetAll(w)
			return
		}

		// ADMIN: GET BY ID
		if role != "admin" {
			utils.RespondError(w, http.StatusForbidden, "forbidden")
			return
		}
		h.handleGetByID(w, path)

	// POST
	case http.MethodPost:

		if path != "me" {
			utils.RespondError(w, http.StatusForbidden, "only /me allowed")
			return
		}

		if role != "user" {
			utils.RespondError(w, http.StatusForbidden, "forbidden")
			return
		}

		h.UpsertMe(w, r)

	default:
		utils.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// ADMIN

// GetAllResponden godoc
// @Summary      Ambil semua responden
// @Description  Hanya admin yang dapat melihat semua data responden
// @Tags         Responden (Admin)
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} dto.RespondenResponse
// @Failure      403 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/survey/responden [get]
func (h *RespondenHandler) handleGetAll(w http.ResponseWriter) {

	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, data)
}

// GetRespondenByID godoc
// @Summary      Ambil responden berdasarkan ID
// @Description  Hanya admin yang dapat melihat detail responden
// @Tags         Responden (Admin)
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Responden ID"
// @Success      200 {object} dto.RespondenResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Router       /api/survey/responden/{id} [get]
func (h *RespondenHandler) handleGetByID(w http.ResponseWriter, id string) {

	idInt, err := strconv.Atoi(id)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "id harus angka")
		return
	}

	data, err := h.service.GetByID(idInt)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, data)
}

// USER

// GetMyResponden godoc
// @Summary      Ambil data responden milik user login
// @Description  User hanya dapat melihat data dirinya sendiri
// @Tags         Responden (User)
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} dto.RespondenResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Router       /api/survey/responden/me [get]
func (h *RespondenHandler) GetMe(w http.ResponseWriter, r *http.Request) {

	userID := middleware.GetUserID(r.Context())

	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	data, err := h.service.GetByUserID(userID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, data)
}

// UpsertMyResponden godoc
// @Summary      Create / Update responden milik user login
// @Description  Jika belum ada maka create, jika sudah ada maka update (upsert)
// @Tags         Responden (User)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateRespondenRequest true "Data Responden"
// @Success      200 {object} dto.RespondenResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Router       /api/survey/responden/me [post]
func (h *RespondenHandler) UpsertMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.CreateRespondenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	resp, err := h.service.UpsertByUserID(userID, req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}
