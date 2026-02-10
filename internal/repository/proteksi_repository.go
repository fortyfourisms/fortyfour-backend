package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"

)

type ProteksiRepository struct {
	db *sql.DB
}

func NewProteksiRepository(db *sql.DB) *ProteksiRepository {
	return &ProteksiRepository{db: db}
}

func (r *ProteksiRepository) Create(req dto.CreateProteksiRequest, id string) error {
	// Hitung nilai_proteksi (rata-rata dari 6 subdomain)
	nilaiProteksi := (req.NilaiSubdomain1 + req.NilaiSubdomain2 +
		req.NilaiSubdomain3 + req.NilaiSubdomain4 + req.NilaiSubdomain5 + req.NilaiSubdomain6) / 6.0

	query := `INSERT INTO proteksi (id, nilai_proteksi, nilai_subdomain1, nilai_subdomain2, 
	          nilai_subdomain3, nilai_subdomain4, nilai_subdomain5, nilai_subdomain6) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		id,
		nilaiProteksi,
		req.NilaiSubdomain1,
		req.NilaiSubdomain2,
		req.NilaiSubdomain3,
		req.NilaiSubdomain4,
		req.NilaiSubdomain5,
		req.NilaiSubdomain6,
	)

	return err
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

func (r *ProteksiRepository) Update(id string, proteksi models.Proteksi) error {
	query := `UPDATE proteksi 
	          SET nilai_proteksi = ?, nilai_subdomain1 = ?, nilai_subdomain2 = ?, 
	              nilai_subdomain3 = ?, nilai_subdomain4 = ?, nilai_subdomain5 = ?, 
	              nilai_subdomain6 = ? 
	          WHERE id = ?`

	_, err := r.db.Exec(query,
		proteksi.NilaiProteksi,
		proteksi.NilaiSubdomain1,
		proteksi.NilaiSubdomain2,
		proteksi.NilaiSubdomain3,
		proteksi.NilaiSubdomain4,
		proteksi.NilaiSubdomain5,
		proteksi.NilaiSubdomain6,
		id,
	)

	return err
}

func (r *ProteksiRepository) Delete(id string) error {
	query := `DELETE FROM proteksi WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
