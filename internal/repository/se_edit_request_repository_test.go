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

var seEditRequestColumns = []string{
	"id", "id_se", "id_user", "status", "catatan_user", "catatan", "data_perubahan", "created_at", "updated_at",
}

func setupSEEditRequestRepoTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *SEEditRequestRepository) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	return db, mock, NewSEEditRequestRepository(db)
}

func TestSEEditRequestRepository_Create(t *testing.T) {
	db, mock, repo := setupSEEditRequestRepoTest(t)
	defer db.Close()

	req := &models.SEEditRequest{
		ID:            "req-1",
		IDSE:          "se-1",
		IDUser:        "user-1",
		Status:        models.SEEditRequestPending,
		DataPerubahan: `{"nama_se":"Updated"}`,
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO se_edit_request").
			WithArgs(req.ID, req.IDSE, req.IDUser, req.Status, req.CatatanUser, req.DataPerubahan).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(req)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO se_edit_request").
			WithArgs(req.ID, req.IDSE, req.IDUser, req.Status, req.CatatanUser, req.DataPerubahan).
			WillReturnError(errors.New("db error"))

		err := repo.Create(req)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSEEditRequestRepository_FindByID(t *testing.T) {
	db, mock, repo := setupSEEditRequestRepoTest(t)
	defer db.Close()

	now := time.Now()
	catatanUser := "perlu update"
	catatanAdmin := "sudah dicek"

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, id_se, id_user, status, catatan_user, catatan, data_perubahan, created_at, updated_at").
			WithArgs("req-1").
			WillReturnRows(sqlmock.NewRows(seEditRequestColumns).
				AddRow("req-1", "se-1", "user-1", "pending", catatanUser, catatanAdmin, `{"nama_se":"Updated"}`, now, now))

		result, err := repo.FindByID("req-1")
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "req-1", result.ID)
		assert.Equal(t, "se-1", result.IDSE)
		assert.Equal(t, models.SEEditRequestPending, result.Status)
		require.NotNil(t, result.CatatanUser)
		assert.Equal(t, catatanUser, *result.CatatanUser)
		require.NotNil(t, result.Catatan)
		assert.Equal(t, catatanAdmin, *result.Catatan)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, id_se, id_user, status, catatan_user, catatan, data_perubahan, created_at, updated_at").
			WithArgs("missing").
			WillReturnError(sql.ErrNoRows)

		result, err := repo.FindByID("missing")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSEEditRequestRepository_FindPendingBySE(t *testing.T) {
	db, mock, repo := setupSEEditRequestRepoTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("FROM se_edit_request WHERE id_se = \\? AND status = 'pending' ORDER BY created_at DESC").
			WithArgs("se-1").
			WillReturnRows(sqlmock.NewRows(seEditRequestColumns).
				AddRow("req-1", "se-1", "user-1", "pending", nil, nil, `{"nama_se":"A"}`, now, now).
				AddRow("req-2", "se-1", "user-2", "pending", nil, nil, `{"nama_se":"B"}`, now, now))

		result, err := repo.FindPendingBySE("se-1")
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "req-1", result[0].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("FROM se_edit_request WHERE id_se = \\? AND status = 'pending' ORDER BY created_at DESC").
			WithArgs("se-2").
			WillReturnError(errors.New("db error"))

		result, err := repo.FindPendingBySE("se-2")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSEEditRequestRepository_FindAllPending(t *testing.T) {
	db, mock, repo := setupSEEditRequestRepoTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("FROM se_edit_request WHERE status = 'pending' ORDER BY created_at DESC").
			WillReturnRows(sqlmock.NewRows(seEditRequestColumns).
				AddRow("req-1", "se-1", "user-1", "pending", nil, nil, `{"nama_se":"A"}`, now, now))

		result, err := repo.FindAllPending()
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "req-1", result[0].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("FROM se_edit_request WHERE status = 'pending' ORDER BY created_at DESC").
			WillReturnError(errors.New("db error"))

		result, err := repo.FindAllPending()
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSEEditRequestRepository_FindByUser(t *testing.T) {
	db, mock, repo := setupSEEditRequestRepoTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("FROM se_edit_request WHERE id_user = \\? ORDER BY created_at DESC").
			WithArgs("user-1").
			WillReturnRows(sqlmock.NewRows(seEditRequestColumns).
				AddRow("req-1", "se-1", "user-1", "pending", nil, nil, `{"nama_se":"A"}`, now, now))

		result, err := repo.FindByUser("user-1")
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "user-1", result[0].IDUser)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("FROM se_edit_request WHERE id_user = \\? ORDER BY created_at DESC").
			WithArgs("user-2").
			WillReturnError(errors.New("db error"))

		result, err := repo.FindByUser("user-2")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSEEditRequestRepository_UpdateStatus(t *testing.T) {
	db, mock, repo := setupSEEditRequestRepoTest(t)
	defer db.Close()

	catatan := "disetujui"

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE se_edit_request SET status = \\?, catatan = \\?, updated_at = NOW\\(\\) WHERE id = \\?").
			WithArgs(models.SEEditRequestApproved, &catatan, "req-1").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateStatus("req-1", models.SEEditRequestApproved, &catatan)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("UPDATE se_edit_request SET status = \\?, catatan = \\?, updated_at = NOW\\(\\) WHERE id = \\?").
			WithArgs(models.SEEditRequestRejected, nil, "req-2").
			WillReturnError(errors.New("db error"))

		err := repo.UpdateStatus("req-2", models.SEEditRequestRejected, nil)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
