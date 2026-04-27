package dto_event

import (
	"fortyfour-backend/internal/dto"
	"time"
)

type EventCreatedEvent struct {
	Request   dto.CreateEventRequest `json:"request"`
	CreatedAt time.Time              `json:"created_at"`
}

type EventUpdatedEvent struct {
	ID        int64                  `json:"id"`
	Request   dto.UpdateEventRequest `json:"request"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type EventDeletedEvent struct {
	ID        int64     `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}
