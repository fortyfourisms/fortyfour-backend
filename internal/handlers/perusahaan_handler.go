package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
)

type PerusahaanHandler struct {
	perusahaanService *services.PerusahaanService
}

func NewPerusahaanHandler(service *services.PerusahaanService) *PerusahaanHandler {
	return &PerusahaanHandler{perusahaanService: service}
}

func (h *PerusahaanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.PerusahaanRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	p, err := h.perusahaanService.Create(req)
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}

	utils.RespondJSON(w, 201, p)
}

func (h *PerusahaanHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	perusahaan, err := h.perusahaanService.GetAll()
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, perusahaan)
}

func (h *PerusahaanHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("id")
	id, err := strconv.Atoi(query)
	if err != nil {
		utils.RespondError(w, 400, "Invalid ID")
		return
	}

	p, err := h.perusahaanService.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			utils.RespondError(w, 404, "Data not found")
			return
		}
		utils.RespondError(w, 500, err.Error())
		return
	}

	utils.RespondJSON(w, 200, p)
}

func (h *PerusahaanHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req dto.PerusahaanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	queryID := r.URL.Query().Get("id")
	id, err := strconv.Atoi(queryID)
	if err != nil {
		utils.RespondError(w, 400, "Invalid ID")
		return
	}

	err = h.perusahaanService.Update(id, req)
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}

	utils.RespondJSON(w, 200, map[string]string{"message": "Update success"})
}

func (h *PerusahaanHandler) Delete(w http.ResponseWriter, r *http.Request) {
	queryID := r.URL.Query().Get("id")
	id, err := strconv.Atoi(queryID)
	if err != nil {
		utils.RespondError(w, 400, "Invalid ID")
		return
	}

	err = h.perusahaanService.Delete(id)
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}

	utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}
