package handlers

import (
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"net/http"
	"strings"
)

type GulihHandler struct {
	service *services.GulihService
}

func NewGulihHandler(service *services.GulihService) *GulihHandler {
	return &GulihHandler{service: service}
}

func (h *GulihHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateGulihRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	g, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 201, g)
}

func (h *GulihHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

func (h *GulihHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/gulih/")
	if id == "" {
		utils.RespondError(w, 400, "ID wajib")
		return
	}

	g, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}

	utils.RespondJSON(w, 200, g)
}

func (h *GulihHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/gulih/")
	if id == "" {
		utils.RespondError(w, 400, "ID wajib")
		return
	}

	var req dto.UpdateGulihRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	g, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 200, g)
}

func (h *GulihHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/gulih/")
	if id == "" {
		utils.RespondError(w, 400, "ID wajib")
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}
