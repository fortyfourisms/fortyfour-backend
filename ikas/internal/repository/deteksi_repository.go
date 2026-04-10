package repository

import (
	"database/sql"
	"ikas/internal/models"
)

type DeteksiRepository struct {
	db *sql.DB
}

func NewDeteksiRepository(db *sql.DB) *DeteksiRepository {
	return &DeteksiRepository{db: db}
}

func (r *DeteksiRepository) GetAll() ([]models.Deteksi, error) {
	rows, err := r.db.Query(`
		SELECT id, perusahaan_id, nilai_deteksi, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3 
		FROM deteksi
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Deteksi
	for rows.Next() {
		var d models.Deteksi
		rows.Scan(
			&d.ID,
			&d.PerusahaanID,
			&d.NilaiDeteksi,
			&d.NilaiSubdomain1,
			&d.NilaiSubdomain2,
			&d.NilaiSubdomain3,
		)
		result = append(result, d)
	}
	return result, nil
}

func (r *DeteksiRepository) GetByPerusahaan(perusahaanID string) ([]models.Deteksi, error) {
	rows, err := r.db.Query(`
		SELECT id, perusahaan_id, nilai_deteksi, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3 
		FROM deteksi WHERE perusahaan_id = ?`, perusahaanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Deteksi
	for rows.Next() {
		var d models.Deteksi
		rows.Scan(
			&d.ID,
			&d.PerusahaanID,
			&d.NilaiDeteksi,
			&d.NilaiSubdomain1,
			&d.NilaiSubdomain2,
			&d.NilaiSubdomain3,
		)
		result = append(result, d)
	}
	return result, nil
}

func (r *DeteksiRepository) GetByID(id string) (*models.Deteksi, error) {
	row := r.db.QueryRow(`
		SELECT id, perusahaan_id, nilai_deteksi, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3
		FROM deteksi WHERE id = ?`, id)

	var d models.Deteksi
	err := row.Scan(
		&d.ID,
		&d.PerusahaanID,
		&d.NilaiDeteksi,
		&d.NilaiSubdomain1,
		&d.NilaiSubdomain2,
		&d.NilaiSubdomain3,
	)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
