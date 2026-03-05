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

// Summary godoc
// @Summary      Get dashboard summary
// @Description  Mengambil ringkasan data dashboard. Mendukung berbagai filter opsional.
// @Description  Prioritas filter tanggal: from+to > year+quarter > year.
// @Tags         Dashboard
// @Security     BearerAuth
// @Produce      json
// @Param        from          query  string  false  "Start date (YYYY-MM-DD)"
// @Param        to            query  string  false  "End date (YYYY-MM-DD)"
// @Param        year          query  string  false  "Filter per tahun, misal 2025"
// @Param        quarter       query  string  false  "Filter per kuartal (1-4), harus digunakan bersama year"
// @Param        sub_sektor_id query  string  false  "Filter per sub-sektor (UUID)"
// @Param        kategori_se   query  string  false  "Filter kategori SE: Strategis | Tinggi | Rendah"
// @Success      200  {object}  dto.DashboardSummary
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/dashboard/summary [get]
func (h *DashboardHandler) Summary(w http.ResponseWriter, r *http.Request) {
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
			writeError(w, http.StatusBadRequest, "kategori_se tidak valid, nilai yang diizinkan: Strategis, Tinggi, Rendah")
			return
		}
		f.KategoriSE = &kategoriSE
	}

	res, err := h.svc.GetSummary(r.Context(), f)
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
	if path == "/api/dashboard/summary" || path == "/api/dashboard/summary/" {
		h.Summary(w, r)
		return
	}

	http.NotFound(w, r)
}
