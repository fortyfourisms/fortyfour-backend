package dto

type CreatePerusahaanRequest struct {
	Photo          *string `json:"photo,omitempty"`
	NamaPerusahaan *string `json:"nama_perusahaan,omitempty"`
	IDSubSektor    *string `json:"id_sub_sektor,omitempty"`
	Alamat         *string `json:"alamat,omitempty"`
	Telepon        *string `json:"telepon,omitempty"`
	Email          *string `json:"email,omitempty"`
	Website        *string `json:"website,omitempty"`
}

type UpdatePerusahaanRequest struct {
	Photo          *string `json:"photo,omitempty"`
	NamaPerusahaan *string `json:"nama_perusahaan,omitempty"`
	IDSubSektor    *string `json:"id_sub_sektor,omitempty"`
	Alamat         *string `json:"alamat,omitempty"`
	Telepon        *string `json:"telepon,omitempty"`
	Email          *string `json:"email,omitempty"`
	Website        *string `json:"website,omitempty"`
}

type PerusahaanResponse struct {
	ID             string             `json:"id"`
	Photo          string             `json:"photo"`
	NamaPerusahaan string             `json:"nama_perusahaan"`
	SubSektor      *SubSektorResponse `json:"sub_sektor,omitempty"`
	Alamat         string             `json:"alamat"`
	Telepon        string             `json:"telepon"`
	Email          string             `json:"email"`
	Website        string             `json:"website"`
	CreatedAt      string             `json:"created_at"`
	UpdatedAt      string             `json:"updated_at"`
}

// PerusahaanMinimalResponse untuk dropdown di halaman register (public endpoint)
// Hanya return data minimal yang diperlukan untuk keamanan
type PerusahaanMinimalResponse struct {
	ID             string `json:"id"`
	NamaPerusahaan string `json:"nama_perusahaan"`
}

// ToMinimal converts full PerusahaanResponse to minimal version
func (p *PerusahaanResponse) ToMinimal() PerusahaanMinimalResponse {
	return PerusahaanMinimalResponse{
		ID:             p.ID,
		NamaPerusahaan: p.NamaPerusahaan,
	}
}

// ConvertToMinimalList converts slice of full response to minimal response list
func ConvertToMinimalList(data []PerusahaanResponse) []PerusahaanMinimalResponse {
	result := make([]PerusahaanMinimalResponse, len(data))
	for i, item := range data {
		result[i] = item.ToMinimal()
	}
	return result
}