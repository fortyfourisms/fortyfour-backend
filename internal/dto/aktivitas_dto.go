package dto

import "time"

type CreateAktivitasRequest struct {
	PerusahaanID   string   `json:"perusahaan_id"`
	Judul          string   `json:"judul"`
	Deskripsi      string   `json:"deskripsi"`
	TanggalMulai   string   `json:"tanggal_mulai"`
	TanggalSelesai string   `json:"tanggal_selesai"`
	JenisAktivitas []string `json:"jenis_aktivitas"`
}

type UpdateAktivitasRequest struct {
	PerusahaanID   *string   `json:"perusahaan_id,omitempty"`
	Judul          *string   `json:"judul,omitempty"`
	Deskripsi      *string   `json:"deskripsi,omitempty"`
	TanggalMulai   *string   `json:"tanggal_mulai,omitempty"`
	TanggalSelesai *string   `json:"tanggal_selesai,omitempty"`
	JenisAktivitas *[]string `json:"jenis_aktivitas,omitempty"`
}

type AktivitasResponse struct {
	ID             int       `json:"id"`
	PerusahaanID   string    `json:"perusahaan_id"`
	Judul          string    `json:"judul"`
	Deskripsi      string    `json:"deskripsi"`
	TanggalMulai   string    `json:"tanggal_mulai"`
	TanggalSelesai string    `json:"tanggal_selesai"`
	JenisAktivitas []string  `json:"jenis_aktivitas"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
