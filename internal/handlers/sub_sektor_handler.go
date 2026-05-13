package handlers

import (
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"net/http"
	"strings"
)

// Ensure dto is available for swagger type resolution.
var _ dto.SubSektorResponse

type SubSektorHandler struct {
	service services.SubSektorServiceInterface
}

func NewSubSektorHandler(service services.SubSektorServiceInterface) *SubSektorHandler {
	return &SubSektorHandler{service: service}
}

func (h *SubSektorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/sub_sektor")

	// Check for /by_sektor/:id route
	if strings.HasPrefix(path, "/by_sektor/") {
		sektorID := strings.TrimPrefix(path, "/by_sektor/")
		if r.Method == http.MethodGet {
			h.handleGetBySektorID(w, r, sektorID)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(path, "/")

	switch r.Method {
	case http.MethodGet:
		if id == "" {
			h.handleGetAll(w, r)
		} else {
			h.handleGetByID(w, r, id)
		}
	case http.MethodPost:
		h.handleCreate(w, r)
	case http.MethodPut:
		h.handleUpdate(w, r, id)
	case http.MethodDelete:
		h.handleDelete(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// GetAllSubSektor godoc
//
//	@Summary		List semua sub sektor
//	@Description	Mengambil seluruh data sub sektor
//	@Tags			SubSektor
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		dto.SubSektorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/api/sub_sektor [get]
func (h *SubSektorHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetSubSektorByID godoc
//
//	@Summary		Ambil sub sektor berdasarkan ID
//	@Description	Mengambil satu data sub sektor
//	@Tags			SubSektor
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"SubSektor ID"
//	@Success		200	{object}	dto.SubSektorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/api/sub_sektor/{id} [get]
func (h *SubSektorHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetSubSektorBySektorID godoc
//
//	@Summary		Ambil sub sektor berdasarkan Sektor ID
//	@Description	Mengambil data sub sektor dalam satu sektor
//	@Tags			SubSektor
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Sektor ID"
//	@Success		200	{array}		dto.SubSektorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/api/sub_sektor/by_sektor/{id} [get]
func (h *SubSektorHandler) handleGetBySektorID(w http.ResponseWriter, _ *http.Request, sektorID string) {
	data, err := h.service.GetBySektorID(sektorID)
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// CreateSubSektor godoc
//
//	@Summary		Buat sub sektor baru
//	@Description	Menambahkan data sub sektor baru
//	@Tags			SubSektor
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.SubSektorRequest	true	"Data sub sektor"
//	@Success		201		{object}	dto.SubSektorResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/api/sub_sektor [post]
func (h *SubSektorHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.SubSektorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if strings.TrimSpace(req.NamaSubSektor) == "" {
		utils.RespondError(w, http.StatusBadRequest, "nama_sub_sektor tidak boleh kosong")
		return
	}
	if strings.TrimSpace(req.IDSektor) == "" {
		utils.RespondError(w, http.StatusBadRequest, "id_sektor tidak boleh kosong")
		return
	}

	data, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, data)
}

// UpdateSubSektor godoc
//
//	@Summary		Update sub sektor
//	@Description	Memperbarui data sub sektor berdasarkan ID
//	@Tags			SubSektor
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string					true	"SubSektor ID"
//	@Param			body	body		dto.SubSektorRequest	true	"Data sub sektor"
//	@Success		200		{object}	dto.SubSektorResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/api/sub_sektor/{id} [put]
func (h *SubSektorHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "ID tidak boleh kosong")
		return
	}
	var req dto.SubSektorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if strings.TrimSpace(req.NamaSubSektor) == "" {
		utils.RespondError(w, http.StatusBadRequest, "nama_sub_sektor tidak boleh kosong")
		return
	}
	if strings.TrimSpace(req.IDSektor) == "" {
		utils.RespondError(w, http.StatusBadRequest, "id_sektor tidak boleh kosong")
		return
	}

	data, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, data)
}

// DeleteSubSektor godoc
//
//	@Summary		Hapus sub sektor
//	@Description	Menghapus data sub sektor berdasarkan ID (admin only)
//	@Tags			SubSektor
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"SubSektor ID"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/api/sub_sektor/{id} [delete]
func (h *SubSektorHandler) handleDelete(w http.ResponseWriter, _ *http.Request, id string) {
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "ID tidak boleh kosong")
		return
	}
	err := h.service.Delete(id)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Sub sektor berhasil dihapus"})
}
