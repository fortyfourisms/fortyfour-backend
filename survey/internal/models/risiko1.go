package models

import "time"

// ENUMS
type ImpactLevel int

const (
	ImpactNotSignificant   ImpactLevel = 1
	ImpactFairlySignificant ImpactLevel = 2
	ImpactSignificant      ImpactLevel = 3
	ImpactVerySignificant  ImpactLevel = 4
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

// REQUEST DTO (API INPUT)
type EligibilityRequest struct {
	RespondenID   int  `json:"responden_id"`
	RisikoID      int  `json:"risiko_id"`
	PernahTerjadi bool `json:"pernah_terjadi"`
}

type ReasonRequest struct {
	RespondenID int    `json:"responden_id"`
	RisikoID    int    `json:"risiko_id"`
	Alasan      string `json:"alasan"`
}

type DetailRequest struct {
	RespondenID int `json:"responden_id"`
	RisikoID    int `json:"risiko_id"`

	DampakReputasi    ImpactLevel `json:"dampak_reputasi"`
	DampakOperasional ImpactLevel `json:"dampak_operasional"`
	DampakFinansial   ImpactLevel `json:"dampak_finansial"`
	DampakHukum       ImpactLevel `json:"dampak_hukum"`

	Frekuensi FrequencyLevel `json:"frekuensi"`
}

type ControlRequest struct {
	RespondenID int  `json:"responden_id"`
	RisikoID    int  `json:"risiko_id"`
	AdaKontrol  bool `json:"ada_pengendalian"`

	DeskripsiPengendalian string `json:"deskripsi_pengendalian,omitempty"`
}

type NavigateRequest struct {
	RespondenID int    `json:"responden_id"`
	Direction   string `json:"direction"` // next | previous
	CurrentRisk int    `json:"current_risk"`
}

// RESPONSE DTO
type EligibilityResponse struct {
	RespondenID   int    `json:"responden_id"`
	RisikoID      int    `json:"risiko_id"`
	PernahTerjadi bool   `json:"pernah_terjadi"`
	NextStep      string `json:"next_step"` // show_reason | show_detail
}

type DetailResponse struct {
	RespondenID int    `json:"responden_id"`
	RisikoID    int    `json:"risiko_id"`
	NextStep    string `json:"next_step"` // show_control
}

type ControlResponse struct {
	RespondenID int    `json:"responden_id"`
	RisikoID    int    `json:"risiko_id"`
	NextStep    string `json:"next_step"` // finish
}

// DATABASE MODELS (ENTITY)

// STEP 1
type RisikoEligibility struct {
	ID               int       `db:"id"`
	RespondenID      int       `db:"responden_id"`
	RisikoID         int       `db:"risiko_id"`
	PernahTerjadi    bool      `db:"pernah_terjadi"`
	LangkahSelanjutnya string   `db:"langkah_selanjutnya"`
	CreatedAt        time.Time `db:"created_at"`
}

// STEP 2a (TIDAK)
type RisikoAlasan struct {
	ID          int       `db:"id"`
	RespondenID int       `db:"responden_id"`
	RisikoID    int       `db:"risiko_id"`
	Alasan      string    `db:"alasan"`
	Selesai     bool      `db:"selesai"`
	CreatedAt   time.Time `db:"created_at"`
}

// STEP 2b (YA)
type RisikoDampak struct {
	ID                  int            `db:"id"`
	RespondenID         int            `db:"responden_id"`
	RisikoID            int            `db:"risiko_id"`
	DampakReputasi      ImpactLevel    `db:"dampak_reputasi"`
	DampakOperasional   ImpactLevel    `db:"dampak_operasional"`
	DampakFinansial     ImpactLevel    `db:"dampak_finansial"`
	DampakHukum         ImpactLevel    `db:"dampak_hukum"`
	Frekuensi           FrequencyLevel `db:"frekuensi"`
	LangkahSelanjutnya  string         `db:"langkah_selanjutnya"`
	CreatedAt           time.Time      `db:"created_at"`
}

// STEP 2c
type RisikoPengendalian struct {
	ID                     int       `db:"id"`
	RespondenID            int       `db:"responden_id"`
	RisikoID               int       `db:"risiko_id"`
	AdaPengendalian        bool      `db:"ada_pengendalian"`
	DeskripsiPengendalian  string    `db:"deskripsi_pengendalian"`
	Selesai                bool      `db:"selesai"`
	CreatedAt              time.Time `db:"created_at"`
}

// PROGRESS MODEL
type SurveyProgress struct {
	ID               int       `db:"id"`
	RespondenID      int       `db:"responden_id"`
	RisikoID         int       `db:"risiko_id"`
	LangkahSaatIni   string    `db:"langkah_saat_ini"`
	NomorRisiko      int       `db:"nomor_risiko"`
	Selesai          bool      `db:"selesai"`
	TerakhirUpdate   time.Time `db:"terakhir_update"`
}

// GLOBAL RESPONSE
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}