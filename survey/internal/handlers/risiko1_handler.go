package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"survey/internal/dto"
	"survey/internal/repository"
	"survey/internal/services"
)

type RisikoHandler struct {
	svc *services.RisikoService
}

func NewRisikoHandler(svc *services.RisikoService) *RisikoHandler {
	return &RisikoHandler{svc: svc}
}

// STEP 1: ELIGIBILITY
// SubmitEligibility godoc
// @Summary      Step 1 - Submit Eligibility
// @Description  Menentukan apakah responden memenuhi kriteria risiko
// @Tags         Risiko
// @Accept       json
// @Produce      json
// @Param        request body dto.EligibilityRequest true "Eligibility Request"
// @Success      200 {object} dto.EligibilityResponse
// @Failure      400 {object} dto.ErrorResponse
// @Router       /api/survey/risiko/eligibility [post]
func (h *RisikoHandler) SubmitEligibility(w http.ResponseWriter, r *http.Request) {
	var req dto.EligibilityRequest

	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.svc.ProcessEligibility(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// STEP 2A: ALASAN (JIKA TIDAK)
// SubmitAlasan godoc
// @Summary      Step 2A - Submit Alasan
// @Description  Mengisi alasan jika tidak memenuhi kriteria risiko
// @Tags         Risiko
// @Accept       json
// @Produce      json
// @Param        request body dto.AlasanRequest true "Alasan Request"
// @Success      200 {object} dto.AlasanResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Router       /api/survey/risiko/alasan [post]
func (h *RisikoHandler) SubmitAlasan(w http.ResponseWriter, r *http.Request) {
	var req dto.AlasanRequest

	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.svc.ProcessAlasan(req)
	if err != nil {
		writeError(w, resolveErrorStatus(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// STEP 2B: DAMPAK (JIKA YA)
// SubmitDampak godoc
// @Summary      Step 2B - Submit Dampak
// @Description  Mengisi dampak jika memenuhi kriteria risiko
// @Tags         Risiko
// @Accept       json
// @Produce      json
// @Param        request body dto.DampakRequest true "Dampak Request"
// @Success      200 {object} dto.DampakResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Router       /api/survey/risiko/dampak [post]
func (h *RisikoHandler) SubmitDampak(w http.ResponseWriter, r *http.Request) {
	var req dto.DampakRequest

	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.svc.ProcessDampak(req)
	if err != nil {
		writeError(w, resolveErrorStatus(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// STEP 2C: PENGENDALIAN
// SubmitPengendalian godoc
// @Summary      Step 2C - Submit Pengendalian
// @Description  Mengisi tindakan pengendalian risiko
// @Tags         Risiko
// @Accept       json
// @Produce      json
// @Param        request body dto.PengendalianRequest true "Pengendalian Request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Router       /api/survey/risiko/pengendalian [post]
func (h *RisikoHandler) SubmitPengendalian(w http.ResponseWriter, r *http.Request) {
	var req dto.PengendalianRequest

	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.svc.ProcessPengendalian(req)
	if err != nil {
		writeError(w, resolveErrorStatus(err), err.Error())
		return
	}

	msg := "Tindakan pengendalian berhasil disimpan"
	if !req.AdaPengendalian {
		msg = "Tidak ada pengendalian, risiko selesai"
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message":   msg,
		"data":      result, 
		"next_step": "finish",
	})
}

// GET BY RESPONDEN ID
// GetByRespondentID godoc
// @Summary      Get Risiko by Respondent ID
// @Description  Mengambil data risiko berdasarkan responden_id
// @Tags         Risiko
// @Produce      json
// @Param        responden_id path int true "Respondent ID"
// @Success      200 {object} dto.RisikoResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Router       /api/survey/risiko/{responden_id} [get]
func (h *RisikoHandler) GetByRespondentID(w http.ResponseWriter, r *http.Request) {

	// ambil dari URL: /api/survey/risiko/
	path := strings.TrimPrefix(r.URL.Path, "/api/survey/risiko/")
	
	if path == "" {
		writeError(w, http.StatusBadRequest, "respondent_id diperlukan")
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "respondent_id harus angka")
		return
	}

	result, err := h.svc.GetByRespondentID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "data tidak ditemukan")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GET PROGRESS
// GetProgress godoc
// @Summary      Get Progress Risiko
// @Description  Mengambil progress pengisian risiko berdasarkan responden_id
// @Tags         Risiko
// @Produce      json
// @Param        responden_id path int true "Respondent ID"
// @Success      200 {object} dto.ProgressResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/survey/progress/{responden_id} [get]
func (h *RisikoHandler) GetProgress(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/api/survey/progress/")

	if path == "" {
		writeError(w, http.StatusBadRequest, "respondent_id diperlukan")
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "respondent_id harus angka")
		return
	}

	result, err := h.svc.GetProgress(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// NAVIGATE
// Navigate godoc
// @Summary      Navigate Step Risiko
// @Description  Navigasi antar step pengisian risiko
// @Tags         Risiko
// @Accept       json
// @Produce      json
// @Param        request body dto.NavigateRequest true "Navigate Request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/survey/risiko/navigate [post]
func (h *RisikoHandler) Navigate(w http.ResponseWriter, r *http.Request) {
	var req dto.NavigateRequest

	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.svc.Navigate(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
	"success": true,
	"data": result,
	})
}

// HELPERS
func decodeJSON(r *http.Request, dst interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]interface{}{
		"success": false,
		"message": message,
	})
}

func resolveErrorStatus(err error) int {
	if errors.Is(err, repository.ErrNotFound) {
		return http.StatusNotFound
	}
	return http.StatusBadRequest
}

// LANJUTKAN NANTI
// SaveProgress godoc
// @Summary      Simpan Progress Risiko (Lanjutkan Nanti)
// @Description  Menyimpan progress sementara untuk dilanjutkan nanti
// @Tags         Risiko
// @Accept       json
// @Produce      json
// @Param        request body dto.NavigateRequest true "Save Progress Request"
// @Success      200 {object} dto.ProgressResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/survey/risiko/save-progress [post]
func (h *RisikoHandler) SaveProgress(w http.ResponseWriter, r *http.Request) {
	var req dto.NavigateRequest

	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.svc.SaveProgress(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}