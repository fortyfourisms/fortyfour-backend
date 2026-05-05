package models

import "time"

// BASE 
type Responden struct {
	ID int64 `json:"id" db:"id"`

	UserID       string `json:"user_id" db:"user_id"`
	IdPerusahaan string `json:"id_perusahaan" db:"id_perusahaan"`

	NamaLengkap string `json:"nama_lengkap" db:"nama_lengkap"`
	Jabatan     string `json:"jabatan" db:"jabatan"`
	Email       string `json:"email" db:"email"`
	NoTelepon   string `json:"no_telepon" db:"no_telepon"`

	SertifikatTraining *string `json:"sertifikat_training,omitempty" db:"sertifikat_training"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// DETAIL (JOIN)
type RespondenDetail struct {
	Responden

	NamaPerusahaan *string `json:"nama_perusahaan,omitempty"`
	NamaSubSektor  *string `json:"nama_sub_sektor,omitempty"`
	NamaSektor     *string `json:"nama_sektor,omitempty"`
}