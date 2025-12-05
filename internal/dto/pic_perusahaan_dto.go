package dto

type CreatePICPerusahaanRequest struct {
	Nama         string `json:"nama"`
	Telepon      string `json:"telepon"`
	IDPerusahaan string `json:"id_perusahaan"`
}

type UpdatePICPerusahaanRequest struct {
	Nama         string `json:"nama"`
	Telepon      string `json:"telepon"`
	IDPerusahaan string `json:"id_perusahaan"`
}

type PICPerusahaanResponse struct {
	ID           string `json:"id"`
	Nama         string `json:"nama"`
	Telepon      string `json:"telepon"`
	IDPerusahaan string `json:"id_perusahaan"`
	NamaCompany  string `json:"nama_perusahaan"`
}
