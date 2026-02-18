package models

type Perusahaan struct {
	ID             string  `json:"id"`
	Photo          string  `json:"photo"`
	NamaPerusahaan string  `json:"nama_perusahaan"`
	IDSubSektor    *string `json:"id_sub_sektor"`
	Alamat         string  `json:"alamat"`
	Telepon        string  `json:"telepon"`
	Email          string  `json:"email"`
	Website        string  `json:"website"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}
