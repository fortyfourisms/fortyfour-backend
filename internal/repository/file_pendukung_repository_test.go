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

func setupFPRepoTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *FilePendukungRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return db, mock, NewFilePendukungRepository(db)
}

var fpColumns = []string{"id", "id_materi", "nama_file", "file_path", "ukuran", "created_at"}

func TestFilePendukungRepository_Create(t *testing.T) {
	db, mock, repo := setupFPRepoTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO file_pendukung").
			WithArgs("fp-1", "m-1", "doc.pdf", "/uploads/doc.pdf", int64(1024)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(&models.FilePendukung{ID: "fp-1", IDMateri: "m-1", NamaFile: "doc.pdf", FilePath: "/uploads/doc.pdf", Ukuran: 1024})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO file_pendukung").
			WillReturnError(errors.New("db error"))

		err := repo.Create(&models.FilePendukung{ID: "fp-2"})
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestFilePendukungRepository_FindByMateri(t *testing.T) {
	db, mock, repo := setupFPRepoTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM file_pendukung WHERE id_materi = \\?").
			WithArgs("m-1").
			WillReturnRows(sqlmock.NewRows(fpColumns).
				AddRow("fp-1", "m-1", "doc1.pdf", "/p1", int64(100), now).
				AddRow("fp-2", "m-1", "doc2.pdf", "/p2", int64(200), now))

		result, err := repo.FindByMateri("m-1")
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM file_pendukung WHERE id_materi = \\?").
			WithArgs("empty").
			WillReturnRows(sqlmock.NewRows(fpColumns))

		result, err := repo.FindByMateri("empty")
		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM file_pendukung WHERE id_materi = \\?").
			WithArgs("m-1").
			WillReturnError(errors.New("db error"))

		result, err := repo.FindByMateri("m-1")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestFilePendukungRepository_FindByID(t *testing.T) {
	db, mock, repo := setupFPRepoTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM file_pendukung WHERE id = \\?").
			WithArgs("fp-1").
			WillReturnRows(sqlmock.NewRows(fpColumns).
				AddRow("fp-1", "m-1", "doc.pdf", "/path", int64(500), now))

		result, err := repo.FindByID("fp-1")
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "doc.pdf", result.NamaFile)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM file_pendukung WHERE id = \\?").
			WithArgs("invalid").
			WillReturnError(sql.ErrNoRows)

		result, err := repo.FindByID("invalid")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestFilePendukungRepository_Delete(t *testing.T) {
	db, mock, repo := setupFPRepoTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM file_pendukung WHERE id=\\?").
			WithArgs("fp-1").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete("fp-1")
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM file_pendukung WHERE id=\\?").
			WithArgs("invalid").
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete("invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tidak ditemukan")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
