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

type JawabanIdentifikasiHandler struct {
	service *services.JawabanIdentifikasiService
}

func NewJawabanIdentifikasiHandler(service *services.JawabanIdentifikasiService) *JawabanIdentifikasiHandler {
	return &JawabanIdentifikasiHandler{
		service: service,
	}
}

func (h *JawabanIdentifikasiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/jawaban-identifikasi")
	id := strings.TrimPrefix(path, "/")

	switch r.Method {
	case http.MethodGet:
		if id == "" {
			// Cek query params untuk filter
			perusahaanID := r.URL.Query().Get("perusahaan_id")
			pertanyaanID := r.URL.Query().Get("pertanyaan_identifikasi_id")

			if perusahaanID != "" {
				h.handleGetByPerusahaan(w, r, perusahaanID)
			} else if pertanyaanID != "" {
				h.handleGetByPertanyaan(w, r, pertanyaanID)
			} else {
				h.handleGetAll(w, r)
			}
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

// GetAllJawabanIdentifikasi godoc
//
//	@Summary      List semua jawaban identifikasi
//	@Description  Mengambil seluruh data jawaban identifikasi. Bisa difilter dengan query param perusahaan_id atau pertanyaan_identifikasi_id
//	@Tags         JawabanIdentifikasi
//	@Produce      json
//	@Param        perusahaan_id              query     string  false  "Filter by perusahaan ID"
//	@Param        pertanyaan_identifikasi_id query     string  false  "Filter by pertanyaan identifikasi ID"
//	@Success      200  {array}   dto.JawabanIdentifikasiResponse
//	@Failure      500  {object}  dto.ErrorResponse
//	@Router       /api/jawaban-identifikasi [get]
func (h *JawabanIdentifikasiHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

func (h *JawabanIdentifikasiHandler) handleGetByPerusahaan(w http.ResponseWriter, _ *http.Request, perusahaanID string) {
	data, err := h.service.GetByPerusahaan(perusahaanID)
	if err != nil {
		rollbar.Error(err)
		if err.Error() == "format perusahaan_id tidak valid" {
			utils.RespondError(w, 400, err.Error())
		} else {
			utils.RespondError(w, 500, err.Error())
		}
		return
	}
	utils.RespondJSON(w, 200, data)
}

func (h *JawabanIdentifikasiHandler) handleGetByPertanyaan(w http.ResponseWriter, _ *http.Request, pertanyaanID string) {
	data, err := h.service.GetByPertanyaan(pertanyaanID)
	if err != nil {
		rollbar.Error(err)
		if err.Error() == "format pertanyaan_identifikasi_id tidak valid" {
			utils.RespondError(w, 400, err.Error())
		} else {
			utils.RespondError(w, 500, err.Error())
		}
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetJawabanIdentifikasiByID godoc
//
//	@Summary      Ambil jawaban identifikasi berdasarkan ID
//	@Description  Mengambil satu data jawaban identifikasi
//	@Tags         JawabanIdentifikasi
//	@Produce      json
//	@Param        id   path      string  true  "JawabanIdentifikasi ID"
//	@Success      200  {object}  dto.JawabanIdentifikasiResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Router       /api/jawaban-identifikasi/{id} [get]
func (h *JawabanIdentifikasiHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
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

// CreateJawabanIdentifikasi godoc
//
//	@Summary      Tambah jawaban identifikasi baru
//	@Description  Membuat record jawaban identifikasi baru
//	@Tags         JawabanIdentifikasi
//	@Accept       json
//	@Produce      json
//	@Param        body  body      dto.CreateJawabanIdentifikasiRequest  true  "Data jawaban identifikasi"
//	@Success      201   {object}  dto.JawabanIdentifikasiResponse
//	@Failure      400   {object}  dto.ErrorResponse
//	@Failure      404   {object}  dto.ErrorResponse
//	@Failure      409   {object}  dto.ErrorResponse
//	@Router       /api/jawaban-identifikasi [post]
func (h *JawabanIdentifikasiHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateJawabanIdentifikasiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Create(req)
	if err != nil {
		rollbar.Error(err)
		switch err.Error() {
		case "pertanyaan_identifikasi_id tidak boleh kosong",
			"format pertanyaan_identifikasi_id tidak valid",
			"perusahaan_id tidak boleh kosong",
			"format perusahaan_id tidak valid",
			"jawaban_identifikasi tidak boleh kosong",
			"validasi hanya boleh diisi jika evidence ada",
			"validasi hanya boleh berisi 'yes' atau 'no'":
			utils.RespondError(w, 400, err.Error())
		case "pertanyaan_identifikasi_id tidak ditemukan",
			"perusahaan_id tidak ditemukan":
			utils.RespondError(w, 404, err.Error())
		case "jawaban untuk pertanyaan ini sudah ada untuk perusahaan tersebut":
			utils.RespondError(w, 409, err.Error())
		default:
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 201, resp)
}

// UpdateJawabanIdentifikasi godoc
//
//	@Summary      Update jawaban identifikasi
//	@Description  Mengubah data jawaban identifikasi berdasarkan ID
//	@Tags         JawabanIdentifikasi
//	@Accept       json
//	@Produce      json
//	@Param        id    path      string                                true  "JawabanIdentifikasi ID"
//	@Param        body  body      dto.UpdateJawabanIdentifikasiRequest  true  "Data update"
//	@Success      200   {object}  dto.JawabanIdentifikasiResponse
//	@Failure      400   {object}  dto.ErrorResponse
//	@Failure      404   {object}  dto.ErrorResponse
//	@Router       /api/jawaban-identifikasi/{id} [put]
func (h *JawabanIdentifikasiHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req dto.UpdateJawabanIdentifikasiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Update(id, req)
	if err != nil {
		rollbar.Error(err)
		switch err.Error() {
		case "data tidak ditemukan":
			utils.RespondError(w, 404, err.Error())
		case "format ID tidak valid",
			"jawaban_identifikasi tidak boleh kosong",
			"validasi hanya boleh diisi jika evidence ada",
			"validasi hanya boleh berisi 'yes' atau 'no'":
			utils.RespondError(w, 400, err.Error())
		default:
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 200, resp)
}

// DeleteJawabanIdentifikasi godoc
//
//	@Summary      Hapus jawaban identifikasi
//	@Description  Menghapus data jawaban identifikasi berdasarkan ID
//	@Tags         JawabanIdentifikasi
//	@Produce      json
//	@Param        id   path      string  true  "JawabanIdentifikasi ID"
//	@Success      200  {object}  dto.MessageResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Failure      500  {object}  dto.ErrorResponse
//	@Router       /api/jawaban-identifikasi/{id} [delete]
func (h *JawabanIdentifikasiHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
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
