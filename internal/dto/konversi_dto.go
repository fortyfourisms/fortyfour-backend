package dto

type KonversiResponse struct {
	PerusahaanID   string  `json:"perusahaan_id"`
	NamaPerusahaan string  `json:"nama_perusahaan"`
	PoinIkas       int     `json:"poin_ikas"`   // 1 jika ada, 0 jika tidak
	PoinKse        int     `json:"poin_kse"`    // 1 jika ada, 0 jika tidak
	PoinSurvey     int     `json:"poin_survey"` // 1 jika ada, 0 jika tidak
	PoinCsirt      int     `json:"poin_csirt"`  // 1 jika ada, 0 jika tidak
	TotalPoin      int     `json:"total_poin"`  // Sum dari semua poin (max 4)
	Persentase     float64 `json:"persentase"`  // (TotalPoin / 4) * 100
}
