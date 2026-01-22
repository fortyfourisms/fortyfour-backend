package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"

	"github.com/rollbar/rollbar-go"
)

type PICHandler struct {
	service    *services.PICService
	sseService *services.SSEService
}

func NewPICHandler(service *services.PICService, sseService *services.SSEService) *PICHandler {
	return &PICHandler{
		service:    service,
		sseService: sseService,
	}
}

func (h *PICHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/api/pic"), "/")

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

// GetAllPIC godoc
// @Summary      List semua pic perusahaan
// @Description  Mengambil seluruh data pic perusahaan
// @Tags         PIC
// @Produce      json
// @Success      200  {array}  dto.PICResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /api/pic [get]
func (h *PICHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetPICByID godoc
// @Summary      Ambil pic perusahaan berdasarkan ID
// @Description  Mengambil satu data pic perusahaan
// @Tags         PIC
// @Produce      json
// @Param        id   path      string  true  "PIC ID"
// @Success      200  {object} dto.PICResponse
// @Failure      404  {object} dto.ErrorResponse
// @Router       /api/pic/{id} [get]
func (h *PICHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}
	utils.RespondJSON(w, 200, data)
}

// CreatePIC godoc
// @Summary      Tambah pic perusahaan baru
// @Description  Membuat record pic perusahaan
// @Tags         PIC
// @Accept       json
// @Produce      json
// @Param        pic body dto.CreatePICRequest true "Data pic perusahaan"
// @Success      201  {object} dto.PICResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/pic [post]
func (h *PICHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePICRequest
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
	h.sseService.NotifyCreate("pic", resp, userID)

	utils.RespondJSON(w, 201, resp)
}

// UpdatePIC godoc
// @Summary      Update pic perusahaan
// @Description  Mengubah data pic perusahaan berdasarkan ID
// @Tags         PIC
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "PIC ID"
// @Param        pic body      dto.UpdatePICRequest true "Data update"
// @Success      200  {object} dto.PICResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/pic/{id} [put]
func (h *PICHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req dto.UpdatePICRequest
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
	h.sseService.NotifyUpdate("pic", resp, userID)

	utils.RespondJSON(w, 200, resp)
}

// DeletePIC godoc
// @Summary      Hapus pic perusahaan
// @Description  Menghapus data pic perusahaan berdasarkan ID
// @Tags         PIC
// @Produce      json
// @Param        id  path  string  true  "PIC ID"
// @Success      200  {object} dto.MessageResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/pic/{id} [delete]
func (h *PICHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
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
	h.sseService.NotifyDelete("perusahaan", id, userID)

	utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}
