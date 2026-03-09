package models

type JawabanProteksi struct {
	ID                   int      `json:"id"`
	PertanyaanProteksiID int      `json:"pertanyaan_proteksi_id"`
	PerusahaanID         string   `json:"perusahaan_id"`
	JawabanProteksi      *float64 `json:"jawaban_proteksi"`
	Evidence             *string  `json:"evidence"`
	Validasi             *string  `json:"validasi"`
	Keterangan           *string  `json:"keterangan"`
	CreatedAt            string   `json:"created_at"`
	UpdatedAt            string   `json:"updated_at"`
}
