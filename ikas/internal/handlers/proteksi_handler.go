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

type ProteksiHandler struct {
	service *services.ProteksiService
}

func NewProteksiHandler(service *services.ProteksiService) *ProteksiHandler {
	return &ProteksiHandler{
		service: service,
	}
}

func (h *ProteksiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "proteksi")

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

// GetAllProteksi godoc
// @Summary      List semua proteksi
// @Description  Mengambil seluruh data proteksi
// @Tags         Proteksi
// @Produce      json
// @Success      200  {array}  dto.ProteksiResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /api/proteksi [get]
func (h *ProteksiHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
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
		data, err = h.service.GetByIkasID(ikasID)
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
		case []models.Proteksi:
			total = len(v)
		case *models.Proteksi:
			total = 1
		}
	}

	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil mengambil data",
		"data":    data,
		"total":   total,
	})
}

// GetProteksiByID godoc
// @Summary      Ambil proteksi berdasarkan ID
// @Description  Mengambil satu data proteksi
// @Tags         Proteksi
// @Produce      json
// @Param        id   path      string  true  "Proteksi ID"
// @Success      200  {object} dto.ProteksiResponse
// @Failure      404  {object} dto.ErrorResponse
// @Router       /api/proteksi/{id} [get]
func (h *ProteksiHandler) handleGetByID(w http.ResponseWriter, r *http.Request, id string) {
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
