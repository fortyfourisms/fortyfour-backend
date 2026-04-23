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
	ID        int       `db:"id"`
	Kode      string    `db:"kode"`
	Nama      string    `db:"nama"`
	Deskripsi string    `db:"deskripsi"`
	Urutan    int       `db:"urutan"`
	Aktif     bool      `db:"aktif"`
	CreatedAt time.Time `db:"created_at"`
}

// CUSTOM RISIKO (RISIKO 14)
type CustomRisiko struct {
	ID          int       `db:"id"`
	RespondenID int       `db:"responden_id"`
	NamaRisiko  string    `db:"nama_risiko"`
	CreatedAt   time.Time `db:"created_at"`
}

// STEP 1 — ELIGIBILITY
type RisikoEligibility struct {
	ID             int       `db:"id"`
	RespondenID    int       `db:"responden_id"`
	RisikoID       int       `db:"risiko_id"`
	CustomRisikoID *int      `db:"custom_risiko_id"`
	PernahTerjadi  bool      `db:"pernah_terjadi"`
	CreatedAt      time.Time `db:"created_at"`
}

// STEP 2A — ALASAN
type RisikoAlasan struct {
	ID             int       `db:"id"`
	RespondenID    int       `db:"responden_id"`
	RisikoID       int       `db:"risiko_id"`
	CustomRisikoID *int      `db:"custom_risiko_id"`
	Alasan         string    `db:"alasan"`
	CreatedAt      time.Time `db:"created_at"`
}

// STEP 2B — DAMPAK
type RisikoDampak struct {
	ID                  int            `db:"id"`
	RespondenID         int            `db:"responden_id"`
	RisikoID            int            `db:"risiko_id"`
	CustomRisikoID      *int           `db:"custom_risiko_id"`

	DampakReputasi      ImpactLevel    `db:"dampak_reputasi"`
	DampakOperasional   ImpactLevel    `db:"dampak_operasional"`
	DampakFinansial     ImpactLevel    `db:"dampak_finansial"`
	DampakHukum         ImpactLevel    `db:"dampak_hukum"`

	Frekuensi           FrequencyLevel `db:"frekuensi"`
	CreatedAt           time.Time      `db:"created_at"`
}

// STEP 2C — PENGENDALIAN
type RisikoPengendalian struct {
	ID                     int       `db:"id"`
	RespondenID            int       `db:"responden_id"`
	RisikoID               int       `db:"risiko_id"`
	CustomRisikoID         *int      `db:"custom_risiko_id"`

	AdaPengendalian        bool      `db:"ada_pengendalian"`
	DeskripsiPengendalian  string    `db:"deskripsi_pengendalian"`

	CreatedAt              time.Time `db:"created_at"`
}

// PROGRESS SURVEY
type SurveyProgress struct {
	ID             int            `db:"id"`
	RespondenID    int            `db:"responden_id"`
	RisikoID       sql.NullInt64  `db:"risiko_id"`
	LangkahSaatIni sql.NullString `db:"langkah_saat_ini"`
	Selesai        bool           `db:"selesai"`
	TerakhirUpdate time.Time      `db:"terakhir_update"`
}

// RISIKO RESPONSE 
type RisikoResponse struct {
	ID         int    `json:"id"`
	NamaRisiko string `json:"nama_risiko"`
	Deskripsi  string `json:"deskripsi"`
} 

// GLOBAL API RESPONSE
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}