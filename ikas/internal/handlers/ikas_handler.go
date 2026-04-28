package handlers

import (
	"encoding/json"
	"ikas/internal/dto"
	"ikas/internal/services"
	"ikas/internal/utils"
	"io"
	"net/http"
	"strings"

	"fortyfour-backend/pkg/logger"
	"ikas/internal/middleware"

	"github.com/google/uuid"
)

type IkasHandler struct {
	service *services.IkasService
}

func NewIkasHandler(service *services.IkasService) *IkasHandler {
	return &IkasHandler{
		service: service,
	}
}

func (h *IkasHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	suffix := utils.ExtractID(r.URL.Path, "ikas")

	if suffix == "import" && r.Method == http.MethodPost {
		h.handleImport(w, r)
		return
	}

	if strings.HasSuffix(suffix, "/export") && r.Method == http.MethodGet {
		id := strings.TrimSuffix(suffix, "/export")
		h.handleExportPDF(w, r, id)
		return
	}

	if strings.HasSuffix(suffix, "/request-edit") && r.Method == http.MethodPost {
		id := strings.TrimSuffix(suffix, "/request-edit")
		h.handleRequestEdit(w, r, id)
		return
	}

	if strings.HasSuffix(suffix, "/approve-edit") && r.Method == http.MethodPut {
		id := strings.TrimSuffix(suffix, "/approve-edit")
		h.handleApproveEdit(w, r, id)
		return
	}

	if strings.HasSuffix(suffix, "/reject-edit") && r.Method == http.MethodPut {
		id := strings.TrimSuffix(suffix, "/reject-edit")
		h.handleRejectEdit(w, r, id)
		return
	}

	id := suffix

	switch r.Method {
	case http.MethodGet:
		if id == "" {
			h.handleGetAll(w, r)
		} else {
			h.handleGetByID(w, r, id)
		}
	case http.MethodPost:
		if id != "" {
			utils.RespondError(w, 400, "ID tidak diperlukan untuk create")
			return
		}
		h.handleCreate(w, r)
	case http.MethodPut:
		if id == "" {
			utils.RespondError(w, 400, "ID wajib")
			return
		}
		if strings.HasSuffix(id, "/validate") {
			realID := strings.TrimSuffix(id, "/validate")
			h.handleValidate(w, r, realID)
		} else {
			h.handleUpdate(w, r, id)
		}
	case http.MethodDelete:
		if id == "" {
			utils.RespondError(w, 400, "ID wajib")
			return
		}
		h.handleDelete(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// GetAllIkas godoc
//
//	@Summary		List semua ikas
//	@Description	Mengambil seluruh data ikas
//	@Tags			Ikas
//	@Produce		json
//	@Success		200	{array}		dto.IkasResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/api/maturity/ikas [get]
func (h *IkasHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	perusahaanID := r.URL.Query().Get("perusahaan_id")

	userRole, _ := r.Context().Value(middleware.Role).(string)
	userPerusahaanID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)

	// Implicit filtering for non-admins
	if userRole != "admin" && userRole != "staff" {
		if userPerusahaanID == "" || userPerusahaanID == "null" {
			utils.RespondJSON(w, 200, map[string]interface{}{
				"message": "Berhasil mengambil data",
				"data":    []dto.IkasResponse{},
				"total":   0,
			})
			return
		}
		perusahaanID = userPerusahaanID
	}

	var data []dto.IkasResponse
	var err error

	if perusahaanID != "" && perusahaanID != "null" {
		data, err = h.service.GetByPerusahaan(perusahaanID)
	} else {
		data, err = h.service.GetAll(userRole)
	}

	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 500, err.Error())
		return
	}
	if data == nil {
		data = []dto.IkasResponse{}
	}
	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil mengambil data",
		"data":    data,
		"total":   len(data),
	})
}

// GetIkasByID godoc
//
//	@Summary		Ambil ikas berdasarkan ID
//	@Description	Mengambil satu data ikas
//	@Tags			Ikas
//	@Produce		json
//	@Param			id	path		string	true	"Ikas ID"
//	@Success		200	{object}	dto.IkasResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/api/maturity/ikas/{id} [get]
func (h *IkasHandler) handleGetByID(w http.ResponseWriter, r *http.Request, id string) {
	userRole, _ := r.Context().Value(middleware.Role).(string)
	userPerusahaanID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)

	data, err := h.service.GetByID(id, userRole, userPerusahaanID)
	if err != nil {
		logger.Error(err, "operation failed")
		if strings.Contains(err.Error(), "tidak memiliki akses") {
			utils.RespondError(w, 403, err.Error())
		} else {
			utils.RespondError(w, 404, "Data tidak ditemukan")
		}
		return
	}
	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil mengambil data",
		"data":    data,
	})
}

// CreateIkas godoc
//
//	@Summary		Tambah ikas baru
//	@Description	Membuat record ikas
//	@Tags			Ikas
//	@Accept			json
//	@Produce		json
//	@Param			ikas	body		dto.CreateIkasRequest	true	"Data ikas"
//	@Success		201		{object}	dto.IkasResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Router			/api/maturity/ikas [post]
func (h *IkasHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateIkasRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	newID := uuid.New().String()
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	userRole, _ := r.Context().Value(middleware.Role).(string)
	userPerusahaanID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)

	if err := h.service.Create(r.Context(), req, newID, userID, userRole, userPerusahaanID); err != nil {
		logger.Error(err, "operation failed")
		if strings.Contains(err.Error(), "tidak memiliki akses") || strings.Contains(err.Error(), "belum terasosiasi") {
			utils.RespondError(w, 403, err.Error())
		} else {
			utils.RespondError(w, 400, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 201, map[string]interface{}{
		"message": "Berhasil menyimpan data",
		"id":      newID,
	})
}

// UpdateIkas godoc
//
//	@Summary		Update ikas
//	@Description	Mengubah data ikas berdasarkan ID
//	@Tags			Ikas
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Ikas ID"
//	@Param			ikas	body		dto.UpdateIkasRequest	true	"Data update"
//	@Success		200		{object}	dto.IkasResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Router			/api/maturity/ikas/{id} [put]
func (h *IkasHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req dto.UpdateIkasRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	userRole, _ := r.Context().Value(middleware.Role).(string)
	userPerusahaanID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)

	updatedID, err := h.service.Update(r.Context(), id, req, userID, userRole, userPerusahaanID)
	if err != nil {
		logger.Error(err, "operation failed")
		if strings.Contains(err.Error(), "tidak memiliki akses") {
			utils.RespondError(w, 403, err.Error())
		} else if strings.Contains(err.Error(), "no rows") {
			utils.RespondError(w, 404, "Data tidak ditemukan")
		} else {
			utils.RespondError(w, 400, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil menyimpan data",
		"id":      updatedID,
	})
}

// DeleteIkas godoc
//
//	@Summary		Hapus ikas
//	@Description	Menghapus data ikas berdasarkan ID
//	@Tags			Ikas
//	@Produce		json
//	@Param			id	path		string	true	"Ikas ID"
//	@Success		200	{object}	dto.MessageResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Router			/api/maturity/ikas/{id} [delete]
func (h *IkasHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	userRole, _ := r.Context().Value(middleware.Role).(string)
	userPerusahaanID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)
	if err := h.service.Delete(r.Context(), id, userID, userRole, userPerusahaanID); err != nil {
		logger.Error(err, "operation failed")
		if strings.Contains(err.Error(), "tidak memiliki akses") {
			utils.RespondError(w, 403, err.Error())
		} else if strings.Contains(err.Error(), "no rows") {
			utils.RespondError(w, 404, "Data tidak ditemukan")
		} else {
			utils.RespondError(w, 400, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil menghapus data",
		"id":      id,
	})
}

// ImportIkas godoc
//
//	@Summary		Import IKAS dari Excel
//	@Description	Import data IKAS dari file Excel (sheet ke-7)
//	@Tags			Ikas
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file			formData	file	true	"File Excel (.xlsx)"
//	@Param			id_perusahaan	formData	string	true	"ID Perusahaan"
//	@Param			tanggal			formData	string	true	"Tanggal (YYYY-MM-DD)"
//	@Param			responden		formData	string	true	"Nama Responden"
//	@Param			telepon			formData	string	true	"Nomor Telepon"
//	@Param			jabatan			formData	string	true	"Jabatan"
//	@Success		201				{object}	dto.ImportIkasResponse
//	@Failure		400				{object}	dto.ErrorResponse
//	@Router			/api/maturity/ikas/import [post]
func (h *IkasHandler) handleImport(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, "Gagal parse form data")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, "File 'file' tidak ditemukan")
		return
	}
	defer file.Close()

	if !strings.HasSuffix(strings.ToLower(header.Filename), ".xlsx") {
		utils.RespondError(w, 400, "File harus berformat .xlsx")
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, "Gagal membaca file")
		return
	}

	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	userRole, _ := r.Context().Value(middleware.Role).(string)
	userPerusahaanID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)

	newID, err := h.service.ImportFromExcel(r.Context(), fileBytes, userID, userRole, userPerusahaanID)
	if err != nil {
		logger.Error(err, "operation failed")
		httpStatus := http.StatusBadRequest
		if strings.Contains(err.Error(), "tidak memiliki akses") || strings.Contains(err.Error(), "belum terasosiasi") {
			httpStatus = http.StatusForbidden
		}
		response := struct {
			Success bool     `json:"success"`
			Message string   `json:"message"`
			Errors  []string `json:"errors,omitempty"`
		}{
			Success: false,
			Message: "Import gagal",
			Errors:  []string{err.Error()},
		}
		w.WriteHeader(httpStatus)
		json.NewEncoder(w).Encode(response)
		return
	}

	utils.RespondJSON(w, 201, map[string]interface{}{
		"success": true,
		"message": "Berhasil menyimpan data",
		"id":      newID,
	})
}

// ValidateIkas godoc
//
//	@Summary		Validasi final IKAS
//	@Description	Melakukan validasi final (lock/unlock) data IKAS. Hanya bisa dilakukan oleh Admin atau Staff.
//	@Tags			Ikas
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Ikas ID"
//	@Param			req		body		dto.ValidasiIkasRequest	true	"Status validasi (true=lock, false=unlock)"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		403		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/api/maturity/ikas/{id}/validate [put]
func (h *IkasHandler) handleValidate(w http.ResponseWriter, r *http.Request, id string) {
	userRole, _ := r.Context().Value(middleware.Role).(string)
	if userRole != "admin" && userRole != "staff" {
		utils.RespondError(w, 403, "Hanya admin yang dapat melakukan validasi final")
		return
	}

	var req dto.ValidasiIkasRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Format request tidak valid: "+err.Error())
		return
	}

	if err := h.service.ValidateIkas(r.Context(), id, req.Status); err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}

	msg := "Berhasil melakukan validasi final"
	if !req.Status {
		msg = "Berhasil membuka validasi final"
	}

	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": msg,
	})
}

// ExportIkasPDF godoc
//
//	@Summary		Export IKAS ke PDF
//	@Description	Mengunduh file laporan IKAS dalam format PDF berdasarkan ID
//	@Tags			Ikas
//	@Produce		application/pdf
//	@Param			id	path	string	true	"Ikas ID"
//	@Success		200	{file}	binary
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/api/maturity/ikas/{id}/export [get]
func (h *IkasHandler) handleExportPDF(w http.ResponseWriter, r *http.Request, id string) {
	if !utils.IsValidUUID(id) {
		utils.RespondError(w, 400, "ID tidak valid")
		return
	}

	userRole, _ := r.Context().Value(middleware.Role).(string)
	userPerusahaanID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)

	ikas, pdfBytes, err := h.service.ExportByIDPDF(r.Context(), id, userRole, userPerusahaanID)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			utils.RespondError(w, 404, "Data tidak ditemukan")
		} else {
			utils.RespondError(w, 500, "Gagal generate PDF: "+err.Error())
		}
		return
	}

	// Set headers for PDF download
	filename := "Laporan_IKAS_" + ikas.Perusahaan.NamaPerusahaan + "_" + ikas.Tanggal + ".pdf"
	filename = strings.ReplaceAll(filename, " ", "_")

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.WriteHeader(http.StatusOK)
	w.Write(pdfBytes)
}

// RequestEditIkas godoc
//
//	@Summary		Ajukan permintaan edit IKAS
//	@Description	Mengajukan permintaan pembukaan kunci (unlock) data IKAS yang sudah divalidasi kepada Admin/Staff
//	@Tags			Ikas
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Ikas ID"
//	@Param			req		body		dto.RequestEditRequest	true	"Alasan pengajuan edit"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/api/maturity/ikas/{id}/request-edit [post]
func (h *IkasHandler) handleRequestEdit(w http.ResponseWriter, r *http.Request, id string) {
	if !utils.IsValidUUID(id) {
		utils.RespondError(w, 400, "ID tidak valid")
		return
	}

	var req dto.RequestEditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	userRole, _ := r.Context().Value(middleware.Role).(string)
	userPerusahaanID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)

	err := h.service.RequestEdit(r.Context(), id, req.Reason, userRole, userPerusahaanID)
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}

	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Permintaan edit berhasil diajukan ke admin",
	})
}

// ApproveEditIkas godoc
//
//	@Summary		Setujui permintaan edit IKAS
//	@Description	Menyetujui permintaan edit dan membuka kembali kunci data IKAS. Hanya bisa dilakukan oleh Admin atau Staff.
//	@Tags			Ikas
//	@Produce		json
//	@Param			id	path		string	true	"Ikas ID"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		403	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/api/maturity/ikas/{id}/approve-edit [put]
func (h *IkasHandler) handleApproveEdit(w http.ResponseWriter, r *http.Request, id string) {
	if !utils.IsValidUUID(id) {
		utils.RespondError(w, 400, "ID tidak valid")
		return
	}

	userRole, _ := r.Context().Value(middleware.Role).(string)
	if userRole != "admin" && userRole != "staff" {
		utils.RespondError(w, 403, "Hanya admin yang dapat menyetujui permintaan edit")
		return
	}

	err := h.service.ApproveEdit(r.Context(), id)
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}

	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Permintaan edit disetujui, data IKAS telah dibuka",
	})
}

// RejectEditIkas godoc
//
//	@Summary		Tolak permintaan edit IKAS
//	@Description	Menolak permintaan edit data IKAS. Hanya bisa dilakukan oleh Admin atau Staff.
//	@Tags			Ikas
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Ikas ID"
//	@Param			req		body		dto.RejectEditRequest	true	"Alasan penolakan"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		403		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/api/maturity/ikas/{id}/reject-edit [put]
func (h *IkasHandler) handleRejectEdit(w http.ResponseWriter, r *http.Request, id string) {
	if !utils.IsValidUUID(id) {
		utils.RespondError(w, 400, "ID tidak valid")
		return
	}

	var req dto.RejectEditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	userRole, _ := r.Context().Value(middleware.Role).(string)
	if userRole != "admin" && userRole != "staff" {
		utils.RespondError(w, 403, "Hanya admin yang dapat menolak permintaan edit")
		return
	}

	err := h.service.RejectEdit(r.Context(), id, req.Reason)
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}

	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Permintaan edit ditolak",
	})
}
