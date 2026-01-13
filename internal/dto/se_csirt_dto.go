package dto

// Create request
type CreateSeCsirtRequest struct {
	IdCsirt     *string `json:"id_csirt,omitempty"`
	NamaSe      *string `json:"nama_se,omitempty"`
	IpSe        *string `json:"ip_se,omitempty"`
	AsNumberSe  *string `json:"as_number_se,omitempty"`
	PengelolaSe *string `json:"pengelola_se,omitempty"`
	FiturSe     *string `json:"fitur_se,omitempty"`
	KategoriSe  *string `json:"kategori_se,omitempty"`
}

// Update request
type UpdateSeCsirtRequest struct {
	NamaSe      *string `json:"nama_se,omitempty"`
	IpSe        *string `json:"ip_se,omitempty"`
	AsNumberSe  *string `json:"as_number_se,omitempty"`
	PengelolaSe *string `json:"pengelola_se,omitempty"`
	FiturSe     *string `json:"fitur_se,omitempty"`
	KategoriSe  *string `json:"kategori_se,omitempty"`
}

// Response dengan nested CSIRT
type SeCsirtResponse struct {
	ID         string             `json:"id"`
	NamaSe     string             `json:"nama_se"`
	IpSe       string             `json:"ip_se"`
	AsNumberSe string             `json:"as_number_se"`
	Pengelola  string             `json:"pengelola_se"`
	FiturSe    string             `json:"fitur_se"`
	KategoriSe string             `json:"kategori_se"`
	Csirt      *CsirtMiniResponse `json:"csirt,omitempty"`
	CreatedAt  string             `json:"created_at"`
	UpdatedAt  string             `json:"updated_at"`
}
