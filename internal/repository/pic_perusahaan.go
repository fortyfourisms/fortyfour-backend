package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/utils"
)

type PICRepository struct {
	db *sql.DB
}

func NewPICRepository(db *sql.DB) *PICRepository {
	return &PICRepository{db: db}
}

func (r *PICRepository) Create(req dto.CreatePICRequest, id string) error {
	_, err := r.db.Exec(`
        INSERT INTO pic_perusahaan
        (id, nama, telepon, id_perusahaan)
        VALUES (?, ?, ?, ?)
    `,
		id,
		utils.ValueOrEmpty(req.Nama),
		utils.ValueOrEmpty(req.Telepon),
		utils.ValueOrEmpty(req.IDPerusahaan),
	)
	return err
}

func (r *PICRepository) GetAll() ([]dto.PICResponse, error) {
	rows, err := r.db.Query(`
        SELECT id, nama, telepon, id_perusahaan, created_at, updated_at
        FROM pic_perusahaan
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []dto.PICResponse{}

	for rows.Next() {
		var p dto.PICResponse
		rows.Scan(
			&p.ID,
			&p.Nama,
			&p.Telepon,
			&p.IDPerusahaan,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		result = append(result, p)
	}

	return result, nil
}

func (r *PICRepository) GetByID(id string) (*dto.PICResponse, error) {
	row := r.db.QueryRow(`
        SELECT id, nama, telepon, id_perusahaan, created_at, updated_at
        FROM pic_perusahaan WHERE id=?
    `, id)

	var p dto.PICResponse
	err := row.Scan(
		&p.ID,
		&p.Nama,
		&p.Telepon,
		&p.IDPerusahaan,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *PICRepository) Update(id string, req dto.PICResponse) error {
	_, err := r.db.Exec(`
        UPDATE pic_perusahaan SET
            nama=?, telepon=?, id_perusahaan=?
        WHERE id=?
    `,
		req.Nama,
		req.Telepon,
		req.IDPerusahaan,
		id,
	)
	return err
}

func (r *PICRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM pic_perusahaan WHERE id=?", id)
	return err
}
