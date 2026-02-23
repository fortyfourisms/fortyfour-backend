package handlers

import (
	"encoding/json"
	"ikas/internal/dto"
	"ikas/internal/services"
	"ikas/internal/utils"
	"net/http"
	"strings"

	"github.com/rollbar/rollbar-go"
)

type SubKategoriHandler struct {
	service *services.SubKategoriService
}

func NewSubKategoriHandler(service *services.SubKategoriService) *SubKategoriHandler {
	return &SubKategoriHandler{
		service: service,
	}
}

func (h *SubKategoriHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/sub-kategori")
	id := strings.TrimPrefix(path, "/")

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

// GetAllSubKategori godoc
//
//	@Summary      List semua sub kategori
//	@Description  Mengambil seluruh data sub kategori
//	@Tags         SubKategori
//	@Produce      json
//	@Success      200  {array}   dto.SubKategoriResponse
//	@Failure      500  {object}  dto.ErrorResponse
//	@Router       /api/sub-kategori [get]
func (h *SubKategoriHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetSubKategoriByID godoc
//
//	@Summary      Ambil sub kategori berdasarkan ID
//	@Description  Mengambil satu data sub kategori
//	@Tags         SubKategori
//	@Produce      json
//	@Param        id   path      string  true  "SubKategori ID"
//	@Success      200  {object}  dto.SubKategoriResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Router       /api/sub-kategori/{id} [get]
func (h *SubKategoriHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		if err.Error() == "data tidak ditemukan" {
			utils.RespondError(w, 404, err.Error())
		} else {
			utils.RespondError(w, 500, err.Error())
		}
		return
	}
	utils.RespondJSON(w, 200, data)
}

// CreateSubKategori godoc
//
//	@Summary      Tambah sub kategori baru
//	@Description  Membuat record sub kategori baru
//	@Tags         SubKategori
//	@Accept       json
//	@Produce      json
//	@Param        subKategori  body      dto.CreateSubKategoriRequest  true  "Data sub kategori"
//	@Success      201          {object}  dto.SubKategoriResponse
//	@Failure      400          {object}  dto.ErrorResponse
//	@Failure      409          {object}  dto.ErrorResponse
//	@Router       /api/sub-kategori [post]
func (h *SubKategoriHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSubKategoriRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Create(req)
	if err != nil {
		rollbar.Error(err)
		// Mapping error ke HTTP status code
		switch err.Error() {
		case "kategori_id tidak boleh kosong",
			"format kategori_id tidak valid",
			"nama_sub_kategori tidak boleh kosong",
			"nama_sub_kategori minimal 3 karakter",
			"nama_sub_kategori maksimal 500 karakter",
			"nama_sub_kategori mengandung karakter yang tidak diizinkan",
			"nama_sub_kategori hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&":
			utils.RespondError(w, 400, err.Error())
		case "kategori_id tidak ditemukan":
			utils.RespondError(w, 404, err.Error())
		case "nama_sub_kategori sudah ada dalam kategori ini":
			utils.RespondError(w, 409, err.Error())
		default:
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 201, resp)
}

// UpdateSubKategori godoc
//
//	@Summary      Update sub kategori
//	@Description  Mengubah data sub kategori berdasarkan ID
//	@Tags         SubKategori
//	@Accept       json
//	@Produce      json
//	@Param        id           path      string                       true  "SubKategori ID"
//	@Param        subKategori  body      dto.UpdateSubKategoriRequest true  "Data update"
//	@Success      200          {object}  dto.SubKategoriResponse
//	@Failure      400          {object}  dto.ErrorResponse
//	@Failure      404          {object}  dto.ErrorResponse
//	@Failure      409          {object}  dto.ErrorResponse
//	@Router       /api/sub-kategori/{id} [put]
func (h *SubKategoriHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req dto.UpdateSubKategoriRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Update(id, req)
	if err != nil {
		rollbar.Error(err)
		// Mapping error ke HTTP status code
		switch err.Error() {
		case "data tidak ditemukan",
			"kategori_id tidak ditemukan":
			utils.RespondError(w, 404, err.Error())
		case "kategori_id tidak boleh kosong",
			"format kategori_id tidak valid",
			"format ID tidak valid",
			"nama_sub_kategori tidak boleh kosong",
			"nama_sub_kategori minimal 3 karakter",
			"nama_sub_kategori maksimal 500 karakter",
			"nama_sub_kategori mengandung karakter yang tidak diizinkan",
			"nama_sub_kategori hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&":
			utils.RespondError(w, 400, err.Error())
		case "nama_sub_kategori sudah ada dalam kategori ini":
			utils.RespondError(w, 409, err.Error())
		default:
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 200, resp)
}

// DeleteSubKategori godoc
//
//	@Summary      Hapus sub kategori
//	@Description  Menghapus data sub kategori berdasarkan ID
//	@Tags         SubKategori
//	@Produce      json
//	@Param        id   path      string  true  "SubKategori ID"
//	@Success      200  {object}  dto.MessageResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Failure      500  {object}  dto.ErrorResponse
//	@Router       /api/sub-kategori/{id} [delete]
func (h *SubKategoriHandler) handleDelete(w http.ResponseWriter, _ *http.Request, id string) {
	if err := h.service.Delete(id); err != nil {
		rollbar.Error(err)
		if err.Error() == "data tidak ditemukan" {
			utils.RespondError(w, 404, err.Error())
		} else {
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}
