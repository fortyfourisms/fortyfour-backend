package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"survey/internal/dto"
	"survey/internal/middleware"
	"survey/internal/repository"
)

// MOCK SERVICE
type mockRisikoService struct {
	ProcessEligibilityFunc  func(string, dto.EligibilityRequest) (map[string]interface{}, error)
	ProcessAlasanFunc       func(string, dto.AlasanRequest) (map[string]interface{}, error)
	ProcessDampakFunc       func(string, dto.DampakRequest) (map[string]interface{}, error)
	ProcessPengendalianFunc func(string, dto.PengendalianRequest) (map[string]interface{}, error)
	GetByUserIDFunc         func(string) (map[string]interface{}, error)
	GetByRespondentIDFunc   func(int64) (map[string]interface{}, error)
	GetProgressFunc         func(string) (dto.ProgressResponse, error)
	NavigateFunc            func(string, dto.NavigateRequest) (dto.ProgressResponse, error)
	SaveProgressFunc        func(string, dto.NavigateRequest) (dto.ProgressResponse, error)
	FinishSurveyFunc        func(string) error
}

func (m *mockRisikoService) ProcessEligibility(userID string, r dto.EligibilityRequest) (map[string]interface{}, error) {
	return m.ProcessEligibilityFunc(userID, r)
}
func (m *mockRisikoService) ProcessAlasan(userID string, r dto.AlasanRequest) (map[string]interface{}, error) {
	return m.ProcessAlasanFunc(userID, r)
}
func (m *mockRisikoService) ProcessDampak(userID string, r dto.DampakRequest) (map[string]interface{}, error) {
	return m.ProcessDampakFunc(userID, r)
}
func (m *mockRisikoService) ProcessPengendalian(userID string, r dto.PengendalianRequest) (map[string]interface{}, error) {
	return m.ProcessPengendalianFunc(userID, r)
}
func (m *mockRisikoService) GetByUserID(userID string) (map[string]interface{}, error) {
	return m.GetByUserIDFunc(userID)
}
func (m *mockRisikoService) GetByRespondentID(id int64) (map[string]interface{}, error) {
	return m.GetByRespondentIDFunc(id)
}
func (m *mockRisikoService) GetProgress(userID string) (dto.ProgressResponse, error) {
	return m.GetProgressFunc(userID)
}
func (m *mockRisikoService) Navigate(userID string, r dto.NavigateRequest) (dto.ProgressResponse, error) {
	return m.NavigateFunc(userID, r)
}
func (m *mockRisikoService) SaveProgress(userID string, r dto.NavigateRequest) (dto.ProgressResponse, error) {
	return m.SaveProgressFunc(userID, r)
}
func (m *mockRisikoService) FinishSurvey(userID string) error {
	return m.FinishSurveyFunc(userID)
}

// helper: inject role and userID into request context
func withRisikoCtx(req *http.Request, userID, role string) *http.Request {
	ctx := req.Context()
	ctx = middleware.SetUserID(ctx, userID)
	ctx = middleware.SetRole(ctx, role)
	return req.WithContext(ctx)
}

// ELIGIBILITY
func TestSubmitEligibility_Success(t *testing.T) {
	mock := &mockRisikoService{
		ProcessEligibilityFunc: func(string, dto.EligibilityRequest) (map[string]interface{}, error) {
			return map[string]interface{}{"message": "ok"}, nil
		},
	}

	h := NewRisikoHandler(mock)

	body, _ := json.Marshal(dto.EligibilityRequest{})
	req := httptest.NewRequest(http.MethodPost, "/eligibility", bytes.NewBuffer(body))
	req = withRisikoCtx(req, "user1", "user")
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
		GetByRespondentIDFunc: func(int64) (map[string]interface{}, error) {
			return map[string]interface{}{"data": "ok"}, nil
		},
	}

	h := NewRisikoHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/survey/risiko/1", nil)
	req = withRisikoCtx(req, "admin1", "admin")
	w := httptest.NewRecorder()

	h.GetByRespondentID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetByRespondentID_NotFound(t *testing.T) {
	mock := &mockRisikoService{
		GetByRespondentIDFunc: func(int64) (map[string]interface{}, error) {
			return nil, repository.ErrNotFound
		},
	}

	h := NewRisikoHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/survey/risiko/1", nil)
	req = withRisikoCtx(req, "admin1", "admin")
	w := httptest.NewRecorder()

	h.GetByRespondentID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestGetByRespondentID_InvalidID(t *testing.T) {
	h := NewRisikoHandler(&mockRisikoService{})

	req := httptest.NewRequest(http.MethodGet, "/api/survey/risiko/abc", nil)
	req = withRisikoCtx(req, "admin1", "admin")
	w := httptest.NewRecorder()

	h.GetByRespondentID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// FINISH SURVEY
func TestFinishSurvey_Success(t *testing.T) {
	mock := &mockRisikoService{
		FinishSurveyFunc: func(string) error {
			return nil
		},
	}

	h := NewRisikoHandler(mock)

	req := httptest.NewRequest(http.MethodPost, "/finish", nil)
	req = withRisikoCtx(req, "user1", "user")
	w := httptest.NewRecorder()

	h.FinishSurvey(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestFinishSurvey_Error(t *testing.T) {
	mock := &mockRisikoService{
		FinishSurveyFunc: func(string) error {
			return errors.New("fail")
		},
	}

	h := NewRisikoHandler(mock)

	req := httptest.NewRequest(http.MethodPost, "/finish", nil)
	req = withRisikoCtx(req, "user1", "user")
	w := httptest.NewRecorder()

	h.FinishSurvey(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
