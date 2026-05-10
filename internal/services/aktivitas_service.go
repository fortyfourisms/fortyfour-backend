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

func parseDate(dateStr string) (time.Time, error) {
	// Coba format YYYY-MM-DD
	t, err := time.Parse("2006-01-02", dateStr)
	if err == nil {
		return t, nil
	}
	// Coba format RFC3339 (ISO8601)
	t, err = time.Parse(time.RFC3339, dateStr)
	if err == nil {
		return t, nil
	}
	return time.Time{}, errors.New("format tanggal tidak valid. Gunakan YYYY-MM-DD atau format ISO8601")
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

	start, err := parseDate(req.TanggalMulai)
	if err != nil {
		return fmt.Errorf("tanggal_mulai: %v", err)
	}
	end, err := parseDate(req.TanggalSelesai)
	if err != nil {
		return fmt.Errorf("tanggal_selesai: %v", err)
	}

	if start.After(end) {
		return errors.New("tanggal_mulai tidak boleh melebihi tanggal_selesai")
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

func (s *AktivitasService) validateUpdate(id int, req *dto.UpdateAktivitasRequest) error {
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

	// Validasi Urutan Tanggal jika ada yang diupdate
	if req.TanggalMulai != nil || req.TanggalSelesai != nil {
		existing, err := s.repo.GetByID(id)
		if err != nil {
			return err
		}

		startStr := existing.TanggalMulai
		if req.TanggalMulai != nil {
			startStr = *req.TanggalMulai
		}

		endStr := existing.TanggalSelesai
		if req.TanggalSelesai != nil {
			endStr = *req.TanggalSelesai
		}

		start, err := parseDate(startStr)
		if err != nil {
			return fmt.Errorf("tanggal_mulai: %v", err)
		}
		end, err := parseDate(endStr)
		if err != nil {
			return fmt.Errorf("tanggal_selesai: %v", err)
		}

		if start.After(end) {
			return errors.New("tanggal_mulai tidak boleh melebihi tanggal_selesai")
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

	if err := s.validateUpdate(id, &req); err != nil {
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
