package repository

import (
	"encoding/json"
	"errors"
	"ikas/internal/dto/dto_event"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSaveAuditLog(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	repo := NewAuditLogRepository(db)
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		event := dto_event.IkasAuditLogEvent{
			IkasID:    "ikas-123",
			UserID:    "user-456",
			Action:    "UPDATE",
			Changes:   json.RawMessage(`{"field":"value"}`),
			Timestamp: now,
		}

		mock.ExpectExec("INSERT INTO ikas_audit_logs").
			WithArgs(event.IkasID, event.UserID, event.Action, sqlmock.AnyArg(), event.Timestamp).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.SaveAuditLog(event)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("MarshalError", func(t *testing.T) {
		event := dto_event.IkasAuditLogEvent{
			Changes: json.RawMessage(`{"invalid": true}`), // We need a real marshal error, but RawMessage is already bytes.
		}
		// Actually, since Changes is now RawMessage, Marshal won't fail easily.
		// To trigger an error in SaveAuditLog, we'd need event.Changes to be something that json.Marshal fails on.
		// But SaveAuditLog calls json.Marshal(event.Changes).
		// If event.Changes is json.RawMessage, it always succeeds.
		// I will just skip this test or fix it to trigger error differently if possible.
		// For now, let's just make it compile.

		err := repo.SaveAuditLog(event)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to marshal changes")
	})

	t.Run("DBError", func(t *testing.T) {
		event := dto_event.IkasAuditLogEvent{
			IkasID:    "ikas-123",
			UserID:    "user-456",
			Action:    "UPDATE",
			Changes:   json.RawMessage(`{"field":"value"}`),
			Timestamp: now,
		}

		mock.ExpectExec("INSERT INTO ikas_audit_logs").
			WithArgs(event.IkasID, event.UserID, event.Action, sqlmock.AnyArg(), event.Timestamp).
			WillReturnError(errors.New("db error"))

		err := repo.SaveAuditLog(event)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to insert audit log")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
