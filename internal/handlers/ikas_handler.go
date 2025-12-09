package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
)

type IkasHandler struct {
	service *services.IkasService
}

func NewIkasHandler(service *services.IkasService) *IkasHandler {
	return &IkasHandler{service: service}
}

// CREATE
func (h *IkasHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateIkasRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	ikas, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 201, ikas)
}

// GET ALL
func (h *IkasHandler) GetAllIkas(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}

	utils.RespondJSON(w, 200, data)
}

// GET BY ID
func (h *IkasHandler) GetIkasByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/ikas/")
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

// UPDATE
func (h *IkasHandler) UpdateIkas(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/ikas/")
	if id == "" {
		utils.RespondError(w, 400, "ID wajib")
		return
	}

	var req dto.UpdateIkasRequest
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

// DELETE 
func (h *IkasHandler) DeleteIkas(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/ikas/")
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