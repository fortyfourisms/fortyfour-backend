package handlers

import (
	"encoding/json"
	"ikas/internal/dto"
	"ikas/internal/services"
	"ikas/internal/utils"
	"net/http"

	"github.com/rollbar/rollbar-go"
)

type PertanyaanIdentifikasiHandler struct {
	service *services.PertanyaanIdentifikasiService
}

func NewPertanyaanIdentifikasiHandler(service *services.PertanyaanIdentifikasiService) *PertanyaanIdentifikasiHandler {
	return &PertanyaanIdentifikasiHandler{
		service: service,
	}
}

func (h *PertanyaanIdentifikasiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, _ := utils.ExtractIntID(r.URL.Path, "pertanyaan-identifikasi")

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

// GetAllPertanyaanIdentifikasi godoc
//
//	@Summary      List semua pertanyaan identifikasi
//	@Description  Mengambil seluruh data pertanyaan identifikasi
//	@Tags         PertanyaanIdentifikasi
//	@Produce      json
//	@Success      200  {array}   dto.PertanyaanIdentifikasiResponse
//	@Failure      500  {object}  dto.ErrorResponse
//	@Router       /api/maturity/pertanyaan-identifikasi [get]
func (h *PertanyaanIdentifikasiHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetPertanyaanIdentifikasiByID godoc
//
//	@Summary      Ambil pertanyaan identifikasi berdasarkan ID
//	@Description  Mengambil satu data pertanyaan identifikasi
//	@Tags         PertanyaanIdentifikasi
//	@Produce      json
//	@Param        id   path      int  true  "PertanyaanIdentifikasi ID"
//	@Success      200  {object}  dto.PertanyaanIdentifikasiResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Router       /api/maturity/pertanyaan-identifikasi/{id} [get]
func (h *PertanyaanIdentifikasiHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id int) {
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

// CreatePertanyaanIdentifikasi godoc
//
//	@Summary      Tambah pertanyaan identifikasi baru
//	@Description  Membuat record pertanyaan identifikasi baru
//	@Tags         PertanyaanIdentifikasi
//	@Accept       json
//	@Produce      json
//	@Param        body  body      dto.CreatePertanyaanIdentifikasiRequest  true  "Data pertanyaan identifikasi"
//	@Success      201   {object}  dto.PertanyaanIdentifikasiResponse
//	@Failure      400   {object}  dto.ErrorResponse
//	@Failure      404   {object}  dto.ErrorResponse
//	@Router       /api/maturity/pertanyaan-identifikasi [post]
func (h *PertanyaanIdentifikasiHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePertanyaanIdentifikasiRequest
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
			"pertanyaan_identifikasi tidak boleh kosong",
			"pertanyaan_identifikasi minimal 3 karakter":
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

// UpdatePertanyaanIdentifikasi godoc
//
//	@Summary      Update pertanyaan identifikasi
//	@Description  Mengubah data pertanyaan identifikasi berdasarkan ID
//	@Tags         PertanyaanIdentifikasi
//	@Accept       json
//	@Produce      json
//	@Param        id    path      int                                   true  "PertanyaanIdentifikasi ID"
//	@Param        body  body      dto.UpdatePertanyaanIdentifikasiRequest  true  "Data update"
//	@Success      200   {object}  dto.PertanyaanIdentifikasiResponse
//	@Failure      400   {object}  dto.ErrorResponse
//	@Failure      404   {object}  dto.ErrorResponse
//	@Router       /api/maturity/pertanyaan-identifikasi/{id} [put]
func (h *PertanyaanIdentifikasiHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id int) {
	var req dto.UpdatePertanyaanIdentifikasiRequest
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
			"pertanyaan_identifikasi tidak boleh kosong",
			"pertanyaan_identifikasi minimal 3 karakter":
			utils.RespondError(w, 400, err.Error())
		default:
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 200, resp)
}

// DeletePertanyaanIdentifikasi godoc
//
//	@Summary      Hapus pertanyaan identifikasi
//	@Description  Menghapus data pertanyaan identifikasi berdasarkan ID
//	@Tags         PertanyaanIdentifikasi
//	@Produce      json
//	@Param        id   path      int  true  "PertanyaanIdentifikasi ID"
//	@Success      200  {object}  dto.MessageResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Failure      500  {object}  dto.ErrorResponse
//	@Router       /api/maturity/pertanyaan-identifikasi/{id} [delete]
func (h *PertanyaanIdentifikasiHandler) handleDelete(w http.ResponseWriter, r *http.Request, id int) {
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
