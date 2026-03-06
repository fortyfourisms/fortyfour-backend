package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/pkg/cache"
)

type SubSektorServiceInterface interface {
	GetAll() ([]dto.SubSektorResponse, error)
	GetByID(id string) (*dto.SubSektorResponse, error)
	GetBySektorID(sektorID string) ([]dto.SubSektorResponse, error)
}

type SubSektorService struct {
	repo repository.SubSektorRepositoryInterface
	rc   cache.RedisInterface
}

func NewSubSektorService(repo repository.SubSektorRepositoryInterface, rc cache.RedisInterface) *SubSektorService {
	return &SubSektorService{repo: repo, rc: rc}
}

func (s *SubSektorService) GetAll() ([]dto.SubSektorResponse, error) {
	key := keyList("sub_sektor")
	var result []dto.SubSektorResponse
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

func (s *SubSektorService) GetByID(id string) (*dto.SubSektorResponse, error) {
	key := keyDetail("sub_sektor", id)
	var result dto.SubSektorResponse
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

func (s *SubSektorService) GetBySektorID(sektorID string) ([]dto.SubSektorResponse, error) {
	key := keyDetail("sub_sektor:sektor", sektorID)
	var result []dto.SubSektorResponse
	if cacheGet(s.rc, key, &result) {
		return result, nil
	}

	result, err := s.repo.GetBySektorID(sektorID)
	if err != nil {
		return nil, err
	}

	cacheSet(s.rc, key, result, TTLList)
	return result, nil
}
