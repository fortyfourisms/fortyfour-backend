package handlers

import (
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"net/http"
	"strings"
)

type IdentifikasiHandler struct {
	service *services.IdentifikasiService
}

func NewIdentifikasiHandler(service *services.IdentifikasiService) *IdentifikasiHandler {
	return &IdentifikasiHandler{service: service}
}

func (h *IdentifikasiHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateIdentifikasiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	i, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 201, i)
}

func (h *IdentifikasiHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

func (h *IdentifikasiHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/identifikasi/")
	if id == "" {
		utils.RespondError(w, 400, "ID wajib")
		return
	}

	i, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}

	utils.RespondJSON(w, 200, i)
}

func (h *IdentifikasiHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/identifikasi/")
	if id == "" {
		utils.RespondError(w, 400, "ID wajib")
		return
	}

	var req dto.UpdateIdentifikasiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	i, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 200, i)
}

func (h *IdentifikasiHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/identifikasi/")
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
