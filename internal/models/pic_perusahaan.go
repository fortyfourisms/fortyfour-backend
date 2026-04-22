package models

type PICPerusahaan struct {
	ID           string `json:"id"`
	Nama         string `json:"nama"`
	Telepon      string `json:"telepon"`
	Email        string `json:"email"`
	IDPerusahaan string `json:"id_perusahaan"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
