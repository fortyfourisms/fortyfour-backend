package repository

import (
	"context"
	"database/sql"
	"fortyfour-backend/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/rollbar/rollbar-go"
)

type RoleRepository interface {
	Create(ctx context.Context, role *models.Role) error
	GetByID(ctx context.Context, id string) (*models.Role, error)
	GetAll(ctx context.Context) ([]*models.Role, error)
	Update(ctx context.Context, role *models.Role) error
	Delete(ctx context.Context, id string) error
	GetByName(ctx context.Context, name string) (*models.Role, error)
}

type roleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *models.Role) error {
	role.ID = uuid.New().String()
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	query := `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		role.ID,
		role.Name,
		role.Description,
		role.CreatedAt,
		role.UpdatedAt,
	)

	return err
}

func (r *roleRepository) GetByID(ctx context.Context, id string) (*models.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE id = ?
	`

	role := &models.Role{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return role, nil
}

func (r *roleRepository) GetAll(ctx context.Context) ([]*models.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	roles := []*models.Role{}
	for rows.Next() {
		role := &models.Role{}
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, rows.Err()
}

func (r *roleRepository) Update(ctx context.Context, role *models.Role) error {
	role.UpdatedAt = time.Now()

	query := `
		UPDATE roles
		SET name = ?, description = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		role.Name,
		role.Description,
		role.UpdatedAt,
		role.ID,
	)

	if err != nil {
		rollbar.Error(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *roleRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM roles WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		rollbar.Error(err)
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*models.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE name = ?
	`

	role := &models.Role{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return role, nil
}
