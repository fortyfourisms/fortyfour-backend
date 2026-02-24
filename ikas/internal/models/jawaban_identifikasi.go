package models

type JawabanIdentifikasi struct {
	ID                       string  `json:"id"`
	PertanyaanIdentifikasiID string  `json:"pertanyaan_identifikasi_id"`
	PerusahaanID             string  `json:"perusahaan_id"`
	JawabanIdentifikasi      string  `json:"jawaban_identifikasi"`
	Evidence                 *string `json:"evidence"`
	Validasi                 *string `json:"validasi"`
	Keterangan               *string `json:"keterangan"`
	CreatedAt                string  `json:"created_at"`
	UpdatedAt                string  `json:"updated_at"`
}
