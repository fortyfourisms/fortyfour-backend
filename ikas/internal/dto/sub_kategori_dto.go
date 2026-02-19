package dto

import "time"

type CreateSubKategoriRequest struct {
	KategoriID      string `json:"kategori_id"`
	NamaSubKategori string `json:"nama_sub_kategori"`
}

type UpdateSubKategoriRequest struct {
	KategoriID      *string `json:"kategori_id,omitempty"`
	NamaSubKategori *string `json:"nama_sub_kategori,omitempty"`
}

type SubKategoriResponse struct {
	ID              string    `json:"id"`
	KategoriID      string    `json:"kategori_id"`
	NamaSubKategori string    `json:"nama_sub_kategori"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
