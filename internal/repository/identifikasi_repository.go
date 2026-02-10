package repository

import (
	"database/sql"
	"fmt"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
)

type IdentifikasiRepository struct {
	db *sql.DB
}

func NewIdentifikasiRepository(db *sql.DB) *IdentifikasiRepository {
	return &IdentifikasiRepository{db: db}
}

func (r *IdentifikasiRepository) Create(req dto.CreateIdentifikasiRequest, id string) error {
	// Hitung nilai_identifikasi (rata-rata dari 5 subdomain)
	NilaiIdentifikasi := (req.NilaiSubdomain1 + req.NilaiSubdomain2 + req.NilaiSubdomain3 +
		req.NilaiSubdomain4 + req.NilaiSubdomain5) / 5.0

	query := `INSERT INTO identifikasi 
			(id, nilai_identifikasi, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3, nilai_subdomain4, nilai_subdomain5)
        	VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		id,
		NilaiIdentifikasi,
		req.NilaiIdentifikasi,
		req.NilaiSubdomain1,
		req.NilaiSubdomain2,
		req.NilaiSubdomain3,
		req.NilaiSubdomain4,
		req.NilaiSubdomain5,
	)
	return err
}

func (r *IdentifikasiRepository) GetAll() ([]models.Identifikasi, error) {
	rows, err := r.db.Query(`SELECT id, nilai_identifikasi, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3, nilai_subdomain4, nilai_subdomain5 FROM identifikasi`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Identifikasi
	for rows.Next() {
		var i models.Identifikasi
		rows.Scan(&i.ID, &i.NilaiIdentifikasi, &i.NilaiSubdomain1, &i.NilaiSubdomain2, &i.NilaiSubdomain3, &i.NilaiSubdomain4, &i.NilaiSubdomain5)
		result = append(result, i)
	}
	return result, nil
}

func (r *IdentifikasiRepository) GetByID(id string) (*models.Identifikasi, error) {
	row := r.db.QueryRow(`SELECT id, nilai_identifikasi, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3, nilai_subdomain4, nilai_subdomain5 FROM identifikasi WHERE id = ?`, id)
	var i models.Identifikasi
	err := row.Scan(&i.ID, &i.NilaiIdentifikasi, &i.NilaiSubdomain1, &i.NilaiSubdomain2, &i.NilaiSubdomain3, &i.NilaiSubdomain4, &i.NilaiSubdomain5)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (r *IdentifikasiRepository) Update(id string, i models.Identifikasi) error {
	_, err := r.db.Exec(`UPDATE identifikasi SET
        nilai_identifikasi=?, nilai_subdomain1=?, nilai_subdomain2=?, nilai_subdomain3=?, nilai_subdomain4=?, nilai_subdomain5=?
        WHERE id=?`,
		i.NilaiIdentifikasi, i.NilaiSubdomain1, i.NilaiSubdomain2, i.NilaiSubdomain3, i.NilaiSubdomain4, i.NilaiSubdomain5, id,
	)
	return err
}

func (r *IdentifikasiRepository) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM identifikasi WHERE id=?`, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("data dengan id %s tidak ditemukan", id)
	}

	return nil
}
