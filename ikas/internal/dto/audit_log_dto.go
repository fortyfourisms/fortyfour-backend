package dto

import "encoding/json"

type UserAuditLogResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AuditLogResponse struct {
	ID        int                  `json:"id"`
	IkasID    string               `json:"ikas_id"`
	User      UserAuditLogResponse `json:"user"`
	Action    string               `json:"action"`
	Changes   json.RawMessage      `json:"changes"`
	CreatedAt string               `json:"created_at"`
}

// AuditLogListRequest holds validated pagination and filter parameters.
type AuditLogListRequest struct {
	IkasID string
	Page   int
	Limit  int
}

// PaginatedAuditLogResponse is the standard paginated list envelope.
type PaginatedAuditLogResponse struct {
	Data       []AuditLogResponse `json:"data"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}
