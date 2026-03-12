package services

import (
	"errors"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/pkg/cache"
	"strings"

	"github.com/google/uuid"
)

type PerusahaanServiceInterface interface {
	GetAll() ([]dto.PerusahaanResponse, error)
	GetByID(id string) (*dto.PerusahaanResponse, error)
	GetByNama(nama string) (*dto.PerusahaanResponse, error)
	Create(req dto.CreatePerusahaanRequest) (*dto.PerusahaanResponse, error)
	Update(id string, req dto.UpdatePerusahaanRequest) (*dto.PerusahaanResponse, error)
	Delete(id string) error
}

type PerusahaanService struct {
	repo          repository.PerusahaanRepositoryInterface
	subSektorRepo repository.SubSektorRepositoryInterface
	rc            cache.RedisInterface
}

func NewPerusahaanService(repo repository.PerusahaanRepositoryInterface, subSektorRepo repository.SubSektorRepositoryInterface, rc cache.RedisInterface) *PerusahaanService {
	return &PerusahaanService{
		repo:          repo,
		subSektorRepo: subSektorRepo,
		rc:            rc,
	}
}

// Can be called from admin (full data) OR from registration (minimal data)
func (s *PerusahaanService) Create(req dto.CreatePerusahaanRequest) (*dto.PerusahaanResponse, error) {
	// Validasi nama perusahaan (WAJIB)
	if req.NamaPerusahaan == nil || strings.TrimSpace(*req.NamaPerusahaan) == "" {
		return nil, errors.New("nama_perusahaan wajib diisi")
	}

	// id_sub_sektor now OPTIONAL - only validate if provided
	if req.IDSubSektor != nil && strings.TrimSpace(*req.IDSubSektor) != "" {
		// Cek apakah sub sektor exists (hanya jika diisi)
		_, err := s.subSektorRepo.GetByID(*req.IDSubSektor)
		if err != nil {
			return nil, errors.New("sub sektor tidak ditemukan")
		}
	}
	// If id_sub_sektor is nil or empty, it will be saved as NULL in database

	// Generate ID dan simpan
	id := uuid.New().String()
	if err := s.repo.Create(req, id); err != nil {
		return nil, err
	}

	result, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	cacheSet(s.rc, keyDetail("perusahaan", id), result, TTLDetail)
	cacheDelete(s.rc, keyList("perusahaan"))

	return result, nil
}

func (s *PerusahaanService) GetAll() ([]dto.PerusahaanResponse, error) {
	key := keyList("perusahaan")
	var result []dto.PerusahaanResponse
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

func (s *PerusahaanService) GetByID(id string) (*dto.PerusahaanResponse, error) {
	key := keyDetail("perusahaan", id)
	var result dto.PerusahaanResponse
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

func (s *PerusahaanService) GetByNama(nama string) (*dto.PerusahaanResponse, error) {
	return s.repo.GetByNama(nama)
}

func (s *PerusahaanService) Update(id string, req dto.UpdatePerusahaanRequest) (*dto.PerusahaanResponse, error) {
	perusahaan, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update field jika ada
	if req.NamaPerusahaan != nil {
		perusahaan.NamaPerusahaan = *req.NamaPerusahaan
	}

	// Only validate id_sub_sektor if provided
	if req.IDSubSektor != nil && strings.TrimSpace(*req.IDSubSektor) != "" {
		// Validasi sub sektor saat update (hanya jika diisi)
		subSektor, err := s.subSektorRepo.GetByID(*req.IDSubSektor)
		if err != nil {
			return nil, errors.New("sub sektor tidak ditemukan")
		}
		// Update SubSektor object
		perusahaan.SubSektor = subSektor
	}
	// If id_sub_sektor is nil or empty, keep existing value

	if req.Alamat != nil {
		perusahaan.Alamat = *req.Alamat
	}
	if req.Telepon != nil {
		perusahaan.Telepon = *req.Telepon
	}
	if req.Email != nil {
		perusahaan.Email = *req.Email
	}
	if req.Website != nil {
		perusahaan.Website = *req.Website
	}
	if req.Photo != nil {
		perusahaan.Photo = *req.Photo
	}

	if err := s.repo.Update(id, *perusahaan); err != nil {
		return nil, err
	}

	cacheDelete(s.rc, keyDetail("perusahaan", id))
	cacheDelete(s.rc, keyList("perusahaan"))

	return perusahaan, nil
}

func (s *PerusahaanService) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	cacheDelete(s.rc, keyDetail("perusahaan", id))
	cacheDelete(s.rc, keyList("perusahaan"))

	return nil
}
