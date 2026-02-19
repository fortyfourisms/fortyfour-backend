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
	CountPerSektor(ctx context.Context, from, to *string) ([]dto.SectorCount, error)
	SeGlobalAgg(ctx context.Context) (dto.SeAgg, error)
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

// GetSummary returns aggregated summary (sektor counts + se)
func (s *DashboardService) GetSummary(ctx context.Context, from, to *string) (*dto.DashboardSummary, error) {
	// Key dinamis berdasarkan filter tanggal
	fromStr := "nil"
	toStr := "nil"
	if from != nil {
		fromStr = *from
	}
	if to != nil {
		toStr = *to
	}
	key := fmt.Sprintf("dashboard:summary:%s:%s", fromStr, toStr)

	var result dto.DashboardSummary
	if cacheGet(s.rc, key, &result) {
		return &result, nil
	}

	sectors, err := s.repo.CountPerSektor(ctx, from, to)
	if err != nil {
		return nil, err
	}
	// TODO: re-enable ikas summary when ikas table is ready
	// ikasAgg, err := s.repo.IkasGlobalAgg(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	seAgg, err := s.repo.SeGlobalAgg(ctx)
	if err != nil {
		return nil, err
	}

	summary := &dto.DashboardSummary{
		Sektor: sectors,
		// Ikas:   ikasAgg, // TODO: re-enable ikas summary when ikas table is ready
		SE:     seAgg,
	}

	cacheSet(s.rc, key, summary, TTLList)
	return summary, nil
}
