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
func (h *DeteksiHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	userRole, _ := r.Context().Value(middleware.Role).(string)
	userPerusahaanID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)

	perusahaanID := r.URL.Query().Get("perusahaan_id")

	if userRole != "admin" {
		if userPerusahaanID == "" || userPerusahaanID == "null" {
			utils.RespondJSON(w, 200, map[string]interface{}{
				"message": "Berhasil mengambil data",
				"data":    []interface{}{},
				"total":   0,
			})
			return
		}
		perusahaanID = userPerusahaanID
	}

	var data interface{}
	var err error

	if perusahaanID != "" {
		data, err = h.service.GetByPerusahaan(perusahaanID)
	} else {
		data, err = h.service.GetAll()
	}

	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 500, err.Error())
		return
	}

	total := 0
	if data != nil {
		switch v := data.(type) {
		case []models.Deteksi:
			total = len(v)
		case *models.Deteksi:
			total = 1
		}
	}

	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil mengambil data",
		"data":    data,
		"total":   total,
	})
}

// GetDeteksiByID godoc
// @Summary      Ambil deteksi berdasarkan ID
// @Tags         Deteksi
// @Produce      json
// @Param        id   path      string  true  "Deteksi ID"
// @Success      200  {object} dto.DeteksiResponse
// @Failure      404  {object} dto.ErrorResponse
// @Router       /api/deteksi/{id} [get]
func (h *DeteksiHandler) handleGetByID(w http.ResponseWriter, r *http.Request, id string) {
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
