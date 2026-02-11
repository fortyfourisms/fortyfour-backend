package dto_event

import "time"

// IkasCreatedEvent
type IkasCreatedEvent struct {
	IkasID          string    `json:"ikas_id"`
	IDPerusahaan    string    `json:"id_perusahaan"`
	NamaPerusahaan  string    `json:"nama_perusahaan,omitempty"`
	Tanggal         string    `json:"tanggal"`
	Responden       string    `json:"responden"`
	NilaiKematangan float64   `json:"nilai_kematangan"`
	TargetNilai     float64   `json:"target_nilai"`
	CreatedAt       time.Time `json:"created_at"`
}

// IkasUpdatedEvent
type IkasUpdatedEvent struct {
	IkasID             string    `json:"ikas_id"`
	OldNilaiKematangan float64   `json:"old_nilai_kematangan"`
	NewNilaiKematangan float64   `json:"new_nilai_kematangan"`
	UpdatedFields      []string  `json:"updated_fields"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// IkasDeletedEvent
type IkasDeletedEvent struct {
	IkasID    string    `json:"ikas_id"`
	DeletedAt time.Time `json:"deleted_at"`
}

// IkasImportedEvent
type IkasImportedEvent struct {
	IkasID          string    `json:"ikas_id"`
	IDPerusahaan    string    `json:"id_perusahaan"`
	NamaPerusahaan  string    `json:"nama_perusahaan,omitempty"`
	FileName        string    `json:"file_name,omitempty"`
	NilaiKematangan float64   `json:"nilai_kematangan"`
	ImportedAt      time.Time `json:"imported_at"`
}

// EmailNotificationPayload
type EmailNotificationPayload struct {
	To      string                 `json:"to"`
	Subject string                 `json:"subject"`
	Body    string                 `json:"body"`
	Data    map[string]interface{} `json:"data,omitempty"`
}
