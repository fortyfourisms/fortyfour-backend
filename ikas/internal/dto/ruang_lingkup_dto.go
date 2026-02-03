package dto

import "time"

type CreateRuangLingkupRequest struct {
	NamaRuangLingkup string `json:"nama_ruang_lingkup"`
}

type UpdateRuangLingkupRequest struct {
	NamaRuangLingkup *string `json:"nama_ruang_lingkup,omitempty"`
}

type RuangLingkupResponse struct {
	ID               string    `json:"id"`
	NamaRuangLingkup string    `json:"nama_ruang_lingkup"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
