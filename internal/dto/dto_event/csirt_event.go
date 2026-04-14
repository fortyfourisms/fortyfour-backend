package dto_event

import "time"

// CsirtCreatedEvent payload untuk event csirt created
type CsirtCreatedEvent struct {
	ID                string    `json:"id"`
	IdPerusahaan      string    `json:"id_perusahaan"`
	NamaCsirt         string    `json:"nama_csirt"`
	WebCsirt          string    `json:"web_csirt"`
	TanggalRegistrasi string    `json:"tanggal_registrasi"`
	CreatedAt         time.Time `json:"created_at"`
}

// CsirtUpdatedEvent payload untuk event csirt updated
type CsirtUpdatedEvent struct {
	ID                string    `json:"id"`
	IdPerusahaan      string    `json:"id_perusahaan"`
	NamaCsirt         string    `json:"nama_csirt"`
	WebCsirt          string    `json:"web_csirt"`
	TanggalRegistrasi string    `json:"tanggal_registrasi"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// CsirtDeletedEvent payload untuk event csirt deleted
type CsirtDeletedEvent struct {
	ID        string    `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}
