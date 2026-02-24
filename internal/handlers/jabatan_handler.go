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

type JabatanHandler struct {
	service    *services.JabatanService
	sseService *services.SSEService
}

func NewJabatanHandler(service *services.JabatanService, sseService *services.SSEService) *JabatanHandler {
	return &JabatanHandler{
		service:    service,
		sseService: sseService,
	}
}

func (h *JabatanHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/api/jabatan"), "/")

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

// GetAllJabatan godoc
// @Summary      List semua jabatan
// @Description  Mengambil seluruh data jabatan
// @Tags         Jabatan
// @Produce      json
// @Success      200  {array}  dto.JabatanResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /api/jabatan [get]
func (h *JabatanHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetJabatanByID godoc
// @Summary      Ambil jabatan berdasarkan ID
// @Description  Mengambil satu data jabatan
// @Tags         Jabatan
// @Produce      json
// @Param        id   path      string  true  "Jabatan ID"
// @Success      200  {object} dto.JabatanResponse
// @Failure      404  {object} dto.ErrorResponse
// @Router       /api/jabatan/{id} [get]
func (h *JabatanHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}
	utils.RespondJSON(w, 200, data)
}

// CreateJabatan godoc
// @Summary      Tambah jabatan baru
// @Description  Membuat record jabatan
// @Tags         Jabatan
// @Accept       json
// @Produce      json
// @Param        jabatan body dto.CreateJabatanRequest true "Data jabatan"
// @Success      201  {object} dto.JabatanResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/jabatan [post]
func (h *JabatanHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateJabatanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Create(req)
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, err.Error())
		return
	}

	// SSE Notif Create
	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyCreate("jabatan", resp, userID)

	utils.RespondJSON(w, 201, resp)
}

// UpdateJabatan godoc
// @Summary      Update jabatan
// @Description  Mengubah data jabatan berdasarkan ID
// @Tags         Jabatan
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Jabatan ID"
// @Param        jabatan body      dto.UpdateJabatanRequest true "Data update"
// @Success      200  {object} dto.JabatanResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/jabatan/{id} [put]
func (h *JabatanHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req dto.UpdateJabatanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Update(id, req)
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, err.Error())
		return
	}

	// SSE Notif Update
	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyUpdate("jabatan", resp, userID)

	utils.RespondJSON(w, 200, resp)
}

// DeleteJabatan godoc
// @Summary      Hapus jabatan
// @Description  Menghapus data jabatan berdasarkan ID
// @Tags         Jabatan
// @Produce      json
// @Param        id  path  string  true  "Jabatan ID"
// @Success      200  {object} dto.MessageResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/jabatan/{id} [delete]
func (h *JabatanHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.Delete(id); err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, err.Error())
		return
	}

	// SSE Notif Update
	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyUpdate("jabatan", id, userID)

	utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}
