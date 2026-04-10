package models

type Deteksi struct {
	ID              string  `json:"id"`
	PerusahaanID    string  `json:"perusahaan_id"`
	NilaiDeteksi    float64 `json:"nilai_deteksi"`
	NilaiSubdomain1 float64 `json:"nilai_subdomain1"`
	NilaiSubdomain2 float64 `json:"nilai_subdomain2"`
	NilaiSubdomain3 float64 `json:"nilai_subdomain3"`
}
