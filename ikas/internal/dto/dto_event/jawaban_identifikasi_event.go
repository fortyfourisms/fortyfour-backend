package dto_event

import (
	"ikas/internal/dto"
	"time"
)

// JawabanIdentifikasiUpdatedEvent
type JawabanIdentifikasiUpdatedEvent struct {
	ID        int                                  `json:"id"`
	Request   dto.UpdateJawabanIdentifikasiRequest `json:"request"`
	UpdatedAt time.Time                            `json:"updated_at"`
}

// JawabanIdentifikasiDeletedEvent
type JawabanIdentifikasiDeletedEvent struct {
	ID        int       `json:"id"`
	IkasID    string    `json:"ikas_id"`
	DeletedAt time.Time `json:"deleted_at"`
}
