package models

type SubKategori struct {
	ID              int    `json:"id"`
	KategoriID      int    `json:"kategori_id"`
	NamaSubKategori string `json:"nama_sub_kategori"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}
