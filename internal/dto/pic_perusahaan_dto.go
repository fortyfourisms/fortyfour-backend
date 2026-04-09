package dto

type CreatePICRequest struct {
	Nama         *string `json:"nama,omitempty"`
	Telepon      *string `json:"telepon,omitempty"`
	Email        *string `json:"email,omitempty"`
	IDPerusahaan *string `json:"id_perusahaan,omitempty"`
}

type UpdatePICRequest struct {
	Nama         *string `json:"nama,omitempty"`
	Telepon      *string `json:"telepon,omitempty"`
	Email        *string `json:"email,omitempty"`
	IDPerusahaan *string `json:"id_perusahaan,omitempty"`
}

type PICResponse struct {
	ID         string           `json:"id"`
	Nama       string           `json:"nama"`
	Telepon    string           `json:"telepon"`
	Email      string           `json:"email"`
	CreatedAt  string           `json:"created_at"`
	UpdatedAt  string           `json:"updated_at"`
	Perusahaan *PerusahaanInPIC `json:"perusahaan,omitempty"`
}

// Struct baru untuk data perusahaan di dalam PIC
type PerusahaanInPIC struct {
	ID             string `json:"id"`
	NamaPerusahaan string `json:"nama_perusahaan"`
	// Tambahkan field lain jika diperlukan
}