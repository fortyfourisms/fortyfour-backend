package dto

type CreateRuangLingkupRequest struct {
	NamaRuangLingkup string `json:"nama_ruang_lingkup"`
}

type UpdateRuangLingkupRequest struct {
	NamaRuangLingkup *string `json:"nama_ruang_lingkup,omitempty"`
}

type RuangLingkupResponse struct {
	ID               string `json:"id"`
	NamaRuangLingkup string `json:"nama_ruang_lingkup"`
}
