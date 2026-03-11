package handlers

import (
	"encoding/json"
	"ikas/internal/dto"
	"ikas/internal/services"
	"ikas/internal/utils"
	"net/http"

	"fortyfour-backend/pkg/logger"
)

type RuangLingkupHandler struct {
	service *services.RuangLingkupService
}

func NewRuangLingkupHandler(service *services.RuangLingkupService) *RuangLingkupHandler {
	return &RuangLingkupHandler{
		service: service,
	}
}

func (h *RuangLingkupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, _ := utils.ExtractIntID(r.URL.Path, "ruang-lingkup")

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

// GetAllRuangLingkup godoc
//
//	@Summary      List semua ruang lingkup
//	@Description  Mengambil seluruh data ruang lingkup
//	@Tags         RuangLingkup
//	@Produce      json
//	@Success      200  {array}   dto.RuangLingkupResponse
//	@Failure      500  {object}  dto.ErrorResponse
//	@Router       /api/maturity/ruang-lingkup [get]
func (h *RuangLingkupHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
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

// GetRuangLingkupByID godoc
//
//	@Summary      Ambil ruang lingkup berdasarkan ID
//	@Description  Mengambil satu data ruang lingkup
//	@Tags         RuangLingkup
//	@Produce      json
//	@Param        id   path      int  true  "RuangLingkup ID"
//	@Success      200  {object}  dto.RuangLingkupResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Router       /api/maturity/ruang-lingkup/{id} [get]
func (h *RuangLingkupHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id int) {
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
	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil mengambil data",
		"data":    data,
	})
}

// CreateRuangLingkup godoc
//
//	@Summary      Tambah ruang lingkup baru
//	@Description  Membuat record ruang lingkup baru
//	@Tags         RuangLingkup
//	@Accept       json
//	@Produce      json
//	@Param        ruangLingkup  body      dto.CreateRuangLingkupRequest  true  "Data ruang lingkup"
//	@Success      201           {object}  dto.RuangLingkupResponse
//	@Failure      400           {object}  dto.ErrorResponse
//	@Failure      409           {object}  dto.ErrorResponse
//	@Router       /api/maturity/ruang-lingkup [post]
func (h *RuangLingkupHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRuangLingkupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Create(req)
	if err != nil {
		logger.Error(err, "operation failed")
		// Mapping error ke HTTP status code
		switch err.Error() {
		case "nama_ruang_lingkup tidak boleh kosong",
			"nama_ruang_lingkup minimal 3 karakter",
			"nama_ruang_lingkup maksimal 50 karakter":
			utils.RespondError(w, 400, err.Error())
		case "nama_ruang_lingkup sudah ada":
			utils.RespondError(w, 409, err.Error())
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

// UpdateRuangLingkup godoc
//
//	@Summary      Update ruang lingkup
//	@Description  Mengubah data ruang lingkup berdasarkan ID
//	@Tags         RuangLingkup
//	@Accept       json
//	@Produce      json
//	@Param        id              path      int                       true  "RuangLingkup ID"
//	@Param        ruangLingkup    body      dto.UpdateRuangLingkupRequest true  "Data update"
//	@Success      200             {object}  dto.RuangLingkupResponse
//	@Failure      400             {object}  dto.ErrorResponse
//	@Failure      404             {object}  dto.ErrorResponse
//	@Failure      409             {object}  dto.ErrorResponse
//	@Router       /api/maturity/ruang-lingkup/{id} [put]
func (h *RuangLingkupHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id int) {
	var req dto.UpdateRuangLingkupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Update(id, req)
	if err != nil {
		logger.Error(err, "operation failed")
		// Mapping error ke HTTP status code
		switch err.Error() {
		case "data tidak ditemukan":
			utils.RespondError(w, 404, err.Error())
		case "nama_ruang_lingkup tidak boleh kosong",
			"nama_ruang_lingkup minimal 3 karakter",
			"nama_ruang_lingkup maksimal 50 karakter":
			utils.RespondError(w, 400, err.Error())
		case "nama_ruang_lingkup sudah ada":
			utils.RespondError(w, 409, err.Error())
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

// DeleteRuangLingkup godoc
//
//	@Summary      Hapus ruang lingkup
//	@Description  Menghapus data ruang lingkup berdasarkan ID
//	@Tags         RuangLingkup
//	@Produce      json
//	@Param        id   path      int  true  "RuangLingkup ID"
//	@Success      200  {object}  dto.MessageResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Failure      500  {object}  dto.ErrorResponse
//	@Router       /api/maturity/ruang-lingkup/{id} [delete]
func (h *RuangLingkupHandler) handleDelete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.service.Delete(id); err != nil {
		logger.Error(err, "operation failed")
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
