package models

import (
	"database/sql"
	"time"
)

// ENUMS
type ImpactLevel int

const (
	ImpactNotSignificant    ImpactLevel = 1
	ImpactFairlySignificant ImpactLevel = 2
	ImpactSignificant       ImpactLevel = 3
	ImpactVerySignificant   ImpactLevel = 4
)

func (i ImpactLevel) Valid() bool {
	return i >= 1 && i <= 4
}

type FrequencyLevel int

const (
	FrequencySmall     FrequencyLevel = 1
	FrequencyMedium    FrequencyLevel = 2
	FrequencyLarge     FrequencyLevel = 3
	FrequencyVeryLarge FrequencyLevel = 4
)

func (f FrequencyLevel) Valid() bool {
	return f >= 1 && f <= 4
}

// MASTER RISIKO
type Risiko struct {
	ID        string    `db:"id"`
	Kode      string    `db:"kode"`
	Nama      string    `db:"nama"`
	Deskripsi *string   `db:"deskripsi"`
	Urutan    int       `db:"urutan"`
	Aktif     bool      `db:"aktif"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// STEP 1 — ELIGIBILITY
type RisikoEligibility struct {
	ID            string    `db:"id"`
	RespondenID   string    `db:"responden_id"`
	RisikoID      *string   `db:"risiko_id"`
	PernahTerjadi bool      `db:"pernah_terjadi"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

// STEP 2A — ALASAN
type RisikoAlasan struct {
	ID          string    `db:"id"`
	RespondenID string    `db:"responden_id"`
	RisikoID    *string   `db:"risiko_id"`
	Alasan      string    `db:"alasan"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// STEP 2B — DAMPAK
type RisikoDampak struct {
	ID          string  `db:"id"`
	RespondenID string  `db:"responden_id"`
	RisikoID    *string `db:"risiko_id"`

	DampakReputasi    string `db:"dampak_reputasi"` // ENUM DB (string)
	DampakOperasional string `db:"dampak_operasional"`
	DampakFinansial   string `db:"dampak_finansial"`
	DampakHukum       string `db:"dampak_hukum"`

	Frekuensi string    `db:"frekuensi"` // ENUM DB
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// STEP 2C — PENGENDALIAN
type RisikoPengendalian struct {
	ID          string  `db:"id"`
	RespondenID string  `db:"responden_id"`
	RisikoID    *string `db:"risiko_id"`

	AdaPengendalian       bool    `db:"ada_pengendalian"`
	DeskripsiPengendalian *string `db:"deskripsi_pengendalian"` // nullable

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// PROGRESS SURVEY
type SurveyProgress struct {
	ID              string         `db:"id"`
	RespondenID     string         `db:"responden_id"`
	RisikoID        sql.NullString `db:"risiko_id"`
	LangkahSaatIni  sql.NullString `db:"langkah_saat_ini"`
	Selesai         bool           `db:"selesai"`
	Status          string         `db:"status"`
	EditReason      sql.NullString `db:"edit_request_reason"`
	EditResponse    sql.NullString `db:"edit_request_response"`
	SubmittedAt     sql.NullTime   `db:"submitted_at"`
	EditRequestedAt sql.NullTime   `db:"edit_requested_at"`
	EditApprovedAt  sql.NullTime   `db:"edit_approved_at"`
	EditApprovedBy  sql.NullString `db:"edit_approved_by"`
	EditRejectedAt  sql.NullTime   `db:"edit_rejected_at"`
	EditRejectedBy  sql.NullString `db:"edit_rejected_by"`
	TerakhirUpdate  time.Time      `db:"terakhir_update"`
}

// EDIT REQUEST ITEM (joined: survey_progress + responden + perusahaan)
type EditRequestItem struct {
	RespondenID     string         `db:"responden_id"`
	UserID          string         `db:"user_id"`
	NamaLengkap     string         `db:"nama_lengkap"`
	NamaPerusahaan  sql.NullString `db:"nama_perusahaan"`
	Status          string         `db:"status"`
	EditReason      sql.NullString `db:"edit_request_reason"`
	EditResponse    sql.NullString `db:"edit_request_response"`
	EditRequestedAt sql.NullTime   `db:"edit_requested_at"`
	EditApprovedAt  sql.NullTime   `db:"edit_approved_at"`
	EditApprovedBy  sql.NullString `db:"edit_approved_by"`
	EditRejectedAt  sql.NullTime   `db:"edit_rejected_at"`
	EditRejectedBy  sql.NullString `db:"edit_rejected_by"`
}

// RESPONSE
type RisikoResponse struct {
	ID         string `json:"id"`
	NamaRisiko string `json:"nama_risiko"`
	Deskripsi  string `json:"deskripsi"`
}

// GLOBAL API RESPONSE
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
