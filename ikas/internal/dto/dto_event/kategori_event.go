package dto_event

import (
	"ikas/internal/dto"
	"time"
)

type KategoriCreatedEvent struct {
	Request   dto.CreateKategoriRequest `json:"request"`
	CreatedAt time.Time                 `json:"created_at"`
}

type KategoriUpdatedEvent struct {
	ID        int                        `json:"id"`
	Request   dto.UpdateKategoriRequest `json:"request"`
	UpdatedAt time.Time                `json:"updated_at"`
}

type KategoriDeletedEvent struct {
	ID        int       `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}
