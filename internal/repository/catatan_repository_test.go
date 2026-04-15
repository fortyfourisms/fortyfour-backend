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

func setupCatatanRepoTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *CatatanRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return db, mock, NewCatatanRepository(db)
}

func TestCatatanRepository_Upsert(t *testing.T) {
	db, mock, repo := setupCatatanRepoTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO catatan_pribadi").
			WithArgs("c-1", "m-1", "u-1", "My notes").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Upsert(&models.CatatanPribadi{ID: "c-1", IDMateri: "m-1", IDUser: "u-1", Konten: "My notes"})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO catatan_pribadi").
			WillReturnError(errors.New("db error"))

		err := repo.Upsert(&models.CatatanPribadi{ID: "c-2"})
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCatatanRepository_FindByUserAndMateri(t *testing.T) {
	db, mock, repo := setupCatatanRepoTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM catatan_pribadi WHERE id_user=\\? AND id_materi=\\?").
			WithArgs("u-1", "m-1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "id_materi", "id_user", "konten", "created_at", "updated_at"}).
				AddRow("c-1", "m-1", "u-1", "Notes", now, now))

		result, err := repo.FindByUserAndMateri("u-1", "m-1")
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "Notes", result.Konten)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM catatan_pribadi WHERE id_user=\\? AND id_materi=\\?").
			WithArgs("u-1", "invalid").
			WillReturnError(sql.ErrNoRows)

		result, err := repo.FindByUserAndMateri("u-1", "invalid")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCatatanRepository_Delete(t *testing.T) {
	db, mock, repo := setupCatatanRepoTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM catatan_pribadi WHERE id=\\?").
			WithArgs("c-1").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete("c-1")
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM catatan_pribadi WHERE id=\\?").
			WithArgs("c-1").
			WillReturnError(errors.New("db error"))

		err := repo.Delete("c-1")
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
