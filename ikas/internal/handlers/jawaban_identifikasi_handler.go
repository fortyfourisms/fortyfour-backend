package handlers

import (
	"encoding/json"
	"ikas/internal/dto"
	"ikas/internal/middleware"
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
	id, _ := utils.ExtractIntID(r.URL.Path, "jawaban-identifikasi")

	switch r.Method {
	case http.MethodGet:
		if id == 0 {
			// Cek query params untuk filter
			perusahaanID := r.URL.Query().Get("perusahaan_id")
			pertanyaanID := r.URL.Query().Get("pertanyaan_identifikasi_id")

			if perusahaanID != "" {
				h.handleGetByPerusahaan(w, r, perusahaanID)
			} else if pertanyaanID != "" {
				idInt, err := utils.StringToInt(pertanyaanID)
				if err != nil {
					utils.RespondError(w, 400, "format pertanyaan_identifikasi_id tidak valid")
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
//	@Router       /api/maturity/jawaban-identifikasi [get]
func (h *JawabanIdentifikasiHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
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
	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil mengambil data",
		"data":    data,
		"total":   len(data),
	})
}

func (h *JawabanIdentifikasiHandler) handleGetByPertanyaan(w http.ResponseWriter, _ *http.Request, pertanyaanID int) {
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
	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil mengambil data",
		"data":    data,
		"total":   len(data),
	})
}

// GetJawabanIdentifikasiByID godoc
//
//	@Summary      Ambil jawaban identifikasi berdasarkan ID
//	@Description  Mengambil satu data jawaban identifikasi
//	@Tags         JawabanIdentifikasi
//	@Produce      json
//	@Param        id   path      int  true  "JawabanIdentifikasi ID"
//	@Success      200  {object}  dto.JawabanIdentifikasiResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Router       /api/maturity/jawaban-identifikasi/{id} [get]
func (h *JawabanIdentifikasiHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id int) {
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
//	@Router       /api/maturity/jawaban-identifikasi [post]
func (h *JawabanIdentifikasiHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateJawabanIdentifikasiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rollbar.Error(err)
		errMsg := err.Error()
		if strings.Contains(errMsg, "pertanyaan_identifikasi_id") {
			utils.RespondError(w, 400, "pertanyaan_identifikasi_id harus berupa integer")
		} else if strings.Contains(errMsg, "jawaban_identifikasi") {
			utils.RespondError(w, 400, "jawaban_identifikasi harus berupa angka 0-5 atau null untuk N/A")
		} else {
			utils.RespondError(w, 400, "Invalid request body")
		}
		return
	}

	msg, err := h.service.Create(req)
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
		case "pertanyaan ini sudah pernah diisi oleh perusahaan Anda":
			utils.RespondError(w, 409, err.Error())
		default:
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 201, map[string]interface{}{
		"message": msg,
	})
}

// UpdateJawabanIdentifikasi godoc
//
//	@Summary      Update jawaban identifikasi
//	@Description  Mengubah data jawaban identifikasi berdasarkan ID
//	@Tags         JawabanIdentifikasi
//	@Accept       json
//	@Produce      json
//	@Param        id    path      int                                true  "JawabanIdentifikasi ID"
//	@Param        body  body      dto.UpdateJawabanIdentifikasiRequest  true  "Data update"
//	@Success      200   {object}  dto.JawabanIdentifikasiResponse
//	@Failure      400   {object}  dto.ErrorResponse
//	@Failure      404   {object}  dto.ErrorResponse
//	@Router       /api/maturity/jawaban-identifikasi/{id} [put]
func (h *JawabanIdentifikasiHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id int) {
	var req dto.UpdateJawabanIdentifikasiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rollbar.Error(err)
		if strings.Contains(err.Error(), "cannot unmarshal") {
			utils.RespondError(w, 400, "jawaban_identifikasi harus berupa angka bulat 0-5 or null for N/A")
			return
		}
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	userID := ""
	if val := r.Context().Value(middleware.UserIDKey); val != nil {
		userID = val.(string)
	}

	err := h.service.Update(id, req, userID)
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

	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil menyimpan data",
		"id":      id,
	})
}

// DeleteJawabanIdentifikasi godoc
//
//	@Summary      Hapus jawaban identifikasi
//	@Description  Menghapus data jawaban identifikasi berdasarkan ID
//	@Tags         JawabanIdentifikasi
//	@Produce      json
//	@Param        id   path      int  true  "JawabanIdentifikasi ID"
//	@Success      200  {object}  dto.MessageResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Failure      500  {object}  dto.ErrorResponse
//	@Router       /api/maturity/jawaban-identifikasi/{id} [delete]
func (h *JawabanIdentifikasiHandler) handleDelete(w http.ResponseWriter, r *http.Request, id int) {
	userID := ""
	if val := r.Context().Value(middleware.UserIDKey); val != nil {
		userID = val.(string)
	}

	if err := h.service.Delete(id, userID); err != nil {
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
		"id":      id,
	})
}
