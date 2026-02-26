package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"fortyfour-backend/internal/services"
)

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

// Summary godoc
// @Summary      Get dashboard summary
// @Description  Mengambil ringkasan data dashboard (jumlah perusahaan per sektor dan total SE). Filter tanggal opsional dengan format YYYY-MM-DD; keduanya harus diisi agar filter aktif.
// @Tags         Dashboard
// @Security     BearerAuth
// @Produce      json
// @Param        from  query  string  false  "Start date (YYYY-MM-DD)"
// @Param        to    query  string  false  "End date (YYYY-MM-DD)"
// @Success      200  {object}  dto.DashboardSummary
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/dashboard/summary [get]
func (h *DashboardHandler) Summary(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	from := q.Get("from")
	to := q.Get("to")
	var fromPtr, toPtr *string
	if from != "" && to != "" {
		// validate format YYYY-MM-DD, ignore if invalid
		if _, err := time.Parse("2006-01-02", from); err == nil {
			if _, err2 := time.Parse("2006-01-02", to); err2 == nil {
				fromPtr = &from
				toPtr = &to
			}
		}
	}
	res, err := h.svc.GetSummary(r.Context(), fromPtr, toPtr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// ServeHTTP so DashboardHandler implements http.Handler.
// It only accepts GET and routes the summary path; other dashboard routes are intentionally not supported.
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