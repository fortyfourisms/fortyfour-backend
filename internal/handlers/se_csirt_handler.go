package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
)

type SeCsirtHandler struct {
	service *services.SeCsirtService
}

func NewSeCsirtHandler(service *services.SeCsirtService) *SeCsirtHandler {
	return &SeCsirtHandler{service: service}
}

func (h *SeCsirtHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/se_csirt")
	id := strings.Trim(path, "/")

	switch r.Method {
	case http.MethodGet:
		if id == "" {
			h.handleGetAll(w)
		} else {
			h.handleGetByID(w, id)
		}
	case http.MethodPost:
		h.handleCreate(w, r)
	case http.MethodPut:
		h.handleUpdate(w, r, id)
	case http.MethodDelete:
		h.handleDelete(w, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// GetAllSE godoc
// @Summary      List semua se csirt
// @Description  Mengambil seluruh data se csirt
// @Tags         SE
// @Produce      json
// @Success      200  {array}  dto.SeCsirtResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /api/se_csirt [get]
func (h *SeCsirtHandler) handleGetAll(w http.ResponseWriter) {
	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetSEByID godoc
// @Summary      Ambil se csirt berdasarkan ID
// @Description  Mengambil satu data se csirt
// @Tags         SE
// @Produce      json
// @Param        id   path      string  true  "SE ID"
// @Success      200  {object} dto.SeCsirtResponse
// @Failure      404  {object} dto.ErrorResponse
// @Router       /api/se_csirt/{id} [get]
func (h *SeCsirtHandler) handleGetByID(w http.ResponseWriter, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, 404, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// CreateSE godoc
// @Summary      Tambah se csirt baru
// @Description  Membuat record se csirt
// @Tags         SE
// @Accept       json
// @Produce      json
// @Param        se body dto.CreateSeCsirtRequest true "Data se csirt"
// @Success      201  {object} dto.SeCsirtResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/se_csirt [post]
func (h *SeCsirtHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSeCsirtRequest
	json.NewDecoder(r.Body).Decode(&req)

	id, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 201, map[string]string{"id": id})
}

// UpdateSE godoc
// @Summary      Update se csirt
// @Description  Mengubah data se csirt berdasarkan ID
// @Tags         SE
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "SE ID"
// @Param        se body      dto.UpdateSeCsirtRequest true "Data update"
// @Success      200  {object} dto.SeCsirtResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/se_csirt/{id} [put]
func (h *SeCsirtHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req dto.UpdateSeCsirtRequest
	json.NewDecoder(r.Body).Decode(&req)

	if err := h.service.Update(id, req); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 200, map[string]string{"message": "Update success"})
}

// DeleteSE godoc
// @Summary      Hapus se csirt
// @Description  Menghapus data se csirt berdasarkan ID
// @Tags         SE
// @Produce      json
// @Param        id  path  string  true  "SE ID"
// @Success      200  {object} dto.MessageResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/se_csirt/{id} [delete]
func (h *SeCsirtHandler) handleDelete(w http.ResponseWriter, id string) {
	if err := h.service.Delete(id); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}
