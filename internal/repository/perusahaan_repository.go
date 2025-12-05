package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
)

type PerusahaanRepository struct {
	db *sql.DB
}

func NewPerusahaanRepository(db *sql.DB) *PerusahaanRepository {
	return &PerusahaanRepository{db: db}
}

func (r *PerusahaanRepository) Create(req dto.CreatePerusahaanRequest, id string) error {
	_, err := r.db.Exec(`INSERT INTO perusahaan
        (id, photo, nama_perusahaan, jenis_usaha, alamat, telepon, email, website)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id,
		valueOrEmpty(req.Photo),
		valueOrEmpty(req.NamaPerusahaan),
		valueOrEmpty(req.JenisUsaha),
		valueOrEmpty(req.Alamat),
		valueOrEmpty(req.Telepon),
		valueOrEmpty(req.Email),
		valueOrEmpty(req.Website),
	)
	return err
}

func (r *PerusahaanRepository) GetAll() ([]dto.PerusahaanResponse, error) {
	rows, err := r.db.Query(`SELECT id, photo, nama_perusahaan, jenis_usaha, alamat, telepon, email, website FROM perusahaan`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.PerusahaanResponse
	for rows.Next() {
		var p dto.PerusahaanResponse
		rows.Scan(&p.ID, &p.Photo, &p.NamaPerusahaan, &p.JenisUsaha, &p.Alamat, &p.Telepon, &p.Email, &p.Website)
		result = append(result, p)
	}
	return result, nil
}

func (r *PerusahaanRepository) GetByID(id string) (*dto.PerusahaanResponse, error) {
	row := r.db.QueryRow(`SELECT id, photo, nama_perusahaan, jenis_usaha, alamat, telepon, email, website FROM perusahaan WHERE id=?`, id)
	var p dto.PerusahaanResponse
	err := row.Scan(&p.ID, &p.Photo, &p.NamaPerusahaan, &p.JenisUsaha, &p.Alamat, &p.Telepon, &p.Email, &p.Website)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PerusahaanRepository) Update(id string, p dto.PerusahaanResponse) error {
	_, err := r.db.Exec(`UPDATE perusahaan SET
        photo=?, nama_perusahaan=?, jenis_usaha=?, alamat=?, telepon=?, email=?, website=?
        WHERE id=?`,
		p.Photo, p.NamaPerusahaan, p.JenisUsaha, p.Alamat, p.Telepon, p.Email, p.Website, id,
	)
	return err
}

func (r *PerusahaanRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM perusahaan WHERE id=?`, id)
	return err
}

func valueOrEmpty(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
