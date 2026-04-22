package repository

import (
	"database/sql"
	"fmt"

	"fortyfour-backend/internal/models"
)

type DiskusiRepository struct {
	db *sql.DB
}

func NewDiskusiRepository(db *sql.DB) *DiskusiRepository {
	return &DiskusiRepository{db: db}
}

var _ DiskusiRepositoryInterface = (*DiskusiRepository)(nil)

func (r *DiskusiRepository) Create(d *models.Diskusi) error {
	_, err := r.db.Exec(
		`INSERT INTO diskusi (id, id_materi, id_user, id_parent, konten, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, NOW(), NOW())`,
		d.ID, d.IDMateri, d.IDUser, d.IDParent, d.Konten,
	)
	return err
}

func (r *DiskusiRepository) FindByMateri(idMateri string) ([]models.Diskusi, error) {
	rows, err := r.db.Query(
		`SELECT id, id_materi, id_user, id_parent, konten, created_at, updated_at
		 FROM diskusi WHERE id_materi = ? AND id_parent IS NULL ORDER BY created_at ASC`, idMateri,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDiskusiRows(rows)
}

func (r *DiskusiRepository) FindByID(id string) (*models.Diskusi, error) {
	row := r.db.QueryRow(
		`SELECT id, id_materi, id_user, id_parent, konten, created_at, updated_at
		 FROM diskusi WHERE id = ?`, id,
	)
	return scanDiskusi(row)
}

func (r *DiskusiRepository) Update(d *models.Diskusi) error {
	_, err := r.db.Exec(
		`UPDATE diskusi SET konten=?, updated_at=NOW() WHERE id=?`,
		d.Konten, d.ID,
	)
	return err
}

func (r *DiskusiRepository) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM diskusi WHERE id=?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("diskusi dengan id %s tidak ditemukan", id)
	}
	return nil
}

func (r *DiskusiRepository) FindReplies(idParent string) ([]models.Diskusi, error) {
	rows, err := r.db.Query(
		`SELECT id, id_materi, id_user, id_parent, konten, created_at, updated_at
		 FROM diskusi WHERE id_parent = ? ORDER BY created_at ASC`, idParent,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDiskusiRows(rows)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func scanDiskusi(row *sql.Row) (*models.Diskusi, error) {
	var d models.Diskusi
	var idParent sql.NullString
	err := row.Scan(&d.ID, &d.IDMateri, &d.IDUser, &idParent, &d.Konten, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if idParent.Valid {
		d.IDParent = &idParent.String
	}
	return &d, nil
}

func scanDiskusiRows(rows *sql.Rows) ([]models.Diskusi, error) {
	var result []models.Diskusi
	for rows.Next() {
		var d models.Diskusi
		var idParent sql.NullString
		if err := rows.Scan(&d.ID, &d.IDMateri, &d.IDUser, &idParent, &d.Konten, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		if idParent.Valid {
			d.IDParent = &idParent.String
		}
		result = append(result, d)
	}
	return result, nil
}
