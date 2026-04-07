package repository

import (
	"database/sql"
	"fmt"

	"fortyfour-backend/internal/models"
)

type MateriRepository struct {
	db *sql.DB
}

func NewMateriRepository(db *sql.DB) *MateriRepository {
	return &MateriRepository{db: db}
}

var _ MateriRepositoryInterface = (*MateriRepository)(nil)

func (r *MateriRepository) Create(m *models.Materi) error {
	_, err := r.db.Exec(
		`INSERT INTO materi (id, id_kelas, judul, tipe, urutan, youtube_id, pdf_path, durasi_detik, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`,
		m.ID, m.IDKelas, m.Judul, m.Tipe, m.Urutan, m.YoutubeID, m.PDFPath, m.DurasiDetik,
	)
	return err
}

func (r *MateriRepository) FindByID(id string) (*models.Materi, error) {
	row := r.db.QueryRow(
		`SELECT id, id_kelas, judul, tipe, urutan, youtube_id, pdf_path, durasi_detik, created_at, updated_at
		 FROM materi WHERE id = ?`, id,
	)
	return scanMateri(row)
}

func (r *MateriRepository) FindByKelas(idKelas string) ([]models.Materi, error) {
	rows, err := r.db.Query(
		`SELECT id, id_kelas, judul, tipe, urutan, youtube_id, pdf_path, durasi_detik, created_at, updated_at
		 FROM materi WHERE id_kelas = ? ORDER BY urutan ASC`, idKelas,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMateriRows(rows)
}

func (r *MateriRepository) FindByKelasBeforeUrutan(idKelas string, urutan int) ([]models.Materi, error) {
	rows, err := r.db.Query(
		`SELECT id, id_kelas, judul, tipe, urutan, youtube_id, pdf_path, durasi_detik, created_at, updated_at
		 FROM materi WHERE id_kelas = ? AND urutan < ? ORDER BY urutan ASC`,
		idKelas, urutan,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMateriRows(rows)
}

func (r *MateriRepository) Update(m *models.Materi) error {
	_, err := r.db.Exec(
		`UPDATE materi SET judul=?, urutan=?, youtube_id=?, pdf_path=?, durasi_detik=?, updated_at=NOW()
		 WHERE id=?`,
		m.Judul, m.Urutan, m.YoutubeID, m.PDFPath, m.DurasiDetik, m.ID,
	)
	return err
}

func (r *MateriRepository) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM materi WHERE id=?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("materi dengan id %s tidak ditemukan", id)
	}
	return nil
}

// ReorderUrutan memperbarui urutan materi dalam satu kelas agar tetap berurutan 1,2,3,...
// dipanggil setelah delete supaya tidak ada gap.
func (r *MateriRepository) ReorderUrutan(idKelas string) error {
	rows, err := r.db.Query(
		`SELECT id FROM materi WHERE id_kelas = ? ORDER BY urutan ASC`, idKelas,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return err
		}
		ids = append(ids, id)
	}

	for i, id := range ids {
		if _, err := r.db.Exec(`UPDATE materi SET urutan=? WHERE id=?`, i+1, id); err != nil {
			return err
		}
	}
	return nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func scanMateri(row *sql.Row) (*models.Materi, error) {
	var m models.Materi
	var youtubeID, pdfPath sql.NullString
	var durasiDetik sql.NullInt64
	err := row.Scan(
		&m.ID, &m.IDKelas, &m.Judul, &m.Tipe, &m.Urutan,
		&youtubeID, &pdfPath, &durasiDetik,
		&m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if youtubeID.Valid {
		m.YoutubeID = &youtubeID.String
	}
	if pdfPath.Valid {
		m.PDFPath = &pdfPath.String
	}
	if durasiDetik.Valid {
		d := int(durasiDetik.Int64)
		m.DurasiDetik = &d
	}
	return &m, nil
}

func scanMateriRows(rows *sql.Rows) ([]models.Materi, error) {
	var result []models.Materi
	for rows.Next() {
		var m models.Materi
		var youtubeID, pdfPath sql.NullString
		var durasiDetik sql.NullInt64
		if err := rows.Scan(
			&m.ID, &m.IDKelas, &m.Judul, &m.Tipe, &m.Urutan,
			&youtubeID, &pdfPath, &durasiDetik,
			&m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if youtubeID.Valid {
			m.YoutubeID = &youtubeID.String
		}
		if pdfPath.Valid {
			m.PDFPath = &pdfPath.String
		}
		if durasiDetik.Valid {
			d := int(durasiDetik.Int64)
			m.DurasiDetik = &d
		}
		result = append(result, m)
	}
	return result, nil
}
