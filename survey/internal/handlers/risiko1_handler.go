package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"survey/internal/dto"
	"survey/internal/repository"
)

type RisikoServiceInterface interface {
	ProcessEligibility(req dto.EligibilityRequest) (map[string]interface{}, error)
	ProcessAlasan(req dto.AlasanRequest) (map[string]interface{}, error)
	ProcessDampak(req dto.DampakRequest) (map[string]interface{}, error)
	ProcessPengendalian(req dto.PengendalianRequest) (map[string]interface{}, error)
	GetByRespondentID(id int) (map[string]interface{}, error)
	GetProgress(id int) (dto.ProgressResponse, error)
	Navigate(req dto.NavigateRequest) (dto.ProgressResponse, error)
	SaveProgress(req dto.NavigateRequest) (dto.ProgressResponse, error)
	CreateCustomRisiko(req dto.CustomRisikoRequest) (int, error)
	FinishSurvey(respondenID int) error
}

type RisikoHandler struct {
	svc RisikoServiceInterface
}

func NewRisikoHandler(svc RisikoServiceInterface) *RisikoHandler {
	return &RisikoHandler{svc: svc}
}

// STEP 1 — ELIGIBILITY
// SubmitEligibility godoc
// @Summary      Step 1 - Eligibility
// @Description  Menentukan apakah responden masuk kategori risiko
// @Tags         Risiko
// @Accept       json
// @Produce      json
// @Param        request body dto.EligibilityRequest true "Eligibility Request"
// @Success      200 {object} map[string]interface{} "Success response"
// @Failure      400 {object} dto.ErrorResponse "Invalid request"
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

	writeSuccess(w, result)
}

// STEP 2A — ALASAN
// SubmitAlasan godoc
// @Summary      Step 2A - Alasan
// @Description  Mengisi alasan jika tidak eligible
// @Tags         Risiko
// @Accept       json
// @Produce      json
// @Param        request body dto.AlasanRequest true "Alasan Request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Router       /api/survey/risiko/reason [post]
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
// SubmitDampak godoc
// @Summary      Step 2B - Dampak
// @Description  Mengisi dampak jika eligible
// @Tags         Risiko
// @Accept       json
// @Produce      json
// @Param        request body dto.DampakRequest true "Dampak Request"
// @Success      200 {object} map[string]interface{}
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

	writeSuccess(w, result)
}

// STEP 2C — PENGENDALIAN
// SubmitPengendalian godoc
// @Summary      Step 2C - Pengendalian
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

	writeSuccess(w, result)
}

// GET RISIKO BY RESPONDENT
// GetByRespondentID godoc
// @Summary      Get Risiko by Responden
// @Description  Mengambil data risiko berdasarkan responden_id
// @Tags         Risiko
// @Produce      json
// @Param        responden_id path int true "Responden ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} dto.ErrorResponse "Invalid ID"
// @Failure      404 {object} dto.ErrorResponse "Data tidak ditemukan"
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/survey/risiko/{responden_id} [get]
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
// GetProgress godoc
// @Summary      Get Progress Risiko
// @Description  Mengambil progress pengisian survey
// @Tags         Risiko
// @Produce      json
// @Param        responden_id path int true "Responden ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/survey/progress/{responden_id} [get]
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
// Navigate godoc
// @Summary      Navigasi Step Risiko
// @Description  Mengatur alur step survey (next/back)
// @Tags         Risiko
// @Accept       json
// @Produce      json
// @Param        request body dto.NavigateRequest true "Navigate Request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} dto.ErrorResponse
// @Router       /api/survey/risiko/navigate [post]
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
// SaveProgress godoc
// @Summary      Simpan Progress
// @Description  Menyimpan progress agar bisa dilanjutkan nanti
// @Tags         Risiko
// @Accept       json
// @Produce      json
// @Param        request body dto.NavigateRequest true "Save Progress Request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/survey/risiko/save-progress [post]
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
// CreateCustomRisiko godoc
// @Summary      Tambah Risiko Custom
// @Description  Menambahkan risiko baru di luar daftar sistem
// @Tags         Risiko
// @Accept       json
// @Produce      json
// @Param        request body dto.CustomRisikoRequest true "Custom Risiko Request"
// @Success      200 {object} map[string]int
// @Failure      400 {object} dto.ErrorResponse
// @Router       /api/survey/risiko/custom-risk [post]
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
// FinishSurvey godoc
// @Summary      Selesaikan Survey
// @Description  Menandai survey telah selesai
// @Tags         Risiko
// @Accept       json
// @Produce      json
// @Param        request body object{responden_id=int} true "Finish Request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/survey/risiko/finish [post]
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
