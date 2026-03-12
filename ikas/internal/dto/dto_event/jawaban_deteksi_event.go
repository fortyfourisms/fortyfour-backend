package dto_event

import (
	"ikas/internal/dto"
	"time"
)

// JawabanDeteksiUpdatedEvent
type JawabanDeteksiUpdatedEvent struct {
	ID        int                             `json:"id"`
	Request   dto.UpdateJawabanDeteksiRequest `json:"request"`
	UpdatedAt time.Time                       `json:"updated_at"`
}

// JawabanDeteksiDeletedEvent
type JawabanDeteksiDeletedEvent struct {
	ID           int       `json:"id"`
	PerusahaanID string    `json:"perusahaan_id"`
	DeletedAt    time.Time `json:"deleted_at"`
}
