package repository

import (
	"database/sql"
	"fmt"
	"fortyfour-backend/internal/dto"
)

type JabatanRepository struct {
	db *sql.DB
}

func NewJabatanRepository(db *sql.DB) *JabatanRepository {
	return &JabatanRepository{db: db}
}

func (r *JabatanRepository) Create(req dto.CreateJabatanRequest, id string) error {
	_, err := r.db.Exec(`INSERT INTO jabatan (id, nama_jabatan) VALUES (?, ?)`,
		id, req.NamaJabatan,
	)
	return err
}

func (r *JabatanRepository) GetAll() ([]dto.JabatanResponse, error) {
	rows, err := r.db.Query(`SELECT id, nama_jabatan, created_at, updated_at FROM jabatan`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.JabatanResponse
	for rows.Next() {
		var j dto.JabatanResponse
		rows.Scan(&j.ID, &j.NamaJabatan, &j.CreatedAt, &j.UpdatedAt)
		result = append(result, j)
	}
	return result, nil
}

func (r *JabatanRepository) GetByID(id string) (*dto.JabatanResponse, error) {
	row := r.db.QueryRow(`SELECT id, nama_jabatan, created_at, updated_at FROM jabatan WHERE id=?`, id)
	var j dto.JabatanResponse
	err := row.Scan(&j.ID, &j.NamaJabatan, &j.CreatedAt, &j.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *JabatanRepository) Update(id string, req dto.JabatanResponse) error {
	_, err := r.db.Exec(`UPDATE jabatan SET nama_jabatan=? WHERE id=?`,
		req.NamaJabatan, id,
	)
	return err
}

func (r *JabatanRepository) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM jabatan WHERE id=?`, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("data dengan id %s tidak ditemukan", id)
	}

	return nil
}
