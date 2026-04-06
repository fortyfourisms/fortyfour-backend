package repository

import (
	"database/sql"
	"fmt"
	"time"

	"fortyfour-backend/internal/models"
)

type KelasRepository struct {
	db *sql.DB
}

func NewKelasRepository(db *sql.DB) *KelasRepository {
	return &KelasRepository{db: db}
}

var _ KelasRepositoryInterface = (*KelasRepository)(nil)

func (r *KelasRepository) Create(k *models.Kelas) error {
	_, err := r.db.Exec(
		`INSERT INTO kelas (id, judul, deskripsi, thumbnail, status, created_by, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`,
		k.ID, k.Judul, k.Deskripsi, k.Thumbnail, k.Status, k.CreatedBy,
	)
	return err
}

func (r *KelasRepository) FindByID(id string) (*models.Kelas, error) {
	row := r.db.QueryRow(
		`SELECT id, judul, deskripsi, thumbnail, status, created_by, created_at, updated_at
		 FROM kelas WHERE id = ?`, id,
	)
	return scanKelas(row)
}

func (r *KelasRepository) FindAll(onlyPublished bool) ([]models.Kelas, error) {
	query := `SELECT id, judul, deskripsi, thumbnail, status, created_by, created_at, updated_at FROM kelas`
	if onlyPublished {
		query += ` WHERE status = 'published'`
	}
	query += ` ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Kelas
	for rows.Next() {
		var k models.Kelas
		var deskripsi, thumbnail sql.NullString
		if err := rows.Scan(
			&k.ID, &k.Judul, &deskripsi, &thumbnail,
			&k.Status, &k.CreatedBy, &k.CreatedAt, &k.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if deskripsi.Valid {
			k.Deskripsi = &deskripsi.String
		}
		if thumbnail.Valid {
			k.Thumbnail = &thumbnail.String
		}
		result = append(result, k)
	}
	return result, nil
}

func (r *KelasRepository) Update(k *models.Kelas) error {
	_, err := r.db.Exec(
		`UPDATE kelas SET judul=?, deskripsi=?, thumbnail=?, status=?, updated_at=NOW() WHERE id=?`,
		k.Judul, k.Deskripsi, k.Thumbnail, k.Status, k.ID,
	)
	return err
}

func (r *KelasRepository) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM kelas WHERE id=?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("kelas dengan id %s tidak ditemukan", id)
	}
	return nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func scanKelas(row *sql.Row) (*models.Kelas, error) {
	var k models.Kelas
	var deskripsi, thumbnail sql.NullString
	err := row.Scan(
		&k.ID, &k.Judul, &deskripsi, &thumbnail,
		&k.Status, &k.CreatedBy, &k.CreatedAt, &k.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if deskripsi.Valid {
		k.Deskripsi = &deskripsi.String
	}
	if thumbnail.Valid {
		k.Thumbnail = &thumbnail.String
	}
	return &k, nil
}

// memastikan time.Time tidak zero ketika dipakai
var _ = time.Now