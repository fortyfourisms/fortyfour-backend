package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
	"time"

	"github.com/google/uuid"
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

func (r *SubSektorRepository) Create(req dto.SubSektorRequest) (*dto.SubSektorResponse, error) {
	id := uuid.New().String()
	now := time.Now()

	_, err := r.db.Exec(`INSERT INTO sub_sektor (id, nama_sub_sektor, id_sektor, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		id, req.NamaSubSektor, req.IDSektor, now, now)
	if err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

func (r *SubSektorRepository) Update(id string, req dto.SubSektorRequest) (*dto.SubSektorResponse, error) {
	now := time.Now()
	result, err := r.db.Exec(`UPDATE sub_sektor SET nama_sub_sektor = ?, id_sektor = ?, updated_at = ? WHERE id = ?`,
		req.NamaSubSektor, req.IDSektor, now, id)
	if err != nil {
		return nil, err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}
	return r.GetByID(id)
}

func (r *SubSektorRepository) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM sub_sektor WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
