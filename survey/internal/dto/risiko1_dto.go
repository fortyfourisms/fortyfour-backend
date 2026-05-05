package dto

import "survey/internal/models"

// STEP 1: ELIGIBILITY
type EligibilityRequest struct {
	RespondenID   int64  `json:"responden_id"`
	RisikoID      *int   `json:"risiko_id,omitempty"`
	PernahTerjadi bool   `json:"pernah_terjadi"`
}

// STEP 2A: ALASAN
type AlasanRequest struct {
	RespondenID int64  `json:"responden_id"`
	RisikoID    *int   `json:"risiko_id,omitempty"`
	Alasan      string `json:"alasan"`
}

// STEP 2B: DAMPAK
type DampakRequest struct {
	RespondenID       int64                 `json:"responden_id"`
	RisikoID          *int                  `json:"risiko_id,omitempty"`
	DampakReputasi    models.ImpactLevel    `json:"dampak_reputasi"`
	DampakOperasional models.ImpactLevel    `json:"dampak_operasional"`
	DampakFinansial   models.ImpactLevel    `json:"dampak_finansial"`
	DampakHukum       models.ImpactLevel    `json:"dampak_hukum"`
	Frekuensi         models.FrequencyLevel `json:"frekuensi"`
}

// STEP 2C: PENGENDALIAN
type PengendalianRequest struct {
	RespondenID           int64  `json:"responden_id"`
	RisikoID              *int   `json:"risiko_id,omitempty"`
	AdaPengendalian       bool   `json:"ada_pengendalian"`
	DeskripsiPengendalian string `json:"deskripsi_pengendalian,omitempty"`
}

// NAVIGATION
type NavigateRequest struct {
	RespondenID int64  `json:"responden_id"`
	Direction   string `json:"direction"`
	CurrentRisk int    `json:"current_risk"`
}

// PROGRESS RESPONSE
type ProgressResponse struct {
	RespondenID    int64   `json:"responden_id"`
	RisikoID       *int    `json:"risiko_id"`
	LangkahSaatIni *string `json:"langkah_saat_ini"`
	Selesai        bool    `json:"selesai"`
}

// CUSTOM RISIKO
// TODO: used by services and handlers for custom risk creation.
type CustomRisikoRequest struct {
	RespondenID int64  `json:"responden_id"`
	NamaRisiko  string `json:"nama_risiko"`
}
