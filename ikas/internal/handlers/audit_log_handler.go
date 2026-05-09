package handlers

import (
	"ikas/internal/dto"
	"ikas/internal/services"
	"ikas/internal/utils"
	"net/http"
	"strconv"
)

const (
	defaultAuditLogPage  = 1
	defaultAuditLogLimit = 20
	maxAuditLogLimit     = 100
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

	h.handleGetAll(w, r)
}

// GetAllAuditLogs godoc
//
//	@Summary		List semua audit logs IKAS (paginated)
//	@Description	Mengambil seluruh riwayat perubahan IKAS dengan pagination. Opsional filter berdasarkan ikas_id.
//	@Tags			Audit Logs
//	@Produce		json
//	@Param			ikas_id	query		string	false	"Filter berdasarkan IKAS ID"
//	@Param			page	query		int		false	"Halaman (default: 1)"					example(1)
//	@Param			limit	query		int		false	"Jumlah data per halaman, maks 100 (default: 20)"	example(10)
//	@Success		200		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/api/maturity/ikas-audit-logs [get]
func (h *AuditLogHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page := parsePositiveInt(q.Get("page"), defaultAuditLogPage)
	limit := parsePositiveInt(q.Get("limit"), defaultAuditLogLimit)

	// Enforce maximum limit to prevent DB overload
	if limit > maxAuditLogLimit {
		limit = maxAuditLogLimit
	}

	req := dto.AuditLogListRequest{
		IkasID: q.Get("ikas_id"),
		Page:   page,
		Limit:  limit,
	}

	result, err := h.service.GetAuditLogs(req)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.PaginatedJSONResponse{
		Status:  "success",
		Message: "Berhasil mengambil data audit logs IKAS",
		Pagination: utils.PaginationMeta{
			Total:      result.Total,
			Page:       result.Page,
			Limit:      result.Limit,
			TotalPages: result.TotalPages,
		},
		Data: result.Data,
	})
}

// parsePositiveInt parses a string to a positive int, returning the fallback if invalid.
func parsePositiveInt(s string, fallback int) int {
	if s == "" {
		return fallback
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 1 {
		return fallback
	}
	return v
}
