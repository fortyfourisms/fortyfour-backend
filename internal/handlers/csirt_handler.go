package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
)

type CsirtHandler struct {
	service *services.CsirtService
}

func NewCsirtHandler(service *services.CsirtService) *CsirtHandler {
	return &CsirtHandler{service: service}
}

func (h *CsirtHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/api/csirt"), "/")

	switch r.Method {
	case http.MethodGet:
		if id == "" {
			h.handleGetAll(w)
		} else {
			h.handleGetByID(w, id)
		}
	case http.MethodPost:
		h.handleCreate(w, r)
	case http.MethodPut:
		h.handleUpdate(w, r, id)
	case http.MethodDelete:
		h.handleDelete(w, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// GetAllCsirt godoc
// @Summary      List semua csirt
// @Description  Ambil seluruh data csirt
// @Tags         Csirt
// @Produce      json
// @Success      200  {array}  dto.CsirtResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /api/csirt [get]
func (h *CsirtHandler) handleGetAll(w http.ResponseWriter) {
	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetCsirtByID godoc
// @Summary      Ambil csirt berdasarkan ID
// @Description  Mengambil satu data csirt
// @Tags         Csirt
// @Produce      json
// @Param        id   path      string  true  "Csirt ID"
// @Success      200  {object} dto.CsirtResponse
// @Failure      404  {object} dto.ErrorResponse
// @Router       /api/csirt/{id} [get]
func (h *CsirtHandler) handleGetByID(w http.ResponseWriter, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}
	utils.RespondJSON(w, 200, data)
}

// CreateCsirt godoc
// @Summary      Tambah csirt baru
// @Description  Membuat record csirt
// @Tags         Csirt
// @Accept       json
// @Produce      json
// @Param        csirt body dto.CreateCsirtRequest true "Data csirt"
// @Success      201  {object} dto.CsirtResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/csirt [post]
func (h *CsirtHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.RespondError(w, 400, "Gagal membaca form-data")
		return
	}

	req := dto.CreateCsirtRequest{
		IdPerusahaan: r.FormValue("id_perusahaan"),
		NamaCsirt:    r.FormValue("nama_csirt"),
		WebCsirt:     r.FormValue("web_csirt"),
		TeleponCsirt: r.FormValue("telepon_csirt"),
	}

	photoPath, err := saveUploadedFile(r, "photo_csirt", "uploads/csirt_photo")
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}
	req.PhotoCsirt = photoPath

	rfcPath, err := saveUploadedFile(r, "file_rfc2350", "uploads/rfc2350")
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}
	req.FileRFC2350 = rfcPath

	pgpPath, err := saveUploadedFile(r, "file_public_key_pgp", "uploads/pgp")
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}
	req.FilePublicKeyPGP = pgpPath

	resp, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 201, resp)
}

// UpdateCsirt godoc
// @Summary      Update csirt
// @Description  Mengubah data csirt berdasarkan ID
// @Tags         Csirt
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Csirt ID"
// @Param        csirt body      dto.UpdateCsirtRequest true "Data update"
// @Success      200  {object} dto.CsirtResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/csirt/{id} [put]
func (h *CsirtHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.RespondError(w, 400, "Gagal membaca form-data")
		return
	}

	req := dto.UpdateCsirtRequest{}

	if v := r.FormValue("nama_csirt"); v != "" {
		req.NamaCsirt = &v
	}
	if v := r.FormValue("web_csirt"); v != "" {
		req.WebCsirt = &v
	}
	if v := r.FormValue("telepon_csirt"); v != "" {
		req.TeleponCsirt = &v
	}

	if path, err := saveUploadedFile(r, "photo_csirt", "uploads/csirt_photo"); err == nil && path != "" {
		req.PhotoCsirt = &path
	}

	if path, err := saveUploadedFile(r, "file_rfc2350", "uploads/rfc2350"); err == nil && path != "" {
		req.FileRFC2350 = &path
	}

	if path, err := saveUploadedFile(r, "file_public_key_pgp", "uploads/pgp"); err == nil && path != "" {
		req.FilePublicKeyPGP = &path
	}

	resp, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 200, resp)
}

// DeleteCsirt godoc
// @Summary      Hapus csirt
// @Description  Menghapus data csirt berdasarkan ID
// @Tags         Csirt
// @Produce      json
// @Param        id  path  string  true  "Csirt ID"
// @Success      200  {object} dto.MessageResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/csirt/{id} [delete]
func (h *CsirtHandler) handleDelete(w http.ResponseWriter, id string) {
	if err := h.service.Delete(id); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}

func saveUploadedFile(r *http.Request, fieldName, uploadDir string) (string, error) {
	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return "", nil
	}
	defer file.Close()

	os.MkdirAll(uploadDir, os.ModePerm)

	ext := filepath.Ext(header.Filename)
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
