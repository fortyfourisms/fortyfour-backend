package handlers

import (
	"ikas/internal/dto"
	"ikas/internal/services"
	"ikas/internal/utils"
	"net/http"
)

type AuditLogHandler struct {
	service *services.AuditLogService
}

func NewAuditLogHandler(service *services.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{service: service}
}

func (h *AuditLogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ikasID := r.URL.Query().Get("ikas_id")
	h.handleGetAll(w, r, ikasID)
}

// GetAllAuditLogs godoc
//
//	@Summary		List semua audit logs
//	@Description	Mengambil seluruh riwayat perubahan IKAS (opsional filter berdasarkan ikas_id)
//	@Tags			Audit Logs
//	@Produce		json
//	@Param			ikas_id	query		string	false	"IKAS ID"
//	@Success		200		{array}		dto.AuditLogResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/api/maturity/ikas-audit-logs [get]
func (h *AuditLogHandler) handleGetAll(w http.ResponseWriter, r *http.Request, ikasID string) {
	var data []dto.AuditLogResponse
	var err error

	if ikasID != "" {
		data, err = h.service.GetAuditLogsByIkasID(ikasID)
	} else {
		data, err = h.service.GetAllAuditLogs()
	}

	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}

	if data == nil {
		data = []dto.AuditLogResponse{}
	}

	utils.RespondJSON(w, 200, utils.JSONResponse{
		Status:  "success",
		Message: "Berhasil mengambil data audit logs",
		Data:    data,
		Total:   len(data),
	})
}
