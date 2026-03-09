package models

import (
	"time"
)

type JawabanGulih struct {
	ID                int       `json:"id"`
	PertanyaanGulihID int       `json:"pertanyaan_gulih_id"`
	PerusahaanID      string    `json:"perusahaan_id"`
	JawabanGulih      *float64  `json:"jawaban_gulih"`
	Evidence          *string   `json:"evidence"`
	Validasi          *string   `json:"validasi"`
	Keterangan        *string   `json:"keterangan"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
