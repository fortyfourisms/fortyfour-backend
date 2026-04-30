package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
)

// validKategoriSE adalah nilai yang diizinkan untuk filter kategori_se
var validKategoriSE = map[string]bool{
	"Strategis": true,
	"Tinggi":    true,
	"Rendah":    true,
}

// reYear mencocokkan format YYYY, misal "2025"
var reYear = regexp.MustCompile(`^\d{4}$`)

// reQuarter mencocokkan nilai "1", "2", "3", "4"
var reQuarter = regexp.MustCompile(`^[1-4]$`)

type DashboardHandler struct {
	svc *services.DashboardService
}

func NewDashboardHandler(svc *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// ptrStr mengembalikan pointer ke string jika tidak kosong, nil jika kosong
func ptrStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// parseFilter mengekstrak filter dari query parameter
func (h *DashboardHandler) parseFilter(r *http.Request) (dto.DashboardFilter, bool, string) {
	q := r.URL.Query()
	f := dto.DashboardFilter{}

	// --- from & to ---
	from := q.Get("from")
	to := q.Get("to")
	if from != "" && to != "" {
		if _, err := time.Parse("2006-01-02", from); err == nil {
			if _, err2 := time.Parse("2006-01-02", to); err2 == nil {
				f.From = &from
				f.To = &to
			}
		}
	}

	// --- year ---
	year := q.Get("year")
	if year != "" {
		if reYear.MatchString(year) {
			f.Year = &year
		}
		// year tidak valid → diabaikan
	}

	// --- quarter (hanya valid bila year juga ada) ---
	quarter := q.Get("quarter")
	if quarter != "" && f.Year != nil {
		if reQuarter.MatchString(quarter) {
			f.Quarter = &quarter
		}
	}

	// --- sub_sektor_id ---
	f.SubSektorID = ptrStr(q.Get("sub_sektor_id"))

	// --- kategori_se ---
	kategoriSE := q.Get("kategori_se")
	if kategoriSE != "" {
		if !validKategoriSE[kategoriSE] {
			return f, false, "kategori_se tidak valid, nilai yang diizinkan: Strategis, Tinggi, Rendah"
		}
		f.KategoriSE = &kategoriSE
	}

	return f, true, ""
}

// SummarySektor godoc
//
//	@Summary		Get dashboard sektor
//	@Description	Mengambil data jumlah perusahaan per sektor.
//	@Tags			Dashboard
//	@Security		BearerAuth
//	@Produce		json
//	@Param			from			query		string	false	"Start date (YYYY-MM-DD)"
//	@Param			to				query		string	false	"End date (YYYY-MM-DD)"
//	@Param			year			query		string	false	"Filter per tahun, misal 2025"
//	@Param			quarter			query		string	false	"Filter per kuartal (1-4), harus digunakan bersama year"
//	@Param			sub_sektor_id	query		string	false	"Filter per sub-sektor (UUID)"
//	@Success		200				{object}	dto.DashboardSektorResponse
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/api/dashboard/sektor [get]
func (h *DashboardHandler) SummarySektor(w http.ResponseWriter, r *http.Request) {
	f, ok, errMsg := h.parseFilter(r)
	if !ok {
		writeError(w, http.StatusBadRequest, errMsg)
		return
	}
	res, err := h.svc.GetSummarySektor(r.Context(), f)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// SummaryIkas godoc
//
//	@Summary		Get dashboard IKAS
//	@Description	Mengambil data agregasi IKAS dan status pengisian IKAS.
//	@Tags			Dashboard
//	@Security		BearerAuth
//	@Produce		json
//	@Param			from			query		string	false	"Start date (YYYY-MM-DD)"
//	@Param			to				query		string	false	"End date (YYYY-MM-DD)"
//	@Param			year			query		string	false	"Filter per tahun, misal 2025"
//	@Param			quarter			query		string	false	"Filter per kuartal (1-4), harus digunakan bersama year"
//	@Param			sub_sektor_id	query		string	false	"Filter per sub-sektor (UUID)"
//	@Success		200				{object}	dto.DashboardIkasResponse
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/api/dashboard/ikas [get]
func (h *DashboardHandler) SummaryIkas(w http.ResponseWriter, r *http.Request) {
	f, ok, errMsg := h.parseFilter(r)
	if !ok {
		writeError(w, http.StatusBadRequest, errMsg)
		return
	}
	res, err := h.svc.GetSummaryIkas(r.Context(), f)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// SummarySE godoc
//
//	@Summary		Get dashboard SE
//	@Description	Mengambil data agregasi SE dan status pengisian SE.
//	@Tags			Dashboard
//	@Security		BearerAuth
//	@Produce		json
//	@Param			from			query		string	false	"Start date (YYYY-MM-DD)"
//	@Param			to				query		string	false	"End date (YYYY-MM-DD)"
//	@Param			year			query		string	false	"Filter per tahun, misal 2025"
//	@Param			quarter			query		string	false	"Filter per kuartal (1-4), harus digunakan bersama year"
//	@Param			sub_sektor_id	query		string	false	"Filter per sub-sektor (UUID)"
//	@Param			kategori_se		query		string	false	"Filter kategori SE: Strategis | Tinggi | Rendah"
//	@Success		200				{object}	dto.DashboardSEResponse
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/api/dashboard/se [get]
func (h *DashboardHandler) SummarySE(w http.ResponseWriter, r *http.Request) {
	f, ok, errMsg := h.parseFilter(r)
	if !ok {
		writeError(w, http.StatusBadRequest, errMsg)
		return
	}
	res, err := h.svc.GetSummarySE(r.Context(), f)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// SummaryCSIRT godoc
//
//	@Summary		Get dashboard CSIRT
//	@Description	Mengambil data agregasi CSIRT dan status pembentukan CSIRT.
//	@Tags			Dashboard
//	@Security		BearerAuth
//	@Produce		json
//	@Param			from			query		string	false	"Start date (YYYY-MM-DD)"
//	@Param			to				query		string	false	"End date (YYYY-MM-DD)"
//	@Param			year			query		string	false	"Filter per tahun, misal 2025"
//	@Param			quarter			query		string	false	"Filter per kuartal (1-4), harus digunakan bersama year"
//	@Param			sub_sektor_id	query		string	false	"Filter per sub-sektor (UUID)"
//	@Success		200				{object}	dto.DashboardCSIRTResponse
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/api/dashboard/csirt [get]
func (h *DashboardHandler) SummaryCSIRT(w http.ResponseWriter, r *http.Request) {
	f, ok, errMsg := h.parseFilter(r)
	if !ok {
		writeError(w, http.StatusBadRequest, errMsg)
		return
	}
	res, err := h.svc.GetSummaryCSIRT(r.Context(), f)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// ServeHTTP so DashboardHandler implements http.Handler.
func (h *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	path := r.URL.Path
	switch path {
	case "/api/dashboard/sektor", "/api/dashboard/sektor/":
		h.SummarySektor(w, r)
	case "/api/dashboard/ikas", "/api/dashboard/ikas/":
		h.SummaryIkas(w, r)
	case "/api/dashboard/se", "/api/dashboard/se/":
		h.SummarySE(w, r)
	case "/api/dashboard/csirt", "/api/dashboard/csirt/":
		h.SummaryCSIRT(w, r)
	default:
		http.NotFound(w, r)
	}
}
