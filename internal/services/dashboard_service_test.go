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

func (m *MockDashboardRepository) IkasGlobalAgg(ctx context.Context, f dto.DashboardFilter) (dto.IkasAgg, error) {
	args := m.Called(ctx, f)
	return args.Get(0).(dto.IkasAgg), args.Error(1)
}

func (m *MockDashboardRepository) SeGlobalAgg(ctx context.Context, f dto.DashboardFilter) (dto.SeAgg, error) {
	args := m.Called(ctx, f)
	return args.Get(0).(dto.SeAgg), args.Error(1)
}

func (m *MockDashboardRepository) SeStatusCount(ctx context.Context, f dto.DashboardFilter) (dto.SeStatusCount, error) {
	args := m.Called(ctx, f)
	return args.Get(0).(dto.SeStatusCount), args.Error(1)
}

func (m *MockDashboardRepository) IkasStatusCount(ctx context.Context, f dto.DashboardFilter) (dto.IkasStatusCount, error) {
	args := m.Called(ctx, f)
	return args.Get(0).(dto.IkasStatusCount), args.Error(1)
}

func createServiceWithMockRepo(mockRepo *MockDashboardRepository) *DashboardService {
	return NewDashboardService(mockRepo, nil)
}

func TestDashboardService_GetSummary_Success_NoFilter(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	expectedSectors := []dto.SectorCount{
		{ID: "sektor-1", Nama: "ILMATE", Total: 100, ThisMonth: 10},
		{ID: "sektor-2", Nama: "IKFT", Total: 50, ThisMonth: 5},
	}
	expectedIKAS := dto.IkasAgg{Total: 45, AvgNilaiKematangan: 2.7, AvgTargetNilai: 4}
	expectedSE := dto.SeAgg{TotalSE: 75, ThisMonth: 8, Strategis: 30, Tinggi: 25, Rendah: 20}
	expectedSEStatus := dto.SeStatusCount{TotalPerusahaan: 150, SudahMengisiKSE: 75, BelumMengisiKSE: 75}
	expectedIkasStatus := dto.IkasStatusCount{TotalPerusahaan: 150, SudahMengisiIKAS: 45, BelumMengisiIKAS: 105}

	mockRepo.On("CountPerSektor", ctx, f).Return(expectedSectors, nil)
	mockRepo.On("IkasGlobalAgg", ctx, f).Return(expectedIKAS, nil)
	mockRepo.On("SeGlobalAgg", ctx, f).Return(expectedSE, nil)
	mockRepo.On("SeStatusCount", ctx, f).Return(expectedSEStatus, nil)
	mockRepo.On("IkasStatusCount", ctx, f).Return(expectedIkasStatus, nil)

	result, err := service.GetSummary(ctx, f)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Sektor, 2)
	assert.Equal(t, int64(45), result.Ikas.Total)
	assert.Equal(t, 2.7, result.Ikas.AvgNilaiKematangan)
	assert.Equal(t, int64(75), result.SE.TotalSE)
	assert.Equal(t, int64(75), result.SEStatus.SudahMengisiKSE)
	assert.Equal(t, int64(45), result.IkasStatus.SudahMengisiIKAS)
	assert.Equal(t, int64(105), result.IkasStatus.BelumMengisiIKAS)
	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummary_Success_WithFilter(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()
	year := "2024"
	quarter := "2"
	f := dto.DashboardFilter{Year: &year, Quarter: &quarter}

	mockRepo.On("CountPerSektor", ctx, f).Return([]dto.SectorCount{{ID: "s1", Total: 20, ThisMonth: 5}}, nil)
	mockRepo.On("IkasGlobalAgg", ctx, f).Return(dto.IkasAgg{Total: 8, AvgNilaiKematangan: 2.5, AvgTargetNilai: 4}, nil)
	mockRepo.On("SeGlobalAgg", ctx, f).Return(dto.SeAgg{TotalSE: 11, ThisMonth: 3}, nil)
	mockRepo.On("SeStatusCount", ctx, f).Return(dto.SeStatusCount{TotalPerusahaan: 20, SudahMengisiKSE: 11, BelumMengisiKSE: 9}, nil)
	mockRepo.On("IkasStatusCount", ctx, f).Return(dto.IkasStatusCount{TotalPerusahaan: 20, SudahMengisiIKAS: 8, BelumMengisiIKAS: 12}, nil)

	result, err := service.GetSummary(ctx, f)

	assert.NoError(t, err)
	assert.Equal(t, int64(8), result.Ikas.Total)
	assert.Equal(t, int64(12), result.IkasStatus.BelumMengisiIKAS)
	mockRepo.AssertExpectations(t)
}

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

func TestDashboardService_GetSummary_IkasGlobalAggError(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	mockRepo.On("CountPerSektor", ctx, f).Return([]dto.SectorCount{}, nil)
	mockRepo.On("IkasGlobalAgg", ctx, f).Return(dto.IkasAgg{}, errors.New("ikas query failed"))

	result, err := service.GetSummary(ctx, f)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "ikas query failed", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummary_IkasStatusCountError(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	mockRepo.On("CountPerSektor", ctx, f).Return([]dto.SectorCount{}, nil)
	mockRepo.On("IkasGlobalAgg", ctx, f).Return(dto.IkasAgg{}, nil)
	mockRepo.On("SeGlobalAgg", ctx, f).Return(dto.SeAgg{}, nil)
	mockRepo.On("SeStatusCount", ctx, f).Return(dto.SeStatusCount{}, nil)
	mockRepo.On("IkasStatusCount", ctx, f).Return(dto.IkasStatusCount{}, errors.New("ikas status query failed"))

	result, err := service.GetSummary(ctx, f)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "ikas status query failed", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestNewDashboardService(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := NewDashboardService(mockRepo, nil)
	assert.NotNil(t, service)
	assert.NotNil(t, service.repo)
}

func TestDashboardService_GetSummary_ContextCancellation(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	f := dto.DashboardFilter{}

	mockRepo.On("CountPerSektor", mock.Anything, f).Return(nil, context.Canceled).Maybe()

	result, err := service.GetSummary(ctx, f)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestDashboardService_GetSummary_CacheHit_SkipRepo(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	rc := newDashboardTestRedis()
	service := NewDashboardService(mockRepo, rc)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	cached := dto.DashboardSummary{
		Sektor:     []dto.SectorCount{{ID: "cache-1", Nama: "Dari Cache", Total: 99}},
		Ikas:       dto.IkasAgg{Total: 11},
		SE:         dto.SeAgg{TotalSE: 42},
		SEStatus:   dto.SeStatusCount{TotalPerusahaan: 99},
		IkasStatus: dto.IkasStatusCount{TotalPerusahaan: 99, SudahMengisiIKAS: 11, BelumMengisiIKAS: 88},
	}
	key := buildCacheKey(f)
	setDashboardCache(rc, key, cached)

	result, err := service.GetSummary(ctx, f)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Dari Cache", result.Sektor[0].Nama)
	assert.Equal(t, int64(11), result.Ikas.Total)
	assert.Equal(t, int64(42), result.SE.TotalSE)
	mockRepo.AssertNotCalled(t, "CountPerSektor")
	mockRepo.AssertNotCalled(t, "IkasGlobalAgg")
	mockRepo.AssertNotCalled(t, "SeGlobalAgg")
	mockRepo.AssertNotCalled(t, "SeStatusCount")
	mockRepo.AssertNotCalled(t, "IkasStatusCount")
}

func TestDashboardService_GetSummary_CacheMiss_SetsCache(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	rc := newDashboardTestRedis()
	service := NewDashboardService(mockRepo, rc)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	mockRepo.On("CountPerSektor", ctx, f).Return([]dto.SectorCount{{ID: "s1", Nama: "ILMATE", Total: 10}}, nil)
	mockRepo.On("IkasGlobalAgg", ctx, f).Return(dto.IkasAgg{Total: 4}, nil)
	mockRepo.On("SeGlobalAgg", ctx, f).Return(dto.SeAgg{TotalSE: 10}, nil)
	mockRepo.On("SeStatusCount", ctx, f).Return(dto.SeStatusCount{}, nil)
	mockRepo.On("IkasStatusCount", ctx, f).Return(dto.IkasStatusCount{}, nil)

	_, err := service.GetSummary(ctx, f)
	assert.NoError(t, err)

	exists, _ := rc.Exists(buildCacheKey(f))
	assert.True(t, exists, "hasil harus di-cache setelah GetSummary")
	mockRepo.AssertExpectations(t)
}

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

	year2 := "2024"
	f2 := dto.DashboardFilter{Year: &year2}

	assert.Equal(t, buildCacheKey(f1), buildCacheKey(f2))
}

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
