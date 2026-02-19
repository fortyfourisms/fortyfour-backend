package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"fortyfour-backend/internal/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/*
=====================================
 MOCK DASHBOARD SERVICE
=====================================
*/

type MockDashboardService struct {
	mock.Mock
}

func (m *MockDashboardService) GetSummary(ctx context.Context, from, to *string) (*dto.DashboardSummary, error) {
	args := m.Called(ctx, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DashboardSummary), args.Error(1)
}

/*
=====================================
 TEST DASHBOARD HANDLER WRAPPER
=====================================
*/

type testDashboardHandler struct {
	mockService *MockDashboardService
}

func createTestHandler(mockService *MockDashboardService) *testDashboardHandler {
	return &testDashboardHandler{mockService: mockService}
}

func (h *testDashboardHandler) Summary(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	from := q.Get("from")
	to := q.Get("to")
	var fromPtr, toPtr *string
	if from != "" && to != "" {
		// validate format YYYY-MM-DD, ignore if invalid
		if _, err := time.Parse("2006-01-02", from); err == nil {
			if _, err2 := time.Parse("2006-01-02", to); err2 == nil {
				fromPtr = &from
				toPtr = &to
			}
		}
	}
	res, err := h.mockService.GetSummary(r.Context(), fromPtr, toPtr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *testDashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	path := r.URL.Path
	if path == "/api/dashboard/summary" || path == "/api/dashboard/summary/" {
		h.Summary(w, r)
		return
	}

	http.NotFound(w, r)
}

/*
=====================================
 TEST SUMMARY - SUCCESS CASES
=====================================
*/

func TestDashboardHandler_Summary_Success_NoDateFilter(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	expectedSummary := &dto.DashboardSummary{
		Sektor: []dto.SectorCount{
			{ID: "sektor-1", Nama: "ILMATE", Total: 100, ThisMonth: 10},
			{ID: "sektor-2", Nama: "IKFT", Total: 50, ThisMonth: 5},
		},
		SE: dto.SeAgg{TotalSE: 75},
	}

	mockService.On("GetSummary", mock.Anything, (*string)(nil), (*string)(nil)).
		Return(expectedSummary, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary", nil)
	w := httptest.NewRecorder()

	handler.Summary(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.DashboardSummary
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Len(t, response.Sektor, 2)
	assert.Equal(t, "ILMATE", response.Sektor[0].Nama)
	assert.Equal(t, int64(100), response.Sektor[0].Total)
	assert.Equal(t, int64(75), response.SE.TotalSE)

	mockService.AssertExpectations(t)
}

func TestDashboardHandler_Summary_Success_WithValidDateFilter(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	from := "2024-01-01"
	to := "2024-01-31"

	expectedSummary := &dto.DashboardSummary{
		Sektor: []dto.SectorCount{
			{ID: "sektor-1", Nama: "ILMATE", Total: 100, ThisMonth: 15},
		},
		SE: dto.SeAgg{TotalSE: 50},
	}

	mockService.On("GetSummary", mock.Anything, &from, &to).
		Return(expectedSummary, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?from=2024-01-01&to=2024-01-31", nil)
	w := httptest.NewRecorder()

	handler.Summary(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.DashboardSummary
	json.NewDecoder(w.Body).Decode(&response)

	assert.Len(t, response.Sektor, 1)
	assert.Equal(t, int64(15), response.Sektor[0].ThisMonth)

	mockService.AssertExpectations(t)
}

func TestDashboardHandler_Summary_Success_EmptyResults(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	expectedSummary := &dto.DashboardSummary{
		Sektor: []dto.SectorCount{},
		SE:     dto.SeAgg{TotalSE: 0},
	}

	mockService.On("GetSummary", mock.Anything, (*string)(nil), (*string)(nil)).
		Return(expectedSummary, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary", nil)
	w := httptest.NewRecorder()

	handler.Summary(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.DashboardSummary
	json.NewDecoder(w.Body).Decode(&response)

	assert.Len(t, response.Sektor, 0)
	assert.Equal(t, int64(0), response.SE.TotalSE)

	mockService.AssertExpectations(t)
}

/*
=====================================
 TEST SUMMARY - DATE VALIDATION
=====================================
*/

func TestDashboardHandler_Summary_InvalidDateFormat_Ignored(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	expectedSummary := &dto.DashboardSummary{
		Sektor: []dto.SectorCount{
			{ID: "sektor-1", Nama: "ILMATE", Total: 50, ThisMonth: 5},
		},
		SE: dto.SeAgg{TotalSE: 25},
	}

	// Invalid date format should be ignored, passed as nil
	mockService.On("GetSummary", mock.Anything, (*string)(nil), (*string)(nil)).
		Return(expectedSummary, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?from=invalid&to=date", nil)
	w := httptest.NewRecorder()

	handler.Summary(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestDashboardHandler_Summary_OnlyFromDate_Ignored(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	expectedSummary := &dto.DashboardSummary{
		Sektor: []dto.SectorCount{},
		SE:     dto.SeAgg{TotalSE: 0},
	}

	mockService.On("GetSummary", mock.Anything, (*string)(nil), (*string)(nil)).
		Return(expectedSummary, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?from=2024-01-01", nil)
	w := httptest.NewRecorder()

	handler.Summary(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestDashboardHandler_Summary_OnlyToDate_Ignored(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	expectedSummary := &dto.DashboardSummary{
		Sektor: []dto.SectorCount{},
		SE:     dto.SeAgg{TotalSE: 0},
	}

	mockService.On("GetSummary", mock.Anything, (*string)(nil), (*string)(nil)).
		Return(expectedSummary, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?to=2024-01-31", nil)
	w := httptest.NewRecorder()

	handler.Summary(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestDashboardHandler_Summary_PartiallyInvalidDate(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	expectedSummary := &dto.DashboardSummary{
		Sektor: []dto.SectorCount{},
		SE:     dto.SeAgg{TotalSE: 0},
	}

	mockService.On("GetSummary", mock.Anything, (*string)(nil), (*string)(nil)).
		Return(expectedSummary, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?from=2024-01-01&to=invalid", nil)
	w := httptest.NewRecorder()

	handler.Summary(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

/*
=====================================
 TEST SUMMARY - ERROR CASES
=====================================
*/

func TestDashboardHandler_Summary_ServiceError(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	mockService.On("GetSummary", mock.Anything, (*string)(nil), (*string)(nil)).
		Return(nil, errors.New("database connection failed"))

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary", nil)
	w := httptest.NewRecorder()

	handler.Summary(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	assert.Contains(t, response["error"], "database connection failed")

	mockService.AssertExpectations(t)
}

/*
=====================================
 TEST SERVE HTTP - ROUTING
=====================================
*/

func TestDashboardHandler_ServeHTTP_SummaryPath(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	expectedSummary := &dto.DashboardSummary{
		Sektor: []dto.SectorCount{},
		SE:     dto.SeAgg{TotalSE: 0},
	}

	mockService.On("GetSummary", mock.Anything, (*string)(nil), (*string)(nil)).
		Return(expectedSummary, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestDashboardHandler_ServeHTTP_SummaryPathWithTrailingSlash(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	expectedSummary := &dto.DashboardSummary{
		Sektor: []dto.SectorCount{},
		SE:     dto.SeAgg{TotalSE: 0},
	}

	mockService.On("GetSummary", mock.Anything, (*string)(nil), (*string)(nil)).
		Return(expectedSummary, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestDashboardHandler_ServeHTTP_InvalidPath_NotFound(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/invalid", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertNotCalled(t, "GetSummary")
}

func TestDashboardHandler_ServeHTTP_MethodNotAllowed(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/dashboard/summary", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, http.StatusMethodNotAllowed, w.Code)

			var response map[string]string
			json.NewDecoder(w.Body).Decode(&response)
			assert.Equal(t, "method not allowed", response["error"])
		})
	}

	mockService.AssertNotCalled(t, "GetSummary")
}

/*
=====================================
 TEST HELPER FUNCTIONS
=====================================
*/

func TestDashboardHandler_WriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"message": "test"}

	writeJSON(w, http.StatusOK, data)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, "test", response["message"])
}

func TestDashboardHandler_WriteError(t *testing.T) {
	w := httptest.NewRecorder()

	writeError(w, http.StatusBadRequest, "validation failed")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, "validation failed", response["error"])
}

/*
=====================================
 TEST EDGE CASES
=====================================
*/

func TestDashboardHandler_Summary_ManySektor(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	sectors := make([]dto.SectorCount, 10)
	for i := 0; i < 10; i++ {
		sectors[i] = dto.SectorCount{
			ID:        "sektor-" + string(rune(i)),
			Nama:      "Sektor " + string(rune(i)),
			Total:     int64(100 * (i + 1)),
			ThisMonth: int64(10 * (i + 1)),
		}
	}

	expectedSummary := &dto.DashboardSummary{
		Sektor: sectors,
		SE:     dto.SeAgg{TotalSE: 5500},
	}

	mockService.On("GetSummary", mock.Anything, (*string)(nil), (*string)(nil)).
		Return(expectedSummary, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary", nil)
	w := httptest.NewRecorder()

	handler.Summary(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.DashboardSummary
	json.NewDecoder(w.Body).Decode(&response)

	assert.Len(t, response.Sektor, 10)
	assert.Equal(t, int64(5500), response.SE.TotalSE)

	mockService.AssertExpectations(t)
}

func TestDashboardHandler_Summary_DateFormats(t *testing.T) {
	tests := []struct {
		name        string
		from        string
		to          string
		shouldParse bool
	}{
		{
			name:        "Valid YYYY-MM-DD",
			from:        "2024-01-01",
			to:          "2024-01-31",
			shouldParse: true,
		},
		{
			name:        "Invalid format DD-MM-YYYY",
			from:        "01-01-2024",
			to:          "31-01-2024",
			shouldParse: false,
		},
		{
			name:        "Invalid format YYYY/MM/DD",
			from:        "2024/01/01",
			to:          "2024/01/31",
			shouldParse: false,
		},
		{
			name:        "Invalid month",
			from:        "2024-13-01",
			to:          "2024-13-31",
			shouldParse: false,
		},
		{
			name:        "Invalid day",
			from:        "2024-01-32",
			to:          "2024-01-40",
			shouldParse: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockDashboardService)
			handler := createTestHandler(mockService)

			expectedSummary := &dto.DashboardSummary{
				Sektor: []dto.SectorCount{},
				SE:     dto.SeAgg{TotalSE: 0},
			}

			if tt.shouldParse {
				mockService.On("GetSummary", mock.Anything, &tt.from, &tt.to).
					Return(expectedSummary, nil)
			} else {
				// Invalid dates should be ignored and passed as nil
				mockService.On("GetSummary", mock.Anything, (*string)(nil), (*string)(nil)).
					Return(expectedSummary, nil)
			}

			req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?from="+tt.from+"&to="+tt.to, nil)
			w := httptest.NewRecorder()

			handler.Summary(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}

/*
=====================================
 TEST CONTENT TYPE HEADERS
=====================================
*/

func TestDashboardHandler_Summary_ContentTypeHeader(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	expectedSummary := &dto.DashboardSummary{
		Sektor: []dto.SectorCount{},
		SE:     dto.SeAgg{TotalSE: 0},
	}

	mockService.On("GetSummary", mock.Anything, (*string)(nil), (*string)(nil)).
		Return(expectedSummary, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary", nil)
	w := httptest.NewRecorder()

	handler.Summary(w, req)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	mockService.AssertExpectations(t)
}
