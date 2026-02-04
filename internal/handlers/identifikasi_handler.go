package handlers

import (
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"net/http"
	"strings"

	"github.com/rollbar/rollbar-go"
)

type IdentifikasiHandler struct {
	service    *services.IdentifikasiService
	sseService *services.SSEService
}

func NewIdentifikasiHandler(service *services.IdentifikasiService, sseService *services.SSEService) *IdentifikasiHandler {
	return &IdentifikasiHandler{
		service:    service,
		sseService: sseService,
	}
}

func (h *IdentifikasiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/api/identifikasi"), "/")

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

// CreateIdentifikasi godoc
// @Summary      Tambah identifikasi baru
// @Description  Membuat record identifikasi
// @Tags         Identifikasi
// @Accept       json
// @Produce      json
// @Param        identifikasi body dto.CreateIdentifikasiRequest true "Data identifikasi"
// @Success      201  {object} dto.IdentifikasiResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/identifikasi [post]
func (h *IdentifikasiHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateIdentifikasiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Create(req)
	if err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, err.Error())
		return
	}

	// SSE Notif Create
	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyCreate("identifikasi", resp, userID)

	utils.RespondJSON(w, 201, resp)
}

// GetAllIdentifikasi godoc
// @Summary      List semua identifikasi
// @Description  Mengambil seluruh data identifikasi
// @Tags         Identifikasi
// @Produce      json
// @Success      200  {array}  dto.IdentifikasiResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /api/identifikasi [get]
func (h *IdentifikasiHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 500, err.Error())
		return
	}

	utils.RespondJSON(w, 200, data)
}

// GetIdentifikasiByID godoc
// @Summary      Ambil identifikasi berdasarkan ID
// @Description  Mengambil satu data identifikasi
// @Tags         Identifikasi
// @Produce      json
// @Param        id   path      string  true  "Identifikasi ID"
// @Success      200  {object} dto.IdentifikasiResponse
// @Failure      404  {object} dto.ErrorResponse
// @Router       /api/identifikasi/{id} [get]
func (h *IdentifikasiHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}
	utils.RespondJSON(w, 200, data)
}

// UpdateIdentifikasi godoc
// @Summary      Update identifikasi
// @Description  Mengubah data identifikasi berdasarkan ID
// @Tags         Identifikasi
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Identifikasi ID"
// @Param        identifikasi body      dto.UpdateIdentifikasiRequest true "Data update"
// @Success      200  {object} dto.IdentifikasiResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/identifikasi/{id} [put]
func (h *IdentifikasiHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req dto.UpdateIdentifikasiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Update(id, req)
	if err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, err.Error())
		return
	}

	// SSE Notif Update
	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyUpdate("identifikasi", resp, userID)

	utils.RespondJSON(w, 200, resp)
}

// DeleteIdentifikasi godoc
// @Summary      Hapus identifikasi
// @Description  Menghapus data identifikasi berdasarkan ID
// @Tags         Identifikasi
// @Produce      json
// @Param        id  path  string  true  "Identifikasi ID"
// @Success      200  {object} dto.MessageResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/identifikasi/{id} [delete]
func (h *IdentifikasiHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.Delete(id); err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, err.Error())
		return
	}

	// SSE Notif Delete
	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyDelete("identifikasi", id, userID)

	utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}
