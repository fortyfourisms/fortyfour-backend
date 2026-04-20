package repository
 
 import (
 	"database/sql"
 	"fortyfour-backend/internal/models"
 )
 
 type NotificationRepository struct {
 	db *sql.DB
 }
 
 func NewNotificationRepository(db *sql.DB) *NotificationRepository {
 	return &NotificationRepository{db: db}
 }
 
 func (r *NotificationRepository) Create(notif *models.Notification) error {
 	query := `
 		INSERT INTO notifications (user_id, type, message, is_read, created_at)
 		VALUES (?, ?, ?, ?, NOW())
 	`
 	res, err := r.db.Exec(query, notif.UserID, notif.Type, notif.Message, notif.Read)
 	if err != nil {
 		return err
 	}
 
 	id, err := res.LastInsertId()
 	if err == nil {
 		notif.ID = id
 	}
 	return nil
 }
 
 func (r *NotificationRepository) FindAllByUserID(userID string) ([]models.Notification, error) {
 	query := `
 		SELECT id, user_id, type, message, is_read, created_at
 		FROM notifications
 		WHERE user_id = ?
 		ORDER BY created_at DESC
 	`
 	rows, err := r.db.Query(query, userID)
 	if err != nil {
 		return nil, err
 	}
 	defer rows.Close()
 
 	var notifs []models.Notification
 	for rows.Next() {
 		var n models.Notification
 		err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Message, &n.Read, &n.CreatedAt)
 		if err != nil {
 			return nil, err
 		}
 		notifs = append(notifs, n)
 	}
 	return notifs, nil
 }
 
 func (r *NotificationRepository) MarkRead(userID string, notifID int64) error {
 	query := `UPDATE notifications SET is_read = TRUE WHERE id = ? AND user_id = ?`
 	_, err := r.db.Exec(query, notifID, userID)
 	return err
 }
 
 func (r *NotificationRepository) MarkAllRead(userID string) error {
 	query := `UPDATE notifications SET is_read = TRUE WHERE user_id = ?`
 	_, err := r.db.Exec(query, userID)
 	return err
 }
 
 func (r *NotificationRepository) Delete(userID string, notifID int64) error {
 	query := `DELETE FROM notifications WHERE id = ? AND user_id = ?`
 	_, err := r.db.Exec(query, notifID, userID)
 	return err
 }
 
 func (r *NotificationRepository) DeleteAllByUserID(userID string) error {
 	query := `DELETE FROM notifications WHERE user_id = ?`
 	_, err := r.db.Exec(query, userID)
 	return err
 }
