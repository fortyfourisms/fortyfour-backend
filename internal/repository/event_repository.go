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
	query := `INSERT INTO events (judul, deskripsi, tanggal) VALUES (?, ?, ?)`
	res, err := r.db.Exec(query, event.Judul, event.Deskripsi, event.Tanggal)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	event.ID = id
	return nil
}

func (r *EventRepository) FindAll() ([]models.Event, error) {
	query := `
		SELECT id, judul, deskripsi, tanggal, created_at, updated_at
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
			&e.ID, &e.Judul, &e.Deskripsi, &e.Tanggal, &e.CreatedAt, &e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (r *EventRepository) FindByID(id int64) (*models.Event, error) {
	query := `
		SELECT id, judul, deskripsi, tanggal, created_at, updated_at
		FROM events
		WHERE id = ?`

	var e models.Event
	err := r.db.QueryRow(query, id).Scan(
		&e.ID, &e.Judul, &e.Deskripsi, &e.Tanggal, &e.CreatedAt, &e.UpdatedAt,
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
	query := `UPDATE events SET judul = ?, deskripsi = ?, tanggal = ? WHERE id = ?`
	_, err := r.db.Exec(query, event.Judul, event.Deskripsi, event.Tanggal, event.ID)
	return err
}

func (r *EventRepository) Delete(id int64) error {
	query := `DELETE FROM events WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
