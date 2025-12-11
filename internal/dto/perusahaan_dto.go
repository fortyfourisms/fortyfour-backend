package dto

type CreatePerusahaanRequest struct {
	Photo          *string `json:"photo,omitempty"`
	NamaPerusahaan *string `json:"nama_perusahaan,omitempty"`
	Sektor         *string `json:"sektor,omitempty"`
	Alamat         *string `json:"alamat,omitempty"`
	Telepon        *string `json:"telepon,omitempty"`
	Email          *string `json:"email,omitempty"`
	Website        *string `json:"website,omitempty"`
}

type UpdatePerusahaanRequest struct {
	Photo          *string `json:"photo,omitempty"`
	NamaPerusahaan *string `json:"nama_perusahaan,omitempty"`
	Sektor         *string `json:"sektor,omitempty"`
	Alamat         *string `json:"alamat,omitempty"`
	Telepon        *string `json:"telepon,omitempty"`
	Email          *string `json:"email,omitempty"`
	Website        *string `json:"website,omitempty"`
}

type PerusahaanResponse struct {
	ID             string `json:"id"`
	Photo          string `json:"photo"`
	NamaPerusahaan string `json:"nama_perusahaan"`
	Sektor         string `json:"sektor"`
	Alamat         string `json:"alamat"`
	Telepon        string `json:"telepon"`
	Email          string `json:"email"`
	Website        string `json:"website"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
