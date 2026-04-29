package handlers

import (
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"fortyfour-backend/internal/validator"
	"net/http"
	"strconv"
	"strings"

	"fortyfour-backend/pkg/logger"
)

type BeritaHandler struct {
	service services.BeritaServiceInterface
}

func NewBeritaHandler(service services.BeritaServiceInterface) *BeritaHandler {
	return &BeritaHandler{service: service}
}

func (h *BeritaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/berita")
	idStr := strings.TrimPrefix(path, "/")

	switch r.Method {
	case http.MethodGet:
		if idStr == "" {
			h.handleGetAll(w, r)
		} else {
			h.handleGetByID(w, r, idStr)
		}
	case http.MethodPost:
		h.handleCreate(w, r)
	case http.MethodPut:
		h.handleUpdate(w, r, idStr)
	case http.MethodDelete:
		h.handleDelete(w, r, idStr)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *BeritaHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		logger.Error(err, "failed to get all berita")
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.JSONResponse{
		Status:  "success",
		Message: "Berhasil mengambil data berita",
		Total:   len(data),
		Data:    data,
	})
}

func (h *BeritaHandler) handleGetByID(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	data, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.JSONResponse{
		Status:  "success",
		Message: "Berhasil mengambil detail berita",
		Data:    data,
	})
}

func (h *BeritaHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBeritaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}

	// Validation
	if err := validator.Validate(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	authorID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// XSS Sanitization for Deskripsi (Best practice for longtext/HTML)
	req.Deskripsi = validator.SanitizeHTML(req.Deskripsi)

	if err := h.service.Create(authorID, req); err != nil {
		logger.Error(err, "failed to create berita")
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, utils.JSONResponse{
		Status:  "success",
		Message: "Permintaan pembuatan berita telah diterima dan sedang diproses",
	})
}

func (h *BeritaHandler) handleUpdate(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	var req dto.UpdateBeritaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}

	// Validation
	if err := validator.Validate(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Deskripsi != nil {
		sanitized := validator.SanitizeHTML(*req.Deskripsi)
		req.Deskripsi = &sanitized
	}

	if err := h.service.Update(id, req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.JSONResponse{
		Status:  "success",
		Message: "Permintaan pembaruan berita telah diterima dan sedang diproses",
	})
}

func (h *BeritaHandler) handleDelete(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.JSONResponse{
		Status:  "success",
		Message: "Berhasil menghapus berita",
	})
}
