package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"fortyfour-backend/pkg/logger"
)

// LMSHandler menangani semua endpoint LMS:
//   /api/kelas          → kelas (admin CRUD, user read)
//   /api/materi         → materi (admin CRUD, user progress)
//   /api/soal           → soal kuis (admin CRUD)
//   /api/kuis           → kuis (user: start & submit)

type LMSHandler struct {
	kelasSvc *services.KelasService
	materiSvc *services.MateriService
	soalSvc  *services.SoalService
	kuisSvc  *services.KuisService
	sseSvc   *services.SSEService
}

func NewLMSHandler(
	kelasSvc *services.KelasService,
	materiSvc *services.MateriService,
	soalSvc *services.SoalService,
	kuisSvc *services.KuisService,
	sseSvc *services.SSEService,
) *LMSHandler {
	return &LMSHandler{
		kelasSvc:  kelasSvc,
		materiSvc: materiSvc,
		soalSvc:   soalSvc,
		kuisSvc:   kuisSvc,
		sseSvc:    sseSvc,
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
// KELAS  —  /api/kelas  dan  /api/kelas/{id}
// ════════════════════════════════════════════════════════════════════════════

// ServeKelas godoc
// @Tags LMS-Kelas
// @Router /api/kelas [get]
// @Router /api/kelas/{id} [get]
// @Router /api/kelas [post]
// @Router /api/kelas/{id} [put]
// @Router /api/kelas/{id} [delete]
func (h *LMSHandler) ServeKelas(w http.ResponseWriter, r *http.Request) {
	id := trimID(r.URL.Path, "/api/kelas")

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

func (h *LMSHandler) kelasGetAll(w http.ResponseWriter, r *http.Request) {
	// Admin melihat semua (draft+published), user hanya published
	userID := getUserID(r)
	onlyPublished := userID == "" // jika tidak ada userID, fallback published
	// Lebih tepatnya bisa cek role dari context, tapi untuk simplisitas:
	// handler ini akan dipanggil dengan middleware yang berbeda untuk admin vs user

	data, err := h.kelasSvc.GetAll(onlyPublished)
	if err != nil {
		logger.Error(err, "kelasGetAll failed")
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

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
// MATERI  —  /api/kelas/{id_kelas}/materi  dan  /api/materi/{id}
//
// Pattern URL:
//   POST   /api/kelas/{id_kelas}/materi          → tambah materi ke kelas
//   PUT    /api/materi/{id}                       → update materi
//   DELETE /api/materi/{id}                       → hapus materi
//   POST   /api/materi/{id}/progress              → update progress user
// ════════════════════════════════════════════════════════════════════════════

func (h *LMSHandler) ServeMateri(w http.ResponseWriter, r *http.Request) {
	// /api/materi/{id}/progress
	path := r.URL.Path
	if strings.HasSuffix(path, "/progress") {
		id := trimID(strings.TrimSuffix(path, "/progress"), "/api/materi")
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		h.materiUpdateProgress(w, r, id)
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
	// Ekstrak id_kelas dari path: /api/kelas/{id_kelas}/materi
	path := strings.TrimPrefix(r.URL.Path, "/api/kelas/")
	idKelas := strings.TrimSuffix(path, "/materi")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	h.materiCreate(w, r, idKelas)
}

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

func (h *LMSHandler) materiDelete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.materiSvc.Delete(id); err != nil {
		logger.Error(err, "materiDelete failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	h.sseSvc.NotifyDelete("materi", id, getUserID(r))
	utils.RespondJSON(w, 200, map[string]string{"message": "Materi berhasil dihapus"})
}

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
// SOAL  —  /api/materi/{id_materi}/soal  dan  /api/soal/{id}
//
// Pattern URL:
//   POST   /api/materi/{id_materi}/soal   → tambah soal ke kuis
//   GET    /api/materi/{id_materi}/soal   → list soal (admin)
//   PUT    /api/soal/{id}                 → update soal
//   DELETE /api/soal/{id}                 → hapus soal
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

// ServeSoalByMateri dipakai untuk GET & POST /api/materi/{id_materi}/soal
func (h *LMSHandler) ServeSoalByMateri(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/materi/")
	idMateri := strings.TrimSuffix(path, "/soal")

	switch r.Method {
	case http.MethodGet:
		h.soalGetByMateri(w, r, idMateri)
	case http.MethodPost:
		h.soalCreate(w, r, idMateri)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *LMSHandler) soalGetByMateri(w http.ResponseWriter, _ *http.Request, idMateri string) {
	data, err := h.soalSvc.GetByMateri(idMateri)
	if err != nil {
		logger.Error(err, "soalGetByMateri failed")
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

func (h *LMSHandler) soalCreate(w http.ResponseWriter, r *http.Request, idMateri string) {
	var req dto.CreateSoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}
	resp, err := h.soalSvc.Create(idMateri, req)
	if err != nil {
		logger.Error(err, "soalCreate failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 201, resp)
}

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

func (h *LMSHandler) soalDelete(w http.ResponseWriter, _ *http.Request, id string) {
	if err := h.soalSvc.Delete(id); err != nil {
		logger.Error(err, "soalDelete failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 200, map[string]string{"message": "Soal berhasil dihapus"})
}

// ════════════════════════════════════════════════════════════════════════════
// KUIS  —  /api/kuis
//
// Pattern URL:
//   POST /api/kuis/{id_materi}/start          → mulai kuis
//   POST /api/kuis/attempt/{id_attempt}/submit → submit jawaban
//   GET  /api/kuis/attempt/{id_attempt}/result → lihat hasil
// ════════════════════════════════════════════════════════════════════════════

func (h *LMSHandler) ServeKuis(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/kuis/")

	switch {
	// POST /api/kuis/{id_materi}/start
	case strings.HasSuffix(path, "/start") && r.Method == http.MethodPost:
		idMateri := strings.TrimSuffix(path, "/start")
		h.kuisStart(w, r, idMateri)

	// POST /api/kuis/attempt/{id}/submit
	case strings.HasPrefix(path, "attempt/") && strings.HasSuffix(path, "/submit") && r.Method == http.MethodPost:
		idAttempt := strings.TrimSuffix(strings.TrimPrefix(path, "attempt/"), "/submit")
		h.kuisSubmit(w, r, idAttempt)

	// GET /api/kuis/attempt/{id}/result
	case strings.HasPrefix(path, "attempt/") && strings.HasSuffix(path, "/result") && r.Method == http.MethodGet:
		idAttempt := strings.TrimSuffix(strings.TrimPrefix(path, "attempt/"), "/result")
		h.kuisResult(w, r, idAttempt)

	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (h *LMSHandler) kuisStart(w http.ResponseWriter, r *http.Request, idMateri string) {
	userID := getUserID(r)
	if userID == "" {
		utils.RespondError(w, 401, "Unauthorized")
		return
	}
	resp, err := h.kuisSvc.Start(userID, idMateri)
	if err != nil {
		logger.Error(err, "kuisStart failed")
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 200, resp)
}

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