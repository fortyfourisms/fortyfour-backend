package dto_event

import "time"

type PicCreatedEvent struct {
	ID           string    `json:"id"`
	Nama         string    `json:"nama"`
	Telepon      string    `json:"telepon"`
	IDPerusahaan string    `json:"id_perusahaan"`
	CreatedAt    time.Time `json:"created_at"`
}

type PicUpdatedEvent struct {
	ID           string    `json:"id"`
	Nama         string    `json:"nama"`
	Telepon      string    `json:"telepon"`
	IDPerusahaan string    `json:"id_perusahaan"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PicDeletedEvent struct {
	ID        string    `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}
