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
	Create(m models.Responden) error
	GetDetailByUserID(userID string) (*models.RespondenDetail, error)
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
// CACHE (WAJIB SATU SOURCE)
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

// =======================
// CONSTRUCTOR
// =======================
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

// =======================
// CREATE
// =======================
func (s *RespondenService) Create(req dto.CreateRespondenRequest) (*dto.RespondenResponse, error) {

	if err := s.validator.ValidateCreate(req); err != nil {
		return nil, err
	}

	model := models.Responden{
		UserID:             strings.TrimSpace(req.UserID),
		NoTelepon:          strings.TrimSpace(req.NoTelepon),
		Sektor:             strings.TrimSpace(req.Sektor),
		SektorLainnya:      toStringPtr(strings.TrimSpace(req.SektorLainnya)),
		SertifikatTraining: strings.TrimSpace(req.SertifikatTraining),
	}

	if err := s.repo.Create(model); err != nil {
		return nil, err
	}

	data, err := s.repo.GetDetailByUserID(model.UserID)
	if err != nil {
		return nil, err
	}

	resp := s.toResponse(data)

	_ = s.cache.Del(s.ctx, "responden:all")

	return &resp, nil
}

// =======================
// GET ALL (CACHE)
// =======================
func (s *RespondenService) GetAll() ([]dto.RespondenResponse, error) {

	cacheKey := "responden:all"

	// CACHE HIT
	if val, ok, err := s.cache.Get(s.ctx, cacheKey); err == nil && ok {
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

	if b, err := json.Marshal(result); err == nil {
		_ = s.cache.Set(s.ctx, cacheKey, string(b), int(cacheTTL.Seconds()))
	}

	return result, nil
}

// =======================
// GET BY ID (CACHE)
// =======================
func (s *RespondenService) GetByID(id int) (*dto.RespondenResponse, error) {

	if id <= 0 {
		return nil, errors.New("id tidak valid")
	}

	cacheKey := "responden:id:" + strconv.Itoa(id)

	// CACHE HIT
	if val, ok, err := s.cache.Get(s.ctx, cacheKey); err == nil && ok {
		var cached dto.RespondenResponse
		if json.Unmarshal([]byte(val), &cached) == nil {
			return &cached, nil
		}
	}

	data, err := s.repo.GetDetailByID(id)
	if err != nil {
		return nil, err
	}

	resp := s.toResponse(data)

	if b, err := json.Marshal(resp); err == nil {
		_ = s.cache.Set(s.ctx, cacheKey, string(b), int(cacheTTL.Seconds()))
	}

	return &resp, nil
}

// =======================
// UPDATE
// =======================
func (s *RespondenService) Update(id int, req dto.UpdateRespondenRequest) (*dto.RespondenResponse, error) {

	if id <= 0 {
		return nil, errors.New("id tidak valid")
	}

	_, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.validator.ValidateUpdate(req); err != nil {
		return nil, err
	}

	model := models.Responden{
		NoTelepon:          strings.TrimSpace(req.NoTelepon),
		Sektor:             strings.TrimSpace(req.Sektor),
		SektorLainnya:      toStringPtr(strings.TrimSpace(req.SektorLainnya)),
		SertifikatTraining: strings.TrimSpace(req.SertifikatTraining),
	}

	if err := s.repo.Update(id, model); err != nil {
		return nil, err
	}

	updated, err := s.repo.GetDetailByID(id)
	if err != nil {
		return nil, err
	}

	resp := s.toResponse(updated)

	_ = s.cache.Del(s.ctx, "responden:all")
	_ = s.cache.Del(s.ctx, "responden:id:"+strconv.Itoa(id))

	return &resp, nil
}

// =======================
// HELPERS
// =======================
func toStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// =======================
// MAPPER
// =======================
func (s *RespondenService) toResponse(m *models.RespondenDetail) dto.RespondenResponse {

	return dto.RespondenResponse{
		ID:              m.ID,
		UserID:          m.UserID,
		NamaLengkap:     safeString(m.NamaLengkap),
		Email:           safeString(m.Email),
		Jabatan:         safeString(m.Jabatan),
		NamaPerusahaan:  safeString(m.NamaPerusahaan),
		PerusahaanID:    safeString(m.PerusahaanID),
		NoTelepon:       m.NoTelepon,
		Sektor:          m.Sektor,
		SektorLainnya:   safeString(m.SektorLainnya),
		SertifikatTraining: m.SertifikatTraining,
		CreatedAt:       m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       m.UpdatedAt.Format(time.RFC3339),
	}
}