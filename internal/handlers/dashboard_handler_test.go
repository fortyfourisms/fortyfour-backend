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

func emptyFilter() dto.DashboardFilter { return dto.DashboardFilter{} }

func emptySummary() *dto.DashboardSummary {
	return &dto.DashboardSummary{
		Sektor:     []dto.SectorCount{},
		Ikas:       dto.IkasAgg{},
		SE:         dto.SeAgg{},
		SEStatus:   dto.SeStatusCount{},
		IkasStatus: dto.IkasStatusCount{},
	}
}

func TestDashboardHandler_Summary_Success_NoFilter(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	expected := &dto.DashboardSummary{
		Sektor: []dto.SectorCount{
			{ID: "sektor-1", Nama: "ILMATE", Total: 100, ThisMonth: 10},
			{ID: "sektor-2", Nama: "IKFT", Total: 50, ThisMonth: 5},
		},
		Ikas:       dto.IkasAgg{Total: 45, AvgNilaiKematangan: 2.7, AvgTargetNilai: 4},
		SE:         dto.SeAgg{TotalSE: 75, ThisMonth: 8, Strategis: 30, Tinggi: 25, Rendah: 20},
		SEStatus:   dto.SeStatusCount{TotalPerusahaan: 150, SudahMengisiKSE: 75, BelumMengisiKSE: 75},
		IkasStatus: dto.IkasStatusCount{TotalPerusahaan: 150, SudahMengisiIKAS: 45, BelumMengisiIKAS: 105},
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
	assert.Equal(t, int64(45), response.Ikas.Total)
	assert.Equal(t, int64(75), response.SE.TotalSE)
	assert.Equal(t, int64(75), response.SEStatus.SudahMengisiKSE)
	assert.Equal(t, int64(45), response.IkasStatus.SudahMengisiIKAS)
	mockService.AssertExpectations(t)
}

func TestDashboardHandler_Summary_WithFilter(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	year := "2024"
	quarter := "3"
	subSektorID := "sub-uuid-abc"
	kategori := "Tinggi"
	f := dto.DashboardFilter{Year: &year, Quarter: &quarter, SubSektorID: &subSektorID, KategoriSE: &kategori}
	mockService.On("GetSummary", mock.Anything, f).Return(emptySummary(), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?year=2024&quarter=3&sub_sektor_id=sub-uuid-abc&kategori_se=Tinggi", nil)
	w := httptest.NewRecorder()
	handler.Summary(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
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

func TestDashboardHandler_Summary_ServiceError(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	mockService.On("GetSummary", mock.Anything, emptyFilter()).Return(nil, errors.New("database connection failed"))

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary", nil)
	w := httptest.NewRecorder()
	handler.Summary(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	assert.Contains(t, response["error"], "database connection failed")
	mockService.AssertExpectations(t)
}

func TestDashboardHandler_ServeHTTP_MethodNotAllowed(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/api/dashboard/summary", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	mockService.AssertNotCalled(t, "GetSummary")
}

func TestDashboardHandler_Summary_ResponseStructure_AllFields(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := createTestHandler(mockService)

	expected := &dto.DashboardSummary{
		Sektor: []dto.SectorCount{
			{ID: "s-1", Nama: "ILMATE", Total: 97, ThisMonth: 8},
			{ID: "s-2", Nama: "Industri Agro", Total: 91, ThisMonth: 10},
		},
		Ikas: dto.IkasAgg{
			Total:              80,
			AvgNilaiKematangan: 2.65,
			AvgTargetNilai:     4.00,
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
		IkasStatus: dto.IkasStatusCount{
			TotalPerusahaan:  251,
			SudahMengisiIKAS: 80,
			BelumMengisiIKAS: 171,
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
	assert.Equal(t, int64(80), response.Ikas.Total)
	assert.Equal(t, 2.65, response.Ikas.AvgNilaiKematangan)
	assert.Equal(t, int64(171), response.IkasStatus.BelumMengisiIKAS)
	assert.Equal(t, response.SEStatus.TotalPerusahaan, response.SEStatus.SudahMengisiKSE+response.SEStatus.BelumMengisiKSE)
	assert.Equal(t, response.IkasStatus.TotalPerusahaan, response.IkasStatus.SudahMengisiIKAS+response.IkasStatus.BelumMengisiIKAS)
	mockService.AssertExpectations(t)
}

func TestPtrStr_EmptyString_ReturnsNil(t *testing.T) {
	result := ptrStr("")
	assert.Nil(t, result)
}

func TestPtrStr_NonEmptyString_ReturnsPointer(t *testing.T) {
	result := ptrStr("hello")
	assert.NotNil(t, result)
	assert.Equal(t, "hello", *result)
}

func TestNewDashboardHandler_ReturnsNonNil(t *testing.T) {
	h := NewDashboardHandler(nil)
	assert.NotNil(t, h)
}
