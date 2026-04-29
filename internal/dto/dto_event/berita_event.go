package dto_event

import (
	"fortyfour-backend/internal/dto"
	"time"
)

type BeritaCreatedEvent struct {
	AuthorID  string                  `json:"author_id"`
	Request   dto.CreateBeritaRequest `json:"request"`
	CreatedAt time.Time               `json:"created_at"`
}

type BeritaUpdatedEvent struct {
	ID        int64                   `json:"id"`
	Request   dto.UpdateBeritaRequest `json:"request"`
	UpdatedAt time.Time               `json:"updated_at"`
}

type BeritaDeletedEvent struct {
	ID        int64     `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}
