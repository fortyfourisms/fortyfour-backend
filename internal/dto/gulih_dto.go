package dto

type CreateGulihRequest struct {
	NilaiGulih      float64 `json:"nilai_gulih"`
	NilaiSubdomain1 float64 `json:"nilai_subdomain1"`
	NilaiSubdomain2 float64 `json:"nilai_subdomain2"`
	NilaiSubdomain3 float64 `json:"nilai_subdomain3"`
	NilaiSubdomain4 float64 `json:"nilai_subdomain4"`
}

type UpdateGulihRequest struct {
	NilaiGulih      *float64 `json:"nilai_gulih,omitempty"`
	NilaiSubdomain1 *float64 `json:"nilai_subdomain1,omitempty"`
	NilaiSubdomain2 *float64 `json:"nilai_subdomain2,omitempty"`
	NilaiSubdomain3 *float64 `json:"nilai_subdomain3,omitempty"`
	NilaiSubdomain4 *float64 `json:"nilai_subdomain4,omitempty"`
}

type GulihResponse struct {
	ID 				string	`json:"id"`
	NilaiGulih   	float64 `json:"nilai_deteksi"`
	NilaiSubdomain1 float64 `json:"nilai_subdomain1"`
	NilaiSubdomain2 float64 `json:"nilai_subdomain2"`
	NilaiSubdomain3 float64 `json:"nilai_subdomain3"`
	NilaiSubdomain4	float64	`json:"nilai_subdomain4"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}