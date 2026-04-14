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

// STEP 1 — ELIGIBILITY
func (s *RisikoService) ProcessEligibility(req models.EligibilityRequest) (*models.EligibilityResponse, error) {
	if err := validation.ValidateEligibilityRequest(req); err != nil {
		return nil, err
	}

	err := s.repo.UpsertEligibility(req.RespondentID, req.HasExperienced)
	if err != nil {
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

// STEP 2a — REASON (Tidak)
func (s *RisikoService) ProcessReason(req models.ReasonRequest) error {
	if err := validation.ValidateReasonRequest(req); err != nil {
		return err
	}

	hasExp, err := s.repo.GetEligibility(req.RespondentID)
	if err != nil {
		return err
	}

	if hasExp {
		return errors.New("tidak bisa isi alasan jika memilih 'Ya'")
	}

	err = s.repo.InsertReason(req.RespondentID, req.Reason)
	if err != nil {
		return err
	}

	s.repo.MarkCompleted(req.RespondentID, 1)

	return nil
}

// STEP 2b — DETAIL (Ya)
func (s *RisikoService) ProcessDetail(req models.DetailRequest) (*models.DetailResponse, error) {
	if err := validation.ValidateDetailRequest(req); err != nil {
		return nil, err
	}

	hasExp, err := s.repo.GetEligibility(req.RespondentID)
	if err != nil {
		return nil, err
	}

	if !hasExp {
		return nil, errors.New("tidak bisa isi detail jika memilih 'Tidak'")
	}

	err = s.repo.InsertDetail(
		req.RespondentID,
		int(req.Impact.Reputation),
		int(req.Impact.Operational),
		int(req.Impact.Financial),
		int(req.Impact.Legal),
		int(req.Frequency),
	)
	if err != nil {
		return nil, err
	}

	return &models.DetailResponse{
		RespondentID: req.RespondentID,
		NextStep:     "show_control",
	}, nil
}

// STEP 2c — CONTROL (Ya)
func (s *RisikoService) ProcessControl(req models.ControlRequest) (*models.ControlResponse, error) {
	if err := validation.ValidateControlRequest(req); err != nil {
		return nil, err
	}

	hasExp, err := s.repo.GetEligibility(req.RespondentID)
	if err != nil {
		return nil, err
	}

	if !hasExp {
		return nil, errors.New("tidak bisa isi pengendalian jika memilih 'Tidak'")
	}

	err = s.repo.InsertControl(
		req.RespondentID,
		req.HasControl,
		req.ControlMeasures,
	)
	if err != nil {
		return nil, err
	}

	s.repo.MarkCompleted(req.RespondentID, 1)

	return &models.ControlResponse{
		RespondentID:    req.RespondentID,
		HasControl:      req.HasControl,
		ControlMeasures: req.ControlMeasures,
		NextStep:        "finish",
	}, nil
}

// GET RESPONSE (Gabungan)
func (s *RisikoService) GetResponse(respondentID string) (*models.IPTheftResponse, error) {
	return s.repo.GetFullResponse(respondentID)
}

// PROGRESS
func (s *RisikoService) GetProgress(respondentID string) (*models.SurveyProgress, error) {
	return s.repo.GetProgress(respondentID)
}

// NAVIGATE
func (s *RisikoService) Navigate(req models.NavigateRequest) (*models.SurveyProgress, error) {
	progress, err := s.repo.GetProgress(req.RespondentID)
	if err != nil {
		return nil, err
	}

	switch req.Direction {
	case "next":
		if progress.CurrentRisk < progress.TotalRisks {
			progress.CurrentRisk++
		}
	case "previous":
		if progress.CurrentRisk > 1 {
			progress.CurrentRisk--
		}
	}

	err = s.repo.UpdateCurrentRisk(req.RespondentID, progress.CurrentRisk)
	if err != nil {
		return nil, err
	}

	return progress, nil
}