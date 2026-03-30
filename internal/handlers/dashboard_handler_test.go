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

func (m *MockDashboardService) GetSummary(ctx context.Context, f dto.DashboardFilter) (*dto.DashboardSummary, error) {
	args := m.Called(ctx, f)
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
	f := dto.DashboardFilter{}

	from := q.Get("from")
	to := q.Get("to")
	if from != "" && to != "" {
		if _, err := time.Parse("2006-01-02", from); err == nil {
			if _, err2 := time.Parse("2006-01-02", to); err2 == nil {
				f.From = &from
				f.To = &to
			}
		}
	}

	year := q.Get("year")
	if year != "" && reYear.MatchString(year) {
		f.Year = &year
	}

	quarter := q.Get("quarter")
	if quarter != "" && f.Year != nil && reQuarter.MatchString(quarter) {
		f.Quarter = &quarter
	}

	f.SubSektorID = ptrStr(q.Get("sub_sektor_id"))

	kategoriSE := q.Get("kategori_se")
	if kategoriSE != "" {
		if !validKategoriSE[kategoriSE] {
			writeError(w, http.StatusBadRequest, "kategori_se tidak valid, nilai yang diizinkan: Strategis, Tinggi, Rendah")
			return
		}
		f.KategoriSE = &kategoriSE
	}

	res, err := h.mockService.GetSummary(r.Context(), f)
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
 HELPERS
=====================================
*/

func emptyFilter() dto.DashboardFilter { return dto.DashboardFilter{} }

func emptySummary() *dto.DashboardSummary {
	return &dto.DashboardSummary{Sektor: []dto.SectorCount{}, SE: dto.SeAgg{}}
}

/*
=====================================
 TEST SUMMARY - SUCCESS CASES
=====================================
*/

func TestDashboardHandler_Summary_Success_NoFilter(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	expected := &dto.DashboardSummary{
		Sektor: []dto.SectorCount{
			{ID: "sektor-1", Nama: "ILMATE", Total: 100, ThisMonth: 10},
			{ID: "sektor-2", Nama: "IKFT", Total: 50, ThisMonth: 5},
		},
		SE:       dto.SeAgg{TotalSE: 75, ThisMonth: 8, Strategis: 30, Tinggi: 25, Rendah: 20},
		SEStatus: dto.SeStatusCount{TotalPerusahaan: 150, SudahMengisiKSE: 75, BelumMengisiKSE: 75},
	}

	mockService.On("GetSummary", mock.Anything, emptyFilter()).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary", nil)
	w := httptest.NewRecorder()
	handler.Summary(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.DashboardSummary
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response.Sektor, 2)
	assert.Equal(t, int64(75), response.SE.TotalSE)
	assert.Equal(t, int64(30), response.SE.Strategis)
	assert.Equal(t, int64(75), response.SEStatus.SudahMengisiKSE)

	mockService.AssertExpectations(t)
}

func TestDashboardHandler_Summary_Success_EmptyResults(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	mockService.On("GetSummary", mock.Anything, emptyFilter()).Return(emptySummary(), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary", nil)
	w := httptest.NewRecorder()
	handler.Summary(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDashboardHandler_Summary_WithYearFilter(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	year := "2024"
	f := dto.DashboardFilter{Year: &year}
	mockService.On("GetSummary", mock.Anything, f).Return(emptySummary(), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?year=2024", nil)
	w := httptest.NewRecorder()
	handler.Summary(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDashboardHandler_Summary_WithQuarterFilter(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	year := "2024"
	quarter := "3"
	f := dto.DashboardFilter{Year: &year, Quarter: &quarter}
	mockService.On("GetSummary", mock.Anything, f).Return(emptySummary(), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?year=2024&quarter=3", nil)
	w := httptest.NewRecorder()
	handler.Summary(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDashboardHandler_Summary_WithKategoriSE_Valid(t *testing.T) {
	for _, kategori := range []string{"Strategis", "Tinggi", "Rendah"} {
		t.Run(kategori, func(t *testing.T) {
			mockService := new(MockDashboardService)
			handler := createTestHandler(mockService)

			k := kategori
			f := dto.DashboardFilter{KategoriSE: &k}
			mockService.On("GetSummary", mock.Anything, f).Return(emptySummary(), nil)

			req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?kategori_se="+kategori, nil)
			w := httptest.NewRecorder()
			handler.Summary(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestDashboardHandler_Summary_WithKategoriSE_Invalid(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?kategori_se=Invalid", nil)
	w := httptest.NewRecorder()
	handler.Summary(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	assert.Contains(t, response["error"], "kategori_se tidak valid")

	mockService.AssertNotCalled(t, "GetSummary")
}

func TestDashboardHandler_Summary_QuarterIgnoredWithoutYear(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	// quarter tanpa year → quarter harus diabaikan, filter kosong
	mockService.On("GetSummary", mock.Anything, emptyFilter()).Return(emptySummary(), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?quarter=2", nil)
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

	mockService.On("GetSummary", mock.Anything, emptyFilter()).
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

	mockService.On("GetSummary", mock.Anything, emptyFilter()).Return(emptySummary(), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDashboardHandler_ServeHTTP_SummaryPathWithTrailingSlash(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	mockService.On("GetSummary", mock.Anything, emptyFilter()).Return(emptySummary(), nil)

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

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch} {
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
	writeJSON(w, http.StatusOK, map[string]string{"message": "test"})

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

func TestDashboardHandler_Summary_ContentTypeHeader(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	mockService.On("GetSummary", mock.Anything, emptyFilter()).Return(emptySummary(), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary", nil)
	w := httptest.NewRecorder()
	handler.Summary(w, req)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	mockService.AssertExpectations(t)
}

/*

/*
=====================================
 TEST FILTER KOMBINASI
=====================================
*/

func TestDashboardHandler_Summary_CombinedFilter_YearAndKategoriSE(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	year := "2025"
	kategori := "Tinggi"
	f := dto.DashboardFilter{Year: &year, KategoriSE: &kategori}
	mockService.On("GetSummary", mock.Anything, f).Return(&dto.DashboardSummary{
		SE: dto.SeAgg{TotalSE: 28, Tinggi: 28},
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?year=2025&kategori_se=Tinggi", nil)
	w := httptest.NewRecorder()
	handler.Summary(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.DashboardSummary
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, int64(28), response.SE.Tinggi)

	mockService.AssertExpectations(t)
}

func TestDashboardHandler_Summary_CombinedFilter_YearQuarterSubSektor(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	year := "2024"
	quarter := "1"
	subSektorID := "sub-uuid-abc"
	f := dto.DashboardFilter{Year: &year, Quarter: &quarter, SubSektorID: &subSektorID}
	mockService.On("GetSummary", mock.Anything, f).Return(emptySummary(), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?year=2024&quarter=1&sub_sektor_id=sub-uuid-abc", nil)
	w := httptest.NewRecorder()
	handler.Summary(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDashboardHandler_Summary_FromTo_TakesPriorityOverYear(t *testing.T) {
	// Jika from+to dan year keduanya dikirim, from+to harus menang di filter
	// (year tetap di-set tapi buildDateRange akan prioritaskan from+to)
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	from := "2024-06-01"
	to := "2024-06-30"
	year := "2024"
	f := dto.DashboardFilter{From: &from, To: &to, Year: &year}
	mockService.On("GetSummary", mock.Anything, f).Return(emptySummary(), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?from=2024-06-01&to=2024-06-30&year=2024", nil)
	w := httptest.NewRecorder()
	handler.Summary(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

/*
=====================================
 TEST VALIDASI FROM / TO EDGE CASES
=====================================
*/

func TestDashboardHandler_Summary_InvalidFromDate(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"only from without to", "/api/dashboard/summary?from=2024-01-01"},
		{"only to without from", "/api/dashboard/summary?to=2024-01-31"},
		{"from invalid format DD-MM-YYYY", "/api/dashboard/summary?from=01-01-2024&to=31-01-2024"},
		{"from invalid format YYYY/MM/DD", "/api/dashboard/summary?from=2024/01/01&to=2024/01/31"},
		{"from has invalid month 13", "/api/dashboard/summary?from=2024-13-01&to=2024-13-31"},
		{"from has invalid day 32", "/api/dashboard/summary?from=2024-01-32&to=2024-01-32"},
		{"empty string", "/api/dashboard/summary?from=&to="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockDashboardService)
			handler := createTestHandler(mockService)

			// Semua kasus di atas harus jatuh ke filter kosong (from/to nil)
			mockService.On("GetSummary", mock.Anything, emptyFilter()).Return(emptySummary(), nil)

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()
			handler.Summary(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "URL: %s", tt.url)
			mockService.AssertExpectations(t)
		})
	}
}

func TestDashboardHandler_Summary_InvalidYearFormat(t *testing.T) {
	tests := []struct {
		name string
		year string
	}{
		{"2 digits", "24"},
		{"3 digits", "202"},
		{"letters", "abcd"},
		{"alphanumeric", "20a4"},
		{"with slash", "2024/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockDashboardService)
			handler := createTestHandler(mockService)

			// year invalid → diabaikan → filter kosong
			mockService.On("GetSummary", mock.Anything, emptyFilter()).Return(emptySummary(), nil)

			req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?year="+tt.year, nil)
			w := httptest.NewRecorder()
			handler.Summary(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestDashboardHandler_Summary_InvalidQuarterFormat(t *testing.T) {
	tests := []struct {
		name    string
		quarter string
	}{
		{"quarter 0", "0"},
		{"quarter 5", "5"},
		{"quarter letter", "q"},
		{"quarter negative", "-1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockDashboardService)
			handler := createTestHandler(mockService)

			year := "2024"
			// quarter invalid → hanya year yang diset
			f := dto.DashboardFilter{Year: &year}
			mockService.On("GetSummary", mock.Anything, f).Return(emptySummary(), nil)

			req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?year=2024&quarter="+tt.quarter, nil)
			w := httptest.NewRecorder()
			handler.Summary(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

/*
=====================================
 TEST RESPONSE BODY STRUCTURE
=====================================
*/

func TestDashboardHandler_Summary_ResponseStructure_AllFields(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	expected := &dto.DashboardSummary{
		Sektor: []dto.SectorCount{
			{ID: "s-1", Nama: "ILMATE", Total: 97, ThisMonth: 8},
			{ID: "s-2", Nama: "Industri Agro", Total: 91, ThisMonth: 10},
			{ID: "s-3", Nama: "IKFT", Total: 63, ThisMonth: 6},
		},
		SE: dto.SeAgg{
			TotalSE:   77,
			ThisMonth: 8,
			Strategis: 30,
			Tinggi:    28,
			Rendah:    19,
		},
		SEStatus: dto.SeStatusCount{
			TotalPerusahaan: 251,
			SudahMengisiKSE: 77,
			BelumMengisiKSE: 174,
		},
	}

	mockService.On("GetSummary", mock.Anything, emptyFilter()).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary", nil)
	w := httptest.NewRecorder()
	handler.Summary(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.DashboardSummary
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	// Verifikasi struktur sektor
	assert.Len(t, response.Sektor, 3)
	assert.Equal(t, "ILMATE", response.Sektor[0].Nama)
	assert.Equal(t, int64(97), response.Sektor[0].Total)
	assert.Equal(t, int64(8), response.Sektor[0].ThisMonth)

	// Verifikasi KSE breakdown
	assert.Equal(t, int64(77), response.SE.TotalSE)
	assert.Equal(t, int64(8), response.SE.ThisMonth)
	assert.Equal(t, int64(30), response.SE.Strategis)
	assert.Equal(t, int64(28), response.SE.Tinggi)
	assert.Equal(t, int64(19), response.SE.Rendah)
	// Strategis+Tinggi+Rendah harus sama dengan TotalSE
	assert.Equal(t, response.SE.TotalSE, response.SE.Strategis+response.SE.Tinggi+response.SE.Rendah)

	// Verifikasi status pengisian
	assert.Equal(t, int64(251), response.SEStatus.TotalPerusahaan)
	assert.Equal(t, int64(77), response.SEStatus.SudahMengisiKSE)
	assert.Equal(t, int64(174), response.SEStatus.BelumMengisiKSE)
	// Sudah+Belum harus sama dengan total
	assert.Equal(t, response.SEStatus.TotalPerusahaan, response.SEStatus.SudahMengisiKSE+response.SEStatus.BelumMengisiKSE)

	mockService.AssertExpectations(t)
}

func TestDashboardHandler_Summary_ManySektor(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	sectors := make([]dto.SectorCount, 10)
	totalSE := int64(0)
	for i := 0; i < 10; i++ {
		sectors[i] = dto.SectorCount{
			ID:        "sektor-" + string(rune('A'+i)),
			Nama:      "Sektor " + string(rune('A'+i)),
			Total:     int64(100 * (i + 1)),
			ThisMonth: int64(10 * (i + 1)),
		}
		totalSE += int64(10 * (i + 1))
	}

	expected := &dto.DashboardSummary{
		Sektor:   sectors,
		SE:       dto.SeAgg{TotalSE: totalSE},
		SEStatus: dto.SeStatusCount{TotalPerusahaan: 550, SudahMengisiKSE: totalSE},
	}

	mockService.On("GetSummary", mock.Anything, emptyFilter()).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary", nil)
	w := httptest.NewRecorder()
	handler.Summary(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.DashboardSummary
	json.NewDecoder(w.Body).Decode(&response)
	assert.Len(t, response.Sektor, 10)
	assert.Equal(t, totalSE, response.SE.TotalSE)

	mockService.AssertExpectations(t)
}

/*
=====================================
 TEST FILTER FROM / TO — cabang yang belum dicakup
=====================================
*/

// from+to keduanya valid TANPA year → From & To diset, Year nil
func TestDashboardHandler_Summary_FromTo_Valid_WithoutYear(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	from := "2024-03-01"
	to := "2024-03-31"
	f := dto.DashboardFilter{From: &from, To: &to}
	mockService.On("GetSummary", mock.Anything, f).Return(emptySummary(), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?from=2024-03-01&to=2024-03-31", nil)
	w := httptest.NewRecorder()
	handler.Summary(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// from valid tapi to invalid → keduanya diabaikan (filter kosong)
func TestDashboardHandler_Summary_FromValid_ToInvalid_BothIgnored(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"to format salah DD-MM-YYYY", "/api/dashboard/summary?from=2024-01-01&to=31-01-2024"},
		{"to bukan tanggal", "/api/dashboard/summary?from=2024-01-01&to=not-a-date"},
		{"to bulan tidak valid", "/api/dashboard/summary?from=2024-01-01&to=2024-13-01"},
		{"to hari tidak valid", "/api/dashboard/summary?from=2024-01-01&to=2024-01-32"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockDashboardService)
			handler := createTestHandler(mockService)

			// from+to tidak keduanya valid → filter kosong (From & To nil)
			mockService.On("GetSummary", mock.Anything, emptyFilter()).Return(emptySummary(), nil)

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()
			handler.Summary(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "URL: %s", tt.url)
			mockService.AssertExpectations(t)
		})
	}
}

// from+to valid + year + kategori_se → semua filter aktif bersamaan
func TestDashboardHandler_Summary_FromTo_WithKategoriSE(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	from := "2025-01-01"
	to := "2025-03-31"
	year := "2025"
	kategori := "Strategis"
	f := dto.DashboardFilter{From: &from, To: &to, Year: &year, KategoriSE: &kategori}
	mockService.On("GetSummary", mock.Anything, f).Return(emptySummary(), nil)

	req := httptest.NewRequest(http.MethodGet,
		"/api/dashboard/summary?from=2025-01-01&to=2025-03-31&year=2025&kategori_se=Strategis", nil)
	w := httptest.NewRecorder()
	handler.Summary(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

/*
=====================================
 TEST FILTER SUB_SEKTOR_ID — cabang standalone
=====================================
*/

// sub_sektor_id saja tanpa filter lain
func TestDashboardHandler_Summary_SubSektorID_Only(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	subID := "sub-sektor-uuid-99"
	f := dto.DashboardFilter{SubSektorID: &subID}
	mockService.On("GetSummary", mock.Anything, f).Return(emptySummary(), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?sub_sektor_id=sub-sektor-uuid-99", nil)
	w := httptest.NewRecorder()
	handler.Summary(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// sub_sektor_id kosong ("") → ptrStr mengembalikan nil → filter kosong
func TestDashboardHandler_Summary_SubSektorID_Empty_TreatedAsNil(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	mockService.On("GetSummary", mock.Anything, emptyFilter()).Return(emptySummary(), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?sub_sektor_id=", nil)
	w := httptest.NewRecorder()
	handler.Summary(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

/*
=====================================
 TEST VALIDASI QUARTER — semua nilai valid (1,2,4)
=====================================
*/

// Semua nilai quarter valid (1–4) harus diterima
func TestDashboardHandler_Summary_AllValidQuarters(t *testing.T) {
	for _, q := range []string{"1", "2", "4"} {
		t.Run("quarter="+q, func(t *testing.T) {
			mockService := new(MockDashboardService)
			handler := createTestHandler(mockService)

			year := "2025"
			quarter := q
			f := dto.DashboardFilter{Year: &year, Quarter: &quarter}
			mockService.On("GetSummary", mock.Anything, f).Return(emptySummary(), nil)

			req := httptest.NewRequest(http.MethodGet,
				"/api/dashboard/summary?year=2025&quarter="+q, nil)
			w := httptest.NewRecorder()
			handler.Summary(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

/*
=====================================
 TEST ptrStr HELPER
=====================================
*/

func TestPtrStr_EmptyString_ReturnsNil(t *testing.T) {
	result := ptrStr("")
	assert.Nil(t, result)
}

func TestPtrStr_NonEmptyString_ReturnsPointer(t *testing.T) {
	result := ptrStr("hello")
	assert.NotNil(t, result)
	assert.Equal(t, "hello", *result)
}

/*
=====================================
 TEST NewDashboardHandler CONSTRUCTOR
=====================================
*/

func TestNewDashboardHandler_ReturnsNonNil(t *testing.T) {
	h := NewDashboardHandler(nil)
	assert.NotNil(t, h)
}