package models

import (
	"time"
)

type Event struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Judul     string    `gorm:"type:varchar(255);not null" json:"judul"`
	Deskripsi string    `gorm:"type:longtext;not null" json:"deskripsi"`
	Tanggal   time.Time `gorm:"not null" json:"tanggal"`
	Lokasi    string    `gorm:"type:varchar(255);not null" json:"lokasi"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
