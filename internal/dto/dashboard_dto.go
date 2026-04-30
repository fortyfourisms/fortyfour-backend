package dto

// SectorCount represents counts per sektor.
type SectorCount struct {
	ID        string `json:"id"`
	Nama      string `json:"nama_sektor"`
	Total     int64  `json:"total"`
	ThisMonth int64  `json:"this_month"`
}

// IkasAgg summary for global ikas (keperluan summary).
type IkasAgg struct {
	Total              int64   `json:"total_ikas"`
	AvgNilaiKematangan float64 `json:"avg_nilai_kematangan"`
	AvgTargetNilai     float64 `json:"avg_target_nilai"`
}

// SeAgg summary for global se — termasuk breakdown per kategori dan this_month
type SeAgg struct {
	TotalSE   int64 `json:"total_se"`
	ThisMonth int64 `json:"this_month"`
	Strategis int64 `json:"strategis"`
	Tinggi    int64 `json:"tinggi"`
	Rendah    int64 `json:"rendah"`
}

// SeStatusCount menghitung perusahaan yang sudah/belum mengisi SE
type SeStatusCount struct {
	TotalPerusahaan int64 `json:"total_perusahaan"`
	SudahMengisiSE  int64 `json:"sudah_mengisi_se"`
	BelumMengisiSE  int64 `json:"belum_mengisi_se"`
}

// IkasStatusCount menghitung perusahaan yang sudah/belum mengisi IKAS.
type IkasStatusCount struct {
	TotalPerusahaan  int64 `json:"total_perusahaan"`
	SudahMengisiIKAS int64 `json:"sudah_mengisi_ikas"`
	BelumMengisiIKAS int64 `json:"belum_mengisi_ikas"`
}

// CsirtAgg summary statistik CSIRT
type CsirtAgg struct {
	TotalCSIRT int64 `json:"total_csirt"`
	ThisMonth  int64 `json:"this_month"`
}

// CsirtStatusCount menghitung perusahaan yang sudah/belum membentuk CSIRT
type CsirtStatusCount struct {
	TotalPerusahaan     int64 `json:"total_perusahaan"`
	SudahMembentukCSIRT int64 `json:"sudah_membentuk_csirt"`
	BelumMembentukCSIRT int64 `json:"belum_membentuk_csirt"`
}

// DashboardFilter menyimpan semua parameter filter yang diterima dari handler
type DashboardFilter struct {
	From        *string // YYYY-MM-DD
	To          *string // YYYY-MM-DD
	Year        *string // YYYY
	Quarter     *string // 1, 2, 3, 4
	SubSektorID *string
	KategoriSE  *string // Strategis | Tinggi | Rendah
}

// ── Response per section ─────────────────────────────────────────────────────

// DashboardSektorResponse — response untuk GET /api/dashboard/sektor
type DashboardSektorResponse struct {
	Sektor []SectorCount `json:"sektor_counts"`
}

// DashboardIkasResponse — response untuk GET /api/dashboard/ikas
type DashboardIkasResponse struct {
	Ikas       IkasAgg         `json:"ikas"`
	IkasStatus IkasStatusCount `json:"ikas_status"`
}

// DashboardSEResponse — response untuk GET /api/dashboard/se
type DashboardSEResponse struct {
	SE       SeAgg         `json:"se"`
	SEStatus SeStatusCount `json:"se_status"`
}

// DashboardCSIRTResponse — response untuk GET /api/dashboard/csirt
type DashboardCSIRTResponse struct {
	CSIRT       CsirtAgg         `json:"csirt"`
	CSIRTStatus CsirtStatusCount `json:"csirt_status"`
}
