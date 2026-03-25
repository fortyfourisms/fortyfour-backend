package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"ikas/internal/dto/dto_event"
)

type AuditLogRepositoryInterface interface {
	SaveAuditLog(event dto_event.IkasAuditLogEvent) error
}

type AuditLogRepository struct {
	db *sql.DB
}

func NewAuditLogRepository(db *sql.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) SaveAuditLog(event dto_event.IkasAuditLogEvent) error {
	changesJSON, err := json.Marshal(event.Changes)
	if err != nil {
		return fmt.Errorf("failed to marshal changes: %w", err)
	}

	query := `
		INSERT INTO ikas_audit_logs (ikas_id, user_id, action, changes, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err = r.db.Exec(query, event.IkasID, event.UserID, event.Action, changesJSON, event.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to insert audit log: %w", err)
	}

	return nil
}
