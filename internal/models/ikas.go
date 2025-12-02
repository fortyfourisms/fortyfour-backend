package models

type Ikas struct {
	ID              int     `json:"id"`
	IDStakeholder   int     `json:"id_stakeholder"`
	Tanggal         string  `json:"tanggal"`
	Responden       string  `json:"responden"`
	Telepon         int     `json:"telepon"`
	Jabatan         string  `json:"jabatan"`
	NilaiKematangan float64 `json:"nilai_kematangan"`
	TargetNilai     float64 `json:"target_nilai"`
	IDIdentifikasi  int     `json:"id_identifikasi"`
	IDProteksi      int     `json:"id_proteksi"`
	IDDeteksi       int     `json:"id_deteksi"`
	IDGulih         int     `json:"id_gulih"`
}