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

// SeAgg summary for global se
type SeAgg struct {
	TotalSE int64 `json:"total_se"`
}

// DashboardSummary top-level
type DashboardSummary struct {
	Sektor []SectorCount `json:"sektor_counts"`
	// Ikas   IkasAgg       `json:"ikas"` // TODO: re-enable ikas summary when ikas table is ready
	SE     SeAgg         `json:"kse"`
}
