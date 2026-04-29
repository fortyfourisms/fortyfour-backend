package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"survey/internal/dto"
	"survey/internal/repository"
)

// MOCK SERVICE
type mockRisikoService struct {
	ProcessEligibilityFunc func(dto.EligibilityRequest) (interface{}, error)
	ProcessAlasanFunc func(dto.AlasanRequest) (interface{}, error)
	ProcessDampakFunc func(dto.DampakRequest) (interface{}, error)
	ProcessPengendalianFunc func(dto.PengendalianRequest) (interface{}, error)
	GetByRespondentIDFunc func(int) (interface{}, error)
	GetProgressFunc func(int) (interface{}, error)
	NavigateFunc func(dto.NavigateRequest) (interface{}, error)
	SaveProgressFunc func(dto.NavigateRequest) (interface{}, error)
	CreateCustomRisikoFunc func(dto.CustomRisikoRequest) (int, error)
	FinishSurveyFunc func(int) error
}

func (m *mockRisikoService) ProcessEligibility(r dto.EligibilityRequest) (interface{}, error) {
	return m.ProcessEligibilityFunc(r)
}
func (m *mockRisikoService) ProcessAlasan(r dto.AlasanRequest) (interface{}, error) {
	return m.ProcessAlasanFunc(r)
}
func (m *mockRisikoService) ProcessDampak(r dto.DampakRequest) (interface{}, error) {
	return m.ProcessDampakFunc(r)
}
func (m *mockRisikoService) ProcessPengendalian(r dto.PengendalianRequest) (interface{}, error) {
	return m.ProcessPengendalianFunc(r)
}
func (m *mockRisikoService) GetByRespondentID(id int) (interface{}, error) {
	return m.GetByRespondentIDFunc(id)
}
func (m *mockRisikoService) GetProgress(id int) (interface{}, error) {
	return m.GetProgressFunc(id)
}
func (m *mockRisikoService) Navigate(r dto.NavigateRequest) (interface{}, error) {
	return m.NavigateFunc(r)
}
func (m *mockRisikoService) SaveProgress(r dto.NavigateRequest) (interface{}, error) {
	return m.SaveProgressFunc(r)
}
func (m *mockRisikoService) CreateCustomRisiko(r dto.CustomRisikoRequest) (int, error) {
	return m.CreateCustomRisikoFunc(r)
}
func (m *mockRisikoService) FinishSurvey(id int) error {
	return m.FinishSurveyFunc(id)
}

// ELIGIBILITY
func TestSubmitEligibility_Success(t *testing.T) {
	mock := &mockRisikoService{
		ProcessEligibilityFunc: func(dto.EligibilityRequest) (interface{}, error) {
			return "ok", nil
		},
	}

	h := NewRisikoHandler(mock)

	body, _ := json.Marshal(dto.EligibilityRequest{})
	req := httptest.NewRequest(http.MethodPost, "/eligibility", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.SubmitEligibility(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestSubmitEligibility_InvalidBody(t *testing.T) {
	h := NewRisikoHandler(&mockRisikoService{})

	req := httptest.NewRequest(http.MethodPost, "/eligibility", bytes.NewBuffer([]byte("invalid")))
	w := httptest.NewRecorder()

	h.SubmitEligibility(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400")
	}
}

// GET BY RESPONDENT
func TestGetByRespondentID_Success(t *testing.T) {
	mock := &mockRisikoService{
		GetByRespondentIDFunc: func(int) (interface{}, error) {
			return "data", nil
		},
	}

	h := NewRisikoHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/survey/risiko/1", nil)
	w := httptest.NewRecorder()

	h.GetByRespondentID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200")
	}
}

func TestGetByRespondentID_NotFound(t *testing.T) {
	mock := &mockRisikoService{
		GetByRespondentIDFunc: func(int) (interface{}, error) {
			return nil, repository.ErrNotFound
		},
	}

	h := NewRisikoHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/survey/risiko/1", nil)
	w := httptest.NewRecorder()

	h.GetByRespondentID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404")
	}
}

func TestGetByRespondentID_InvalidID(t *testing.T) {
	h := NewRisikoHandler(&mockRisikoService{})

	req := httptest.NewRequest(http.MethodGet, "/api/survey/risiko/abc", nil)
	w := httptest.NewRecorder()

	h.GetByRespondentID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400")
	}
}

// FINISH SURVEY
func TestFinishSurvey_Success(t *testing.T) {
	mock := &mockRisikoService{
		FinishSurveyFunc: func(int) error {
			return nil
		},
	}

	h := NewRisikoHandler(mock)

	body := `{"responden_id":1}`
	req := httptest.NewRequest(http.MethodPost, "/finish", bytes.NewBuffer([]byte(body)))
	w := httptest.NewRecorder()

	h.FinishSurvey(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200")
	}
}

func TestFinishSurvey_Error(t *testing.T) {
	mock := &mockRisikoService{
		FinishSurveyFunc: func(int) error {
			return errors.New("fail")
		},
	}

	h := NewRisikoHandler(mock)

	body := `{"responden_id":1}`
	req := httptest.NewRequest(http.MethodPost, "/finish", bytes.NewBuffer([]byte(body)))
	w := httptest.NewRecorder()

	h.FinishSurvey(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500")
	}
}