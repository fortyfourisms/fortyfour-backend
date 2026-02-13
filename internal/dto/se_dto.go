package dto

// CreateSERequest represents the request to create a new SE
type CreateSERequest struct {
	IDPerusahaan string `json:"id_perusahaan" validate:"required,uuid4"`
	IDSubSektor  string `json:"id_sub_sektor" validate:"omitempty,uuid4"`
	IDCsirt      string `json:"id_csirt" validate:"omitempty,uuid4"`

	// Karakteristik Instansi (A/B/C)
	NilaiInvestasi             string `json:"nilai_investasi" validate:"required,oneof=A B C"`
	AnggaranOperasional        string `json:"anggaran_operasional" validate:"required,oneof=A B C"`
	KepatuhanPeraturan         string `json:"kepatuhan_peraturan" validate:"required,oneof=A B C"`
	TeknikKriptografi          string `json:"teknik_kriptografi" validate:"required,oneof=A B C"`
	JumlahPengguna             string `json:"jumlah_pengguna" validate:"required,oneof=A B C"`
	DataPribadi                string `json:"data_pribadi" validate:"required,oneof=A B C"`
	KlasifikasiData            string `json:"klasifikasi_data" validate:"required,oneof=A B C"`
	KekritisanProses           string `json:"kekritisan_proses" validate:"required,oneof=A B C"`
	DampakKegagalan            string `json:"dampak_kegagalan" validate:"required,oneof=A B C"`
	PotensiKerugiandanDampakNegatif string `json:"potensi_kerugian_dan_dampak_negatif" validate:"required,oneof=A B C"`

	// Informasi SE
	NamaSE      string `json:"nama_se" validate:"required"`
	IpSE        string `json:"ip_se" validate:"required"`
	AsNumberSE  string `json:"as_number_se" validate:"required"`
	PengelolaSE string `json:"pengelola_se" validate:"required"`
	FiturSE     string `json:"fitur_se"`
}

// UpdateSERequest represents the request to update an existing SE
type UpdateSERequest struct {
	IDPerusahaan *string `json:"id_perusahaan" validate:"omitempty,uuid4"`
	IDSubSektor  *string `json:"id_sub_sektor" validate:"omitempty,uuid4"`
	IDCsirt      *string `json:"id_csirt" validate:"omitempty,uuid4"`

	// Karakteristik Instansi (A/B/C)
	NilaiInvestasi             *string `json:"nilai_investasi" validate:"omitempty,oneof=A B C"`
	AnggaranOperasional        *string `json:"anggaran_operasional" validate:"omitempty,oneof=A B C"`
	KepatuhanPeraturan         *string `json:"kepatuhan_peraturan" validate:"omitempty,oneof=A B C"`
	TeknikKriptografi          *string `json:"teknik_kriptografi" validate:"omitempty,oneof=A B C"`
	JumlahPengguna             *string `json:"jumlah_pengguna" validate:"omitempty,oneof=A B C"`
	DataPribadi                *string `json:"data_pribadi" validate:"omitempty,oneof=A B C"`
	KlasifikasiData            *string `json:"klasifikasi_data" validate:"omitempty,oneof=A B C"`
	KekritisanProses           *string `json:"kekritisan_proses" validate:"omitempty,oneof=A B C"`
	DampakKegagalan            *string `json:"dampak_kegagalan" validate:"omitempty,oneof=A B C"`
	PotensiKerugiandanDampakNegatif *string `json:"potensi_kerugian_dan_dampak_negatif" validate:"omitempty,oneof=A B C"`

	// Informasi SE
	NamaSE      *string `json:"nama_se"`
	IpSE        *string `json:"ip_se"`
	AsNumberSE  *string `json:"as_number_se"`
	PengelolaSE *string `json:"pengelola_se"`
	FiturSE     *string `json:"fitur_se"`
}

// SEResponse represents the response for SE data
type SEResponse struct {
	ID           string `json:"id"`
	IDPerusahaan string `json:"id_perusahaan"`
	IDSubSektor  string `json:"id_sub_sektor"`
	IDCsirt      string `json:"id_csirt"`

	// Karakteristik Instansi
	NilaiInvestasi             string `json:"nilai_investasi"`
	AnggaranOperasional        string `json:"anggaran_operasional"`
	KepatuhanPeraturan         string `json:"kepatuhan_peraturan"`
	TeknikKriptografi          string `json:"teknik_kriptografi"`
	JumlahPengguna             string `json:"jumlah_pengguna"`
	DataPribadi                string `json:"data_pribadi"`
	KlasifikasiData            string `json:"klasifikasi_data"`
	KekritisanProses           string `json:"kekritisan_proses"`
	DampakKegagalan            string `json:"dampak_kegagalan"`
	PotensiKerugiandanDampakNegatif string `json:"potensi_kerugian_dan_dampak_negatif"`

	// Informasi SE
	NamaSE      string `json:"nama_se"`
	IpSE        string `json:"ip_se"`
	AsNumberSE  string `json:"as_number_se"`
	PengelolaSE string `json:"pengelola_se"`
	FiturSE     string `json:"fitur_se"`

	// Hasil Kalkulasi
	TotalBobot int    `json:"total_bobot"`
	KategoriSE string `json:"kategori_se"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	// Nested objects
	Perusahaan *PerusahaanMiniResponse `json:"perusahaan,omitempty"`
	SubSektor  *SubSektorMiniResponse  `json:"sub_sektor,omitempty"`
	Csirt      *CsirtMiniResponse      `json:"csirt,omitempty"`
}

// SEListResponse represents the list response for SE
type SEListResponse struct {
	Data       []SEResponse `json:"data"`
	TotalCount int          `json:"total_count"`
}

// PerusahaanMiniResponse for nested response
type PerusahaanMiniResponse struct {
	ID             string `json:"id"`
	NamaPerusahaan string `json:"nama_perusahaan"`
}

// SubSektorMiniResponse for nested response
type SubSektorMiniResponse struct {
	ID            string `json:"id"`
	NamaSubSektor string `json:"nama_sub_sektor"`
	IDSektor      string `json:"id_sektor"`
	NamaSektor    string `json:"nama_sektor"`
}