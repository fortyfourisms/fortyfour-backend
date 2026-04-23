package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSEEditRequestService struct {
	mock.Mock
}

func (m *mockSEEditRequestService) CreateRequest(userID, idSE string, req dto.CreateSEEditRequestDTO) (*dto.SEEditRequestResponse, error) {
	args := m.Called(userID, idSE, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SEEditRequestResponse), args.Error(1)
}

func (m *mockSEEditRequestService) GetPending() ([]dto.SEEditRequestResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.SEEditRequestResponse), args.Error(1)
}

func (m *mockSEEditRequestService) GetByUser(userID string) ([]dto.SEEditRequestResponse, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.SEEditRequestResponse), args.Error(1)
}

func (m *mockSEEditRequestService) Review(id string, req dto.ReviewSEEditRequestDTO) (*dto.SEEditRequestResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SEEditRequestResponse), args.Error(1)
}

func withRoleAndUser(req *http.Request, role, userID string) *http.Request {
	ctx := context.WithValue(req.Context(), middleware.RoleKey, role)
	if userID != "" {
		ctx = context.WithValue(ctx, middleware.UserIDKey, userID)
	}
	return req.WithContext(ctx)
}

func TestSEEditRequestHandler_ListAsAdmin(t *testing.T) {
	svc := new(mockSEEditRequestService)
	h := NewSEEditRequestHandler(svc)

	expected := []dto.SEEditRequestResponse{
		{ID: "req-1", Status: models.SEEditRequestPending},
	}
	svc.On("GetPending").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/se/edit-requests", nil)
	req = withRoleAndUser(req, "admin", "")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []dto.SEEditRequestResponse
	_ = json.NewDecoder(w.Body).Decode(&resp)
	assert.Len(t, resp, 1)
	assert.Equal(t, "req-1", resp[0].ID)
	svc.AssertExpectations(t)
}

func TestSEEditRequestHandler_ListAsUser(t *testing.T) {
	svc := new(mockSEEditRequestService)
	h := NewSEEditRequestHandler(svc)

	expected := []dto.SEEditRequestResponse{
		{ID: "req-1", IDUser: "user-1"},
	}
	svc.On("GetByUser", "user-1").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/se/edit-requests", nil)
	req = withRoleAndUser(req, "user", "user-1")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []dto.SEEditRequestResponse
	_ = json.NewDecoder(w.Body).Decode(&resp)
	assert.Len(t, resp, 1)
	assert.Equal(t, "user-1", resp[0].IDUser)
	svc.AssertExpectations(t)
}

func TestSEEditRequestHandler_ListAsUserUnauthorized(t *testing.T) {
	h := NewSEEditRequestHandler(new(mockSEEditRequestService))

	req := httptest.NewRequest(http.MethodGet, "/api/se/edit-requests", nil)
	req = withRoleAndUser(req, "user", "")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSEEditRequestHandler_ReviewSuccess(t *testing.T) {
	svc := new(mockSEEditRequestService)
	h := NewSEEditRequestHandler(svc)
	catatan := "ok"

	svc.On("Review", "req-1", dto.ReviewSEEditRequestDTO{
		Status:  "approved",
		Catatan: &catatan,
	}).Return(&dto.SEEditRequestResponse{
		ID:     "req-1",
		Status: models.SEEditRequestApproved,
	}, nil)

	body := []byte(`{"status":"approved","catatan":"ok"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/se/edit-requests/req-1/review", bytes.NewReader(body))
	req = withRoleAndUser(req, "admin", "")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.SEEditRequestResponse
	_ = json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, models.SEEditRequestApproved, resp.Status)
	svc.AssertExpectations(t)
}

func TestSEEditRequestHandler_ReviewForbiddenForNonAdmin(t *testing.T) {
	h := NewSEEditRequestHandler(new(mockSEEditRequestService))

	body := []byte(`{"status":"approved"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/se/edit-requests/req-1/review", bytes.NewReader(body))
	req = withRoleAndUser(req, "user", "user-1")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestSEEditRequestHandler_ReviewInvalidBody(t *testing.T) {
	h := NewSEEditRequestHandler(new(mockSEEditRequestService))

	req := httptest.NewRequest(http.MethodPut, "/api/se/edit-requests/req-1/review", bytes.NewBufferString("{invalid"))
	req = withRoleAndUser(req, "admin", "")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSEEditRequestHandler_HandleRequestEditSuccess(t *testing.T) {
	svc := new(mockSEEditRequestService)
	h := NewSEEditRequestHandler(svc)

	namaSE := "SE Baru"
	catatan := "tolong update"
	reqDTO := dto.CreateSEEditRequestDTO{
		Catatan: &catatan,
		DataPerubahan: dto.UpdateSERequest{
			NamaSE: &namaSE,
		},
	}

	svc.On("CreateRequest", "user-1", "se-1", reqDTO).Return(&dto.SEEditRequestResponse{
		ID:     "req-1",
		IDSE:   "se-1",
		IDUser: "user-1",
		Status: models.SEEditRequestPending,
	}, nil)

	body, _ := json.Marshal(reqDTO)
	req := httptest.NewRequest(http.MethodPost, "/api/se/se-1/request-edit", bytes.NewReader(body))
	req = withRoleAndUser(req, "user", "user-1")
	w := httptest.NewRecorder()

	h.HandleRequestEdit(w, req, "se-1")

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp dto.SEEditRequestResponse
	_ = json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "req-1", resp.ID)
	svc.AssertExpectations(t)
}

func TestSEEditRequestHandler_HandleRequestEditUnauthorized(t *testing.T) {
	h := NewSEEditRequestHandler(new(mockSEEditRequestService))

	req := httptest.NewRequest(http.MethodPost, "/api/se/se-1/request-edit", bytes.NewBufferString(`{}`))
	w := httptest.NewRecorder()

	h.HandleRequestEdit(w, req, "se-1")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSEEditRequestHandler_HandleRequestEditServiceError(t *testing.T) {
	svc := new(mockSEEditRequestService)
	h := NewSEEditRequestHandler(svc)

	reqDTO := dto.CreateSEEditRequestDTO{}
	svc.On("CreateRequest", "user-1", "se-1", reqDTO).Return(nil, errors.New("invalid request"))

	body, _ := json.Marshal(reqDTO)
	req := httptest.NewRequest(http.MethodPost, "/api/se/se-1/request-edit", bytes.NewReader(body))
	req = withRoleAndUser(req, "user", "user-1")
	w := httptest.NewRecorder()

	h.HandleRequestEdit(w, req, "se-1")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertExpectations(t)
}

func TestSEEditRequestHandler_ServeHTTP_NotFound(t *testing.T) {
	h := NewSEEditRequestHandler(new(mockSEEditRequestService))

	req := httptest.NewRequest(http.MethodDelete, "/api/se/edit-requests/req-1", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
