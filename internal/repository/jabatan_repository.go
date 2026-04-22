package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
)

type JabatanRepository struct {
	db *sql.DB
}

func NewJabatanRepository(db *sql.DB) *JabatanRepository {
	return &JabatanRepository{db: db}
}

func (r *JabatanRepository) Create(req dto.CreateJabatanRequest, id string) error {
	query := `INSERT INTO jabatan (id, nama_jabatan) VALUES (?, ?)`
	_, err := r.db.Exec(query, id, *req.NamaJabatan)
	return err
}

func (r *JabatanRepository) GetByID(id string) (*dto.JabatanResponse, error) {
	row := r.db.QueryRow(`SELECT id, nama_jabatan, created_at, updated_at FROM jabatan WHERE id=?`, id)

	var jabatan dto.JabatanResponse
	err := row.Scan(&jabatan.ID, &jabatan.NamaJabatan, &jabatan.CreatedAt, &jabatan.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &jabatan, nil
}

func (r *JabatanRepository) GetAll() ([]dto.JabatanResponse, error) {
	rows, err := r.db.Query(`SELECT id, nama_jabatan, created_at, updated_at FROM jabatan ORDER BY nama_jabatan`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.JabatanResponse
	for rows.Next() {
		var j dto.JabatanResponse
		if err := rows.Scan(&j.ID, &j.NamaJabatan, &j.CreatedAt, &j.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, j)
	}
	return result, rows.Err()
}

func (r *JabatanRepository) Update(id string, jabatan dto.JabatanResponse) error {
	query := `UPDATE jabatan SET nama_jabatan=? WHERE id=?`
	result, err := r.db.Exec(query, jabatan.NamaJabatan, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *JabatanRepository) Delete(id string) error {
	query := `DELETE FROM jabatan WHERE id=?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
