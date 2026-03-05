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
	SeGlobalAgg(ctx context.Context, f dto.DashboardFilter) (dto.SeAgg, error)
	SeStatusCount(ctx context.Context, f dto.DashboardFilter) (dto.SeStatusCount, error)
	// TODO: re-enable ikas status when ikas table is ready
	// IkasStatusCount(ctx context.Context, f dto.DashboardFilter) (dto.IkasStatusCount, error)
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
func buildCacheKey(f dto.DashboardFilter) string {
	str := func(p *string) string {
		if p == nil {
			return "nil"
		}
		return *p
	}
	return fmt.Sprintf("dashboard:summary:%s:%s:%s:%s:%s:%s",
		str(f.From),
		str(f.To),
		str(f.Year),
		str(f.Quarter),
		str(f.SubSektorID),
		str(f.KategoriSE),
	)
}

// GetSummary returns aggregated summary (sektor counts + se agg + se status)
func (s *DashboardService) GetSummary(ctx context.Context, f dto.DashboardFilter) (*dto.DashboardSummary, error) {
	key := buildCacheKey(f)

	var result dto.DashboardSummary
	if cacheGet(s.rc, key, &result) {
		return &result, nil
	}

	sectors, err := s.repo.CountPerSektor(ctx, f)
	if err != nil {
		return nil, err
	}

	// TODO: re-enable ikas summary when ikas table is ready
	// ikasAgg, err := s.repo.IkasGlobalAgg(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	seAgg, err := s.repo.SeGlobalAgg(ctx, f)
	if err != nil {
		return nil, err
	}

	seStatus, err := s.repo.SeStatusCount(ctx, f)
	if err != nil {
		return nil, err
	}

	// TODO: re-enable ikas status when ikas table is ready
	// ikasStatus, err := s.repo.IkasStatusCount(ctx, f)
	// if err != nil {
	// 	return nil, err
	// }

	summary := &dto.DashboardSummary{
		Sektor: sectors,
		// Ikas:       ikasAgg,   // TODO: re-enable ikas summary when ikas table is ready
		SE:       seAgg,
		SEStatus: seStatus,
		// IkasStatus: ikasStatus, // TODO: re-enable ikas status when ikas table is ready
	}

	cacheSet(s.rc, key, summary, TTLList)
	return summary, nil
}