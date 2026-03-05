package handlers

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"net/http"
	"strings"
)

// Ensure dto is available for swagger type resolution.
var _ dto.SubSektorResponse

type SubSektorHandler struct {
	service services.SubSektorServiceInterface
}

func NewSubSektorHandler(service services.SubSektorServiceInterface) *SubSektorHandler {
	return &SubSektorHandler{service: service}
}

func (h *SubSektorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/sub_sektor")

	// Check for /by-sektor/:id route
	if strings.HasPrefix(path, "/by_sektor/") {
		sektorID := strings.TrimPrefix(path, "/by_sektor/")
		if r.Method == http.MethodGet {
			h.handleGetBySektorID(w, r, sektorID)
			return
		}
	}

	id := strings.TrimPrefix(path, "/")

	switch r.Method {
	case http.MethodGet:
		if id == "" {
			h.handleGetAll(w, r)
		} else {
			h.handleGetByID(w, r, id)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// GetAllSubSektor godoc
// @Summary      List semua sub sektor
// @Description  Mengambil seluruh data sub sektor
// @Tags         SubSektor
// @Produce      json
// @Success      200  {array}  dto.SubSektorResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /api/sub_sektor [get]
func (h *SubSektorHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetSubSektorByID godoc
// @Summary      Ambil sub sektor berdasarkan ID
// @Description  Mengambil satu data sub sektor
// @Tags         SubSektor
// @Produce      json
// @Param        id   path      string  true  "SubSektor ID"
// @Success      200  {object} dto.SubSektorResponse
// @Failure      404  {object} dto.ErrorResponse
// @Router       /api/sub_sektor/{id} [get]
func (h *SubSektorHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetSubSektorBySektorID godoc
// @Summary      Ambil sub sektor berdasarkan Sektor ID
// @Description  Mengambil data sub sektor dalam satu sektor
// @Tags         SubSektor
// @Produce      json
// @Param        id   path      string  true  "Sektor ID"
// @Success      200  {array}  dto.SubSektorResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /api/sub_sektor/by_sektor/{id} [get]
func (h *SubSektorHandler) handleGetBySektorID(w http.ResponseWriter, _ *http.Request, sektorID string) {
	data, err := h.service.GetBySektorID(sektorID)
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}