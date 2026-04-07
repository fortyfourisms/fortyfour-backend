package models

import "time"

type Responden struct {
	ID                 int       `json:"id"`
	NamaLengkap        string    `json:"nama_lengkap"`
	Jabatan            string    `json:"jabatan"`
	Perusahaan         string    `json:"perusahaan"`
	Email              string    `json:"email"`
	NoTelepon          string    `json:"no_telepon"`
	Sektor             string    `json:"sektor"`
	SektorLainnya      string    `json:"sektor_lainnya"`
	SertifikatTraining string    `json:"sertifikat_training"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
