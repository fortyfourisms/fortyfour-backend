package dto_event

import "time"

type JabatanCreatedEvent struct {
	ID          string    `json:"id"`
	NamaJabatan string    `json:"nama_jabatan"`
	CreatedAt   time.Time `json:"created_at"`
}

type JabatanUpdatedEvent struct {
	ID          string    `json:"id"`
	NamaJabatan string    `json:"nama_jabatan"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type JabatanDeletedEvent struct {
	ID        string    `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}
