package handlers

import (
	"ikas/internal/services"
	"ikas/internal/utils"
	"net/http"

	"fortyfour-backend/pkg/logger"
)

type DeteksiHandler struct {
	service *services.DeteksiService
}

func NewDeteksiHandler(service *services.DeteksiService) *DeteksiHandler {
	return &DeteksiHandler{
		service: service,
	}
}

func (h *DeteksiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "deteksi")

	if r.Method == http.MethodGet {
		if id == "" {
			h.handleGetAll(w, r)
		} else {
			h.handleGetByID(w, r, id)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// GetAllDeteksi godoc
// @Summary      List semua deteksi
// @Description  Mengambil seluruh data deteksi
// @Tags         Deteksi
// @Produce      json
// @Success      200  {array}  dto.DeteksiResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /api/deteksi [get]
func (h *DeteksiHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetDeteksiByID godoc
// @Summary      Ambil deteksi berdasarkan ID
// @Tags         Deteksi
// @Produce      json
// @Param        id   path      string  true  "Deteksi ID"
// @Success      200  {object} dto.DeteksiResponse
// @Failure      404  {object} dto.ErrorResponse
// @Router       /api/deteksi/{id} [get]
func (h *DeteksiHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}
	utils.RespondJSON(w, 200, data)
}
