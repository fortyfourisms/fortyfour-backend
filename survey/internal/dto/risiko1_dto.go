package dto

type RisikoResponse struct {
	ID         int    `json:"id"`
	NamaRisiko string `json:"nama_risiko"`
	Deskripsi  string `json:"deskripsi"`
}

// INPUT JAWABAN RISIKO
type CreateRisikoJawabanRequest struct {
	RespondenID int `json:"responden_id"`
	RisikoID    int `json:"risiko_id"`

	PernahTerjadi string `json:"pernah_terjadi"`

	DampakReputasi    string `json:"dampak_reputasi"`
	DampakOperasional string `json:"dampak_operasional"`
	DampakFinansial   string `json:"dampak_finansial"`
	DampakHukum       string `json:"dampak_hukum"`

	Frekuensi string `json:"frekuensi"`

	AdaPengendalian       string `json:"ada_pengendalian"`
	DeskripsiPengendalian string `json:"deskripsi_pengendalian"`
}
