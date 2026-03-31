package dto_event

import (
	"ikas/internal/dto"
	"time"
)

type PertanyaanProteksiCreatedEvent struct {
	Request   dto.CreatePertanyaanProteksiRequest `json:"request"`
	CreatedAt time.Time                           `json:"created_at"`
}

type PertanyaanProteksiUpdatedEvent struct {
	ID        int                                 `json:"id"`
	Request   dto.UpdatePertanyaanProteksiRequest `json:"request"`
	UpdatedAt time.Time                           `json:"updated_at"`
}

type PertanyaanProteksiDeletedEvent struct {
	ID        int       `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}
