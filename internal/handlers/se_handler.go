package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
)

type SEHandler struct {
	service    services.SEService
	sseService services.SSEServiceInterface
}

func NewSEHandler(
	service services.SEService,
	sseService services.SSEServiceInterface,
) *SEHandler {
	return &SEHandler{
		service:    service,
		sseService: sseService,
	}
}

func (h *SEHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/api/se"), "/")

	switch r.Method {
	case http.MethodGet:
		if id == "" {
			h.handleGetAll(w)
		} else {
			h.handleGetByID(w, id)
		}
	case http.MethodPost:
		h.handleCreate(w, r)
	case http.MethodPut:
		h.handleUpdate(w, r, id)
	case http.MethodDelete:
		h.handleDelete(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// @Summary Get all SE
// @Description Get all sistem elektronik with kategorisasi
// @Tags SE
// @Accept json
// @Produce json
// @Success 200 {object} dto.SEListResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/se [get]
func (h *SEHandler) handleGetAll(w http.ResponseWriter) {
	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// @Summary Get SE by ID
// @Description Get sistem elektronik by ID
// @Tags SE
// @Accept json
// @Produce json
// @Param id path string true "SE ID"
// @Success 200 {object} dto.SEResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /api/se/{id} [get]
func (h *SEHandler) handleGetByID(w http.ResponseWriter, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}
	utils.RespondJSON(w, 200, data)
}

// @Summary Create SE
// @Description Create new sistem elektronik with kategorisasi
// @Tags SE
// @Accept json
// @Produce json
// @Param request body dto.SECreateRequest true "SE Create Request"
// @Success 201 {object} dto.SEResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /api/se [post]
func (h *SEHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSERequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyCreate("se", resp, userID)

	utils.RespondJSON(w, 201, resp)
}

// @Summary Update SE
// @Description Update sistem elektronik
// @Tags SE
// @Accept json
// @Produce json
// @Param id path string true "SE ID"
// @Param request body dto.SEUpdateRequest true "SE Update Request"
// @Success 200 {object} dto.SEResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /api/se/{id} [put]
func (h *SEHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req dto.UpdateSERequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyUpdate("se", resp, userID)

	utils.RespondJSON(w, 200, resp)
}

// @Summary Delete SE
// @Description Delete sistem elektronik
// @Tags SE
// @Accept json
// @Produce json
// @Param id path string true "SE ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} utils.ErrorResponse
// @Router /api/se/{id} [delete]
func (h *SEHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.Delete(id); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyDelete("se", id, userID)

	utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}
