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

// handleGetAll handles GET /api/survey/responden
func (h *RespondenHandler) handleGetAll(w http.ResponseWriter) {
	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Berhasil mengambil data responden", data)
}

// handleGetByID handles GET /api/survey/responden/{id}
func (h *RespondenHandler) handleGetByID(w http.ResponseWriter, id string) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID harus angka")
		return
	}

	data, err := h.service.GetByID(idInt)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Berhasil mengambil detail responden", data)
}

// GetMe handles GET /api/survey/responden/me
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

// UpsertMe handles POST /api/survey/responden/me
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

	resp, err := h.service.UpsertByUserID(userID, req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Berhasil memperbarui data profil survey", resp)
}
