package handlers

import (
	"encoding/json"
	"ikas/internal/dto"
	"ikas/internal/services"
	"ikas/internal/utils"

	"net/http"

	"fortyfour-backend/pkg/logger"
)

type KategoriHandler struct {
	service *services.KategoriService
}

func NewKategoriHandler(service *services.KategoriService) *KategoriHandler {
	return &KategoriHandler{
		service: service,
	}
}

func (h *KategoriHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, _ := utils.ExtractIntID(r.URL.Path, "kategori")

	switch r.Method {
	case http.MethodGet:
		if id == 0 {
			h.handleGetAll(w, r)
		} else {
			h.handleGetByID(w, r, id)
		}
	case http.MethodPost:
		if id != 0 {
			utils.RespondError(w, 400, "ID tidak diperlukan untuk create")
			return
		}
		h.handleCreate(w, r)
	case http.MethodPut:
		if id == 0 {
			utils.RespondError(w, 400, "ID wajib")
			return
		}
		h.handleUpdate(w, r, id)
	case http.MethodDelete:
		if id == 0 {
			utils.RespondError(w, 400, "ID wajib")
			return
		}
		h.handleDelete(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// GetAllKategori godoc
//
//	@Summary      List semua kategori
//	@Description  Mengambil seluruh data kategori
//	@Tags         Kategori
//	@Produce      json
//	@Success      200  {array}   dto.KategoriResponse
//	@Failure      500  {object}  dto.ErrorResponse
//	@Router       /api/maturity/kategori [get]
func (h *KategoriHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil mengambil data",
		"data":    data,
		"total":   len(data),
	})
}

// GetKategoriByID godoc
//
//	@Summary      Ambil kategori berdasarkan ID
//	@Description  Mengambil satu data kategori
//	@Tags         Kategori
//	@Produce      json
//	@Param        id   path      int  true  "Kategori ID"
//	@Success      200  {object}  dto.KategoriResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Router       /api/maturity/kategori/{id} [get]
func (h *KategoriHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id int) {
	data, err := h.service.GetByID(id)
	if err != nil {
		logger.Error(err, "operation failed")
		if err.Error() == "data tidak ditemukan" {
			utils.RespondError(w, 404, err.Error())
		} else {
			utils.RespondError(w, 500, err.Error())
		}
		return
	}
	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil mengambil data",
		"data":    data,
	})
}

// CreateKategori godoc
//
//	@Summary      Tambah kategori baru
//	@Description  Membuat record kategori baru
//	@Tags         Kategori
//	@Accept       json
//	@Produce      json
//	@Param        kategori  body      dto.CreateKategoriRequest  true  "Data kategori"
//	@Success      201       {object}  dto.KategoriResponse
//	@Failure      400       {object}  dto.ErrorResponse
//	@Failure      409       {object}  dto.ErrorResponse
//	@Router       /api/maturity/kategori [post]
func (h *KategoriHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateKategoriRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	_, err := h.service.Create(req)
	if err != nil {
		logger.Error(err, "operation failed")
		// Mapping error ke HTTP status code
		switch err.Error() {
		case "domain_id tidak boleh kosong",
			"format domain_id tidak valid",
			"nama_kategori tidak boleh kosong",
			"nama_kategori minimal 3 karakter",
			"nama_kategori maksimal 500 karakter",
			"nama_kategori mengandung karakter yang tidak diizinkan",
			"nama_kategori hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&":
			utils.RespondError(w, 400, err.Error())
		case "domain_id tidak ditemukan":
			utils.RespondError(w, 404, err.Error())
		case "nama_kategori sudah ada dalam domain ini":
			utils.RespondError(w, 409, err.Error())
		default:
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 201, dto.KategoriMessageResponse{
		Message: "Berhasil menyimpan data",
	})
}

// UpdateKategori godoc
//
//	@Summary      Update kategori
//	@Description  Mengubah data kategori berdasarkan ID
//	@Tags         Kategori
//	@Accept       json
//	@Produce      json
//	@Param        id        path      int                    true  "Kategori ID"
//	@Param        kategori  body      dto.UpdateKategoriRequest true  "Data update"
//	@Success      200       {object}  dto.KategoriResponse
//	@Failure      400       {object}  dto.ErrorResponse
//	@Failure      404       {object}  dto.ErrorResponse
//	@Failure      409       {object}  dto.ErrorResponse
//	@Router       /api/maturity/kategori/{id} [put]
func (h *KategoriHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id int) {
	var req dto.UpdateKategoriRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	_, err := h.service.Update(id, req)
	if err != nil {
		logger.Error(err, "operation failed")
		// Mapping error ke HTTP status code
		switch err.Error() {
		case "data tidak ditemukan",
			"domain_id tidak ditemukan":
			utils.RespondError(w, 404, err.Error())
		case "domain_id tidak boleh kosong",
			"format domain_id tidak valid",
			"format ID tidak valid",
			"nama_kategori tidak boleh kosong",
			"nama_kategori minimal 3 karakter",
			"nama_kategori maksimal 500 karakter",
			"nama_kategori mengandung karakter yang tidak diizinkan",
			"nama_kategori hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&":
			utils.RespondError(w, 400, err.Error())
		case "nama_kategori sudah ada dalam domain ini":
			utils.RespondError(w, 409, err.Error())
		default:
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 200, dto.KategoriMessageResponse{
		ID:      id,
		Message: "Berhasil memperbarui data",
	})
}

// DeleteKategori godoc
//
//	@Summary      Hapus kategori
//	@Description  Menghapus data kategori berdasarkan ID
//	@Tags         Kategori
//	@Produce      json
//	@Param        id   path      int  true  "Kategori ID"
//	@Success      200  {object}  dto.KategoriMessageResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Failure      500  {object}  dto.ErrorResponse
//	@Router       /api/maturity/kategori/{id} [delete]
func (h *KategoriHandler) handleDelete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.service.Delete(id); err != nil {
		logger.Error(err, "operation failed")
		if err.Error() == "data tidak ditemukan" {
			utils.RespondError(w, 404, err.Error())
		} else {
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 200, dto.KategoriMessageResponse{
		ID:      id,
		Message: "Berhasil menghapus data",
	})
}
