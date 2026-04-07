package handlers

import (
	"ikas/internal/services"
	"ikas/internal/utils"
	"ikas/internal/dto"
	"net/http"

	"fortyfour-backend/pkg/logger"
)

type IdentifikasiHandler struct {
	service *services.IdentifikasiService
}

func NewIdentifikasiHandler(service *services.IdentifikasiService) *IdentifikasiHandler {
	return &IdentifikasiHandler{
		service: service,
	}
}

func (h *IdentifikasiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "identifikasi")

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

// GetAllIdentifikasi godoc
// @Summary      List semua identifikasi
// @Description  Mengambil seluruh data identifikasi
// @Tags         Identifikasi
// @Produce      json
// @Success      200  {array}  dto.IdentifikasiResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /api/identifikasi [get]
func (h *IdentifikasiHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 500, err.Error())
		return
	}

	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil mengambil data",
		"data":    data,
		"total":   len(data),
	})
}

// GetIdentifikasiByID godoc
// @Summary      Ambil identifikasi berdasarkan ID
// @Description  Mengambil satu data identifikasi
// @Tags         Identifikasi
// @Produce      json
// @Param        id   path      string  true  "Identifikasi ID"
// @Success      200  {object} dto.IdentifikasiResponse
// @Failure      404  {object} dto.ErrorResponse
// @Router       /api/identifikasi/{id} [get]
func (h *IdentifikasiHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}
	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil mengambil data",
		"data":    data,
	})
}
