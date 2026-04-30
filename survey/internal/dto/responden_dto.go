package dto

type CreateRespondenRequest struct {
	IdPerusahaan       string `json:"id_perusahaan"`
	NamaLengkap        string `json:"nama_lengkap"`
	Jabatan            string `json:"jabatan"`
	Email              string `json:"email"`
	NoTelepon          string `json:"no_telepon"`
	SertifikatTraining string `json:"sertifikat_training"`
}

type UpdateRespondenRequest struct {
	IdPerusahaan       string `json:"id_perusahaan"`
	NamaLengkap        string `json:"nama_lengkap"`
	Jabatan            string `json:"jabatan"`
	Email              string `json:"email"`
	NoTelepon          string `json:"no_telepon"`
	SertifikatTraining string `json:"sertifikat_training"`
}

type RespondenResponse struct {
	ID int `json:"id"`

	// dari RESPONDEN
	NamaLengkap        string `json:"nama_lengkap"`
	Jabatan            string `json:"jabatan"`
	Email              string `json:"email"`
	NoTelepon          string `json:"no_telepon"`
	SertifikatTraining string `json:"sertifikat_training"`

	// dari PERUSAHAAN
	IdPerusahaan   string `json:"id_perusahaan"`
	NamaPerusahaan string `json:"nama_perusahaan"`

	// dari SUB SEKTOR & SEKTOR (hasil JOIN)
	NamaSubSektor string `json:"nama_sub_sektor,omitempty"`
	NamaSektor    string `json:"nama_sektor,omitempty"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
