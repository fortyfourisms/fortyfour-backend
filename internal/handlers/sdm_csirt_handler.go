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
    id := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/api/sdm_csirt"), "/")

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

// GetAllSDM godoc
// @Summary      List semua sdm csirt
// @Description  Mengambil seluruh data sdm csirt
// @Tags         SDM
// @Produce      json
// @Success      200  {array}  dto.SdmCsirtResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /api/sdm_csirt [get]
func (h *SdmCsirtHandler) handleGetAll(w http.ResponseWriter) {
    data, err := h.service.GetAll()
    if err != nil {
        utils.RespondError(w, 400, err.Error())
        return
    }
    utils.RespondJSON(w, 200, data)
}

// GetSDMByID godoc
// @Summary      Ambil sdm csirt berdasarkan ID
// @Description  Mengambil satu data sdm csirt
// @Tags         SDM
// @Produce      json
// @Param        id   path      string  true  "SDM ID"
// @Success      200  {object} dto.SdmCsirtResponse
// @Failure      404  {object} dto.ErrorResponse
// @Router       /api/sdm_csirt/{id} [get]
func (h *SdmCsirtHandler) handleGetByID(w http.ResponseWriter, id string) {
    data, err := h.service.GetByID(id)
    if err != nil {
        utils.RespondError(w, 400, err.Error())
        return
    }
    utils.RespondJSON(w, 200, data)
}

// CreateSDM godoc
// @Summary      Tambah sdm csirt baru
// @Description  Membuat record sdm csirt
// @Tags         SDM
// @Accept       json
// @Produce      json
// @Param        sdm body dto.CreateSdmCsirtRequest true "Data sdm csirt"
// @Success      201  {object} dto.SdmCsirtResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/sdm_csirt [post]
func (h *SdmCsirtHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
    var req dto.CreateSdmCsirtRequest
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

// UpdateSDM godoc
// @Summary      Update sdm csirt
// @Description  Mengubah data sdm csirt berdasarkan ID
// @Tags         SDM
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "SDM ID"
// @Param        sdm body      dto.UpdateSdmCsirtRequest true "Data update"
// @Success      200  {object} dto.SdmCsirtResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/sdm_csirt/{id} [put]
func (h *SdmCsirtHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
    var req dto.UpdateSdmCsirtRequest
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

// DeleteSDM godoc
// @Summary      Hapus sdm csirt
// @Description  Menghapus data sdm csirt berdasarkan ID
// @Tags         SDM
// @Produce      json
// @Param        id  path  string  true  "SDM ID"
// @Success      200  {object} dto.MessageResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/sdm_csirt/{id} [delete]
func (h *SdmCsirtHandler) handleDelete(w http.ResponseWriter, _ *http.Request, id string) {
    if err := h.service.Delete(id); err != nil {
        utils.RespondError(w, 400, err.Error())
        return
    }
    utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}
