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
	ID        string                 `json:"id"`
	Request   dto.UpdateEventRequest `json:"request"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type EventDeletedEvent struct {
	ID        string    `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}

type EventRegistrationCreatedEvent struct {
	EventID   string                             `json:"event_id"`
	Request   dto.CreateEventRegistrationRequest `json:"request"`
	ID        string                             `json:"id"`
	QRToken   string                             `json:"qr_token"`
	QRPayload string                             `json:"qr_payload"`
	CreatedAt time.Time                          `json:"created_at"`
}
