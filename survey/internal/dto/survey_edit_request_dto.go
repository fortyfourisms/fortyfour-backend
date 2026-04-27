package dto

import "survey/internal/models"

// CREATE REQUEST (USER)
type CreateEditRequestDTO struct {
	RespondenID int                    `json:"responden_id"`
	RisikoID    int                    `json:"risiko_id"`
	DataPerubahan map[string]interface{} `json:"data_perubahan"` // FLEXIBLE JSON
	CatatanUser *string                `json:"catatan_user,omitempty"`
}

// REVIEW (ADMIN)
type ReviewEditRequestDTO struct {
	Status  string  `json:"status"` // approved / rejected
	Catatan *string `json:"catatan,omitempty"`
}

// RESPONSE
type EditRequestResponse struct {
	ID           string `json:"id"`
	RespondenID  int    `json:"responden_id"`
	RisikoID     int    `json:"risiko_id"`
	UserID       string `json:"user_id"`

	Status       models.EditRequestStatus `json:"status"`

	CatatanUser  *string `json:"catatan_user,omitempty"`
	Catatan      *string `json:"catatan,omitempty"`

	DataPerubahan map[string]interface{} `json:"data_perubahan,omitempty"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}