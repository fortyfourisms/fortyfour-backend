package repository

import (
	"database/sql"
	"ikas/internal/models"
	"strconv"
)

type DeteksiRepository struct {
	db *sql.DB
}

func NewDeteksiRepository(db *sql.DB) *DeteksiRepository {
	return &DeteksiRepository{db: db}
}

func (r *DeteksiRepository) GetAll() ([]models.Deteksi, error) {
	rows, err := r.db.Query(`
		SELECT id, ikas_id, nilai_deteksi, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3 
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
			&d.IkasID,
			&d.NilaiDeteksi,
			&d.NilaiSubdomain1,
			&d.NilaiSubdomain2,
			&d.NilaiSubdomain3,
		)
		result = append(result, d)
	}
	return result, nil
}

func (r *DeteksiRepository) GetByIkasID(ikasID string) ([]models.Deteksi, error) {
	rows, err := r.db.Query(`
		SELECT id, ikas_id, nilai_deteksi, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3 
		FROM deteksi WHERE ikas_id = ?`, ikasID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Deteksi
	for rows.Next() {
		var d models.Deteksi
		rows.Scan(
			&d.ID,
			&d.IkasID,
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
		SELECT id, ikas_id, nilai_deteksi, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3
		FROM deteksi WHERE id = ?`, id)

	var d models.Deteksi
	err := row.Scan(
		&d.ID,
		&d.IkasID,
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
func (r *DeteksiRepository) GetByPerusahaanID(perusahaanID string) ([]models.Deteksi, error) {
	query := `
		SELECT t.id, t.ikas_id, t.nilai_deteksi, t.nilai_subdomain1, t.nilai_subdomain2, t.nilai_subdomain3 
		FROM deteksi t
		JOIN ikas i ON t.ikas_id = i.id
		WHERE i.id_perusahaan = ?
	`
	rows, err := r.db.Query(query, perusahaanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Deteksi
	for rows.Next() {
		var d models.Deteksi
		rows.Scan(
			&d.ID,
			&d.IkasID,
			&d.NilaiDeteksi,
			&d.NilaiSubdomain1,
			&d.NilaiSubdomain2,
			&d.NilaiSubdomain3,
		)
		result = append(result, d)
	}
	return result, nil
}

func (r *DeteksiRepository) CloneByIkasID(sourceIkasID, targetIkasID string) (string, error) {
	query := `
		INSERT INTO deteksi 
			(ikas_id, nilai_deteksi, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3)
		SELECT 
			?, nilai_deteksi, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3
		FROM deteksi 
		WHERE ikas_id = ?`

	res, err := r.db.Exec(query, targetIkasID, sourceIkasID)
	if err != nil {
		return "", err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(lastID, 10), nil
}
