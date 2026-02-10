package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
)

type GulihRepository struct {
	db *sql.DB
}

func NewGulihRepository(db *sql.DB) *GulihRepository {
	return &GulihRepository{db: db}
}

func (r *GulihRepository) Create(req dto.CreateGulihRequest, id string) error {
	// Hitung nilai_gulih (rata-rata dari 4 subdomain)
	NilaiGulih := (req.NilaiSubdomain1 + req.NilaiSubdomain2 +
		req.NilaiSubdomain3 + req.NilaiSubdomain4) / 4.0

	query := `INSERT INTO gulih 
		(id, nilai_gulih, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3, nilai_subdomain4)
		VALUES (?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		id,
		NilaiGulih,
		req.NilaiSubdomain1,
		req.NilaiSubdomain2,
		req.NilaiSubdomain3,
		req.NilaiSubdomain4,
	)
	return err
}

func (r *GulihRepository) GetAll() ([]models.Gulih, error) {
	rows, err := r.db.Query(`
		SELECT id, nilai_gulih, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3, nilai_subdomain4
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
		SELECT id, nilai_gulih, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3, nilai_subdomain4
		FROM gulih WHERE id = ?`, id)

	var g models.Gulih
	if err := row.Scan(
		&g.ID,
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

func (r *GulihRepository) Update(id string, g models.Gulih) error {
	_, err := r.db.Exec(`
		UPDATE gulih SET
		nilai_gulih=?, nilai_subdomain1=?, nilai_subdomain2=?, nilai_subdomain3=?, nilai_subdomain4=?
		WHERE id=?`,
		g.NilaiGulih,
		g.NilaiSubdomain1,
		g.NilaiSubdomain2,
		g.NilaiSubdomain3,
		g.NilaiSubdomain4,
		id,
	)
	return err
}

func (r *GulihRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM gulih WHERE id=?`, id)
	return err
}
