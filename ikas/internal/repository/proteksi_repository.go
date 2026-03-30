package repository

import (
	"database/sql"
	"ikas/internal/models"
)

type ProteksiRepository struct {
	db *sql.DB
}

func NewProteksiRepository(db *sql.DB) *ProteksiRepository {
	return &ProteksiRepository{db: db}
}

func (r *ProteksiRepository) GetAll() ([]models.Proteksi, error) {
	query := `SELECT id, nilai_proteksi, nilai_subdomain1, nilai_subdomain2, 
	          nilai_subdomain3, nilai_subdomain4, nilai_subdomain5, nilai_subdomain6 
	          FROM proteksi 
	          ORDER BY id DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var proteksiList []models.Proteksi
	for rows.Next() {
		var proteksi models.Proteksi
		err := rows.Scan(
			&proteksi.ID,
			&proteksi.NilaiProteksi,
			&proteksi.NilaiSubdomain1,
			&proteksi.NilaiSubdomain2,
			&proteksi.NilaiSubdomain3,
			&proteksi.NilaiSubdomain4,
			&proteksi.NilaiSubdomain5,
			&proteksi.NilaiSubdomain6,
		)
		if err != nil {
			return nil, err
		}
		proteksiList = append(proteksiList, proteksi)
	}

	return proteksiList, nil
}

func (r *ProteksiRepository) GetByID(id string) (*models.Proteksi, error) {
	var proteksi models.Proteksi
	query := `SELECT id, nilai_proteksi, nilai_subdomain1, nilai_subdomain2, 
	          nilai_subdomain3, nilai_subdomain4, nilai_subdomain5, nilai_subdomain6 
	          FROM proteksi 
	          WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(
		&proteksi.ID,
		&proteksi.NilaiProteksi,
		&proteksi.NilaiSubdomain1,
		&proteksi.NilaiSubdomain2,
		&proteksi.NilaiSubdomain3,
		&proteksi.NilaiSubdomain4,
		&proteksi.NilaiSubdomain5,
		&proteksi.NilaiSubdomain6,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &proteksi, nil
}
