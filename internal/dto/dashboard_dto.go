package dto

// SectorCount represents counts per sektor.
type SectorCount struct {
	ID        string `json:"id"`
	Nama      string `json:"nama_sektor"`
	Total     int64  `json:"total"`
	ThisMonth int64  `json:"this_month"`
}

// TODO: re-enable ikas summary when ikas table is ready
// IkasAgg summary for global ikas (keperluan summary)
// type IkasAgg struct {
// 	Total                int64   `json:"total_ikas"`
// 	AvgNilaiKematangan   float64 `json:"avg_nilai_kematangan"`
// 	AvgTargetNilai       float64 `json:"avg_target_nilai"`
// }

// SeAgg summary for global se — termasuk breakdown per kategori dan this_month
type SeAgg struct {
	TotalSE   int64 `json:"total_se"`
	ThisMonth int64 `json:"this_month"`
	Strategis int64 `json:"strategis"`
	Tinggi    int64 `json:"tinggi"`
	Rendah    int64 `json:"rendah"`
}

// SeStatusCount menghitung perusahaan yang sudah/belum mengisi KSE
type SeStatusCount struct {
	TotalPerusahaan int64 `json:"total_perusahaan"`
	SudahMengisiKSE int64 `json:"sudah_mengisi_kse"`
	BelumMengisiKSE int64 `json:"belum_mengisi_kse"`
}

// TODO: re-enable ikas status when ikas table is ready
// IkasStatusCount menghitung perusahaan yang sudah/belum mengisi IKAS
// type IkasStatusCount struct {
// 	TotalPerusahaan   int64 `json:"total_perusahaan"`
// 	SudahMengisiIKAS  int64 `json:"sudah_mengisi_ikas"`
// 	BelumMengisiIKAS  int64 `json:"belum_mengisi_ikas"`
// }

// DashboardFilter menyimpan semua parameter filter yang diterima dari handler
type DashboardFilter struct {
	From        *string // YYYY-MM-DD
	To          *string // YYYY-MM-DD
	Year        *string // YYYY
	Quarter     *string // 1, 2, 3, 4
	SubSektorID *string
	KategoriSE  *string // Strategis | Tinggi | Rendah
}

// DashboardSummary top-level
type DashboardSummary struct {
	Sektor []SectorCount `json:"sektor_counts"`
	// Ikas   IkasAgg       `json:"ikas"` // TODO: re-enable ikas summary when ikas table is ready
	SE       SeAgg         `json:"kse"`
	SEStatus SeStatusCount `json:"kse_status"`
	// IkasStatus IkasStatusCount `json:"ikas_status"` // TODO: re-enable ikas status when ikas table is ready
}
