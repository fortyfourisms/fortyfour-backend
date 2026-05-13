package dto

type CreateEventRequest struct {
	Judul     string `json:"judul" validate:"required,min=5,max=255"`
	Deskripsi string `json:"deskripsi" validate:"required"`
	Lokasi    string `json:"lokasi" validate:"required"`
	Tanggal   string `json:"tanggal" validate:"required" example:"2024-12-31T15:00:00Z"`
}

type UpdateEventRequest struct {
	Judul     *string `json:"judul,omitempty" validate:"omitempty,min=5,max=255"`
	Deskripsi *string `json:"deskripsi,omitempty" validate:"omitempty"`
	Lokasi    *string `json:"lokasi,omitempty" validate:"omitempty"`
	Tanggal   *string `json:"tanggal,omitempty" validate:"omitempty"`
}

type EventResponse struct {
	ID        string `json:"id"`
	Slug      string `json:"slug"`
	Judul     string `json:"judul"`
	Deskripsi string `json:"deskripsi"`
	Lokasi    string `json:"lokasi"`
	Tanggal   string `json:"tanggal"`
	Status    string `json:"status"` // upcoming atau past
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateEventRegistrationRequest struct {
	Nama           string `json:"nama" validate:"required,min=3,max=255"`
	Email          string `json:"email" validate:"required,email,max=255"`
	Perusahaan     string `json:"perusahaan" validate:"required,min=2,max=255"`
	Jabatan        string `json:"jabatan" validate:"required,min=2,max=255"`
	NoHP           string `json:"no_hp" validate:"required,min=8,max=50"`
	Sektor         string `json:"sektor" validate:"required,min=2,max=255"`
	TurnstileToken string `json:"cf-turnstile-response" validate:"required"`
}

type EventRegistrationResponse struct {
	ID           int64  `json:"id"`
	EventID      string `json:"event_id"`
	Nama         string `json:"nama"`
	Email        string `json:"email"`
	Perusahaan   string `json:"perusahaan"`
	Jabatan      string `json:"jabatan"`
	NoHP         string `json:"no_hp"`
	Sektor       string `json:"sektor"`
	QRToken      string `json:"qr_token"`
	QRPayload    string `json:"qr_payload"`
	QRCodeBase64 string `json:"qr_code_base64"`
	DownloadURL  string `json:"download_url"`
	CreatedAt    string `json:"created_at"`
}
