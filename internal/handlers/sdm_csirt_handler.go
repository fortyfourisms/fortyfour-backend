package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
)

type SdmCsirtHandler struct {
	service *services.SdmCsirtService
}

func NewSdmCsirtHandler(service *services.SdmCsirtService) *SdmCsirtHandler {
	return &SdmCsirtHandler{service: service}
}

func (h *SdmCsirtHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/sdm_csirt")
	id := strings.Trim(path, "/")

	switch r.Method {
	case http.MethodGet:
		if id == "" {
			h.handleGetAll(w)
		} else {
			h.handleGetByID(w, id)
		}
	case http.MethodPost:
		h.handleCreate(w, r)
	case http.MethodPut:
		h.handleUpdate(w, r, id)
	case http.MethodDelete:
		h.handleDelete(w, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *SdmCsirtHandler) handleGetAll(w http.ResponseWriter) {
	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

func (h *SdmCsirtHandler) handleGetByID(w http.ResponseWriter, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, 404, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

func (h *SdmCsirtHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSdmCsirtRequest
	json.NewDecoder(r.Body).Decode(&req)

	id, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 201, map[string]string{"id": id})
}

func (h *SdmCsirtHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req dto.UpdateSdmCsirtRequest
	json.NewDecoder(r.Body).Decode(&req)

	if err := h.service.Update(id, req); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 200, map[string]string{"message": "Update success"})
}

func (h *SdmCsirtHandler) handleDelete(w http.ResponseWriter, id string) {
	if err := h.service.Delete(id); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}
