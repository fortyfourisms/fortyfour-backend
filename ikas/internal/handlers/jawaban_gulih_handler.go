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

type JawabanGulihHandler struct {
	service *services.JawabanGulihService
}

func NewJawabanGulihHandler(service *services.JawabanGulihService) *JawabanGulihHandler {
	return &JawabanGulihHandler{service: service}
}

func (h *JawabanGulihHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method

	switch {
	case method == http.MethodPost && path == "/api/maturity/jawaban-gulih":
		h.handleCreate(w, r)
	case method == http.MethodGet && path == "/api/maturity/jawaban-gulih":
		h.handleGetAll(w, r)
	case method == http.MethodGet && strings.HasPrefix(path, "/api/maturity/jawaban-gulih/"):
		h.handleGetByID(w, r)
	case method == http.MethodPut && strings.HasPrefix(path, "/api/maturity/jawaban-gulih/"):
		h.handleUpdate(w, r)
	case method == http.MethodDelete && strings.HasPrefix(path, "/api/maturity/jawaban-gulih/"):
		h.handleDelete(w, r)
	default:
		utils.RespondError(w, http.StatusNotFound, "Endpoint tidak ditemukan")
	}
}

// @Summary Create Jawaban Gulih
// @Description Create a new answer for gulih question
// @Tags Jawaban Gulih
// @Accept json
// @Produce json
// @Param request body dto.CreateJawabanGulihRequest true "Jawaban Gulih Request"
// @Success 201 {object} map[string]interface{}
// @Router /api/maturity/jawaban-gulih [post]
func (h *JawabanGulihHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateJawabanGulihRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	msg, err := h.service.Create(req)
	if err != nil {
		rollbar.Error(err)
		switch err.Error() {
		case "pertanyaan_gulih_id tidak valid",
			"perusahaan_id tidak boleh kosong",
			"format perusahaan_id tidak valid",
			"jawaban_gulih harus bernilai antara 0 sampai 5, atau null untuk N/A",
			"validasi hanya boleh diisi jika evidence ada",
			"validasi hanya boleh berisi 'yes' atau 'no'":
			utils.RespondError(w, 400, err.Error())
		case "pertanyaan_gulih_id tidak ditemukan",
			"perusahaan_id tidak ditemukan":
			utils.RespondError(w, 404, err.Error())
		case "pertanyaan ini sudah pernah diisi oleh perusahaan Anda":
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

// @Summary Get All Jawaban Gulih
// @Description Get all answers for gulih questions, optionally filtered by perusahaan_id or pertanyaan_gulih_id
// @Tags Jawaban Gulih
// @Produce json
// @Param perusahaan_id query string false "Filter by Perusahaan ID"
// @Param pertanyaan_gulih_id query int false "Filter by Pertanyaan Gulih ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/maturity/jawaban-gulih [get]
func (h *JawabanGulihHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	perusahaanID := r.URL.Query().Get("perusahaan_id")
	pertanyaanIDStr := r.URL.Query().Get("pertanyaan_gulih_id")

	var data []dto.JawabanGulihResponse
	var err error

	if perusahaanID != "" {
		data, err = h.service.GetByPerusahaan(perusahaanID)
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

// @Summary Get Jawaban Gulih by ID
// @Description Get a specific gulih answer by its ID
// @Tags Jawaban Gulih
// @Produce json
// @Param id path int true "Jawaban Gulih ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/maturity/jawaban-gulih/{id} [get]
func (h *JawabanGulihHandler) handleGetByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractIntID(r.URL.Path, "jawaban-gulih")
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	resp, err := h.service.GetByID(id)
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

// @Summary Update Jawaban Gulih
// @Description Update an existing gulih answer
// @Tags Jawaban Gulih
// @Accept json
// @Produce json
// @Param id path int true "Jawaban Gulih ID"
// @Param request body dto.UpdateJawabanGulihRequest true "Update Request"
// @Success 200 {object} map[string]interface{}
// @Router /api/maturity/jawaban-gulih/{id} [put]
func (h *JawabanGulihHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractIntID(r.URL.Path, "jawaban-gulih")
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	var req dto.UpdateJawabanGulihRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID := ""
	if val := r.Context().Value(middleware.UserIDKey); val != nil {
		userID = val.(string)
	}

	err = h.service.Update(id, req, userID)
	if err != nil {
		rollbar.Error(err)
		switch err.Error() {
		case "data tidak ditemukan":
			utils.RespondError(w, 404, err.Error())
		case "format ID tidak valid",
			"jawaban_gulih harus bernilai antara 0 sampai 5, atau null untuk N/A",
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

// @Summary Delete Jawaban Gulih
// @Description Delete a specific gulih answer
// @Tags Jawaban Gulih
// @Produce json
// @Param id path int true "Jawaban Gulih ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/maturity/jawaban-gulih/{id} [delete]
func (h *JawabanGulihHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractIntID(r.URL.Path, "jawaban-gulih")
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	userID := ""
	if val := r.Context().Value(middleware.UserIDKey); val != nil {
		userID = val.(string)
	}

	if err := h.service.Delete(id, userID); err != nil {
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
