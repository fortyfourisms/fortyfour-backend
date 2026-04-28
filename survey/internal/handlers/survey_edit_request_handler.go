package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"survey/internal/dto"
	"survey/internal/middleware"
	"survey/internal/services"
	"survey/internal/utils"
)

type EditRequestHandler struct {
	service services.EditRequestService
}

func NewEditRequestHandler(service services.EditRequestService) *EditRequestHandler {
	return &EditRequestHandler{service: service}
}

// ROUTER ENTRY
func (h *EditRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/survey/edit-requests")
	path = strings.Trim(path, "/")

	switch {
	case path == "" && r.Method == http.MethodGet:
		h.handleList(w, r)

	case path == "" && r.Method == http.MethodPost:
		h.handleCreate(w, r)

	case strings.HasSuffix(path, "/review") && r.Method == http.MethodPut:
		id := strings.TrimSuffix(path, "/review")
		id = strings.Trim(id, "/")
		h.handleReview(w, r, id)

	default:
		utils.RespondError(w, http.StatusNotFound, "Endpoint tidak ditemukan")
	}
}

// CREATE EDIT REQUEST
// @Summary		Submit edit request risiko
// @Description	User mengajukan perubahan data risiko
// @Tags			Risiko - Edit Request
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			request	body		dto.CreateEditRequestDTO	true	"Data perubahan"
// @Success		201		{object}	dto.EditRequestResponse
// @Failure		400		{object}	dto.ErrorResponse
// @Failure		401		{object}	dto.ErrorResponse
// @Router			/api/survey/edit-requests [post]
func (h *EditRequestHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.CreateEditRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.service.Create(userID, req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, resp)
}

// LIST EDIT REQUEST
// @Summary		List edit request risiko
// @Description	Admin melihat semua pending, user melihat miliknya
// @Tags			Risiko - Edit Request
// @Produce		json
// @Security		BearerAuth
// @Success		200	{array}	dto.EditRequestResponse
// @Failure		401	{object}	dto.ErrorResponse
// @Failure		500	{object}	dto.ErrorResponse
// @Router			/api/survey/edit-requests [get]
func (h *EditRequestHandler) handleList(w http.ResponseWriter, r *http.Request) {

	role := middleware.GetRole(r.Context())

	// ADMIN
	if role == "admin" {
		data, err := h.service.GetPending()
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.RespondJSON(w, http.StatusOK, data)
		return
	}

	// USER
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	data, err := h.service.GetByUser(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, data)
}

// REVIEW EDIT REQUEST
// @Summary		Review edit request risiko
// @Description	Admin approve atau reject request edit risiko
// @Tags			Risiko - Edit Request
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path		string						true	"ID Edit Request"
// @Param			request	body		dto.ReviewEditRequestDTO	true	"Status dan catatan"
// @Success		200		{object}	dto.EditRequestResponse
// @Failure		400		{object}	dto.ErrorResponse
// @Failure		403		{object}	dto.ErrorResponse
// @Router			/api/survey/edit-requests/{id}/review [put]
func (h *EditRequestHandler) handleReview(w http.ResponseWriter, r *http.Request, id string) {
	defer r.Body.Close()

	role := middleware.GetRole(r.Context())
	if role != "admin" {
		utils.RespondError(w, http.StatusForbidden, "Hanya admin yang dapat mereview")
		return
	}

	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	var req dto.ReviewEditRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.service.Review(id, req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}
