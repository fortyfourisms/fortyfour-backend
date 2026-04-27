package models

import "time"

type EditRequestStatus string

const (
	EditRequestPending  EditRequestStatus = "pending"
	EditRequestApproved EditRequestStatus = "approved"
	EditRequestRejected EditRequestStatus = "rejected"
)

type EditRequest struct {
	ID           string            `json:"id"`

	// relasi utama
	RespondenID  int               `json:"responden_id"`
	RisikoID     int               `json:"risiko_id"`

	// siapa yang request (user login)
	UserID       string            `json:"user_id"`

	// status workflow
	Status       EditRequestStatus `json:"status"`

	// catatan
	CatatanUser  *string           `json:"catatan_user,omitempty"` // dari user
	Catatan      *string           `json:"catatan,omitempty"`      // dari admin

	// JSON perubahan
	DataPerubahan string           `json:"data_perubahan"`

	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}