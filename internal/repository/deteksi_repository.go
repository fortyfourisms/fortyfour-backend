package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"

	"github.com/rollbar/rollbar-go"
)

type DeteksiRepository struct {
	db *sql.DB
}

func NewDeteksiRepository(db *sql.DB) *DeteksiRepository {
	return &DeteksiRepository{db: db}
}

func (r *DeteksiRepository) Create(req dto.CreateDeteksiRequest, id string) error {
	// Hitung nilai_deteksi (rata-rata dari 3 subdomain)
	NilaiDeteksi := (req.NilaiSubdomain1 + req.NilaiSubdomain2 + req.NilaiSubdomain3) / 3.0

	query := `INSERT INTO deteksi (id, nilai_deteksi, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3) 
			  VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		id,
		NilaiDeteksi,
		req.NilaiSubdomain1,
		req.NilaiSubdomain2,
		req.NilaiSubdomain3,
	)
	return err
}

func (r *DeteksiRepository) GetAll() ([]models.Deteksi, error) {
	rows, err := r.db.Query(`
		SELECT id, nilai_deteksi, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3 
		FROM deteksi
	`)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var result []models.Deteksi
	for rows.Next() {
		var d models.Deteksi
		rows.Scan(
			&d.ID,
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
		SELECT id, nilai_deteksi, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3
		FROM deteksi WHERE id = ?`, id)

	var d models.Deteksi
	err := row.Scan(
		&d.ID,
		&d.NilaiDeteksi,
		&d.NilaiSubdomain1,
		&d.NilaiSubdomain2,
		&d.NilaiSubdomain3,
	)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	return &d, nil
}

func (r *DeteksiRepository) Update(id string, d models.Deteksi) error {
	_, err := r.db.Exec(`
		UPDATE deteksi SET
			nilai_deteksi=?,
			nilai_subdomain1=?,
			nilai_subdomain2=?,
			nilai_subdomain3=?
		WHERE id=?`,
		d.NilaiDeteksi,
		d.NilaiSubdomain1,
		d.NilaiSubdomain2,
		d.NilaiSubdomain3,
		id,
	)
	return err
}

func (r *DeteksiRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM deteksi WHERE id=?`, id)
	return err
}
