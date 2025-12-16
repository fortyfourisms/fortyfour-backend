package models

type SdmCsirt struct {
	ID                string `json:"id"`
	IdCsirt           string `json:"id_csirt"`
	NamaPersonel      string `json:"nama_personel"`
	JabatanCsirt      string `json:"jabatan_csirt"`
	JabatanPerusahaan string `json:"jabatan_perusahaan"`
	Skill             string `json:"skill"`
	Sertifikasi       string `json:"sertifikasi"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}
