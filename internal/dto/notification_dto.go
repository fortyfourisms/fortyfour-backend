package dto

// NotificationUserResponse adalah detail user dalam notifikasi
type NotificationUserResponse struct {
	UserID      string  `json:"user_id"`
	Username    string  `json:"username"`
	DisplayName string  `json:"display_name"`
	FotoProfile *string `json:"foto_profile"`
}

// NotificationResponse adalah response notifikasi untuk user
type NotificationResponse struct {
	ID        int64                    `json:"id"`
	User      NotificationUserResponse `json:"user"`
	Type      string                   `json:"type"`
	Message   string                   `json:"message"`
	Read      bool                     `json:"read"`
	CreatedAt string                   `json:"created_at"`
}

// NotificationListResponse adalah response list notifikasi
type NotificationListResponse struct {
	Message       string                 `json:"message"`
	TotalData     int                    `json:"total_data"`
	ReadCount     int                    `json:"read_count"`
	UnreadCount   int                    `json:"unread_count"`
	Notifications []NotificationResponse `json:"notifications"`
}
