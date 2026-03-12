package dto_event

import (
	"ikas/internal/dto"
	"time"
)

// JawabanGulihUpdatedEvent
type JawabanGulihUpdatedEvent struct {
	ID        int                           `json:"id"`
	Request   dto.UpdateJawabanGulihRequest `json:"request"`
	UpdatedAt time.Time                     `json:"updated_at"`
}

// JawabanGulihDeletedEvent
type JawabanGulihDeletedEvent struct {
	ID           int       `json:"id"`
	PerusahaanID string    `json:"perusahaan_id"`
	DeletedAt    time.Time `json:"deleted_at"`
}
