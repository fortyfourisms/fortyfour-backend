package dto

type CreateSdmCsirtRequest struct {
	IdCsirt           *string `json:"id_csirt,omitempty"`
	NamaPersonel      *string `json:"nama_personel,omitempty"`
	JabatanCsirt      *string `json:"jabatan_csirt,omitempty"`
	JabatanPerusahaan *string `json:"jabatan_perusahaan,omitempty"`
	Skill             *string `json:"skill,omitempty"`
	Sertifikasi       *string `json:"sertifikasi,omitempty"`
}

type UpdateSdmCsirtRequest struct {
	NamaPersonel      *string `json:"nama_personel,omitempty"`
	JabatanCsirt      *string `json:"jabatan_csirt,omitempty"`
	JabatanPerusahaan *string `json:"jabatan_perusahaan,omitempty"`
	Skill             *string `json:"skill,omitempty"`
	Sertifikasi       *string `json:"sertifikasi,omitempty"`
}

type SdmCsirtResponse struct {
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
