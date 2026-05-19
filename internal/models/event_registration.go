package models

import "time"

type EventRegistration struct {
	ID         string    `json:"id"`
	EventID    string    `json:"event_id"`
	Nama       string    `json:"nama"`
	Email      string    `json:"email"`
	Perusahaan string    `json:"perusahaan"`
	Jabatan    string    `json:"jabatan"`
	NoHP       string    `json:"no_hp"`
	Sektor     string    `json:"sektor"`
	QRPayload  string    `json:"qr_payload"`
	QRToken    string    `json:"qr_token"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
