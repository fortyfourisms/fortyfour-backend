package dto

import "time"

type CreatePertanyaanGulihRequest struct {
	SubKategoriID   int     `json:"sub_kategori_id"`
	RuangLingkupID  int     `json:"ruang_lingkup_id"`
	PertanyaanGulih string  `json:"pertanyaan_gulih"`
	Index0          *string `json:"index0,omitempty"`
	Index1          *string `json:"index1,omitempty"`
	Index2          *string `json:"index2,omitempty"`
	Index3          *string `json:"index3,omitempty"`
	Index4          *string `json:"index4,omitempty"`
	Index5          *string `json:"index5,omitempty"`
}

type UpdatePertanyaanGulihRequest struct {
	SubKategoriID   *int    `json:"sub_kategori_id,omitempty"`
	RuangLingkupID  *int    `json:"ruang_lingkup_id,omitempty"`
	PertanyaanGulih *string `json:"pertanyaan_gulih,omitempty"`
	Index0          *string `json:"index0,omitempty"`
	Index1          *string `json:"index1,omitempty"`
	Index2          *string `json:"index2,omitempty"`
	Index3          *string `json:"index3,omitempty"`
	Index4          *string `json:"index4,omitempty"`
	Index5          *string `json:"index5,omitempty"`
}

type PertanyaanGulihResponse struct {
	ID              int              `json:"id"`
	SubKategori     SubKategoriInfo  `json:"sub_kategori"`
	RuangLingkup    RuangLingkupInfo `json:"ruang_lingkup"`
	PertanyaanGulih string           `json:"pertanyaan_gulih"`
	Index0          *string          `json:"index0"`
	Index1          *string          `json:"index1"`
	Index2          *string          `json:"index2"`
	Index3          *string          `json:"index3"`
	Index4          *string          `json:"index4"`
	Index5          *string          `json:"index5"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}
