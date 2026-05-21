package dto

import "survey/internal/models"

// STEP 1: ELIGIBILITY
type EligibilityRequest struct {
	RespondenID   string  `json:"responden_id"`
	RisikoID      *string `json:"risiko_id,omitempty"`
	PernahTerjadi bool    `json:"pernah_terjadi"`
}

// STEP 2A: ALASAN
type AlasanRequest struct {
	RespondenID string  `json:"responden_id"`
	RisikoID    *string `json:"risiko_id,omitempty"`
	Alasan      string  `json:"alasan"`
}

// STEP 2B: DAMPAK
type DampakRequest struct {
	RespondenID       string                `json:"responden_id"`
	RisikoID          *string               `json:"risiko_id,omitempty"`
	DampakReputasi    models.ImpactLevel    `json:"dampak_reputasi"`
	DampakOperasional models.ImpactLevel    `json:"dampak_operasional"`
	DampakFinansial   models.ImpactLevel    `json:"dampak_finansial"`
	DampakHukum       models.ImpactLevel    `json:"dampak_hukum"`
	Frekuensi         models.FrequencyLevel `json:"frekuensi"`
}

// STEP 2C: PENGENDALIAN
type PengendalianRequest struct {
	RespondenID           string  `json:"responden_id"`
	RisikoID              *string `json:"risiko_id,omitempty"`
	AdaPengendalian       bool    `json:"ada_pengendalian"`
	DeskripsiPengendalian string  `json:"deskripsi_pengendalian,omitempty"`
}

// NAVIGATION
type NavigateRequest struct {
	RespondenID string `json:"responden_id"`
	Direction   string `json:"direction"`
	CurrentRisk int    `json:"current_risk"`
}

// PROGRESS RESPONSE
type ProgressResponse struct {
	RespondenID     string  `json:"responden_id"`
	RisikoID        *string `json:"risiko_id"`
	LangkahSaatIni  *string `json:"langkah_saat_ini"`
	Selesai         bool    `json:"selesai"`
	Status          string  `json:"status"`
	IsRejected      bool    `json:"is_rejected"`
	EditReason      *string `json:"edit_request_reason,omitempty"`
	EditResponse    *string `json:"edit_request_response,omitempty"`
	SubmittedAt     *string `json:"submitted_at,omitempty"`
	EditRequestedAt *string `json:"edit_requested_at,omitempty"`
	EditApprovedAt  *string `json:"edit_approved_at,omitempty"`
	EditApprovedBy  *string `json:"edit_approved_by,omitempty"`
	EditRejectedAt  *string `json:"edit_rejected_at,omitempty"`
	EditRejectedBy  *string `json:"edit_rejected_by,omitempty"`
}

// CUSTOM RISIKO
// TODO: used by services and handlers for custom risk creation.
type CustomRisikoRequest struct {
	RespondenID string `json:"responden_id"`
	NamaRisiko  string `json:"nama_risiko"`
}

type RequestEditRequest struct {
	Reason string `json:"reason"`
}

type ReviewEditRequest struct {
	Action   string `json:"action" example:"approve"` // "approve" or "reject"
	Response string `json:"response,omitempty" example:"Alasan persetujuan atau penolakan"`
}

// EDIT REQUEST LIST RESPONSE
type EditRequestItemResponse struct {
	RespondenID     string  `json:"responden_id"`
	UserID          string  `json:"user_id"`
	NamaLengkap     string  `json:"nama_lengkap"`
	NamaPerusahaan  *string `json:"nama_perusahaan,omitempty"`
	Status          string  `json:"status"`
	EditReason      *string `json:"edit_request_reason,omitempty"`
	EditResponse    *string `json:"edit_request_response,omitempty"`
	EditRequestedAt *string `json:"edit_requested_at,omitempty"`
	EditApprovedAt  *string `json:"edit_approved_at,omitempty"`
	EditApprovedBy  *string `json:"edit_approved_by,omitempty"`
	EditRejectedAt  *string `json:"edit_rejected_at,omitempty"`
	EditRejectedBy  *string `json:"edit_rejected_by,omitempty"`
}
