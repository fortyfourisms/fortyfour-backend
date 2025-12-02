package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
)

type PerusahaanHandler struct {
	service *services.PerusahaanService
}

func NewPerusahaanHandler(service *services.PerusahaanService) *PerusahaanHandler {
	return &PerusahaanHandler{service: service}
}

func (h *PerusahaanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.PerusahaanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	p, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 201, p)
}

func (h *PerusahaanHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

func (h *PerusahaanHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/perusahaan/")
	if id == "" {
		utils.RespondError(w, 400, "ID wajib")
		return
	}

	p, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}

	utils.RespondJSON(w, 200, p)
}

func (h *PerusahaanHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/perusahaan/")
	if id == "" {
		http.Error(w, "ID wajib di path", http.StatusBadRequest)
		return
	}

	var req dto.PerusahaanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	p, err := h.service.Update(id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(p)
}

func (h *PerusahaanHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/perusahaan/")
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
