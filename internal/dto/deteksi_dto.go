package dto

type CreateDeteksiRequest struct {
	NilaiDeteksi    float64 `json:"nilai_deteksi"`
	NilaiSubdomain1 float64 `json:"nilai_subdomain1"`
	NilaiSubdomain2 float64 `json:"nilai_subdomain2"`
	NilaiSubdomain3 float64 `json:"nilai_subdomain3"`
}

type UpdateDeteksiRequest struct {
	NilaiDeteksi    *float64 `json:"nilai_deteksi,omitempty"`
	NilaiSubdomain1 *float64 `json:"nilai_subdomain1,omitempty"`
	NilaiSubdomain2 *float64 `json:"nilai_subdomain2,omitempty"`
	NilaiSubdomain3 *float64 `json:"nilai_subdomain3,omitempty"`
}

type DeteksiResponse struct {
	ID              string  `json:"id"`
	NilaiDeteksi    float64 `json:"nilai_deteksi"`
	NilaiSubdomain1 float64 `json:"nilai_subdomain1"`
	NilaiSubdomain2 float64 `json:"nilai_subdomain2"`
	NilaiSubdomain3 float64 `json:"nilai_subdomain3"`
}