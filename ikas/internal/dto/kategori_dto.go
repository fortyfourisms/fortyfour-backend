package dto

import "time"

type CreateKategoriRequest struct {
	DomainID     int    `json:"domain_id"`
	NamaKategori string `json:"nama_kategori"`
}

type UpdateKategoriRequest struct {
	DomainID     *int    `json:"domain_id,omitempty"`
	NamaKategori *string `json:"nama_kategori,omitempty"`
}

type KategoriResponse struct {
	ID           int       `json:"id"`
	DomainID     int       `json:"domain_id"`
	NamaKategori string    `json:"nama_kategori"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type KategoriMessageResponse struct {
	ID      int    `json:"id,omitempty"`
	Message string `json:"message"`
}
