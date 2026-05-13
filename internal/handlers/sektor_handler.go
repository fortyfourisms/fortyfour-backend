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
var _ dto.SektorResponse

type SektorHandler struct {
	service services.SektorServiceInterface
}

func NewSektorHandler(service services.SektorServiceInterface) *SektorHandler {
	return &SektorHandler{service: service}
}

func (h *SektorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/api/sektor"), "/")

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
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// GetAllSektor godoc
//
//	@Summary		List semua sektor
//	@Description	Mengambil seluruh data sektor beserta sub sektor
//	@Tags			Sektor
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		dto.SektorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/api/sektor [get]
func (h *SektorHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetSektorByID godoc
//
//	@Summary		Ambil sektor berdasarkan ID
//	@Description	Mengambil satu data sektor beserta sub sektor
//	@Tags			Sektor
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Sektor ID"
//	@Success		200	{object}	dto.SektorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/api/sektor/{id} [get]
func (h *SektorHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}
	utils.RespondJSON(w, 200, data)
}

// CreateSektor godoc
//
//	@Summary		Buat sektor baru
//	@Description	Menambahkan data sektor baru
//	@Tags			Sektor
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.SektorRequest	true	"Data sektor"
//	@Success		201		{object}	dto.SektorResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/api/sektor [post]
func (h *SektorHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.SektorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if strings.TrimSpace(req.NamaSektor) == "" {
		utils.RespondError(w, http.StatusBadRequest, "nama_sektor tidak boleh kosong")
		return
	}

	data, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, data)
}

// UpdateSektor godoc
//
//	@Summary		Update sektor
//	@Description	Memperbarui data sektor berdasarkan ID
//	@Tags			Sektor
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"Sektor ID"
//	@Param			body	body		dto.SektorRequest	true	"Data sektor"
//	@Success		200		{object}	dto.SektorResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/api/sektor/{id} [put]
func (h *SektorHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "ID tidak boleh kosong")
		return
	}
	var req dto.SektorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if strings.TrimSpace(req.NamaSektor) == "" {
		utils.RespondError(w, http.StatusBadRequest, "nama_sektor tidak boleh kosong")
		return
	}

	data, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, data)
}
