package services

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"fortyfour-backend/internal/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/*
=====================================
 MOCK DASHBOARD REPOSITORY
=====================================
*/

type MockDashboardRepository struct {
	mock.Mock
}

func (m *MockDashboardRepository) CountPerSektor(ctx context.Context, f dto.DashboardFilter) ([]dto.SectorCount, error) {
	args := m.Called(ctx, f)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.SectorCount), args.Error(1)
}

func (m *MockDashboardRepository) SeGlobalAgg(ctx context.Context, f dto.DashboardFilter) (dto.SeAgg, error) {
	args := m.Called(ctx, f)
	return args.Get(0).(dto.SeAgg), args.Error(1)
}

func (m *MockDashboardRepository) SeStatusCount(ctx context.Context, f dto.DashboardFilter) (dto.SeStatusCount, error) {
	args := m.Called(ctx, f)
	return args.Get(0).(dto.SeStatusCount), args.Error(1)
}

// TODO: re-enable ikas status when ikas table is ready
// func (m *MockDashboardRepository) IkasStatusCount(ctx context.Context, f dto.DashboardFilter) (dto.IkasStatusCount, error) {
// 	args := m.Called(ctx, f)
// 	return args.Get(0).(dto.IkasStatusCount), args.Error(1)
// }

/*
=====================================
 HELPER FUNCTIONS
=====================================
*/

func createServiceWithMockRepo(mockRepo *MockDashboardRepository) *DashboardService {
	return NewDashboardService(mockRepo, nil)
}

/*
=====================================
 TEST GET SUMMARY - SUCCESS CASES
=====================================
*/

func TestDashboardService_GetSummary_Success_NoFilter(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	expectedSectors := []dto.SectorCount{
		{ID: "sektor-1", Nama: "ILMATE", Total: 100, ThisMonth: 10},
		{ID: "sektor-2", Nama: "IKFT", Total: 50, ThisMonth: 5},
	}
	expectedSE := dto.SeAgg{TotalSE: 75, ThisMonth: 8, Strategis: 30, Tinggi: 25, Rendah: 20}
	expectedStatus := dto.SeStatusCount{TotalPerusahaan: 150, SudahMengisiKSE: 75, BelumMengisiKSE: 75}

	mockRepo.On("CountPerSektor", ctx, f).Return(expectedSectors, nil)
	mockRepo.On("SeGlobalAgg", ctx, f).Return(expectedSE, nil)
	mockRepo.On("SeStatusCount", ctx, f).Return(expectedStatus, nil)

	result, err := service.GetSummary(ctx, f)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Sektor, 2)
	assert.Equal(t, int64(75), result.SE.TotalSE)
	assert.Equal(t, int64(8), result.SE.ThisMonth)
	assert.Equal(t, int64(30), result.SE.Strategis)
	assert.Equal(t, int64(75), result.SEStatus.SudahMengisiKSE)
	assert.Equal(t, int64(75), result.SEStatus.BelumMengisiKSE)

	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummary_Success_WithDateFilter(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()
	from := "2024-01-01"
	to := "2024-01-31"
	f := dto.DashboardFilter{From: &from, To: &to}

	mockRepo.On("CountPerSektor", ctx, f).Return([]dto.SectorCount{
		{ID: "sektor-1", Nama: "ILMATE", Total: 100, ThisMonth: 15},
	}, nil)
	mockRepo.On("SeGlobalAgg", ctx, f).Return(dto.SeAgg{TotalSE: 50, ThisMonth: 10}, nil)
	mockRepo.On("SeStatusCount", ctx, f).Return(dto.SeStatusCount{TotalPerusahaan: 100, SudahMengisiKSE: 50, BelumMengisiKSE: 50}, nil)

	result, err := service.GetSummary(ctx, f)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(15), result.Sektor[0].ThisMonth)
	assert.Equal(t, int64(10), result.SE.ThisMonth)

	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummary_Success_WithYearFilter(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()
	year := "2024"
	f := dto.DashboardFilter{Year: &year}

	mockRepo.On("CountPerSektor", ctx, f).Return([]dto.SectorCount{}, nil)
	mockRepo.On("SeGlobalAgg", ctx, f).Return(dto.SeAgg{TotalSE: 120}, nil)
	mockRepo.On("SeStatusCount", ctx, f).Return(dto.SeStatusCount{}, nil)

	result, err := service.GetSummary(ctx, f)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(120), result.SE.TotalSE)

	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummary_Success_WithQuarterFilter(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()
	year := "2024"
	quarter := "2"
	f := dto.DashboardFilter{Year: &year, Quarter: &quarter}

	mockRepo.On("CountPerSektor", ctx, f).Return([]dto.SectorCount{}, nil)
	mockRepo.On("SeGlobalAgg", ctx, f).Return(dto.SeAgg{TotalSE: 40}, nil)
	mockRepo.On("SeStatusCount", ctx, f).Return(dto.SeStatusCount{}, nil)

	result, err := service.GetSummary(ctx, f)

	assert.NoError(t, err)
	assert.Equal(t, int64(40), result.SE.TotalSE)

	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummary_Success_WithKategoriSEFilter(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()
	kategori := "Strategis"
	f := dto.DashboardFilter{KategoriSE: &kategori}

	mockRepo.On("CountPerSektor", ctx, f).Return([]dto.SectorCount{}, nil)
	mockRepo.On("SeGlobalAgg", ctx, f).Return(dto.SeAgg{TotalSE: 30, Strategis: 30}, nil)
	mockRepo.On("SeStatusCount", ctx, f).Return(dto.SeStatusCount{}, nil)

	result, err := service.GetSummary(ctx, f)

	assert.NoError(t, err)
	assert.Equal(t, int64(30), result.SE.Strategis)

	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummary_Success_EmptyResults(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	mockRepo.On("CountPerSektor", ctx, f).Return([]dto.SectorCount{}, nil)
	mockRepo.On("SeGlobalAgg", ctx, f).Return(dto.SeAgg{}, nil)
	mockRepo.On("SeStatusCount", ctx, f).Return(dto.SeStatusCount{}, nil)

	result, err := service.GetSummary(ctx, f)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Sektor, 0)
	assert.Equal(t, int64(0), result.SE.TotalSE)
	assert.Equal(t, int64(0), result.SEStatus.TotalPerusahaan)

	mockRepo.AssertExpectations(t)
}

/*
=====================================
 TEST GET SUMMARY - ERROR CASES
=====================================
*/

func TestDashboardService_GetSummary_CountPerSektorError(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	mockRepo.On("CountPerSektor", ctx, f).Return(nil, errors.New("database connection failed"))

	result, err := service.GetSummary(ctx, f)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database connection failed", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummary_SeGlobalAggError(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	mockRepo.On("CountPerSektor", ctx, f).Return([]dto.SectorCount{
		{ID: "sektor-1", Nama: "ILMATE", Total: 100, ThisMonth: 10},
	}, nil)
	mockRepo.On("SeGlobalAgg", ctx, f).Return(dto.SeAgg{}, errors.New("se table not found"))

	result, err := service.GetSummary(ctx, f)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "se table not found", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummary_SeStatusCountError(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	mockRepo.On("CountPerSektor", ctx, f).Return([]dto.SectorCount{}, nil)
	mockRepo.On("SeGlobalAgg", ctx, f).Return(dto.SeAgg{TotalSE: 10}, nil)
	mockRepo.On("SeStatusCount", ctx, f).Return(dto.SeStatusCount{}, errors.New("query failed"))

	result, err := service.GetSummary(ctx, f)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

/*
=====================================
 TEST NEW DASHBOARD SERVICE
=====================================
*/

func TestNewDashboardService(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := NewDashboardService(mockRepo, nil)
	assert.NotNil(t, service)
	assert.NotNil(t, service.repo)
}

/*
=====================================
 TEST CONTEXT CANCELLATION
=====================================
*/

func TestDashboardService_GetSummary_ContextCancellation(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	f := dto.DashboardFilter{}

	mockRepo.On("CountPerSektor", mock.Anything, f).
		Return(nil, context.Canceled).Maybe()

	result, err := service.GetSummary(ctx, f)

	assert.Error(t, err)
	assert.Nil(t, result)
}

/*
=====================================
 TEST CACHE — GetSummary
=====================================
*/

func TestDashboardService_GetSummary_CacheHit_SkipRepo(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	rc := newDashboardTestRedis()
	service := NewDashboardService(mockRepo, rc)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	// Pre-populate cache dengan key yang sama yang digunakan service
	cached := dto.DashboardSummary{
		Sektor: []dto.SectorCount{{ID: "cache-1", Nama: "Dari Cache", Total: 99}},
		SE:     dto.SeAgg{TotalSE: 42},
	}
	key := buildCacheKey(f)
	setDashboardCache(rc, key, cached)

	result, err := service.GetSummary(ctx, f)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Dari Cache", result.Sektor[0].Nama)
	assert.Equal(t, int64(42), result.SE.TotalSE)

	mockRepo.AssertNotCalled(t, "CountPerSektor")
	mockRepo.AssertNotCalled(t, "SeGlobalAgg")
	mockRepo.AssertNotCalled(t, "SeStatusCount")
}

func TestDashboardService_GetSummary_CacheMiss_SetsCache(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	rc := newDashboardTestRedis()
	service := NewDashboardService(mockRepo, rc)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	mockRepo.On("CountPerSektor", ctx, f).Return([]dto.SectorCount{{ID: "s1", Nama: "ILMATE", Total: 10}}, nil)
	mockRepo.On("SeGlobalAgg", ctx, f).Return(dto.SeAgg{TotalSE: 10}, nil)
	mockRepo.On("SeStatusCount", ctx, f).Return(dto.SeStatusCount{}, nil)

	_, err := service.GetSummary(ctx, f)
	assert.NoError(t, err)

	exists, _ := rc.Exists(buildCacheKey(f))
	assert.True(t, exists, "hasil harus di-cache setelah GetSummary")

	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummary_CacheKey_BerdasarkanFilter(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	rc := newDashboardTestRedis()
	service := NewDashboardService(mockRepo, rc)
	ctx := context.Background()

	from := "2024-01-01"
	to := "2024-01-31"
	f := dto.DashboardFilter{From: &from, To: &to}

	mockRepo.On("CountPerSektor", ctx, f).Return([]dto.SectorCount{{ID: "s1", Total: 5}}, nil)
	mockRepo.On("SeGlobalAgg", ctx, f).Return(dto.SeAgg{TotalSE: 5}, nil)
	mockRepo.On("SeStatusCount", ctx, f).Return(dto.SeStatusCount{}, nil)

	_, err := service.GetSummary(ctx, f)
	assert.NoError(t, err)

	keyWithFilter := buildCacheKey(f)
	keyNoFilter := buildCacheKey(dto.DashboardFilter{})

	existsWith, _ := rc.Exists(keyWithFilter)
	existsWithout, _ := rc.Exists(keyNoFilter)

	assert.True(t, existsWith, "cache harus ada untuk key dengan filter tanggal")
	assert.False(t, existsWithout, "cache untuk key tanpa filter tidak boleh ada")

	mockRepo.AssertExpectations(t)
}

/*
=====================================
 HELPERS REDIS UNTUK DASHBOARD TEST
=====================================
*/

func newDashboardTestRedis() *dashboardTestRedis {
	return &dashboardTestRedis{data: make(map[string]string)}
}

func setDashboardCache(rc *dashboardTestRedis, key string, value interface{}) {
	b, _ := json.Marshal(value)
	rc.data[key] = string(b)
}

type dashboardTestRedis struct {
	data map[string]string
}

func (r *dashboardTestRedis) Set(key string, value interface{}, ttl time.Duration) error {
	if v, ok := value.(string); ok {
		r.data[key] = v
	}
	return nil
}
func (r *dashboardTestRedis) Get(key string) (string, error) {
	v, ok := r.data[key]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}
func (r *dashboardTestRedis) Delete(key string) error { delete(r.data, key); return nil }
func (r *dashboardTestRedis) Exists(key string) (bool, error) {
	_, ok := r.data[key]
	return ok, nil
}
func (r *dashboardTestRedis) Scan(pattern string) ([]string, error) { return nil, nil }
func (r *dashboardTestRedis) Close() error                          { return nil }

/*

/*
=====================================
 TEST FILTER KOMBINASI
=====================================
*/

func TestDashboardService_GetSummary_CombinedFilter_YearQuarter(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()

	year := "2024"
	quarter := "4"
	f := dto.DashboardFilter{Year: &year, Quarter: &quarter}

	mockRepo.On("CountPerSektor", ctx, f).Return([]dto.SectorCount{
		{ID: "s1", Nama: "IKFT", Total: 63, ThisMonth: 15},
	}, nil)
	mockRepo.On("SeGlobalAgg", ctx, f).Return(dto.SeAgg{TotalSE: 20, ThisMonth: 20}, nil)
	mockRepo.On("SeStatusCount", ctx, f).Return(dto.SeStatusCount{TotalPerusahaan: 63, SudahMengisiKSE: 20, BelumMengisiKSE: 43}, nil)

	result, err := service.GetSummary(ctx, f)

	assert.NoError(t, err)
	assert.Equal(t, int64(63), result.Sektor[0].Total)
	assert.Equal(t, int64(20), result.SE.TotalSE)
	assert.Equal(t, int64(43), result.SEStatus.BelumMengisiKSE)

	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummary_CombinedFilter_SubSektorAndKategoriSE(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()

	subSektor := "sub-uuid-xyz"
	kategori := "Strategis"
	f := dto.DashboardFilter{SubSektorID: &subSektor, KategoriSE: &kategori}

	mockRepo.On("CountPerSektor", ctx, f).Return([]dto.SectorCount{}, nil)
	mockRepo.On("SeGlobalAgg", ctx, f).Return(dto.SeAgg{TotalSE: 15, Strategis: 15}, nil)
	mockRepo.On("SeStatusCount", ctx, f).Return(dto.SeStatusCount{TotalPerusahaan: 30, SudahMengisiKSE: 15, BelumMengisiKSE: 15}, nil)

	result, err := service.GetSummary(ctx, f)

	assert.NoError(t, err)
	assert.Equal(t, int64(15), result.SE.Strategis)
	assert.Equal(t, int64(0), result.SE.Tinggi)
	assert.Equal(t, int64(0), result.SE.Rendah)

	mockRepo.AssertExpectations(t)
}

/*
=====================================
 TEST INVARIANT DATA
=====================================
*/

func TestDashboardService_GetSummary_Invariant_KategoriSumEqualsTotal(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	mockRepo.On("CountPerSektor", ctx, f).Return([]dto.SectorCount{}, nil)
	mockRepo.On("SeGlobalAgg", ctx, f).Return(dto.SeAgg{
		TotalSE: 77, ThisMonth: 8, Strategis: 30, Tinggi: 28, Rendah: 19,
	}, nil)
	mockRepo.On("SeStatusCount", ctx, f).Return(dto.SeStatusCount{}, nil)

	result, err := service.GetSummary(ctx, f)

	assert.NoError(t, err)
	// Strategis + Tinggi + Rendah harus == TotalSE
	assert.Equal(t, result.SE.TotalSE, result.SE.Strategis+result.SE.Tinggi+result.SE.Rendah)

	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummary_Invariant_SEStatusSumEqualsTotal(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	mockRepo.On("CountPerSektor", ctx, f).Return([]dto.SectorCount{}, nil)
	mockRepo.On("SeGlobalAgg", ctx, f).Return(dto.SeAgg{TotalSE: 77}, nil)
	mockRepo.On("SeStatusCount", ctx, f).Return(dto.SeStatusCount{
		TotalPerusahaan: 251, SudahMengisiKSE: 77, BelumMengisiKSE: 174,
	}, nil)

	result, err := service.GetSummary(ctx, f)

	assert.NoError(t, err)
	// Sudah + Belum harus == Total
	assert.Equal(t, result.SEStatus.TotalPerusahaan, result.SEStatus.SudahMengisiKSE+result.SEStatus.BelumMengisiKSE)

	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummary_Invariant_ThisMonthNotExceedTotal(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	mockRepo.On("CountPerSektor", ctx, f).Return([]dto.SectorCount{
		{ID: "s1", Nama: "ILMATE", Total: 97, ThisMonth: 8},
		{ID: "s2", Nama: "IKFT", Total: 63, ThisMonth: 6},
	}, nil)
	mockRepo.On("SeGlobalAgg", ctx, f).Return(dto.SeAgg{TotalSE: 77, ThisMonth: 8}, nil)
	mockRepo.On("SeStatusCount", ctx, f).Return(dto.SeStatusCount{}, nil)

	result, err := service.GetSummary(ctx, f)

	assert.NoError(t, err)
	// ThisMonth di setiap sektor tidak boleh melebihi Total
	for _, s := range result.Sektor {
		assert.LessOrEqual(t, s.ThisMonth, s.Total, "sektor %s: ThisMonth melebihi Total", s.Nama)
	}
	// ThisMonth SE tidak boleh melebihi TotalSE
	assert.LessOrEqual(t, result.SE.ThisMonth, result.SE.TotalSE)

	mockRepo.AssertExpectations(t)
}

/*
=====================================
 TEST buildCacheKey - UNIQUENESS
=====================================
*/

func TestBuildCacheKey_UniquePerFilter(t *testing.T) {
	year := "2024"
	quarter := "2"
	from := "2024-01-01"
	to := "2024-01-31"
	sub := "sub-uuid"
	kat := "Strategis"

	filters := []dto.DashboardFilter{
		{},
		{Year: &year},
		{Year: &year, Quarter: &quarter},
		{From: &from, To: &to},
		{SubSektorID: &sub},
		{KategoriSE: &kat},
		{Year: &year, KategoriSE: &kat},
		{Year: &year, Quarter: &quarter, SubSektorID: &sub, KategoriSE: &kat},
	}

	keys := make(map[string]bool)
	for _, f := range filters {
		key := buildCacheKey(f)
		assert.False(t, keys[key], "cache key tidak unik untuk filter: %+v, key: %s", f, key)
		keys[key] = true
	}
}

func TestBuildCacheKey_SameFilterSameKey(t *testing.T) {
	year := "2024"
	f1 := dto.DashboardFilter{Year: &year}

	year2 := "2024" // pointer berbeda, nilai sama
	f2 := dto.DashboardFilter{Year: &year2}

	assert.Equal(t, buildCacheKey(f1), buildCacheKey(f2))
}
