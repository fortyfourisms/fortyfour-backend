package models

type PICPerusahaan struct {
	ID           string `db:"id" json:"id"`
	Nama         string `db:"nama" json:"nama"`
	Telepon      string `db:"telepon" json:"telepon"`
	IDPerusahaan string `db:"id_perusahaan" json:"id_perusahaan"`
}
