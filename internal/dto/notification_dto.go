package dto

// NotificationResponse adalah response notifikasi untuk user
type NotificationResponse struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Message   string `json:"message"`
	Read      bool   `json:"read"`
	CreatedAt string `json:"created_at"`
}

// NotificationListResponse adalah response list notifikasi
type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	UnreadCount   int                    `json:"unread_count"`
}
