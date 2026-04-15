package repository

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"fortyfour-backend/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupKelasTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *KelasRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	repo := NewKelasRepository(db)
	return db, mock, repo
}

var kelasColumns = []string{"id", "judul", "deskripsi", "thumbnail", "status", "created_by", "created_at", "updated_at"}

func TestKelasRepository_Create(t *testing.T) {
	db, mock, repo := setupKelasTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO kelas").
			WithArgs("k-1", "Go Class", nil, nil, models.KelasStatusDraft, "admin-1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(&models.Kelas{ID: "k-1", Judul: "Go Class", Status: models.KelasStatusDraft, CreatedBy: "admin-1"})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO kelas").
			WithArgs("k-2", "Fail", nil, nil, models.KelasStatusDraft, "admin-1").
			WillReturnError(errors.New("db error"))

		err := repo.Create(&models.Kelas{ID: "k-2", Judul: "Fail", Status: models.KelasStatusDraft, CreatedBy: "admin-1"})
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestKelasRepository_FindByID(t *testing.T) {
	db, mock, repo := setupKelasTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM kelas WHERE id = \\?").
			WithArgs("k-1").
			WillReturnRows(sqlmock.NewRows(kelasColumns).
				AddRow("k-1", "Go", nil, nil, "published", "admin", now, now))

		result, err := repo.FindByID("k-1")
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "k-1", result.ID)
		assert.Equal(t, "Go", result.Judul)
		assert.Nil(t, result.Deskripsi)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("with nullable fields", func(t *testing.T) {
		desc := "A Go course"
		thumb := "/img/go.png"
		mock.ExpectQuery("SELECT .+ FROM kelas WHERE id = \\?").
			WithArgs("k-2").
			WillReturnRows(sqlmock.NewRows(kelasColumns).
				AddRow("k-2", "Go Advanced", desc, thumb, "draft", "admin", now, now))

		result, err := repo.FindByID("k-2")
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, &desc, result.Deskripsi)
		assert.Equal(t, &thumb, result.Thumbnail)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM kelas WHERE id = \\?").
			WithArgs("invalid").
			WillReturnError(sql.ErrNoRows)

		result, err := repo.FindByID("invalid")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestKelasRepository_FindAll(t *testing.T) {
	db, mock, repo := setupKelasTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("all", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM kelas ORDER BY").
			WillReturnRows(sqlmock.NewRows(kelasColumns).
				AddRow("k-1", "A", nil, nil, "published", "admin", now, now).
				AddRow("k-2", "B", nil, nil, "draft", "admin", now, now))

		result, err := repo.FindAll(false)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("only published", func(t *testing.T) {
		mock.ExpectQuery("WHERE status = 'published'").
			WillReturnRows(sqlmock.NewRows(kelasColumns).
				AddRow("k-1", "A", nil, nil, "published", "admin", now, now))

		result, err := repo.FindAll(true)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM kelas").
			WillReturnRows(sqlmock.NewRows(kelasColumns))

		result, err := repo.FindAll(false)
		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM kelas").
			WillReturnError(errors.New("db error"))

		result, err := repo.FindAll(false)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestKelasRepository_Update(t *testing.T) {
	db, mock, repo := setupKelasTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE kelas SET").
			WithArgs("Go Updated", nil, nil, "published", "k-1").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(&models.Kelas{ID: "k-1", Judul: "Go Updated", Status: "published"})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("UPDATE kelas SET").
			WithArgs("Fail", nil, nil, "draft", "k-1").
			WillReturnError(errors.New("db error"))

		err := repo.Update(&models.Kelas{ID: "k-1", Judul: "Fail", Status: "draft"})
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestKelasRepository_Delete(t *testing.T) {
	db, mock, repo := setupKelasTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM kelas WHERE id=\\?").
			WithArgs("k-1").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete("k-1")
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM kelas WHERE id=\\?").
			WithArgs("invalid").
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete("invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tidak ditemukan")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM kelas WHERE id=\\?").
			WithArgs("k-1").
			WillReturnError(errors.New("db error"))

		err := repo.Delete("k-1")
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
