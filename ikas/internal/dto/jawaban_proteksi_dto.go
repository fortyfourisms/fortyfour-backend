package dto

import "time"

type CreateJawabanProteksiRequest struct {
	PertanyaanProteksiID int      `json:"pertanyaan_proteksi_id"`
	PerusahaanID         string   `json:"perusahaan_id"`
	JawabanProteksi      *float64 `json:"jawaban_proteksi"`
	Evidence             *string  `json:"evidence,omitempty"`
	Validasi             *string  `json:"validasi,omitempty"`
	Keterangan           *string  `json:"keterangan,omitempty"`
}

type UpdateJawabanProteksiRequest struct {
	JawabanProteksi *float64 `json:"jawaban_proteksi,omitempty"`
	Evidence        *string  `json:"evidence,omitempty"`
	Validasi        *string  `json:"validasi,omitempty"`
	Keterangan      *string  `json:"keterangan,omitempty"`
}

type PertanyaanProteksiInfo struct {
	ID                 int             `json:"id"`
	PertanyaanProteksi string          `json:"pertanyaan_proteksi"`
	SubKategori        SubKategoriInfo `json:"sub_kategori"`
}

type JawabanProteksiResponse struct {
	ID                 int                    `json:"id"`
	PertanyaanProteksi PertanyaanProteksiInfo `json:"pertanyaan_proteksi"`
	PerusahaanID       string                 `json:"perusahaan_id"`
	JawabanProteksi    *float64               `json:"jawaban_proteksi"`
	Evidence           *string                `json:"evidence"`
	Validasi           *string                `json:"validasi"`
	Keterangan         *string                `json:"keterangan"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
}
