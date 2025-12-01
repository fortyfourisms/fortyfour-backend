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

func (repo *PerusahaanRepository) Create(req dto.PerusahaanRequest) (int, error) {
	result, err := repo.db.Exec("INSERT INTO perusahaan (nama_perusahaan, jenis_usaha) VALUES (?, ?)", req.NamaPerusahaan, req.JenisUsaha)
	if err != nil {
		return 0, err
	}

	id, _ := result.LastInsertId()
	return int(id), nil
}

func (repo *PerusahaanRepository) GetAll() ([]dto.PerusahaanResponse, error) {
	rows, err := repo.db.Query("SELECT id, nama_perusahaan, jenis_usaha FROM perusahaan")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perusahaan []dto.PerusahaanResponse
	for rows.Next() {
		var p dto.PerusahaanResponse
		rows.Scan(&p.ID, &p.NamaPerusahaan, &p.JenisUsaha)
		perusahaan = append(perusahaan, p)
	}
	return perusahaan, nil
}

func (repo *PerusahaanRepository) GetByID(id int) (*dto.PerusahaanResponse, error) {
	row := repo.db.QueryRow("SELECT id, nama_perusahaan, jenis_usaha FROM perusahaan WHERE id = ?", id)

	var p dto.PerusahaanResponse
	err := row.Scan(&p.ID, &p.NamaPerusahaan, &p.JenisUsaha)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *PerusahaanRepository) Update(id int, req dto.PerusahaanRequest) error {
	_, err := repo.db.Exec("UPDATE perusahaan SET nama_perusahaan = ?, jenis_usaha = ? WHERE id = ?", req.NamaPerusahaan, req.JenisUsaha, id)
	return err
}

func (repo *PerusahaanRepository) Delete(id int) error {
	_, err := repo.db.Exec("DELETE FROM perusahaan WHERE id = ?", id)
	return err
}
