package models

type Kategori struct {
	ID           int    `json:"id"`
	DomainID     int    `json:"domain_id"`
	NamaKategori string `json:"nama_kategori"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
