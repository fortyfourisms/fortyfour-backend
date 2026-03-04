package models

type PertanyaanDeteksi struct {
	ID                int     `json:"id"`
	SubKategoriID     string  `json:"sub_kategori_id"`
	RuangLingkupID    string  `json:"ruang_lingkup_id"`
	PertanyaanDeteksi string  `json:"pertanyaan_deteksi"`
	Index0            *string `json:"index0"`
	Index1            *string `json:"index1"`
	Index2            *string `json:"index2"`
	Index3            *string `json:"index3"`
	Index4            *string `json:"index4"`
	Index5            *string `json:"index5"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}
