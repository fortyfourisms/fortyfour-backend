package models

type RisikoSurvey struct {
	ID int `json:"id"`
	RespondenID int `json:"responden_id"`
	RisikoIP bool `json:"risiko_ip"`
	DampakReputasi string `json:"dampak_reputasi"`
	DampakOperasional string `json:"dampak_operasional"`
	DampakFinansial string `json:"dampak_finansial"`
	DampakHukum string `json:"dampak_hukum"`
	Frekuensi string `json:"frekuensi"`
	AdaPengendalian bool `json:"ada_pengendalian"`
	TindakanPengendalian string `json:"tindakan_pengendalian"`
}