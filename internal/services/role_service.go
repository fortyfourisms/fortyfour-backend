package services

import (
	"context"
	"database/sql"
	"errors"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/dto/dto_event"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/rabbitmq"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/pkg/cache"
	"time"
)

type RoleService struct {
	repo     repository.RoleRepository
	rc       cache.RedisInterface
	producer *rabbitmq.Producer
}

func NewRoleService(repo repository.RoleRepository, rc cache.RedisInterface, producer *rabbitmq.Producer) *RoleService {
	return &RoleService{
		repo:     repo,
		rc:       rc,
		producer: producer,
	}
}

func (s *RoleService) Create(req dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	ctx := context.Background()

	// Check if role name already exists
	existingRole, err := s.repo.GetByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if existingRole != nil {
		return nil, errors.New("role name already exists")
	}

	role := &models.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.repo.Create(ctx, role); err != nil {
		return nil, err
	}

	result := s.toResponse(role)
	cacheSet(s.rc, keyDetail("role", role.ID), result, TTLDetail)
	cacheDelete(s.rc, keyList("role"))

	if s.producer != nil {
		go func() {
			event := dto_event.RoleCreatedEvent{
				ID:          result.ID,
				Name:        result.Name,
				Description: result.Description,
				CreatedAt:   time.Now(),
			}
			s.producer.PublishRoleCreated(context.Background(), event)
		}()
	}

	return result, nil
}

func (s *RoleService) GetByID(id string) (*dto.RoleResponse, error) {
	key := keyDetail("role", id)
	var result dto.RoleResponse
	if cacheGet(s.rc, key, &result) {
		return &result, nil
	}

	ctx := context.Background()
	role, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	data := s.toResponse(role)
	cacheSet(s.rc, key, data, TTLDetail)
	return data, nil
}

func (s *RoleService) GetAll() ([]*dto.RoleResponse, error) {
	key := keyList("role")
	var result []*dto.RoleResponse
	if cacheGet(s.rc, key, &result) {
		return result, nil
	}

	ctx := context.Background()
	roles, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.RoleResponse, 0, len(roles))
	for _, role := range roles {
		responses = append(responses, s.toResponse(role))
	}

	cacheSet(s.rc, key, responses, TTLList)
	return responses, nil
}

func (s *RoleService) Update(id string, req dto.UpdateRoleRequest) (*dto.RoleResponse, error) {
	ctx := context.Background()

	role, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	// Check if new name already exists (if name is being changed)
	if req.Name != "" && req.Name != role.Name {
		existingRole, err := s.repo.GetByName(ctx, req.Name)
		if err != nil {
			return nil, err
		}
		if existingRole != nil {
			return nil, errors.New("role name already exists")
		}
		role.Name = req.Name
	}

	if req.Description != "" {
		role.Description = req.Description
	}

	if err := s.repo.Update(ctx, role); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	cacheDelete(s.rc, keyDetail("role", id))
	cacheDelete(s.rc, keyList("role"))

	result := s.toResponse(role)

	if s.producer != nil {
		go func() {
			event := dto_event.RoleUpdatedEvent{
				ID:          result.ID,
				Name:        result.Name,
				Description: result.Description,
				UpdatedAt:   time.Now(),
			}
			s.producer.PublishRoleUpdated(context.Background(), event)
		}()
	}

	return result, nil
}

func (s *RoleService) Delete(id string) error {
	ctx := context.Background()

	err := s.repo.Delete(ctx, id)
	if err == sql.ErrNoRows {
		return errors.New("role not found")
	}
	if err != nil {
		return err
	}

	cacheDelete(s.rc, keyDetail("role", id))
	cacheDelete(s.rc, keyList("role"))

	if s.producer != nil {
		go func() {
			event := dto_event.RoleDeletedEvent{
				ID:        id,
				DeletedAt: time.Now(),
			}
			s.producer.PublishRoleDeleted(context.Background(), event)
		}()
	}

	return nil
}

func (s *RoleService) toResponse(role *models.Role) *dto.RoleResponse {
	return &dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
