package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"ikas/internal/models"
	"ikas/internal/repository"
	"ikas/pkg/cache"
)

type DeteksiService struct {
	repo     repository.DeteksiRepositoryInterface
	ikasRepo repository.IkasRepositoryInterface
	cache    cache.RedisInterface
}

func NewDeteksiService(
	repo repository.DeteksiRepositoryInterface,
	ikasRepo repository.IkasRepositoryInterface,
	cache cache.RedisInterface,
) *DeteksiService {
	return &DeteksiService{
		repo:     repo,
		ikasRepo: ikasRepo,
		cache:    cache,
	}
}

func (s *DeteksiService) GetAll(userRole string) ([]models.Deteksi, error) {
	if userRole != "admin" && userRole != "staff" {
		return nil, errors.New("anda tidak memiliki akses untuk melihat semua data")
	}
	return s.repo.GetAll()
}

func (s *DeteksiService) GetByIkasID(ikasID string, userRole string, userPerusahaanID string) ([]models.Deteksi, error) {
	if userRole != "admin" && userRole != "staff" {
		owned, err := s.ikasRepo.CheckOwnership(ikasID, userPerusahaanID)
		if err != nil {
			return nil, err
		}
		if !owned {
			return nil, errors.New("anda tidak memiliki akses ke data asesmen ini")
		}
	}

	cacheKey := fmt.Sprintf("%s%s", cache.CacheKeyPrefixDeteksi, ikasID)
	if s.cache != nil {
		cachedData, err := s.cache.Get(cacheKey)
		if err == nil && cachedData != "" {
			var data []models.Deteksi
			if err := json.Unmarshal([]byte(cachedData), &data); err == nil {
				return data, nil
			}
		}
	}

	data, err := s.repo.GetByIkasID(ikasID)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		go func(key string, dataToCache []models.Deteksi) {
			jsonData, err := json.Marshal(dataToCache)
			if err == nil {
				_ = s.cache.Set(key, string(jsonData), cache.DefaultCacheExpiration)
			}
		}(cacheKey, data)
	}

	return data, nil
}

func (s *DeteksiService) GetByID(id string, userRole string, userPerusahaanID string) (*models.Deteksi, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if userRole != "admin" && userRole != "staff" {
		owned, err := s.ikasRepo.CheckOwnership(data.IkasID, userPerusahaanID)
		if err != nil {
			return nil, err
		}
		if !owned {
			return nil, errors.New("anda tidak memiliki akses ke data ini")
		}
	}

	return data, nil
}

func (s *DeteksiService) GetByPerusahaanID(perusahaanID string, userRole string, userPerusahaanID string) ([]models.Deteksi, error) {
	if userRole != "admin" && userRole != "staff" {
		if perusahaanID != userPerusahaanID {
			return nil, errors.New("anda tidak memiliki akses ke data perusahaan ini")
		}
	}
	return s.repo.GetByPerusahaanID(perusahaanID)
}
