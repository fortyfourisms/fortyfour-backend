package repository

import (
	"database/sql"
	"fmt"

	"fortyfour-backend/internal/models"
)

type KuisRepository struct {
	db *sql.DB
}

func NewKuisRepository(db *sql.DB) *KuisRepository {
	return &KuisRepository{db: db}
}

var _ KuisRepositoryInterface = (*KuisRepository)(nil)

func (r *KuisRepository) Create(k *models.Kuis) error {
	_, err := r.db.Exec(
		`INSERT INTO kuis (id, id_kelas, id_materi, judul, deskripsi, durasi_menit, passing_grade, is_final, urutan, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`,
		k.ID, k.IDKelas, k.IDMateri, k.Judul, k.Deskripsi, k.DurasiMenit, k.PassingGrade, k.IsFinal, k.Urutan,
	)
	return err
}

func (r *KuisRepository) FindByID(id string) (*models.Kuis, error) {
	row := r.db.QueryRow(
		`SELECT id, id_kelas, id_materi, judul, deskripsi, durasi_menit, passing_grade, is_final, urutan, created_at, updated_at
		 FROM kuis WHERE id = ?`, id,
	)
	return scanKuis(row)
}

func (r *KuisRepository) FindByKelas(idKelas string) ([]models.Kuis, error) {
	rows, err := r.db.Query(
		`SELECT id, id_kelas, id_materi, judul, deskripsi, durasi_menit, passing_grade, is_final, urutan, created_at, updated_at
		 FROM kuis WHERE id_kelas = ? ORDER BY is_final ASC, urutan ASC`, idKelas,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanKuisRows(rows)
}

func (r *KuisRepository) FindByMateri(idMateri string) (*models.Kuis, error) {
	row := r.db.QueryRow(
		`SELECT id, id_kelas, id_materi, judul, deskripsi, durasi_menit, passing_grade, is_final, urutan, created_at, updated_at
		 FROM kuis WHERE id_materi = ? LIMIT 1`, idMateri,
	)
	return scanKuis(row)
}

func (r *KuisRepository) FindFinalByKelas(idKelas string) (*models.Kuis, error) {
	row := r.db.QueryRow(
		`SELECT id, id_kelas, id_materi, judul, deskripsi, durasi_menit, passing_grade, is_final, urutan, created_at, updated_at
		 FROM kuis WHERE id_kelas = ? AND is_final = 1 LIMIT 1`, idKelas,
	)
	return scanKuis(row)
}

func (r *KuisRepository) Update(k *models.Kuis) error {
	_, err := r.db.Exec(
		`UPDATE kuis SET judul=?, deskripsi=?, durasi_menit=?, passing_grade=?, is_final=?, urutan=?, updated_at=NOW()
		 WHERE id=?`,
		k.Judul, k.Deskripsi, k.DurasiMenit, k.PassingGrade, k.IsFinal, k.Urutan, k.ID,
	)
	return err
}

func (r *KuisRepository) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM kuis WHERE id=?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("kuis dengan id %s tidak ditemukan", id)
	}
	return nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func scanKuis(row *sql.Row) (*models.Kuis, error) {
	var k models.Kuis
	var idMateri, deskripsi sql.NullString
	var durasiMenit sql.NullInt64
	err := row.Scan(
		&k.ID, &k.IDKelas, &idMateri, &k.Judul, &deskripsi, &durasiMenit,
		&k.PassingGrade, &k.IsFinal, &k.Urutan, &k.CreatedAt, &k.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if idMateri.Valid {
		k.IDMateri = &idMateri.String
	}
	if deskripsi.Valid {
		k.Deskripsi = &deskripsi.String
	}
	if durasiMenit.Valid {
		d := int(durasiMenit.Int64)
		k.DurasiMenit = &d
	}
	return &k, nil
}

func scanKuisRows(rows *sql.Rows) ([]models.Kuis, error) {
	var result []models.Kuis
	for rows.Next() {
		var k models.Kuis
		var idMateri, deskripsi sql.NullString
		var durasiMenit sql.NullInt64
		if err := rows.Scan(
			&k.ID, &k.IDKelas, &idMateri, &k.Judul, &deskripsi, &durasiMenit,
			&k.PassingGrade, &k.IsFinal, &k.Urutan, &k.CreatedAt, &k.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if idMateri.Valid {
			k.IDMateri = &idMateri.String
		}
		if deskripsi.Valid {
			k.Deskripsi = &deskripsi.String
		}
		if durasiMenit.Valid {
			d := int(durasiMenit.Int64)
			k.DurasiMenit = &d
		}
		result = append(result, k)
	}
	return result, nil
}
