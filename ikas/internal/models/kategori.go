package models

type Kategori struct {
	ID           string `json:"id"`
	DomainID     string `json:"domain_id"`
	NamaKategori string `json:"nama_kategori"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
