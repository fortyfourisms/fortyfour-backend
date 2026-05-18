package dto_event

import (
	"ikas/internal/dto"
	"time"
)

// JawabanGulihUpdatedEvent
type JawabanGulihUpdatedEvent struct {
	UUID      string                        `json:"id"`
	Request   dto.UpdateJawabanGulihRequest `json:"request"`
	UpdatedAt time.Time                     `json:"updated_at"`
}

// JawabanGulihDeletedEvent
type JawabanGulihDeletedEvent struct {
	UUID      string    `json:"id"`
	IkasID    string    `json:"ikas_id"`
	DeletedAt time.Time `json:"deleted_at"`
}
