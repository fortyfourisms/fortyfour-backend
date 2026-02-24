package handlers

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"net/http"
	"strings"
)

// Ensure dto is available for swagger type resolution.
var _ dto.SektorResponse

type SektorHandler struct {
	service services.SektorServiceInterface
}

func NewSektorHandler(service services.SektorServiceInterface) *SektorHandler {
	return &SektorHandler{service: service}
}

func (h *SektorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/api/sektor"), "/")

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

// GetAllSektor godoc
// @Summary      List semua sektor
// @Description  Mengambil seluruh data sektor beserta sub sektor
// @Tags         Sektor
// @Produce      json
// @Success      200  {array}  dto.SektorResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /api/sektor [get]
func (h *SektorHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetSektorByID godoc
// @Summary      Ambil sektor berdasarkan ID
// @Description  Mengambil satu data sektor beserta sub sektor
// @Tags         Sektor
// @Produce      json
// @Param        id   path      string  true  "Sektor ID"
// @Success      200  {object} dto.SektorResponse
// @Failure      404  {object} dto.ErrorResponse
// @Router       /api/sektor/{id} [get]
func (h *SektorHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}
	utils.RespondJSON(w, 200, data)
}
