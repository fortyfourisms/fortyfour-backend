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

func setupKuisRepoTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *KuisRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return db, mock, NewKuisRepository(db)
}

var kuisColumns = []string{"id", "id_kelas", "id_materi", "judul", "deskripsi", "durasi_menit", "passing_grade", "is_final", "urutan", "created_at", "updated_at"}

func TestKuisRepository_Create(t *testing.T) {
	db, mock, repo := setupKuisRepoTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO kuis").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(&models.Kuis{ID: "q-1", IDKelas: "k-1", Judul: "Kuis 1", PassingGrade: 70, Urutan: 1})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO kuis").
			WillReturnError(errors.New("db error"))

		err := repo.Create(&models.Kuis{ID: "q-2"})
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestKuisRepository_FindByID(t *testing.T) {
	db, mock, repo := setupKuisRepoTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM kuis WHERE id = \\?").
			WithArgs("q-1").
			WillReturnRows(sqlmock.NewRows(kuisColumns).
				AddRow("q-1", "k-1", nil, "Kuis 1", nil, nil, 70.0, false, 1, now, now))

		result, err := repo.FindByID("q-1")
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "q-1", result.ID)
		assert.Equal(t, float64(70), result.PassingGrade)
		assert.Nil(t, result.IDMateri)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM kuis WHERE id = \\?").
			WithArgs("invalid").
			WillReturnError(sql.ErrNoRows)

		result, err := repo.FindByID("invalid")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestKuisRepository_FindByKelas(t *testing.T) {
	db, mock, repo := setupKuisRepoTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM kuis WHERE id_kelas = \\?").
			WithArgs("k-1").
			WillReturnRows(sqlmock.NewRows(kuisColumns).
				AddRow("q-1", "k-1", "m-1", "Quiz 1", nil, nil, 70.0, false, 1, now, now).
				AddRow("q-2", "k-1", nil, "Final", nil, nil, 80.0, true, 2, now, now))

		result, err := repo.FindByKelas("k-1")
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM kuis WHERE id_kelas = \\?").
			WithArgs("empty").
			WillReturnRows(sqlmock.NewRows(kuisColumns))

		result, err := repo.FindByKelas("empty")
		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestKuisRepository_FindFinalByKelas(t *testing.T) {
	db, mock, repo := setupKuisRepoTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM kuis WHERE id_kelas = \\? AND is_final = 1").
			WithArgs("k-1").
			WillReturnRows(sqlmock.NewRows(kuisColumns).
				AddRow("q-final", "k-1", nil, "Kuis Akhir", nil, nil, 80.0, true, 1, now, now))

		result, err := repo.FindFinalByKelas("k-1")
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.IsFinal)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM kuis WHERE id_kelas = \\? AND is_final = 1").
			WithArgs("k-nofinal").
			WillReturnError(sql.ErrNoRows)

		result, err := repo.FindFinalByKelas("k-nofinal")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestKuisRepository_Update(t *testing.T) {
	db, mock, repo := setupKuisRepoTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE kuis SET").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(&models.Kuis{ID: "q-1", Judul: "Updated", PassingGrade: 75, Urutan: 1})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestKuisRepository_Delete(t *testing.T) {
	db, mock, repo := setupKuisRepoTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM kuis WHERE id=\\?").
			WithArgs("q-1").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete("q-1")
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM kuis WHERE id=\\?").
			WithArgs("invalid").
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete("invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tidak ditemukan")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
