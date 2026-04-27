package dto

type CreateEventRequest struct {
	Judul     string `json:"judul" validate:"required,min=5,max=255"`
	Deskripsi string `json:"deskripsi" validate:"required"`
	Tanggal   string `json:"tanggal" validate:"required" example:"2024-12-31T15:00:00Z"`
}

type UpdateEventRequest struct {
	Judul     *string `json:"judul,omitempty" validate:"omitempty,min=5,max=255"`
	Deskripsi *string `json:"deskripsi,omitempty" validate:"omitempty"`
	Tanggal   *string `json:"tanggal,omitempty" validate:"omitempty"`
}

type EventResponse struct {
	ID        int64  `json:"id"`
	Judul     string `json:"judul"`
	Deskripsi string `json:"deskripsi"`
	Tanggal   string `json:"tanggal"`
	Status    string `json:"status"` // upcoming atau past
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
