package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"fortyfour-backend/pkg/logger"
)

// LMSHandler menangani semua endpoint LMS:
//   /api/kelas             → kelas (admin CRUD, user read)
//   /api/materi            → materi (admin CRUD, user progress)
//   /api/soal              → soal kuis (admin CRUD)
//   /api/kuis              → kuis (admin CRUD + user: start & submit)
//   /api/file-pendukung    → file pendukung (admin upload, user download)
//   /api/sertifikat        → sertifikat user

type LMSHandler struct {
	kelasSvc      *services.KelasService
	materiSvc     *services.MateriService
	soalSvc       *services.SoalService
	kuisSvc       *services.KuisService
	fpSvc         *services.FilePendukungService
	feedbackSvc   *services.FeedbackService
	sertifikatSvc *services.SertifikatService
	sseSvc        *services.SSEService
}

func NewLMSHandler(
	kelasSvc *services.KelasService,
	materiSvc *services.MateriService,
	soalSvc *services.SoalService,
	kuisSvc *services.KuisService,
	fpSvc *services.FilePendukungService,
	feedbackSvc *services.FeedbackService,
	sertifikatSvc *services.SertifikatService,
	sseSvc *services.SSEService,
) *LMSHandler {
	return &LMSHandler{
		kelasSvc:      kelasSvc,
		materiSvc:     materiSvc,
		soalSvc:       soalSvc,
		kuisSvc:       kuisSvc,
		fpSvc:         fpSvc,
		feedbackSvc:   feedbackSvc,
		sertifikatSvc: sertifikatSvc,
		sseSvc:        sseSvc,
	}
}

// ── helper ────────────────────────────────────────────────────────────────────

func getUserID(r *http.Request) string {
	if uid, ok := r.Context().Value(middleware.UserIDKey).(string); ok {
		return uid
	}
	return ""
}

func trimID(path, prefix string) string {
	return strings.TrimPrefix(strings.TrimPrefix(path, prefix), "/")
}

// ════════════════════════════════════════════════════════════════════════════
// PUBLIC KELAS  —  /api/public/kelas (tanpa auth, landing page)
// ════════════════════════════════════════════════════════════════════════════

// ServePublicKelas menangani endpoint publik untuk kelas (tanpa auth).
// GET /api/public/kelas       → list kelas published
// GET /api/public/kelas/{id}  → detail kelas (tanpa progress user)
func (h *LMSHandler) ServePublicKelas(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := trimID(r.URL.Path, "/api/public/kelas")
	if id == "" {
		// List semua kelas yang statusnya 'published'
		data, err := h.kelasSvc.GetAll(true)
		if err != nil {
			logger.Error(err, "publicKelasGetAll failed")
			utils.RespondError(w, 500, err.Error())
			return
		}
		utils.RespondJSON(w, 200, data)
	} else {
		// Detail kelas tanpa userID → progress di-skip oleh service
		data, err := h.kelasSvc.GetDetail(id, "")
		if err != nil {
			logger.Error(err, "publicKelasGetDetail failed")
			utils.RespondError(w, 404, err.Error())
			return
		}
		utils.RespondJSON(w, 200, data)
	}
}

// ════════════════════════════════════════════════════════════════════════════
// KELAS  —  /api/kelas  dan  /api/kelas/{id}
// ════════════════════════════════════════════════════════════════════════════

func (h *LMSHandler) ServeKelas(w http.ResponseWriter, r *http.Request) {
	id := trimID(r.URL.Path, "/api/kelas")

	// Handle nested routes: /api/kelas/{id}/materi or /api/kelas/{id}/materi/{id_materi}
	if strings.Contains(id, "/materi") {
		parts := strings.SplitN(id, "/materi", 2)
		idKelas := parts[0]
		materiSuffix := ""
		if len(parts) > 1 {
			materiSuffix = strings.TrimPrefix(parts[1], "/")
		}

		switch r.Method {
		case http.MethodPost:
			if materiSuffix != "" {
				utils.RespondError(w, 400, "ID materi tidak diperlukan untuk create")
				return
			}
			h.materiCreate(w, r, idKelas)
		case http.MethodPut:
			if materiSuffix == "" {
				utils.RespondError(w, 400, "ID materi wajib untuk update")
				return
			}
			h.materiUpdate(w, r, materiSuffix)
		case http.MethodDelete:
			if materiSuffix == "" {
				utils.RespondError(w, 400, "ID materi wajib untuk delete")
				return
			}
			h.materiDelete(w, r, materiSuffix)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}

	// Handle /api/kelas/{id}/kuis
	if strings.Contains(id, "/kuis") {
		parts := strings.SplitN(id, "/kuis", 2)
		idKelas := parts[0]
		kuisSuffix := ""
		if len(parts) > 1 {
			kuisSuffix = strings.TrimPrefix(parts[1], "/")
		}
		h.serveKuisByKelas(w, r, idKelas, kuisSuffix)
		return
	}

	// Handle /api/kelas/{id}/soal
	if strings.Contains(id, "/soal") {
		parts := strings.SplitN(id, "/soal", 2)
		idKuis := parts[0]
		h.serveSoalByKuis(w, r, idKuis)
		return
	}

	// Handle /api/kelas/{id}/sertifikat
	if strings.Contains(id, "/sertifikat") {
		parts := strings.SplitN(id, "/sertifikat", 2)
		idKelas := parts[0]
		sertSuffix := ""
		if len(parts) > 1 {
			sertSuffix = strings.TrimPrefix(parts[1], "/")
		}
		h.serveSertifikatByKelas(w, r, idKelas, sertSuffix)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if id == "" {
			h.kelasGetAll(w, r)
		} else {
			h.kelasGetDetail(w, r, id)
		}
	case http.MethodPost:
		if id != "" {
			utils.RespondError(w, 400, "ID tidak diperlukan untuk create")
			return
		}
		h.kelasCreate(w, r)
	case http.MethodPut:
		if id == "" {
			utils.RespondError(w, 400, "ID wajib")
			return
		}
		h.kelasUpdate(w, r, id)
	case http.MethodDelete:
		if id == "" {
			utils.RespondError(w, 400, "ID wajib")
			return
		}
		h.kelasDelete(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// @Summary		List semua kelas
// @Description	Admin melihat semua kelas. User hanya melihat kelas yang published.
// @Tags			LMS - Kelas
// @Produce		json
// @Security		BearerAuth
// @Success		200	{array}		dto.KelasResponse
// @Failure		500	{object}	dto.ErrorResponse
// @Router			/api/kelas [get]
func (h *LMSHandler) kelasGetAll(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	onlyPublished := role != "admin"

	data, err := h.kelasSvc.GetAll(onlyPublished)
	if err != nil {
		logger.Error(err, "kelasGetAll failed")
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// @Summary		Detail kelas
// @Description	Mendapatkan detail kelas beserta materi dan progress user.
// @Tags			LMS - Kelas
// @Produce		json
// @Security		BearerAuth
// @Param			id	path		string	true	"ID Kelas"
// @Success		200	{object}	dto.KelasResponse
// @Failure		404	{object}	dto.ErrorResponse
// @Router			/api/kelas/{id} [get]
func (h *LMSHandler) kelasGetDetail(w http.ResponseWriter, r *http.Request, id string) {
	userID := getUserID(r)
	data, err := h.kelasSvc.GetDetail(id, userID)
	if err != nil {
		logger.Error(err, "kelasGetDetail failed")
		utils.RespondError(w, 404, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// @Summary		Buat kelas baru
// @Description	Admin membuat kelas baru.
// @Tags			LMS - Kelas
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			request	body		dto.CreateKelasRequest	true	"Data kelas"
// @Success		201		{object}	dto.KelasResponse
// @Failure		400		{object}	dto.ErrorResponse
// @Router			/api/kelas [post]
func (h *LMSHandler) kelasCreate(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	var req dto.CreateKelasRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}
	resp, err := h.kelasSvc.Create(req, userID)
	if err != nil {
		logger.Error(err, "kelasCreate failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	h.sseSvc.NotifyCreate("kelas", resp, userID)
	utils.RespondJSON(w, 201, resp)
}

// @Summary		Update kelas
// @Description	Admin mengupdate kelas.
// @Tags			LMS - Kelas
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path		string					true	"ID Kelas"
// @Param			request	body		dto.UpdateKelasRequest	true	"Data update"
// @Success		200		{object}	dto.KelasResponse
// @Failure		400		{object}	dto.ErrorResponse
// @Router			/api/kelas/{id} [put]
func (h *LMSHandler) kelasUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req dto.UpdateKelasRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}
	resp, err := h.kelasSvc.Update(id, req)
	if err != nil {
		logger.Error(err, "kelasUpdate failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	h.sseSvc.NotifyUpdate("kelas", resp, getUserID(r))
	utils.RespondJSON(w, 200, resp)
}

// @Summary		Hapus kelas
// @Description	Admin menghapus kelas.
// @Tags			LMS - Kelas
// @Produce		json
// @Security		BearerAuth
// @Param			id	path		string	true	"ID Kelas"
// @Success		200	{object}	map[string]string
// @Failure		400	{object}	dto.ErrorResponse
// @Router			/api/kelas/{id} [delete]
func (h *LMSHandler) kelasDelete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.kelasSvc.Delete(id); err != nil {
		logger.Error(err, "kelasDelete failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	h.sseSvc.NotifyDelete("kelas", id, getUserID(r))
	utils.RespondJSON(w, 200, map[string]string{"message": "Kelas berhasil dihapus"})
}

// ════════════════════════════════════════════════════════════════════════════
// MATERI  —  /api/materi/{id}
// ════════════════════════════════════════════════════════════════════════════

func (h *LMSHandler) ServeMateri(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// /api/materi/{id}/progress
	if strings.HasSuffix(path, "/progress") {
		id := trimID(strings.TrimSuffix(path, "/progress"), "/api/materi")
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		h.materiUpdateProgress(w, r, id)
		return
	}

	// /api/materi/{id}/file-pendukung
	if strings.Contains(path, "/file-pendukung") {
		parts := strings.SplitN(strings.TrimPrefix(path, "/api/materi/"), "/file-pendukung", 2)
		idMateri := parts[0]
		switch r.Method {
		case http.MethodGet:
			h.filePendukungGetByMateri(w, r, idMateri)
		case http.MethodPost:
			h.filePendukungUpload(w, r, idMateri)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}

	// /api/materi/{id}/feedback
	if strings.Contains(path, "/feedback") {
		parts := strings.SplitN(strings.TrimPrefix(path, "/api/materi/"), "/feedback", 2)
		idMateri := parts[0]
		switch r.Method {
		case http.MethodGet:
			h.feedbackGet(w, r, idMateri)
		case http.MethodPut:
			h.feedbackUpsert(w, r, idMateri)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}

	id := trimID(path, "/api/materi")
	switch r.Method {
	case http.MethodPut:
		if id == "" {
			utils.RespondError(w, 400, "ID wajib")
			return
		}
		h.materiUpdate(w, r, id)
	case http.MethodDelete:
		if id == "" {
			utils.RespondError(w, 400, "ID wajib")
			return
		}
		h.materiDelete(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// ServeMateriByKelas dipakai untuk POST /api/kelas/{id_kelas}/materi
func (h *LMSHandler) ServeMateriByKelas(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/kelas/")
	idKelas := strings.TrimSuffix(path, "/materi")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	h.materiCreate(w, r, idKelas)
}

// @Summary		Tambah materi ke kelas
// @Description	Admin menambahkan materi baru ke dalam kelas.
// @Tags			LMS - Materi
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path		string					true	"ID Kelas"
// @Param			request	body		dto.CreateMateriRequest	true	"Data materi"
// @Success		201		{object}	dto.MateriResponse
// @Failure		400		{object}	dto.ErrorResponse
// @Router			/api/kelas/{id}/materi [post]
func (h *LMSHandler) materiCreate(w http.ResponseWriter, r *http.Request, idKelas string) {
	var req dto.CreateMateriRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}
	resp, err := h.materiSvc.Create(idKelas, req)
	if err != nil {
		logger.Error(err, "materiCreate failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	h.sseSvc.NotifyCreate("materi", resp, getUserID(r))
	utils.RespondJSON(w, 201, resp)
}

// @Summary		Update materi
// @Description	Admin mengupdate materi.
// @Tags			LMS - Materi
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path		string					true	"ID Materi"
// @Param			request	body		dto.UpdateMateriRequest	true	"Data update"
// @Success		200		{object}	dto.MateriResponse
// @Failure		400		{object}	dto.ErrorResponse
// @Router			/api/materi/{id} [put]
func (h *LMSHandler) materiUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req dto.UpdateMateriRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}
	resp, err := h.materiSvc.Update(id, req)
	if err != nil {
		logger.Error(err, "materiUpdate failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	h.sseSvc.NotifyUpdate("materi", resp, getUserID(r))
	utils.RespondJSON(w, 200, resp)
}

// @Summary		Hapus materi
// @Description	Admin menghapus materi.
// @Tags			LMS - Materi
// @Produce		json
// @Security		BearerAuth
// @Param			id	path		string	true	"ID Materi"
// @Success		200	{object}	map[string]string
// @Failure		400	{object}	dto.ErrorResponse
// @Router			/api/materi/{id} [delete]
func (h *LMSHandler) materiDelete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.materiSvc.Delete(id); err != nil {
		logger.Error(err, "materiDelete failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	h.sseSvc.NotifyDelete("materi", id, getUserID(r))
	utils.RespondJSON(w, 200, map[string]string{"message": "Materi berhasil dihapus"})
}

// @Summary		Update progress materi
// @Description	User menandai progress materi (video/teks) telah diselesaikan.
// @Tags			LMS - Materi
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path		string						true	"ID Materi"
// @Param			request	body		dto.UpdateProgressRequest	true	"Data progress"
// @Success		200		{object}	dto.ProgressResponse
// @Failure		400		{object}	dto.ErrorResponse
// @Router			/api/materi/{id}/progress [post]
func (h *LMSHandler) materiUpdateProgress(w http.ResponseWriter, r *http.Request, id string) {
	userID := getUserID(r)
	if userID == "" {
		utils.RespondError(w, 401, "Unauthorized")
		return
	}
	var req dto.UpdateProgressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}
	resp, err := h.materiSvc.UpdateProgress(userID, id, req)
	if err != nil {
		logger.Error(err, "materiUpdateProgress failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 200, resp)
}

// ════════════════════════════════════════════════════════════════════════════
// FILE PENDUKUNG  —  /api/file-pendukung/{id}
// ════════════════════════════════════════════════════════════════════════════

func (h *LMSHandler) ServeFilePendukung(w http.ResponseWriter, r *http.Request) {
	id := trimID(r.URL.Path, "/api/file-pendukung")

	// /api/file-pendukung/{id}/download → download PDF
	if strings.HasSuffix(id, "/download") {
		fpID := strings.TrimSuffix(id, "/download")
		h.filePendukungDownload(w, r, fpID)
		return
	}

	switch r.Method {
	case http.MethodDelete:
		if id == "" {
			utils.RespondError(w, 400, "ID wajib")
			return
		}
		h.filePendukungDelete(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// @Summary		List file pendukung per materi
// @Description	Mendapatkan daftar file pendukung untuk suatu materi.
// @Tags			LMS - File Pendukung
// @Produce		json
// @Security		BearerAuth
// @Param			id	path		string	true	"ID Materi"
// @Success		200	{array}		dto.FilePendukungResponse
// @Failure		500	{object}	dto.ErrorResponse
// @Router			/api/materi/{id}/file-pendukung [get]
func (h *LMSHandler) filePendukungGetByMateri(w http.ResponseWriter, _ *http.Request, idMateri string) {
	data, err := h.fpSvc.GetByMateri(idMateri)
	if err != nil {
		logger.Error(err, "filePendukungGetByMateri failed")
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// @Summary		Upload file pendukung
// @Description	Admin mengupload file PDF sebagai pendukung materi.
// @Tags			LMS - File Pendukung
// @Accept			multipart/form-data
// @Produce		json
// @Security		BearerAuth
// @Param			id		path		string	true	"ID Materi"
// @Param			file	formData	file	true	"File PDF"
// @Success		201		{object}	dto.FilePendukungResponse
// @Failure		400		{object}	dto.ErrorResponse
// @Router			/api/materi/{id}/file-pendukung [post]
func (h *LMSHandler) filePendukungUpload(w http.ResponseWriter, r *http.Request, idMateri string) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.RespondError(w, 400, "Gagal membaca form-data")
		return
	}

	filePath, err := saveUploadedFile(r, "file", "uploads/file_pendukung", ".pdf")
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}
	if filePath == "" {
		utils.RespondError(w, 400, "File PDF wajib diupload")
		return
	}

	// Ambil info file
	_, header, _ := r.FormFile("file")
	namaFile := header.Filename
	ukuran := header.Size

	resp, err := h.fpSvc.Create(idMateri, namaFile, filePath, ukuran)
	if err != nil {
		logger.Error(err, "filePendukungUpload failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 201, resp)
}

// @Summary		Hapus file pendukung
// @Description	Admin menghapus file pendukung.
// @Tags			LMS - File Pendukung
// @Produce		json
// @Security		BearerAuth
// @Param			id	path		string	true	"ID File Pendukung"
// @Success		200	{object}	map[string]string
// @Failure		400	{object}	dto.ErrorResponse
// @Router			/api/file-pendukung/{id} [delete]
func (h *LMSHandler) filePendukungDelete(w http.ResponseWriter, _ *http.Request, id string) {
	if err := h.fpSvc.Delete(id); err != nil {
		logger.Error(err, "filePendukungDelete failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 200, map[string]string{"message": "File pendukung berhasil dihapus"})
}

// @Summary		Download file pendukung
// @Description	User mendownload file PDF pendukung materi.
// @Tags			LMS - File Pendukung
// @Produce		application/pdf
// @Security		BearerAuth
// @Param			id	path		string	true	"ID File Pendukung"
// @Success		200	{file}		binary
// @Failure		404	{object}	dto.ErrorResponse
// @Router			/api/file-pendukung/{id}/download [get]
func (h *LMSHandler) filePendukungDownload(w http.ResponseWriter, _ *http.Request, id string) {
	fp, err := h.fpSvc.FindByID(id)
	if err != nil {
		utils.RespondError(w, 404, "File tidak ditemukan")
		return
	}

	f, err := os.Open(fp.FilePath)
	if err != nil {
		utils.RespondError(w, 404, "File tidak ditemukan di server")
		return
	}
	defer f.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fp.NamaFile))
	w.Header().Set("Content-Type", "application/pdf")
	w.WriteHeader(http.StatusOK)
	io.Copy(w, f)
}

// ════════════════════════════════════════════════════════════════════════════
// SOAL  —  /api/kuis/{id_kuis}/soal  dan  /api/soal/{id}
// ════════════════════════════════════════════════════════════════════════════

func (h *LMSHandler) ServeSoal(w http.ResponseWriter, r *http.Request) {
	id := trimID(r.URL.Path, "/api/soal")
	switch r.Method {
	case http.MethodPut:
		if id == "" {
			utils.RespondError(w, 400, "ID wajib")
			return
		}
		h.soalUpdate(w, r, id)
	case http.MethodDelete:
		if id == "" {
			utils.RespondError(w, 400, "ID wajib")
			return
		}
		h.soalDelete(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// ServeSoalByKuis dipakai untuk GET & POST /api/kuis/{id_kuis}/soal
func (h *LMSHandler) ServeSoalByKuis(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/kuis/")
	idKuis := strings.TrimSuffix(path, "/soal")
	h.serveSoalByKuis(w, r, idKuis)
}

func (h *LMSHandler) serveSoalByKuis(w http.ResponseWriter, r *http.Request, idKuis string) {
	switch r.Method {
	case http.MethodGet:
		h.soalGetByKuis(w, r, idKuis)
	case http.MethodPost:
		h.soalCreate(w, r, idKuis)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// @Summary		List soal per kuis
// @Description	Admin mendapatkan daftar soal dalam suatu kuis.
// @Tags			LMS - Soal
// @Produce		json
// @Security		BearerAuth
// @Param			id_kuis	path		string	true	"ID Kuis"
// @Success		200		{array}		dto.SoalResponse
// @Failure		500		{object}	dto.ErrorResponse
// @Router			/api/kuis/{id_kuis}/soal [get]
func (h *LMSHandler) soalGetByKuis(w http.ResponseWriter, _ *http.Request, idKuis string) {
	data, err := h.soalSvc.GetByKuis(idKuis)
	if err != nil {
		logger.Error(err, "soalGetByKuis failed")
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// @Summary		Tambah soal ke kuis
// @Description	Admin menambahkan soal beserta pilihan jawaban ke kuis.
// @Tags			LMS - Soal
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id_kuis	path		string					true	"ID Kuis"
// @Param			request	body		dto.CreateSoalRequest	true	"Data soal + pilihan"
// @Success		201		{object}	dto.SoalResponse
// @Failure		400		{object}	dto.ErrorResponse
// @Router			/api/kuis/{id_kuis}/soal [post]
func (h *LMSHandler) soalCreate(w http.ResponseWriter, r *http.Request, idKuis string) {
	var req dto.CreateSoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}
	resp, err := h.soalSvc.Create(idKuis, req)
	if err != nil {
		logger.Error(err, "soalCreate failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 201, resp)
}

// @Summary		Update soal
// @Description	Admin mengupdate soal dan pilihan jawaban.
// @Tags			LMS - Soal
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path		string					true	"ID Soal"
// @Param			request	body		dto.UpdateSoalRequest	true	"Data update"
// @Success		200		{object}	dto.SoalResponse
// @Failure		400		{object}	dto.ErrorResponse
// @Router			/api/soal/{id} [put]
func (h *LMSHandler) soalUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req dto.UpdateSoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}
	resp, err := h.soalSvc.Update(id, req)
	if err != nil {
		logger.Error(err, "soalUpdate failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 200, resp)
}

// @Summary		Hapus soal
// @Description	Admin menghapus soal.
// @Tags			LMS - Soal
// @Produce		json
// @Security		BearerAuth
// @Param			id	path		string	true	"ID Soal"
// @Success		200	{object}	map[string]string
// @Failure		400	{object}	dto.ErrorResponse
// @Router			/api/soal/{id} [delete]
func (h *LMSHandler) soalDelete(w http.ResponseWriter, _ *http.Request, id string) {
	if err := h.soalSvc.Delete(id); err != nil {
		logger.Error(err, "soalDelete failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 200, map[string]string{"message": "Soal berhasil dihapus"})
}

// ════════════════════════════════════════════════════════════════════════════
// KUIS  —  /api/kuis (admin CRUD + user flow)
// ════════════════════════════════════════════════════════════════════════════

func (h *LMSHandler) serveKuisByKelas(w http.ResponseWriter, r *http.Request, idKelas, suffix string) {
	switch r.Method {
	case http.MethodGet:
		h.kuisGetByKelas(w, r, idKelas)
	case http.MethodPost:
		h.kuisCreate(w, r, idKelas)
	case http.MethodPut:
		if suffix == "" {
			utils.RespondError(w, 400, "ID kuis wajib untuk update")
			return
		}
		h.kuisUpdate(w, r, suffix)
	case http.MethodDelete:
		if suffix == "" {
			utils.RespondError(w, 400, "ID kuis wajib untuk delete")
			return
		}
		h.kuisDeleteHandler(w, r, suffix)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *LMSHandler) ServeKuis(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/kuis/")

	// POST /api/kuis/{id_kuis}/start
	if strings.HasSuffix(path, "/start") && r.Method == http.MethodPost {
		idKuis := strings.TrimSuffix(path, "/start")
		h.kuisStart(w, r, idKuis)
		return
	}

	// /api/kuis/{id_kuis}/soal
	if strings.HasSuffix(path, "/soal") || strings.Contains(path, "/soal") {
		idKuis := strings.Split(path, "/soal")[0]
		h.serveSoalByKuis(w, r, idKuis)
		return
	}

	// POST /api/kuis/attempt/{id}/submit
	if strings.HasPrefix(path, "attempt/") && strings.HasSuffix(path, "/submit") && r.Method == http.MethodPost {
		idAttempt := strings.TrimSuffix(strings.TrimPrefix(path, "attempt/"), "/submit")
		h.kuisSubmit(w, r, idAttempt)
		return
	}

	// GET /api/kuis/attempt/{id}/result
	if strings.HasPrefix(path, "attempt/") && strings.HasSuffix(path, "/result") && r.Method == http.MethodGet {
		idAttempt := strings.TrimSuffix(strings.TrimPrefix(path, "attempt/"), "/result")
		h.kuisResult(w, r, idAttempt)
		return
	}

	// PUT/DELETE /api/kuis/{id}
	id := trimID(r.URL.Path, "/api/kuis")
	switch r.Method {
	case http.MethodPut:
		if id == "" {
			utils.RespondError(w, 400, "ID wajib")
			return
		}
		h.kuisUpdate(w, r, id)
	case http.MethodDelete:
		if id == "" {
			utils.RespondError(w, 400, "ID wajib")
			return
		}
		h.kuisDeleteHandler(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// @Summary		List kuis per kelas
// @Description	Mendapatkan daftar kuis dalam suatu kelas beserta progress user yang login.
// @Tags			LMS - Kuis
// @Produce		json
// @Security		BearerAuth
// @Param			id	path		string	true	"ID Kelas"
// @Success		200	{array}		dto.KuisResponse
// @Failure		500	{object}	dto.ErrorResponse
// @Router			/api/kelas/{id}/kuis [get]
func (h *LMSHandler) kuisGetByKelas(w http.ResponseWriter, r *http.Request, idKelas string) {
	userID := getUserID(r)
	data, err := h.kuisSvc.GetKuisByKelas(idKelas, userID)
	if err != nil {
		logger.Error(err, "kuisGetByKelas failed")
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// @Summary		Buat kuis baru
// @Description	Admin membuat kuis baru dalam kelas. Kuis bisa per-materi atau final.
// @Tags			LMS - Kuis
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path		string					true	"ID Kelas"
// @Param			request	body		dto.CreateKuisRequest	true	"Data kuis"
// @Success		201		{object}	dto.KuisResponse
// @Failure		400		{object}	dto.ErrorResponse
// @Router			/api/kelas/{id}/kuis [post]
func (h *LMSHandler) kuisCreate(w http.ResponseWriter, r *http.Request, idKelas string) {
	var req dto.CreateKuisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}
	resp, err := h.kuisSvc.CreateKuis(idKelas, req)
	if err != nil {
		logger.Error(err, "kuisCreate failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 201, resp)
}

// @Summary		Update kuis
// @Description	Admin mengupdate data kuis.
// @Tags			LMS - Kuis
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path		string					true	"ID Kuis"
// @Param			request	body		dto.UpdateKuisRequest	true	"Data update"
// @Success		200		{object}	dto.KuisResponse
// @Failure		400		{object}	dto.ErrorResponse
// @Router			/api/kuis/{id} [put]
func (h *LMSHandler) kuisUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req dto.UpdateKuisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}
	resp, err := h.kuisSvc.UpdateKuis(id, req)
	if err != nil {
		logger.Error(err, "kuisUpdate failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 200, resp)
}

// @Summary		Hapus kuis
// @Description	Admin menghapus kuis.
// @Tags			LMS - Kuis
// @Produce		json
// @Security		BearerAuth
// @Param			id	path		string	true	"ID Kuis"
// @Success		200	{object}	map[string]string
// @Failure		400	{object}	dto.ErrorResponse
// @Router			/api/kuis/{id} [delete]
func (h *LMSHandler) kuisDeleteHandler(w http.ResponseWriter, _ *http.Request, id string) {
	if err := h.kuisSvc.DeleteKuis(id); err != nil {
		logger.Error(err, "kuisDelete failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 200, map[string]string{"message": "Kuis berhasil dihapus"})
}

// @Summary		Mulai mengerjakan kuis
// @Description	User memulai attempt kuis. Dibatasi oleh max_attempt jika ada.
// @Tags			LMS - Kuis
// @Produce		json
// @Security		BearerAuth
// @Param			id_kuis	path		string	true	"ID Kuis"
// @Success		200		{object}	dto.StartKuisResponse
// @Failure		400		{object}	dto.ErrorResponse
// @Router			/api/kuis/{id_kuis}/start [post]
func (h *LMSHandler) kuisStart(w http.ResponseWriter, r *http.Request, idKuis string) {
	userID := getUserID(r)
	if userID == "" {
		utils.RespondError(w, 401, "Unauthorized")
		return
	}
	resp, err := h.kuisSvc.Start(userID, idKuis)
	if err != nil {
		logger.Error(err, "kuisStart failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 200, resp)
}

// @Summary		Submit jawaban kuis
// @Description	User mengirim jawaban untuk attempt kuis yang sedang berjalan.
// @Tags			LMS - Kuis
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id_attempt	path		string					true	"ID Attempt"
// @Param			request		body		dto.SubmitKuisRequest	true	"Jawaban kuis"
// @Success		200			{object}	dto.KuisResultResponse
// @Failure		400			{object}	dto.ErrorResponse
// @Router			/api/kuis/attempt/{id_attempt}/submit [post]
func (h *LMSHandler) kuisSubmit(w http.ResponseWriter, r *http.Request, idAttempt string) {
	userID := getUserID(r)
	if userID == "" {
		utils.RespondError(w, 401, "Unauthorized")
		return
	}
	var req dto.SubmitKuisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}
	resp, err := h.kuisSvc.Submit(userID, idAttempt, req)
	if err != nil {
		logger.Error(err, "kuisSubmit failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 200, resp)
}

// @Summary		Lihat hasil kuis
// @Description	User melihat hasil dari attempt kuis yang sudah selesai.
// @Tags			LMS - Kuis
// @Produce		json
// @Security		BearerAuth
// @Param			id_attempt	path		string	true	"ID Attempt"
// @Success		200			{object}	dto.KuisResultResponse
// @Failure		400			{object}	dto.ErrorResponse
// @Router			/api/kuis/attempt/{id_attempt}/result [get]
func (h *LMSHandler) kuisResult(w http.ResponseWriter, r *http.Request, idAttempt string) {
	userID := getUserID(r)
	if userID == "" {
		utils.RespondError(w, 401, "Unauthorized")
		return
	}
	resp, err := h.kuisSvc.GetResult(userID, idAttempt)
	if err != nil {
		logger.Error(err, "kuisResult failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 200, resp)
}

// ════════════════════════════════════════════════════════════════════════════
// FEEDBACK  —  /api/materi/{id}/feedback
// ════════════════════════════════════════════════════════════════════════════

// @Summary		Get feedback
// @Description	User mendapatkan feedback untuk suatu materi.
// @Tags			LMS - Feedback
// @Produce		json
// @Security		BearerAuth
// @Param			id	path		string	true	"ID Materi"
// @Success		200	{object}	dto.FeedbackResponse
// @Failure		404	{object}	dto.ErrorResponse
// @Router			/api/materi/{id}/feedback [get]
func (h *LMSHandler) feedbackGet(w http.ResponseWriter, r *http.Request, idMateri string) {
	userID := getUserID(r)
	if userID == "" {
		utils.RespondError(w, 401, "Unauthorized")
		return
	}
	resp, err := h.feedbackSvc.GetByUserAndMateri(userID, idMateri)
	if err != nil {
		utils.RespondError(w, 404, err.Error())
		return
	}
	utils.RespondJSON(w, 200, resp)
}

// @Summary		Upsert feedback
// @Description	User membuat atau mengupdate feedback untuk suatu materi.
// @Tags			LMS - Feedback
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path		string						true	"ID Materi"
// @Param			request	body		dto.UpsertFeedbackRequest	true	"Data feedback"
// @Success		200		{object}	dto.FeedbackResponse
// @Failure		400		{object}	dto.ErrorResponse
// @Router			/api/materi/{id}/feedback [put]
func (h *LMSHandler) feedbackUpsert(w http.ResponseWriter, r *http.Request, idMateri string) {
	userID := getUserID(r)
	if userID == "" {
		utils.RespondError(w, 401, "Unauthorized")
		return
	}
	var req dto.UpsertFeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}
	resp, err := h.feedbackSvc.Upsert(idMateri, userID, req)
	if err != nil {
		logger.Error(err, "feedbackUpsert failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 200, resp)
}

// ════════════════════════════════════════════════════════════════════════════
// SERTIFIKAT  —  /api/sertifikat
// ════════════════════════════════════════════════════════════════════════════

func (h *LMSHandler) serveSertifikatByKelas(w http.ResponseWriter, r *http.Request, idKelas, suffix string) {
	userID := getUserID(r)
	if userID == "" {
		utils.RespondError(w, 401, "Unauthorized")
		return
	}

	if suffix == "generate" && r.Method == http.MethodPost {
		h.sertifikatGenerate(w, r, userID, idKelas)
		return
	}

	if r.Method == http.MethodGet {
		h.sertifikatGetByKelas(w, r, userID, idKelas)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *LMSHandler) ServeSertifikat(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/sertifikat")
	path = strings.TrimPrefix(path, "/")

	userID := getUserID(r)
	if userID == "" {
		utils.RespondError(w, 401, "Unauthorized")
		return
	}

	// GET /api/sertifikat/me
	if path == "me" && r.Method == http.MethodGet {
		h.sertifikatGetByUser(w, r, userID)
		return
	}

	// GET /api/sertifikat/{id}/download
	if strings.HasSuffix(path, "/download") && r.Method == http.MethodGet {
		id := strings.TrimSuffix(path, "/download")
		h.sertifikatDownload(w, r, id)
		return
	}

	// GET /api/sertifikat/{id}
	if path != "" && r.Method == http.MethodGet {
		h.sertifikatGetByID(w, r, path)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

// @Summary		Generate sertifikat
// @Description	User men-generate sertifikat setelah menyelesaikan kelas.
// @Tags			LMS - Sertifikat
// @Produce		json
// @Security		BearerAuth
// @Param			id	path		string	true	"ID Kelas"
// @Success		201	{object}	dto.SertifikatResponse
// @Failure		400	{object}	dto.ErrorResponse
// @Router			/api/kelas/{id}/sertifikat/generate [post]
func (h *LMSHandler) sertifikatGenerate(w http.ResponseWriter, _ *http.Request, userID, idKelas string) {
	resp, err := h.sertifikatSvc.Generate(userID, idKelas)
	if err != nil {
		logger.Error(err, "sertifikatGenerate failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 201, resp)
}

// @Summary		Cek sertifikat user per kelas
// @Description	Mengecek apakah user sudah punya sertifikat untuk kelas tertentu.
// @Tags			LMS - Sertifikat
// @Produce		json
// @Security		BearerAuth
// @Param			id	path		string	true	"ID Kelas"
// @Success		200	{object}	dto.SertifikatResponse
// @Failure		404	{object}	dto.ErrorResponse
// @Router			/api/kelas/{id}/sertifikat [get]
func (h *LMSHandler) sertifikatGetByKelas(w http.ResponseWriter, _ *http.Request, userID, idKelas string) {
	resp, err := h.sertifikatSvc.GetByUserAndKelas(userID, idKelas)
	if err != nil {
		utils.RespondError(w, 404, err.Error())
		return
	}
	utils.RespondJSON(w, 200, resp)
}

// @Summary		List sertifikat user
// @Description	User mendapatkan semua sertifikat miliknya.
// @Tags			LMS - Sertifikat
// @Produce		json
// @Security		BearerAuth
// @Success		200	{array}		dto.SertifikatResponse
// @Failure		500	{object}	dto.ErrorResponse
// @Router			/api/sertifikat/me [get]
func (h *LMSHandler) sertifikatGetByUser(w http.ResponseWriter, _ *http.Request, userID string) {
	data, err := h.sertifikatSvc.GetByUser(userID)
	if err != nil {
		logger.Error(err, "sertifikatGetByUser failed")
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// @Summary		Detail sertifikat
// @Description	Mendapatkan detail sertifikat berdasarkan ID.
// @Tags			LMS - Sertifikat
// @Produce		json
// @Security		BearerAuth
// @Param			id	path		string	true	"ID Sertifikat"
// @Success		200	{object}	dto.SertifikatResponse
// @Failure		404	{object}	dto.ErrorResponse
// @Router			/api/sertifikat/{id} [get]
func (h *LMSHandler) sertifikatGetByID(w http.ResponseWriter, _ *http.Request, id string) {
	resp, err := h.sertifikatSvc.GetByID(id)
	if err != nil {
		utils.RespondError(w, 404, err.Error())
		return
	}
	utils.RespondJSON(w, 200, resp)
}

// @Summary		Download sertifikat PDF
// @Description	User mendownload sertifikat dalam format PDF.
// @Tags			LMS - Sertifikat
// @Produce		application/pdf
// @Security		BearerAuth
// @Param			id	path		string	true	"ID Sertifikat"
// @Success		200	{file}		binary
// @Failure		404	{object}	dto.ErrorResponse
// @Router			/api/sertifikat/{id}/download [get]
func (h *LMSHandler) sertifikatDownload(w http.ResponseWriter, _ *http.Request, id string) {
	resp, err := h.sertifikatSvc.GetByID(id)
	if err != nil {
		utils.RespondError(w, 404, "Sertifikat tidak ditemukan")
		return
	}

	if resp.PDFPath == nil || *resp.PDFPath == "" {
		utils.RespondError(w, 404, "File sertifikat belum tersedia")
		return
	}

	f, err := os.Open(*resp.PDFPath)
	if err != nil {
		utils.RespondError(w, 404, "File sertifikat tidak ditemukan di server")
		return
	}
	defer f.Close()

	downloadName := fmt.Sprintf("Sertifikat_%s.pdf", resp.NomorSertifikat)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", downloadName))
	w.Header().Set("Content-Type", "application/pdf")
	w.WriteHeader(http.StatusOK)
	io.Copy(w, f)
}
