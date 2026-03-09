package handlers

import (
	"encoding/json"
	"ikas/internal/dto"
	"ikas/internal/services"
	"ikas/internal/utils"
	"net/http"
	"strconv"
	"strings"
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

	resp, err := h.service.Create(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "pertanyaan_deteksi_id tidak ditemukan" ||
			err.Error() == "perusahaan_id tidak ditemukan" ||
			err.Error() == "jawaban untuk pertanyaan ini sudah ada untuk perusahaan tersebut" {
			status = http.StatusBadRequest
		}
		utils.RespondError(w, status, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Jawaban berhasil disimpan",
		"data":    resp,
	})
}

// @Summary Get All Jawaban Deteksi
// @Description Get all answers for detection questions, optionally filtered by perusahaan_id or pertanyaan_deteksi_id
// @Tags Jawaban Deteksi
// @Produce json
// @Param perusahaan_id query string false "Filter by Perusahaan ID"
// @Param pertanyaan_deteksi_id query int false "Filter by Pertanyaan Deteksi ID"
// @Success 200 {array} dto.JawabanDeteksiResponse
// @Router /api/maturity/jawaban-deteksi [get]
func (h *JawabanDeteksiHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	perusahaanID := r.URL.Query().Get("perusahaan_id")
	pertanyaanIDStr := r.URL.Query().Get("pertanyaan_deteksi_id")

	var data []dto.JawabanDeteksiResponse
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

	resp, err := h.service.Update(id, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "data tidak ditemukan" {
			status = http.StatusNotFound
		}
		utils.RespondError(w, status, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Berhasil memperbarui data",
		"data":    resp,
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

	if err := h.service.Delete(id); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "data tidak ditemukan" {
			status = http.StatusNotFound
		}
		utils.RespondError(w, status, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Berhasil menghapus data",
	})
}
