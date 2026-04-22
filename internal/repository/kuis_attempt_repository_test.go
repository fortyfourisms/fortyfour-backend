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

func setupAttemptRepoTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *KuisAttemptRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return db, mock, NewKuisAttemptRepository(db)
}

var attemptColumns = []string{"id", "id_user", "id_kuis", "skor", "total_soal", "total_benar", "is_passed", "started_at", "finished_at"}

func TestAttemptRepository_Create(t *testing.T) {
	db, mock, repo := setupAttemptRepoTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO kuis_attempt").
			WithArgs("a-1", "u-1", "q-1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(&models.KuisAttempt{ID: "a-1", IDUser: "u-1", IDKuis: "q-1"})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO kuis_attempt").
			WillReturnError(errors.New("db error"))

		err := repo.Create(&models.KuisAttempt{ID: "a-2"})
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAttemptRepository_FindByID(t *testing.T) {
	db, mock, repo := setupAttemptRepoTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("success unfinished", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM kuis_attempt WHERE id=\\?").
			WithArgs("a-1").
			WillReturnRows(sqlmock.NewRows(attemptColumns).
				AddRow("a-1", "u-1", "q-1", 0, 0, 0, false, now, nil))

		result, err := repo.FindByID("a-1")
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Nil(t, result.FinishedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success finished", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM kuis_attempt WHERE id=\\?").
			WithArgs("a-2").
			WillReturnRows(sqlmock.NewRows(attemptColumns).
				AddRow("a-2", "u-1", "q-1", 100, 5, 5, true, now, now))

		result, err := repo.FindByID("a-2")
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.FinishedAt)
		assert.True(t, result.IsPassed)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM kuis_attempt WHERE id=\\?").
			WithArgs("invalid").
			WillReturnError(sql.ErrNoRows)

		result, err := repo.FindByID("invalid")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAttemptRepository_Finish(t *testing.T) {
	db, mock, repo := setupAttemptRepoTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE kuis_attempt").
			WithArgs(80.0, 2, 1, true, "a-1").
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("INSERT INTO kuis_jawaban").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO kuis_jawaban").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Finish("a-1", 80.0, 1, true, []models.KuisJawaban{
			{ID: "j1", IDAttempt: "a-1", IDSoal: "s1", IDPilihan: "p1", IsCorrect: true},
			{ID: "j2", IDAttempt: "a-1", IDSoal: "s2", IDPilihan: "p3", IsCorrect: false},
		})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("attempt not found or already finished", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE kuis_attempt").
			WillReturnResult(sqlmock.NewResult(0, 0)) // no rows affected
		mock.ExpectRollback()

		err := repo.Finish("a-gone", 0, 0, false, []models.KuisJawaban{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tidak ditemukan atau sudah selesai")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAttemptRepository_HasPassedAllKuisInKelas(t *testing.T) {
	db, mock, repo := setupAttemptRepoTest(t)
	defer db.Close()

	t.Run("all passed", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM kuis WHERE id_kelas=\\? AND is_final=0").
			WithArgs("k-1").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
		mock.ExpectQuery("SELECT COUNT\\(DISTINCT ka.id_kuis\\)").
			WithArgs("u-1", "k-1").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		passed, err := repo.HasPassedAllKuisInKelas("u-1", "k-1")
		assert.NoError(t, err)
		assert.True(t, passed)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not all passed", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM kuis WHERE id_kelas=\\? AND is_final=0").
			WithArgs("k-1").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
		mock.ExpectQuery("SELECT COUNT\\(DISTINCT ka.id_kuis\\)").
			WithArgs("u-1", "k-1").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		passed, err := repo.HasPassedAllKuisInKelas("u-1", "k-1")
		assert.NoError(t, err)
		assert.False(t, passed)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no non-final kuis returns true", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM kuis WHERE id_kelas=\\? AND is_final=0").
			WithArgs("k-empty").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		passed, err := repo.HasPassedAllKuisInKelas("u-1", "k-empty")
		assert.NoError(t, err)
		assert.True(t, passed)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAttemptRepository_FindJawabanByAttempt(t *testing.T) {
	db, mock, repo := setupAttemptRepoTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM kuis_jawaban WHERE id_attempt=\\?").
			WithArgs("a-1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "id_attempt", "id_soal", "id_pilihan", "is_correct"}).
				AddRow("j1", "a-1", "s1", "p1", true).
				AddRow("j2", "a-1", "s2", "p3", false))

		result, err := repo.FindJawabanByAttempt("a-1")
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM kuis_jawaban WHERE id_attempt=\\?").
			WithArgs("empty").
			WillReturnRows(sqlmock.NewRows([]string{"id", "id_attempt", "id_soal", "id_pilihan", "is_correct"}))

		result, err := repo.FindJawabanByAttempt("empty")
		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
