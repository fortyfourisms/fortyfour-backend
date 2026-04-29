package handlers

import (
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"fortyfour-backend/internal/validator"
	"net/http"
	"strconv"
	"strings"

	"fortyfour-backend/pkg/logger"
)

type EventHandler struct {
	service services.EventServiceInterface
}

func NewEventHandler(service services.EventServiceInterface) *EventHandler {
	return &EventHandler{service: service}
}

func (h *EventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/kegiatan")
	path = strings.TrimPrefix(path, "/")

	switch {
	case strings.HasSuffix(path, "/registrasi") && r.Method == http.MethodPost:
		h.handleRegister(w, r, strings.TrimSuffix(path, "/registrasi"))
	case strings.HasPrefix(path, "registrasi/") && strings.HasSuffix(path, "/download") && r.Method == http.MethodGet:
		registrationID := strings.TrimPrefix(path, "registrasi/")
		registrationID = strings.TrimSuffix(registrationID, "/download")
		h.handleDownloadRegistrationPDF(w, r, registrationID)
	case path == "" && r.Method == http.MethodGet:
		h.handleGetAll(w, r)
	case path == "" && r.Method == http.MethodPost:
		h.handleCreate(w, r)
	case path != "" && r.Method == http.MethodGet:
		h.handleGetByID(w, r, path)
	case path != "" && r.Method == http.MethodPut:
		h.handleUpdate(w, r, path)
	case path != "" && r.Method == http.MethodDelete:
		h.handleDelete(w, r, path)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *EventHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		logger.Error(err, "failed to get all events")
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.JSONResponse{
		Status:  "success",
		Message: "Berhasil mengambil data event",
		Total:   len(data),
		Data:    data,
	})
}

func (h *EventHandler) handleGetByID(w http.ResponseWriter, r *http.Request, idStr string) {
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
		Message: "Berhasil mengambil detail event",
		Data:    data,
	})
}

func (h *EventHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validator.Validate(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// XSS Sanitization for Deskripsi
	req.Deskripsi = validator.SanitizeHTML(req.Deskripsi)

	if err := h.service.Create(req); err != nil {
		logger.Error(err, "failed to create event")
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, utils.JSONResponse{
		Status:  "success",
		Message: "Permintaan pembuatan event telah diterima dan sedang diproses",
	})
}

func (h *EventHandler) handleUpdate(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	var req dto.UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

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
		Message: "Permintaan pembaruan event telah diterima dan sedang diproses",
	})
}

func (h *EventHandler) handleDelete(w http.ResponseWriter, r *http.Request, idStr string) {
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
		Message: "Berhasil menghapus event",
	})
}

func (h *EventHandler) handleRegister(w http.ResponseWriter, r *http.Request, eventIDStr string) {
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID event tidak valid")
		return
	}

	var req dto.CreateEventRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.Nama = validator.SanitizeHTML(req.Nama)
	req.Email = validator.SanitizeHTML(req.Email)
	req.Perusahaan = validator.SanitizeHTML(req.Perusahaan)
	req.Jabatan = validator.SanitizeHTML(req.Jabatan)
	req.NoHP = validator.SanitizeHTML(req.NoHP)
	req.Sektor = validator.SanitizeHTML(req.Sektor)

	if err := validator.Validate(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.service.Register(eventID, req)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "event tidak ditemukan":
			status = http.StatusNotFound
		case "email sudah terdaftar pada event ini":
			status = http.StatusConflict
		default:
			if strings.Contains(err.Error(), "tidak valid") {
				status = http.StatusBadRequest
			}
		}
		utils.RespondError(w, status, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, utils.JSONResponse{
		Status:  "success",
		Message: "Berhasil registrasi event",
		Data:    resp,
	})
}

func (h *EventHandler) handleDownloadRegistrationPDF(w http.ResponseWriter, r *http.Request, registrationIDStr string) {
	registrationID, err := strconv.ParseInt(registrationIDStr, 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID registrasi tidak valid")
		return
	}

	pdfBytes, filename, err := h.service.DownloadRegistrationPDF(registrationID)
	if err != nil {
		if err.Error() == "registrasi event tidak ditemukan" {
			utils.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	w.Header().Set("Content-Length", strconv.Itoa(len(pdfBytes)))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pdfBytes)
}
