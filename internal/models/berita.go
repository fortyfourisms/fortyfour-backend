package models

import "time"

type Berita struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Judul     string    `gorm:"type:varchar(255);not null" json:"judul"`
	Deskripsi string    `gorm:"type:longtext;not null" json:"deskripsi"`
	AuthorID  string    `gorm:"type:char(36);not null" json:"author_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Optional: Link to Author (User)
	Author *User `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}
