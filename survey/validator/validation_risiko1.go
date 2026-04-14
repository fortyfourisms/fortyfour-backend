package validation

import (
	"errors"
	"strings"
	"survey/internal/models"
)

var (
	ErrMissingRespondentID = errors.New("respondent_id wajib diisi")

	// Reason (alur Tidak)
	ErrMissingReason = errors.New("alasan wajib diisi ketika perusahaan tidak mengalami insiden")

	// Dampak & Frekuensi (alur Ya)
	ErrInvalidImpactReputation  = errors.New("dampak reputasi tidak valid")
	ErrInvalidImpactOperational = errors.New("dampak operasional tidak valid")
	ErrInvalidImpactFinancial   = errors.New("dampak finansial tidak valid")
	ErrInvalidImpactLegal       = errors.New("dampak hukum tidak valid")
	ErrInvalidFrequency         = errors.New("frekuensi tidak valid")

	// Pengendalian
	ErrMissingControlInfo = errors.New("deskripsi pengendalian wajib diisi jika has_control = true")
)

// STEP 1 — ELIGIBILITY
func ValidateEligibilityRequest(req models.EligibilityRequest) error {
	if strings.TrimSpace(req.RespondentID) == "" {
		return ErrMissingRespondentID
	}

	// Tidak perlu validasi boolean karena default Go sudah aman
	return nil
}

// STEP 2a — REASON (Tidak)
func ValidateReasonRequest(req models.ReasonRequest) error {
	if strings.TrimSpace(req.RespondentID) == "" {
		return ErrMissingRespondentID
	}

	if strings.TrimSpace(req.Reason) == "" {
		return ErrMissingReason
	}

	return nil
}

// STEP 2b — DETAIL (Ya)
func ValidateDetailRequest(req models.DetailRequest) error {
	if strings.TrimSpace(req.RespondentID) == "" {
		return ErrMissingRespondentID
	}

	// Validasi matrix dampak (sesuai CHECK DB 1–4)
	if !req.Impact.Reputation.Valid() {
		return ErrInvalidImpactReputation
	}
	if !req.Impact.Operational.Valid() {
		return ErrInvalidImpactOperational
	}
	if !req.Impact.Financial.Valid() {
		return ErrInvalidImpactFinancial
	}
	if !req.Impact.Legal.Valid() {
		return ErrInvalidImpactLegal
	}

	// Validasi frekuensi (1–4)
	if !req.Frequency.Valid() {
		return ErrInvalidFrequency
	}

	return nil
}

// STEP 2c — CONTROL (Ya)
func ValidateControlRequest(req models.ControlRequest) error {
	if strings.TrimSpace(req.RespondentID) == "" {
		return ErrMissingRespondentID
	}

	// Constraint DB:
	// ada_pengendalian = true → deskripsi wajib
	if req.HasControl {
		if strings.TrimSpace(req.ControlMeasures) == "" {
			return ErrMissingControlInfo
		}
	}

	return nil
}