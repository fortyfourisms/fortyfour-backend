package handlers

import (
	_ "fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"net/http"
)

type KonversiHandler struct {
	service services.KonversiServiceInterface
}

func NewKonversiHandler(service services.KonversiServiceInterface) *KonversiHandler {
	return &KonversiHandler{service: service}
}

// ServeHTTP handles GET /api/konversi and GET /api/konversi/{id}
//
//	@Summary		Get Konversi Data
//	@Description	Mengambil data poin partisipasi (IKAS, KSE, Survey, CSIRT) untuk perusahaan.
//	@Tags			Konversi
//	@Produce		json
//	@Param			id				path		string	false	"ID Perusahaan (UUID) jika melalui path /api/konversi/{id}"
//	@Param			perusahaan_id	query		string	false	"ID Perusahaan (UUID) jika melalui query /api/konversi?perusahaan_id="
//	@Success		200				{object}	utils.JSONResponse{data=[]dto.KonversiResponse}
//	@Router			/api/konversi [get]
//	@Router			/api/konversi/{id} [get]
func (h *KonversiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 1. Coba ambil dari Query Param
	perusahaanID := r.URL.Query().Get("perusahaan_id")

	// 2. Jika query kosong, coba ambil dari Path Param
	if perusahaanID == "" {
		// Mengambil string setelah "/api/konversi/"
		path := r.URL.Path
		if len(path) > len("/api/konversi/") {
			perusahaanID = path[len("/api/konversi/"):]
		}
	}

	data, err := h.service.GetKonversi(perusahaanID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.JSONResponse{
		Status:  "success",
		Message: "Berhasil mengambil data konversi",
		Total:   len(data),
		Data:    data,
	})
}
