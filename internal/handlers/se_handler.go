package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
)

type SEHandler struct {
	service    services.SEService
	sseService services.SSEServiceInterface
}

func NewSEHandler(
	service services.SEService,
	sseService services.SSEServiceInterface,
) *SEHandler {
	return &SEHandler{
		service:    service,
		sseService: sseService,
	}
}

func (h *SEHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/api/se"), "/")

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

// @Summary Get all SE
// @Description Admin mendapat semua SE. User hanya mendapat SE milik perusahaannya.
// @Tags SE
// @Accept json
// @Produce json
// @Success 200 {object} dto.SEListResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/se [get]
func (h *SEHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())

	if role == "admin" {
		// Admin → seluruh data
		data, err := h.service.GetAll()
		if err != nil {
			utils.RespondError(w, 500, err.Error())
			return
		}
		utils.RespondJSON(w, 200, data)
		return
	}

	// User biasa → hanya data perusahaannya
	idPerusahaan := middleware.GetIDPerusahaan(r.Context())
	if idPerusahaan == "" {
		utils.RespondError(w, 403, "Akun Anda belum terhubung ke perusahaan")
		return
	}

	data, err := h.service.GetByPerusahaan(idPerusahaan)
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// @Summary Get SE by ID
// @Description Get sistem elektronik by ID. User hanya bisa akses SE milik perusahaannya.
// @Tags SE
// @Accept json
// @Produce json
// @Param id path string true "SE ID"
// @Success 200 {object} dto.SEResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/se/{id} [get]
func (h *SEHandler) handleGetByID(w http.ResponseWriter, r *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}

	// Validasi ownership untuk user
	role := middleware.GetRole(r.Context())
	if role == "user" {
		idPerusahaan := middleware.GetIDPerusahaan(r.Context())
		if data.IDPerusahaan != idPerusahaan {
			utils.RespondError(w, 403, "Anda tidak memiliki akses ke data ini")
			return
		}
	}

	utils.RespondJSON(w, 200, data)
}

// @Summary Create SE
// @Description Create new sistem elektronik. User biasa otomatis menggunakan id_perusahaan dari token.
// @Tags SE
// @Accept json
// @Produce json
// @Param request body dto.CreateSERequest true "SE Create Request"
// @Success 201 {object} dto.SEResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Router /api/se [post]
func (h *SEHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSERequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	// User biasa: paksa id_perusahaan dari JWT, tidak bisa diisi sembarangan
	role := middleware.GetRole(r.Context())
	if role == "user" {
		idPerusahaan := middleware.GetIDPerusahaan(r.Context())
		if idPerusahaan == "" {
			utils.RespondError(w, 403, "Akun Anda belum terhubung ke perusahaan")
			return
		}
		req.IDPerusahaan = idPerusahaan
	}

	resp, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyCreate("se", resp, userID)

	utils.RespondJSON(w, 201, resp)
}

// @Summary Update SE
// @Description Update sistem elektronik. User biasa hanya bisa update SE milik perusahaannya.
// @Tags SE
// @Accept json
// @Produce json
// @Param id path string true "SE ID"
// @Param request body dto.UpdateSERequest true "SE Update Request"
// @Success 200 {object} dto.SEResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Router /api/se/{id} [put]
func (h *SEHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	// Validasi ownership sebelum update untuk user
	role := middleware.GetRole(r.Context())
	if role == "user" {
		existing, err := h.service.GetByID(id)
		if err != nil {
			utils.RespondError(w, 404, "Data tidak ditemukan")
			return
		}
		idPerusahaan := middleware.GetIDPerusahaan(r.Context())
		if existing.IDPerusahaan != idPerusahaan {
			utils.RespondError(w, 403, "Anda tidak memiliki akses ke data ini")
			return
		}
	}

	var req dto.UpdateSERequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyUpdate("se", resp, userID)

	utils.RespondJSON(w, 200, resp)
}

// @Summary Delete SE
// @Description Delete sistem elektronik. User biasa hanya bisa hapus SE milik perusahaannya.
// @Tags SE
// @Accept json
// @Produce json
// @Param id path string true "SE ID"
// @Success 200 {object} map[string]string
// @Failure 403 {object} dto.ErrorResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/se/{id} [delete]
func (h *SEHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	// Validasi ownership sebelum delete untuk user
	role := middleware.GetRole(r.Context())
	if role == "user" {
		existing, err := h.service.GetByID(id)
		if err != nil {
			utils.RespondError(w, 404, "Data tidak ditemukan")
			return
		}
		idPerusahaan := middleware.GetIDPerusahaan(r.Context())
		if existing.IDPerusahaan != idPerusahaan {
			utils.RespondError(w, 403, "Anda tidak memiliki akses ke data ini")
			return
		}
	}

	if err := h.service.Delete(id); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyDelete("se", id, userID)

	utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}
