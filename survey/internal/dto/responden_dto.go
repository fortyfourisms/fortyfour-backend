package dto

// CREATE REQUEST
type CreateRespondenRequest struct {
	IdPerusahaan       string  `json:"id_perusahaan"`
	NamaLengkap        string  `json:"nama_lengkap"`
	Jabatan            string  `json:"jabatan"`
	Email              string  `json:"email"`
	NoTelepon          string  `json:"no_telepon"`
	SertifikatTraining *string `json:"sertifikat_training,omitempty"`
}

// UPDATE REQUEST
type UpdateRespondenRequest struct {
	IdPerusahaan       string  `json:"id_perusahaan"`
	NamaLengkap        string  `json:"nama_lengkap"`
	Jabatan            string  `json:"jabatan"`
	Email              string  `json:"email"`
	NoTelepon          string  `json:"no_telepon"`
	SertifikatTraining *string `json:"sertifikat_training,omitempty"`
}

// RESPONSE 
type RespondenResponse struct {
	ID int64 `json:"id"`

	// dari backend (JWT)
	UserID string `json:"user_id"`

	// dari RESPONDEN
	NamaLengkap        string  `json:"nama_lengkap"`
	Jabatan            string  `json:"jabatan"`
	Email              string  `json:"email"`
	NoTelepon          string  `json:"no_telepon"`
	SertifikatTraining *string `json:"sertifikat_training,omitempty"`

	// dari PERUSAHAAN
	IdPerusahaan   string `json:"id_perusahaan"`
	NamaPerusahaan *string `json:"nama_perusahaan"`

	// dari JOIN (opsional)
	NamaSubSektor *string `json:"nama_sub_sektor,omitempty"`
	NamaSektor    *string `json:"nama_sektor,omitempty"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}