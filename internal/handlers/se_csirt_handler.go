package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
)

type SeCsirtHandler struct {
    service *services.SeCsirtService
}

func NewSeCsirtHandler(service *services.SeCsirtService) *SeCsirtHandler {
    return &SeCsirtHandler{service: service}
}

func (h *SeCsirtHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/api/se_csirt"), "/")

    switch r.Method {
    case http.MethodGet:
        if id == "" {
            h.handleGetAll(w)
        } else {
            h.handleGetByID(w, id)
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

func (h *SeCsirtHandler) handleGetAll(w http.ResponseWriter) {
    data, err := h.service.GetAll()
    if err != nil {
        utils.RespondError(w, 400, err.Error())
        return
    }
    utils.RespondJSON(w, 200, data)
}

func (h *SeCsirtHandler) handleGetByID(w http.ResponseWriter, id string) {
    data, err := h.service.GetByID(id)
    if err != nil {
        utils.RespondError(w, 400, err.Error())
        return
    }
    utils.RespondJSON(w, 200, data)
}

func (h *SeCsirtHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
    var req dto.CreateSeCsirtRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.RespondError(w, 400, err.Error())
        return
    }
    id, err := h.service.Create(req)
    if err != nil {
        utils.RespondError(w, 400, err.Error())
        return
    }
    utils.RespondJSON(w, 200, map[string]string{"id": id})
}

func (h *SeCsirtHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
    var req dto.SeCsirtResponse
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.RespondError(w, 400, err.Error())
        return
    }
    if err := h.service.Update(id, req); err != nil {
        utils.RespondError(w, 400, err.Error())
        return
    }
    utils.RespondJSON(w, 200, map[string]string{"message": "Update success"})
}

func (h *SeCsirtHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
    if err := h.service.Delete(id); err != nil {
        utils.RespondError(w, 400, err.Error())
        return
    }
    utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}
