package repository

import (
	"database/sql"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
)

type FeedbackRepository struct {
	db *sql.DB
}

func NewFeedbackRepository(db *sql.DB) *FeedbackRepository {
	return &FeedbackRepository{db: db}
}

var _ FeedbackRepositoryInterface = (*FeedbackRepository)(nil)

func (r *FeedbackRepository) Upsert(f *models.Feedback) error {
	_, err := r.db.Exec(
		`INSERT INTO catatan_pribadi (id, id_materi, id_user, konten, created_at, updated_at)
		 VALUES (?, ?, ?, ?, NOW(), NOW())
		 ON DUPLICATE KEY UPDATE
		    konten     = VALUES(konten),
		    updated_at = NOW()`,
		f.ID, f.IDMateri, f.IDUser, f.Konten,
	)
	return err
}

func (r *FeedbackRepository) FindByUserAndMateri(idUser, idMateri string) (*models.Feedback, error) {
	row := r.db.QueryRow(
		`SELECT id, id_materi, id_user, konten, created_at, updated_at
		 FROM catatan_pribadi WHERE id_user=? AND id_materi=?`,
		idUser, idMateri,
	)
	var f models.Feedback
	err := row.Scan(&f.ID, &f.IDMateri, &f.IDUser, &f.Konten, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// FindByMateri mengembalikan semua feedback untuk materi tertentu (untuk admin).
// Menyertakan username dari tabel users.
func (r *FeedbackRepository) FindByMateri(idMateri string) ([]dto.FeedbackListItem, error) {
	rows, err := r.db.Query(
		`SELECT cp.id, cp.id_materi, cp.id_user, COALESCE(u.username, '') AS username,
		        cp.konten, cp.created_at, cp.updated_at
		 FROM catatan_pribadi cp
		 LEFT JOIN users u ON cp.id_user = u.id
		 WHERE cp.id_materi = ?
		 ORDER BY cp.created_at DESC`,
		idMateri,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.FeedbackListItem
	for rows.Next() {
		var item dto.FeedbackListItem
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&item.ID, &item.IDMateri, &item.IDUser, &item.Username,
			&item.Konten, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *FeedbackRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM catatan_pribadi WHERE id=?`, id)
	return err
}
