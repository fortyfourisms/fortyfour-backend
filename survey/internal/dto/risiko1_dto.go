package dto

// RESPONSE MASTER RISIKO
type RisikoResponse struct {
	ID         int    `json:"id"`
	Kode       string `json:"kode"`
	Nama       string `json:"nama"`
	Deskripsi  string `json:"deskripsi"`
}

// STEP 1: ELIGIBILITY
type EligibilityRequest struct {
	RespondenID   int  `json:"responden_id"`
	RisikoID      int  `json:"risiko_id"`
	PernahTerjadi bool `json:"pernah_terjadi"`
}

// STEP 2A: ALASAN (JIKA TIDAK)
type AlasanRequest struct {
	RespondenID int    `json:"responden_id"`
	RisikoID    int    `json:"risiko_id"`
	Alasan      string `json:"alasan"`
}

// STEP 2B: DAMPAK (JIKA YA)
type DampakRequest struct {
	RespondenID int `json:"responden_id"`
	RisikoID    int `json:"risiko_id"`

	DampakReputasi    string `json:"dampak_reputasi"`    // ENUM
	DampakOperasional string `json:"dampak_operasional"` // ENUM
	DampakFinansial   string `json:"dampak_finansial"`   // ENUM
	DampakHukum       string `json:"dampak_hukum"`       // ENUM

	Frekuensi string `json:"frekuensi"` // ENUM
}

// STEP 2C: PENGENDALIAN
type PengendalianRequest struct {
	RespondenID int  `json:"responden_id"`
	RisikoID    int  `json:"risiko_id"`

	AdaPengendalian bool   `json:"ada_pengendalian"`
	Deskripsi       string `json:"deskripsi_pengendalian,omitempty"`
}

// RESPONSE GENERIC
type RisikoStepResponse struct {
	NextStep string `json:"next_step"`
	Message  string `json:"message,omitempty"`
}