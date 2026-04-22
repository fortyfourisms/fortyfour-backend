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

func setupSoalRepoTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *SoalRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return db, mock, NewSoalRepository(db)
}

func TestSoalRepository_Create(t *testing.T) {
	db, mock, repo := setupSoalRepoTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO soal").
			WithArgs("s-1", "q-1", "Apa itu Go?", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO pilihan_jawaban").
			WithArgs("p-1", "s-1", "Bahasa", true, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO pilihan_jawaban").
			WithArgs("p-2", "s-1", "Framework", false, 2).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(
			&models.Soal{ID: "s-1", IDKuis: "q-1", Pertanyaan: "Apa itu Go?", Urutan: 1},
			[]models.PilihanJawaban{
				{ID: "p-1", IDSoal: "s-1", Teks: "Bahasa", IsCorrect: true, Urutan: 1},
				{ID: "p-2", IDSoal: "s-1", Teks: "Framework", IsCorrect: false, Urutan: 2},
			},
		)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("soal insert error rolls back", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO soal").
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		err := repo.Create(&models.Soal{ID: "s-2"}, nil)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSoalRepository_FindByID(t *testing.T) {
	db, mock, repo := setupSoalRepoTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("success with pilihan", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM soal WHERE id = \\?").
			WithArgs("s-1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "id_kuis", "pertanyaan", "urutan", "created_at"}).
				AddRow("s-1", "q-1", "Q?", 1, now))
		mock.ExpectQuery("SELECT .+ FROM pilihan_jawaban").
			WithArgs("s-1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "id_soal", "teks", "is_correct", "urutan"}).
				AddRow("p-1", "s-1", "A", true, 1).
				AddRow("p-2", "s-1", "B", false, 2))

		result, err := repo.FindByID("s-1")
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Pilihan, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM soal WHERE id = \\?").
			WithArgs("invalid").
			WillReturnError(sql.ErrNoRows)

		result, err := repo.FindByID("invalid")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSoalRepository_FindByKuis(t *testing.T) {
	db, mock, repo := setupSoalRepoTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM soal WHERE id_kuis = \\?").
			WithArgs("q-1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "id_kuis", "pertanyaan", "urutan", "created_at"}).
				AddRow("s-1", "q-1", "Q1?", 1, now))
		mock.ExpectQuery("SELECT .+ FROM pilihan_jawaban").
			WithArgs("s-1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "id_soal", "teks", "is_correct", "urutan"}).
				AddRow("p-1", "s-1", "A", true, 1))

		result, err := repo.FindByKuis("q-1")
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Len(t, result[0].Pilihan, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM soal WHERE id_kuis = \\?").
			WithArgs("empty").
			WillReturnRows(sqlmock.NewRows([]string{"id", "id_kuis", "pertanyaan", "urutan", "created_at"}))

		result, err := repo.FindByKuis("empty")
		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSoalRepository_Delete(t *testing.T) {
	db, mock, repo := setupSoalRepoTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM soal WHERE id=\\?").
			WithArgs("s-1").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete("s-1")
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM soal WHERE id=\\?").
			WithArgs("invalid").
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete("invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tidak ditemukan")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSoalRepository_FindPilihanByID(t *testing.T) {
	db, mock, repo := setupSoalRepoTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM pilihan_jawaban WHERE id=\\?").
			WithArgs("p-1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "id_soal", "teks", "is_correct", "urutan"}).
				AddRow("p-1", "s-1", "Correct", true, 1))

		result, err := repo.FindPilihanByID("p-1")
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.IsCorrect)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM pilihan_jawaban WHERE id=\\?").
			WithArgs("invalid").
			WillReturnError(sql.ErrNoRows)

		result, err := repo.FindPilihanByID("invalid")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSoalRepository_FindCorrectPilihan(t *testing.T) {
	db, mock, repo := setupSoalRepoTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT .+ FROM pilihan_jawaban.+is_correct=1").
			WithArgs("s-1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "id_soal", "teks", "is_correct", "urutan"}).
				AddRow("p-1", "s-1", "Right", true, 1))

		result, err := repo.FindCorrectPilihan("s-1")
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.IsCorrect)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
