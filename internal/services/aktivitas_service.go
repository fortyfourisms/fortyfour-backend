package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/dto/dto_event"
	"fortyfour-backend/internal/rabbitmq"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/pkg/cache"
	"strings"
	"time"
)

var AllowedJenisAktivitas = []string{"dinas", "workshop", "seminar", "webinar", "rapat koordinasi", "rapat kerja", "kunjungan kerja", "rapat evaluasi", "pelatihan/diklat", "bimbingan teknis"}

type AktivitasService struct {
	repo           repository.AktivitasRepositoryInterface
	perusahaanRepo repository.PerusahaanRepositoryInterface
	producer       *rabbitmq.Producer
	rc             cache.RedisInterface
}

func NewAktivitasService(
	repo repository.AktivitasRepositoryInterface,
	perusahaanRepo repository.PerusahaanRepositoryInterface,
	producer *rabbitmq.Producer,
	rc cache.RedisInterface,
) *AktivitasService {
	return &AktivitasService{
		repo:           repo,
		perusahaanRepo: perusahaanRepo,
		producer:       producer,
		rc:             rc,
	}
}

func (s *AktivitasService) GetAllowedJenis() []string {
	return AllowedJenisAktivitas
}

func (s *AktivitasService) validateCreate(req *dto.CreateAktivitasRequest) error {
	if req.PerusahaanID == "" {
		return errors.New("perusahaan_id tidak boleh kosong")
	}
	if req.Judul == "" {
		return errors.New("judul tidak boleh kosong")
	}
	if req.TanggalMulai == "" {
		return errors.New("tanggal_mulai tidak boleh kosong")
	}
	if req.TanggalSelesai == "" {
		return errors.New("tanggal_selesai tidak boleh kosong")
	}
	if len(req.JenisAktivitas) == 0 {
		return errors.New("jenis_aktivitas tidak boleh kosong")
	}

	validJenis := make(map[string]bool)
	for _, v := range AllowedJenisAktivitas {
		validJenis[v] = true
	}

	for _, j := range req.JenisAktivitas {
		if !validJenis[j] {
			return errors.New("jenis_aktivitas '" + j + "' tidak valid. Harus salah satu dari: " + strings.Join(AllowedJenisAktivitas, ", "))
		}
	}

	// Validasi perusahaan_id
	if _, err := s.perusahaanRepo.GetByID(req.PerusahaanID); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("perusahaan_id tidak ditemukan")
		}
		return fmt.Errorf("gagal memvalidasi perusahaan_id: %v", err)
	}

	return nil
}

func (s *AktivitasService) validateUpdate(req *dto.UpdateAktivitasRequest) error {
	if req.JenisAktivitas != nil {
		validJenis := make(map[string]bool)
		for _, v := range AllowedJenisAktivitas {
			validJenis[v] = true
		}

		for _, j := range *req.JenisAktivitas {
			if !validJenis[j] {
				return errors.New("jenis_aktivitas '" + j + "' tidak valid. Harus salah satu dari: " + strings.Join(AllowedJenisAktivitas, ", "))
			}
		}
	}

	if req.PerusahaanID != nil {
		if _, err := s.perusahaanRepo.GetByID(*req.PerusahaanID); err != nil {
			if err == sql.ErrNoRows {
				return errors.New("perusahaan_id tidak ditemukan")
			}
			return fmt.Errorf("gagal memvalidasi perusahaan_id: %v", err)
		}
	}

	return nil
}

func (s *AktivitasService) Create(req dto.CreateAktivitasRequest) (*dto.AktivitasResponse, error) {
	if err := s.validateCreate(&req); err != nil {
		return nil, err
	}

	if s.producer != nil {
		err := s.producer.PublishAktivitasCreated(context.Background(), dto_event.AktivitasCreatedEvent{
			Request:   req,
			CreatedAt: time.Now(),
		})
		if err != nil {
			return nil, err
		}
	}

	cacheDelete(s.rc, keyList("aktivitas"))
	cacheDelete(s.rc, fmt.Sprintf("aktivitas:perusahaan:%s", req.PerusahaanID))

	return nil, nil
}

func (s *AktivitasService) GetAll() ([]dto.AktivitasResponse, error) {
	key := keyList("aktivitas")
	var result []dto.AktivitasResponse
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

func (s *AktivitasService) GetByID(id int) (*dto.AktivitasResponse, error) {
	key := keyDetail("aktivitas", fmt.Sprintf("%d", id))
	var result dto.AktivitasResponse
	if cacheGet(s.rc, key, &result) {
		return &result, nil
	}

	data, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}

	cacheSet(s.rc, key, data, TTLDetail)
	return data, nil
}

func (s *AktivitasService) GetByPerusahaanID(perusahaanID string) ([]dto.AktivitasResponse, error) {
	key := fmt.Sprintf("aktivitas:perusahaan:%s", perusahaanID)
	var result []dto.AktivitasResponse
	if cacheGet(s.rc, key, &result) {
		return result, nil
	}

	result, err := s.repo.GetByPerusahaanID(perusahaanID)
	if err != nil {
		return nil, err
	}

	cacheSet(s.rc, key, result, TTLList)
	return result, nil
}

func (s *AktivitasService) Update(id int, req dto.UpdateAktivitasRequest) (*dto.AktivitasResponse, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}

	if err := s.validateUpdate(&req); err != nil {
		return nil, err
	}

	if s.producer != nil {
		err = s.producer.PublishAktivitasUpdated(context.Background(), dto_event.AktivitasUpdatedEvent{
			ID:        id,
			Request:   req,
			UpdatedAt: time.Now(),
		})
		if err != nil {
			return nil, err
		}
	}

	cacheDelete(s.rc, keyList("aktivitas"))
	cacheDelete(s.rc, keyDetail("aktivitas", fmt.Sprintf("%d", id)))
	cacheDelete(s.rc, fmt.Sprintf("aktivitas:perusahaan:%s", data.PerusahaanID))
	if req.PerusahaanID != nil {
		cacheDelete(s.rc, fmt.Sprintf("aktivitas:perusahaan:%s", *req.PerusahaanID))
	}

	return nil, nil
}

func (s *AktivitasService) Delete(id int) error {
	data, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("data tidak ditemukan")
		}
		return err
	}

	if s.producer != nil {
		err = s.producer.PublishAktivitasDeleted(context.Background(), dto_event.AktivitasDeletedEvent{
			ID:        id,
			DeletedAt: time.Now(),
		})
		if err != nil {
			return err
		}
	}

	cacheDelete(s.rc, keyList("aktivitas"))
	cacheDelete(s.rc, keyDetail("aktivitas", fmt.Sprintf("%d", id)))
	cacheDelete(s.rc, fmt.Sprintf("aktivitas:perusahaan:%s", data.PerusahaanID))

	return nil
}
