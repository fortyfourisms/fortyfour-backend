package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"fortyfour-backend/pkg/logger"
)

// SEExportHandler handles PDF export endpoints for SE data.
type SEExportHandler struct {
	exportService services.SEExportServiceInterface
}

func NewSEExportHandler(exportService services.SEExportServiceInterface) *SEExportHandler {
	return &SEExportHandler{exportService: exportService}
}

// ServeHTTP routes export requests:
//
//	GET /api/se/export-pdf              → export semua (admin) atau milik perusahaan sendiri (user)
//	GET /api/se/export-pdf?id_perusahaan=xxx → admin: filter perusahaan tertentu | user: diabaikan
//	GET /api/se/{id}/export-pdf         → export satu SE by ID
func (h *SEExportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Trim prefix and detect sub-path
	// Possible paths arriving here:
	//   /api/se/export-pdf
	//   /api/se/{id}/export-pdf
	path := strings.TrimPrefix(r.URL.Path, "/api/se/")
	path = strings.TrimSuffix(path, "/")

	if path == "export-pdf" {
		// List export
		h.handleExportAll(w, r)
		return
	}

	// Pattern: {id}/export-pdf
	if strings.HasSuffix(path, "/export-pdf") {
		id := strings.TrimSuffix(path, "/export-pdf")
		if id == "" {
			utils.RespondError(w, 400, "ID tidak valid")
			return
		}
		h.handleExportByID(w, r, id)
		return
	}

	utils.RespondError(w, 404, "Route tidak ditemukan")
}

// handleExportAll
// @Summary      Export semua SE ke PDF
// @Description  Admin: export semua SE, atau filter by id_perusahaan. User: hanya milik perusahaannya.
// @Tags         SE Export
// @Produce      application/pdf
// @Param        id_perusahaan query string false "Filter by ID Perusahaan (admin only)"
// @Success      200 {file} binary
// @Failure      403 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/se/export-pdf [get]
func (h *SEExportHandler) handleExportAll(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())

	var (
		pdfBytes []byte
		err      error
	)

	if role == "admin" {
		// Admin: cek query param id_perusahaan (opsional)
		idPerusahaan := strings.TrimSpace(r.URL.Query().Get("id_perusahaan"))
		if idPerusahaan != "" {
			// Export filtered by perusahaan tertentu
			pdfBytes, err = h.exportService.ExportByPerusahaanPDF(idPerusahaan)
		} else {
			// Export semua data
			pdfBytes, err = h.exportService.ExportAllPDF()
		}
	} else {
		// User biasa: paksa pakai id_perusahaan dari JWT, query param diabaikan
		idPerusahaan := middleware.GetIDPerusahaan(r.Context())
		if idPerusahaan == "" {
			utils.RespondError(w, 403, "Akun Anda belum terhubung ke perusahaan")
			return
		}
		pdfBytes, err = h.exportService.ExportByPerusahaanPDF(idPerusahaan)
	}

	if err != nil {
		logger.Error(err, "gagal generate PDF SE")
		utils.RespondError(w, 500, err.Error())
		return
	}

	servePDF(w, pdfBytes, "laporan-se.pdf")
}

// handleExportByID
// @Summary      Export satu SE ke PDF
// @Description  Export data SE berdasarkan ID. User hanya bisa akses SE milik perusahaannya.
// @Tags         SE Export
// @Produce      application/pdf
// @Param        id path string true "SE ID"
// @Success      200 {file} binary
// @Failure      403 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/se/{id}/export-pdf [get]
func (h *SEExportHandler) handleExportByID(w http.ResponseWriter, r *http.Request, id string) {
	se, pdfBytes, err := h.exportService.ExportByIDPDF(id)
	if err != nil {
		if err.Error() == "data tidak ditemukan" {
			utils.RespondError(w, 404, err.Error())
			return
		}
		logger.Error(err, "gagal generate PDF SE by ID")
		utils.RespondError(w, 500, err.Error())
		return
	}

	// Validasi ownership untuk user
	role := middleware.GetRole(r.Context())
	if role != "admin" {
		idPerusahaan := middleware.GetIDPerusahaan(r.Context())
		if se.IDPerusahaan != idPerusahaan {
			utils.RespondError(w, 403, "Anda tidak memiliki akses ke data ini")
			return
		}
	}

	filename := fmt.Sprintf("laporan-se-%s.pdf", id)
	servePDF(w, pdfBytes, filename)
}

// servePDF writes PDF bytes to the response with appropriate headers.
func servePDF(w http.ResponseWriter, data []byte, filename string) {
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// ════════════════════════════════════════════════════════════════════════════
// CSIRT Export Handler
// ════════════════════════════════════════════════════════════════════════════

// CsirtExportHandler handles PDF export endpoints for CSIRT data.
type CsirtExportHandler struct {
	exportService services.CsirtExportServiceInterface
}

func NewCsirtExportHandler(exportService services.CsirtExportServiceInterface) *CsirtExportHandler {
	return &CsirtExportHandler{exportService: exportService}
}

// ServeHTTP routes CSIRT export requests:
//
//	GET /api/csirt/export-pdf                        → semua (admin) atau milik perusahaan sendiri (user)
//	GET /api/csirt/export-pdf?id_perusahaan=xxx      → admin: filter perusahaan tertentu
//	GET /api/csirt/{id}/export-pdf                   → export satu CSIRT by ID
func (h *CsirtExportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/csirt/")
	path = strings.TrimSuffix(path, "/")

	if path == "export-pdf" {
		h.handleExportAll(w, r)
		return
	}

	if strings.HasSuffix(path, "/export-pdf") {
		id := strings.TrimSuffix(path, "/export-pdf")
		if id == "" {
			utils.RespondError(w, 400, "ID tidak valid")
			return
		}
		h.handleExportByID(w, r, id)
		return
	}

	utils.RespondError(w, 404, "Route tidak ditemukan")
}

// handleExportAll
// @Summary      Export semua CSIRT ke PDF
// @Description  Admin: export semua CSIRT, atau filter by id_perusahaan. User: hanya milik perusahaannya.
// @Tags         CSIRT Export
// @Produce      application/pdf
// @Param        id_perusahaan query string false "Filter by ID Perusahaan (admin only)"
// @Success      200 {file} binary
// @Failure      403 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/csirt/export-pdf [get]
func (h *CsirtExportHandler) handleExportAll(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())

	var (
		pdfBytes []byte
		err      error
	)

	if role == "admin" {
		idPerusahaan := strings.TrimSpace(r.URL.Query().Get("id_perusahaan"))
		if idPerusahaan != "" {
			pdfBytes, err = h.exportService.ExportByPerusahaanPDF(idPerusahaan)
		} else {
			pdfBytes, err = h.exportService.ExportAllPDF()
		}
	} else {
		idPerusahaan := middleware.GetIDPerusahaan(r.Context())
		if idPerusahaan == "" {
			utils.RespondError(w, 403, "Akun Anda belum terhubung ke perusahaan")
			return
		}
		pdfBytes, err = h.exportService.ExportByPerusahaanPDF(idPerusahaan)
	}

	if err != nil {
		logger.Error(err, "gagal generate PDF CSIRT")
		utils.RespondError(w, 500, err.Error())
		return
	}

	servePDF(w, pdfBytes, "laporan-csirt.pdf")
}

// handleExportByID
// @Summary      Export satu CSIRT ke PDF
// @Description  Export data CSIRT berdasarkan ID. User hanya bisa akses CSIRT milik perusahaannya.
// @Tags         CSIRT Export
// @Produce      application/pdf
// @Param        id path string true "CSIRT ID"
// @Success      200 {file} binary
// @Failure      403 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/csirt/{id}/export-pdf [get]
func (h *CsirtExportHandler) handleExportByID(w http.ResponseWriter, r *http.Request, id string) {
	csirt, pdfBytes, err := h.exportService.ExportByIDPDF(id)
	if err != nil {
		if err.Error() == "data tidak ditemukan" {
			utils.RespondError(w, 404, err.Error())
			return
		}
		logger.Error(err, "gagal generate PDF CSIRT by ID")
		utils.RespondError(w, 500, err.Error())
		return
	}

	// Validasi ownership untuk user
	role := middleware.GetRole(r.Context())
	if role != "admin" {
		idPerusahaan := middleware.GetIDPerusahaan(r.Context())
		if csirt.Perusahaan.ID != idPerusahaan {
			utils.RespondError(w, 403, "Anda tidak memiliki akses ke data ini")
			return
		}
	}

	filename := fmt.Sprintf("laporan-csirt-%s.pdf", id)
	servePDF(w, pdfBytes, filename)
}