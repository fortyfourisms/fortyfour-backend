package services

import (
	"context"
	"fmt"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/dto/dto_event"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/rabbitmq"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/pkg/cache"

	"github.com/google/uuid"
)

type CsirtServiceInterface interface {
	GetAll() ([]dto.CsirtResponse, error)
	GetByID(id string) (*dto.CsirtResponse, error)
	GetByPerusahaan(idPerusahaan string) ([]dto.CsirtResponse, error)
	Create(req dto.CreateCsirtRequest) (*models.Csirt, error)
	Update(id string, req dto.UpdateCsirtRequest) (*models.Csirt, error)
	Delete(id string) error
}

type CsirtService struct {
	repo     repository.CsirtRepositoryInterface
	rc       cache.RedisInterface
	producer *rabbitmq.Producer
}

func NewCsirtService(repo repository.CsirtRepositoryInterface, rc cache.RedisInterface, producer *rabbitmq.Producer) *CsirtService {
	return &CsirtService{repo: repo, rc: rc, producer: producer}
}

func (s *CsirtService) Create(req dto.CreateCsirtRequest) (*models.Csirt, error) {
	exists, err := s.repo.ExistsByPerusahaan(req.IdPerusahaan)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("perusahaan ini sudah memiliki data CSIRT")
	}

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
	cacheDelete(s.rc, "csirt:perusahaan:"+req.IdPerusahaan)

	// Publish CsirtCreated event
	if s.producer != nil {
		go func() {
			tglReg := ""
			if result.TanggalRegistrasi != nil {
				tglReg = *result.TanggalRegistrasi
			}
			event := dto_event.CsirtCreatedEvent{
				ID:                result.ID,
				IdPerusahaan:      result.IdPerusahaan,
				NamaCsirt:         result.NamaCsirt,
				WebCsirt:          result.WebCsirt,
				TanggalRegistrasi: tglReg,
				CreatedAt:         time.Now(),
			}
			s.producer.PublishCsirtCreated(context.Background(), event)
		}()
	}

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

func (s *CsirtService) GetByPerusahaan(idPerusahaan string) ([]dto.CsirtResponse, error) {
	key := "csirt:perusahaan:" + idPerusahaan
	var result []dto.CsirtResponse
	if cacheGet(s.rc, key, &result) {
		return result, nil
	}

	data, err := s.repo.GetByPerusahaan(idPerusahaan)
	if err != nil {
		return nil, err
	}

	cacheSet(s.rc, key, data, TTLList)
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
		c.PhotoCsirt = req.PhotoCsirt
	}
	if req.FileRFC2350 != nil {
		c.FileRFC2350 = req.FileRFC2350
	}
	if req.FilePublicKeyPGP != nil {
		c.FilePublicKeyPGP = req.FilePublicKeyPGP
	}
	if req.FileStr != nil {
		c.FileStr = req.FileStr
	}
	if req.TanggalRegistrasi != nil {
		c.TanggalRegistrasi = req.TanggalRegistrasi
	}
	if req.TanggalKadaluarsa != nil {
		c.TanggalKadaluarsa = req.TanggalKadaluarsa
	}
	if req.TanggalRegistrasiUlang != nil {
		c.TanggalRegistrasiUlang = req.TanggalRegistrasiUlang
	}

	if err := s.repo.Update(id, *c); err != nil {
		return nil, err
	}

	cacheDelete(s.rc, keyDetail("csirt", id))
	cacheDelete(s.rc, keyList("csirt"))
	cacheDelete(s.rc, "csirt:perusahaan:"+c.IdPerusahaan)

	// Publish CsirtUpdated event
	if s.producer != nil {
		go func() {
			tglReg := ""
			if c.TanggalRegistrasi != nil {
				tglReg = *c.TanggalRegistrasi
			}
			event := dto_event.CsirtUpdatedEvent{
				ID:                c.ID,
				IdPerusahaan:      c.IdPerusahaan,
				NamaCsirt:         c.NamaCsirt,
				WebCsirt:          c.WebCsirt,
				TanggalRegistrasi: tglReg,
				UpdatedAt:         time.Now(),
			}
			s.producer.PublishCsirtUpdated(context.Background(), event)
		}()
	}

	return c, nil
}

func (s *CsirtService) Delete(id string) error {
	// Ambil data dulu untuk invalidate cache per perusahaan
	existing, _ := s.repo.GetByID(id)

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	cacheDelete(s.rc, keyDetail("csirt", id))
	cacheDelete(s.rc, keyList("csirt"))
	if existing != nil {
		cacheDelete(s.rc, "csirt:perusahaan:"+existing.IdPerusahaan)
	}

	// Publish CsirtDeleted event
	if s.producer != nil {
		go func() {
			event := dto_event.CsirtDeletedEvent{
				ID:        id,
				DeletedAt: time.Now(),
			}
			s.producer.PublishCsirtDeleted(context.Background(), event)
		}()
	}

	return nil
}
