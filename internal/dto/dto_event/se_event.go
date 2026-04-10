package dto_event

import (
	"time"
)

type SeCreatedEvent struct {
	ID                              string    `json:"id"`
	IDPerusahaan                    string    `json:"id_perusahaan"`
	IDSubSektor                     string    `json:"id_sub_sektor"`
	IDCsirt                         string    `json:"id_csirt"`
	NilaiInvestasi                  string    `json:"nilai_investasi"`
	AnggaranOperasional             string    `json:"anggaran_operasional"`
	KepatuhanPeraturan              string    `json:"kepatuhan_peraturan"`
	TeknikKriptografi               string    `json:"teknik_kriptografi"`
	JumlahPengguna                  string    `json:"jumlah_pengguna"`
	DataPribadi                     string    `json:"data_pribadi"`
	KlasifikasiData                 string    `json:"klasifikasi_data"`
	KekritisanProses                string    `json:"kekritisan_proses"`
	DampakKegagalan                 string    `json:"dampak_kegagalan"`
	PotensiKerugiandanDampakNegatif string    `json:"potensi_kerugian_dan_dampak_negatif"`
	NamaSE                          string    `json:"nama_se"`
	IpSE                            string    `json:"ip_se"`
	AsNumberSE                      string    `json:"as_number_se"`
	PengelolaSE                     string    `json:"pengelola_se"`
	FiturSE                         string    `json:"fitur_se"`
	TotalBobot                      int       `json:"total_bobot"`
	KategoriSE                      string    `json:"kategori_se"`
	CreatedAt                       time.Time `json:"created_at"`
}

type SeUpdatedEvent struct {
	ID                              string    `json:"id"`
	IDPerusahaan                    string    `json:"id_perusahaan"`
	IDSubSektor                     string    `json:"id_sub_sektor"`
	IDCsirt                         string    `json:"id_csirt"`
	NilaiInvestasi                  string    `json:"nilai_investasi"`
	AnggaranOperasional             string    `json:"anggaran_operasional"`
	KepatuhanPeraturan              string    `json:"kepatuhan_peraturan"`
	TeknikKriptografi               string    `json:"teknik_kriptografi"`
	JumlahPengguna                  string    `json:"jumlah_pengguna"`
	DataPribadi                     string    `json:"data_pribadi"`
	KlasifikasiData                 string    `json:"klasifikasi_data"`
	KekritisanProses                string    `json:"kekritisan_proses"`
	DampakKegagalan                 string    `json:"dampak_kegagalan"`
	PotensiKerugiandanDampakNegatif string    `json:"potensi_kerugian_dan_dampak_negatif"`
	NamaSE                          string    `json:"nama_se"`
	IpSE                            string    `json:"ip_se"`
	AsNumberSE                      string    `json:"as_number_se"`
	PengelolaSE                     string    `json:"pengelola_se"`
	FiturSE                         string    `json:"fitur_se"`
	TotalBobot                      int       `json:"total_bobot"`
	KategoriSE                      string    `json:"kategori_se"`
	UpdatedAt                       time.Time `json:"updated_at"`
}

type SeDeletedEvent struct {
	ID           string    `json:"id"`
	IDPerusahaan string    `json:"id_perusahaan"`
	DeletedAt    time.Time `json:"deleted_at"`
}
