package handlers

import (
	"encoding/json"
	"ikas/internal/dto"
	"ikas/internal/services"
	"ikas/internal/utils"
	"net/http"

	"fortyfour-backend/pkg/logger"
)

type DomainHandler struct {
	service *services.DomainService
}

func NewDomainHandler(service *services.DomainService) *DomainHandler {
	return &DomainHandler{
		service: service,
	}
}

func (h *DomainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, _ := utils.ExtractIntID(r.URL.Path, "domain")

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

// GetAllDomain godoc
//
//	@Summary      List semua domain
//	@Description  Mengambil seluruh data domain
//	@Tags         Domain
//	@Produce      json
//	@Success      200  {array}   dto.DomainResponse
//	@Failure      500  {object}  dto.ErrorResponse
//	@Router       /api/maturity/domain [get]
func (h *DomainHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
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

// GetDomainByID godoc
//
//	@Summary      Ambil domain berdasarkan ID
//	@Description  Mengambil satu data domain
//	@Tags         Domain
//	@Produce      json
//	@Param        id   path      int  true  "Domain ID"
//	@Success      200  {object}  dto.DomainResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Router       /api/maturity/domain/{id} [get]
func (h *DomainHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id int) {
	data, err := h.service.GetByID(id)
	if err != nil {
		logger.Error(err, "operation failed")
		if err.Error() == "data tidak ditemukan" || err.Error() == "format ID tidak valid" {
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

// CreateDomain godoc
//
//	@Summary      Tambah domain baru
//	@Description  Membuat record domain baru
//	@Tags         Domain
//	@Accept       json
//	@Produce      json
//	@Param        domain  body      dto.CreateDomainRequest  true  "Data domain"
//	@Success      201     {object}  dto.DomainResponse
//	@Failure      400     {object}  dto.ErrorResponse
//	@Failure      409     {object}  dto.ErrorResponse
//	@Router       /api/maturity/domain [post]
func (h *DomainHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateDomainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	_, err := h.service.Create(req)
	if err != nil {
		logger.Error(err, "operation failed")
		switch err.Error() {
		case "nama_domain tidak boleh kosong",
			"nama_domain minimal 3 karakter",
			"nama_domain maksimal 50 karakter",
			"nama_domain mengandung karakter tidak diizinkan",
			"nama_domain hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&":
			utils.RespondError(w, 400, err.Error())
		case "nama_domain sudah ada":
			utils.RespondError(w, 409, err.Error())
		default:
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 201, dto.DomainMessageResponse{
		Message: "Berhasil menyimpan data",
	})
}

// UpdateDomain godoc
//
//	@Summary      Update domain
//	@Description  Mengubah data domain berdasarkan ID
//	@Tags         Domain
//	@Accept       json
//	@Produce      json
//	@Param        id      path      string                     true  "Domain ID"
//	@Param        domain  body      dto.UpdateDomainRequest    true  "Data update"
//	@Success      200     {object}  dto.DomainResponse
//	@Failure      400     {object}  dto.ErrorResponse
//	@Failure      404     {object}  dto.ErrorResponse
//	@Failure      409     {object}  dto.ErrorResponse
//	@Router       /api/maturity/domain/{id} [put]
func (h *DomainHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id int) {
	var req dto.UpdateDomainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	_, err := h.service.Update(id, req)
	if err != nil {
		logger.Error(err, "operation failed")
		switch err.Error() {
		case "data tidak ditemukan":
			utils.RespondError(w, 404, err.Error())
		case "nama_domain tidak boleh kosong",
			"nama_domain minimal 3 karakter",
			"nama_domain maksimal 50 karakter",
			"nama_domain mengandung karakter tidak diizinkan",
			"nama_domain hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&":
			utils.RespondError(w, 400, err.Error())
		case "nama_domain sudah ada":
			utils.RespondError(w, 409, err.Error())
		default:
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 200, dto.DomainMessageResponse{
		ID:      id,
		Message: "Berhasil memperbarui data",
	})
}

// DeleteDomain godoc
//
//	@Summary      Hapus domain
//	@Description  Menghapus data domain berdasarkan ID
//	@Tags         Domain
//	@Produce      json
//	@Param        id   path      int  true  "Domain ID"
//	@Success      200  {object}  dto.MessageResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Failure      500  {object}  dto.ErrorResponse
//	@Router       /api/maturity/domain/{id} [delete]
func (h *DomainHandler) handleDelete(w http.ResponseWriter, _ *http.Request, id int) {
	if err := h.service.Delete(id); err != nil {
		logger.Error(err, "operation failed")
		if err.Error() == "data tidak ditemukan" {
			utils.RespondError(w, 404, err.Error())
		} else {
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 200, dto.DomainMessageResponse{
		ID:      id,
		Message: "Berhasil menghapus data",
	})
}
