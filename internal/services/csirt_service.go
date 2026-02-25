package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/pkg/cache"
	"github.com/google/uuid"
)

type CsirtServiceInterface interface {
	GetAll() ([]dto.CsirtResponse, error)
	GetByID(id string) (*dto.CsirtResponse, error)
	Create(req dto.CreateCsirtRequest) (*models.Csirt, error)
	Update(id string, req dto.UpdateCsirtRequest) (*models.Csirt, error)
	Delete(id string) error
}

type CsirtService struct {
	repo repository.CsirtRepositoryInterface
	rc   cache.RedisInterface
}

func NewCsirtService(repo repository.CsirtRepositoryInterface, rc cache.RedisInterface) *CsirtService {
	return &CsirtService{repo: repo, rc: rc}
}

func (s *CsirtService) Create(req dto.CreateCsirtRequest) (*models.Csirt, error) {
	id := uuid.New().String()
	if err := s.repo.Create(req, id); err != nil {
		return nil, err
	}

	result, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	cacheSet(s.rc, keyDetail("csirt", id), result, TTLDetail)
	cacheDelete(s.rc, keyList("csirt"))

	return result, nil
}

func (s *CsirtService) GetAll() ([]dto.CsirtResponse, error) {
	key := keyList("csirt")
	var result []dto.CsirtResponse
	if cacheGet(s.rc, key, &result) {
		return result, nil
	}

	result, err := s.repo.GetAllWithPerusahaan()
	if err != nil {
		return nil, err
	}

	cacheSet(s.rc, key, result, TTLList)
	return result, nil
}

func (s *CsirtService) GetByID(id string) (*dto.CsirtResponse, error) {
	key := keyDetail("csirt", id)
	var result dto.CsirtResponse
	if cacheGet(s.rc, key, &result) {
		return &result, nil
	}

	data, err := s.repo.GetByIDWithPerusahaan(id)
	if err != nil {
		return nil, err
	}

	cacheSet(s.rc, key, data, TTLDetail)
	return data, nil
}

func (s *CsirtService) Update(id string, req dto.UpdateCsirtRequest) (*models.Csirt, error) {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.NamaCsirt != nil {
		c.NamaCsirt = *req.NamaCsirt
	}
	if req.WebCsirt != nil {
		c.WebCsirt = *req.WebCsirt
	}
	if req.TeleponCsirt != nil {
		c.TeleponCsirt = req.TeleponCsirt
	}
	if req.PhotoCsirt != nil {
		c.PhotoCsirt = *req.PhotoCsirt
	}
	if req.FileRFC2350 != nil {
		c.FileRFC2350 = *req.FileRFC2350
	}
	if req.FilePublicKeyPGP != nil {
		c.FilePublicKeyPGP = *req.FilePublicKeyPGP
	}

	if err := s.repo.Update(id, *c); err != nil {
		return nil, err
	}

	cacheDelete(s.rc, keyDetail("csirt", id))
	cacheDelete(s.rc, keyList("csirt"))

	return c, nil
}

func (s *CsirtService) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	cacheDelete(s.rc, keyDetail("csirt", id))
	cacheDelete(s.rc, keyList("csirt"))

	return nil
}
