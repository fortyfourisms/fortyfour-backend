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
 		INSERT INTO notifications (id, user_id, type, message, is_read, created_at)
 		VALUES (?, ?, ?, ?, ?, NOW())
 	`
 	_, err := r.db.Exec(query, notif.ID, notif.UserID, notif.Type, notif.Message, notif.Read)
 	return err
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
 
 func (r *NotificationRepository) MarkRead(userID, notifID string) error {
 	query := `UPDATE notifications SET is_read = TRUE WHERE id = ? AND user_id = ?`
 	_, err := r.db.Exec(query, notifID, userID)
 	return err
 }
 
 func (r *NotificationRepository) MarkAllRead(userID string) error {
 	query := `UPDATE notifications SET is_read = TRUE WHERE user_id = ?`
 	_, err := r.db.Exec(query, userID)
 	return err
 }
 
 func (r *NotificationRepository) Delete(userID, notifID string) error {
 	query := `DELETE FROM notifications WHERE id = ? AND user_id = ?`
 	_, err := r.db.Exec(query, notifID, userID)
 	return err
 }
 
 func (r *NotificationRepository) DeleteAllByUserID(userID string) error {
 	query := `DELETE FROM notifications WHERE user_id = ?`
 	_, err := r.db.Exec(query, userID)
 	return err
 }
