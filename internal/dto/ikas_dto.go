package dto

type CreateIkasRequest struct {
	IDPerusahaan    string  `json:"id_perusahaan"`
	Tanggal         string  `json:"tanggal"`
	Responden       string  `json:"responden"`
	Telepon         string  `json:"telepon"`
	Jabatan         string  `json:"jabatan"`
	NilaiKematangan float64 `json:"nilai_kematangan"`
	TargetNilai     float64 `json:"target_nilai"`
	IDIdentifikasi  string  `json:"id_identifikasi"`
	IDProteksi      string  `json:"id_proteksi"`
	IDDeteksi       string  `json:"id_deteksi"`
	IDGulih         string  `json:"id_gulih"`
}

type UpdateIkasRequest struct {
	IDPerusahaan    *string  `json:"id_perusahaan,omitempty"`
	Tanggal         *string  `json:"tanggal,omitempty"`
	Responden       *string  `json:"responden,omitempty"`
	Telepon         *string  `json:"telepon,omitempty"`
	Jabatan         *string  `json:"jabatan,omitempty"`
	NilaiKematangan *float64 `json:"nilai_kematangan,omitempty"`
	TargetNilai     *float64 `json:"target_nilai,omitempty"`
	IDIdentifikasi  *string  `json:"id_identifikasi,omitempty"`
	IDProteksi      *string  `json:"id_proteksi,omitempty"`
	IDDeteksi       *string  `json:"id_deteksi,omitempty"`
	IDGulih         *string  `json:"id_gulih,omitempty"`
}

// Response dengan nested objects
type IkasResponse struct {
	ID              string              `json:"id"`
	Tanggal         string              `json:"tanggal"`
	Responden       string              `json:"responden"`
	Telepon         string              `json:"telepon"`
	Jabatan         string              `json:"jabatan"`
	NilaiKematangan float64             `json:"nilai_kematangan"`
	TargetNilai     float64             `json:"target_nilai"`
	Perusahaan      *PerusahaanInIkas   `json:"perusahaan,omitempty"`
	Identifikasi    *IdentifikasiInIkas `json:"identifikasi,omitempty"`
	Proteksi        *ProteksiInIkas     `json:"proteksi,omitempty"`
	Deteksi         *DeteksiInIkas      `json:"deteksi,omitempty"`
	Gulih           *GulihInIkas        `json:"gulih,omitempty"`
}

// Nested structs untuk foreign keys
type PerusahaanInIkas struct {
	ID             string `json:"id"`
	NamaPerusahaan string `json:"nama_perusahaan"`
}

type IdentifikasiInIkas struct {
	ID                string  `json:"id"`
	NilaiIdentifikasi float64 `json:"nilai_identifikasi"`
	// NilaiSubdomain1   float64 `json:"nilai_subdomain1"`
	// NilaiSubdomain2   float64 `json:"nilai_subdomain2"`
	// NilaiSubdomain3   float64 `json:"nilai_subdomain3"`
	// NilaiSubdomain4   float64 `json:"nilai_subdomain4"`
	// NilaiSubdomain5   float64 `json:"nilai_subdomain5"`
}

type ProteksiInIkas struct {
	ID            string  `json:"id"`
	NilaiProteksi float64 `json:"nilai_proteksi"`
	// NilaiSubdomain1 float64 `json:"nilai_subdomain1"`
	// NilaiSubdomain2 float64 `json:"nilai_subdomain2"`
	// NilaiSubdomain3 float64 `json:"nilai_subdomain3"`
	// NilaiSubdomain4 float64 `json:"nilai_subdomain4"`
	// NilaiSubdomain5 float64 `json:"nilai_subdomain5"`
	// NilaiSubdomain6 float64 `json:"nilai_subdomain6"`
}

type DeteksiInIkas struct {
	ID           string  `json:"id"`
	NilaiDeteksi float64 `json:"nilai_deteksi"`
	// NilaiSubdomain1 float64 `json:"nilai_subdomain1"`
	// NilaiSubdomain2 float64 `json:"nilai_subdomain2"`
	// NilaiSubdomain3 float64 `json:"nilai_subdomain3"`
}

type GulihInIkas struct {
	ID         string  `json:"id"`
	NilaiGulih float64 `json:"nilai_gulih"`
	// NilaiSubdomain1 float64 `json:"nilai_subdomain1"`
	// NilaiSubdomain2 float64 `json:"nilai_subdomain2"`
	// NilaiSubdomain3 float64 `json:"nilai_subdomain3"`
	// NilaiSubdomain4 float64 `json:"nilai_subdomain4"`
}
