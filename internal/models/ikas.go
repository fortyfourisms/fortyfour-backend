package models

type Ikas struct {
	ID              string     `json:"id"`
	IDStakeholder   string     `json:"id_stakeholder"`
	Tanggal         string  `json:"tanggal"`
	Responden       string  `json:"responden"`
	Telepon         string     `json:"telepon"`
	Jabatan         string  `json:"jabatan"`
	NilaiKematangan float64 `json:"nilai_kematangan"`
	TargetNilai     float64 `json:"target_nilai"`
	IDIdentifikasi  string     `json:"id_identifikasi"`
	IDProteksi      string     `json:"id_proteksi"`
	IDDeteksi       string     `json:"id_deteksi"`
	IDGulih         string     `json:"id_gulih"`
}