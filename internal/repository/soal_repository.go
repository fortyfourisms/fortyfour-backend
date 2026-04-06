package repository

import (
	"database/sql"
	"fmt"

	"fortyfour-backend/internal/models"
)

type SoalRepository struct {
	db *sql.DB
}

func NewSoalRepository(db *sql.DB) *SoalRepository {
	return &SoalRepository{db: db}
}

var _ SoalRepositoryInterface = (*SoalRepository)(nil)

// Create menyimpan soal beserta semua pilihan jawabannya dalam satu transaksi.
func (r *SoalRepository) Create(soal *models.Soal, pilihan []models.PilihanJawaban) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`INSERT INTO soal (id, id_materi, pertanyaan, urutan, created_at)
		 VALUES (?, ?, ?, ?, NOW())`,
		soal.ID, soal.IDMateri, soal.Pertanyaan, soal.Urutan,
	)
	if err != nil {
		return err
	}

	for _, p := range pilihan {
		_, err = tx.Exec(
			`INSERT INTO pilihan_jawaban (id, id_soal, teks, is_correct, urutan)
			 VALUES (?, ?, ?, ?, ?)`,
			p.ID, p.IDSoal, p.Teks, p.IsCorrect, p.Urutan,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *SoalRepository) FindByID(id string) (*models.Soal, error) {
	row := r.db.QueryRow(
		`SELECT id, id_materi, pertanyaan, urutan, created_at FROM soal WHERE id = ?`, id,
	)
	soal, err := scanSoal(row)
	if err != nil {
		return nil, err
	}

	pilihan, err := r.findPilihanBySoal(soal.ID)
	if err != nil {
		return nil, err
	}
	soal.Pilihan = pilihan
	return soal, nil
}

func (r *SoalRepository) FindByMateri(idMateri string) ([]models.Soal, error) {
	rows, err := r.db.Query(
		`SELECT id, id_materi, pertanyaan, urutan, created_at
		 FROM soal WHERE id_materi = ? ORDER BY urutan ASC`, idMateri,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var soalList []models.Soal
	for rows.Next() {
		var s models.Soal
		if err := rows.Scan(&s.ID, &s.IDMateri, &s.Pertanyaan, &s.Urutan, &s.CreatedAt); err != nil {
			return nil, err
		}
		pilihan, err := r.findPilihanBySoal(s.ID)
		if err != nil {
			return nil, err
		}
		s.Pilihan = pilihan
		soalList = append(soalList, s)
	}
	return soalList, nil
}

// Update mengganti data soal dan — jika pilihan dikirim — menghapus pilihan lama
// lalu menyisipkan yang baru, semua dalam satu transaksi.
func (r *SoalRepository) Update(soal *models.Soal, pilihan []models.PilihanJawaban) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`UPDATE soal SET pertanyaan=?, urutan=? WHERE id=?`,
		soal.Pertanyaan, soal.Urutan, soal.ID,
	)
	if err != nil {
		return err
	}

	if len(pilihan) > 0 {
		if _, err = tx.Exec(`DELETE FROM pilihan_jawaban WHERE id_soal=?`, soal.ID); err != nil {
			return err
		}
		for _, p := range pilihan {
			_, err = tx.Exec(
				`INSERT INTO pilihan_jawaban (id, id_soal, teks, is_correct, urutan)
				 VALUES (?, ?, ?, ?, ?)`,
				p.ID, p.IDSoal, p.Teks, p.IsCorrect, p.Urutan,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *SoalRepository) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM soal WHERE id=?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("soal dengan id %s tidak ditemukan", id)
	}
	return nil
}

func (r *SoalRepository) FindPilihanByID(idPilihan string) (*models.PilihanJawaban, error) {
	row := r.db.QueryRow(
		`SELECT id, id_soal, teks, is_correct, urutan FROM pilihan_jawaban WHERE id=?`, idPilihan,
	)
	return scanPilihan(row)
}

func (r *SoalRepository) FindCorrectPilihan(idSoal string) (*models.PilihanJawaban, error) {
	row := r.db.QueryRow(
		`SELECT id, id_soal, teks, is_correct, urutan FROM pilihan_jawaban
		 WHERE id_soal=? AND is_correct=1 LIMIT 1`, idSoal,
	)
	return scanPilihan(row)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func (r *SoalRepository) findPilihanBySoal(idSoal string) ([]models.PilihanJawaban, error) {
	rows, err := r.db.Query(
		`SELECT id, id_soal, teks, is_correct, urutan FROM pilihan_jawaban
		 WHERE id_soal=? ORDER BY urutan ASC`, idSoal,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.PilihanJawaban
	for rows.Next() {
		var p models.PilihanJawaban
		if err := rows.Scan(&p.ID, &p.IDSoal, &p.Teks, &p.IsCorrect, &p.Urutan); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

func scanSoal(row *sql.Row) (*models.Soal, error) {
	var s models.Soal
	err := row.Scan(&s.ID, &s.IDMateri, &s.Pertanyaan, &s.Urutan, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func scanPilihan(row *sql.Row) (*models.PilihanJawaban, error) {
	var p models.PilihanJawaban
	err := row.Scan(&p.ID, &p.IDSoal, &p.Teks, &p.IsCorrect, &p.Urutan)
	if err != nil {
		return nil, err
	}
	return &p, nil
}