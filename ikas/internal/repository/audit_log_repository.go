package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"ikas/internal/dto/dto_event"
	"ikas/internal/models"
)

type AuditLogRepositoryInterface interface {
	SaveAuditLog(event dto_event.IkasAuditLogEvent) error
	GetAuditLogsByIkasID(ikasID string) ([]models.AuditLog, error)
	GetAllAuditLogs() ([]models.AuditLog, error)
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

func (r *AuditLogRepository) GetAuditLogsByIkasID(ikasID string) ([]models.AuditLog, error) {
	query := `
		SELECT l.id, l.ikas_id, l.user_id, COALESCE(u.display_name, u.username, 'System') as user_name, l.action, l.changes, l.created_at
		FROM ikas_audit_logs l
		LEFT JOIN users u ON l.user_id = u.id
		WHERE l.ikas_id = ?
		ORDER BY l.created_at DESC
	`

	rows, err := r.db.Query(query, ikasID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch audit logs: %w", err)
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		err := rows.Scan(&log.ID, &log.IkasID, &log.UserID, &log.UserName, &log.Action, &log.Changes, &log.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (r *AuditLogRepository) GetAllAuditLogs() ([]models.AuditLog, error) {
	query := `
		SELECT l.id, l.ikas_id, l.user_id, COALESCE(u.display_name, u.username, 'System') as user_name, l.action, l.changes, l.created_at
		FROM ikas_audit_logs l
		LEFT JOIN users u ON l.user_id = u.id
		ORDER BY l.created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all audit logs: %w", err)
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		err := rows.Scan(&log.ID, &log.IkasID, &log.UserID, &log.UserName, &log.Action, &log.Changes, &log.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}
