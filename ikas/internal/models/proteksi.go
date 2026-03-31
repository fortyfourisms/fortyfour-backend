package models

type Proteksi struct {
	ID              string  `json:"id"`
	NilaiProteksi   float64 `json:"nilai_proteksi"`
	NilaiSubdomain1 float64 `json:"nilai_subdomain1"`
	NilaiSubdomain2 float64 `json:"nilai_subdomain2"`
	NilaiSubdomain3 float64 `json:"nilai_subdomain3"`
	NilaiSubdomain4 float64 `json:"nilai_subdomain4"`
	NilaiSubdomain5 float64 `json:"nilai_subdomain5"`
	NilaiSubdomain6 float64 `json:"nilai_subdomain6"`
}
