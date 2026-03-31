package handlers

import (
	"encoding/json"
	"ikas/internal/dto"
	"ikas/internal/services"
	"ikas/internal/utils"
	"net/http"

	"github.com/rollbar/rollbar-go"
)

type PertanyaanProteksiHandler struct {
	service *services.PertanyaanProteksiService
}

func NewPertanyaanProteksiHandler(service *services.PertanyaanProteksiService) *PertanyaanProteksiHandler {
	return &PertanyaanProteksiHandler{
		service: service,
	}
}

func (h *PertanyaanProteksiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, _ := utils.ExtractIntID(r.URL.Path, "pertanyaan-proteksi")

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

// GetAllPertanyaanProteksi godoc
//
//	@Summary      List semua pertanyaan proteksi
//	@Description  Mengambil seluruh data pertanyaan proteksi
//	@Tags         PertanyaanProteksi
//	@Produce      json
//	@Success      200  {array}   dto.PertanyaanProteksiResponse
//	@Failure      500  {object}  dto.ErrorResponse
//	@Router       /api/maturity/pertanyaan-proteksi [get]
func (h *PertanyaanProteksiHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
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

// GetPertanyaanProteksiByID godoc
//
//	@Summary      Ambil pertanyaan proteksi berdasarkan ID
//	@Description  Mengambil satu data pertanyaan proteksi
//	@Tags         PertanyaanProteksi
//	@Produce      json
//	@Param        id   path      int  true  "PertanyaanProteksi ID"
//	@Success      200  {object}  dto.PertanyaanProteksiResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Router       /api/maturity/pertanyaan-proteksi/{id} [get]
func (h *PertanyaanProteksiHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id int) {
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

// CreatePertanyaanProteksi godoc
//
//	@Summary      Tambah pertanyaan proteksi baru
//	@Description  Membuat record pertanyaan proteksi baru
//	@Tags         PertanyaanProteksi
//	@Accept       json
//	@Produce      json
//	@Param        body  body      dto.CreatePertanyaanProteksiRequest  true  "Data pertanyaan proteksi"
//	@Success      201   {object}  dto.PertanyaanProteksiResponse
//	@Failure      400   {object}  dto.ErrorResponse
//	@Failure      404   {object}  dto.ErrorResponse
//	@Router       /api/maturity/pertanyaan-proteksi [post]
func (h *PertanyaanProteksiHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePertanyaanProteksiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	_, err := h.service.Create(req)
	if err != nil {
		rollbar.Error(err)
		switch err.Error() {
		case "sub_kategori_id tidak boleh kosong",
			"format sub_kategori_id tidak valid",
			"ruang_lingkup_id tidak boleh kosong",
			"format ruang_lingkup_id tidak valid",
			"pertanyaan_proteksi tidak boleh kosong",
			"pertanyaan_proteksi minimal 3 karakter",
			"pertanyaan_proteksi mengandung karakter yang tidak diizinkan",
			"pertanyaan_proteksi hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&",
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

	utils.RespondJSON(w, 201, dto.PertanyaanProteksiMessageResponse{
		Message: "Berhasil menyimpan data",
	})
}

// UpdatePertanyaanProteksi godoc
//
//	@Summary      Update pertanyaan proteksi
//	@Description  Mengubah data pertanyaan proteksi berdasarkan ID
//	@Tags         PertanyaanProteksi
//	@Accept       json
//	@Produce      json
//	@Param        id    path      int                               true  "PertanyaanProteksi ID"
//	@Param        body  body      dto.UpdatePertanyaanProteksiRequest  true  "Data update"
//	@Success      200   {object}  dto.PertanyaanProteksiResponse
//	@Failure      400   {object}  dto.ErrorResponse
//	@Failure      404   {object}  dto.ErrorResponse
//	@Router       /api/maturity/pertanyaan-proteksi/{id} [put]
func (h *PertanyaanProteksiHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id int) {
	var req dto.UpdatePertanyaanProteksiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	_, err := h.service.Update(id, req)
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
			"pertanyaan_proteksi tidak boleh kosong",
			"pertanyaan_proteksi minimal 3 karakter",
			"pertanyaan_proteksi mengandung karakter yang tidak diizinkan",
			"pertanyaan_proteksi hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&",
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

	utils.RespondJSON(w, 200, dto.PertanyaanProteksiMessageResponse{
		ID:      id,
		Message: "Berhasil memperbarui data",
	})
}

// DeletePertanyaanProteksi godoc
//
//	@Summary      Hapus pertanyaan proteksi
//	@Description  Menghapus data pertanyaan proteksi berdasarkan ID
//	@Tags         PertanyaanProteksi
//	@Produce      json
//	@Param        id   path      int  true  "PertanyaanProteksi ID"
//	@Success      200  {object}  dto.PertanyaanProteksiMessageResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Failure      500  {object}  dto.ErrorResponse
//	@Router       /api/maturity/pertanyaan-proteksi/{id} [delete]
func (h *PertanyaanProteksiHandler) handleDelete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.service.Delete(id); err != nil {
		rollbar.Error(err)
		if err.Error() == "data tidak ditemukan" {
			utils.RespondError(w, 404, err.Error())
		} else {
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 200, dto.PertanyaanProteksiMessageResponse{
		ID:      id,
		Message: "Berhasil menghapus data",
	})
}
