package models

import "time"

// NotificationType mendefinisikan jenis notifikasi
type NotificationType string

const (
	NotifLoginFailed        NotificationType = "login_failed"
	NotifPasswordExpirySoon NotificationType = "password_expiry_soon"
	NotifPasswordExpired    NotificationType = "password_expired"
	NotifAccountSuspended   NotificationType = "account_suspended"
	NotifSTRExpirySoon      NotificationType = "str_expiry_soon"
	NotifSTRExpired         NotificationType = "str_expired"
)

// Notification adalah struktur notifikasi yang disimpan di Redis
type Notification struct {
	ID        string           `json:"id"`
	UserID    string           `json:"user_id"`
	Type      NotificationType `json:"type"`
	Message   string           `json:"message"`
	Read      bool             `json:"read"`
	CreatedAt time.Time        `json:"created_at"`
}
