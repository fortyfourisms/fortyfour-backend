package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"fortyfour-backend/internal/dto"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupJabatanTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *JabatanRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewJabatanRepository(db)
	return db, mock, repo
}

func TestNewJabatanRepository(t *testing.T) {
	db, _, _ := setupJabatanTest(t)
	defer db.Close()

	repo := NewJabatanRepository(db)
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestJabatanRepository_Create(t *testing.T) {
	db, mock, repo := setupJabatanTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		namaJabatan := "Manager IT"
		req := dto.CreateJabatanRequest{
			NamaJabatan: &namaJabatan,
		}
		id := "001"

		mock.ExpectExec("INSERT INTO jabatan").
			WithArgs(id, req.NamaJabatan).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(req, id)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		namaJabatan := "Manager IT"
		req := dto.CreateJabatanRequest{
			NamaJabatan: &namaJabatan,
		}
		id := "001"

		mock.ExpectExec("INSERT INTO jabatan").
			WithArgs(id, req.NamaJabatan).
			WillReturnError(errors.New("database error"))

		err := repo.Create(req, id)
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate key error", func(t *testing.T) {
		namaJabatan := "Manager IT"
		req := dto.CreateJabatanRequest{
			NamaJabatan: &namaJabatan,
		}
		id := "001"

		mock.ExpectExec("INSERT INTO jabatan").
			WithArgs(id, req.NamaJabatan).
			WillReturnError(errors.New("duplicate entry"))

		err := repo.Create(req, id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestJabatanRepository_GetAll(t *testing.T) {
	db, mock, repo := setupJabatanTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		createdAt := "2024-01-01 10:00:00"
		updatedAt := "2024-01-02 15:30:00"

		rows := sqlmock.NewRows([]string{
			"id", "nama_jabatan", "created_at", "updated_at",
		}).
			AddRow("001", "Manager IT", createdAt, updatedAt).
			AddRow("002", "Staff Admin", createdAt, updatedAt).
			AddRow("003", "Direktur", createdAt, updatedAt)

		mock.ExpectQuery("SELECT id, nama_jabatan, created_at, updated_at FROM jabatan").
			WillReturnRows(rows)

		result, err := repo.GetAll()
		assert.NoError(t, err)
		require.Len(t, result, 3)

		// Check first record
		assert.Equal(t, "001", result[0].ID)
		assert.Equal(t, "Manager IT", result[0].NamaJabatan)
		assert.Equal(t, createdAt, result[0].CreatedAt)
		assert.Equal(t, updatedAt, result[0].UpdatedAt)

		// Check second record
		assert.Equal(t, "002", result[1].ID)
		assert.Equal(t, "Staff Admin", result[1].NamaJabatan)

		// Check third record
		assert.Equal(t, "003", result[2].ID)
		assert.Equal(t, "Direktur", result[2].NamaJabatan)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, nama_jabatan, created_at, updated_at FROM jabatan").
			WillReturnError(errors.New("query error"))

		result, err := repo.GetAll()
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "query error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "nama_jabatan", "created_at", "updated_at",
		})

		mock.ExpectQuery("SELECT id, nama_jabatan, created_at, updated_at FROM jabatan").
			WillReturnRows(rows)

		result, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error - continues processing", func(t *testing.T) {
		// Note: The original code doesn't handle scan errors in the loop
		// It will just append zero-valued structs
		createdAt := "2024-01-01 10:00:00"
		updatedAt := "2024-01-02 15:30:00"

		rows := sqlmock.NewRows([]string{
			"id", "nama_jabatan", "created_at", "updated_at",
		}).
			AddRow("001", "Manager IT", createdAt, updatedAt).
			AddRow("invalid", nil, nil, nil) // This will cause scan error but code continues

		mock.ExpectQuery("SELECT id, nama_jabatan, created_at, updated_at FROM jabatan").
			WillReturnRows(rows)

		result, err := repo.GetAll()
		// No error returned because scan error is not checked
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestJabatanRepository_GetByID(t *testing.T) {
	db, mock, repo := setupJabatanTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		id := "001"
		createdAt := "2024-01-01 10:00:00"
		updatedAt := "2024-01-02 15:30:00"

		row := sqlmock.NewRows([]string{
			"id", "nama_jabatan", "created_at", "updated_at",
		}).
			AddRow(id, "Manager IT", createdAt, updatedAt)

		mock.ExpectQuery("SELECT id, nama_jabatan, created_at, updated_at FROM jabatan WHERE id=?").
			WithArgs(id).
			WillReturnRows(row)

		result, err := repo.GetByID(id)
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, id, result.ID)
		assert.Equal(t, "Manager IT", result.NamaJabatan)
		assert.Equal(t, createdAt, result.CreatedAt)
		assert.Equal(t, updatedAt, result.UpdatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		id := "jabatan-nonexistent"

		mock.ExpectQuery("SELECT id, nama_jabatan, created_at, updated_at FROM jabatan WHERE id=?").
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)

		result, err := repo.GetByID(id)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		id := "001"
		row := sqlmock.NewRows([]string{"id"}).
			AddRow(id)

		mock.ExpectQuery("SELECT id, nama_jabatan, created_at, updated_at FROM jabatan WHERE id=?").
			WithArgs(id).
			WillReturnRows(row)

		result, err := repo.GetByID(id)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		id := "001"

		mock.ExpectQuery("SELECT id, nama_jabatan, created_at, updated_at FROM jabatan WHERE id=?").
			WithArgs(id).
			WillReturnError(errors.New("database connection error"))

		result, err := repo.GetByID(id)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "database connection error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestJabatanRepository_Update(t *testing.T) {
	db, mock, repo := setupJabatanTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		id := "001"
		createdAt := "2024-01-01 10:00:00"
		updatedAt := "2024-01-02 15:30:00"

		req := dto.JabatanResponse{
			ID:          id,
			NamaJabatan: "Senior Manager IT",
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		mock.ExpectExec("UPDATE jabatan SET nama_jabatan=\\? WHERE id=\\?").
			WithArgs(req.NamaJabatan, id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(id, req)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		id := "001"
		createdAt := "2024-01-01 10:00:00"
		updatedAt := "2024-01-02 15:30:00"

		req := dto.JabatanResponse{
			ID:          id,
			NamaJabatan: "Senior Manager IT",
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		mock.ExpectExec("UPDATE jabatan SET nama_jabatan=\\? WHERE id=\\?").
			WithArgs(req.NamaJabatan, id).
			WillReturnError(errors.New("update error"))

		err := repo.Update(id, req)
		assert.Error(t, err)
		assert.Equal(t, "update error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no rows affected - not checked by function", func(t *testing.T) {
		id := "jabatan-nonexistent"
		createdAt := "2024-01-01 10:00:00"
		updatedAt := "2024-01-02 15:30:00"

		req := dto.JabatanResponse{
			ID:          id,
			NamaJabatan: "Senior Manager IT",
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		mock.ExpectExec("UPDATE jabatan SET nama_jabatan=\\? WHERE id=\\?").
			WithArgs(req.NamaJabatan, id).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Update(id, req)
		// Function doesn't check rows affected, so no error
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestJabatanRepository_Delete(t *testing.T) {
	db, mock, repo := setupJabatanTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		id := "001"

		mock.ExpectExec("DELETE FROM jabatan WHERE id=\\?").
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(id)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found - no rows affected", func(t *testing.T) {
		id := "jabatan-nonexistent"

		mock.ExpectExec("DELETE FROM jabatan WHERE id=\\?").
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete(id)
		assert.Error(t, err)
		assert.Equal(t, fmt.Sprintf("data dengan id %s tidak ditemukan", id), err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error on exec", func(t *testing.T) {
		id := "001"

		mock.ExpectExec("DELETE FROM jabatan WHERE id=\\?").
			WithArgs(id).
			WillReturnError(errors.New("delete error"))

		err := repo.Delete(id)
		assert.Error(t, err)
		assert.Equal(t, "delete error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error on rows affected check", func(t *testing.T) {
		id := "001"

		mock.ExpectExec("DELETE FROM jabatan WHERE id=\\?").
			WithArgs(id).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))

		err := repo.Delete(id)
		assert.Error(t, err)
		assert.Equal(t, "rows affected error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("multiple rows deleted", func(t *testing.T) {
		id := "001"

		mock.ExpectExec("DELETE FROM jabatan WHERE id=\\?").
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 5))

		err := repo.Delete(id)
		// Function doesn't check if more than 1 row was deleted
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
