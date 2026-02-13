package models

type SubSektor struct {
	ID            string `json:"id"`
	NamaSubSektor string `json:"nama_sub_sektor"`
	IDSektor      string `json:"id_sektor"`
	NamaSektor    string `json:"nama_sektor,omitempty"` // Untuk join dengan tabel sektor
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}