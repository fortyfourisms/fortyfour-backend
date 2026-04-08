package validation

import (
	"errors"
	"strings"
	"survey/internal/models"
)

var (
	ErrMissingRespondentID = errors.New("respondent_id wajib diisi")
	ErrMissingReason       = errors.New("alasan wajib diisi ketika perusahaan tidak mengalami insiden")
	ErrInvalidImpactLevel  = errors.New("nilai dampak tidak valid (harus 1–4)")
	ErrInvalidFrequency    = errors.New("nilai frekuensi tidak valid (harus 1–4)")
	ErrMissingControlInfo  = errors.New("tindakan pengendalian wajib diisi jika has_control = true")
)

func ValidateEligibilityRequest(req models.EligibilityRequest) error {
	if strings.TrimSpace(req.RespondentID) == "" {
		return ErrMissingRespondentID
	}
	return nil
}

func ValidateReasonRequest(req models.ReasonRequest) error {
	if strings.TrimSpace(req.RespondentID) == "" {
		return ErrMissingRespondentID
	}
	if strings.TrimSpace(req.Reason) == "" {
		return ErrMissingReason
	}
	return nil
}

// ValidateDetailRequest — validasi Step 2b: matrix dampak + frekuensi
// (tidak lagi termasuk validasi pengendalian — itu ada di ValidateControlRequest)
func ValidateDetailRequest(req models.DetailRequest) error {
	if strings.TrimSpace(req.RespondentID) == "" {
		return ErrMissingRespondentID
	}
	if err := validateImpactMatrix(req.Impact); err != nil {
		return err
	}
	if !req.Frequency.Valid() {
		return ErrInvalidFrequency
	}
	return nil
}

// ValidateControlRequest — validasi Step 2c: tindakan pengendalian
//
//	has_control = true  → ControlMeasures wajib diisi
//	has_control = false → ControlMeasures diabaikan, langsung Berikutnya
func ValidateControlRequest(req models.ControlRequest) error {
	if strings.TrimSpace(req.RespondentID) == "" {
		return ErrMissingRespondentID
	}
	if req.HasControl && strings.TrimSpace(req.ControlMeasures) == "" {
		return ErrMissingControlInfo
	}
	return nil
}

func validateImpactMatrix(m models.ImpactMatrix) error {
	if !m.Reputation.Valid() {
		return errors.New("dampak reputasi tidak valid (harus 1–4)")
	}
	if !m.Operational.Valid() {
		return errors.New("dampak operasional tidak valid (harus 1–4)")
	}
	if !m.Financial.Valid() {
		return errors.New("dampak finansial tidak valid (harus 1–4)")
	}
	if !m.Legal.Valid() {
		return errors.New("dampak hukum tidak valid (harus 1–4)")
	}
	return nil
}
