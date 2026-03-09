package dto

import "time"

type CreateJawabanGulihRequest struct {
	PertanyaanGulihID int      `json:"pertanyaan_gulih_id"`
	PerusahaanID      string   `json:"perusahaan_id"`
	JawabanGulih      *float64 `json:"jawaban_gulih"`
	Evidence          *string  `json:"evidence,omitempty"`
	Validasi          *string  `json:"validasi,omitempty"`
	Keterangan        *string  `json:"keterangan,omitempty"`
}

type UpdateJawabanGulihRequest struct {
	JawabanGulih *float64 `json:"jawaban_gulih,omitempty"`
	Evidence     *string  `json:"evidence,omitempty"`
	Validasi     *string  `json:"validasi,omitempty"`
	Keterangan   *string  `json:"keterangan,omitempty"`
}

type PertanyaanGulihInfo struct {
	ID              int             `json:"id"`
	PertanyaanGulih string          `json:"pertanyaan_gulih"`
	SubKategori     SubKategoriInfo `json:"sub_kategori"`
}

type JawabanGulihResponse struct {
	ID              int                 `json:"id"`
	PertanyaanGulih PertanyaanGulihInfo `json:"pertanyaan_gulih"`
	PerusahaanID    string              `json:"perusahaan_id"`
	JawabanGulih    *float64            `json:"jawaban_gulih"`
	Evidence        *string             `json:"evidence"`
	Validasi        *string             `json:"validasi"`
	Keterangan      *string             `json:"keterangan"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
}
