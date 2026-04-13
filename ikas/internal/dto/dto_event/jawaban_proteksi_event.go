package dto_event

import (
	"ikas/internal/dto"
	"time"
)

// JawabanProteksiUpdatedEvent
type JawabanProteksiUpdatedEvent struct {
	ID        int                              `json:"id"`
	Request   dto.UpdateJawabanProteksiRequest `json:"request"`
	UpdatedAt time.Time                        `json:"updated_at"`
}

// JawabanProteksiDeletedEvent
type JawabanProteksiDeletedEvent struct {
	ID        int       `json:"id"`
	IkasID    string    `json:"ikas_id"`
	DeletedAt time.Time `json:"deleted_at"`
}
