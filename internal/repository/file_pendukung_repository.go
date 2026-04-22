package repository

import (
	"database/sql"
	"fmt"

	"fortyfour-backend/internal/models"
)

type FilePendukungRepository struct {
	db *sql.DB
}

func NewFilePendukungRepository(db *sql.DB) *FilePendukungRepository {
	return &FilePendukungRepository{db: db}
}

var _ FilePendukungRepositoryInterface = (*FilePendukungRepository)(nil)

func (r *FilePendukungRepository) Create(fp *models.FilePendukung) error {
	_, err := r.db.Exec(
		`INSERT INTO file_pendukung (id, id_materi, nama_file, file_path, ukuran, created_at)
		 VALUES (?, ?, ?, ?, ?, NOW())`,
		fp.ID, fp.IDMateri, fp.NamaFile, fp.FilePath, fp.Ukuran,
	)
	return err
}

func (r *FilePendukungRepository) FindByMateri(idMateri string) ([]models.FilePendukung, error) {
	rows, err := r.db.Query(
		`SELECT id, id_materi, nama_file, file_path, ukuran, created_at
		 FROM file_pendukung WHERE id_materi = ? ORDER BY created_at ASC`, idMateri,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.FilePendukung
	for rows.Next() {
		var fp models.FilePendukung
		if err := rows.Scan(&fp.ID, &fp.IDMateri, &fp.NamaFile, &fp.FilePath, &fp.Ukuran, &fp.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, fp)
	}
	return result, nil
}

func (r *FilePendukungRepository) FindByID(id string) (*models.FilePendukung, error) {
	row := r.db.QueryRow(
		`SELECT id, id_materi, nama_file, file_path, ukuran, created_at
		 FROM file_pendukung WHERE id = ?`, id,
	)
	var fp models.FilePendukung
	err := row.Scan(&fp.ID, &fp.IDMateri, &fp.NamaFile, &fp.FilePath, &fp.Ukuran, &fp.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &fp, nil
}

func (r *FilePendukungRepository) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM file_pendukung WHERE id=?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("file pendukung dengan id %s tidak ditemukan", id)
	}
	return nil
}
