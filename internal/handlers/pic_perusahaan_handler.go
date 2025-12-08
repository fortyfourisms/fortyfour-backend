package handlers

import (
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"net/http"
	"strings"
)

type PICHandler struct {
	service *services.PICService
}

func NewPICHandler(service *services.PICService) *PICHandler {
	return &PICHandler{service: service}
}

func (h *PICHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Correct the path to match the router
	id := strings.TrimPrefix(r.URL.Path, "/api/pic")
	id = strings.TrimPrefix(id, "/")

	switch r.Method {
	case http.MethodGet:
		if id == "" {
			data, err := h.service.GetAll()
			if err != nil {
				utils.RespondError(w, 500, err.Error())
				return
			}
			utils.RespondJSON(w, 200, data)
		} else {
			p, err := h.service.GetByID(id)
			if err != nil {
				utils.RespondError(w, 404, "Data tidak ditemukan")
				return
			}
			utils.RespondJSON(w, 200, p)
		}

	case http.MethodPost:
		var req dto.CreatePICRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondError(w, 400, "Invalid JSON")
			return
		}

		resp, err := h.service.Create(req)
		if err != nil {
			utils.RespondError(w, 400, err.Error())
			return
		}

		utils.RespondJSON(w, 201, resp)

	case http.MethodPut:
		if id == "" {
			utils.RespondError(w, 400, "ID wajib")
			return
		}

		var req dto.UpdatePICRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondError(w, 400, "Invalid JSON")
			return
		}

		resp, err := h.service.Update(id, req)
		if err != nil {
			utils.RespondError(w, 400, err.Error())
			return
		}

		utils.RespondJSON(w, 200, resp)

	case http.MethodDelete:
		if id == "" {
			utils.RespondError(w, 400, "ID wajib")
			return
		}

		if err := h.service.Delete(id); err != nil {
			utils.RespondError(w, 400, err.Error())
			return
		}

		utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
