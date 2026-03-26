package dto

import "time"

type CreateRuangLingkupRequest struct {
	NamaRuangLingkup string `json:"nama_ruang_lingkup"`
}

type UpdateRuangLingkupRequest struct {
	NamaRuangLingkup *string `json:"nama_ruang_lingkup,omitempty"`
}

type RuangLingkupResponse struct {
	ID               int       `json:"id"`
	NamaRuangLingkup string    `json:"nama_ruang_lingkup"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type RuangLingkupMessageResponse struct {
	ID      int    `json:"id,omitempty"`
	Message string `json:"message"`
}
