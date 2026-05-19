package models

import "time"

type Event struct {
	ID        string    `json:"id"`
	Slug      string    `json:"slug"`
	Judul     string    `json:"judul"`
	Deskripsi string    `json:"deskripsi"`
	Lokasi    string    `json:"lokasi"`
	Tanggal   time.Time `json:"tanggal"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
