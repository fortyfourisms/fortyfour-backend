package dto

// Nested structs untuk relasi
type DomainInfo struct {
	ID         int    `json:"id"`
	NamaDomain string `json:"nama_domain"`
}

type KategoriInfo struct {
	ID           int        `json:"id"`
	NamaKategori string     `json:"nama_kategori"`
	Domain       DomainInfo `json:"domain"`
}

type SubKategoriInfo struct {
	ID              int          `json:"id"`
	NamaSubKategori string       `json:"nama_sub_kategori"`
	Kategori        KategoriInfo `json:"kategori"`
}

type RuangLingkupInfo struct {
	ID               int    `json:"id"`
	NamaRuangLingkup string `json:"nama_ruang_lingkup"`
}
