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

	resp, err := h.service.Create(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "pertanyaan_gulih_id tidak ditemukan" ||
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
