package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"ikas/internal/dto/dto_event"
	"ikas/internal/models"
)

// AuditLogRepositoryInterface defines the contract for audit log persistence.
type AuditLogRepositoryInterface interface {
	SaveAuditLog(event dto_event.IkasAuditLogEvent) error
	// GetAuditLogs returns a paginated list of all audit logs and the total count.
	GetAuditLogs(offset, limit int) ([]models.AuditLog, int, error)
	// GetAuditLogsByIkasID returns a paginated list of audit logs filtered by IKAS ID and the total count.
	GetAuditLogsByIkasID(ikasID string, offset, limit int) ([]models.AuditLog, int, error)
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

// GetAuditLogs retrieves a paginated slice of all audit logs and their total count.
func (r *AuditLogRepository) GetAuditLogs(offset, limit int) ([]models.AuditLog, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM ikas_audit_logs`
	if err := r.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	query := `
		SELECT l.id, l.ikas_id, l.user_id, COALESCE(u.display_name, u.username, 'System') as user_name,
		       l.action, l.changes, l.created_at
		FROM ikas_audit_logs l
		LEFT JOIN users u ON l.user_id = u.id
		ORDER BY l.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch audit logs: %w", err)
	}
	defer rows.Close()

	logs, err := scanAuditLogs(rows)
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetAuditLogsByIkasID retrieves a paginated slice of audit logs for a specific IKAS record and their total count.
func (r *AuditLogRepository) GetAuditLogsByIkasID(ikasID string, offset, limit int) ([]models.AuditLog, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM ikas_audit_logs WHERE ikas_id = ?`
	if err := r.db.QueryRow(countQuery, ikasID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	query := `
		SELECT l.id, l.ikas_id, l.user_id, COALESCE(u.display_name, u.username, 'System') as user_name,
		       l.action, l.changes, l.created_at
		FROM ikas_audit_logs l
		LEFT JOIN users u ON l.user_id = u.id
		WHERE l.ikas_id = ?
		ORDER BY l.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, ikasID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch audit logs by ikas_id: %w", err)
	}
	defer rows.Close()

	logs, err := scanAuditLogs(rows)
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// scanAuditLogs is a shared helper to scan sql.Rows into []models.AuditLog.
func scanAuditLogs(rows *sql.Rows) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		if err := rows.Scan(&log.ID, &log.IkasID, &log.UserID, &log.UserName, &log.Action, &log.Changes, &log.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}
		logs = append(logs, log)
	}
	return logs, nil
}
