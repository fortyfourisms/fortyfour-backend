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
//
//	@Summary		List semua deteksi
//	@Description	Mengambil seluruh data deteksi
//	@Tags			Deteksi
//	@Produce		json
//	@Success		200	{array}		dto.DeteksiResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/api/maturity/deteksi [get]
func (h *DeteksiHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	userRole, _ := r.Context().Value(middleware.Role).(string)
	userPerusahaanID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)

	ikasID := r.URL.Query().Get("ikas_id")

	if userRole != "admin" && userRole != "staff" && (userPerusahaanID == "" || userPerusahaanID == "null") {
		utils.RespondListData(w, 200, "Berhasil mengambil data", []interface{}{}, 0)
		return
	}

	var data interface{}
	var err error

	if ikasID != "" {
		data, err = h.service.GetByIkasID(ikasID, userRole, userPerusahaanID)
	} else {
		if userRole != "admin" && userRole != "staff" {
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
		case []models.Deteksi:
			total = len(v)
		case *models.Deteksi:
			total = 1
		}
	}

	utils.RespondListData(w, 200, "Berhasil mengambil data deteksi", data, total)
}

// GetDeteksiByID godoc
//
//	@Summary	Ambil deteksi berdasarkan ID
//	@Tags		Deteksi
//	@Produce	json
//	@Param		id	path		string	true	"Deteksi ID"
//	@Success	200	{object}	dto.DeteksiResponse
//	@Failure	404	{object}	dto.ErrorResponse
//	@Router		/api/maturity/deteksi/{id} [get]
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
	utils.RespondSuccess(w, 200, "Berhasil mengambil data deteksi", data)
}
