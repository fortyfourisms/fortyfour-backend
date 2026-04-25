package dto

type CreateRespondenRequest struct {
	UserID             string `json:"user_id"` 
	NoTelepon          string `json:"no_telepon"`
	Sektor             string `json:"sektor"`
	SektorLainnya      string `json:"sektor_lainnya,omitempty"`
	SertifikatTraining string `json:"sertifikat_training"`
}

type UpdateRespondenRequest struct {
	UserID             string `json:"user_id"`
	NoTelepon          string `json:"no_telepon"`
	Sektor             string `json:"sektor"`
	SektorLainnya      string `json:"sektor_lainnya,omitempty"`
	SertifikatTraining string `json:"sertifikat_training"`
}

type RespondenResponse struct {
	ID int `json:"id"`

	// dari USERS
	UserID 		string `json:"user_id"`
	NamaLengkap string `json:"nama_lengkap"`
	Email       string `json:"email"`

	// dari JABATAN
	Jabatan string `json:"jabatan"`

	// dari PERUSAHAAN
	PerusahaanID   string `json:"perusahaan_id,omitempty"`
	NamaPerusahaan string `json:"nama_perusahaan,omitempty"`

	// dari RESPONDEN
	NoTelepon          string `json:"no_telepon"`
	Sektor             string `json:"sektor"`
	SektorLainnya      string `json:"sektor_lainnya,omitempty"`
	SertifikatTraining string `json:"sertifikat_training"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}