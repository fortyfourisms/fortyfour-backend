package dto_event

import (
	"fortyfour-backend/internal/dto"
	"time"
)

type AktivitasCreatedEvent struct {
	Request   dto.CreateAktivitasRequest `json:"request"`
	CreatedAt time.Time                  `json:"created_at"`
}

type AktivitasUpdatedEvent struct {
	ID        int                        `json:"id"`
	Request   dto.UpdateAktivitasRequest `json:"request"`
	UpdatedAt time.Time                  `json:"updated_at"`
}

type AktivitasDeletedEvent struct {
	ID        int       `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}
