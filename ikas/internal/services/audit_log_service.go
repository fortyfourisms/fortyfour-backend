package services

import (
	"encoding/json"
	"fmt"
	"ikas/internal/dto"
	"ikas/internal/models"
	"ikas/internal/repository"
	"ikas/pkg/cache"
	"math"
	"time"

	"github.com/rollbar/rollbar-go"
)

// AuditLogService handles business logic for IKAS audit logs.
type AuditLogService struct {
	repo  repository.AuditLogRepositoryInterface
	cache cache.RedisInterface
}

func NewAuditLogService(repo repository.AuditLogRepositoryInterface, cache cache.RedisInterface) *AuditLogService {
	return &AuditLogService{repo: repo, cache: cache}
}

// GetAuditLogs returns a paginated list of audit logs, optionally filtered by ikas_id.
// Results are served from Redis cache when available (Cache-Aside pattern).
func (s *AuditLogService) GetAuditLogs(req dto.AuditLogListRequest) (*dto.PaginatedAuditLogResponse, error) {
	offset := (req.Page - 1) * req.Limit

	// Build a deterministic cache key that encodes all query dimensions.
	var cacheKey string
	if req.IkasID != "" {
		cacheKey = fmt.Sprintf("%s%s:%d:%d", cache.CacheKeyPrefixAuditLogsByIkas, req.IkasID, req.Page, req.Limit)
	} else {
		cacheKey = fmt.Sprintf("%s%d:%d", cache.CacheKeyPrefixAuditLogs, req.Page, req.Limit)
	}

	// --- Cache-Aside: Read ---
	if s.cache != nil {
		if cached, err := s.cache.Get(cacheKey); err == nil && cached != "" {
			var result dto.PaginatedAuditLogResponse
			if err := json.Unmarshal([]byte(cached), &result); err == nil {
				return &result, nil
			}
		}
	}

	// --- Cache Miss: Query DB ---
	var (
		rawLogs []models.AuditLog
		total   int
		err     error
	)

	if req.IkasID != "" {
		rawLogs, total, err = s.repo.GetAuditLogsByIkasID(req.IkasID, offset, req.Limit)
	} else {
		rawLogs, total, err = s.repo.GetAuditLogs(offset, req.Limit)
	}
	if err != nil {
		return nil, err
	}

	logs := mapAuditLogs(rawLogs)

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))
	result := &dto.PaginatedAuditLogResponse{
		Data:       logs,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}

	// --- Cache-Aside: Write (async to not block response) ---
	if s.cache != nil {
		go func(key string, data *dto.PaginatedAuditLogResponse) {
			jsonData, err := json.Marshal(data)
			if err == nil {
				if err := s.cache.Set(key, string(jsonData), cache.AuditLogsCacheExpiration); err != nil {
					rollbar.Error(fmt.Errorf("audit log redis cache write error: %v", err))
				}
			}
		}(cacheKey, result)
	}

	return result, nil
}

// mapAuditLogs converts a slice of models.AuditLog to a slice of dto.AuditLogResponse.
func mapAuditLogs(rawLogs []models.AuditLog) []dto.AuditLogResponse {
	result := make([]dto.AuditLogResponse, 0, len(rawLogs))
	for _, log := range rawLogs {
		result = append(result, dto.AuditLogResponse{
			ID:     log.ID,
			IkasID: log.IkasID,
			User: dto.UserAuditLogResponse{
				ID:   log.UserID,
				Name: log.UserName,
			},
			Action:    log.Action,
			Changes:   log.Changes,
			CreatedAt: log.CreatedAt.Format(time.RFC3339),
		})
	}
	return result
}
