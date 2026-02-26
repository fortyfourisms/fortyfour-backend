package handlers

import (
	"encoding/json"
	"ikas/internal/dto"
	"ikas/internal/services"
	"ikas/internal/utils"
	"net/http"
	"strings"

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
	path := strings.TrimPrefix(r.URL.Path, "/api/maturity/ruang-lingkup")
	id := strings.TrimPrefix(path, "/")

	switch r.Method {
	case http.MethodGet:
		if id == "" {
			h.handleGetAll(w, r)
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
	utils.RespondJSON(w, 200, data)
}

// GetRuangLingkupByID godoc
//
//	@Summary      Ambil ruang lingkup berdasarkan ID
//	@Description  Mengambil satu data ruang lingkup
//	@Tags         RuangLingkup
//	@Produce      json
//	@Param        id   path      string  true  "RuangLingkup ID"
//	@Success      200  {object}  dto.RuangLingkupResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Router       /api/maturity/ruang-lingkup/{id} [get]
func (h *RuangLingkupHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
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
	utils.RespondJSON(w, 200, data)
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

	utils.RespondJSON(w, 201, resp)
}

// UpdateRuangLingkup godoc
//
//	@Summary      Update ruang lingkup
//	@Description  Mengubah data ruang lingkup berdasarkan ID
//	@Tags         RuangLingkup
//	@Accept       json
//	@Produce      json
//	@Param        id              path      string                       true  "RuangLingkup ID"
//	@Param        ruangLingkup    body      dto.UpdateRuangLingkupRequest true  "Data update"
//	@Success      200             {object}  dto.RuangLingkupResponse
//	@Failure      400             {object}  dto.ErrorResponse
//	@Failure      404             {object}  dto.ErrorResponse
//	@Failure      409             {object}  dto.ErrorResponse
//	@Router       /api/maturity/ruang-lingkup/{id} [put]
func (h *RuangLingkupHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
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

	utils.RespondJSON(w, 200, resp)
}

// DeleteRuangLingkup godoc
//
//	@Summary      Hapus ruang lingkup
//	@Description  Menghapus data ruang lingkup berdasarkan ID
//	@Tags         RuangLingkup
//	@Produce      json
//	@Param        id   path      string  true  "RuangLingkup ID"
//	@Success      200  {object}  dto.MessageResponse
//	@Failure      404  {object}  dto.ErrorResponse
//	@Failure      500  {object}  dto.ErrorResponse
//	@Router       /api/maturity/ruang-lingkup/{id} [delete]
func (h *RuangLingkupHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.Delete(id); err != nil {
		logger.Error(err, "operation failed")
		if err.Error() == "data tidak ditemukan" {
			utils.RespondError(w, 404, err.Error())
		} else {
			utils.RespondError(w, 500, err.Error())
		}
		return
	}
	utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}
