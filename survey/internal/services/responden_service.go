package services

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"survey/internal/dto"
	"survey/internal/models"
)

// CONFIG
const cacheTTL = 10 * time.Minute

// REPOSITORY
type RespondenRepositoryInterface interface {
	GetAllDetail() ([]models.RespondenDetail, error)
	GetDetailByID(id string) (*models.RespondenDetail, error)

	GetByUserID(userID string) (*models.RespondenDetail, error)

	Create(m models.Responden) (int64, error)
	UpsertByUserID(userID string, m models.Responden) error
	CanEditByUserID(userID string) (bool, string, error)
}

// VALIDATOR
type Validator interface {
	ValidateCreate(dto.CreateRespondenRequest) error
}

// CACHE
type CacheRepository interface {
	Get(ctx context.Context, key string) (string, bool, error)
	Set(ctx context.Context, key string, value string, ttlSeconds int) error
	Del(ctx context.Context, key string) error
}

// SERVICE
type RespondenService struct {
	repo      RespondenRepositoryInterface
	validator Validator
	cache     CacheRepository
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
	}
}

// USER FLOW

// ADAPTER
func (s *RespondenService) GetByUserID(userID string) (*dto.RespondenResponse, error) {
	return s.GetMe(userID)
}

// GET ME
func (s *RespondenService) GetMe(userID string) (*dto.RespondenResponse, error) {

	ctx := context.Background()
	cacheKey := "responden:user:" + userID

	if val, ok, _ := s.cache.Get(ctx, cacheKey); ok {
		var cached dto.RespondenResponse
		if json.Unmarshal([]byte(val), &cached) == nil {
			return &cached, nil
		}
	}

	data, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	resp := s.toResponse(data)

	s.setCache(ctx, cacheKey, resp)
	return &resp, nil
}

// UPSERT ME
func (s *RespondenService) UpsertByUserID(userID string, req dto.CreateRespondenRequest) (*dto.RespondenResponse, error) {
	canEdit, status, err := s.repo.CanEditByUserID(userID)
	if err != nil {
		return nil, err
	}
	if !canEdit {
		return nil, errors.New("survey sudah selesai dengan status " + status + ", ajukan request edit ke admin")
	}

	if err := s.validator.ValidateCreate(req); err != nil {
		return nil, err
	}

	model := models.Responden{
		UserID:       userID,
		IdPerusahaan: strings.TrimSpace(req.IdPerusahaan),
		NamaLengkap:  strings.TrimSpace(req.NamaLengkap),
		Jabatan:      strings.TrimSpace(req.Jabatan),
		Email:        strings.TrimSpace(req.Email),
		NoTelepon:    strings.TrimSpace(req.NoTelepon),
	}

	// FIX POINTER
	if req.SertifikatTraining != nil {
		val := strings.TrimSpace(*req.SertifikatTraining)
		model.SertifikatTraining = &val
	}

	err = s.repo.UpsertByUserID(userID, model)
	if err != nil {
		return nil, err
	}

	updated, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	resp := s.toResponse(updated)

	s.invalidateUserCache(userID, updated.ID)
	return &resp, nil
}

// ADMIN FLOW

// GET ALL
func (s *RespondenService) GetAll() ([]dto.RespondenResponse, error) {

	ctx := context.Background()
	cacheKey := "responden:all"

	if val, ok, _ := s.cache.Get(ctx, cacheKey); ok {
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

	s.setCache(ctx, cacheKey, result)
	return result, nil
}

// GET BY ID
func (s *RespondenService) GetByID(id string) (*dto.RespondenResponse, error) {

	ctx := context.Background()
	cacheKey := "responden:id:" + id

	if val, ok, _ := s.cache.Get(ctx, cacheKey); ok {
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

	s.setCache(ctx, cacheKey, resp)
	return &resp, nil
}

// CACHE

func (s *RespondenService) setCache(ctx context.Context, key string, data any) {
	b, err := json.Marshal(data)
	if err != nil {
		return
	}
	_ = s.cache.Set(ctx, key, string(b), int(cacheTTL.Seconds()))
}

func (s *RespondenService) invalidateUserCache(userID string, id string) {
	ctx := context.Background()

	_ = s.cache.Del(ctx, "responden:all")

	if userID != "" {
		_ = s.cache.Del(ctx, "responden:user:"+userID)
	}

	_ = s.cache.Del(ctx, "responden:id:"+id)
}

// MAPPER

func (s *RespondenService) toResponse(m *models.RespondenDetail) dto.RespondenResponse {

	return dto.RespondenResponse{
		ID:                 m.ID,
		UserID:             m.UserID,
		IdPerusahaan:       m.IdPerusahaan,
		NamaLengkap:        m.NamaLengkap,
		Jabatan:            m.Jabatan,
		Email:              m.Email,
		NoTelepon:          m.NoTelepon,
		SertifikatTraining: safeString(m.SertifikatTraining),

		NamaPerusahaan: safeString(m.NamaPerusahaan),
		NamaSubSektor:  safeString(m.NamaSubSektor),
		NamaSektor:     safeString(m.NamaSektor),

		CreatedAt: m.CreatedAt.Format(time.RFC3339),
		UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
	}
}

// HELPER
func safeString(s *string) *string {
	if s == nil {
		return nil
	}
	val := *s
	return &val
}