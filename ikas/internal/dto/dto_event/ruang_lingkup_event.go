package dto_event

import (
	"ikas/internal/dto"
	"time"
)

type RuangLingkupCreatedEvent struct {
	Request   dto.CreateRuangLingkupRequest `json:"request"`
	CreatedAt time.Time                     `json:"created_at"`
}

type RuangLingkupUpdatedEvent struct {
	ID        int                           `json:"id"`
	Request   dto.UpdateRuangLingkupRequest `json:"request"`
	UpdatedAt time.Time                     `json:"updated_at"`
}

type RuangLingkupDeletedEvent struct {
	ID        int       `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}
