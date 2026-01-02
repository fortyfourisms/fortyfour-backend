package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"

	"github.com/nfnt/resize"
)

const (
	maxUploadSize = 10 << 20
	maxImageSize  = 1 << 20
	imageQuality  = 80
	resizeHeight  = 1024
)

type PerusahaanHandler struct {
	service    *services.PerusahaanService
	uploadPath string
	sseService *services.SSEService
}

func NewPerusahaanHandler(service *services.PerusahaanService, uploadPath string, sseService *services.SSEService) *PerusahaanHandler {
	return &PerusahaanHandler{
		service:    service,
		uploadPath: uploadPath,
		sseService: sseService,
	}
}

func (h *PerusahaanHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/api/perusahaan"), "/")

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
		h.handleUpdate(w, r, id)
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

// GetAllPerusahaan godoc
// @Summary      List semua perusahaan
// @Description  Mengambil seluruh data perusahaan
// @Tags         Perusahaan
// @Produce      json
// @Success      200  {array}  dto.PerusahaanResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /api/perusahaan [get]
func (h *PerusahaanHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetPerusahaanByID godoc
// @Summary      Ambil perusahaan berdasarkan ID
// @Description  Mengambil satu data perusahaan
// @Tags         Perusahaan
// @Produce      json
// @Param        id   path      string  true  "Perusahaan ID"
// @Success      200  {object} dto.PerusahaanResponse
// @Failure      404  {object} dto.ErrorResponse
// @Router       /api/perusahaan/{id} [get]
func (h *PerusahaanHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}
	utils.RespondJSON(w, 200, data)
}

// CreatePerusahaan godoc
// @Summary      Tambah perusahaan baru
// @Description  Membuat record perusahaan
// @Tags         Perusahaan
// @Accept       json
// @Produce      json
// @Param        perusahaan body dto.CreatePerusahaanRequest true "Data perusahaan"
// @Success      201  {object} dto.PerusahaanResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/perusahaan [post]
func (h *PerusahaanHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		utils.RespondError(w, 400, "Gagal membaca form data")
		return
	}

	req := h.parseCreateForm(r.MultipartForm)

	// Handle file upload jika ada
	if filename, err := h.processFileUpload(r); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	} else if filename != "" {
		req.Photo = &filename
	}

	resp, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	// SSE Notif Create
	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyCreate("perusahaan", resp, userID)

	utils.RespondJSON(w, 201, resp)
}

// UpdatePerusahaan godoc
// @Summary      Update perusahaan
// @Description  Mengubah data perusahaan berdasarkan ID
// @Tags         Perusahaan
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Perusahaan ID"
// @Param        perusahaan body      dto.UpdatePerusahaanRequest true "Data update"
// @Success      200  {object} dto.PerusahaanResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/perusahaan/{id} [put]
func (h *PerusahaanHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		utils.RespondError(w, 400, "Gagal membaca form data")
		return
	}

	req := h.parseUpdateForm(r.MultipartForm)

	// Handle file upload jika ada
	if filename, err := h.processFileUpload(r); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	} else if filename != "" {
		// Hapus file lama
		h.deleteOldPhoto(id)
		req.Photo = &filename
	}

	resp, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	// SSE Notif Update
	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyUpdate("perusahaan", resp, userID)

	utils.RespondJSON(w, 200, resp)
}

// DeletePerusahaan godoc
// @Summary      Hapus perusahaan
// @Description  Menghapus data perusahaan berdasarkan ID
// @Tags         Perusahaan
// @Produce      json
// @Param        id  path  string  true  "Perusahaan ID"
// @Success      200  {object} dto.MessageResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/perusahaan/{id} [delete]
func (h *PerusahaanHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	// Hapus file photo sebelum delete record
	h.deleteOldPhoto(id)

	if err := h.service.Delete(id); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	// SSE Notif Delete
	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyDelete("perusahaan", id, userID)

	utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}

// processFileUpload menangani upload file dan return filename (empty string jika tidak ada file)
func (h *PerusahaanHandler) processFileUpload(r *http.Request) (string, error) {
	file, header, err := r.FormFile("photo")
	if err != nil {
		// Tidak ada file di-upload, bukan error
		if err == http.ErrMissingFile {
			return "", nil
		}
		return "", fmt.Errorf("gagal membaca file: %v", err)
	}
	defer file.Close()

	return h.saveUploadedFile(file, header)
}

// deleteOldPhoto menghapus file photo lama dari disk
func (h *PerusahaanHandler) deleteOldPhoto(id string) {
	perusahaan, err := h.service.GetByID(id)
	if err == nil && perusahaan.Photo != "" {
		oldPath := filepath.Join(h.uploadPath, perusahaan.Photo)
		os.Remove(oldPath) // ignore error
	}
}

// UBAH INI: Ganti "sektor" menjadi "id_sub_sektor"
func (h *PerusahaanHandler) parseCreateForm(form *multipart.Form) dto.CreatePerusahaanRequest {
	return dto.CreatePerusahaanRequest{
		NamaPerusahaan: getFormValue(form, "nama_perusahaan"),
		IDSubSektor:    getFormValue(form, "id_sub_sektor"), // Changed from Sektor
		Alamat:         getFormValue(form, "alamat"),
		Telepon:        getFormValue(form, "telepon"),
		Email:          getFormValue(form, "email"),
		Website:        getFormValue(form, "website"),
	}
}

// UBAH INI: Ganti "sektor" menjadi "id_sub_sektor"
func (h *PerusahaanHandler) parseUpdateForm(form *multipart.Form) dto.UpdatePerusahaanRequest {
	return dto.UpdatePerusahaanRequest{
		NamaPerusahaan: getFormValue(form, "nama_perusahaan"),
		IDSubSektor:    getFormValue(form, "id_sub_sektor"), // Changed from Sektor
		Alamat:         getFormValue(form, "alamat"),
		Telepon:        getFormValue(form, "telepon"),
		Email:          getFormValue(form, "email"),
		Website:        getFormValue(form, "website"),
	}
}

func (h *PerusahaanHandler) saveUploadedFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Validasi content type
	if err := h.validateImageFile(header); err != nil {
		return "", err
	}

	// Baca file ke buffer
	buff := &bytes.Buffer{}
	size, err := io.Copy(buff, file)
	if err != nil {
		return "", fmt.Errorf("gagal membaca file")
	}

	// Decode dan resize jika perlu
	img, err := h.processImage(buff.Bytes(), size)
	if err != nil {
		return "", err
	}

	// Generate filename unik
	filename := h.generateFilename(header.Filename)

	// Simpan file
	if err := h.saveImageToFile(img, filename); err != nil {
		return "", err
	}

	return filename, nil
}

func (h *PerusahaanHandler) validateImageFile(header *multipart.FileHeader) error {
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return fmt.Errorf("hanya file gambar yang diizinkan")
	}
	return nil
}

func (h *PerusahaanHandler) processImage(data []byte, size int64) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("gagal decode image")
	}

	// Resize jika ukuran lebih dari 1MB
	if size > maxImageSize {
		img = resize.Resize(0, resizeHeight, img, resize.Lanczos3)
	}

	return img, nil
}

func (h *PerusahaanHandler) generateFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	hash := sha256.Sum256([]byte(fmt.Sprintf("%d%s", time.Now().UnixNano(), originalFilename)))
	return hex.EncodeToString(hash[:]) + ext
}

func (h *PerusahaanHandler) saveImageToFile(img image.Image, filename string) error {
	outPath := filepath.Join(h.uploadPath, filename)
	out, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("gagal menyimpan file")
	}
	defer out.Close()

	return jpeg.Encode(out, img, &jpeg.Options{Quality: imageQuality})
}

func getFormValue(form *multipart.Form, key string) *string {
	if values := form.Value[key]; len(values) > 0 {
		return &values[0]
	}
	return nil
}