package dto

type CreateIdentifikasiRequest struct {
	NilaiIdentifikasi float64 `json:"nilai_identifikasi"`
	NilaiSubdomain1   float64 `json:"nilai_subdomain1"`
	NilaiSubdomain2   float64 `json:"nilai_subdomain2"`
	NilaiSubdomain3   float64 `json:"nilai_subdomain3"`
	NilaiSubdomain4   float64 `json:"nilai_subdomain4"`
	NilaiSubdomain5   float64 `json:"nilai_subdomain5"`
}

type UpdateIdentifikasiRequest struct {
	NilaiIdentifikasi *float64 `json:"nilai_identifikasi,omitempty"`
	NilaiSubdomain1   *float64 `json:"nilai_subdomain1,omitempty"`
	NilaiSubdomain2   *float64 `json:"nilai_subdomain2,omitempty"`
	NilaiSubdomain3   *float64 `json:"nilai_subdomain3,omitempty"`
	NilaiSubdomain4   *float64 `json:"nilai_subdomain4,omitempty"`
	NilaiSubdomain5   *float64 `json:"nilai_subdomain5,omitempty"`
}

type IdentifikasiResponse struct {
	ID 					string	`json:"id"`
	NilaiIdentifikasi   float64 `json:"nilai_identiifasi"`
	NilaiSubdomain1 	float64 `json:"nilai_subdomain1"`
	NilaiSubdomain2 	float64 `json:"nilai_subdomain2"`
	NilaiSubdomain3 	float64 `json:"nilai_subdomain3"`
	NilaiSubdomain4		float64	`json:"nilai_subdomain4"`
	NilaiSubdomain5		float64	`json:"nilai_subdomain5"`
	CreatedAt       	string  `json:"created_at"`
	UpdatedAt       	string  `json:"updated_at"`
}
