package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/utils"
	"strings"
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
		utils.ValueOrNull(req.Nama),
		utils.ValueOrNull(req.Telepon),
		utils.ValueOrNull(req.IDPerusahaan),
	)
	return err
}

func (r *PICRepository) GetAll() ([]dto.PICResponse, error) {
	rows, err := r.db.Query(`
        SELECT 
            p.id, 
            p.nama, 
            p.telepon, 
            p.created_at, 
            p.updated_at,
            per.id,
            per.nama_perusahaan
        FROM pic_perusahaan p
        LEFT JOIN perusahaan per ON p.id_perusahaan = per.id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []dto.PICResponse{}

	for rows.Next() {
		var p dto.PICResponse
		var perusahaanID sql.NullString
		var namaPerusahaan sql.NullString

		err := rows.Scan(
			&p.ID,
			&p.Nama,
			&p.Telepon,
			&p.CreatedAt,
			&p.UpdatedAt,
			&perusahaanID,
			&namaPerusahaan,
		)
		if err != nil {
			continue
		}

		// Jika perusahaan ada, tambahkan ke response
		if perusahaanID.Valid && namaPerusahaan.Valid {
			p.Perusahaan = &dto.PerusahaanInPIC{
				ID:             perusahaanID.String,
				NamaPerusahaan: namaPerusahaan.String,
			}
		}

		result = append(result, p)
	}

	return result, nil
}

func (r *PICRepository) GetByID(id string) (*dto.PICResponse, error) {
	row := r.db.QueryRow(`
        SELECT 
            p.id, 
            p.nama, 
            p.telepon, 
            p.created_at, 
            p.updated_at,
            per.id,
            per.nama_perusahaan
        FROM pic_perusahaan p
        LEFT JOIN perusahaan per ON p.id_perusahaan = per.id
        WHERE p.id = ?
    `, id)

	var p dto.PICResponse
	var perusahaanID sql.NullString
	var namaPerusahaan sql.NullString

	err := row.Scan(
		&p.ID,
		&p.Nama,
		&p.Telepon,
		&p.CreatedAt,
		&p.UpdatedAt,
		&perusahaanID,
		&namaPerusahaan,
	)

	if err != nil {
		return nil, err
	}

	// Jika perusahaan ada, tambahkan ke response
	if perusahaanID.Valid && namaPerusahaan.Valid {
		p.Perusahaan = &dto.PerusahaanInPIC{
			ID:             perusahaanID.String,
			NamaPerusahaan: namaPerusahaan.String,
		}
	}

	return &p, nil
}

func (r *PICRepository) Update(id string, req dto.UpdatePICRequest) error {
	query := "UPDATE pic_perusahaan SET "
	args := []interface{}{}
	updates := []string{}

	if req.Nama != nil {
		updates = append(updates, "nama=?")
		args = append(args, *req.Nama)
	}
	if req.Telepon != nil {
		updates = append(updates, "telepon=?")
		args = append(args, *req.Telepon)
	}
	if req.IDPerusahaan != nil {
		updates = append(updates, "id_perusahaan=?")
		args = append(args, *req.IDPerusahaan)
	}

	if len(updates) == 0 {
		// Tidak ada field yang diupdate, cek apakah data ada
		var count int
		err := r.db.QueryRow("SELECT COUNT(*) FROM pic_perusahaan WHERE id=?", id).Scan(&count)
		if err != nil {
			return err
		}
		if count == 0 {
			return sql.ErrNoRows
		}
		return nil
	}

	query += strings.Join(updates, ", ")
	query += " WHERE id=?"
	args = append(args, id)

	result, err := r.db.Exec(query, args...)
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

func (r *PICRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM pic_perusahaan WHERE id=?", id)
	return err
}

func (r *PICRepository) GetByPerusahaan(idPerusahaan string) ([]dto.PICResponse, error) {
	rows, err := r.db.Query(`
        SELECT 
            p.id, 
            p.nama, 
            p.telepon, 
            p.created_at, 
            p.updated_at,
            per.id,
            per.nama_perusahaan
        FROM pic_perusahaan p
        LEFT JOIN perusahaan per ON p.id_perusahaan = per.id
        WHERE p.id_perusahaan = ?
    `, idPerusahaan)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []dto.PICResponse{}
	for rows.Next() {
		var p dto.PICResponse
		var perusahaanID sql.NullString
		var namaPerusahaan sql.NullString

		err := rows.Scan(
			&p.ID,
			&p.Nama,
			&p.Telepon,
			&p.CreatedAt,
			&p.UpdatedAt,
			&perusahaanID,
			&namaPerusahaan,
		)
		if err != nil {
			continue
		}

		if perusahaanID.Valid && namaPerusahaan.Valid {
			p.Perusahaan = &dto.PerusahaanInPIC{
				ID:             perusahaanID.String,
				NamaPerusahaan: namaPerusahaan.String,
			}
		}

		result = append(result, p)
	}

	return result, nil
}
