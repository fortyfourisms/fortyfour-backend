package repository

import (
	"database/sql"
	"fortyfour-backend/internal/models"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

var _ EventRepositoryInterface = (*EventRepository)(nil)

func (r *EventRepository) Create(event *models.Event) error {
	query := `INSERT INTO events (id, slug, judul, deskripsi, lokasi, tanggal) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, event.ID, event.Slug, event.Judul, event.Deskripsi, event.Lokasi, event.Tanggal)
	return err
}

func (r *EventRepository) FindAll() ([]models.Event, error) {
	query := `
		SELECT id, slug, judul, deskripsi, lokasi, tanggal, created_at, updated_at
		FROM events
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		err := rows.Scan(
			&e.ID, &e.Slug, &e.Judul, &e.Deskripsi, &e.Lokasi, &e.Tanggal, &e.CreatedAt, &e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (r *EventRepository) FindByID(id string) (*models.Event, error) {
	query := `
		SELECT id, slug, judul, deskripsi, lokasi, tanggal, created_at, updated_at
		FROM events
		WHERE id = ?`

	var e models.Event
	err := r.db.QueryRow(query, id).Scan(
		&e.ID, &e.Slug, &e.Judul, &e.Deskripsi, &e.Lokasi, &e.Tanggal, &e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &e, nil
}

func (r *EventRepository) FindBySlug(slug string) (*models.Event, error) {
	query := `
		SELECT id, slug, judul, deskripsi, lokasi, tanggal, created_at, updated_at
		FROM events
		WHERE slug = ?`

	var e models.Event
	err := r.db.QueryRow(query, slug).Scan(
		&e.ID, &e.Slug, &e.Judul, &e.Deskripsi, &e.Lokasi, &e.Tanggal, &e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &e, nil
}

func (r *EventRepository) Update(event *models.Event) error {
	query := `UPDATE events SET slug = ?, judul = ?, deskripsi = ?, lokasi = ?, tanggal = ? WHERE id = ?`
	_, err := r.db.Exec(query, event.Slug, event.Judul, event.Deskripsi, event.Lokasi, event.Tanggal, event.ID)
	return err
}

func (r *EventRepository) Delete(id string) error {
	query := `DELETE FROM events WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *EventRepository) CreateRegistration(reg *models.EventRegistration) error {
	query := `
		INSERT INTO event_registrations
			(event_id, nama, email, perusahaan, jabatan, no_hp, sektor, qr_payload, qr_token)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := r.db.Exec(query,
		reg.EventID,
		reg.Nama,
		reg.Email,
		reg.Perusahaan,
		reg.Jabatan,
		reg.NoHP,
		reg.Sektor,
		reg.QRPayload,
		reg.QRToken,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	reg.ID = id

	saved, err := r.FindRegistrationByID(id)
	if err != nil {
		return err
	}
	if saved != nil {
		*reg = *saved
	}

	return nil
}

func (r *EventRepository) FindRegistrationByID(id int64) (*models.EventRegistration, error) {
	query := `
		SELECT id, event_id, nama, email, perusahaan, jabatan, no_hp, sektor, qr_payload, qr_token, created_at, updated_at
		FROM event_registrations
		WHERE id = ?
	`

	var reg models.EventRegistration
	var qrPayload []byte

	err := r.db.QueryRow(query, id).Scan(
		&reg.ID,
		&reg.EventID,
		&reg.Nama,
		&reg.Email,
		&reg.Perusahaan,
		&reg.Jabatan,
		&reg.NoHP,
		&reg.Sektor,
		&qrPayload,
		&reg.QRToken,
		&reg.CreatedAt,
		&reg.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	reg.QRPayload = string(qrPayload)
	return &reg, nil
}

func (r *EventRepository) ExistsRegistrationByEventAndEmail(eventID string, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM event_registrations WHERE event_id = ? AND email = ?)`

	var exists bool
	if err := r.db.QueryRow(query, eventID, email).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (r *EventRepository) UpdateRegistrationPayload(id int64, payload string) error {
	query := `UPDATE event_registrations SET qr_payload = ? WHERE id = ?`
	_, err := r.db.Exec(query, payload, id)
	return err
}
