package dto

import "fortyfour-backend/internal/models"

// CreateSEEditRequestDTO — user mengirim request edit SE
type CreateSEEditRequestDTO struct {
	DataPerubahan UpdateSERequest `json:"data_perubahan" validate:"required"`
	Catatan       *string         `json:"catatan,omitempty"`      // alasan user
	CatatanUser   *string         `json:"catatan_user,omitempty"` // alias dari frontend
}

// ReviewSEEditRequestDTO — admin approve/reject request
type ReviewSEEditRequestDTO struct {
	Status  string  `json:"status" validate:"required,oneof=approved rejected"`
	Catatan *string `json:"catatan,omitempty"`
}

// SEEditRequestResponse — response untuk list/detail request
type SEEditRequestResponse struct {
	ID            string                     `json:"id"`
	IDSE          string                     `json:"id_se"`
	IDUser        string                     `json:"id_user"`
	NamaUser      string                     `json:"nama_user,omitempty"`
	NamaSE        string                     `json:"nama_se,omitempty"`
	Status        models.SEEditRequestStatus `json:"status"`
	CatatanUser   *string                    `json:"catatan_user,omitempty"`
	Catatan       *string                    `json:"catatan,omitempty"` // catatan admin
	DataPerubahan *UpdateSERequest           `json:"data_perubahan,omitempty"`
	CreatedAt     string                     `json:"created_at"`
	UpdatedAt     string                     `json:"updated_at"`
}
