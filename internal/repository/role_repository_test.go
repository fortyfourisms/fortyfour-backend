package repository

import (
	"context"
	"database/sql"
	"errors"
	"fortyfour-backend/internal/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRoleTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, RoleRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	
	repo := NewRoleRepository(db)
	return db, mock, repo
}

func TestNewRoleRepository(t *testing.T) {
	db, _, _ := setupRoleTest(t)
	defer db.Close()

	repo := NewRoleRepository(db)
	assert.NotNil(t, repo)
}

func TestRoleRepository_Create(t *testing.T) {
	db, mock, repo := setupRoleTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		role := &models.Role{
			Name:        "Admin",
			Description: "Administrator role",
		}

		mock.ExpectExec("INSERT INTO roles").
			WithArgs(sqlmock.AnyArg(), role.Name, role.Description, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(ctx, role)
		
		assert.NoError(t, err)
		assert.NotEmpty(t, role.ID) // ID should be generated
		assert.False(t, role.CreatedAt.IsZero()) // CreatedAt should be set
		assert.False(t, role.UpdatedAt.IsZero()) // UpdatedAt should be set
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		ctx := context.Background()
		role := &models.Role{
			Name:        "Admin",
			Description: "Administrator role",
		}

		mock.ExpectExec("INSERT INTO roles").
			WithArgs(sqlmock.AnyArg(), role.Name, role.Description, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(errors.New("database error"))

		err := repo.Create(ctx, role)
		
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate name error", func(t *testing.T) {
		ctx := context.Background()
		role := &models.Role{
			Name:        "Admin",
			Description: "Administrator role",
		}

		mock.ExpectExec("INSERT INTO roles").
			WithArgs(sqlmock.AnyArg(), role.Name, role.Description, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(errors.New("duplicate entry for key 'name'"))

		err := repo.Create(ctx, role)
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRoleRepository_GetByID(t *testing.T) {
	db, mock, repo := setupRoleTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		id := "role-123"
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "name", "description", "created_at", "updated_at",
		}).AddRow(id, "Admin", "Administrator role", now, now)

		mock.ExpectQuery("SELECT (.+) FROM roles WHERE id = ?").
			WithArgs(id).
			WillReturnRows(rows)

		result, err := repo.GetByID(ctx, id)
		
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, id, result.ID)
		assert.Equal(t, "Admin", result.Name)
		assert.Equal(t, "Administrator role", result.Description)
		assert.Equal(t, now, result.CreatedAt)
		assert.Equal(t, now, result.UpdatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found - returns nil", func(t *testing.T) {
		ctx := context.Background()
		id := "role-nonexistent"

		mock.ExpectQuery("SELECT (.+) FROM roles WHERE id = ?").
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)

		result, err := repo.GetByID(ctx, id)
		
		assert.NoError(t, err) // No error, just nil result
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		ctx := context.Background()
		id := "role-123"

		mock.ExpectQuery("SELECT (.+) FROM roles WHERE id = ?").
			WithArgs(id).
			WillReturnError(errors.New("database connection error"))

		result, err := repo.GetByID(ctx, id)
		
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "database connection error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		ctx := context.Background()
		id := "role-123"

		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(id)

		mock.ExpectQuery("SELECT (.+) FROM roles WHERE id = ?").
			WithArgs(id).
			WillReturnRows(rows)

		result, err := repo.GetByID(ctx, id)
		
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRoleRepository_GetAll(t *testing.T) {
	db, mock, repo := setupRoleTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "name", "description", "created_at", "updated_at",
		}).
			AddRow("role-1", "Admin", "Administrator role", now, now).
			AddRow("role-2", "User", "Regular user role", now, now).
			AddRow("role-3", "Guest", "Guest role", now, now)

		mock.ExpectQuery("SELECT (.+) FROM roles ORDER BY created_at DESC").
			WillReturnRows(rows)

		result, err := repo.GetAll(ctx)
		
		assert.NoError(t, err)
		require.Len(t, result, 3)
		
		// Check first record
		assert.Equal(t, "role-1", result[0].ID)
		assert.Equal(t, "Admin", result[0].Name)
		assert.Equal(t, "Administrator role", result[0].Description)
		
		// Check second record
		assert.Equal(t, "role-2", result[1].ID)
		assert.Equal(t, "User", result[1].Name)
		
		// Check third record
		assert.Equal(t, "role-3", result[2].ID)
		assert.Equal(t, "Guest", result[2].Name)
		
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		ctx := context.Background()

		rows := sqlmock.NewRows([]string{
			"id", "name", "description", "created_at", "updated_at",
		})

		mock.ExpectQuery("SELECT (.+) FROM roles ORDER BY created_at DESC").
			WillReturnRows(rows)

		result, err := repo.GetAll(ctx)
		
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		ctx := context.Background()

		mock.ExpectQuery("SELECT (.+) FROM roles ORDER BY created_at DESC").
			WillReturnError(errors.New("query error"))

		result, err := repo.GetAll(ctx)
		
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "query error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		ctx := context.Background()

		rows := sqlmock.NewRows([]string{"id"}).
			AddRow("role-1")

		mock.ExpectQuery("SELECT (.+) FROM roles ORDER BY created_at DESC").
			WillReturnRows(rows)

		result, err := repo.GetAll(ctx)
		
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows close error", func(t *testing.T) {
		ctx := context.Background()
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "name", "description", "created_at", "updated_at",
		}).
			AddRow("role-1", "Admin", "Administrator role", now, now).
			CloseError(errors.New("close error"))

		mock.ExpectQuery("SELECT (.+) FROM roles ORDER BY created_at DESC").
			WillReturnRows(rows)

		_, err := repo.GetAll(ctx)
		
		// The function returns the result and checks rows.Err() at the end
		// CloseError sets an error that will be returned by rows.Err()
		assert.Error(t, err)
		assert.Equal(t, "close error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRoleRepository_Update(t *testing.T) {
	db, mock, repo := setupRoleTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		role := &models.Role{
			ID:          "role-123",
			Name:        "Super Admin",
			Description: "Super administrator role",
		}

		mock.ExpectExec("UPDATE roles SET name = \\?, description = \\?, updated_at = \\? WHERE id = \\?").
			WithArgs(role.Name, role.Description, sqlmock.AnyArg(), role.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(ctx, role)
		
		assert.NoError(t, err)
		assert.False(t, role.UpdatedAt.IsZero()) // UpdatedAt should be set
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found - no rows affected", func(t *testing.T) {
		ctx := context.Background()
		role := &models.Role{
			ID:          "role-nonexistent",
			Name:        "Super Admin",
			Description: "Super administrator role",
		}

		mock.ExpectExec("UPDATE roles SET name = \\?, description = \\?, updated_at = \\? WHERE id = \\?").
			WithArgs(role.Name, role.Description, sqlmock.AnyArg(), role.ID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Update(ctx, role)
		
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		ctx := context.Background()
		role := &models.Role{
			ID:          "role-123",
			Name:        "Super Admin",
			Description: "Super administrator role",
		}

		mock.ExpectExec("UPDATE roles SET name = \\?, description = \\?, updated_at = \\? WHERE id = \\?").
			WithArgs(role.Name, role.Description, sqlmock.AnyArg(), role.ID).
			WillReturnError(errors.New("update error"))

		err := repo.Update(ctx, role)
		
		assert.Error(t, err)
		assert.Equal(t, "update error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows affected error", func(t *testing.T) {
		ctx := context.Background()
		role := &models.Role{
			ID:          "role-123",
			Name:        "Super Admin",
			Description: "Super administrator role",
		}

		mock.ExpectExec("UPDATE roles SET name = \\?, description = \\?, updated_at = \\? WHERE id = \\?").
			WithArgs(role.Name, role.Description, sqlmock.AnyArg(), role.ID).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))

		err := repo.Update(ctx, role)
		
		assert.Error(t, err)
		assert.Equal(t, "rows affected error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRoleRepository_Delete(t *testing.T) {
	db, mock, repo := setupRoleTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		id := "role-123"

		mock.ExpectExec("DELETE FROM roles WHERE id = \\?").
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(ctx, id)
		
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found - no rows affected", func(t *testing.T) {
		ctx := context.Background()
		id := "role-nonexistent"

		mock.ExpectExec("DELETE FROM roles WHERE id = \\?").
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete(ctx, id)
		
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		ctx := context.Background()
		id := "role-123"

		mock.ExpectExec("DELETE FROM roles WHERE id = \\?").
			WithArgs(id).
			WillReturnError(errors.New("delete error"))

		err := repo.Delete(ctx, id)
		
		assert.Error(t, err)
		assert.Equal(t, "delete error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows affected error", func(t *testing.T) {
		ctx := context.Background()
		id := "role-123"

		mock.ExpectExec("DELETE FROM roles WHERE id = \\?").
			WithArgs(id).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))

		err := repo.Delete(ctx, id)
		
		assert.Error(t, err)
		assert.Equal(t, "rows affected error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRoleRepository_GetByName(t *testing.T) {
	db, mock, repo := setupRoleTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		name := "Admin"
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "name", "description", "created_at", "updated_at",
		}).AddRow("role-123", name, "Administrator role", now, now)

		mock.ExpectQuery("SELECT (.+) FROM roles WHERE name = ?").
			WithArgs(name).
			WillReturnRows(rows)

		result, err := repo.GetByName(ctx, name)
		
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "role-123", result.ID)
		assert.Equal(t, name, result.Name)
		assert.Equal(t, "Administrator role", result.Description)
		assert.Equal(t, now, result.CreatedAt)
		assert.Equal(t, now, result.UpdatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found - returns nil", func(t *testing.T) {
		ctx := context.Background()
		name := "NonExistentRole"

		mock.ExpectQuery("SELECT (.+) FROM roles WHERE name = ?").
			WithArgs(name).
			WillReturnError(sql.ErrNoRows)

		result, err := repo.GetByName(ctx, name)
		
		assert.NoError(t, err) // No error, just nil result
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		ctx := context.Background()
		name := "Admin"

		mock.ExpectQuery("SELECT (.+) FROM roles WHERE name = ?").
			WithArgs(name).
			WillReturnError(errors.New("database connection error"))

		result, err := repo.GetByName(ctx, name)
		
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "database connection error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		ctx := context.Background()
		name := "Admin"

		rows := sqlmock.NewRows([]string{"id"}).
			AddRow("role-123")

		mock.ExpectQuery("SELECT (.+) FROM roles WHERE name = ?").
			WithArgs(name).
			WillReturnRows(rows)

		result, err := repo.GetByName(ctx, name)
		
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}