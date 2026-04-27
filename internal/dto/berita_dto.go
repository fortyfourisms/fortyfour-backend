package dto

type CreateBeritaRequest struct {
	Judul     string `json:"judul" validate:"required,min=5,max=255"`
	Deskripsi string `json:"deskripsi" validate:"required"`
}

type UpdateBeritaRequest struct {
	Judul     *string `json:"judul" validate:"omitempty,min=5,max=255"`
	Deskripsi *string `json:"deskripsi" validate:"omitempty"`
}

type BeritaAuthor struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	DisplayName *string `json:"display_name"`
}

type BeritaResponse struct {
	ID        int64         `json:"id"`
	Judul     string        `json:"judul"`
	Deskripsi string        `json:"deskripsi"`
	AuthorID  string        `json:"author_id"`
	Author    *BeritaAuthor `json:"author,omitempty"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
}
