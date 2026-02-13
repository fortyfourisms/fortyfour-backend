package dto

type SektorResponse struct {
	ID         string              `json:"id"`
	NamaSektor string              `json:"nama_sektor"`
	SubSektor  []SubSektorResponse `json:"sub_sektor,omitempty"`
	CreatedAt  string              `json:"created_at"`
	UpdatedAt  string              `json:"updated_at"`
}

type SubSektorResponse struct {
	ID            string `json:"id"`
	NamaSubSektor string `json:"nama_sub_sektor"`
	IDSektor      string `json:"id_sektor"`
	NamaSektor    string `json:"nama_sektor,omitempty"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}