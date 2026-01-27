package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/utils"

	"github.com/rollbar/rollbar-go"
)

type PerusahaanRepository struct {
	db *sql.DB
}

func NewPerusahaanRepository(db *sql.DB) *PerusahaanRepository {
	return &PerusahaanRepository{db: db}
}

func (r *PerusahaanRepository) Create(req dto.CreatePerusahaanRequest, id string) error {
	_, err := r.db.Exec(`INSERT INTO perusahaan
        (id, photo, nama_perusahaan, sektor, alamat, telepon, email, website)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id,
		utils.ValueOrEmpty(req.Photo),
		utils.ValueOrEmpty(req.NamaPerusahaan),
		utils.ValueOrEmpty(req.Sektor),
		utils.ValueOrEmpty(req.Alamat),
		utils.ValueOrEmpty(req.Telepon),
		utils.ValueOrEmpty(req.Email),
		utils.ValueOrEmpty(req.Website),
	)
	return err
}

func (r *PerusahaanRepository) GetAll() ([]dto.PerusahaanResponse, error) {
	rows, err := r.db.Query(`SELECT id, photo, nama_perusahaan, sektor, alamat, telepon, email, website, created_at, updated_at FROM perusahaan`)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var result []dto.PerusahaanResponse
	for rows.Next() {
		var p dto.PerusahaanResponse
		rows.Scan(&p.ID, &p.Photo, &p.NamaPerusahaan, &p.Sektor, &p.Alamat, &p.Telepon, &p.Email, &p.Website, &p.CreatedAt, &p.UpdatedAt)
		result = append(result, p)
	}
	return result, nil
}

func (r *PerusahaanRepository) GetLatest(limit int) ([]dto.PerusahaanResponse, error) {
	rows, err := r.db.Query(`
		SELECT id, photo, nama_perusahaan, sektor, alamat, telepon, email, website, created_at, updated_at
		FROM perusahaan
		ORDER BY created_at DESC
		LIMIT ?
	`, limit)

	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var result []dto.PerusahaanResponse
	for rows.Next() {
		var p dto.PerusahaanResponse
		rows.Scan(
			&p.ID,
			&p.Photo,
			&p.NamaPerusahaan,
			&p.Sektor,
			&p.Alamat,
			&p.Telepon,
			&p.Email,
			&p.Website,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		result = append(result, p)
	}

	return result, nil
}

func (r *PerusahaanRepository) GetByID(id string) (*dto.PerusahaanResponse, error) {
	row := r.db.QueryRow(`SELECT id, photo, nama_perusahaan, sektor, alamat, telepon, email, website, created_at, updated_at FROM perusahaan WHERE id=?`, id)
	var p dto.PerusahaanResponse
	err := row.Scan(&p.ID, &p.Photo, &p.NamaPerusahaan, &p.Sektor, &p.Alamat, &p.Telepon, &p.Email, &p.Website, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	return &p, nil
}

func (r *PerusahaanRepository) Update(id string, p dto.PerusahaanResponse) error {
	_, err := r.db.Exec(`UPDATE perusahaan SET
        photo=?, nama_perusahaan=?, sektor=?, alamat=?, telepon=?, email=?, website=?, updated_at=CURRENT_TIMESTAMP
        WHERE id=?`,
		p.Photo, p.NamaPerusahaan, p.Sektor, p.Alamat, p.Telepon, p.Email, p.Website, id,
	)
	return err
}

func (r *PerusahaanRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM perusahaan WHERE id=?`, id)
	return err
}
