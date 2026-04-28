package models

import "time"

type Event struct {
	ID        int64     `json:"id"`
	Judul     string    `json:"judul"`
	Deskripsi string    `json:"deskripsi"`
	Lokasi    string    `json:"lokasi"`
	Tanggal   time.Time `json:"tanggal"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
