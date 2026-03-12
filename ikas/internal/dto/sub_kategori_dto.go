package dto

import "time"

type CreateSubKategoriRequest struct {
	KategoriID      int    `json:"kategori_id"`
	NamaSubKategori string `json:"nama_sub_kategori"`
}

type UpdateSubKategoriRequest struct {
	KategoriID      *int    `json:"kategori_id,omitempty"`
	NamaSubKategori *string `json:"nama_sub_kategori,omitempty"`
}

type SubKategoriResponse struct {
	ID              int       `json:"id"`
	KategoriID      int       `json:"kategori_id"`
	NamaSubKategori string    `json:"nama_sub_kategori"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
