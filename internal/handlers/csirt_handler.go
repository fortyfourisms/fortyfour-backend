package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"fortyfour-backend/pkg/logger"
)

type CsirtHandler struct {
	service       services.CsirtServiceInterface
	sseService    services.SSEServiceInterface
	exportHandler *CsirtExportHandler
}

func NewCsirtHandler(service services.CsirtServiceInterface, sseService services.SSEServiceInterface) *CsirtHandler {
	return &CsirtHandler{service: service, sseService: sseService}
}

// SetExportHandler injects the export handler so CsirtHandler can delegate
// requests matching /{id}/export-pdf.
func (h *CsirtHandler) SetExportHandler(exportH *CsirtExportHandler) {
	h.exportHandler = exportH
}

func (h *CsirtHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/api/csirt"), "/")

	// Delegate /{id}/export-pdf ke CsirtExportHandler
	if strings.HasSuffix(id, "/export-pdf") || id == "export-pdf" {
		if h.exportHandler != nil {
			h.exportHandler.ServeHTTP(w, r)
		} else {
			utils.RespondError(w, 500, "Export handler tidak tersedia")
		}
		return
	}

	// Handle /{id}/pgp-download → download file PGP sebagai .asc
	if strings.HasSuffix(id, "/pgp-download") {
		csirtID := strings.TrimSuffix(id, "/pgp-download")
		h.handlePGPDownload(w, r, csirtID)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if id == "" {
			h.handleGetAll(w, r)
		} else {
			h.handleGetByID(w, r, id)
		}
	case http.MethodPost:
		h.handleCreate(w, r)
	case http.MethodPut:
		h.handleUpdate(w, r, id)
	case http.MethodDelete:
		h.handleDelete(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *CsirtHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())

	if role == "admin" {
		data, err := h.service.GetAll()
		if err != nil {
			logger.Error(err, "failed to get all CSIRT data")
			utils.RespondError(w, 500, err.Error())
			return
		}
		utils.RespondJSON(w, 200, data)
		return
	}

	idPerusahaan := middleware.GetIDPerusahaan(r.Context())
	if idPerusahaan == "" {
		utils.RespondError(w, 403, "Akun Anda belum terhubung ke perusahaan")
		return
	}

	data, err := h.service.GetByPerusahaan(idPerusahaan)
	if err != nil {
		logger.Error(err, "failed to get CSIRT data by perusahaan")
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// @Summary Ambil CSIRT berdasarkan ID
// @Description Mengambil satu data CSIRT
// @Tags CSIRT
// @Produce json
// @Security BearerAuth
// @Param id path string true "CSIRT ID"
// @Success 200 {object} dto.CsirtResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/csirt/{id} [get]
func (h *CsirtHandler) handleGetByID(w http.ResponseWriter, r *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		logger.Error(err, "failed to get CSIRT by ID")
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}

	role := middleware.GetRole(r.Context())
	if role == "user" {
		idPerusahaan := middleware.GetIDPerusahaan(r.Context())
		if data.Perusahaan.ID != idPerusahaan {
			utils.RespondError(w, 403, "Anda tidak memiliki akses ke data ini")
			return
		}
	}

	utils.RespondJSON(w, 200, data)
}

// @Summary Download file Public Key PGP CSIRT
// @Description Mendownload file PGP (.asc) milik CSIRT berdasarkan ID
// @Tags CSIRT
// @Produce application/octet-stream
// @Security BearerAuth
// @Param id path string true "CSIRT ID"
// @Success 200 {file} binary
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/csirt/{id}/pgp-download [get]
func (h *CsirtHandler) handlePGPDownload(w http.ResponseWriter, r *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		logger.Error(err, "failed to get CSIRT for PGP download")
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}

	// Ownership check untuk user
	role := middleware.GetRole(r.Context())
	if role == "user" {
		idPerusahaan := middleware.GetIDPerusahaan(r.Context())
		if data.Perusahaan.ID != idPerusahaan {
			utils.RespondError(w, 403, "Anda tidak memiliki akses ke data ini")
			return
		}
	}

	if data.FilePublicKeyPGP == "" {
		utils.RespondError(w, 404, "File PGP tidak tersedia")
		return
	}

	filePath := data.FilePublicKeyPGP
	f, err := os.Open(filePath)
	if err != nil {
		logger.Error(err, "failed to open PGP file")
		utils.RespondError(w, 404, "File PGP tidak ditemukan di server")
		return
	}
	defer f.Close()

	// Paksa nama file sebagai .asc saat didownload
	downloadName := fmt.Sprintf("%s_pgp_key.asc", id)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", downloadName))
	w.Header().Set("Content-Type", "application/pgp-keys")
	w.WriteHeader(http.StatusOK)
	io.Copy(w, f)
}

// @Summary Tambah CSIRT baru
// @Description Membuat record CSIRT baru dengan file upload
// @Tags CSIRT
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id_perusahaan formData string true "ID Perusahaan"
// @Param nama_csirt formData string true "Nama CSIRT"
// @Param web_csirt formData string false "Website CSIRT"
// @Param email_csirt formData string false "Email CSIRT"
// @Param telepon_csirt formData string false "Telepon CSIRT"
// @Param photo_csirt formData file false "Photo CSIRT"
// @Param file_rfc2350 formData file false "File RFC2350"
// @Param file_public_key_pgp formData file false "File Public Key PGP (.asc)"
// @Param file_str formData file false "File STR CSIRT (nullable)"
// @Param tanggal_registrasi formData string false "Tanggal Registrasi (YYYY-MM-DD)"
// @Param tanggal_kadaluarsa formData string false "Tanggal Kadaluarsa (YYYY-MM-DD)"
// @Success 201 {object} dto.CsirtResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/csirt [post]
func (h *CsirtHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		logger.Error(err, "failed to parse multipart form for CSIRT create")
		utils.RespondError(w, 400, "Gagal membaca form-data")
		return
	}

	req := dto.CreateCsirtRequest{
		IdPerusahaan:      r.FormValue("id_perusahaan"),
		NamaCsirt:         r.FormValue("nama_csirt"),
		WebCsirt:          r.FormValue("web_csirt"),
		EmailCsirt:        r.FormValue("email_csirt"),
		TeleponCsirt:      r.FormValue("telepon_csirt"),
		TanggalRegistrasi: r.FormValue("tanggal_registrasi"),
		TanggalKadaluarsa: r.FormValue("tanggal_kadaluarsa"),
	}

	// User: paksa id_perusahaan dari JWT, tidak bisa diisi sembarangan
	role := middleware.GetRole(r.Context())
	if role == "user" {
		idPerusahaan := middleware.GetIDPerusahaan(r.Context())
		if idPerusahaan == "" {
			utils.RespondError(w, 403, "Akun Anda belum terhubung ke perusahaan")
			return
		}
		req.IdPerusahaan = idPerusahaan
	}

	photoPath, err := saveUploadedFile(r, "photo_csirt", "uploads/csirt_photo", "")
	if err != nil {
		logger.Error(err, "failed to upload CSIRT photo")
		utils.RespondError(w, 400, err.Error())
		return
	}
	req.PhotoCsirt = photoPath

	rfcPath, err := saveUploadedFile(r, "file_rfc2350", "uploads/rfc2350", "")
	if err != nil {
		logger.Error(err, "failed to upload RFC2350 file")
		utils.RespondError(w, 400, err.Error())
		return
	}
	req.FileRFC2350 = rfcPath

	// PGP: hanya terima file .asc
	pgpPath, err := saveUploadedFile(r, "file_public_key_pgp", "uploads/pgp", ".asc")
	if err != nil {
		logger.Error(err, "failed to upload PGP key file")
		utils.RespondError(w, 400, err.Error())
		return
	}
	req.FilePublicKeyPGP = pgpPath

	// Upload file STR (nullable — tidak wajib)
	strPath, err := saveUploadedFile(r, "file_str", "uploads/str_csirt", "")
	if err != nil {
		logger.Error(err, "failed to upload STR file")
		utils.RespondError(w, 400, err.Error())
		return
	}
	req.FileStr = strPath

	resp, err := h.service.Create(req)
	if err != nil {
		logger.Error(err, "failed to create CSIRT")
		utils.RespondError(w, 400, err.Error())
		return
	}

	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyCreate("csirt", resp, userID)

	utils.RespondJSON(w, 201, resp)
}

// @Summary Update CSIRT
// @Description Mengubah data CSIRT berdasarkan ID
// @Tags CSIRT
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "CSIRT ID"
// @Param nama_csirt formData string false "Nama CSIRT"
// @Param web_csirt formData string false "Website CSIRT"
// @Param email_csirt formData string false "Email CSIRT"
// @Param telepon_csirt formData string false "Telepon CSIRT"
// @Param photo_csirt formData file false "Photo CSIRT"
// @Param file_rfc2350 formData file false "File RFC2350"
// @Param file_public_key_pgp formData file false "File Public Key PGP (.asc)"
// @Param file_str formData file false "File STR CSIRT (nullable)"
// @Param tanggal_registrasi formData string false "Tanggal Registrasi (YYYY-MM-DD)"
// @Param tanggal_kadaluarsa formData string false "Tanggal Kadaluarsa (YYYY-MM-DD)"
// @Success 200 {object} dto.CsirtResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/csirt/{id} [put]
func (h *CsirtHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		logger.Error(err, "failed to parse multipart form for CSIRT update")
		utils.RespondError(w, 400, "Gagal membaca form-data")
		return
	}

	// Ownership check untuk user
	role := middleware.GetRole(r.Context())
	if role == "user" {
		existing, err := h.service.GetByID(id)
		if err != nil {
			utils.RespondError(w, 404, "Data tidak ditemukan")
			return
		}
		idPerusahaan := middleware.GetIDPerusahaan(r.Context())
		if existing.Perusahaan.ID != idPerusahaan {
			utils.RespondError(w, 403, "Anda tidak memiliki akses ke data ini")
			return
		}
	}

	req := dto.UpdateCsirtRequest{}

	if v := r.FormValue("nama_csirt"); v != "" {
		req.NamaCsirt = &v
	}
	if v := r.FormValue("web_csirt"); v != "" {
		req.WebCsirt = &v
	}
	if v := r.FormValue("email_csirt"); v != "" {
		req.EmailCsirt = &v
	}
	if v := r.FormValue("telepon_csirt"); v != "" {
		req.TeleponCsirt = &v
	}
	if v := r.FormValue("tanggal_registrasi"); v != "" {
		req.TanggalRegistrasi = &v
	}
	if v := r.FormValue("tanggal_kadaluarsa"); v != "" {
		req.TanggalKadaluarsa = &v
	}

	if path, err := saveUploadedFile(r, "photo_csirt", "uploads/csirt_photo", ""); err == nil && path != "" {
		req.PhotoCsirt = &path
	}

	if path, err := saveUploadedFile(r, "file_rfc2350", "uploads/rfc2350", ""); err == nil && path != "" {
		req.FileRFC2350 = &path
	}

	// PGP: hanya terima file .asc
	if path, err := saveUploadedFile(r, "file_public_key_pgp", "uploads/pgp", ".asc"); err != nil {
		logger.Error(err, "failed to upload PGP key file")
		utils.RespondError(w, 400, err.Error())
		return
	} else if path != "" {
		req.FilePublicKeyPGP = &path
	}

	if path, err := saveUploadedFile(r, "file_str", "uploads/str_csirt", ""); err == nil && path != "" {
		req.FileStr = &path
	}

	resp, err := h.service.Update(id, req)
	if err != nil {
		logger.Error(err, "failed to update CSIRT")
		utils.RespondError(w, 400, err.Error())
		return
	}

	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyUpdate("csirt", resp, userID)

	utils.RespondJSON(w, 200, resp)
}

// @Summary Hapus CSIRT
// @Description Menghapus data CSIRT berdasarkan ID
// @Tags CSIRT
// @Produce json
// @Security BearerAuth
// @Param id path string true "CSIRT ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/csirt/{id} [delete]
func (h *CsirtHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	// Ownership check untuk user
	role := middleware.GetRole(r.Context())
	if role == "user" {
		existing, err := h.service.GetByID(id)
		if err != nil {
			utils.RespondError(w, 404, "Data tidak ditemukan")
			return
		}
		idPerusahaan := middleware.GetIDPerusahaan(r.Context())
		if existing.Perusahaan.ID != idPerusahaan {
			utils.RespondError(w, 403, "Anda tidak memiliki akses ke data ini")
			return
		}
	}

	if err := h.service.Delete(id); err != nil {
		logger.Error(err, "failed to delete CSIRT")
		utils.RespondError(w, 400, err.Error())
		return
	}

	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyDelete("csirt", id, userID)

	utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}

// saveUploadedFile menyimpan file yang diupload ke direktori tujuan.
// Jika requiredExt tidak kosong (misal ".asc"), hanya file dengan ekstensi tersebut yang diterima.
// Mengembalikan string kosong (tanpa error) jika field tidak ada di form.
func saveUploadedFile(r *http.Request, fieldName, uploadDir, requiredExt string) (string, error) {
	file, header, err := r.FormFile(fieldName)
	if err != nil {
		// Field tidak ada → bukan error, kembalikan string kosong
		return "", nil
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))

	// Validasi ekstensi jika requiredExt diberikan
	if requiredExt != "" && ext != strings.ToLower(requiredExt) {
		return "", fmt.Errorf("file %s harus berformat %s, diterima: %s", fieldName, requiredExt, ext)
	}

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}

	filename := uuid.New().String() + ext
	fullPath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return fullPath, nil
}
