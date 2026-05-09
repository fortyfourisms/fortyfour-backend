package handlers

import (
	"encoding/json"
	"fmt"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"fortyfour-backend/internal/validator"
	"net"
	"net/http"
	"strconv"
	"strings"

	"fortyfour-backend/pkg/logger"
)

type EventHandler struct {
	service            services.EventServiceInterface
	turnstileValidator *utils.TurnstileValidator
}

func NewEventHandler(service services.EventServiceInterface, turnstileValidator *utils.TurnstileValidator) *EventHandler {
	return &EventHandler{
		service:            service,
		turnstileValidator: turnstileValidator,
	}
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

// ListEvents godoc
//
//	@Summary		List semua event
//	@Description	Mengambil daftar semua event/kegiatan yang tersedia
//	@Tags			Event
//	@Produce		json
//	@Success		200	{object}	utils.JSONResponse{data=[]dto.EventResponse}
//	@Router			/api/kegiatan [get]
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

// GetEventByID godoc
//
//	@Summary		Detail event
//	@Description	Mengambil detail satu event berdasarkan ID
//	@Tags			Event
//	@Produce		json
//	@Param			id	path		string	true	"Event ID"
//	@Success		200	{object}	utils.JSONResponse{data=dto.EventResponse}
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/api/kegiatan/{id} [get]
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

// CreateEvent godoc
//
//	@Summary		Tambah event baru
//	@Description	Membuat record event baru (proses via RabbitMQ)
//	@Tags			Event
//	@Accept			json
//	@Produce		json
//	@Param			event	body		dto.CreateEventRequest	true	"Data event"
//	@Success		201		{object}	utils.JSONResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Router			/api/kegiatan [post]
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

// UpdateEvent godoc
//
//	@Summary		Update event
//	@Description	Mengubah data event berdasarkan ID (proses via RabbitMQ)
//	@Tags			Event
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Event ID"
//	@Param			event	body		dto.UpdateEventRequest	true	"Data event yang diubah"
//	@Success		200		{object}	utils.JSONResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Router			/api/kegiatan/{id} [put]
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

// DeleteEvent godoc
//
//	@Summary		Hapus event
//	@Description	Menghapus data event berdasarkan ID
//	@Tags			Event
//	@Produce		json
//	@Param			id	path		string	true	"Event ID"
//	@Success		200	{object}	utils.JSONResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Router			/api/kegiatan/{id} [delete]
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

// RegisterEvent godoc
//
//	@Summary		Registrasi event
//	@Description	Mendaftarkan peserta ke sebuah event
//	@Tags			Event
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string								true	"Event ID"
//	@Param			registration	body		dto.CreateEventRegistrationRequest	true	"Data registrasi"
//	@Success		201				{object}	utils.JSONResponse{data=dto.EventRegistrationResponse}
//	@Failure		400,409			{object}	dto.ErrorResponse
//	@Router			/api/kegiatan/{id}/registrasi [post]
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

	// Validasi Turnstile
	if err := h.validateTurnstile(r, req.TurnstileToken); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validasi dilakukan sebelum sanitasi agar error message merujuk ke input asli
	if err := validator.Validate(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// XSS sanitization — email hanya di-TrimSpace karena html.EscapeString
	// akan merusak karakter email yang valid (misalnya ' menjadi &#39;)
	req.Nama = validator.SanitizeHTML(req.Nama)
	req.Email = strings.TrimSpace(req.Email)
	req.Perusahaan = validator.SanitizeHTML(req.Perusahaan)
	req.Jabatan = validator.SanitizeHTML(req.Jabatan)
	req.NoHP = strings.TrimSpace(req.NoHP)
	req.Sektor = validator.SanitizeHTML(req.Sektor)

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

// DownloadRegistrationPDF godoc
//
//	@Summary		Download PDF Registrasi
//	@Description	Mengunduh bukti registrasi event dalam format PDF
//	@Tags			Event
//	@Produce		application/pdf
//	@Param			id	path	string	true	"Registration ID"
//	@Success		200	{file}	binary
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/api/kegiatan/registrasi/{id}/download [get]
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

// validateTurnstile is a helper to verify Cloudflare Turnstile token
func (h *EventHandler) validateTurnstile(r *http.Request, token string) error {
	if h.turnstileValidator == nil {
		return fmt.Errorf("sistem keamanan Turnstile tidak terkonfigurasi")
	}

	remoteIP := r.Header.Get("CF-Connecting-IP")
	if remoteIP == "" {
		remoteIP = r.Header.Get("X-Forwarded-For")
	}
	if remoteIP != "" {
		// Get the first IP if it's a comma-separated list
		remoteIP = strings.TrimSpace(strings.Split(remoteIP, ",")[0])
	} else {
		remoteIP = r.RemoteAddr
		// Clean up port from RemoteAddr properly
		if host, _, err := net.SplitHostPort(remoteIP); err == nil {
			remoteIP = host
		}
	}

	turnstileResp, err := h.turnstileValidator.Validate(token, remoteIP)
	if err != nil {
		return fmt.Errorf("gagal memvalidasi Turnstile")
	}

	if !turnstileResp.Success {
		return fmt.Errorf("verifikasi Turnstile gagal: %s", strings.Join(turnstileResp.ErrorCodes, ", "))
	}

	return nil
}
