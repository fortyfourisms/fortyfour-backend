package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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
		"next_step": "finish",
	})
}

// GET BY RESPONDEN ID
func (h *RisikoHandler) GetByRespondentID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("respondent_id")

	id, err := strconv.Atoi(idStr)
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
func (h *RisikoHandler) GetProgress(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("respondent_id")

	id, err := strconv.Atoi(idStr)
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

	writeJSON(w, http.StatusOK, result)
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