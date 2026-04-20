package handlers

import (
	"encoding/json"
	"ikas/internal/dto"
	"ikas/internal/middleware"
	"ikas/internal/services"
	"ikas/internal/utils"
	"net/http"
	"strconv"
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
			ikasID := r.URL.Query().Get("ikas_id")
			pertanyaanID := r.URL.Query().Get("pertanyaan_proteksi_id")

			if ikasID != "" {
				h.handleGetByIkasID(w, r, ikasID)
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
func (h *JawabanProteksiHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	ikasID := r.URL.Query().Get("ikas_id")
	pertanyaanIDStr := r.URL.Query().Get("pertanyaan_proteksi_id")

	userRole, _ := r.Context().Value(middleware.Role).(string)
	userPerusahaanID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)

	if userRole != "admin" && (userPerusahaanID == "" || userPerusahaanID == "null") {
		utils.RespondJSON(w, 200, map[string]interface{}{
			"message": "Berhasil mengambil data",
			"data":    []dto.JawabanProteksiResponse{},
			"total":   0,
		})
		return
	}

	var data []dto.JawabanProteksiResponse
	var err error

	if ikasID != "" {
		h.handleGetByIkasID(w, r, ikasID)
		return
	} else if pertanyaanIDStr != "" {
		pID, _ := strconv.Atoi(pertanyaanIDStr)
		data, err = h.service.GetByPertanyaan(pID)
	} else {
		if userRole != "admin" {
			data, err = h.service.GetByPerusahaanID(userPerusahaanID, userRole, userPerusahaanID)
		} else {
			data, err = h.service.GetAll(userRole)
		}
	}

	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "format ikas_id tidak valid" {
			status = http.StatusBadRequest
		}
		utils.RespondError(w, status, err.Error())
		return
	}
	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil mengambil data",
		"data":    data,
		"total":   len(data),
	})
}

func (h *JawabanProteksiHandler) handleGetByIkasID(w http.ResponseWriter, r *http.Request, ikasID string) {
	userRole, _ := r.Context().Value(middleware.Role).(string)
	userPerusahaanID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)
	data, err := h.service.GetByIkasID(ikasID, userRole, userPerusahaanID)
	if err != nil {
		rollbar.Error(err)
		if err.Error() == "format ikas_id tidak valid" {
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
	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil mengambil data",
		"data":    data,
		"total":   len(data),
	})
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
func (h *JawabanProteksiHandler) handleGetByID(w http.ResponseWriter, r *http.Request, id int) {
	userRole, _ := r.Context().Value(middleware.Role).(string)
	userPerusahaanID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)

	data, err := h.service.GetByID(id, userRole, userPerusahaanID)
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

	userRole, _ := r.Context().Value(middleware.Role).(string)
	userPerusahaanID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)

	msg, err := h.service.Create(req, userRole, userPerusahaanID)
	if err != nil {
		rollbar.Error(err)
		switch err.Error() {
		case "pertanyaan_proteksi_id tidak boleh kosong",
			"format pertanyaan_proteksi_id tidak valid",
			"ikas_id tidak boleh kosong",
			"format ikas_id tidak valid",
			"jawaban_proteksi tidak boleh kosong",
			"validasi hanya boleh diisi jika evidence ada",
			"validasi hanya boleh berisi 'yes' atau 'no'":
			utils.RespondError(w, 400, err.Error())
		case "pertanyaan_proteksi_id tidak ditemukan",
			"ikas_id tidak ditemukan":
			utils.RespondError(w, 404, err.Error())
		case "pertanyaan ini sudah pernah diisi untuk asesmen ini":
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

	userID := ""
	if val := r.Context().Value(middleware.UserIDKey); val != nil {
		userID = val.(string)
	}

	userRole := ""
	if val := r.Context().Value(middleware.Role); val != nil {
		userRole = val.(string)
	}

	userPerusahaanID := ""
	if val := r.Context().Value(middleware.PerusahaanIDKey); val != nil {
		userPerusahaanID = val.(string)
	}

	err := h.service.Update(id, req, userID, userRole, userPerusahaanID)
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

	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil menyimpan data",
		"id":      id,
	})
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
	userID := ""
	if val := r.Context().Value(middleware.UserIDKey); val != nil {
		userID = val.(string)
	}

	userRole, _ := r.Context().Value(middleware.Role).(string)
	userPerusahaanID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)

	if err := h.service.Delete(id, userID, userRole, userPerusahaanID); err != nil {
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
