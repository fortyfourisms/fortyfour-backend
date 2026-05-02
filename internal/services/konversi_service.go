package services

import (
	"fmt"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/pkg/cache"
	"math"
)

type KonversiServiceInterface interface {
	GetKonversi(perusahaanID string) ([]dto.KonversiResponse, error)
}

type KonversiService struct {
	repo repository.KonversiRepositoryInterface
	rc   cache.RedisInterface
}

func NewKonversiService(repo repository.KonversiRepositoryInterface, rc cache.RedisInterface) *KonversiService {
	return &KonversiService{
		repo: repo,
		rc:   rc,
	}
}

func (s *KonversiService) GetKonversi(perusahaanID string) ([]dto.KonversiResponse, error) {
	cacheKey := "konversi:all"
	if perusahaanID != "" {
		cacheKey = fmt.Sprintf("konversi:%s", perusahaanID)
	}

	var results []dto.KonversiResponse
	if cacheGet(s.rc, cacheKey, &results) {
		return results, nil
	}

	data, err := s.repo.GetAllKonversi(perusahaanID)
	if err != nil {
		return nil, err
	}

	for i := range data {
		data[i].TotalPoin = data[i].PoinIkas + data[i].PoinKse + data[i].PoinSurvey + data[i].PoinCsirt
		// Hitung persentase: (Total / 4) * 100
		percentage := (float64(data[i].TotalPoin) / 4.0) * 100
		data[i].Persentase = math.Round(percentage*100) / 100 // Round to 2 decimal places
	}

	cacheSet(s.rc, cacheKey, data, TTLList)
	return data, nil
}
