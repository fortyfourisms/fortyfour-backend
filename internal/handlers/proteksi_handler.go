package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
)

type ProteksiHandler struct {
	service *services.ProteksiService
}

func NewProteksiHandler(service *services.ProteksiService) *ProteksiHandler {
	return &ProteksiHandler{service: service}
}

func (h *ProteksiHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProteksiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	proteksi, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 201, proteksi)
}

func (h *ProteksiHandler) GetAllProteksi(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}

	utils.RespondJSON(w, 200, data)
}

func (h *ProteksiHandler) GetProteksiByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/proteksi/")
	if id == "" {
		utils.RespondError(w, 400, "ID wajib")
		return
	}

	result, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}

	utils.RespondJSON(w, 200, result)
}

func (h *ProteksiHandler) UpdateProteksi(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/proteksi/")
	if id == "" {
		utils.RespondError(w, 400, "ID wajib")
		return
	}

	var req dto.UpdateProteksiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	result, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 200, result)
}

func (h *ProteksiHandler) DeleteProteksi(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/proteksi/")
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