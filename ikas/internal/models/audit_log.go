package models

import (
	"encoding/json"
	"time"
)

type AuditLog struct {
	ID        int             `json:"id"`
	IkasID    string          `json:"ikas_id"`
	UserID    string          `json:"user_id"`
	UserName  string          `json:"user_name"`
	Action    string          `json:"action"`
	Changes   json.RawMessage `json:"changes"`
	CreatedAt time.Time       `json:"created_at"`
}
