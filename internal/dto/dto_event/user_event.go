package dto_event

import "time"

// UserCreatedEvent
type UserCreatedEvent struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	RoleID    string    `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

// UserUpdatedEvent
type UserUpdatedEvent struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	RoleID    string    `json:"role_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserDeletedEvent
type UserDeletedEvent struct {
	ID        string    `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}

// UserPasswordUpdatedEvent
type UserPasswordUpdatedEvent struct {
	ID        string    `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
}
