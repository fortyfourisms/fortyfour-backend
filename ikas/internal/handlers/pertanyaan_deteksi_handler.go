package handlers

import (
	"encoding/json"
	"ikas/internal/dto"
	"ikas/internal/services"
	"ikas/internal/utils"
	"net/http"

	"github.com/rollbar/rollbar-go"
)

type PertanyaanDeteksiHandler struct {
	service *services.PertanyaanDeteksiService
}

func NewPertanyaanDeteksiHandler(service *services.PertanyaanDeteksiService) *PertanyaanDeteksiHandler {
	return &PertanyaanDeteksiHandler{
		service: service,
	}
}

func (h *PertanyaanDeteksiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, _ := utils.ExtractIntID(r.URL.Path, "pertanyaan-deteksi")

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

// GetAllPertanyaanDeteksi godoc
//
//	@Summary      List semua pertanyaan deteksi
//	@Description  Mengambil seluruh data pertanyaan deteksi
//	@Tags         PertanyaanDeteksi
//	@Produce      json
//	@Success      200  {array}   dto.PertanyaanDeteksiResponse
//	@Failure      500  {object}  dto.ErrorResponse
//	@Router       /api/maturity/pertanyaan-deteksi [get]
func (h *PertanyaanDeteksiHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil mengambil data",
		"data":    data,
		"total":   len(data),
	})
}

// GetPertanyaanDeteksiByID godoc
//
//	@Summary      Ambil pertanyaan deteksi berdasarkan ID
//	@Description  Mengambil satu data pertanyaan deteksi
//	@Tags         PertanyaanDeteksi
//	@Produce      json
//	@Param        id   path      int  true  "PertanyaanDeteksi ID"
//	@Success      200  {object}  dto.PertanyaanDeteksiResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Router       /api/maturity/pertanyaan-deteksi/{id} [get]
func (h *PertanyaanDeteksiHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id int) {
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
	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil mengambil data",
		"data":    data,
	})
}

// CreatePertanyaanDeteksi godoc
//
//	@Summary      Tambah pertanyaan deteksi baru
//	@Description  Membuat record pertanyaan deteksi baru
//	@Tags         PertanyaanDeteksi
//	@Accept       json
//	@Produce      json
//	@Param        body  body      dto.CreatePertanyaanDeteksiRequest  true  "Data pertanyaan deteksi"
//	@Success      201   {object}  dto.PertanyaanDeteksiResponse
//	@Failure      400   {object}  dto.ErrorResponse
//	@Failure      404   {object}  dto.ErrorResponse
//	@Router       /api/maturity/pertanyaan-deteksi [post]
func (h *PertanyaanDeteksiHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePertanyaanDeteksiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Create(req)
	if err != nil {
		rollbar.Error(err)
		switch err.Error() {
		case "sub_kategori_id tidak boleh kosong",
			"format sub_kategori_id tidak valid",
			"ruang_lingkup_id tidak boleh kosong",
			"format ruang_lingkup_id tidak valid",
			"pertanyaan_deteksi tidak boleh kosong",
			"pertanyaan_deteksi minimal 3 karakter",
			"pertanyaan_deteksi mengandung karakter yang tidak diizinkan",
			"pertanyaan_deteksi hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&",
			"index0 mengandung karakter yang tidak diizinkan",
			"index0 hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&",
			"index1 mengandung karakter yang tidak diizinkan",
			"index1 hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&",
			"index2 mengandung karakter yang tidak diizinkan",
			"index2 hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&",
			"index3 mengandung karakter yang tidak diizinkan",
			"index3 hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&",
			"index4 mengandung karakter yang tidak diizinkan",
			"index4 hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&",
			"index5 mengandung karakter yang tidak diizinkan",
			"index5 hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&":
			utils.RespondError(w, 400, err.Error())
		case "sub_kategori_id tidak ditemukan",
			"ruang_lingkup_id tidak ditemukan":
			utils.RespondError(w, 404, err.Error())
		default:
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 201, map[string]interface{}{
		"message": "Berhasil menyimpan data",
		"data":    resp,
	})
}

// UpdatePertanyaanDeteksi godoc
//
//	@Summary      Update pertanyaan deteksi
//	@Description  Mengubah data pertanyaan deteksi berdasarkan ID
//	@Tags         PertanyaanDeteksi
//	@Accept       json
//	@Produce      json
//	@Param        id    path      int                              true  "PertanyaanDeteksi ID"
//	@Param        body  body      dto.UpdatePertanyaanDeteksiRequest  true  "Data update"
//	@Success      200   {object}  dto.PertanyaanDeteksiResponse
//	@Failure      400   {object}  dto.ErrorResponse
//	@Failure      404   {object}  dto.ErrorResponse
//	@Router       /api/maturity/pertanyaan-deteksi/{id} [put]
func (h *PertanyaanDeteksiHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id int) {
	var req dto.UpdatePertanyaanDeteksiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Update(id, req)
	if err != nil {
		rollbar.Error(err)
		switch err.Error() {
		case "data tidak ditemukan",
			"sub_kategori_id tidak ditemukan",
			"ruang_lingkup_id tidak ditemukan":
			utils.RespondError(w, 404, err.Error())
		case "format ID tidak valid",
			"sub_kategori_id tidak boleh kosong",
			"format sub_kategori_id tidak valid",
			"ruang_lingkup_id tidak boleh kosong",
			"format ruang_lingkup_id tidak valid",
			"pertanyaan_deteksi tidak boleh kosong",
			"pertanyaan_deteksi minimal 3 karakter",
			"pertanyaan_deteksi mengandung karakter yang tidak diizinkan",
			"pertanyaan_deteksi hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&",
			"index0 mengandung karakter yang tidak diizinkan",
			"index0 hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&",
			"index1 mengandung karakter yang tidak diizinkan",
			"index1 hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&",
			"index2 mengandung karakter yang tidak diizinkan",
			"index2 hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&",
			"index3 mengandung karakter yang tidak diizinkan",
			"index3 hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&",
			"index4 mengandung karakter yang tidak diizinkan",
			"index4 hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&",
			"index5 mengandung karakter yang tidak diizinkan",
			"index5 hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&":
			utils.RespondError(w, 400, err.Error())
		default:
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil memperbarui data",
		"data":    resp,
	})
}

// DeletePertanyaanDeteksi godoc
//
//	@Summary      Hapus pertanyaan deteksi
//	@Description  Menghapus data pertanyaan deteksi berdasarkan ID
//	@Tags         PertanyaanDeteksi
//	@Produce      json
//	@Param        id   path      int  true  "PertanyaanDeteksi ID"
//	@Success      200  {object}  dto.MessageResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Failure      500  {object}  dto.ErrorResponse
//	@Router       /api/maturity/pertanyaan-deteksi/{id} [delete]
func (h *PertanyaanDeteksiHandler) handleDelete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.service.Delete(id); err != nil {
		rollbar.Error(err)
		if err.Error() == "data tidak ditemukan" {
			utils.RespondError(w, 404, err.Error())
		} else {
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil menghapus data",
	})
}
