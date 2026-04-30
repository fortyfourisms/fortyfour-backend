package services

import (
	"context"
	"fmt"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/pkg/cache"
)

/*
=====================================
 DASHBOARD REPOSITORY INTERFACE
=====================================
*/

type DashboardRepositoryInterface interface {
	CountPerSektor(ctx context.Context, f dto.DashboardFilter) ([]dto.SectorCount, error)
	IkasGlobalAgg(ctx context.Context, f dto.DashboardFilter) (dto.IkasAgg, error)
	SeGlobalAgg(ctx context.Context, f dto.DashboardFilter) (dto.SeAgg, error)
	SeStatusCount(ctx context.Context, f dto.DashboardFilter) (dto.SeStatusCount, error)
	IkasStatusCount(ctx context.Context, f dto.DashboardFilter) (dto.IkasStatusCount, error)
	CsirtGlobalAgg(ctx context.Context, f dto.DashboardFilter) (dto.CsirtAgg, error)
	CsirtStatusCount(ctx context.Context, f dto.DashboardFilter) (dto.CsirtStatusCount, error)
}

/*
=====================================
 DASHBOARD SERVICE
=====================================
*/

type DashboardService struct {
	repo DashboardRepositoryInterface
	rc   cache.RedisInterface
}

func NewDashboardService(repo DashboardRepositoryInterface, rc cache.RedisInterface) *DashboardService {
	return &DashboardService{repo: repo, rc: rc}
}

// buildCacheKey membuat cache key unik berdasarkan semua parameter filter
func buildCacheKey(prefix string, f dto.DashboardFilter) string {
	str := func(p *string) string {
		if p == nil {
			return "nil"
		}
		return *p
	}
	return fmt.Sprintf("dashboard:%s:%s:%s:%s:%s:%s:%s",
		prefix,
		str(f.From),
		str(f.To),
		str(f.Year),
		str(f.Quarter),
		str(f.SubSektorID),
		str(f.KategoriSE),
	)
}

// GetSummarySektor returns sektor counts
func (s *DashboardService) GetSummarySektor(ctx context.Context, f dto.DashboardFilter) (*dto.DashboardSektorResponse, error) {
	key := buildCacheKey("sektor", f)

	var result dto.DashboardSektorResponse
	if cacheGet(s.rc, key, &result) {
		return &result, nil
	}

	sectors, err := s.repo.CountPerSektor(ctx, f)
	if err != nil {
		return nil, err
	}

	summary := &dto.DashboardSektorResponse{
		Sektor: sectors,
	}

	cacheSet(s.rc, key, summary, TTLList)
	return summary, nil
}

// GetSummaryIkas returns IKAS aggregate + status
func (s *DashboardService) GetSummaryIkas(ctx context.Context, f dto.DashboardFilter) (*dto.DashboardIkasResponse, error) {
	key := buildCacheKey("ikas", f)

	var result dto.DashboardIkasResponse
	if cacheGet(s.rc, key, &result) {
		return &result, nil
	}

	ikasAgg, err := s.repo.IkasGlobalAgg(ctx, f)
	if err != nil {
		return nil, err
	}

	ikasStatus, err := s.repo.IkasStatusCount(ctx, f)
	if err != nil {
		return nil, err
	}

	summary := &dto.DashboardIkasResponse{
		Ikas:       ikasAgg,
		IkasStatus: ikasStatus,
	}

	cacheSet(s.rc, key, summary, TTLList)
	return summary, nil
}

// GetSummarySE returns SE aggregate + status
func (s *DashboardService) GetSummarySE(ctx context.Context, f dto.DashboardFilter) (*dto.DashboardSEResponse, error) {
	key := buildCacheKey("se", f)

	var result dto.DashboardSEResponse
	if cacheGet(s.rc, key, &result) {
		return &result, nil
	}

	seAgg, err := s.repo.SeGlobalAgg(ctx, f)
	if err != nil {
		return nil, err
	}

	seStatus, err := s.repo.SeStatusCount(ctx, f)
	if err != nil {
		return nil, err
	}

	summary := &dto.DashboardSEResponse{
		SE:       seAgg,
		SEStatus: seStatus,
	}

	cacheSet(s.rc, key, summary, TTLList)
	return summary, nil
}

// GetSummaryCSIRT returns CSIRT aggregate + status
func (s *DashboardService) GetSummaryCSIRT(ctx context.Context, f dto.DashboardFilter) (*dto.DashboardCSIRTResponse, error) {
	key := buildCacheKey("csirt", f)

	var result dto.DashboardCSIRTResponse
	if cacheGet(s.rc, key, &result) {
		return &result, nil
	}

	csirtAgg, err := s.repo.CsirtGlobalAgg(ctx, f)
	if err != nil {
		return nil, err
	}

	csirtStatus, err := s.repo.CsirtStatusCount(ctx, f)
	if err != nil {
		return nil, err
	}

	summary := &dto.DashboardCSIRTResponse{
		CSIRT:       csirtAgg,
		CSIRTStatus: csirtStatus,
	}

	cacheSet(s.rc, key, summary, TTLList)
	return summary, nil
}
