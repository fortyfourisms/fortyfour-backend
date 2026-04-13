package repository

import (
	"database/sql"
	"ikas/internal/models"
)

type GulihRepository struct {
	db *sql.DB
}

func NewGulihRepository(db *sql.DB) *GulihRepository {
	return &GulihRepository{db: db}
}

func (r *GulihRepository) GetAll() ([]models.Gulih, error) {
	rows, err := r.db.Query(`
		SELECT id, ikas_id, nilai_gulih, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3, nilai_subdomain4
		FROM gulih`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Gulih
	for rows.Next() {
		var g models.Gulih
		rows.Scan(
			&g.ID,
			&g.IkasID,
			&g.NilaiGulih,
			&g.NilaiSubdomain1,
			&g.NilaiSubdomain2,
			&g.NilaiSubdomain3,
			&g.NilaiSubdomain4,
		)
		result = append(result, g)
	}
	return result, nil
}

func (r *GulihRepository) GetByIkasID(ikasID string) ([]models.Gulih, error) {
	rows, err := r.db.Query(`
		SELECT id, ikas_id, nilai_gulih, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3, nilai_subdomain4
		FROM gulih WHERE ikas_id = ?`, ikasID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Gulih
	for rows.Next() {
		var g models.Gulih
		rows.Scan(
			&g.ID,
			&g.IkasID,
			&g.NilaiGulih,
			&g.NilaiSubdomain1,
			&g.NilaiSubdomain2,
			&g.NilaiSubdomain3,
			&g.NilaiSubdomain4,
		)
		result = append(result, g)
	}
	return result, nil
}

func (r *GulihRepository) GetByID(id string) (*models.Gulih, error) {
	row := r.db.QueryRow(`
		SELECT id, ikas_id, nilai_gulih, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3, nilai_subdomain4
		FROM gulih WHERE id = ?`, id)

	var g models.Gulih
	if err := row.Scan(
		&g.ID,
		&g.IkasID,
		&g.NilaiGulih,
		&g.NilaiSubdomain1,
		&g.NilaiSubdomain2,
		&g.NilaiSubdomain3,
		&g.NilaiSubdomain4,
	); err != nil {
		return nil, err
	}

	return &g, nil
}
