package repository

import (
	"database/sql"

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

func (r *FeedbackRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM catatan_pribadi WHERE id=?`, id)
	return err
}
