package dto_event

import (
	"ikas/internal/dto"
	"time"
)

type PertanyaanDeteksiCreatedEvent struct {
	Request   dto.CreatePertanyaanDeteksiRequest `json:"request"`
	CreatedAt time.Time                          `json:"created_at"`
}

type PertanyaanDeteksiUpdatedEvent struct {
	ID        int                                `json:"id"`
	Request   dto.UpdatePertanyaanDeteksiRequest `json:"request"`
	UpdatedAt time.Time                          `json:"updated_at"`
}

type PertanyaanDeteksiDeletedEvent struct {
	ID        int       `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}
