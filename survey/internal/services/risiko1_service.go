package service

import (
	"errors"

	"survey-backend/model"
	"survey-backend/repository"
	"survey-backend/validation"
)

type IPTheftService struct {
	repo     *repository.IPTheftRepository
	progress *repository.ProgressRepository
}

func NewIPTheftService(
	repo *repository.IPTheftRepository,
	progress *repository.ProgressRepository,
) *IPTheftService {
	return &IPTheftService{repo: repo, progress: progress}
}

// STEP 1 — Eligibility
// POST /api/survey/risk/ip-theft/eligibility
// Pertanyaan:
//   "Apakah perusahaan Anda berpotensi mengalami atau pernah mengalami
//    insiden pencurian Intellectual Property?"
// Branching:
//   has_experienced = true  → next_step: "show_detail"   (alur Ya)
//   has_experienced = false → next_step: "show_reason"   (alur Tidak)

func (s *IPTheftService) ProcessEligibility(req model.EligibilityRequest) (*model.EligibilityResponse, error) {
	if err := validation.ValidateEligibilityRequest(req); err != nil {
		return nil, err
	}

	s.progress.GetOrCreate(req.RespondentID)

	record := &model.IPTheftResponse{
		RespondentID:   req.RespondentID,
		HasExperienced: req.HasExperienced,
		CurrentStep:    model.StepEligibility,
	}
	if err := s.repo.Upsert(record); err != nil {
		return nil, err
	}

	nextStep := "show_reason"
	if req.HasExperienced {
		nextStep = "show_detail"
	}

	return &model.EligibilityResponse{
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

func (s *IPTheftService) ProcessReason(req model.ReasonRequest) (*model.IPTheftResponse, error) {
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
	existing.CurrentStep = model.StepDone

	if err := s.repo.Upsert(existing); err != nil {
		return nil, err
	}

	s.progress.MarkCompleted(req.RespondentID, 1)
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

func (s *IPTheftService) ProcessDetail(req model.DetailRequest) (*model.DetailResponse, error) {
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
	existing.CurrentStep = model.StepDetail

	if err := s.repo.Upsert(existing); err != nil {
		return nil, err
	}

	// Selalu lanjut ke step pengendalian — belum selesai
	return &model.DetailResponse{
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

func (s *IPTheftService) ProcessControl(req model.ControlRequest) (*model.ControlResponse, error) {
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
	if existing.CurrentStep != model.StepDetail {
		return nil, errors.New("langkah dampak & frekuensi (detail) harus diisi sebelum tindakan pengendalian")
	}

	// Simpan jawaban pengendalian
	hasControl := req.HasControl
	existing.HasControl = &hasControl
	existing.CurrentStep = model.StepControl

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
	existing.CurrentStep = model.StepDone
	_ = s.repo.Upsert(existing)
	s.progress.MarkCompleted(req.RespondentID, 1)

	return &model.ControlResponse{
		RespondentID:    req.RespondentID,
		HasControl:      req.HasControl,
		ControlMeasures: req.ControlMeasures,
		NextStep:        "finish", // tombol Berikutnya aktif
	}, nil
}

// Query
func (s *IPTheftService) GetResponse(respondentID string) (*model.IPTheftResponse, error) {
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

// Navigation Service
type NavigationService struct {
	progress *repository.ProgressRepository
}

func NewNavigationService(progress *repository.ProgressRepository) *NavigationService {
	return &NavigationService{progress: progress}
}

func (s *NavigationService) Navigate(req model.NavigateRequest) (*model.SurveyProgress, error) {
	p := s.progress.GetOrCreate(req.RespondentID)
	switch req.Direction {
	case "next":
		if p.CurrentRisk < p.TotalRisks {
			s.progress.SetCurrentRisk(req.RespondentID, p.CurrentRisk+1)
		}
	case "previous":
		if p.CurrentRisk > 1 {
			s.progress.SetCurrentRisk(req.RespondentID, p.CurrentRisk-1)
		}
	}
	return s.progress.Get(req.RespondentID)
}

func (s *NavigationService) GetProgress(respondentID string) (*model.SurveyProgress, error) {
	return s.progress.GetOrCreate(respondentID), nil
}