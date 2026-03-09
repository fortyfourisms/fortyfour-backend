package models

import (
	"time"
)

type JawabanDeteksi struct {
	ID                  int       `json:"id"`
	PertanyaanDeteksiID int       `json:"pertanyaan_deteksi_id"`
	PerusahaanID        string    `json:"perusahaan_id"`
	JawabanDeteksi      *float64  `json:"jawaban_deteksi"`
	Evidence            *string   `json:"evidence"`
	Validasi            *string   `json:"validasi"`
	Keterangan          *string   `json:"keterangan"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
