package dto_event

import (
	"ikas/internal/dto"
	"time"
)

type DomainCreatedEvent struct {
	Request   dto.CreateDomainRequest `json:"request"`
	CreatedAt time.Time              `json:"created_at"`
}

type DomainUpdatedEvent struct {
	ID        int                     `json:"id"`
	Request   dto.UpdateDomainRequest `json:"request"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type DomainDeletedEvent struct {
	ID        int       `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}
