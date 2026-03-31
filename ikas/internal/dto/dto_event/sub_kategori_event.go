package dto_event

import (
	"ikas/internal/dto"
	"time"
)

type SubKategoriCreatedEvent struct {
	Request   dto.CreateSubKategoriRequest `json:"request"`
	CreatedAt time.Time                    `json:"created_at"`
}

type SubKategoriUpdatedEvent struct {
	ID        int                        `json:"id"`
	Request   dto.UpdateSubKategoriRequest `json:"request"`
	UpdatedAt time.Time                  `json:"updated_at"`
}

type SubKategoriDeletedEvent struct {
	ID        int       `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}
