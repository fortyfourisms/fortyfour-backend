package validation

import (
	"errors"
	"survey/internal/dto"
)

var (
	ErrMissingRespondentID = errors.New("responden_id wajib diisi")
	ErrMissingRisikoID     = errors.New("risiko_id wajib diisi")

	ErrMissingReason = errors.New("alasan wajib diisi")

	ErrInvalidImpact  = errors.New("nilai dampak harus 1-4")
	ErrInvalidFreq    = errors.New("frekuensi harus 1-4")

	ErrMissingControl = errors.New("deskripsi pengendalian wajib diisi jika ada_pengendalian = true")
)

// STEP 1 — ELIGIBILITY
func ValidateEligibilityRequest(req dto.EligibilityRequest) error {
	if req.RespondenID <= 0 {
		return ErrMissingRespondentID
	}
	if req.RisikoID <= 0 {
		return ErrMissingRisikoID
	}
	return nil
}

// STEP 2A — ALASAN
func ValidateAlasanRequest(req dto.AlasanRequest) error {
	if req.RespondenID <= 0 {
		return ErrMissingRespondentID
	}
	if req.RisikoID <= 0 {
		return ErrMissingRisikoID
	}
	if req.Alasan == "" {
		return ErrMissingReason
	}
	return nil
}

// STEP 2B — DAMPAK
func ValidateDampakRequest(req dto.DampakRequest) error {
	if req.RespondenID <= 0 {
		return ErrMissingRespondentID
	}
	if req.RisikoID <= 0 {
		return ErrMissingRisikoID
	}

	// karena ImpactLevel = int (1-4)
	if req.DampakReputasi < 1 || req.DampakReputasi > 4 {
		return ErrInvalidImpact
	}
	if req.DampakOperasional < 1 || req.DampakOperasional > 4 {
		return ErrInvalidImpact
	}
	if req.DampakFinansial < 1 || req.DampakFinansial > 4 {
		return ErrInvalidImpact
	}
	if req.DampakHukum < 1 || req.DampakHukum > 4 {
		return ErrInvalidImpact
	}

	if req.Frekuensi < 1 || req.Frekuensi > 4 {
		return ErrInvalidFreq
	}

	return nil
}

// STEP 2C — PENGENDALIAN
func ValidatePengendalianRequest(req dto.PengendalianRequest) error {
	if req.RespondenID <= 0 {
		return ErrMissingRespondentID
	}
	if req.RisikoID <= 0 {
		return ErrMissingRisikoID
	}

	if req.AdaPengendalian && req.DeskripsiPengendalian == "" {
		return ErrMissingControl
	}

	return nil
}