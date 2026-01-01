package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
)

type GulihHandler struct {
	service    *services.GulihService
	sseService *services.SSEService
}

func NewGulihHandler(service *services.GulihService, sseService *services.SSEService) *GulihHandler {
	return &GulihHandler{
		service:    service,
		sseService: sseService,
	}
}

func (h *GulihHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/api/gulih"), "/")

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

// GetAllGulih godoc
// @Summary      List semua gulih
// @Description  Mengambil seluruh data gulih
// @Tags         Gulih
// @Produce      json
// @Success      200  {array}  dto.GulihResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /api/gulih [get]
func (h *GulihHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetGulihByID godoc
// @Summary      Ambil gulih berdasarkan ID
// @Tags         Gulih
// @Produce      json
// @Param        id   path      string  true  "Gulih ID"
// @Success      200  {object} dto.GulihResponse
// @Failure      404  {object} dto.ErrorResponse
// @Router       /api/gulih/{id} [get]
func (h *GulihHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}
	utils.RespondJSON(w, 200, data)
}

// CreateGulih godoc
// @Summary      Tambah gulih baru
// @Description  Membuat record gulih
// @Tags         Gulih
// @Accept       json
// @Produce      json
// @Param        deteksi body dto.CreateGulihRequest true "Data gulih"
// @Success      201  {object} dto.GulihResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/gulih [post]
func (h *GulihHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateGulihRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
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
	h.sseService.NotifyCreate("gulih", resp, userID)

	utils.RespondJSON(w, 201, resp)
}

// UpdateGulih godoc
// @Summary      Update gulih
// @Description  Mengubah data gulih berdasarkan ID
// @Tags         Gulih
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Gulih ID"
// @Param        deteksi body      dto.UpdateGulihRequest true "Data update"
// @Success      200  {object} dto.DeteksiResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/gulih/{id} [put]
func (h *GulihHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req dto.UpdateGulihRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
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
	h.sseService.NotifyUpdate("gulih", resp, userID)

	utils.RespondJSON(w, 200, resp)
}

// DeleteGulih godoc
// @Summary      Hapus gulih
// @Description  Menghapus data gulih berdasarkan ID
// @Tags         Gulih
// @Produce      json
// @Param        id  path  string  true  "Gulih ID"
// @Success      200  {object} dto.MessageResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/gulih/{id} [delete]
func (h *GulihHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.Delete(id); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	// SSE Notif Delete
	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyDelete("gulih", id, userID)

	utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}
