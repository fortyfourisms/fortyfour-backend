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

type RuangLingkupHandler struct {
	service *services.RuangLingkupService
}

func NewRuangLingkupHandler(service *services.RuangLingkupService) *RuangLingkupHandler {
	return &RuangLingkupHandler{
		service: service,
	}
}

func (h *RuangLingkupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/ruang-lingkup")
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
			utils.RespondError(w, 400, "ID wajib untuk update")
			return
		}
		h.handleUpdate(w, r, id)
	case http.MethodDelete:
		if id == "" {
			utils.RespondError(w, 400, "ID wajib untuk hapus")
			return
		}
		h.handleDelete(w, r, id)
	default:
		utils.RespondError(w, 405, "Method tidak diizinkan")
	}
}

// @Summary      List semua ruang lingkup
// @Description  Mengambil seluruh data ruang lingkup
// @Tags         RuangLingkup
// @Produce      json
// @Success      200  {array}   dto.RuangLingkupResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/ruang-lingkup [get]
func (h *RuangLingkupHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 500, "Gagal mengambil data ruang lingkup")
		return
	}

	// Kalau kosong return empty array, bukan null
	if len(data) == 0 {
		utils.RespondJSON(w, 200, []dto.RuangLingkupResponse{})
		return
	}

	utils.RespondJSON(w, 200, data)
}

// @Summary      Ambil ruang lingkup berdasarkan ID
// @Description  Mengambil satu data ruang lingkup
// @Tags         RuangLingkup
// @Produce      json
// @Param        id   path      string  true  "RuangLingkup ID"
// @Success      200  {object}  dto.RuangLingkupResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Router       /api/ruang-lingkup/{id} [get]
func (h *RuangLingkupHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
	if !isValidUUID(id) {
		utils.RespondError(w, 400, "Format ID tidak valid")
		return
	}

	data, err := h.service.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 404, "Data ruang lingkup tidak ditemukan")
		return
	}

	utils.RespondJSON(w, 200, data)
}

// @Summary      Tambah ruang lingkup baru
// @Description  Membuat record ruang lingkup baru
// @Tags         RuangLingkup
// @Accept       json
// @Produce      json
// @Param        ruangLingkup  body      dto.CreateRuangLingkupRequest  true  "Data ruang lingkup"
// @Success      201           {object}  dto.RuangLingkupResponse
// @Failure      400           {object}  dto.ErrorResponse
// @Router       /api/ruang-lingkup [post]
func (h *RuangLingkupHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRuangLingkupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, "Format request body tidak valid, pastikan format JSON yang benar")
		return
	}

	// Validasi nama
	if errMsg := validateNamaRuangLingkup(req.NamaRuangLingkup); errMsg != "" {
		utils.RespondError(w, 400, errMsg)
		return
	}

	resp, err := h.service.Create(req)
	if err != nil {
		rollbar.Error(err)
		if err.Error() == "nama ruang lingkup sudah ada" {
			utils.RespondError(w, 409, "nama ruang lingkup sudah ada, gunakan nama yang berbeda")
			return
		}
		utils.RespondError(w, 500, "Gagal menambahkan data ruang lingkup")
		return
	}

	utils.RespondJSON(w, 201, resp)
}

// @Summary      Update ruang lingkup
// @Description  Mengubah data ruang lingkup berdasarkan ID
// @Tags         RuangLingkup
// @Accept       json
// @Produce      json
// @Param        id              path      string                       true  "RuangLingkup ID"
// @Param        ruangLingkup    body      dto.UpdateRuangLingkupRequest true  "Data update"
// @Success      200             {object}  dto.RuangLingkupResponse
// @Failure      400             {object}  dto.ErrorResponse
// @Router       /api/ruang-lingkup/{id} [put]
func (h *RuangLingkupHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	if !isValidUUID(id) {
		utils.RespondError(w, 400, "Format ID tidak valid")
		return
	}

	var req dto.UpdateRuangLingkupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, "Format request body tidak valid, pastikan format JSON yang benar")
		return
	}

	// Cek body kosong atau tidak ada field yang dikirim
	if req.NamaRuangLingkup == nil {
		utils.RespondError(w, 400, "Tidak ada field yang akan diupdate")
		return
	}

	// Validasi nama kalau dikirim
	if req.NamaRuangLingkup != nil {
		if errMsg := validateNamaRuangLingkup(*req.NamaRuangLingkup); errMsg != "" {
			utils.RespondError(w, 400, errMsg)
			return
		}
	}

	resp, err := h.service.Update(id, req)
	if err != nil {
		rollbar.Error(err)
		if err.Error() == "nama ruang lingkup sudah ada" {
			utils.RespondError(w, 409, "nama ruang lingkup sudah ada, gunakan nama yang berbeda")
			return
		}
		utils.RespondError(w, 404, "Data ruang lingkup tidak ditemukan atau gagal diupdate")
		return
	}

	utils.RespondJSON(w, 200, resp)
}

// @Summary      Hapus ruang lingkup
// @Description  Menghapus data ruang lingkup berdasarkan ID
// @Tags         RuangLingkup
// @Produce      json
// @Param        id   path      string  true  "RuangLingkup ID"
// @Success      200  {object}  dto.MessageResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /api/ruang-lingkup/{id} [delete]
func (h *RuangLingkupHandler) handleDelete(w http.ResponseWriter, _ *http.Request, id string) {
	if !isValidUUID(id) {
		utils.RespondError(w, 400, "Format ID tidak valid")
		return
	}

	// Cek dulu data ada apa tidak sebelum hapus
	existing, err := h.service.GetByID(id)
	if err != nil || existing == nil {
		rollbar.Error(err)
		utils.RespondError(w, 404, "Data ruang lingkup tidak ditemukan, tidak bisa dihapus")
		return
	}

	// Hapus
	if err := h.service.Delete(id); err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 500, "Gagal menghapus data ruang lingkup")
		return
	}

	utils.RespondJSON(w, 200, map[string]string{"message": "Data ruang lingkup berhasil dihapus"})
}

// HELPER / VALIDASI
func validateNamaRuangLingkup(nama string) string {
	nama = strings.TrimSpace(nama)

	if nama == "" {
		return "nama ruang lingkup tidak boleh kosong"
	}
	if len(nama) < 2 {
		return "nama ruang lingkup minimal 2 karakter"
	}
	if len(nama) > 50 {
		return "nama ruang lingkup maksimal 50 karakter"
	}

	return ""
}

// isValidUUID — cek format UUID 8-4-4-4-12
func isValidUUID(id string) bool {
	if len(id) != 36 {
		return false
	}
	for i, c := range id {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if c != '-' {
				return false
			}
		} else {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
				return false
			}
		}
	}
	return true
}
