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

type JawabanProteksiHandler struct {
	service *services.JawabanProteksiService
}

func NewJawabanProteksiHandler(service *services.JawabanProteksiService) *JawabanProteksiHandler {
	return &JawabanProteksiHandler{
		service: service,
	}
}

func (h *JawabanProteksiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, _ := utils.ExtractIntID(r.URL.Path, "jawaban-proteksi")

	switch r.Method {
	case http.MethodGet:
		if id == 0 {
			// Cek query params untuk filter
			perusahaanID := r.URL.Query().Get("perusahaan_id")
			pertanyaanID := r.URL.Query().Get("pertanyaan_proteksi_id")

			if perusahaanID != "" {
				h.handleGetByPerusahaan(w, r, perusahaanID)
			} else if pertanyaanID != "" {
				idInt, err := utils.StringToInt(pertanyaanID)
				if err != nil {
					utils.RespondError(w, 400, "format pertanyaan_proteksi_id tidak valid")
					return
				}
				h.handleGetByPertanyaan(w, r, idInt)
			} else {
				h.handleGetAll(w, r)
			}
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

// GetAllJawabanProteksi godoc
//
//	@Summary      List semua jawaban proteksi
//	@Description  Mengambil seluruh data jawaban proteksi. Bisa difilter dengan query param perusahaan_id atau pertanyaan_proteksi_id
//	@Tags         JawabanProteksi
//	@Produce      json
//	@Param        perusahaan_id           query     string  false  "Filter by perusahaan ID"
//	@Param        pertanyaan_proteksi_id  query     string  false  "Filter by pertanyaan proteksi ID"
//	@Success      200  {array}   dto.JawabanProteksiResponse
//	@Failure      500  {object}  dto.ErrorResponse
//	@Router       /api/maturity/jawaban-proteksi [get]
func (h *JawabanProteksiHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

func (h *JawabanProteksiHandler) handleGetByPerusahaan(w http.ResponseWriter, _ *http.Request, perusahaanID string) {
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

func (h *JawabanProteksiHandler) handleGetByPertanyaan(w http.ResponseWriter, _ *http.Request, pertanyaanID int) {
	data, err := h.service.GetByPertanyaan(pertanyaanID)
	if err != nil {
		rollbar.Error(err)
		if err.Error() == "format pertanyaan_proteksi_id tidak valid" {
			utils.RespondError(w, 400, err.Error())
		} else {
			utils.RespondError(w, 500, err.Error())
		}
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetJawabanProteksiByID godoc
//
//	@Summary      Ambil jawaban proteksi berdasarkan ID
//	@Description  Mengambil satu data jawaban proteksi
//	@Tags         JawabanProteksi
//	@Produce      json
//	@Param        id   path      int  true  "JawabanProteksi ID"
//	@Success      200  {object}  dto.JawabanProteksiResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Router       /api/maturity/jawaban-proteksi/{id} [get]
func (h *JawabanProteksiHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id int) {
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

// CreateJawabanProteksi godoc
//
//	@Summary      Tambah jawaban proteksi baru
//	@Description  Membuat record jawaban proteksi baru
//	@Tags         JawabanProteksi
//	@Accept       json
//	@Produce      json
//	@Param        body  body      dto.CreateJawabanProteksiRequest  true  "Data jawaban proteksi"
//	@Success      201   {object}  dto.JawabanProteksiResponse
//	@Failure      400   {object}  dto.ErrorResponse
//	@Failure      404   {object}  dto.ErrorResponse
//	@Failure      409   {object}  dto.ErrorResponse
//	@Router       /api/maturity/jawaban-proteksi [post]
func (h *JawabanProteksiHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateJawabanProteksiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rollbar.Error(err)
		errMsg := err.Error()
		if strings.Contains(errMsg, "pertanyaan_proteksi_id") {
			utils.RespondError(w, 400, "pertanyaan_proteksi_id harus berupa integer")
		} else if strings.Contains(errMsg, "jawaban_proteksi") {
			utils.RespondError(w, 400, "jawaban_proteksi harus berupa angka 0-5 atau null untuk N/A")
		} else {
			utils.RespondError(w, 400, "Invalid request body")
		}
		return
	}

	resp, err := h.service.Create(req)
	if err != nil {
		rollbar.Error(err)
		switch err.Error() {
		case "pertanyaan_proteksi_id tidak boleh kosong",
			"format pertanyaan_proteksi_id tidak valid",
			"perusahaan_id tidak boleh kosong",
			"format perusahaan_id tidak valid",
			"jawaban_proteksi tidak boleh kosong",
			"validasi hanya boleh diisi jika evidence ada",
			"validasi hanya boleh berisi 'yes' atau 'no'":
			utils.RespondError(w, 400, err.Error())
		case "pertanyaan_proteksi_id tidak ditemukan",
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

// UpdateJawabanProteksi godoc
//
//	@Summary      Update jawaban proteksi
//	@Description  Mengubah data jawaban proteksi berdasarkan ID
//	@Tags         JawabanProteksi
//	@Accept       json
//	@Produce      json
//	@Param        id    path      int                              true  "JawabanProteksi ID"
//	@Param        body  body      dto.UpdateJawabanProteksiRequest true  "Data update"
//	@Success      200   {object}  dto.JawabanProteksiResponse
//	@Failure      400   {object}  dto.ErrorResponse
//	@Failure      404   {object}  dto.ErrorResponse
//	@Router       /api/maturity/jawaban-proteksi/{id} [put]
func (h *JawabanProteksiHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id int) {
	var req dto.UpdateJawabanProteksiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rollbar.Error(err)
		if strings.Contains(err.Error(), "cannot unmarshal") {
			utils.RespondError(w, 400, "jawaban_proteksi harus berupa angka bulat 0-5 or null for N/A")
			return
		}
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
			"jawaban_proteksi tidak boleh kosong",
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

// DeleteJawabanProteksi godoc
//
//	@Summary      Hapus jawaban proteksi
//	@Description  Menghapus data jawaban proteksi berdasarkan ID
//	@Tags         JawabanProteksi
//	@Produce      json
//	@Param        id   path      int  true  "JawabanProteksi ID"
//	@Success      200  {object}  dto.MessageResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Failure      500  {object}  dto.ErrorResponse
//	@Router       /api/maturity/jawaban-proteksi/{id} [delete]
func (h *JawabanProteksiHandler) handleDelete(w http.ResponseWriter, r *http.Request, id int) {
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
