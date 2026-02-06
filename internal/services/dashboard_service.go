package services

import (
	"context"

	"fortyfour-backend/internal/dto"
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
}

func NewDashboardService(repo DashboardRepositoryInterface) *DashboardService {
	return &DashboardService{repo: repo}
}

// GetSummary returns aggregated summary (sektor counts + se)
func (s *DashboardService) GetSummary(ctx context.Context, from, to *string) (*dto.DashboardSummary, error) {
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
	return &dto.DashboardSummary{
		Sektor: sectors,
		// Ikas:   ikasAgg, // TODO: re-enable ikas summary when ikas table is ready
		SE:     seAgg,
	}, nil
}