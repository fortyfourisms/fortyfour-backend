package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/pkg/cache"

	"github.com/google/uuid"
)

type SdmCsirtServiceInterface interface {
	Create(req dto.CreateSdmCsirtRequest) (string, error)
	GetAll() ([]dto.SdmCsirtResponse, error)
	GetByID(id string) (*dto.SdmCsirtResponse, error)
	Update(id string, req dto.UpdateSdmCsirtRequest) error
	Delete(id string) error
}

type SdmCsirtService struct {
	repo repository.SdmCsirtRepositoryInterface
	rc   cache.RedisInterface
}

func NewSdmCsirtService(repo repository.SdmCsirtRepositoryInterface, rc cache.RedisInterface) *SdmCsirtService {
	return &SdmCsirtService{repo: repo, rc: rc}
}

func (s *SdmCsirtService) Create(req dto.CreateSdmCsirtRequest) (string, error) {
	id := uuid.New().String()

	if err := s.repo.Create(req, id); err != nil {
		return "", err
	}

	result, err := s.repo.GetByID(id)
	if err != nil {
		return "", err
	}

	cacheSet(s.rc, keyDetail("sdm", id), result, TTLDetail)
	cacheDelete(s.rc, keyList("sdm"))

	return id, nil
}

func (s *SdmCsirtService) GetAll() ([]dto.SdmCsirtResponse, error) {
	key := keyList("sdm")
	var result []dto.SdmCsirtResponse

	if cacheGet(s.rc, key, &result) {
		return result, nil
	}

	data, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	cacheSet(s.rc, key, data, TTLList)
	return data, nil
}

func (s *SdmCsirtService) GetByID(id string) (*dto.SdmCsirtResponse, error) {
	key := keyDetail("sdm", id)
	var result dto.SdmCsirtResponse

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

func (s *SdmCsirtService) Update(id string, req dto.UpdateSdmCsirtRequest) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if req.NamaPersonel != nil {
		existing.NamaPersonel = *req.NamaPersonel
	}
	if req.JabatanCsirt != nil {
		existing.JabatanCsirt = *req.JabatanCsirt
	}
	if req.JabatanPerusahaan != nil {
		existing.JabatanPerusahaan = *req.JabatanPerusahaan
	}
	if req.Skill != nil {
		existing.Skill = *req.Skill
	}
	if req.Sertifikasi != nil {
		existing.Sertifikasi = *req.Sertifikasi
	}

	// ✅ FIX UTAMA DI SINI
	if err := s.repo.Update(id, *existing); err != nil {
		return err
	}

	cacheDelete(s.rc, keyDetail("sdm", id))
	cacheDelete(s.rc, keyList("sdm"))

	return nil
}

func (s *SdmCsirtService) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	cacheDelete(s.rc, keyDetail("sdm", id))
	cacheDelete(s.rc, keyList("sdm"))

	return nil
}
