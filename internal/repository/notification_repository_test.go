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

var notificationColumns = []string{
	"id", "user_id", "username", "display_name", "foto_profile", "type", "message", "is_read", "created_at",
}

func setupNotificationRepoTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *NotificationRepository) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	return db, mock, NewNotificationRepository(db)
}

func TestNotificationRepository_Create(t *testing.T) {
	db, mock, repo := setupNotificationRepoTest(t)
	defer db.Close()

	notif := &models.Notification{
		UserID:  "user-1",
		Type:    models.NotifResourceCreated,
		Message: "resource created",
		Read:    false,
	}

	t.Run("success sets inserted id", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO notifications").
			WithArgs(notif.UserID, notif.Type, notif.Message, notif.Read).
			WillReturnResult(sqlmock.NewResult(42, 1))

		err := repo.Create(notif)
		assert.NoError(t, err)
		assert.Equal(t, int64(42), notif.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO notifications").
			WithArgs(notif.UserID, notif.Type, notif.Message, notif.Read).
			WillReturnError(errors.New("db error"))

		err := repo.Create(notif)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNotificationRepository_FindAll(t *testing.T) {
	db, mock, repo := setupNotificationRepoTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT n.id, n.user_id, u.username, COALESCE\\(u.display_name, ''\\) as display_name").
			WillReturnRows(sqlmock.NewRows(notificationColumns).
				AddRow(1, "user-1", "alice", "Alice", "alice.jpg", "resource_created", "created", false, now).
				AddRow(2, "user-2", "bob", "", nil, "resource_deleted", "deleted", true, now))

		result, err := repo.FindAll()
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "alice", result[0].Username)
		require.NotNil(t, result[0].FotoProfile)
		assert.Equal(t, "alice.jpg", *result[0].FotoProfile)
		assert.Nil(t, result[1].FotoProfile)
		assert.Equal(t, "", result[1].DisplayName)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT n.id, n.user_id, u.username, COALESCE\\(u.display_name, ''\\) as display_name").
			WillReturnError(errors.New("db error"))

		result, err := repo.FindAll()
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNotificationRepository_FindAllByUserID(t *testing.T) {
	db, mock, repo := setupNotificationRepoTest(t)
	defer db.Close()

	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("WHERE n.user_id = \\?").
			WithArgs("user-1").
			WillReturnRows(sqlmock.NewRows(notificationColumns).
				AddRow(1, "user-1", "alice", "Alice", "alice.jpg", "resource_created", "created", false, now))

		result, err := repo.FindAllByUserID("user-1")
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "user-1", result[0].UserID)
		require.NotNil(t, result[0].FotoProfile)
		assert.Equal(t, "alice.jpg", *result[0].FotoProfile)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("WHERE n.user_id = \\?").
			WithArgs("user-2").
			WillReturnError(errors.New("db error"))

		result, err := repo.FindAllByUserID("user-2")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNotificationRepository_MarkRead(t *testing.T) {
	db, mock, repo := setupNotificationRepoTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE notifications SET is_read = TRUE WHERE id = \\? AND user_id = \\?").
			WithArgs(int64(10), "user-1").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.MarkRead("user-1", 10)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("UPDATE notifications SET is_read = TRUE WHERE id = \\? AND user_id = \\?").
			WithArgs(int64(11), "user-2").
			WillReturnError(errors.New("db error"))

		err := repo.MarkRead("user-2", 11)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNotificationRepository_MarkAllRead(t *testing.T) {
	db, mock, repo := setupNotificationRepoTest(t)
	defer db.Close()

	mock.ExpectExec("UPDATE notifications SET is_read = TRUE WHERE user_id = \\?").
		WithArgs("user-1").
		WillReturnResult(sqlmock.NewResult(0, 3))

	err := repo.MarkAllRead("user-1")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationRepository_Delete(t *testing.T) {
	db, mock, repo := setupNotificationRepoTest(t)
	defer db.Close()

	mock.ExpectExec("DELETE FROM notifications WHERE id = \\? AND user_id = \\?").
		WithArgs(int64(5), "user-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete("user-1", 5)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationRepository_DeleteAllByUserID(t *testing.T) {
	db, mock, repo := setupNotificationRepoTest(t)
	defer db.Close()

	mock.ExpectExec("DELETE FROM notifications WHERE user_id = \\?").
		WithArgs("user-1").
		WillReturnResult(sqlmock.NewResult(0, 4))

	err := repo.DeleteAllByUserID("user-1")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
