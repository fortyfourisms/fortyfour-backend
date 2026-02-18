package services

import (
	"context"
	"database/sql"
	"errors"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
)

type RoleService struct {
	repo repository.RoleRepository
}

func NewRoleService(repo repository.RoleRepository) *RoleService {
	return &RoleService{repo: repo}
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

	return s.toResponse(role), nil
}

func (s *RoleService) GetByID(id string) (*dto.RoleResponse, error) {
	ctx := context.Background()

	role, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	return s.toResponse(role), nil
}

func (s *RoleService) GetAll() ([]*dto.RoleResponse, error) {
	ctx := context.Background()

	roles, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.RoleResponse, 0, len(roles))
	for _, role := range roles {
		responses = append(responses, s.toResponse(role))
	}

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

	return s.toResponse(role), nil
}

func (s *RoleService) Delete(id string) error {
	ctx := context.Background()

	err := s.repo.Delete(ctx, id)
	if err == sql.ErrNoRows {
		return errors.New("role not found")
	}
	return err
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
