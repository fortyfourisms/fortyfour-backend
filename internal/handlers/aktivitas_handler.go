package handlers

import (
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"net/http"

	"fortyfour-backend/pkg/logger"
)

type AktivitasHandler struct {
	service *services.AktivitasService
}

func NewAktivitasHandler(service *services.AktivitasService) *AktivitasHandler {
	return &AktivitasHandler{
		service: service,
	}
}

func (h *AktivitasHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/aktivitas/jenis" {
		h.HandleGetJenis(w, r)
		return
	}

	id, _ := utils.ExtractIntID(r.URL.Path, "aktivitas")

	switch r.Method {
	case http.MethodGet:
		if id == 0 {
			perusahaanID := r.URL.Query().Get("perusahaan_id")
			if perusahaanID != "" {
				h.handleGetByPerusahaanID(w, r, perusahaanID)
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

// GetAllAktivitas godoc
//
//	@Summary		List semua aktivitas
//	@Description	Mengambil seluruh data aktivitas
//	@Tags			Aktivitas
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	utils.JSONResponse{data=[]dto.AktivitasResponse}
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/api/aktivitas [get]
func (h *AktivitasHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, utils.JSONResponse{
		Status:  "success",
		Message: "Berhasil mengambil data aktivitas",
		Data:    data,
		Total:   len(data),
	})
}

// GetAktivitasByID godoc
//
//	@Summary		Ambil aktivitas berdasarkan ID
//	@Description	Mengambil satu data aktivitas
//	@Tags			Aktivitas
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Aktivitas ID"
//	@Success		200	{object}	utils.JSONResponse{data=dto.AktivitasResponse}
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/api/aktivitas/{id} [get]
func (h *AktivitasHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id int) {
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
	utils.RespondJSON(w, 200, utils.JSONResponse{
		Status:  "success",
		Message: "Berhasil mengambil data aktivitas",
		Data:    data,
	})
}

// GetAktivitasByPerusahaanID godoc
//
//	@Summary		List aktivitas berdasarkan Perusahaan ID
//	@Description	Mengambil data aktivitas per perusahaan
//	@Tags			Aktivitas
//	@Produce		json
//	@Security		BearerAuth
//	@Param			perusahaan_id	query		string	true	"Perusahaan ID"
//	@Success		200				{object}	utils.JSONResponse{data=[]dto.AktivitasResponse}
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/api/aktivitas [get]
func (h *AktivitasHandler) handleGetByPerusahaanID(w http.ResponseWriter, _ *http.Request, perusahaanID string) {
	data, err := h.service.GetByPerusahaanID(perusahaanID)
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, utils.JSONResponse{
		Status:  "success",
		Message: "Berhasil mengambil data aktivitas perusahaan",
		Data:    data,
		Total:   len(data),
	})
}

// CreateAktivitas godoc
//
//	@Summary		Tambah aktivitas baru
//	@Description	Membuat record aktivitas baru
//	@Tags			Aktivitas
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			aktivitas	body		dto.CreateAktivitasRequest	true	"Data aktivitas"
//	@Success		201			{object}	utils.JSONResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Router			/api/aktivitas [post]
func (h *AktivitasHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAktivitasRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	_, err := h.service.Create(req)
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 201, utils.JSONResponse{
		Status:  "success",
		Message: "Permintaan pembuatan aktivitas telah diterima dan sedang diproses",
	})
}

// UpdateAktivitas godoc
//
//	@Summary		Update aktivitas
//	@Description	Mengubah data aktivitas berdasarkan ID
//	@Tags			Aktivitas
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			path		int							true	"Aktivitas ID"
//	@Param			aktivitas	body		dto.UpdateAktivitasRequest	true	"Data update"
//	@Success		200			{object}	utils.JSONResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Router			/api/aktivitas/{id} [put]
func (h *AktivitasHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id int) {
	var req dto.UpdateAktivitasRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	_, err := h.service.Update(id, req)
	if err != nil {
		logger.Error(err, "operation failed")
		if err.Error() == "data tidak ditemukan" {
			utils.RespondError(w, 404, err.Error())
		} else {
			utils.RespondError(w, 400, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 200, utils.JSONResponse{
		Status:  "success",
		Message: "Permintaan pembaruan aktivitas telah diterima dan sedang diproses",
	})
}

// DeleteAktivitas godoc
//
//	@Summary		Hapus aktivitas
//	@Description	Menghapus data aktivitas berdasarkan ID
//	@Tags			Aktivitas
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Aktivitas ID"
//	@Success		200	{object}	utils.JSONResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/api/aktivitas/{id} [delete]
func (h *AktivitasHandler) handleDelete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.service.Delete(id); err != nil {
		logger.Error(err, "operation failed")
		if err.Error() == "data tidak ditemukan" {
			utils.RespondError(w, 404, err.Error())
		} else {
			utils.RespondError(w, 500, err.Error())
		}
		return
	}
	utils.RespondJSON(w, 200, utils.JSONResponse{
		Status:  "success",
		Message: "Permintaan penghapusan aktivitas telah diterima dan sedang diproses",
	})
}

// GetJenisAktivitas godoc
//
//	@Summary		List jenis aktivitas
//	@Description	Mengambil daftar jenis aktivitas yang diperbolehkan (untuk dropdown)
//	@Tags			Aktivitas
//	@Produce		json
//	@Success		200	{object}	utils.JSONResponse{data=[]string}
//	@Router			/api/aktivitas/jenis [get]
func (h *AktivitasHandler) HandleGetJenis(w http.ResponseWriter, _ *http.Request) {
	data := h.service.GetAllowedJenis()
	utils.RespondJSON(w, 200, utils.JSONResponse{
		Status:  "success",
		Message: "Berhasil mengambil daftar jenis aktivitas",
		Data:    data,
	})
}
