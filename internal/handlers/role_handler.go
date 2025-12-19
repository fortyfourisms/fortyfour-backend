// File: internal/handlers/role_handler.go
package handlers

import (
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"net/http"
	"strings"
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

func (h *RoleHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

func (h *RoleHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}
	utils.RespondJSON(w, 200, data)
}

func (h *RoleHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	// SSE Notif Create
	userID := ""
	if uid := r.Context().Value("user_id"); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyCreate("role", resp, userID)

	utils.RespondJSON(w, 201, resp)
}

func (h *RoleHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req dto.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	// SSE Notif Update
	userID := ""
	if uid := r.Context().Value("user_id"); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyUpdate("role", resp, userID)

	utils.RespondJSON(w, 200, resp)
}

func (h *RoleHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.Delete(id); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	// SSE Notif Delete
	userID := ""
	if uid := r.Context().Value("user_id"); uid != nil {
		userID = uid.(string)
	}
	h.sseService.NotifyDelete("role", id, userID)

	utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}
