package handlers

import (
	"ikas/internal/services"
	"ikas/internal/utils"
	"net/http"

	"fortyfour-backend/pkg/logger"
)

type GulihHandler struct {
	service *services.GulihService
}

func NewGulihHandler(service *services.GulihService) *GulihHandler {
	return &GulihHandler{
		service: service,
	}
}

func (h *GulihHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "gulih")

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

// GetAllGulih godoc
// @Summary      List semua gulih
// @Description  Mengambil seluruh data gulih
// @Tags         Gulih
// @Produce      json
// @Success      200  {array}  dto.GulihResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /api/gulih [get]
func (h *GulihHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
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

// GetGulihByID godoc
// @Summary      Ambil gulih berdasarkan ID
// @Tags         Gulih
// @Produce      json
// @Param        id   path      string  true  "Gulih ID"
// @Success      200  {object} dto.GulihResponse
// @Failure      404  {object} dto.ErrorResponse
// @Router       /api/gulih/{id} [get]
func (h *GulihHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
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
