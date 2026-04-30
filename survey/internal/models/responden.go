package models

import "time"

type Responden struct {
	ID                 int       `json:"id"`
	IdPerusahaan       string    `json:"id_perusahaan"`
	NamaLengkap        string    `json:"nama_lengkap"`
	Jabatan            string    `json:"jabatan"`
	Email              string    `json:"email"`
	NoTelepon          string    `json:"no_telepon"`
	SertifikatTraining string    `json:"sertifikat_training"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type RespondenDetail struct {
	ID int `json:"id"`

	// dari responden
	IdPerusahaan       string    `json:"id_perusahaan"`
	NamaLengkap        string    `json:"nama_lengkap"`
	Jabatan            string    `json:"jabatan"`
	Email              string    `json:"email"`
	NoTelepon          string    `json:"no_telepon"`
	SertifikatTraining string    `json:"sertifikat_training"`

	// dari perusahaan
	NamaPerusahaan *string `json:"nama_perusahaan,omitempty"`

	// dari sub sektor & sektor
	NamaSubSektor *string `json:"nama_sub_sektor,omitempty"`
	NamaSektor    *string `json:"nama_sektor,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}