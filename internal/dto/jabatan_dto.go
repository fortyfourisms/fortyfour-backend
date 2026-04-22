package dto

type CreateJabatanRequest struct {
	NamaJabatan *string `json:"nama_jabatan"`
}

type UpdateJabatanRequest struct {
	NamaJabatan *string `json:"nama_jabatan"`
}

type JabatanResponse struct {
	ID          string `json:"id"`
	NamaJabatan string `json:"nama_jabatan"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
