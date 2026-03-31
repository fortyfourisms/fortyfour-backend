package dto_event

import (
	"ikas/internal/dto"
	"time"
)

type PertanyaanGulihCreatedEvent struct {
	Request   dto.CreatePertanyaanGulihRequest `json:"request"`
	CreatedAt time.Time                        `json:"created_at"`
}

type PertanyaanGulihUpdatedEvent struct {
	ID        int                              `json:"id"`
	Request   dto.UpdatePertanyaanGulihRequest `json:"request"`
	UpdatedAt time.Time                        `json:"updated_at"`
}

type PertanyaanGulihDeletedEvent struct {
	ID        int       `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}
