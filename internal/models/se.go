package models

import "time"

type SE struct {
	ID           string `json:"id"`
	IDPerusahaan string `json:"id_perusahaan"`
	IDSubSektor  string `json:"id_sub_sektor"`
	IDCsirt      string `json:"id_csirt"`

	// Dari form kategorisasi - Karakteristik Instansi
	NilaiInvestasi             string `json:"nilai_investasi"`              // Q1: Nilai investasi sistem elektronik
	AnggaranOperasional        string `json:"anggaran_operasional"`         // Q2: Total anggaran operasional tahunan
	KepatuhanPeraturan         string `json:"kepatuhan_peraturan"`          // Q3: Kewajiban kepatuhan terhadap peraturan/standar
	TeknikKriptografi          string `json:"teknik_kriptografi"`           // Q4: Penggunaan teknik kriptografi
	JumlahPengguna             string `json:"jumlah_pengguna"`              // Q5: Jumlah pengguna sistem elektronik
	DataPribadi                string `json:"data_pribadi"`                 // Q6: Data pribadi yang dikelola
	KlasifikasiData            string `json:"klasifikasi_data"`             // Q7: Tingkat klasifikasi/kekritisan data
	KekritisanProses           string `json:"kekritisan_proses"`            // Q8: Tingkat kekritisan proses
	DampakKegagalan            string `json:"dampak_kegagalan"`             // Q9: Dampak dari kegagalan sistem elektronik
	PotensiKerugiandanDampakNegatif string `json:"potensi_kerugian_dan_dampak_negatif"` // Q10: Potensi kerugian dari insiden keamanan

	// Dari SE CSIRT
	NamaSE       string `json:"nama_se"`
	IpSE         string `json:"ip_se"`
	AsNumberSE   string `json:"as_number_se"`
	PengelolaSE  string `json:"pengelola_se"`
	FiturSE      string `json:"fitur_se"`

	// Hasil kalkulasi
	TotalBobot int    `json:"total_bobot"`
	KategoriSE string `json:"kategori_se"` // ENUM: 'Strategis', 'Tinggi', 'Rendah'

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}