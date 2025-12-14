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