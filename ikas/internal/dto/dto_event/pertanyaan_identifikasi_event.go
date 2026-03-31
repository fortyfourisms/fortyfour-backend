package dto_event

import (
	"ikas/internal/dto"
	"time"
)

type PertanyaanIdentifikasiCreatedEvent struct {
	Request   dto.CreatePertanyaanIdentifikasiRequest `json:"request"`
	CreatedAt time.Time                               `json:"created_at"`
}

type PertanyaanIdentifikasiUpdatedEvent struct {
	ID        int                                     `json:"id"`
	Request   dto.UpdatePertanyaanIdentifikasiRequest `json:"request"`
	UpdatedAt time.Time                               `json:"updated_at"`
}

type PertanyaanIdentifikasiDeletedEvent struct {
	ID        int       `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}
