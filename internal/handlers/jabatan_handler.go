package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
)

type JabatanHandler struct {
	service *services.JabatanService
}

func NewJabatanHandler(service *services.JabatanService) *JabatanHandler {
	return &JabatanHandler{service: service}
}

func (h *JabatanHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/jabatan")
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
			j, err := h.service.GetByID(id)
			if err != nil {
				utils.RespondError(w, 404, "Data tidak ditemukan")
				return
			}
			utils.RespondJSON(w, 200, j)
		}

	case http.MethodPost:
		var req dto.CreateJabatanRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondError(w, 400, "Invalid request body")
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

		var req dto.UpdateJabatanRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondError(w, 400, "Invalid request body")
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
