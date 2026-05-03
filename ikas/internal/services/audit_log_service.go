package services

import (
	"ikas/internal/dto"
	"ikas/internal/repository"
	"time"
)

type AuditLogService struct {
	repo repository.AuditLogRepositoryInterface
}

func NewAuditLogService(repo repository.AuditLogRepositoryInterface) *AuditLogService {
	return &AuditLogService{repo: repo}
}

func (s *AuditLogService) GetAuditLogsByIkasID(ikasID string) ([]dto.AuditLogResponse, error) {
	logs, err := s.repo.GetAuditLogsByIkasID(ikasID)
	if err != nil {
		return nil, err
	}

	var response []dto.AuditLogResponse
	for _, log := range logs {
		response = append(response, dto.AuditLogResponse{
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

	return response, nil
}

func (s *AuditLogService) GetAllAuditLogs() ([]dto.AuditLogResponse, error) {
	logs, err := s.repo.GetAllAuditLogs()
	if err != nil {
		return nil, err
	}

	var response []dto.AuditLogResponse
	for _, log := range logs {
		response = append(response, dto.AuditLogResponse{
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

	return response, nil
}
