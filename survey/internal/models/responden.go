package models

import "time"

type Responden struct {
	ID                 int       `json:"id"`
	UserID             string    `json:"user_id"`
	NoTelepon          string    `json:"no_telepon"`
	Sektor             string    `json:"sektor"`
	SektorLainnya      *string   `json:"sektor_lainnya,omitempty"`
	SertifikatTraining string    `json:"sertifikat_training"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type RespondenDetail struct {
	ID     int    `json:"id"`
	UserID string `json:"user_id"`

	// dari users
	NamaLengkap *string `json:"nama_lengkap,omitempty"`
	Email       *string `json:"email,omitempty"`

	// dari jabatan
	Jabatan *string `json:"jabatan,omitempty"`

	// dari perusahaan
	NamaPerusahaan *string `json:"nama_perusahaan,omitempty"`
	PerusahaanID   *string `json:"perusahaan_id,omitempty"`

	// dari responden
	NoTelepon          string  `json:"no_telepon"`
	Sektor             string  `json:"sektor"`
	SektorLainnya      *string `json:"sektor_lainnya,omitempty"`
	SertifikatTraining string  `json:"sertifikat_training"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
