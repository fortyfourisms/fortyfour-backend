package dto

import "time"

type CreateKategoriRequest struct {
	DomainID     string `json:"domain_id"`
	NamaKategori string `json:"nama_kategori"`
}

type UpdateKategoriRequest struct {
	DomainID     *string `json:"domain_id,omitempty"`
	NamaKategori *string `json:"nama_kategori,omitempty"`
}

type KategoriResponse struct {
	ID           string    `json:"id"`
	DomainID     string    `json:"domain_id"`
	NamaKategori string    `json:"nama_kategori"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
