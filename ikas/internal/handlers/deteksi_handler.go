package handlers

import (
	"encoding/json"
	"ikas/internal/dto"
	"ikas/internal/services"
	"ikas/internal/utils"
	"net/http"

	"fortyfour-backend/pkg/logger"
)

type DeteksiHandler struct {
	service *services.DeteksiService
}

func NewDeteksiHandler(service *services.DeteksiService) *DeteksiHandler {
	return &DeteksiHandler{
		service: service,
	}
}

func (h *DeteksiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "deteksi")

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

// GetAllDeteksi godoc
// @Summary      List semua deteksi
// @Description  Mengambil seluruh data deteksi
// @Tags         Deteksi
// @Produce      json
// @Success      200  {array}  dto.DeteksiResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /api/deteksi [get]
func (h *DeteksiHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetDeteksiByID godoc
// @Summary      Ambil deteksi berdasarkan ID
// @Tags         Deteksi
// @Produce      json
// @Param        id   path      string  true  "Deteksi ID"
// @Success      200  {object} dto.DeteksiResponse
// @Failure      404  {object} dto.ErrorResponse
// @Router       /api/deteksi/{id} [get]
func (h *DeteksiHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}
	utils.RespondJSON(w, 200, data)
}

// CreateDeteksi godoc
// @Summary      Tambah deteksi baru
// @Description  Membuat record deteksi
// @Tags         Deteksi
// @Accept       json
// @Produce      json
// @Param        deteksi body dto.CreateDeteksiRequest true "Data deteksi"
// @Success      201  {object} dto.DeteksiResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/deteksi [post]
func (h *DeteksiHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateDeteksiRequest
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

	utils.RespondJSON(w, 201, resp)
}

// UpdateDeteksi godoc
// @Summary      Update deteksi
// @Description  Mengubah data deteksi berdasarkan ID
// @Tags         Deteksi
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Deteksi ID"
// @Param        deteksi body      dto.UpdateDeteksiRequest true "Data update"
// @Success      200  {object} dto.DeteksiResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/deteksi/{id} [put]
func (h *DeteksiHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req dto.UpdateDeteksiRequest
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

	utils.RespondJSON(w, 200, resp)
}

// DeleteDeteksi godoc
// @Summary      Hapus deteksi
// @Description  Menghapus data deteksi berdasarkan ID
// @Tags         Deteksi
// @Produce      json
// @Param        id  path  string  true  "Deteksi ID"
// @Success      200  {object} dto.MessageResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/deteksi/{id} [delete]
func (h *DeteksiHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.Delete(id); err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}
