package repository

import (
	"database/sql"

	"fortyfour-backend/internal/models"
)

type CatatanRepository struct {
	db *sql.DB
}

func NewCatatanRepository(db *sql.DB) *CatatanRepository {
	return &CatatanRepository{db: db}
}

var _ CatatanRepositoryInterface = (*CatatanRepository)(nil)

func (r *CatatanRepository) Upsert(c *models.CatatanPribadi) error {
	_, err := r.db.Exec(
		`INSERT INTO catatan_pribadi (id, id_materi, id_user, konten, created_at, updated_at)
		 VALUES (?, ?, ?, ?, NOW(), NOW())
		 ON DUPLICATE KEY UPDATE
		    konten     = VALUES(konten),
		    updated_at = NOW()`,
		c.ID, c.IDMateri, c.IDUser, c.Konten,
	)
	return err
}

func (r *CatatanRepository) FindByUserAndMateri(idUser, idMateri string) (*models.CatatanPribadi, error) {
	row := r.db.QueryRow(
		`SELECT id, id_materi, id_user, konten, created_at, updated_at
		 FROM catatan_pribadi WHERE id_user=? AND id_materi=?`,
		idUser, idMateri,
	)
	var c models.CatatanPribadi
	err := row.Scan(&c.ID, &c.IDMateri, &c.IDUser, &c.Konten, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CatatanRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM catatan_pribadi WHERE id=?`, id)
	return err
}
