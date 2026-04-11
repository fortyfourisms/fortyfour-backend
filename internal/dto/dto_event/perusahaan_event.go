package dto_event

import "time"

type PerusahaanCreatedEvent struct {
	ID             string    `json:"id"`
	Photo          string    `json:"photo"`
	NamaPerusahaan string    `json:"nama_perusahaan"`
	IDSubSektor    string    `json:"id_sub_sektor"`
	Alamat         string    `json:"alamat"`
	Telepon        string    `json:"telepon"`
	Email          string    `json:"email"`
	Website        string    `json:"website"`
	CreatedAt      time.Time `json:"created_at"`
}

type PerusahaanUpdatedEvent struct {
	ID             string    `json:"id"`
	Photo          string    `json:"photo"`
	NamaPerusahaan string    `json:"nama_perusahaan"`
	IDSubSektor    string    `json:"id_sub_sektor"`
	Alamat         string    `json:"alamat"`
	Telepon        string    `json:"telepon"`
	Email          string    `json:"email"`
	Website        string    `json:"website"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type PerusahaanDeletedEvent struct {
	ID        string    `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}
