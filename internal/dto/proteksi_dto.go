package dto

type CreateProteksiRequest struct {
	NilaiProteksi   float64 `json:"nilai_proteksi"`
	NilaiSubdomain1 float64 `json:"nilai_subdomain1"`
	NilaiSubdomain2 float64 `json:"nilai_subdomain2"`
	NilaiSubdomain3 float64 `json:"nilai_subdomain3"`
	NilaiSubdomain4 float64 `json:"nilai_subdomain4"`
	NilaiSubdomain5 float64 `json:"nilai_subdomain5"`
	NilaiSubdomain6 float64 `json:"nilai_subdomain6"`
}

type UpdateProteksiRequest struct {
	NilaiProteksi   *float64 `json:"nilai_proteksi,omitempty"`
	NilaiSubdomain1 *float64 `json:"nilai_subdomain1,omitempty"`
	NilaiSubdomain2 *float64 `json:"nilai_subdomain2,omitempty"`
	NilaiSubdomain3 *float64 `json:"nilai_subdomain3,omitempty"`
	NilaiSubdomain4 *float64 `json:"nilai_subdomain4,omitempty"`
	NilaiSubdomain5 *float64 `json:"nilai_subdomain5,omitempty"`
	NilaiSubdomain6 *float64 `json:"nilai_subdomain6,omitempty"`
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
}
