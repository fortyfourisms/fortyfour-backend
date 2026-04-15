package models

type Gulih struct {
	ID              string  `json:"id"`
	IkasID          string  `json:"ikas_id"`
	NilaiGulih      float64 `json:"nilai_gulih"`
	NilaiSubdomain1 float64 `json:"nilai_subdomain1"`
	NilaiSubdomain2 float64 `json:"nilai_subdomain2"`
	NilaiSubdomain3 float64 `json:"nilai_subdomain3"`
	NilaiSubdomain4 float64 `json:"nilai_subdomain4"`
}
