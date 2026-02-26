package services

import (
	"errors"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/pkg/cache"
	"github.com/google/uuid"
	"strings"
)

type PICService struct {
	repo repository.PICRepositoryInterface
	rc   cache.RedisInterface
}

func NewPICService(repo repository.PICRepositoryInterface, rc cache.RedisInterface) *PICService {
	return &PICService{repo: repo, rc: rc}
}

func (s *PICService) Create(req dto.CreatePICRequest) (*dto.PICResponse, error) {
	if req.Nama == nil || strings.TrimSpace(*req.Nama) == "" {
		return nil, errors.New("nama wajib diisi")
	}

	if req.IDPerusahaan == nil || strings.TrimSpace(*req.IDPerusahaan) == "" {
		return nil, errors.New("id_perusahaan wajib diisi")
	}

	id := uuid.New().String()

	if err := s.repo.Create(req, id); err != nil {
		return nil, err
	}

	result, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	cacheSet(s.rc, keyDetail("pic", id), result, TTLDetail)
	cacheDelete(s.rc, keyList("pic"))

	return result, nil
}

func (s *PICService) GetAll() ([]dto.PICResponse, error) {
	key := keyList("pic")
	var result []dto.PICResponse
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

func (s *PICService) GetByID(id string) (*dto.PICResponse, error) {
	key := keyDetail("pic", id)
	var result dto.PICResponse
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

func (s *PICService) Update(id string, req dto.UpdatePICRequest) (*dto.PICResponse, error) {
	if err := s.repo.Update(id, req); err != nil {
		return nil, err
	}

	updated, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	cacheDelete(s.rc, keyDetail("pic", id))
	cacheDelete(s.rc, keyList("pic"))
	return updated, nil
}

func (s *PICService) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	cacheDelete(s.rc, keyDetail("pic", id))
	cacheDelete(s.rc, keyList("pic"))

	return nil
}
