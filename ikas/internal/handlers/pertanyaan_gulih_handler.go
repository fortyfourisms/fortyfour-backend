package handlers

import (
	"encoding/json"
	"ikas/internal/dto"
	"ikas/internal/services"
	"ikas/internal/utils"
	"net/http"

	"github.com/rollbar/rollbar-go"
)

type PertanyaanGulihHandler struct {
	service *services.PertanyaanGulihService
}

func NewPertanyaanGulihHandler(service *services.PertanyaanGulihService) *PertanyaanGulihHandler {
	return &PertanyaanGulihHandler{
		service: service,
	}
}

func (h *PertanyaanGulihHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, _ := utils.ExtractIntID(r.URL.Path, "pertanyaan-gulih")

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

// GetAllPertanyaanGulih godoc
//
//	@Summary      List semua pertanyaan gulih
//	@Description  Mengambil seluruh data pertanyaan gulih
//	@Tags         PertanyaanGulih
//	@Produce      json
//	@Success      200  {array}   dto.PertanyaanGulihResponse
//	@Failure      500  {object}  dto.ErrorResponse
//	@Router       /api/maturity/pertanyaan-gulih [get]
func (h *PertanyaanGulihHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetPertanyaanGulihByID godoc
//
//	@Summary      Ambil pertanyaan gulih berdasarkan ID
//	@Description  Mengambil satu data pertanyaan gulih
//	@Tags         PertanyaanGulih
//	@Produce      json
//	@Param        id   path      int  true  "PertanyaanGulih ID"
//	@Success      200  {object}  dto.PertanyaanGulihResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Router       /api/maturity/pertanyaan-gulih/{id} [get]
func (h *PertanyaanGulihHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id int) {
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

// CreatePertanyaanGulih godoc
//
//	@Summary      Tambah pertanyaan gulih baru
//	@Description  Membuat record pertanyaan gulih baru
//	@Tags         PertanyaanGulih
//	@Accept       json
//	@Produce      json
//	@Param        body  body      dto.CreatePertanyaanGulihRequest  true  "Data pertanyaan gulih"
//	@Success      201   {object}  dto.PertanyaanGulihResponse
//	@Failure      400   {object}  dto.ErrorResponse
//	@Failure      404   {object}  dto.ErrorResponse
//	@Router       /api/maturity/pertanyaan-gulih [post]
func (h *PertanyaanGulihHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePertanyaanGulihRequest
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
			"pertanyaan_gulih tidak boleh kosong",
			"pertanyaan_gulih minimal 3 karakter",
			"pertanyaan_gulih mengandung karakter yang tidak diizinkan",
			"pertanyaan_gulih hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&",
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

	utils.RespondJSON(w, 201, resp)
}

// UpdatePertanyaanGulih godoc
//
//	@Summary      Update pertanyaan gulih
//	@Description  Mengubah data pertanyaan gulih berdasarkan ID
//	@Tags         PertanyaanGulih
//	@Accept       json
//	@Produce      json
//	@Param        id    path      int                            true  "PertanyaanGulih ID"
//	@Param        body  body      dto.UpdatePertanyaanGulihRequest  true  "Data update"
//	@Success      200   {object}  dto.PertanyaanGulihResponse
//	@Failure      400   {object}  dto.ErrorResponse
//	@Failure      404   {object}  dto.ErrorResponse
//	@Router       /api/maturity/pertanyaan-gulih/{id} [put]
func (h *PertanyaanGulihHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id int) {
	var req dto.UpdatePertanyaanGulihRequest
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
			"pertanyaan_gulih tidak boleh kosong",
			"pertanyaan_gulih minimal 3 karakter",
			"pertanyaan_gulih mengandung karakter yang tidak diizinkan",
			"pertanyaan_gulih hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&",
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

	utils.RespondJSON(w, 200, resp)
}

// DeletePertanyaanGulih godoc
//
//	@Summary      Hapus pertanyaan gulih
//	@Description  Menghapus data pertanyaan gulih berdasarkan ID
//	@Tags         PertanyaanGulih
//	@Produce      json
//	@Param        id   path      int  true  "PertanyaanGulih ID"
//	@Success      200  {object}  dto.MessageResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Failure      500  {object}  dto.ErrorResponse
//	@Router       /api/maturity/pertanyaan-gulih/{id} [delete]
func (h *PertanyaanGulihHandler) handleDelete(w http.ResponseWriter, r *http.Request, id int) {
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
