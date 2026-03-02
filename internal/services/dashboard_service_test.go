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

func (m *MockDashboardRepository) CountPerSektor(ctx context.Context, from, to *string) ([]dto.SectorCount, error) {
	args := m.Called(ctx, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.SectorCount), args.Error(1)
}

func (m *MockDashboardRepository) SeGlobalAgg(ctx context.Context) (dto.SeAgg, error) {
	args := m.Called(ctx)
	return args.Get(0).(dto.SeAgg), args.Error(1)
}

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

func TestDashboardService_GetSummary_Success_NoDateFilter(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)

	ctx := context.Background()

	expectedSectors := []dto.SectorCount{
		{
			ID:        "sektor-1",
			Nama:      "ILMATE",
			Total:     100,
			ThisMonth: 10,
		},
		{
			ID:        "sektor-2",
			Nama:      "IKFT",
			Total:     50,
			ThisMonth: 5,
		},
	}

	expectedSE := dto.SeAgg{
		TotalSE: 75,
	}

	// Mock expectations - no date filter (nil, nil)
	mockRepo.On("CountPerSektor", ctx, (*string)(nil), (*string)(nil)).Return(expectedSectors, nil)
	mockRepo.On("SeGlobalAgg", ctx).Return(expectedSE, nil)

	result, err := service.GetSummary(ctx, nil, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Sektor, 2)
	assert.Equal(t, "ILMATE", result.Sektor[0].Nama)
	assert.Equal(t, int64(100), result.Sektor[0].Total)
	assert.Equal(t, int64(10), result.Sektor[0].ThisMonth)
	assert.Equal(t, int64(75), result.SE.TotalSE)

	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummary_Success_WithDateFilter(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)

	ctx := context.Background()
	from := "2024-01-01"
	to := "2024-01-31"

	expectedSectors := []dto.SectorCount{
		{
			ID:        "sektor-1",
			Nama:      "ILMATE",
			Total:     100,
			ThisMonth: 15,
		},
	}

	expectedSE := dto.SeAgg{
		TotalSE: 50,
	}

	mockRepo.On("CountPerSektor", ctx, &from, &to).Return(expectedSectors, nil)
	mockRepo.On("SeGlobalAgg", ctx).Return(expectedSE, nil)

	result, err := service.GetSummary(ctx, &from, &to)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Sektor, 1)
	assert.Equal(t, int64(15), result.Sektor[0].ThisMonth)
	assert.Equal(t, int64(50), result.SE.TotalSE)

	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummary_Success_EmptyResults(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)

	ctx := context.Background()

	// Empty sectors
	expectedSectors := []dto.SectorCount{}
	expectedSE := dto.SeAgg{TotalSE: 0}

	mockRepo.On("CountPerSektor", ctx, (*string)(nil), (*string)(nil)).Return(expectedSectors, nil)
	mockRepo.On("SeGlobalAgg", ctx).Return(expectedSE, nil)

	result, err := service.GetSummary(ctx, nil, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Sektor, 0)
	assert.Equal(t, int64(0), result.SE.TotalSE)

	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummary_Success_MultipleSektor(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)

	ctx := context.Background()

	expectedSectors := []dto.SectorCount{
		{ID: "sektor-1", Nama: "ILMATE", Total: 100, ThisMonth: 10},
		{ID: "sektor-2", Nama: "IKFT", Total: 80, ThisMonth: 8},
		{ID: "sektor-3", Nama: "ENERGI", Total: 60, ThisMonth: 6},
		{ID: "sektor-4", Nama: "TRANSPORTASI", Total: 40, ThisMonth: 4},
	}

	expectedSE := dto.SeAgg{TotalSE: 280}

	mockRepo.On("CountPerSektor", ctx, (*string)(nil), (*string)(nil)).Return(expectedSectors, nil)
	mockRepo.On("SeGlobalAgg", ctx).Return(expectedSE, nil)

	result, err := service.GetSummary(ctx, nil, nil)

	assert.NoError(t, err)
	assert.Len(t, result.Sektor, 4)
	assert.Equal(t, int64(280), result.SE.TotalSE)

	// Verify all sectors
	for i, sector := range result.Sektor {
		assert.Equal(t, expectedSectors[i].ID, sector.ID)
		assert.Equal(t, expectedSectors[i].Nama, sector.Nama)
		assert.Equal(t, expectedSectors[i].Total, sector.Total)
		assert.Equal(t, expectedSectors[i].ThisMonth, sector.ThisMonth)
	}

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

	mockRepo.On("CountPerSektor", ctx, (*string)(nil), (*string)(nil)).
		Return(nil, errors.New("database connection failed"))

	result, err := service.GetSummary(ctx, nil, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database connection failed", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummary_SeGlobalAggError(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)

	ctx := context.Background()

	expectedSectors := []dto.SectorCount{
		{ID: "sektor-1", Nama: "ILMATE", Total: 100, ThisMonth: 10},
	}

	mockRepo.On("CountPerSektor", ctx, (*string)(nil), (*string)(nil)).Return(expectedSectors, nil)
	mockRepo.On("SeGlobalAgg", ctx).Return(dto.SeAgg{}, errors.New("se table not found"))

	result, err := service.GetSummary(ctx, nil, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "se table not found", err.Error())

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
 TEST EDGE CASES
=====================================
*/

func TestDashboardService_GetSummary_NilContext(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)

	ctx := context.Background()

	expectedSectors := []dto.SectorCount{}
	expectedSE := dto.SeAgg{TotalSE: 0}

	mockRepo.On("CountPerSektor", ctx, (*string)(nil), (*string)(nil)).Return(expectedSectors, nil)
	mockRepo.On("SeGlobalAgg", ctx).Return(expectedSE, nil)

	result, err := service.GetSummary(ctx, nil, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummary_OnlyFromDate(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)

	ctx := context.Background()
	from := "2024-01-01"

	expectedSectors := []dto.SectorCount{
		{ID: "sektor-1", Nama: "ILMATE", Total: 50, ThisMonth: 5},
	}
	expectedSE := dto.SeAgg{TotalSE: 25}

	mockRepo.On("CountPerSektor", ctx, &from, (*string)(nil)).Return(expectedSectors, nil)
	mockRepo.On("SeGlobalAgg", ctx).Return(expectedSE, nil)

	result, err := service.GetSummary(ctx, &from, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummary_OnlyToDate(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)

	ctx := context.Background()
	to := "2024-01-31"

	expectedSectors := []dto.SectorCount{
		{ID: "sektor-1", Nama: "ILMATE", Total: 50, ThisMonth: 5},
	}
	expectedSE := dto.SeAgg{TotalSE: 25}

	mockRepo.On("CountPerSektor", ctx, (*string)(nil), &to).Return(expectedSectors, nil)
	mockRepo.On("SeGlobalAgg", ctx).Return(expectedSE, nil)

	result, err := service.GetSummary(ctx, nil, &to)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockRepo.AssertExpectations(t)
}

/*
=====================================
 TEST DATA INTEGRITY
=====================================
*/

func TestDashboardService_GetSummary_DataIntegrity(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)

	ctx := context.Background()

	expectedSectors := []dto.SectorCount{
		{
			ID:        "sektor-123",
			Nama:      "ILMATE",
			Total:     150,
			ThisMonth: 25,
		},
	}

	expectedSE := dto.SeAgg{
		TotalSE: 150,
	}

	mockRepo.On("CountPerSektor", ctx, (*string)(nil), (*string)(nil)).Return(expectedSectors, nil)
	mockRepo.On("SeGlobalAgg", ctx).Return(expectedSE, nil)

	result, err := service.GetSummary(ctx, nil, nil)

	assert.NoError(t, err)

	// Verify data integrity - ThisMonth should not exceed Total
	for _, sector := range result.Sektor {
		assert.LessOrEqual(t, sector.ThisMonth, sector.Total,
			"ThisMonth count should not exceed Total count")
	}

	// Verify SE count matches
	assert.Equal(t, expectedSE.TotalSE, result.SE.TotalSE)

	mockRepo.AssertExpectations(t)
}

/*
=====================================
 TEST WITH CONTEXT CANCELLATION
=====================================
*/

func TestDashboardService_GetSummary_ContextCancellation(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mockRepo.On("CountPerSektor", mock.Anything, (*string)(nil), (*string)(nil)).
		Return(nil, context.Canceled).Maybe()

	result, err := service.GetSummary(ctx, nil, nil)

	// Should return error (either from context or from repo)
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

	// Pre-populate cache dengan key yang sama yang digunakan service (from=nil → "nil", to=nil → "nil")
	cached := dto.DashboardSummary{
		Sektor: []dto.SectorCount{{ID: "cache-1", Nama: "Dari Cache", Total: 99}},
		SE:     dto.SeAgg{TotalSE: 42},
	}
	setDashboardCache(rc, "dashboard:summary:nil:nil", cached)

	result, err := service.GetSummary(ctx, nil, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Sektor, 1)
	assert.Equal(t, "Dari Cache", result.Sektor[0].Nama)
	assert.Equal(t, int64(42), result.SE.TotalSE)

	// Repo tidak boleh dipanggil sama sekali
	mockRepo.AssertNotCalled(t, "CountPerSektor")
	mockRepo.AssertNotCalled(t, "SeGlobalAgg")
}

func TestDashboardService_GetSummary_CacheMiss_SetsCache(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	rc := newDashboardTestRedis()
	service := NewDashboardService(mockRepo, rc)

	ctx := context.Background()

	mockRepo.On("CountPerSektor", ctx, (*string)(nil), (*string)(nil)).
		Return([]dto.SectorCount{{ID: "s1", Nama: "ILMATE", Total: 10}}, nil)
	mockRepo.On("SeGlobalAgg", ctx).Return(dto.SeAgg{TotalSE: 10}, nil)

	_, err := service.GetSummary(ctx, nil, nil)

	assert.NoError(t, err)

	exists, _ := rc.Exists("dashboard:summary:nil:nil")
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

	mockRepo.On("CountPerSektor", ctx, &from, &to).
		Return([]dto.SectorCount{{ID: "s1", Total: 5}}, nil)
	mockRepo.On("SeGlobalAgg", ctx).Return(dto.SeAgg{TotalSE: 5}, nil)

	_, err := service.GetSummary(ctx, &from, &to)

	assert.NoError(t, err)

	// Key harus spesifik berdasarkan tanggal — berbeda dari key tanpa filter
	keyWithFilter := "dashboard:summary:2024-01-01:2024-01-31"
	keyNoFilter := "dashboard:summary:nil:nil"

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