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

func setupDiskusiRepoTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *DiskusiRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return db, mock, NewDiskusiRepository(db)
}

var diskusiColumns = []string{"id", "id_materi", "id_user", "id_parent", "konten", "created_at", "updated_at"}

func TestDiskusiRepository_Create(t *testing.T) {
	db, mock, repo := setupDiskusiRepoTest(t)
	defer db.Close()

	t.Run("success top-level", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO diskusi").
			WithArgs("d-1", "m-1", "u-1", nil, "Hello").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(&models.Diskusi{ID: "d-1", IDMateri: "m-1", IDUser: "u-1", Konten: "Hello"})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success reply", func(t *testing.T) {
		parentID := "d-1"
		mock.ExpectExec("INSERT INTO diskusi").
			WithArgs("d-2", "m-1", "u-2", &parentID, "Reply!").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(&models.Diskusi{ID: "d-2", IDMateri: "m-1", IDUser: "u-2", IDParent: &parentID, Konten: "Reply!"})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO diskusi").
			WillReturnError(errors.New("db error"))

		err := repo.Create(&models.Diskusi{ID: "d-3"})
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDiskusiRepository_FindByMateri(t *testing.T) {
	db, mock, repo := setupDiskusiRepoTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM diskusi WHERE id_materi = \\? AND id_parent IS NULL").
			WithArgs("m-1").
			WillReturnRows(sqlmock.NewRows(diskusiColumns).
				AddRow("d-1", "m-1", "u-1", nil, "Hello", now, now))

		result, err := repo.FindByMateri("m-1")
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Nil(t, result[0].IDParent) // top-level
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM diskusi WHERE id_materi = \\? AND id_parent IS NULL").
			WithArgs("empty").
			WillReturnRows(sqlmock.NewRows(diskusiColumns))

		result, err := repo.FindByMateri("empty")
		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDiskusiRepository_FindByID(t *testing.T) {
	db, mock, repo := setupDiskusiRepoTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM diskusi WHERE id = \\?").
			WithArgs("d-1").
			WillReturnRows(sqlmock.NewRows(diskusiColumns).
				AddRow("d-1", "m-1", "u-1", nil, "Hello", now, now))

		result, err := repo.FindByID("d-1")
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "Hello", result.Konten)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM diskusi WHERE id = \\?").
			WithArgs("invalid").
			WillReturnError(sql.ErrNoRows)

		result, err := repo.FindByID("invalid")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDiskusiRepository_Update(t *testing.T) {
	db, mock, repo := setupDiskusiRepoTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE diskusi SET konten=\\?, updated_at=NOW\\(\\) WHERE id=\\?").
			WithArgs("Updated", "d-1").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(&models.Diskusi{ID: "d-1", Konten: "Updated"})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDiskusiRepository_Delete(t *testing.T) {
	db, mock, repo := setupDiskusiRepoTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM diskusi WHERE id=\\?").
			WithArgs("d-1").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete("d-1")
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM diskusi WHERE id=\\?").
			WithArgs("invalid").
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete("invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tidak ditemukan")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDiskusiRepository_FindReplies(t *testing.T) {
	db, mock, repo := setupDiskusiRepoTest(t)
	defer db.Close()

	now := time.Now()
	parentID := "d-1"

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM diskusi WHERE id_parent = \\?").
			WithArgs("d-1").
			WillReturnRows(sqlmock.NewRows(diskusiColumns).
				AddRow("r-1", "m-1", "u-2", parentID, "Reply 1", now, now).
				AddRow("r-2", "m-1", "u-3", parentID, "Reply 2", now, now))

		result, err := repo.FindReplies("d-1")
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM diskusi WHERE id_parent = \\?").
			WithArgs("d-noreplies").
			WillReturnRows(sqlmock.NewRows(diskusiColumns))

		result, err := repo.FindReplies("d-noreplies")
		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
