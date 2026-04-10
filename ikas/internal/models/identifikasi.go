package models

type Identifikasi struct {
	ID                string  `json:"id"`
	PerusahaanID      string  `json:"perusahaan_id"`
	NilaiIdentifikasi float64 `json:"nilai_identifikasi"`
	NilaiSubdomain1   float64 `json:"nilai_subdomain1"`
	NilaiSubdomain2   float64 `json:"nilai_subdomain2"`
	NilaiSubdomain3   float64 `json:"nilai_subdomain3"`
	NilaiSubdomain4   float64 `json:"nilai_subdomain4"`
	NilaiSubdomain5   float64 `json:"nilai_subdomain5"`
}
