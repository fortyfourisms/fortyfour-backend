package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"survey/internal/dto"
	"survey/internal/middleware"
	"survey/internal/models"
	"survey/internal/repository"
)

// MOCK SERVICE
type mockRisikoService struct {
	GetAllRisikoFunc        func() ([]models.RisikoResponse, error)
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
	RequestEditFunc         func(string, dto.RequestEditRequest) (dto.ProgressResponse, error)
	ReviewEditRequestFunc   func(string, int64, dto.ReviewEditRequest) (dto.ProgressResponse, error)
	GetAllEditRequestsFunc  func() ([]dto.EditRequestItemResponse, error)
	GetMyEditRequestFunc    func(string) (*dto.EditRequestItemResponse, error)
}

func (m *mockRisikoService) GetAllRisiko() ([]models.RisikoResponse, error) {
	if m.GetAllRisikoFunc != nil {
		return m.GetAllRisikoFunc()
	}
	return nil, nil
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
func (m *mockRisikoService) RequestEdit(userID string, r dto.RequestEditRequest) (dto.ProgressResponse, error) {
	return m.RequestEditFunc(userID, r)
}
func (m *mockRisikoService) ReviewEditRequest(adminID string, respondenID int64, r dto.ReviewEditRequest) (dto.ProgressResponse, error) {
	return m.ReviewEditRequestFunc(adminID, respondenID, r)
}
func (m *mockRisikoService) GetAllEditRequests() ([]dto.EditRequestItemResponse, error) {
	if m.GetAllEditRequestsFunc != nil {
		return m.GetAllEditRequestsFunc()
	}
	return nil, nil
}
func (m *mockRisikoService) GetMyEditRequest(userID string) (*dto.EditRequestItemResponse, error) {
	if m.GetMyEditRequestFunc != nil {
		return m.GetMyEditRequestFunc(userID)
	}
	return nil, nil
}

// helper: inject role and userID into request context
func withRisikoCtx(req *http.Request, userID, role string) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, middleware.UserIDKey, userID)
	ctx = context.WithValue(ctx, middleware.RoleKey, role)
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
