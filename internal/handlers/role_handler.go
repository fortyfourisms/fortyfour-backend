// File: internal/handlers/role_handler.go
package handlers

import (
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"net/http"
	"strings"

	"fortyfour-backend/pkg/logger"
)

type RoleHandler struct {
	service    *services.RoleService
	sseService *services.SSEService
}

func NewRoleHandler(service *services.RoleService, sseService *services.SSEService) *RoleHandler {
	return &RoleHandler{
		service:    service,
		sseService: sseService,
	}
}

func (h *RoleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/api/role"), "/")

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

// GetAllRole godoc
// @Summary      List semua role
// @Description  Mengambil seluruh data role
// @Tags         Role
// @Produce      json
// @Success      200  {array}  dto.RoleResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /api/role [get]
func (h *RoleHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetRoleByID godoc
// @Summary      Ambil role berdasarkan ID
// @Description  Mengambil satu data role
// @Tags         Role
// @Produce      json
// @Param        id   path      string  true  "Role ID"
// @Success      200  {object} dto.RoleResponse
// @Failure      404  {object} dto.ErrorResponse
// @Router       /api/role/{id} [get]
func (h *RoleHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}
	utils.RespondJSON(w, 200, data)
}

// CreateRole godoc
// @Summary      Tambah role baru
// @Description  Membuat record role
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        role body dto.CreateRoleRequest true "Data role"
// @Success      201  {object} dto.RoleResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/role [post]
func (h *RoleHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Create(req)
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, err.Error())
		return
	}

	// SSE Notif Create
	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyCreate("role", resp, userID)

	utils.RespondJSON(w, 201, resp)
}

// UpdateRole godoc
// @Summary      Update role
// @Description  Mengubah data role berdasarkan ID
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Role ID"
// @Param        role body      dto.UpdateRoleRequest true "Data update"
// @Success      200  {object} dto.RoleResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/role/{id} [put]
func (h *RoleHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req dto.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Update(id, req)
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, err.Error())
		return
	}

	// SSE Notif Update
	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyUpdate("role", resp, userID)

	utils.RespondJSON(w, 200, resp)
}

// DeleteRole godoc
// @Summary      Hapus role
// @Description  Menghapus data role berdasarkan ID
// @Tags         Role
// @Produce      json
// @Param        id  path  string  true  "Role ID"
// @Success      200  {object} dto.MessageResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/role/{id} [delete]
func (h *RoleHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.Delete(id); err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, err.Error())
		return
	}

	// SSE Notif Delete
	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyDelete("role", id, userID)

	utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}
