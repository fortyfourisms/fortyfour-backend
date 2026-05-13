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
	Create(req dto.SubSektorRequest) (*dto.SubSektorResponse, error)
	Update(id string, req dto.SubSektorRequest) (*dto.SubSektorResponse, error)
	Delete(id string) error
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

func (s *SubSektorService) Create(req dto.SubSektorRequest) (*dto.SubSektorResponse, error) {
	data, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}
	cacheDelete(s.rc, keyList("sub_sektor"))
	cacheDeletePattern(s.rc, "sub_sektor:sektor:*")
	return data, nil
}

func (s *SubSektorService) Update(id string, req dto.SubSektorRequest) (*dto.SubSektorResponse, error) {
	data, err := s.repo.Update(id, req)
	if err != nil {
		return nil, err
	}
	cacheDelete(s.rc, keyList("sub_sektor"))
	cacheDelete(s.rc, keyDetail("sub_sektor", id))
	cacheDeletePattern(s.rc, "sub_sektor:sektor:*")
	return data, nil
}

func (s *SubSektorService) Delete(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	cacheDelete(s.rc, keyList("sub_sektor"))
	cacheDelete(s.rc, keyDetail("sub_sektor", id))
	cacheDeletePattern(s.rc, "sub_sektor:sektor:*")
	return nil
}
