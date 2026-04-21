package dto
 
 // NotificationUserResponse adalah detail user dalam notifikasi
 type NotificationUserResponse struct {
 	UserID      string `json:"user_id"`
 	Username    string `json:"username"`
 	DisplayName string `json:"display_name"`
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
 	Notifications []NotificationResponse `json:"notifications"`
 	UnreadCount   int                    `json:"unread_count"`
 }
