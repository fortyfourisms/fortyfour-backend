package dto

import "time"

type CreateJawabanDeteksiRequest struct {
	PertanyaanDeteksiID int      `json:"pertanyaan_deteksi_id"`
	IkasID              string   `json:"ikas_id"`
	JawabanDeteksi      *float64 `json:"jawaban_deteksi"`
	Evidence            *string  `json:"evidence,omitempty"`
	Validasi            *string  `json:"validasi,omitempty"`
	Keterangan          *string  `json:"keterangan,omitempty"`
}

type UpdateJawabanDeteksiRequest struct {
	JawabanDeteksi *float64 `json:"jawaban_deteksi,omitempty"`
	Evidence       *string  `json:"evidence,omitempty"`
	Validasi       *string  `json:"validasi,omitempty"`
	Keterangan     *string  `json:"keterangan,omitempty"`
}

type PertanyaanDeteksiInfo struct {
	ID                int             `json:"id"`
	PertanyaanDeteksi string          `json:"pertanyaan_deteksi"`
	SubKategori       SubKategoriInfo `json:"sub_kategori"`
}

type JawabanDeteksiResponse struct {
	ID                int                   `json:"id"`
	PertanyaanDeteksi PertanyaanDeteksiInfo `json:"pertanyaan_deteksi"`
	IkasID            string                `json:"ikas_id"`
	JawabanDeteksi    *float64              `json:"jawaban_deteksi"`
	Evidence          *string               `json:"evidence"`
	Validasi          *string               `json:"validasi"`
	Keterangan        *string               `json:"keterangan"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
}
