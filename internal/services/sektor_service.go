package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/pkg/cache"
)

type SektorServiceInterface interface {
	GetAll() ([]dto.SektorResponse, error)
	GetByID(id string) (*dto.SektorResponse, error)
}

type SektorService struct {
	repo repository.SektorRepositoryInterface
	rc   cache.RedisInterface
}

func NewSektorService(repo repository.SektorRepositoryInterface, rc cache.RedisInterface) *SektorService {
	return &SektorService{repo: repo, rc: rc}
}

func (s *SektorService) GetAll() ([]dto.SektorResponse, error) {
	key := keyList("sektor")
	var result []dto.SektorResponse
	if cacheGet(s.rc, key, &result) {
		return result, nil
	}

	result, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	cacheSet(s.rc, key, result, TTLList)
	return result, nil
}

func (s *SektorService) GetByID(id string) (*dto.SektorResponse, error) {
	key := keyDetail("sektor", id)
	var result dto.SektorResponse
	if cacheGet(s.rc, key, &result) {
		return &result, nil
	}

	data, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	cacheSet(s.rc, key, data, TTLDetail)
	return data, nil
}
