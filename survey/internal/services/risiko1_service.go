package services

import (
	"errors"

	"survey/internal/models"
	"survey/internal/repository"
	"survey/validator"
)

type RisikoService struct {
	repo *repository.RisikoRepository
}

func NewRisikoService(repo *repository.RisikoRepository) *RisikoService {
	return &RisikoService{repo: repo}
}

// STEP 1 — Eligibility
// POST /api/survey/risk/ip-theft/eligibility
// Pertanyaan:
//   "Apakah perusahaan Anda berpotensi mengalami atau pernah mengalami
//    insiden pencurian Intellectual Property?"
// Branching:
//   has_experienced = true  → next_step: "show_detail"   (alur Ya)
//   has_experienced = false → next_step: "show_reason"   (alur Tidak)

func (s *RisikoService) ProcessEligibility(req models.EligibilityRequest) (*models.EligibilityResponse, error) {
	if err := validation.ValidateEligibilityRequest(req); err != nil {
		return nil, err
	}

	s.repo.GetOrCreate(req.RespondentID)

	record := &models.IPTheftResponse{
		RespondentID:   req.RespondentID,
		HasExperienced: req.HasExperienced,
		CurrentStep:    models.StepEligibility,
	}
	if err := s.repo.Upsert(record); err != nil {
		return nil, err
	}

	nextStep := "show_reason"
	if req.HasExperienced {
		nextStep = "show_detail"
	}

	return &models.EligibilityResponse{
		RespondentID:   req.RespondentID,
		HasExperienced: req.HasExperienced,
		NextStep:       nextStep,
	}, nil
}

// STEP 2a — Reason   (alur "Tidak")
// POST /api/survey/risk/ip-theft/reason
// Syarat masuk: has_experienced = false
// Pertanyaan:
//   "Mengapa perusahaan Anda tidak berpotensi mengalami atau tidak pernah
//    mengalami insiden pencurian Intellectual Property?"
// Setelah step ini → Risiko 1 SELESAI (tombol Berikutnya aktif)

func (s *RisikoService) ProcessReason(req models.ReasonRequest) (*models.IPTheftResponse, error) {
	if err := validation.ValidateReasonRequest(req); err != nil {
		return nil, err
	}

	existing, err := s.repo.FindByRespondentID(req.RespondentID)
	if err != nil {
		return nil, err
	}

	if existing.HasExperienced {
		return nil, errWrongBranch("reason", "Ya")
	}

	existing.Reason = req.Reason
	existing.CurrentStep = models.StepDone

	if err := s.repo.Upsert(existing); err != nil {
		return nil, err
	}

	s.repo.MarkCompleted(req.RespondentID, 1)
	return existing, nil
}

// STEP 2b — Detail   (alur "Ya")
// POST /api/survey/risk/ip-theft/detail
// Syarat masuk: has_experienced = true
// Pertanyaan:
//   1. "Seberapa besar dampak dari pencurian Intellectual Property perusahaan?"
//      → Matrix 4 dimensi (Reputasi / Operasional / Finansial / Hukum), nilai 1–4
//   2. "Seberapa sering dalam setahun risiko pencurian IP berpotensi terjadi?"
//      → Frekuensi (Kecil=1 / Sedang=2 / Besar=3 / Sangat Besar=4)
// Setelah step ini → next_step: "show_control"
//   UI wajib melanjutkan ke Step 2c (ControlRequest) sebelum tombol Berikutnya aktif

func (s *RisikoService) ProcessDetail(req models.DetailRequest) (*models.DetailResponse, error) {
	if err := validation.ValidateDetailRequest(req); err != nil {
		return nil, err
	}

	existing, err := s.repo.FindByRespondentID(req.RespondentID)
	if err != nil {
		return nil, err
	}

	if !existing.HasExperienced {
		return nil, errWrongBranch("detail", "Tidak")
	}

	existing.Impact = &req.Impact
	existing.Frequency = &req.Frequency
	existing.CurrentStep = models.StepDetail

	if err := s.repo.Upsert(existing); err != nil {
		return nil, err
	}

	// Selalu lanjut ke step pengendalian — belum selesai
	return &models.DetailResponse{
		RespondentID: req.RespondentID,
		NextStep:     "show_control",
	}, nil
}

// STEP 2c — Control   (alur "Ya", sub-branching pengendalian)
// POST /api/survey/risk/ip-theft/control
// Syarat masuk: has_experienced = true AND step sebelumnya = "detail"
// Pertanyaan wajib:
//   "Apa perusahaan Anda telah memiliki tindakan pengendalian terhadap
//    risiko pencurian Intellectual Property?"
//   ● Ya  → muncul pertanyaan lanjutan:
//            "Apa tindakan pengendalian yang telah dilakukan oleh perusahaan Anda
//             terhadap risiko pencurian Intellectual Property perusahaan?"
//            (wajib diisi, field: control_measures)
//   ● Tidak → langsung tombol Berikutnya aktif (tidak ada pertanyaan tambahan)
// Setelah step ini → Risiko 1 SELESAI (next_step: "finish")

func (s *RisikoService) ProcessControl(req models.ControlRequest) (*models.ControlResponse, error) {
	if err := validation.ValidateControlRequest(req); err != nil {
		return nil, err
	}

	existing, err := s.repo.FindByRespondentID(req.RespondentID)
	if err != nil {
		return nil, err
	}

	// Guard 1: hanya bisa diakses dari alur "Ya"
	if !existing.HasExperienced {
		return nil, errWrongBranch("control", "Tidak")
	}

	// Guard 2: step 2b (detail) harus sudah diisi terlebih dahulu
	if existing.CurrentStep != models.StepDetail {
		return nil, errors.New("langkah dampak & frekuensi (detail) harus diisi sebelum tindakan pengendalian")
	}

	// Simpan jawaban pengendalian
	hasControl := req.HasControl
	existing.HasControl = &hasControl
	existing.CurrentStep = models.StepControl

	if req.HasControl {
		// Alur Ya → simpan tindakan pengendalian yang diisi user
		existing.ControlMeasures = req.ControlMeasures
	} else {
		// Alur Tidak → kosongkan (jika sebelumnya ada nilai)
		existing.ControlMeasures = ""
	}

	if err := s.repo.Upsert(existing); err != nil {
		return nil, err
	}

	// Risiko 1 selesai
	existing.CurrentStep = models.StepDone
	_ = s.repo.Upsert(existing)
	s.repo.MarkCompleted(req.RespondentID, 1)

	return &models.ControlResponse{
		RespondentID:    req.RespondentID,
		HasControl:      req.HasControl,
		ControlMeasures: req.ControlMeasures,
		NextStep:        "finish", // tombol Berikutnya aktif
	}, nil
}

// Query
func (s *RisikoService) GetResponse(respondentID string) (*models.IPTheftResponse, error) {
	return s.repo.FindByRespondentID(respondentID)
}

// Helpers
type branchError struct{ endpoint, branch string }

func (e *branchError) Error() string {
	return "endpoint '" + e.endpoint + "' tidak tersedia untuk jawaban '" + e.branch + "'"
}

func errWrongBranch(endpoint, branch string) error {
	return &branchError{endpoint: endpoint, branch: branch}
}

func (s *RisikoService) Navigate(req models.NavigateRequest) (*models.SurveyProgress, error) {
	p := s.repo.GetOrCreate(req.RespondentID)
	switch req.Direction {
	case "next":
		if p.CurrentRisk < p.TotalRisks {
			s.repo.SetCurrentRisk(req.RespondentID, p.CurrentRisk+1)
		}
	case "previous":
		if p.CurrentRisk > 1 {
			s.repo.SetCurrentRisk(req.RespondentID, p.CurrentRisk-1)
		}
	}
	return s.repo.Get(req.RespondentID)
}

func (s *RisikoService) GetProgress(respondentID string) (*models.SurveyProgress, error) {
	return s.repo.GetOrCreate(respondentID), nil
}
