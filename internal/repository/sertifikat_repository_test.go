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

func setupSertifikatRepoTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *SertifikatRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return db, mock, NewSertifikatRepository(db)
}

var sertifikatColumns = []string{"id", "nomor_sertifikat", "id_kelas", "id_user", "nama_peserta", "nama_kelas", "tanggal_terbit", "pdf_path", "created_at"}

func TestSertifikatRepository_Create(t *testing.T) {
	db, mock, repo := setupSertifikatRepoTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO sertifikat").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(&models.Sertifikat{
			ID: "c-1", NomorSertifikat: "CERT/001", IDKelas: "k-1", IDUser: "u-1",
			NamaPeserta: "John", NamaKelas: "Go", TanggalTerbit: time.Now(),
		})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO sertifikat").
			WillReturnError(errors.New("db error"))

		err := repo.Create(&models.Sertifikat{ID: "c-2"})
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSertifikatRepository_FindByUserAndKelas(t *testing.T) {
	db, mock, repo := setupSertifikatRepoTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("found", func(t *testing.T) {
		pdfPath := "/path/to/cert.pdf"
		mock.ExpectQuery("SELECT .+ FROM sertifikat WHERE id_user=\\? AND id_kelas=\\?").
			WithArgs("u-1", "k-1").
			WillReturnRows(sqlmock.NewRows(sertifikatColumns).
				AddRow("c-1", "CERT/001", "k-1", "u-1", "John", "Go", now, pdfPath, now))

		result, err := repo.FindByUserAndKelas("u-1", "k-1")
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "CERT/001", result.NomorSertifikat)
		assert.Equal(t, &pdfPath, result.PDFPath)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM sertifikat WHERE id_user=\\? AND id_kelas=\\?").
			WithArgs("u-1", "invalid").
			WillReturnError(sql.ErrNoRows)

		result, err := repo.FindByUserAndKelas("u-1", "invalid")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSertifikatRepository_FindByID(t *testing.T) {
	db, mock, repo := setupSertifikatRepoTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM sertifikat WHERE id=\\?").
			WithArgs("c-1").
			WillReturnRows(sqlmock.NewRows(sertifikatColumns).
				AddRow("c-1", "CERT/001", "k-1", "u-1", "John", "Go", now, nil, now))

		result, err := repo.FindByID("c-1")
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Nil(t, result.PDFPath) // nullable
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM sertifikat WHERE id=\\?").
			WithArgs("invalid").
			WillReturnError(sql.ErrNoRows)

		result, err := repo.FindByID("invalid")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSertifikatRepository_FindByUser(t *testing.T) {
	db, mock, repo := setupSertifikatRepoTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM sertifikat WHERE id_user=\\?").
			WithArgs("u-1").
			WillReturnRows(sqlmock.NewRows(sertifikatColumns).
				AddRow("c-1", "CERT/001", "k-1", "u-1", "John", "Go", now, nil, now).
				AddRow("c-2", "CERT/002", "k-2", "u-1", "John", "Rust", now, nil, now))

		result, err := repo.FindByUser("u-1")
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM sertifikat WHERE id_user=\\?").
			WithArgs("u-new").
			WillReturnRows(sqlmock.NewRows(sertifikatColumns))

		result, err := repo.FindByUser("u-new")
		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM sertifikat WHERE id_user=\\?").
			WithArgs("u-1").
			WillReturnError(errors.New("db error"))

		result, err := repo.FindByUser("u-1")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
