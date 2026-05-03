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
