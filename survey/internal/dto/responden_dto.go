package dto

type CreateRespondenRequest struct {
	NamaLengkap        string `json:"nama_lengkap"`
	Jabatan            string `json:"jabatan"`
	Perusahaan         string `json:"perusahaan"`
	Email              string `json:"email"`
	NoTelepon          string `json:"no_telepon"`
	Sektor             string `json:"sektor"`
	SektorLainnya      string `json:"sektor_lainnya"`
	SertifikatTraining string `json:"sertifikat_training"`
}

type UpdateRespondenRequest struct {
	NamaLengkap        string `json:"nama_lengkap"`
	Jabatan            string `json:"jabatan"`
	Perusahaan         string `json:"perusahaan"`
	Email              string `json:"email"`
	NoTelepon          string `json:"no_telepon"`
	Sektor             string `json:"sektor"`
	SektorLainnya      string `json:"sektor_lainnya"`
	SertifikatTraining string `json:"sertifikat_training"`
}

type RespondenResponse struct {
	ID                 int    `json:"id"`
	NamaLengkap        string `json:"nama_lengkap"`
	Jabatan            string `json:"jabatan"`
	Perusahaan         string `json:"perusahaan"`
	Email              string `json:"email"`
	NoTelepon          string `json:"no_telepon"`
	Sektor             string `json:"sektor"`
	SektorLainnya      string `json:"sektor_lainnya"`
	SertifikatTraining string `json:"sertifikat_training"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}
