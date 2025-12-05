package models

type Perusahaan struct {
	ID             string `json:"id"`
	Photo          string `json:"photo"`
	NamaPerusahaan string `json:"nama_perusahaan"`
	JenisUsaha     string `json:"jenis_usaha"`
	Alamat         string `json:"alamat"`
	Telepon        string `json:"telepon"`
	Email          string `json:"email"`
	Website        string `json:"website"`
}
