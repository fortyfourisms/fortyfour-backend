package repository

import (
	"database/sql"
	"fmt"

	"fortyfour-backend/internal/models"
)

type KuisAttemptRepository struct {
	db *sql.DB
}

func NewKuisAttemptRepository(db *sql.DB) *KuisAttemptRepository {
	return &KuisAttemptRepository{db: db}
}

var _ KuisAttemptRepositoryInterface = (*KuisAttemptRepository)(nil)

func (r *KuisAttemptRepository) Create(a *models.KuisAttempt) error {
	_, err := r.db.Exec(
		`INSERT INTO kuis_attempt (id, id_user, id_materi, skor, total_soal, total_benar, started_at)
		 VALUES (?, ?, ?, 0, 0, 0, NOW())`,
		a.ID, a.IDUser, a.IDMateri,
	)
	return err
}

func (r *KuisAttemptRepository) FindByID(id string) (*models.KuisAttempt, error) {
	row := r.db.QueryRow(
		`SELECT id, id_user, id_materi, skor, total_soal, total_benar, started_at, finished_at
		 FROM kuis_attempt WHERE id=?`, id,
	)
	return scanAttempt(row)
}

func (r *KuisAttemptRepository) FindByUserAndMateri(idUser, idMateri string) ([]models.KuisAttempt, error) {
	rows, err := r.db.Query(
		`SELECT id, id_user, id_materi, skor, total_soal, total_benar, started_at, finished_at
		 FROM kuis_attempt WHERE id_user=? AND id_materi=? ORDER BY started_at DESC`,
		idUser, idMateri,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.KuisAttempt
	for rows.Next() {
		var a models.KuisAttempt
		var finishedAt sql.NullTime
		if err := rows.Scan(
			&a.ID, &a.IDUser, &a.IDMateri, &a.Skor,
			&a.TotalSoal, &a.TotalBenar, &a.StartedAt, &finishedAt,
		); err != nil {
			return nil, err
		}
		if finishedAt.Valid {
			a.FinishedAt = &finishedAt.Time
		}
		result = append(result, a)
	}
	return result, nil
}

// FindLatestByUserAndMateri mengembalikan attempt terakhir user untuk materi kuis tertentu.
// Dipakai untuk cek apakah ada attempt yang belum selesai (finished_at IS NULL).
func (r *KuisAttemptRepository) FindLatestByUserAndMateri(idUser, idMateri string) (*models.KuisAttempt, error) {
	row := r.db.QueryRow(
		`SELECT id, id_user, id_materi, skor, total_soal, total_benar, started_at, finished_at
		 FROM kuis_attempt WHERE id_user=? AND id_materi=?
		 ORDER BY started_at DESC LIMIT 1`,
		idUser, idMateri,
	)
	return scanAttempt(row)
}

// Finish menyimpan hasil kuis: update attempt + insert semua jawaban, dalam satu transaksi.
func (r *KuisAttemptRepository) Finish(id string, skor float64, totalBenar int, jawaban []models.KuisJawaban) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	totalSoal := len(jawaban)
	res, err := tx.Exec(
		`UPDATE kuis_attempt
		 SET skor=?, total_soal=?, total_benar=?, finished_at=NOW()
		 WHERE id=? AND finished_at IS NULL`,
		skor, totalSoal, totalBenar, id,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("attempt tidak ditemukan atau sudah selesai")
	}

	for _, j := range jawaban {
		_, err = tx.Exec(
			`INSERT INTO kuis_jawaban (id, id_attempt, id_soal, id_pilihan, is_correct)
			 VALUES (?, ?, ?, ?, ?)`,
			j.ID, j.IDAttempt, j.IDSoal, j.IDPilihan, j.IsCorrect,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *KuisAttemptRepository) FindJawabanByAttempt(idAttempt string) ([]models.KuisJawaban, error) {
	rows, err := r.db.Query(
		`SELECT id, id_attempt, id_soal, id_pilihan, is_correct
		 FROM kuis_jawaban WHERE id_attempt=?`, idAttempt,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.KuisJawaban
	for rows.Next() {
		var j models.KuisJawaban
		if err := rows.Scan(&j.ID, &j.IDAttempt, &j.IDSoal, &j.IDPilihan, &j.IsCorrect); err != nil {
			return nil, err
		}
		result = append(result, j)
	}
	return result, nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func scanAttempt(row *sql.Row) (*models.KuisAttempt, error) {
	var a models.KuisAttempt
	var finishedAt sql.NullTime
	err := row.Scan(
		&a.ID, &a.IDUser, &a.IDMateri, &a.Skor,
		&a.TotalSoal, &a.TotalBenar, &a.StartedAt, &finishedAt,
	)
	if err != nil {
		return nil, err
	}
	if finishedAt.Valid {
		a.FinishedAt = &finishedAt.Time
	}
	return &a, nil
}