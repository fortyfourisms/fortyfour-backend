package repository

import (
	"database/sql"
	"ikas/internal/models"
)

type IdentifikasiRepository struct {
	db *sql.DB
}

func NewIdentifikasiRepository(db *sql.DB) *IdentifikasiRepository {
	return &IdentifikasiRepository{db: db}
}

func (r *IdentifikasiRepository) GetAll() ([]models.Identifikasi, error) {
	rows, err := r.db.Query(`SELECT id, ikas_id, nilai_identifikasi, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3, nilai_subdomain4, nilai_subdomain5 FROM identifikasi`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Identifikasi
	for rows.Next() {
		var i models.Identifikasi
		rows.Scan(&i.ID, &i.IkasID, &i.NilaiIdentifikasi, &i.NilaiSubdomain1, &i.NilaiSubdomain2, &i.NilaiSubdomain3, &i.NilaiSubdomain4, &i.NilaiSubdomain5)
		result = append(result, i)
	}
	return result, nil
}

func (r *IdentifikasiRepository) GetByIkasID(ikasID string) ([]models.Identifikasi, error) {
	rows, err := r.db.Query(`SELECT id, ikas_id, nilai_identifikasi, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3, nilai_subdomain4, nilai_subdomain5 FROM identifikasi WHERE ikas_id = ?`, ikasID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Identifikasi
	for rows.Next() {
		var i models.Identifikasi
		rows.Scan(&i.ID, &i.IkasID, &i.NilaiIdentifikasi, &i.NilaiSubdomain1, &i.NilaiSubdomain2, &i.NilaiSubdomain3, &i.NilaiSubdomain4, &i.NilaiSubdomain5)
		result = append(result, i)
	}
	return result, nil
}

func (r *IdentifikasiRepository) GetByID(id string) (*models.Identifikasi, error) {
	row := r.db.QueryRow(`SELECT id, ikas_id, nilai_identifikasi, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3, nilai_subdomain4, nilai_subdomain5 FROM identifikasi WHERE id = ?`, id)
	var i models.Identifikasi
	err := row.Scan(&i.ID, &i.IkasID, &i.NilaiIdentifikasi, &i.NilaiSubdomain1, &i.NilaiSubdomain2, &i.NilaiSubdomain3, &i.NilaiSubdomain4, &i.NilaiSubdomain5)
	if err != nil {
		return nil, err
	}
	return &i, nil
}
