package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"fortyfour-backend/pkg/logger"
)

type SdmCsirtHandler struct {
	service      services.SdmCsirtServiceInterface
	csirtService services.CsirtServiceInterface
	sseService   *services.SSEService
}

func NewSdmCsirtHandler(service services.SdmCsirtServiceInterface, csirtService services.CsirtServiceInterface, sseService *services.SSEService) *SdmCsirtHandler {
	return &SdmCsirtHandler{service: service, csirtService: csirtService, sseService: sseService}
}

func (h *SdmCsirtHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/sdm_csirt")
	id := strings.Trim(path, "/")

	switch r.Method {
	case http.MethodGet:
		if id == "" {
			h.handleGetAll(w, r)
		} else {
			h.handleGetByID(w, r, id)
		}
	case http.MethodPost:
		h.handleCreate(w, r)
	case http.MethodPut:
		h.handleUpdate(w, r, id)
	case http.MethodDelete:
		h.handleDelete(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// GetAllSDM godoc
// @Summary      List semua sdm csirt
// @Description  Mengambil seluruh data sdm csirt
// @Tags         SDM
// @Produce      json
// @Success      200  {array}  dto.SdmCsirtResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /api/sdm_csirt [get]
func (h *SdmCsirtHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())

	if role != "user" {
		// admin atau no-context: return semua
		data, err := h.service.GetAll()
		if err != nil {
			logger.Error(err, "failed to get all SDM CSIRT data")
			utils.RespondError(w, 400, err.Error())
			return
		}
		utils.RespondJSON(w, 200, data)
		return
	}

	// user: ambil SDM berdasarkan CSIRT milik perusahaannya
	idPerusahaan := middleware.GetIDPerusahaan(r.Context())
	if idPerusahaan == "" {
		utils.RespondError(w, 403, "Akun Anda belum terhubung ke perusahaan")
		return
	}

	// Cari CSIRT milik perusahaan ini
	csirtList, err := h.csirtService.GetByPerusahaan(idPerusahaan)
	if err != nil || len(csirtList) == 0 {
		utils.RespondJSON(w, 200, []dto.SdmCsirtResponse{})
		return
	}

	// Ambil SDM dari CSIRT pertama (1 perusahaan = 1 CSIRT)
	data, err := h.service.GetByCsirt(csirtList[0].ID)
	if err != nil {
		logger.Error(err, "failed to get SDM CSIRT by csirt")
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetSDMByID godoc
// @Summary      Ambil sdm csirt berdasarkan ID
// @Description  Mengambil satu data sdm csirt
// @Tags         SDM
// @Produce      json
// @Param        id   path      string  true  "SDM ID"
// @Success      200  {object} dto.SdmCsirtResponse
// @Failure      404  {object} dto.ErrorResponse
// @Router       /api/sdm_csirt/{id} [get]
func (h *SdmCsirtHandler) handleGetByID(w http.ResponseWriter, r *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		logger.Error(err, "failed to get SDM CSIRT by ID")
		utils.RespondError(w, 404, err.Error())
		return
	}

	role := middleware.GetRole(r.Context())
	if role == "user" {
		idPerusahaan := middleware.GetIDPerusahaan(r.Context())
		if !h.sdmBelongsToPerusahaan(data, idPerusahaan) {
			utils.RespondError(w, 403, "Anda tidak memiliki akses ke data ini")
			return
		}
	}

	utils.RespondJSON(w, 200, data)
}

// CreateSDM godoc
// @Summary      Tambah sdm csirt baru
// @Description  Membuat record sdm csirt
// @Tags         SDM
// @Accept       json
// @Produce      json
// @Param        sdm body dto.CreateSdmCsirtRequest true "Data sdm csirt"
// @Success      201  {object} dto.SdmCsirtResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/sdm_csirt [post]
func (h *SdmCsirtHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSdmCsirtRequest
	json.NewDecoder(r.Body).Decode(&req)

	// Ownership check: user hanya bisa tambah SDM ke CSIRT miliknya
	role := middleware.GetRole(r.Context())
	if role == "user" {
		idPerusahaan := middleware.GetIDPerusahaan(r.Context())
		if idPerusahaan == "" {
			utils.RespondError(w, 403, "Akun Anda belum terhubung ke perusahaan")
			return
		}

		// Pastikan id_csirt yang dikirim milik perusahaan ini
		if req.IdCsirt == nil || *req.IdCsirt == "" {
			utils.RespondError(w, 400, "id_csirt wajib diisi")
			return
		}

		if !h.csirtBelongsToPerusahaan(*req.IdCsirt, idPerusahaan) {
			utils.RespondError(w, 403, "CSIRT tidak ditemukan atau bukan milik perusahaan Anda")
			return
		}
	}

	id, err := h.service.Create(req)
	if err != nil {
		logger.Error(err, "failed to create SDM CSIRT")
		utils.RespondError(w, 400, err.Error())
		return
	}

	resp, err := h.service.GetByID(id)
	if err != nil {
		logger.Error(err, "failed to get SDM CSIRT after create")
		utils.RespondError(w, 500, err.Error())
		return
	}

	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyCreate("sdm_csirt", resp, userID)

	utils.RespondJSON(w, 201, resp)
}

// UpdateSDM godoc
// @Summary      Update sdm csirt
// @Description  Mengubah data sdm csirt berdasarkan ID
// @Tags         SDM
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "SDM ID"
// @Param        sdm body      dto.UpdateSdmCsirtRequest true "Data update"
// @Success      200  {object} dto.SdmCsirtResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/sdm_csirt/{id} [put]
func (h *SdmCsirtHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	role := middleware.GetRole(r.Context())
	if role == "user" {
		existing, err := h.service.GetByID(id)
		if err != nil {
			utils.RespondError(w, 404, "Data tidak ditemukan")
			return
		}
		idPerusahaan := middleware.GetIDPerusahaan(r.Context())
		if !h.sdmBelongsToPerusahaan(existing, idPerusahaan) {
			utils.RespondError(w, 403, "Anda tidak memiliki akses ke data ini")
			return
		}
	}

	var req dto.UpdateSdmCsirtRequest
	json.NewDecoder(r.Body).Decode(&req)

	if err := h.service.Update(id, req); err != nil {
		logger.Error(err, "failed to update SDM CSIRT")
		utils.RespondError(w, 400, err.Error())
		return
	}

	resp, err := h.service.GetByID(id)
	if err != nil {
		logger.Error(err, "failed to get SDM CSIRT after update")
		utils.RespondError(w, 500, err.Error())
		return
	}

	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyUpdate("sdm_csirt", resp, userID)

	utils.RespondJSON(w, 200, resp)
}

// DeleteSDM godoc
// @Summary      Hapus sdm csirt
// @Description  Menghapus data sdm csirt berdasarkan ID
// @Tags         SDM
// @Produce      json
// @Param        id  path  string  true  "SDM ID"
// @Success      200  {object} dto.MessageResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/sdm_csirt/{id} [delete]
func (h *SdmCsirtHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	role := middleware.GetRole(r.Context())
	if role == "user" {
		existing, err := h.service.GetByID(id)
		if err != nil {
			utils.RespondError(w, 404, "Data tidak ditemukan")
			return
		}
		idPerusahaan := middleware.GetIDPerusahaan(r.Context())
		if !h.sdmBelongsToPerusahaan(existing, idPerusahaan) {
			utils.RespondError(w, 403, "Anda tidak memiliki akses ke data ini")
			return
		}
	}

	if err := h.service.Delete(id); err != nil {
		logger.Error(err, "failed to delete SDM CSIRT")
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}

// sdmBelongsToPerusahaan cek apakah SDM ini milik perusahaan user
func (h *SdmCsirtHandler) sdmBelongsToPerusahaan(sdm *dto.SdmCsirtResponse, idPerusahaan string) bool {
	if sdm.Csirt == nil {
		return false
	}
	csirtList, err := h.csirtService.GetByPerusahaan(idPerusahaan)
	if err != nil {
		return false
	}
	for _, c := range csirtList {
		if c.ID == sdm.Csirt.ID {
			return true
		}
	}
	return false
}

// csirtBelongsToPerusahaan cek apakah CSIRT tertentu milik perusahaan user
func (h *SdmCsirtHandler) csirtBelongsToPerusahaan(idCsirt, idPerusahaan string) bool {
	csirtList, err := h.csirtService.GetByPerusahaan(idPerusahaan)
	if err != nil {
		return false
	}
	for _, c := range csirtList {
		if c.ID == idCsirt {
			return true
		}
	}
	return false
}