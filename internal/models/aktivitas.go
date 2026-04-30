package models

type Aktivitas struct {
	ID             int    `json:"id"`
	PerusahaanID   string `json:"perusahaan_id"`
	Judul          string `json:"judul"`
	Deskripsi      string `json:"deskripsi"`
	TanggalMulai   string `json:"tanggal_mulai"`
	TanggalSelesai string `json:"tanggal_selesai"`
	JenisAktivitas string `json:"jenis_aktivitas"` // JSON string
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
