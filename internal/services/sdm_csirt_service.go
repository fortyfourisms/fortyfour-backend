package services

import (
	"context"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/dto/dto_event"
	"fortyfour-backend/internal/rabbitmq"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/pkg/cache"
	"time"

	"github.com/google/uuid"
)

type SdmCsirtServiceInterface interface {
	Create(req dto.CreateSdmCsirtRequest) (string, error)
	GetAll() ([]dto.SdmCsirtResponse, error)
	GetByID(id string) (*dto.SdmCsirtResponse, error)
	GetByCsirt(idCsirt string) ([]dto.SdmCsirtResponse, error)
	Update(id string, req dto.UpdateSdmCsirtRequest) error
	Delete(id string) error
}

type SdmCsirtService struct {
	repo     repository.SdmCsirtRepositoryInterface
	rc       cache.RedisInterface
	producer *rabbitmq.Producer
}

func NewSdmCsirtService(repo repository.SdmCsirtRepositoryInterface, rc cache.RedisInterface, producer *rabbitmq.Producer) *SdmCsirtService {
	return &SdmCsirtService{
		repo:     repo,
		rc:       rc,
		producer: producer,
	}
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
	if req.IdCsirt != nil {
		cacheDelete(s.rc, "sdm_csirt_"+*req.IdCsirt)
	}

	if s.producer != nil {
		go func() {
			var idCsirt string
			if result.Csirt != nil {
				idCsirt = result.Csirt.ID
			}

			event := dto_event.SdmCsirtCreatedEvent{
				ID:                result.ID,
				IdCsirt:           idCsirt,
				NamaPersonel:      result.NamaPersonel,
				JabatanCsirt:      result.JabatanCsirt,
				JabatanPerusahaan: result.JabatanPerusahaan,
				Skill:             result.Skill,
				Sertifikasi:       result.Sertifikasi,
				CreatedAt:         time.Now(),
			}
			s.producer.PublishSdmCsirtCreated(context.Background(), event)
		}()
	}

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

func (s *SdmCsirtService) GetByCsirt(idCsirt string) ([]dto.SdmCsirtResponse, error) {
	key := "sdm_csirt_" + idCsirt
	var result []dto.SdmCsirtResponse

	if cacheGet(s.rc, key, &result) {
		return result, nil
	}

	data, err := s.repo.GetByCsirt(idCsirt)
	if err != nil {
		return nil, err
	}

	cacheSet(s.rc, key, data, TTLList)
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

	if err := s.repo.Update(id, *existing); err != nil {
		return err
	}

	cacheDelete(s.rc, keyDetail("sdm", id))
	cacheDelete(s.rc, keyList("sdm"))
	if existing.Csirt != nil {
		cacheDelete(s.rc, "sdm_csirt_"+existing.Csirt.ID)
	}

	if s.producer != nil {
		go func() {
			var idCsirt string
			if existing.Csirt != nil {
				idCsirt = existing.Csirt.ID
			}

			event := dto_event.SdmCsirtUpdatedEvent{
				ID:                existing.ID,
				IdCsirt:           idCsirt,
				NamaPersonel:      existing.NamaPersonel,
				JabatanCsirt:      existing.JabatanCsirt,
				JabatanPerusahaan: existing.JabatanPerusahaan,
				Skill:             existing.Skill,
				Sertifikasi:       existing.Sertifikasi,
				UpdatedAt:         time.Now(),
			}
			s.producer.PublishSdmCsirtUpdated(context.Background(), event)
		}()
	}

	return nil
}

func (s *SdmCsirtService) Delete(id string) error {
	// Ambil data dulu untuk invalidate cache per csirt
	existing, _ := s.repo.GetByID(id)

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	cacheDelete(s.rc, keyDetail("sdm", id))
	cacheDelete(s.rc, keyList("sdm"))
	if existing != nil && existing.Csirt != nil {
		cacheDelete(s.rc, "sdm_csirt_"+existing.Csirt.ID)
	}

	if s.producer != nil {
		go func() {
			var idCsirt string
			if existing != nil && existing.Csirt != nil {
				idCsirt = existing.Csirt.ID
			}

			event := dto_event.SdmCsirtDeletedEvent{
				ID:        id,
				IdCsirt:   idCsirt,
				DeletedAt: time.Now(),
			}
			s.producer.PublishSdmCsirtDeleted(context.Background(), event)
		}()
	}

	return nil
}
