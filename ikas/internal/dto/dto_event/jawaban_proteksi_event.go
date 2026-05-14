package dto_event

import (
	"ikas/internal/dto"
	"time"
)

// JawabanProteksiUpdatedEvent
type JawabanProteksiUpdatedEvent struct {
	UUID      string                           `json:"id"`
	Request   dto.UpdateJawabanProteksiRequest `json:"request"`
	UpdatedAt time.Time                        `json:"updated_at"`
}

// JawabanProteksiDeletedEvent
type JawabanProteksiDeletedEvent struct {
	UUID      string    `json:"id"`
	IkasID    string    `json:"ikas_id"`
	DeletedAt time.Time `json:"deleted_at"`
}
