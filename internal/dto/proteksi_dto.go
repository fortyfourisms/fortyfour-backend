package dto

type CreateProteksiRequest struct {
	NilaiSubdomain1 float64 `json:"nilai_subdomain1" validate:"required,min=0,max=10"`
	NilaiSubdomain2 float64 `json:"nilai_subdomain2" validate:"required,min=0,max=10"`
	NilaiSubdomain3 float64 `json:"nilai_subdomain3" validate:"required,min=0,max=10"`
	NilaiSubdomain4 float64 `json:"nilai_subdomain4" validate:"required,min=0,max=10"`
	NilaiSubdomain5 float64 `json:"nilai_subdomain5" validate:"required,min=0,max=10"`
	NilaiSubdomain6 float64 `json:"nilai_subdomain6" validate:"required,min=0,max=10"`
}

type UpdateProteksiRequest struct {
	NilaiSubdomain1 *float64 `json:"nilai_subdomain1,omitempty" validate:"omitempty,min=0,max=10"`
	NilaiSubdomain2 *float64 `json:"nilai_subdomain2,omitempty" validate:"omitempty,min=0,max=10"`
	NilaiSubdomain3 *float64 `json:"nilai_subdomain3,omitempty" validate:"omitempty,min=0,max=10"`
	NilaiSubdomain4 *float64 `json:"nilai_subdomain4,omitempty" validate:"omitempty,min=0,max=10"`
	NilaiSubdomain5 *float64 `json:"nilai_subdomain5,omitempty" validate:"omitempty,min=0,max=10"`
	NilaiSubdomain6 *float64 `json:"nilai_subdomain6,omitempty" validate:"omitempty,min=0,max=10"`
}

type ProteksiResponse struct {
	ID              string  `json:"id"`
	NilaiProteksi   float64 `json:"nilai_proteksi"`
	NilaiSubdomain1 float64 `json:"nilai_subdomain1"`
	NilaiSubdomain2 float64 `json:"nilai_subdomain2"`
	NilaiSubdomain3 float64 `json:"nilai_subdomain3"`
	NilaiSubdomain4 float64 `json:"nilai_subdomain4"`
	NilaiSubdomain5 float64 `json:"nilai_subdomain5"`
	NilaiSubdomain6 float64 `json:"nilai_subdomain6"`
	JumlahSubdomain int     `json:"jumlah_subdomain"`
}