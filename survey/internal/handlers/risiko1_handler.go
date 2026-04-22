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

// STEP 1 — ELIGIBILITY
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

	writeSuccess(w, result)
}

// STEP 2A — ALASAN
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

	writeSuccess(w, result)
}

// STEP 2B — DAMPAK
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

	writeSuccess(w, result)
}

// STEP 2C — PENGENDALIAN
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

	writeSuccess(w, result)
}

// GET RISIKO BY RESPONDENT
func (h *RisikoHandler) GetByRespondentID(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/api/survey/risiko/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	result, err := h.svc.GetByRespondentID(id) // ⬅ pastikan service ADA
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "data tidak ditemukan")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, result)
}

// GET PROGRESS
func (h *RisikoHandler) GetProgress(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/api/survey/progress/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	progress, err := h.svc.GetProgress(id) // ⬅ HARUS ADA DI SERVICE
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, progress)
}

// NAVIGATE
func (h *RisikoHandler) Navigate(w http.ResponseWriter, r *http.Request) {
	var req dto.NavigateRequest

	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	result, err := h.svc.Navigate(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccess(w, result)
}

// SAVE PROGRESS
func (h *RisikoHandler) SaveProgress(w http.ResponseWriter, r *http.Request) {
	var req dto.NavigateRequest

	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	result, err := h.svc.SaveProgress(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, result)
}

// CUSTOM RISIKO
func (h *RisikoHandler) CreateCustomRisiko(w http.ResponseWriter, r *http.Request) {
	var req dto.CustomRisikoRequest

	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	id, err := h.svc.CreateCustomRisiko(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccess(w, map[string]int{
		"custom_risiko_id": id,
	})
}

// FINISH
func (h *RisikoHandler) FinishSurvey(w http.ResponseWriter, r *http.Request) {

	var req struct {
		RespondenID int `json:"responden_id"`
	}

	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if err := h.svc.FinishSurvey(req.RespondenID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, "survey selesai")
}

// HELPERS
func decodeJSON(r *http.Request, dst interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

func writeSuccess(w http.ResponseWriter, data interface{}) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]interface{}{
		"success": false,
		"message": message,
	})
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func resolveErrorStatus(err error) int {
	if errors.Is(err, repository.ErrNotFound) {
		return http.StatusNotFound
	}
	return http.StatusBadRequest
}