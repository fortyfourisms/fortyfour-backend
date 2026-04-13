package handlers

import (
	"encoding/json"
	"ikas/internal/dto"
	"ikas/internal/middleware"
	"ikas/internal/services"
	"ikas/internal/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/rollbar/rollbar-go"
)

type JawabanDeteksiHandler struct {
	service *services.JawabanDeteksiService
}

func NewJawabanDeteksiHandler(service *services.JawabanDeteksiService) *JawabanDeteksiHandler {
	return &JawabanDeteksiHandler{service: service}
}

func (h *JawabanDeteksiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method

	switch {
	case method == http.MethodPost && path == "/api/maturity/jawaban-deteksi":
		h.handleCreate(w, r)
	case method == http.MethodGet && path == "/api/maturity/jawaban-deteksi":
		// Implicitly filter by PerusahaanID if user is not admin
		userRole, _ := r.Context().Value(middleware.Role).(string)
		userPerusahaanID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)

		perusahaanID := r.URL.Query().Get("perusahaan_id")

		if userRole != "admin" {
			if userPerusahaanID == "" || userPerusahaanID == "null" {
				utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
					"message": "Berhasil mengambil data",
					"data":    []dto.JawabanDeteksiResponse{},
					"total":   0,
				})
				return
			}
			// Override for non-admins
			values := r.URL.Query()
			values.Set("perusahaan_id", userPerusahaanID)
			r.URL.RawQuery = values.Encode()
		} else if perusahaanID != "" {
			// Leave as is for admin filtering
		}

		h.handleGetAll(w, r)
	case method == http.MethodGet && strings.HasPrefix(path, "/api/maturity/jawaban-deteksi/"):
		h.handleGetByID(w, r)
	case method == http.MethodPut && strings.HasPrefix(path, "/api/maturity/jawaban-deteksi/"):
		h.handleUpdate(w, r)
	case method == http.MethodDelete && strings.HasPrefix(path, "/api/maturity/jawaban-deteksi/"):
		h.handleDelete(w, r)
	default:
		utils.RespondError(w, http.StatusNotFound, "Endpoint tidak ditemukan")
	}
}

// @Summary Create Jawaban Deteksi
// @Description Create a new answer for detection question
// @Tags Jawaban Deteksi
// @Accept json
// @Produce json
// @Param request body dto.CreateJawabanDeteksiRequest true "Jawaban Deteksi Request"
// @Success 201 {object} dto.JawabanDeteksiResponse
// @Router /api/maturity/jawaban-deteksi [post]
func (h *JawabanDeteksiHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateJawabanDeteksiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	userRole := ""
	if val := r.Context().Value(middleware.Role); val != nil {
		userRole = val.(string)
	}

	msg, err := h.service.Create(req, userRole)
	if err != nil {
		rollbar.Error(err)
		switch err.Error() {
		case "pertanyaan_deteksi_id tidak valid",
			"ikas_id tidak boleh kosong",
			"format ikas_id tidak valid",
			"jawaban_deteksi harus bernilai antara 0 sampai 5, atau null untuk N/A",
			"validasi hanya boleh diisi jika evidence ada",
			"validasi hanya boleh berisi 'yes' atau 'no'":
			utils.RespondError(w, 400, err.Error())
		case "pertanyaan_deteksi_id tidak ditemukan",
			"ikas_id tidak ditemukan":
			utils.RespondError(w, 404, err.Error())
		case "pertanyaan ini sudah pernah diisi untuk asesmen ini":
			utils.RespondError(w, 409, err.Error())
		default:
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 201, map[string]interface{}{
		"message": msg,
	})
}

// @Summary Get All Jawaban Deteksi
// @Description Get all answers for detection questions, optionally filtered by ikas_id or pertanyaan_deteksi_id
// @Tags Jawaban Deteksi
// @Produce json
// @Param ikas_id query string false "Filter by Ikas ID"
// @Param pertanyaan_deteksi_id query int false "Filter by Pertanyaan Deteksi ID"
// @Success 200 {array} dto.JawabanDeteksiResponse
// @Router /api/maturity/jawaban-deteksi [get]
func (h *JawabanDeteksiHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	ikasID := r.URL.Query().Get("ikas_id")
	pertanyaanIDStr := r.URL.Query().Get("pertanyaan_deteksi_id")

	var data []dto.JawabanDeteksiResponse
	var err error

	if ikasID != "" {
		h.handleGetByIkasID(w, r, ikasID)
		return
	} else if pertanyaanIDStr != "" {
		pID, _ := strconv.Atoi(pertanyaanIDStr)
		data, err = h.service.GetByPertanyaan(pID)
	} else {
		data, err = h.service.GetAll()
	}

	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Berhasil mengambil data",
		"data":    data,
		"total":   len(data),
	})
}

func (h *JawabanDeteksiHandler) handleGetByIkasID(w http.ResponseWriter, _ *http.Request, ikasID string) {
	data, err := h.service.GetByIkasID(ikasID)
	if err != nil {
		rollbar.Error(err)
		if err.Error() == "format ikas_id tidak valid" {
			utils.RespondError(w, 400, err.Error())
		} else {
			utils.RespondError(w, 500, err.Error())
		}
		return
	}
	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil mengambil data",
		"data":    data,
		"total":   len(data),
	})
}

// @Summary Get Jawaban Deteksi by ID
// @Description Get a specific detection answer by its ID
// @Tags Jawaban Deteksi
// @Produce json
// @Param id path int true "Jawaban Deteksi ID"
// @Success 200 {object} dto.JawabanDeteksiResponse
// @Router /api/maturity/jawaban-deteksi/{id} [get]
func (h *JawabanDeteksiHandler) handleGetByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractIntID(r.URL.Path, "jawaban-deteksi")
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	userRole, _ := r.Context().Value(middleware.Role).(string)
	userIkasID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)

	resp, err := h.service.GetByID(id, userRole, userIkasID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "data tidak ditemukan" {
			status = http.StatusNotFound
		}
		utils.RespondError(w, status, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Berhasil mengambil data",
		"data":    resp,
	})
}

// @Summary Update Jawaban Deteksi
// @Description Update an existing detection answer
// @Tags Jawaban Deteksi
// @Accept json
// @Produce json
// @Param id path int true "Jawaban Deteksi ID"
// @Param request body dto.UpdateJawabanDeteksiRequest true "Update Request"
// @Success 200 {object} dto.JawabanDeteksiResponse
// @Router /api/maturity/jawaban-deteksi/{id} [put]
func (h *JawabanDeteksiHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractIntID(r.URL.Path, "jawaban-deteksi")
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	var req dto.UpdateJawabanDeteksiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID := ""
	if val := r.Context().Value(middleware.UserIDKey); val != nil {
		userID = val.(string)
	}

	userRole := ""
	if val := r.Context().Value(middleware.Role); val != nil {
		userRole = val.(string)
	}

	userIkasID := ""
	if val := r.Context().Value(middleware.PerusahaanIDKey); val != nil {
		userIkasID = val.(string)
	}

	err = h.service.Update(id, req, userID, userRole, userIkasID)
	if err != nil {
		rollbar.Error(err)
		switch err.Error() {
		case "data tidak ditemukan":
			utils.RespondError(w, 404, err.Error())
		case "format ID tidak valid",
			"jawaban_deteksi harus bernilai antara 0 sampai 5, atau null untuk N/A",
			"validasi hanya boleh berisi 'yes' atau 'no'",
			"validasi hanya boleh diisi jika evidence ada":
			utils.RespondError(w, 400, err.Error())
		default:
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil menyimpan data",
		"id":      id,
	})
}

// @Summary Delete Jawaban Deteksi
// @Description Delete a specific detection answer
// @Tags Jawaban Deteksi
// @Produce json
// @Param id path int true "Jawaban Deteksi ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/maturity/jawaban-deteksi/{id} [delete]
func (h *JawabanDeteksiHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractIntID(r.URL.Path, "jawaban-deteksi")
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	userID := ""
	if val := r.Context().Value(middleware.UserIDKey); val != nil {
		userID = val.(string)
	}

	userRole, _ := r.Context().Value(middleware.Role).(string)
	userIkasID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)

	if err := h.service.Delete(id, userID, userRole, userIkasID); err != nil {
		rollbar.Error(err)
		if err.Error() == "data tidak ditemukan" {
			utils.RespondError(w, 404, err.Error())
		} else {
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil menghapus data",
		"id":      id,
	})
}
