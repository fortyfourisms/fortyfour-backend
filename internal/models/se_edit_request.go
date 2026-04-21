package models

import "time"

type SEEditRequestStatus string

const (
	SEEditRequestPending  SEEditRequestStatus = "pending"
	SEEditRequestApproved SEEditRequestStatus = "approved"
	SEEditRequestRejected SEEditRequestStatus = "rejected"
)

type SEEditRequest struct {
	ID            string              `json:"id"`
	IDSE          string              `json:"id_se"`
	IDUser        string              `json:"id_user"`
	Status        SEEditRequestStatus `json:"status"`
	CatatanUser   *string             `json:"catatan_user,omitempty"`
	Catatan       *string             `json:"catatan,omitempty"` // catatan admin saat review
	DataPerubahan string              `json:"data_perubahan"`    // JSON string
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
}
