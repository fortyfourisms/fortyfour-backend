package dto_event

import "time"

type SdmCsirtCreatedEvent struct {
	ID                string    `json:"id"`
	IdCsirt           string    `json:"id_csirt"`
	NamaPersonel      string    `json:"nama_personel"`
	JabatanCsirt      string    `json:"jabatan_csirt"`
	JabatanPerusahaan string    `json:"jabatan_perusahaan"`
	Skill             string    `json:"skill"`
	Sertifikasi       string    `json:"sertifikasi"`
	CreatedAt         time.Time `json:"created_at"`
}

type SdmCsirtUpdatedEvent struct {
	ID                string    `json:"id"`
	IdCsirt           string    `json:"id_csirt"`
	NamaPersonel      string    `json:"nama_personel"`
	JabatanCsirt      string    `json:"jabatan_csirt"`
	JabatanPerusahaan string    `json:"jabatan_perusahaan"`
	Skill             string    `json:"skill"`
	Sertifikasi       string    `json:"sertifikasi"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type SdmCsirtDeletedEvent struct {
	ID        string    `json:"id"`
	IdCsirt   string    `json:"id_csirt"`
	DeletedAt time.Time `json:"deleted_at"`
}
