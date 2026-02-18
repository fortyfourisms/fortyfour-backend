package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
)

type SubSektorRepository struct {
	db *sql.DB
}

func NewSubSektorRepository(db *sql.DB) *SubSektorRepository {
	return &SubSektorRepository{db: db}
}

func (r *SubSektorRepository) GetAll() ([]dto.SubSektorResponse, error) {
	rows, err := r.db.Query(`
		SELECT ss.id, ss.nama_sub_sektor, ss.id_sektor, s.nama_sektor, ss.created_at, ss.updated_at
		FROM sub_sektor ss
		JOIN sektor s ON ss.id_sektor = s.id
		ORDER BY s.nama_sektor, ss.nama_sub_sektor
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.SubSektorResponse
	for rows.Next() {
		var sub dto.SubSektorResponse
		rows.Scan(&sub.ID, &sub.NamaSubSektor, &sub.IDSektor, &sub.NamaSektor, &sub.CreatedAt, &sub.UpdatedAt)
		result = append(result, sub)
	}
	return result, nil
}

func (r *SubSektorRepository) GetByID(id string) (*dto.SubSektorResponse, error) {
	row := r.db.QueryRow(`
		SELECT ss.id, ss.nama_sub_sektor, ss.id_sektor, s.nama_sektor, ss.created_at, ss.updated_at
		FROM sub_sektor ss
		JOIN sektor s ON ss.id_sektor = s.id
		WHERE ss.id=?
	`, id)

	var sub dto.SubSektorResponse
	err := row.Scan(&sub.ID, &sub.NamaSubSektor, &sub.IDSektor, &sub.NamaSektor, &sub.CreatedAt, &sub.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *SubSektorRepository) GetBySektorID(sektorID string) ([]dto.SubSektorResponse, error) {
	rows, err := r.db.Query(`
		SELECT ss.id, ss.nama_sub_sektor, ss.id_sektor, s.nama_sektor, ss.created_at, ss.updated_at
		FROM sub_sektor ss
		JOIN sektor s ON ss.id_sektor = s.id
		WHERE ss.id_sektor=?
		ORDER BY ss.nama_sub_sektor
	`, sektorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.SubSektorResponse
	for rows.Next() {
		var sub dto.SubSektorResponse
		rows.Scan(&sub.ID, &sub.NamaSubSektor, &sub.IDSektor, &sub.NamaSektor, &sub.CreatedAt, &sub.UpdatedAt)
		result = append(result, sub)
	}
	return result, nil
}
