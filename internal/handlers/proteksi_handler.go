package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
)

type ProteksiHandler struct {
	service *services.ProteksiService
}

func NewProteksiHandler(service *services.ProteksiService) *ProteksiHandler {
	return &ProteksiHandler{service: service}
}

func (h *ProteksiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/api/proteksi"), "/")

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

// GetAllProteksi godoc
// @Summary      List semua proteksi
// @Description  Mengambil seluruh data proteksi
// @Tags         Proteksi
// @Produce      json
// @Success      200  {array}  dto.ProteksiResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /api/proteksi [get]
func (h *ProteksiHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetProteksiByID godoc
// @Summary      Ambil proteksi berdasarkan ID
// @Description  Mengambil satu data proteksi
// @Tags         Proteksi
// @Produce      json
// @Param        id   path      string  true  "Proteksi ID"
// @Success      200  {object} dto.ProteksiResponse
// @Failure      404  {object} dto.ErrorResponse
// @Router       /api/proteksi/{id} [get]
func (h *ProteksiHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}
	utils.RespondJSON(w, 200, data)
}

// CreateProteksi godoc
// @Summary      Tambah proteksi baru
// @Description  Membuat record proteksi
// @Tags         Proteksi
// @Accept       json
// @Produce      json
// @Param        proteksi body dto.CreateProteksiRequest true "Data proteksi"
// @Success      201  {object} dto.ProteksiResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/proteksi [post]
func (h *ProteksiHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProteksiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 201, resp)
}

// UpdateProteksi godoc
// @Summary      Update proteksi
// @Description  Mengubah data proteksi berdasarkan ID
// @Tags         Proteksi
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Proteksi ID"
// @Param        proteksi body      dto.UpdateProteksiRequest true "Data update"
// @Success      200  {object} dto.ProteksiResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/proteksi/{id} [put]
func (h *ProteksiHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req dto.UpdateProteksiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 200, resp)
}

// DeleteProteksi godoc
// @Summary      Hapus proteksi
// @Description  Menghapus data proteksi berdasarkan ID
// @Tags         Proteksi
// @Produce      json
// @Param        id  path  string  true  "Proteksi ID"
// @Success      200  {object} dto.MessageResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/proteksi/{id} [delete]
func (h *ProteksiHandler) handleDelete(w http.ResponseWriter, _ *http.Request, id string) {
	if err := h.service.Delete(id); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}
