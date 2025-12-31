package dto

type CreateIkasRequest struct {
	IDPerusahaan string  `json:"id_perusahaan"`
	Tanggal      string  `json:"tanggal"`
	Responden    string  `json:"responden"`
	Telepon      string  `json:"telepon"`
	Jabatan      string  `json:"jabatan"`
	TargetNilai  float64 `json:"target_nilai"`

	// Nested data untuk auto-create
	Identifikasi *CreateIdentifikasiData `json:"identifikasi,omitempty"`
	Proteksi     *CreateProteksiData     `json:"proteksi,omitempty"`
	Deteksi      *CreateDeteksiData      `json:"deteksi,omitempty"`
	Gulih        *CreateGulihData        `json:"gulih,omitempty"`

	// ID manual (backward compatibility)
	IDIdentifikasi string `json:"id_identifikasi,omitempty"`
	IDProteksi     string `json:"id_proteksi,omitempty"`
	IDDeteksi      string `json:"id_deteksi,omitempty"`
	IDGulih        string `json:"id_gulih,omitempty"`
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

	// Nested data untuk update subdomain
	Identifikasi *UpdateIdentifikasiData `json:"identifikasi,omitempty"`
	Proteksi     *UpdateProteksiData     `json:"proteksi,omitempty"`
	Deteksi      *UpdateDeteksiData      `json:"deteksi,omitempty"`
	Gulih        *UpdateGulihData        `json:"gulih,omitempty"`
}

// Response dengan nested objects
type IkasResponse struct {
	ID                              string              `json:"id"`
	Tanggal                         string              `json:"tanggal"`
	Responden                       string              `json:"responden"`
	Telepon                         string              `json:"telepon"`
	Jabatan                         string              `json:"jabatan"`
	NilaiKematangan                 float64             `json:"nilai_kematangan"`
	KategoriKematanganKeamananSiber string              `json:"kategori_kematangan_keamanan_siber"`
	TargetNilai                     float64             `json:"target_nilai"`
	Perusahaan                      *PerusahaanInIkas   `json:"perusahaan,omitempty"`
	Identifikasi                    *IdentifikasiInIkas `json:"identifikasi,omitempty"`
	Proteksi                        *ProteksiInIkas     `json:"proteksi,omitempty"`
	Deteksi                         *DeteksiInIkas      `json:"deteksi,omitempty"`
	Gulih                           *GulihInIkas        `json:"gulih,omitempty"`
}

// Nested structs untuk foreign keys
type PerusahaanInIkas struct {
	ID             string `json:"id"`
	NamaPerusahaan string `json:"nama_perusahaan"`
}

type IdentifikasiInIkas struct {
	ID                              string  `json:"id"`
	NilaiIdentifikasi               float64 `json:"nilai_identifikasi"`
	KategoriTingkatKematanganDomain string  `json:"kategori_tingkat_kematangan_domain"`
	NilaiSubdomain1                 float64 `json:"nilai_subdomain1"`
	NilaiSubdomain2                 float64 `json:"nilai_subdomain2"`
	NilaiSubdomain3                 float64 `json:"nilai_subdomain3"`
	NilaiSubdomain4                 float64 `json:"nilai_subdomain4"`
	NilaiSubdomain5                 float64 `json:"nilai_subdomain5"`
}

type ProteksiInIkas struct {
	ID                              string  `json:"id"`
	NilaiProteksi                   float64 `json:"nilai_proteksi"`
	KategoriTingkatKematanganDomain string  `json:"kategori_tingkat_kematangan_domain"`
	NilaiSubdomain1                 float64 `json:"nilai_subdomain1"`
	NilaiSubdomain2                 float64 `json:"nilai_subdomain2"`
	NilaiSubdomain3                 float64 `json:"nilai_subdomain3"`
	NilaiSubdomain4                 float64 `json:"nilai_subdomain4"`
	NilaiSubdomain5                 float64 `json:"nilai_subdomain5"`
	NilaiSubdomain6                 float64 `json:"nilai_subdomain6"`
}

type DeteksiInIkas struct {
	ID                              string  `json:"id"`
	NilaiDeteksi                    float64 `json:"nilai_deteksi"`
	KategoriTingkatKematanganDomain string  `json:"kategori_tingkat_kematangan_domain"`
	NilaiSubdomain1                 float64 `json:"nilai_subdomain1"`
	NilaiSubdomain2                 float64 `json:"nilai_subdomain2"`
	NilaiSubdomain3                 float64 `json:"nilai_subdomain3"`
}

type GulihInIkas struct {
	ID                              string  `json:"id"`
	NilaiGulih                      float64 `json:"nilai_gulih"`
	KategoriTingkatKematanganDomain string  `json:"kategori_tingkat_kematangan_domain"`
	NilaiSubdomain1                 float64 `json:"nilai_subdomain1"`
	NilaiSubdomain2                 float64 `json:"nilai_subdomain2"`
	NilaiSubdomain3                 float64 `json:"nilai_subdomain3"`
	NilaiSubdomain4                 float64 `json:"nilai_subdomain4"`
}

// Tambahkan struct baru untuk nested create
type CreateIdentifikasiData struct {
	NilaiSubdomain1 float64 `json:"nilai_subdomain1"`
	NilaiSubdomain2 float64 `json:"nilai_subdomain2"`
	NilaiSubdomain3 float64 `json:"nilai_subdomain3"`
	NilaiSubdomain4 float64 `json:"nilai_subdomain4"`
	NilaiSubdomain5 float64 `json:"nilai_subdomain5"`
}

type CreateProteksiData struct {
	NilaiSubdomain1 float64 `json:"nilai_subdomain1"`
	NilaiSubdomain2 float64 `json:"nilai_subdomain2"`
	NilaiSubdomain3 float64 `json:"nilai_subdomain3"`
	NilaiSubdomain4 float64 `json:"nilai_subdomain4"`
	NilaiSubdomain5 float64 `json:"nilai_subdomain5"`
	NilaiSubdomain6 float64 `json:"nilai_subdomain6"`
}

type CreateDeteksiData struct {
	NilaiSubdomain1 float64 `json:"nilai_subdomain1"`
	NilaiSubdomain2 float64 `json:"nilai_subdomain2"`
	NilaiSubdomain3 float64 `json:"nilai_subdomain3"`
}

type CreateGulihData struct {
	NilaiSubdomain1 float64 `json:"nilai_subdomain1"`
	NilaiSubdomain2 float64 `json:"nilai_subdomain2"`
	NilaiSubdomain3 float64 `json:"nilai_subdomain3"`
	NilaiSubdomain4 float64 `json:"nilai_subdomain4"`
}

// Struct untuk update nested data
type UpdateIdentifikasiData struct {
	NilaiSubdomain1 *float64 `json:"nilai_subdomain1,omitempty"`
	NilaiSubdomain2 *float64 `json:"nilai_subdomain2,omitempty"`
	NilaiSubdomain3 *float64 `json:"nilai_subdomain3,omitempty"`
	NilaiSubdomain4 *float64 `json:"nilai_subdomain4,omitempty"`
	NilaiSubdomain5 *float64 `json:"nilai_subdomain5,omitempty"`
}

type UpdateProteksiData struct {
	NilaiSubdomain1 *float64 `json:"nilai_subdomain1,omitempty"`
	NilaiSubdomain2 *float64 `json:"nilai_subdomain2,omitempty"`
	NilaiSubdomain3 *float64 `json:"nilai_subdomain3,omitempty"`
	NilaiSubdomain4 *float64 `json:"nilai_subdomain4,omitempty"`
	NilaiSubdomain5 *float64 `json:"nilai_subdomain5,omitempty"`
	NilaiSubdomain6 *float64 `json:"nilai_subdomain6,omitempty"`
}

type UpdateDeteksiData struct {
	NilaiSubdomain1 *float64 `json:"nilai_subdomain1,omitempty"`
	NilaiSubdomain2 *float64 `json:"nilai_subdomain2,omitempty"`
	NilaiSubdomain3 *float64 `json:"nilai_subdomain3,omitempty"`
}

type UpdateGulihData struct {
	NilaiSubdomain1 *float64 `json:"nilai_subdomain1,omitempty"`
	NilaiSubdomain2 *float64 `json:"nilai_subdomain2,omitempty"`
	NilaiSubdomain3 *float64 `json:"nilai_subdomain3,omitempty"`
	NilaiSubdomain4 *float64 `json:"nilai_subdomain4,omitempty"`
}

// Import Excel
type ImportIkasRequest struct {
	IDPerusahaan string `json:"id_perusahaan"`
	Tanggal      string `json:"tanggal"`
	Responden    string `json:"responden"`
	Telepon      string `json:"telepon"`
	Jabatan      string `json:"jabatan"`
}

type ImportIkasResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    *IkasResponse `json:"data,omitempty"`
	Errors  []string      `json:"errors,omitempty"`
}
