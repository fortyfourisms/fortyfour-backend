package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
)

type IkasHandler struct {
	service *services.IkasService
}

func NewIkasHandler(service *services.IkasService) *IkasHandler {
	return &IkasHandler{service: service}
}

func (h *IkasHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/api/ikas"), "/")

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

// GetAllIkas godoc
// @Summary      List semua ikas
// @Description  Mengambil seluruh data deteksi
// @Tags         Ikas
// @Produce      json
// @Success      200  {array}  dto.IkasResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /api/ikas [get]
func (h *IkasHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetIkasByID godoc
// @Summary      Ambil ikas berdasarkan ID
// @Description  Mengambil satu data ikas
// @Tags         Ikas
// @Produce      json
// @Param        id   path      string  true  "Ikas ID"
// @Success      200  {object} dto.IkasResponse
// @Failure      404  {object} dto.ErrorResponse
// @Router       /api/ikas/{id} [get]
func (h *IkasHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}
	utils.RespondJSON(w, 200, data)
}

// CreateIkas godoc
// @Summary      Tambah ikas baru
// @Description  Membuat record ikas
// @Tags         Ikas
// @Accept       json
// @Produce      json
// @Param        ikas body dto.CreateIkasRequest true "Data ikas"
// @Success      201  {object} dto.IkasResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/ikas [post]
func (h *IkasHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateIkasRequest
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

// UpdateIkas godoc
// @Summary      Update ikas
// @Description  Mengubah data ikas berdasarkan ID
// @Tags         Ikas
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Ikas ID"
// @Param        ikas body      dto.UpdateIkasRequest true "Data update"
// @Success      200  {object} dto.IkasResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/ikas/{id} [put]
func (h *IkasHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req dto.UpdateIkasRequest
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

// DeleteIkas godoc
// @Summary      Hapus ikas
// @Description  Menghapus data ikas berdasarkan ID
// @Tags         Ikas
// @Produce      json
// @Param        id  path  string  true  "Ikas ID"
// @Success      200  {object} dto.MessageResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/ikas/{id} [delete]
func (h *IkasHandler) handleDelete(w http.ResponseWriter, _ *http.Request, id string) {
	if err := h.service.Delete(id); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}
