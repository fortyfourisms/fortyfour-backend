package services

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"survey/internal/dto"
	"survey/internal/models"
)

const cacheTTL = 10 * time.Minute

// =======================
// REPOSITORY
// =======================
type RespondenRepositoryInterface interface {
	Create(m models.Responden) (int64, error) // ✅ return insert ID
	GetAllDetail() ([]models.RespondenDetail, error)
	GetDetailByID(id int) (*models.RespondenDetail, error)
	GetByID(id int) (*models.Responden, error)
	Update(id int, m models.Responden) error
}

// =======================
// VALIDATOR
// =======================
type Validator interface {
	ValidateCreate(dto.CreateRespondenRequest) error
	ValidateUpdate(dto.UpdateRespondenRequest) error
}

// =======================
// CACHE
// =======================
type CacheRepository interface {
	Get(ctx context.Context, key string) (string, bool, error)
	Set(ctx context.Context, key string, value string, ttlSeconds int) error
	Del(ctx context.Context, key string) error
}

// =======================
// SERVICE
// =======================
type RespondenService struct {
	repo      RespondenRepositoryInterface
	validator Validator
	cache     CacheRepository
	ctx       context.Context
}

func NewRespondenService(
	repo RespondenRepositoryInterface,
	v Validator,
	cache CacheRepository,
) *RespondenService {
	return &RespondenService{
		repo:      repo,
		validator: v,
		cache:     cache,
		ctx:       context.Background(),
	}
}

//
// =======================
// CREATE
// =======================
//
func (s *RespondenService) Create(req dto.CreateRespondenRequest) (*dto.RespondenResponse, error) {

	if err := s.validator.ValidateCreate(req); err != nil {
		return nil, err
	}

	model := models.Responden{
		IdPerusahaan:       strings.TrimSpace(req.IdPerusahaan),
		NamaLengkap:        strings.TrimSpace(req.NamaLengkap),
		Jabatan:            strings.TrimSpace(req.Jabatan),
		Email:              strings.TrimSpace(req.Email),
		NoTelepon:          strings.TrimSpace(req.NoTelepon),
		SertifikatTraining: strings.TrimSpace(req.SertifikatTraining),
	}

	// ✅ insert + ambil ID
	insertID, err := s.repo.Create(model)
	if err != nil {
		return nil, err
	}

	// ✅ ambil data berdasarkan ID
	data, err := s.repo.GetDetailByID(int(insertID))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, errors.New("data tidak ditemukan setelah insert")
	}

	resp := s.toResponse(data)

	// invalidate cache
	_ = s.cache.Del(s.ctx, "responden:all")

	return &resp, nil
}

//
// =======================
// GET ALL (CACHE)
// =======================
//
func (s *RespondenService) GetAll() ([]dto.RespondenResponse, error) {

	cacheKey := "responden:all"

	// ✅ CACHE HIT
	if val, ok, err := s.cache.Get(s.ctx, cacheKey); err == nil && ok && val != "" && val != "null" {
		var cached []dto.RespondenResponse
		if json.Unmarshal([]byte(val), &cached) == nil {
			return cached, nil
		}
	}

	data, err := s.repo.GetAllDetail()
	if err != nil {
		return nil, err
	}

	result := make([]dto.RespondenResponse, 0, len(data))
	for i := range data {
		result = append(result, s.toResponse(&data[i]))
	}

	// ✅ hanya cache kalau ada data
	if len(result) > 0 {
		if b, err := json.Marshal(result); err == nil {
			_ = s.cache.Set(s.ctx, cacheKey, string(b), int(cacheTTL.Seconds()))
		}
	}

	return result, nil
}

//
// =======================
// GET BY ID (CACHE)
// =======================
//
func (s *RespondenService) GetByID(id int) (*dto.RespondenResponse, error) {

	if id <= 0 {
		return nil, errors.New("id tidak valid")
	}

	cacheKey := "responden:id:" + strconv.Itoa(id)

	// ✅ CACHE HIT
	if val, ok, err := s.cache.Get(s.ctx, cacheKey); err == nil && ok && val != "" && val != "null" {
		var cached dto.RespondenResponse
		if json.Unmarshal([]byte(val), &cached) == nil {
			return &cached, nil
		}
	}

	data, err := s.repo.GetDetailByID(id)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, errors.New("data tidak ditemukan")
	}

	resp := s.toResponse(data)

	// cache
	if b, err := json.Marshal(resp); err == nil {
		_ = s.cache.Set(s.ctx, cacheKey, string(b), int(cacheTTL.Seconds()))
	}

	return &resp, nil
}

//
// =======================
// UPDATE
// =======================
//
func (s *RespondenService) Update(id int, req dto.UpdateRespondenRequest) (*dto.RespondenResponse, error) {

	if id <= 0 {
		return nil, errors.New("id tidak valid")
	}

	// cek exist
	_, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.validator.ValidateUpdate(req); err != nil {
		return nil, err
	}

	model := models.Responden{
		IdPerusahaan:       strings.TrimSpace(req.IdPerusahaan),
		NamaLengkap:        strings.TrimSpace(req.NamaLengkap),
		Jabatan:            strings.TrimSpace(req.Jabatan),
		Email:              strings.TrimSpace(req.Email),
		NoTelepon:          strings.TrimSpace(req.NoTelepon),
		SertifikatTraining: strings.TrimSpace(req.SertifikatTraining),
	}

	if err := s.repo.Update(id, model); err != nil {
		return nil, err
	}

	updated, err := s.repo.GetDetailByID(id)
	if err != nil {
		return nil, err
	}

	if updated == nil {
		return nil, errors.New("data tidak ditemukan setelah update")
	}

	resp := s.toResponse(updated)

	// invalidate cache
	_ = s.cache.Del(s.ctx, "responden:all")
	_ = s.cache.Del(s.ctx, "responden:id:"+strconv.Itoa(id))

	return &resp, nil
}

//
// =======================
// MAPPER
// =======================
//
func (s *RespondenService) toResponse(m *models.RespondenDetail) dto.RespondenResponse {

	return dto.RespondenResponse{
		ID:                 m.ID,
		IdPerusahaan:       m.IdPerusahaan,
		NamaLengkap:        m.NamaLengkap,
		Jabatan:            m.Jabatan,
		Email:              m.Email,
		NoTelepon:          m.NoTelepon,
		SertifikatTraining: m.SertifikatTraining,
		NamaPerusahaan:     safeString(m.NamaPerusahaan),
		NamaSubSektor:      safeString(m.NamaSubSektor),
		NamaSektor:         safeString(m.NamaSektor),
		CreatedAt:          m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          m.UpdatedAt.Format(time.RFC3339),
	}
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}