package dto

import "time"

type CreateJawabanIdentifikasiRequest struct {
	PertanyaanIdentifikasiID string  `json:"pertanyaan_identifikasi_id"`
	PerusahaanID             string  `json:"perusahaan_id"`
	JawabanIdentifikasi      float64 `json:"jawaban_identifikasi"`
	Evidence                 *string `json:"evidence,omitempty"`
	Validasi                 *string `json:"validasi,omitempty"`
	Keterangan               *string `json:"keterangan,omitempty"`
}

type UpdateJawabanIdentifikasiRequest struct {
	JawabanIdentifikasi *float64 `json:"jawaban_identifikasi,omitempty"`
	Evidence            *string  `json:"evidence,omitempty"`
	Validasi            *string  `json:"validasi,omitempty"`
	Keterangan          *string  `json:"keterangan,omitempty"`
}

type JawabanIdentifikasiResponse struct {
	ID                       string    `json:"id"`
	PertanyaanIdentifikasiID string    `json:"pertanyaan_identifikasi_id"`
	PerusahaanID             string    `json:"perusahaan_id"`
	JawabanIdentifikasi      float64   `json:"jawaban_identifikasi"`
	Evidence                 *string   `json:"evidence"`
	Validasi                 *string   `json:"validasi"`
	Keterangan               *string   `json:"keterangan"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}
