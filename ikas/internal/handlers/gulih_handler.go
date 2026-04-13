package handlers

import (
	_ "ikas/internal/dto"
	"ikas/internal/services"
	"ikas/internal/utils"
	"net/http"

	"fortyfour-backend/pkg/logger"
	"ikas/internal/middleware"
	"ikas/internal/models"
	"strings"
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
func (h *GulihHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	userRole, _ := r.Context().Value(middleware.Role).(string)
	userPerusahaanID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)

	ikasID := r.URL.Query().Get("ikas_id")

	if userRole != "admin" && (userPerusahaanID == "" || userPerusahaanID == "null") {
		utils.RespondJSON(w, 200, map[string]interface{}{
			"message": "Berhasil mengambil data",
			"data":    []interface{}{},
			"total":   0,
		})
		return
	}

	var data interface{}
	var err error

	if ikasID != "" {
		data, err = h.service.GetByIkasID(ikasID, userRole, userPerusahaanID)
	} else {
		if userRole != "admin" {
			data, err = h.service.GetByPerusahaanID(userPerusahaanID, userRole, userPerusahaanID)
		} else {
			data, err = h.service.GetAll(userRole)
		}
	}

	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 500, err.Error())
		return
	}

	total := 0
	if data != nil {
		switch v := data.(type) {
		case []models.Gulih:
			total = len(v)
		case *models.Gulih:
			total = 1
		}
	}

	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil mengambil data",
		"data":    data,
		"total":   total,
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
func (h *GulihHandler) handleGetByID(w http.ResponseWriter, r *http.Request, id string) {
	userRole, _ := r.Context().Value(middleware.Role).(string)
	userPerusahaanID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)

	data, err := h.service.GetByID(id, userRole, userPerusahaanID)
	if err != nil {
		logger.Error(err, "operation failed")
		if strings.Contains(err.Error(), "tidak memiliki akses") {
			utils.RespondError(w, 403, err.Error())
		} else {
			utils.RespondError(w, 404, "Data tidak ditemukan")
		}
		return
	}
	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil mengambil data",
		"data":    data,
	})
}
