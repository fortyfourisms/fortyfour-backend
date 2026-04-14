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

func setupProgressRepoTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *ProgressRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return db, mock, NewProgressRepository(db)
}

var progressColumns = []string{"id", "id_user", "id_materi", "is_completed", "last_watched_seconds", "completed_at", "created_at", "updated_at"}

func TestProgressRepository_Upsert(t *testing.T) {
	db, mock, repo := setupProgressRepoTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		mock.ExpectExec("INSERT INTO user_materi_progress").
			WithArgs("p-1", "u-1", "m-1", true, 100, &now).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Upsert(&models.UserMateriProgress{
			ID: "p-1", IDUser: "u-1", IDMateri: "m-1",
			IsCompleted: true, LastWatchedSeconds: 100, CompletedAt: &now,
		})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO user_materi_progress").
			WillReturnError(errors.New("db error"))

		err := repo.Upsert(&models.UserMateriProgress{ID: "p-2"})
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProgressRepository_FindByUserAndMateri(t *testing.T) {
	db, mock, repo := setupProgressRepoTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("success completed", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM user_materi_progress.+WHERE id_user=\\? AND id_materi=\\?").
			WithArgs("u-1", "m-1").
			WillReturnRows(sqlmock.NewRows(progressColumns).
				AddRow("p-1", "u-1", "m-1", true, 100, now, now, now))

		result, err := repo.FindByUserAndMateri("u-1", "m-1")
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.IsCompleted)
		assert.NotNil(t, result.CompletedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success not completed", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM user_materi_progress.+WHERE id_user=\\? AND id_materi=\\?").
			WithArgs("u-1", "m-2").
			WillReturnRows(sqlmock.NewRows(progressColumns).
				AddRow("p-2", "u-1", "m-2", false, 30, nil, now, now))

		result, err := repo.FindByUserAndMateri("u-1", "m-2")
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.False(t, result.IsCompleted)
		assert.Nil(t, result.CompletedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM user_materi_progress.+WHERE id_user=\\? AND id_materi=\\?").
			WithArgs("u-1", "invalid").
			WillReturnError(sql.ErrNoRows)

		result, err := repo.FindByUserAndMateri("u-1", "invalid")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProgressRepository_HasCompletedAllMateri(t *testing.T) {
	db, mock, repo := setupProgressRepoTest(t)
	defer db.Close()

	t.Run("all completed", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM materi WHERE id_kelas = \\?").
			WithArgs("k-1").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM user_materi_progress").
			WithArgs("u-1", "k-1").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

		completed, err := repo.HasCompletedAllMateri("u-1", "k-1")
		assert.NoError(t, err)
		assert.True(t, completed)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not all completed", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM materi WHERE id_kelas = \\?").
			WithArgs("k-1").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM user_materi_progress").
			WithArgs("u-1", "k-1").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		completed, err := repo.HasCompletedAllMateri("u-1", "k-1")
		assert.NoError(t, err)
		assert.False(t, completed)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no materi returns false", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM materi WHERE id_kelas = \\?").
			WithArgs("k-empty").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		completed, err := repo.HasCompletedAllMateri("u-1", "k-empty")
		assert.NoError(t, err)
		assert.False(t, completed)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
