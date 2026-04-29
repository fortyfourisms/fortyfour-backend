package models

import "time"

type Berita struct {
	ID        int64     `json:"id"`
	Judul     string    `json:"judul"`
	Deskripsi string    `json:"deskripsi"`
	AuthorID  string    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Optional: Link to Author (User)
	Author *User `json:"author,omitempty"`
}
