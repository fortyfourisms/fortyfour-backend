package dto

type PerusahaanRequest struct {
	NamaPerusahaan string `json:"nama_perusahaan"`
	JenisUsaha     string `json:"jenis_usaha"`
}

type PerusahaanResponse struct {
	ID             int    `json:"id"`
	NamaPerusahaan string `json:"nama_perusahaan"`
	JenisUsaha     string `json:"jenis_usaha"`
}
