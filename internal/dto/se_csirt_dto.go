package dto

type CreateSeCsirtRequest struct {
	IdCsirt     *string `json:"id_csirt,omitempty"`
	NamaSe      *string `json:"nama_se,omitempty"`
	IpSe        *string `json:"ip_se,omitempty"`
	AsNumberSe  *string `json:"as_number_se,omitempty"`
	PengelolaSe *string `json:"pengelola_se,omitempty"`
	FiturSe     *string `json:"fitur_se,omitempty"`
	KategoriSe  *string `json:"kategori_se,omitempty"`
}

type UpdateSeCsirtRequest struct {
	NamaSe      *string `json:"nama_se,omitempty"`
	IpSe        *string `json:"ip_se,omitempty"`
	AsNumberSe  *string `json:"as_number_se,omitempty"`
	PengelolaSe *string `json:"pengelola_se,omitempty"`
	FiturSe     *string `json:"fitur_se,omitempty"`
	KategoriSe  *string `json:"kategori_se,omitempty"`
}

type SeCsirtResponse struct {
	ID          string `json:"id"`
	IdCsirt     string `json:"id_csirt"`
	NamaSe      string `json:"nama_se"`
	IpSe        string `json:"ip_se"`
	AsNumberSe  string `json:"as_number_se"`
	PengelolaSe string `json:"pengelola_se"`
	FiturSe     string `json:"fitur_se"`
	KategoriSe  string `json:"kategori_se"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
