package models

import "time"

// Enums / Constants

// ImpactLevel maps to: Tidak Signifikan, Cukup Signifikan, Signifikan, Sangat Signifikan
type ImpactLevel int

const (
	ImpactNotSignificant   ImpactLevel = 1
	ImpactFairlySignifcant ImpactLevel = 2
	ImpactSignificant      ImpactLevel = 3
	ImpactVerySignificant  ImpactLevel = 4
)

func (i ImpactLevel) Label() string {
	switch i {
	case ImpactNotSignificant:
		return "Tidak Signifikan"
	case ImpactFairlySignifcant:
		return "Cukup Signifikan"
	case ImpactSignificant:
		return "Signifikan"
	case ImpactVerySignificant:
		return "Sangat Signifikan"
	default:
		return "Unknown"
	}
}

func (i ImpactLevel) Valid() bool {
	return i >= ImpactNotSignificant && i <= ImpactVerySignificant
}

// FrequencyLevel maps to: Kecil (1), Sedang (2), Besar (3), Sangat Besar (4)
type FrequencyLevel int

const (
	FrequencySmall     FrequencyLevel = 1 // < 2 kali/tahun
	FrequencyMedium    FrequencyLevel = 2 // 2-5 kali/tahun
	FrequencyLarge     FrequencyLevel = 3 // 5-10 kali/tahun
	FrequencyVeryLarge FrequencyLevel = 4 // > 10 kali/tahun
)

func (f FrequencyLevel) Label() string {
	switch f {
	case FrequencySmall:
		return "Kecil"
	case FrequencyMedium:
		return "Sedang"
	case FrequencyLarge:
		return "Besar"
	case FrequencyVeryLarge:
		return "Sangat Besar"
	default:
		return "Unknown"
	}
}

func (f FrequencyLevel) Description() string {
	switch f {
	case FrequencySmall:
		return "Kemungkinan terjadinya tidak lebih dari 2 kali per tahun"
	case FrequencyMedium:
		return "Kemungkinan terjadinya lebih dari 2 kali / tahun, namun tidak lebih dari 5 kali / tahun"
	case FrequencyLarge:
		return "Kemungkinan terjadinya lebih dari 5 kali / tahun, namun tidak lebih dari 10 kali / tahun"
	case FrequencyVeryLarge:
		return "Kemungkinan terjadinya lebih dari 10 kali / tahun"
	default:
		return ""
	}
}

func (f FrequencyLevel) Valid() bool {
	return f >= FrequencySmall && f <= FrequencyVeryLarge
}

// Request Payloads

// EligibilityRequest — Pertanyaan pertama:
// "Apakah perusahaan Anda berpotensi mengalami atau pernah mengalami insiden pencurian IP?"
type EligibilityRequest struct {
	RespondentID string `json:"respondent_id"`
	// true = Ya, false = Tidak
	HasExperienced bool `json:"has_experienced"`
}

// ReasonRequest — Ditampilkan HANYA jika has_experienced = false (alur "Tidak")
// "Mengapa perusahaan Anda tidak berpotensi mengalami atau tidak pernah mengalami insiden pencurian IP?"
type ReasonRequest struct {
	RespondentID string `json:"respondent_id"`
	Reason       string `json:"reason"`
}

// ImpactMatrix — Matrix dampak 4 dimensi (ditampilkan jika has_experienced = true)
type ImpactMatrix struct {
	Reputation  ImpactLevel `json:"reputation"`  // Reputasi
	Operational ImpactLevel `json:"operational"` // Operasional
	Financial   ImpactLevel `json:"financial"`   // Finansial
	Legal       ImpactLevel `json:"legal"`       // Hukum
}

// DetailRequest — Step 2b (alur "Ya")
// Pertanyaan: matrix dampak (4 dimensi) + frekuensi kejadian.
// Setelah ini, backend mengembalikan next_step untuk sub-branching pengendalian.
type DetailRequest struct {
	RespondentID string `json:"respondent_id"`

	// Seberapa besar dampak dari pencurian IP? (matrix 4 dimensi, nilai 1–4)
	Impact ImpactMatrix `json:"impact"`

	// Seberapa sering dalam setahun risiko berpotensi terjadi? (nilai 1–4)
	Frequency FrequencyLevel `json:"frequency"`
}

// DetailResponse — response setelah submit detail
// next_step menentukan sub-branching berikutnya di sisi UI:
//
//	"show_control_measures" → user menjawab has_control = true  → tampilkan textarea
//	"finish"                → user menjawab has_control = false → langsung tombol Berikutnya
type DetailResponse struct {
	RespondentID string `json:"respondent_id"`
	NextStep     string `json:"next_step"`
}

// ControlRequest — Step 2c: sub-branching pengendalian risiko
// Pertanyaan wajib: "Apa perusahaan Anda telah memiliki tindakan pengendalian?"
//
//	has_control = true  → wajib isi ControlMeasures
//	has_control = false → langsung selesai (tombol berikutnya aktif)
type ControlRequest struct {
	RespondentID string `json:"respondent_id"`

	// true = Ya, false = Tidak
	HasControl bool `json:"has_control"`

	// Wajib diisi HANYA jika has_control = true
	// "Apa tindakan pengendalian yang telah dilakukan?"
	ControlMeasures string `json:"control_measures,omitempty"`
}

// ControlResponse — response setelah submit kontrol
// next_step:
//
//	"finish" → risiko 1 selesai, tampilkan tombol Berikutnya
type ControlResponse struct {
	RespondentID    string `json:"respondent_id"`
	HasControl      bool   `json:"has_control"`
	ControlMeasures string `json:"control_measures,omitempty"`
	NextStep        string `json:"next_step"` // selalu "finish"
}

// NavigateRequest — untuk navigasi Sebelumnya / Berikutnya
type NavigateRequest struct {
	RespondentID string `json:"respondent_id"`
	Direction    string `json:"direction"` // "next" | "previous"
	CurrentRisk  int    `json:"current_risk"`
}

// Domain Models (Storage)

// SurveyStep merepresentasikan tahapan pengisian dalam satu risiko.
// Digunakan untuk melacak posisi user dalam alur multi-step.
type SurveyStep string

const (
	StepEligibility SurveyStep = "eligibility" // Step 1  — sudah jawab Ya/Tidak
	StepDetail      SurveyStep = "detail"      // Step 2b — sudah isi dampak & frekuensi (alur Ya)
	StepControl     SurveyStep = "control"     // Step 2c — sudah isi tindakan pengendalian (alur Ya)
	StepReason      SurveyStep = "reason"      // Step 2a — sudah isi alasan (alur Tidak)
	StepDone        SurveyStep = "done"        // Semua step selesai
)

// IPTheftResponse — rekaman jawaban lengkap Risiko 1
type IPTheftResponse struct {
	ID           string `json:"id"`
	RespondentID string `json:"respondent_id"`
	RiskNumber   int    `json:"risk_number"` // 1
	RiskName     string `json:"risk_name"`

	// Melacak step terakhir yang sudah diisi
	CurrentStep SurveyStep `json:"current_step"`

	// Step 1 — Eligibility
	HasExperienced bool `json:"has_experienced"`

	// Step 2a — Alur "Tidak"
	Reason string `json:"reason,omitempty"`

	// Step 2b — Alur "Ya": dampak & frekuensi
	Impact    *ImpactMatrix   `json:"impact,omitempty"`
	Frequency *FrequencyLevel `json:"frequency,omitempty"`

	// Step 2c — Alur "Ya": tindakan pengendalian
	// Pointer bool agar bisa dibedakan antara "belum diisi" (nil) vs "diisi false"
	HasControl      *bool  `json:"has_control,omitempty"`
	ControlMeasures string `json:"control_measures,omitempty"` // hanya ada jika HasControl = true

	// Meta
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SurveyProgress — status progres keseluruhan survey (14 risiko)
type SurveyProgress struct {
	RespondentID    string  `json:"respondent_id"`
	CurrentRisk     int     `json:"current_risk"`
	TotalRisks      int     `json:"total_risks"`
	PercentComplete float64 `json:"percent_complete"`
	CompletedRisks  []int   `json:"completed_risks"`
}

// Response Wrappers
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type EligibilityResponse struct {
	RespondentID   string `json:"respondent_id"`
	HasExperienced bool   `json:"has_experienced"`
	// NextStep menentukan form apa yang ditampilkan UI berikutnya:
	//   "show_reason" → tampilkan form alasan (alur Tidak)
	//   "show_detail" → tampilkan matrix dampak + frekuensi (alur Ya /Step 2b)
	NextStep string `json:"next_step"`
}
