package repository

import (
	"database/sql"

	"fortyfour-backend/internal/models"
)

type ProgressRepository struct {
	db *sql.DB
}

func NewProgressRepository(db *sql.DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

var _ ProgressRepositoryInterface = (*ProgressRepository)(nil)

// Upsert melakukan INSERT jika belum ada, UPDATE jika sudah ada,
// menggunakan ON DUPLICATE KEY UPDATE (memanfaatkan unique key user_materi).
func (r *ProgressRepository) Upsert(p *models.UserMateriProgress) error {
	_, err := r.db.Exec(
		`INSERT INTO user_materi_progress
		    (id, id_user, id_materi, is_completed, last_watched_seconds, completed_at, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
		 ON DUPLICATE KEY UPDATE
		    is_completed         = VALUES(is_completed),
		    last_watched_seconds = VALUES(last_watched_seconds),
		    completed_at         = VALUES(completed_at),
		    updated_at           = NOW()`,
		p.ID, p.IDUser, p.IDMateri,
		p.IsCompleted, p.LastWatchedSeconds, p.CompletedAt,
	)
	return err
}

func (r *ProgressRepository) FindByUserAndMateri(idUser, idMateri string) (*models.UserMateriProgress, error) {
	row := r.db.QueryRow(
		`SELECT id, id_user, id_materi, is_completed, last_watched_seconds,
		        completed_at, created_at, updated_at
		 FROM user_materi_progress
		 WHERE id_user=? AND id_materi=?`,
		idUser, idMateri,
	)
	return scanProgress(row)
}

func (r *ProgressRepository) FindByUserAndKelas(idUser, idKelas string) ([]models.UserMateriProgress, error) {
	rows, err := r.db.Query(
		`SELECT p.id, p.id_user, p.id_materi, p.is_completed, p.last_watched_seconds,
		        p.completed_at, p.created_at, p.updated_at
		 FROM user_materi_progress p
		 JOIN materi m ON m.id = p.id_materi
		 WHERE p.id_user=? AND m.id_kelas=?`,
		idUser, idKelas,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.UserMateriProgress
	for rows.Next() {
		var p models.UserMateriProgress
		var completedAt sql.NullTime
		if err := rows.Scan(
			&p.ID, &p.IDUser, &p.IDMateri, &p.IsCompleted, &p.LastWatchedSeconds,
			&completedAt, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if completedAt.Valid {
			p.CompletedAt = &completedAt.Time
		}
		result = append(result, p)
	}
	return result, nil
}

// HasCompletedAllMateri mengecek apakah user sudah menyelesaikan semua materi dalam kelas.
func (r *ProgressRepository) HasCompletedAllMateri(idUser, idKelas string) (bool, error) {
	// Hitung total materi dalam kelas
	var totalMateri int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM materi WHERE id_kelas = ?`, idKelas,
	).Scan(&totalMateri)
	if err != nil {
		return false, err
	}
	if totalMateri == 0 {
		return false, nil // tidak ada materi berarti belum selesai
	}

	// Hitung materi yang sudah diselesaikan user
	var completedCount int
	err = r.db.QueryRow(
		`SELECT COUNT(*) FROM user_materi_progress p
		 JOIN materi m ON m.id = p.id_materi
		 WHERE p.id_user   = ?
		   AND m.id_kelas  = ?
		   AND p.is_completed = 1`,
		idUser, idKelas,
	).Scan(&completedCount)
	if err != nil {
		return false, err
	}
	return completedCount >= totalMateri, nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func scanProgress(row *sql.Row) (*models.UserMateriProgress, error) {
	var p models.UserMateriProgress
	var completedAt sql.NullTime
	err := row.Scan(
		&p.ID, &p.IDUser, &p.IDMateri, &p.IsCompleted, &p.LastWatchedSeconds,
		&completedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if completedAt.Valid {
		p.CompletedAt = &completedAt.Time
	}
	return &p, nil
}
