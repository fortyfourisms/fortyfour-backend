package models

type SeCsirt struct {
	ID          string `json:"id"`
	IdCsirt     string `json:"id_csirt"`
	NamaSe      string `json:"nama_se"`
	IpSe        string `json:"ip_se"`
	AsNumberSe  string `json:"as_number_se"`
	PengelolaSe string `json:"pengelola_se"`
	FiturSe     string `json:"fitur_se"`
	KategoriSe  string `json:"kategori_se"` // rendah, tinggi, strategis
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
