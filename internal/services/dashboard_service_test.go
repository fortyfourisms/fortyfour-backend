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

func (m *MockDashboardRepository) CsirtGlobalAgg(ctx context.Context, f dto.DashboardFilter) (dto.CsirtAgg, error) {
	args := m.Called(ctx, f)
	return args.Get(0).(dto.CsirtAgg), args.Error(1)
}

func (m *MockDashboardRepository) CsirtStatusCount(ctx context.Context, f dto.DashboardFilter) (dto.CsirtStatusCount, error) {
	args := m.Called(ctx, f)
	return args.Get(0).(dto.CsirtStatusCount), args.Error(1)
}

func createServiceWithMockRepo(mockRepo *MockDashboardRepository) *DashboardService {
	return NewDashboardService(mockRepo, nil)
}

func TestNewDashboardService(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := NewDashboardService(mockRepo, nil)
	assert.NotNil(t, service)
	assert.NotNil(t, service.repo)
}

// ── Sektor ──────────────────────────────────────────────────────────────────

func TestDashboardService_GetSummarySektor_Success(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	mockRepo.On("CountPerSektor", ctx, f).Return([]dto.SectorCount{
		{ID: "s1", Nama: "ILMATE", Total: 100, ThisMonth: 10},
	}, nil)

	result, err := service.GetSummarySektor(ctx, f)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Sektor, 1)
	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummarySektor_Error(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	mockRepo.On("CountPerSektor", ctx, f).Return(nil, errors.New("db error"))

	result, err := service.GetSummarySektor(ctx, f)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// ── IKAS ────────────────────────────────────────────────────────────────────

func TestDashboardService_GetSummaryIkas_Success(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	mockRepo.On("IkasGlobalAgg", ctx, f).Return(dto.IkasAgg{Total: 45}, nil)
	mockRepo.On("IkasStatusCount", ctx, f).Return(dto.IkasStatusCount{TotalPerusahaan: 100, SudahMengisiIKAS: 45, BelumMengisiIKAS: 55}, nil)

	result, err := service.GetSummaryIkas(ctx, f)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(45), result.Ikas.Total)
	assert.Equal(t, int64(55), result.IkasStatus.BelumMengisiIKAS)
	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummaryIkas_IkasAggError(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	mockRepo.On("IkasGlobalAgg", ctx, f).Return(dto.IkasAgg{}, errors.New("ikas error"))

	result, err := service.GetSummaryIkas(ctx, f)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// ── SE ──────────────────────────────────────────────────────────────────────

func TestDashboardService_GetSummarySE_Success(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	mockRepo.On("SeGlobalAgg", ctx, f).Return(dto.SeAgg{TotalSE: 75}, nil)
	mockRepo.On("SeStatusCount", ctx, f).Return(dto.SeStatusCount{TotalPerusahaan: 150, SudahMengisiSE: 75, BelumMengisiSE: 75}, nil)

	result, err := service.GetSummarySE(ctx, f)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(75), result.SE.TotalSE)
	assert.Equal(t, int64(75), result.SEStatus.BelumMengisiSE)
	mockRepo.AssertExpectations(t)
}

// ── CSIRT ────────────────────────────────────────────────────────────────────

func TestDashboardService_GetSummaryCSIRT_Success(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := createServiceWithMockRepo(mockRepo)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	mockRepo.On("CsirtGlobalAgg", ctx, f).Return(dto.CsirtAgg{TotalCSIRT: 30, ThisMonth: 5}, nil)
	mockRepo.On("CsirtStatusCount", ctx, f).Return(dto.CsirtStatusCount{TotalPerusahaan: 150, SudahMembentukCSIRT: 30, BelumMembentukCSIRT: 120}, nil)

	result, err := service.GetSummaryCSIRT(ctx, f)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(30), result.CSIRT.TotalCSIRT)
	assert.Equal(t, int64(120), result.CSIRTStatus.BelumMembentukCSIRT)
	mockRepo.AssertExpectations(t)
}

// ── Cache ────────────────────────────────────────────────────────────────────

func TestDashboardService_GetSummarySektor_CacheHit(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	rc := newDashboardTestRedis()
	service := NewDashboardService(mockRepo, rc)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	cached := dto.DashboardSektorResponse{
		Sektor: []dto.SectorCount{{ID: "cache-1", Nama: "Dari Cache", Total: 99}},
	}
	key := buildCacheKey("sektor", f)
	setDashboardCache(rc, key, cached)

	result, err := service.GetSummarySektor(ctx, f)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Dari Cache", result.Sektor[0].Nama)
	mockRepo.AssertNotCalled(t, "CountPerSektor")
}

func TestBuildCacheKey_UniquePerSection(t *testing.T) {
	f := dto.DashboardFilter{}
	k1 := buildCacheKey("sektor", f)
	k2 := buildCacheKey("ikas", f)
	k3 := buildCacheKey("kse", f)
	k4 := buildCacheKey("csirt", f)

	assert.NotEqual(t, k1, k2)
	assert.NotEqual(t, k2, k3)
	assert.NotEqual(t, k3, k4)
}

func TestBuildCacheKey_SameFilterSameKey(t *testing.T) {
	year := "2024"
	f1 := dto.DashboardFilter{Year: &year}
	year2 := "2024"
	f2 := dto.DashboardFilter{Year: &year2}

	assert.Equal(t, buildCacheKey("sektor", f1), buildCacheKey("sektor", f2))
}

// ── helpers ─────────────────────────────────────────────────────────────────

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
