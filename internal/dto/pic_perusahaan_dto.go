package dto

type CreatePICRequest struct {
	Nama         *string `json:"nama,omitempty"`
	Telepon      *string `json:"telepon,omitempty"`
	IDPerusahaan *string `json:"id_perusahaan,omitempty"`
}

type UpdatePICRequest struct {
	Nama         *string `json:"nama,omitempty"`
	Telepon      *string `json:"telepon,omitempty"`
	IDPerusahaan *string `json:"id_perusahaan,omitempty"`
}

type PICResponse struct {
	ID           string `json:"id"`
	Nama         string `json:"nama"`
	Telepon      string `json:"telepon"`
	IDPerusahaan string `json:"id_perusahaan"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
