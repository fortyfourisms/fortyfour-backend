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
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"

	"github.com/nfnt/resize"
)

type PerusahaanHandler struct {
	service    *services.PerusahaanService
	uploadPath string
}

func NewPerusahaanHandler(service *services.PerusahaanService, uploadPath string) *PerusahaanHandler {
	return &PerusahaanHandler{service: service, uploadPath: uploadPath}
}

func (h *PerusahaanHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/perusahaan")
	id = strings.TrimPrefix(id, "/")

	switch r.Method {
	case http.MethodGet:
		if id == "" {
			data, err := h.service.GetAll()
			if err != nil {
				utils.RespondError(w, 500, err.Error())
				return
			}
			utils.RespondJSON(w, 200, data)
		} else {
			p, err := h.service.GetByID(id)
			if err != nil {
				utils.RespondError(w, 404, "Data tidak ditemukan")
				return
			}
			utils.RespondJSON(w, 200, p)
		}

	case http.MethodPost:
		var req dto.CreatePerusahaanRequest
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			utils.RespondError(w, 400, "Gagal membaca form data")
			return
		}
		req = h.parseCreateForm(r.MultipartForm)

		file, header, err := r.FormFile("photo")
		if err == nil {
			defer file.Close()
			filename, err := h.handleFileUpload(file, header)
			if err != nil {
				utils.RespondError(w, 400, err.Error())
				return
			}
			req.Photo = &filename
		}

		resp, err := h.service.Create(req)
		if err != nil {
			utils.RespondError(w, 400, err.Error())
			return
		}
		utils.RespondJSON(w, 201, resp)

	case http.MethodPut:
		if id == "" {
			utils.RespondError(w, 400, "ID wajib")
			return
		}

		var req dto.UpdatePerusahaanRequest
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			utils.RespondError(w, 400, "Gagal membaca form data")
			return
		}
		req = h.parseUpdateForm(r.MultipartForm)

		file, header, err := r.FormFile("photo")
		if err == nil {
			defer file.Close()

			// ambil data lama dari DB
			perusahaan, errGet := h.service.GetByID(id)
			if errGet == nil && perusahaan.Photo != "" {
				oldPath := filepath.Join(h.uploadPath, perusahaan.Photo)
				os.Remove(oldPath)
			}

			filename, err := h.handleFileUpload(file, header)
			if err != nil {
				utils.RespondError(w, 400, err.Error())
				return
			}
			req.Photo = &filename
		}

		resp, err := h.service.Update(id, req)
		if err != nil {
			utils.RespondError(w, 400, err.Error())
			return
		}
		utils.RespondJSON(w, 200, resp)

	case http.MethodDelete:
		if id == "" {
			utils.RespondError(w, 400, "ID wajib")
			return
		}
		if err := h.service.Delete(id); err != nil {
			utils.RespondError(w, 400, err.Error())
			return
		}
		utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *PerusahaanHandler) parseCreateForm(form *multipart.Form) dto.CreatePerusahaanRequest {
	return dto.CreatePerusahaanRequest{
		NamaPerusahaan: getFormValue(form, "nama_perusahaan"),
		JenisUsaha:     getFormValue(form, "jenis_usaha"),
		Alamat:         getFormValue(form, "alamat"),
		Telepon:        getFormValue(form, "telepon"),
		Email:          getFormValue(form, "email"),
		Website:        getFormValue(form, "website"),
	}
}

func (h *PerusahaanHandler) parseUpdateForm(form *multipart.Form) dto.UpdatePerusahaanRequest {
	return dto.UpdatePerusahaanRequest{
		NamaPerusahaan: getFormValue(form, "nama_perusahaan"),
		JenisUsaha:     getFormValue(form, "jenis_usaha"),
		Alamat:         getFormValue(form, "alamat"),
		Telepon:        getFormValue(form, "telepon"),
		Email:          getFormValue(form, "email"),
		Website:        getFormValue(form, "website"),
	}
}

func getFormValue(form *multipart.Form, key string) *string {
	values := form.Value[key]
	if len(values) > 0 {
		return &values[0]
	}
	return nil
}

func (h *PerusahaanHandler) handleFileUpload(file multipart.File, header *multipart.FileHeader) (string, error) {
	buff := &bytes.Buffer{}
	size, err := io.Copy(buff, file)
	if err != nil {
		return "", fmt.Errorf("gagal membaca file")
	}

	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("hanya file gambar yang diizinkan")
	}

	img, _, err := image.Decode(bytes.NewReader(buff.Bytes()))
	if err != nil {
		return "", fmt.Errorf("gagal decode image")
	}

	if size > 1<<20 { // >1MB compress
		img = resize.Resize(0, 1024, img, resize.Lanczos3)
	}

	ext := filepath.Ext(header.Filename)
	hash := sha256.Sum256([]byte(fmt.Sprintf("%d%s", time.Now().UnixNano(), header.Filename)))
	filename := hex.EncodeToString(hash[:]) + ext

	outPath := filepath.Join(h.uploadPath, filename)
	out, err := os.Create(outPath)
	if err != nil {
		return "", fmt.Errorf("gagal menyimpan file")
	}
	defer out.Close()
	jpeg.Encode(out, img, &jpeg.Options{Quality: 80})

	return filename, nil
}
