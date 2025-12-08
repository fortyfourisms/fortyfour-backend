package dto

type PerusahaanRequest struct {
	Photo          *string `json:"photo,omitempty"`
	NamaPerusahaan *string `json:"nama_perusahaan,omitempty"`
	JenisUsaha     *string `json:"jenis_usaha,omitempty"`
	Alamat         *string `json:"alamat,omitempty"`
	Telepon        *string `json:"telepon,omitempty"`
	Email          *string `json:"email,omitempty"`
	Website        *string `json:"website,omitempty"`
}

type PerusahaanResponse struct {
	ID             string `json:"id"`
	Photo          string `json:"photo"`
	NamaPerusahaan string `json:"nama_perusahaan"`
	JenisUsaha     string `json:"jenis_usaha"`
	Alamat         string `json:"alamat"`
	Telepon        string `json:"telepon"`
	Email          string `json:"email"`
	Website        string `json:"website"`
}
