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

func setupMateriTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *MateriRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return db, mock, NewMateriRepository(db)
}

var materiColumns = []string{"id", "id_kelas", "judul", "tipe", "urutan", "youtube_id", "durasi_detik", "konten_html", "deskripsi_singkat", "kategori", "created_at", "updated_at"}

func TestMateriRepository_Create(t *testing.T) {
	db, mock, repo := setupMateriTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		ytID := "abc123"
		durasi := 600
		mock.ExpectExec("INSERT INTO materi").
			WithArgs("m-1", "k-1", "Intro", models.MateriTipeVideo, 1, &ytID, &durasi, nil, nil, nil).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(&models.Materi{ID: "m-1", IDKelas: "k-1", Judul: "Intro", Tipe: models.MateriTipeVideo, Urutan: 1, YoutubeID: &ytID, DurasiDetik: &durasi})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO materi").
			WillReturnError(errors.New("db error"))

		err := repo.Create(&models.Materi{ID: "m-2"})
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMateriRepository_FindByID(t *testing.T) {
	db, mock, repo := setupMateriTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("success video", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM materi WHERE id = \\?").
			WithArgs("m-1").
			WillReturnRows(sqlmock.NewRows(materiColumns).
				AddRow("m-1", "k-1", "Intro", "video", 1, "yt123", 600, nil, nil, nil, now, now))

		result, err := repo.FindByID("m-1")
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "m-1", result.ID)
		assert.Equal(t, models.MateriTipeVideo, result.Tipe)
		assert.Equal(t, "yt123", *result.YoutubeID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success teks", func(t *testing.T) {
		html := "<p>Hello</p>"
		mock.ExpectQuery("SELECT .+ FROM materi WHERE id = \\?").
			WithArgs("m-2").
			WillReturnRows(sqlmock.NewRows(materiColumns).
				AddRow("m-2", "k-1", "Artikel", "teks", 2, nil, nil, html, nil, nil, now, now))

		result, err := repo.FindByID("m-2")
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, &html, result.KontenHTML)
		assert.Nil(t, result.YoutubeID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM materi WHERE id = \\?").
			WithArgs("invalid").
			WillReturnError(sql.ErrNoRows)

		result, err := repo.FindByID("invalid")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMateriRepository_FindByKelas(t *testing.T) {
	db, mock, repo := setupMateriTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM materi WHERE id_kelas = \\?").
			WithArgs("k-1").
			WillReturnRows(sqlmock.NewRows(materiColumns).
				AddRow("m-1", "k-1", "A", "video", 1, "yt1", 100, nil, nil, nil, now, now).
				AddRow("m-2", "k-1", "B", "teks", 2, nil, nil, "<p>Hi</p>", nil, nil, now, now))

		result, err := repo.FindByKelas("k-1")
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM materi WHERE id_kelas = \\?").
			WithArgs("empty").
			WillReturnRows(sqlmock.NewRows(materiColumns))

		result, err := repo.FindByKelas("empty")
		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM materi WHERE id_kelas = \\?").
			WithArgs("k-1").
			WillReturnError(errors.New("db error"))

		result, err := repo.FindByKelas("k-1")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMateriRepository_Update(t *testing.T) {
	db, mock, repo := setupMateriTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE materi SET").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(&models.Materi{ID: "m-1", Judul: "Updated", Urutan: 1})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("UPDATE materi SET").
			WillReturnError(errors.New("db error"))

		err := repo.Update(&models.Materi{ID: "m-1"})
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMateriRepository_Delete(t *testing.T) {
	db, mock, repo := setupMateriTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM materi WHERE id=\\?").
			WithArgs("m-1").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete("m-1")
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM materi WHERE id=\\?").
			WithArgs("invalid").
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete("invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tidak ditemukan")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMateriRepository_ReorderUrutan(t *testing.T) {
	db, mock, repo := setupMateriTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT id FROM materi WHERE id_kelas = \\?").
			WithArgs("k-1").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow("m-3").AddRow("m-1").AddRow("m-2"))

		mock.ExpectExec("UPDATE materi SET urutan=\\? WHERE id=\\?").WithArgs(1, "m-3").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("UPDATE materi SET urutan=\\? WHERE id=\\?").WithArgs(2, "m-1").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("UPDATE materi SET urutan=\\? WHERE id=\\?").WithArgs(3, "m-2").WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.ReorderUrutan("k-1")
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id FROM materi WHERE id_kelas = \\?").
			WithArgs("k-1").
			WillReturnError(errors.New("db error"))

		err := repo.ReorderUrutan("k-1")
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
