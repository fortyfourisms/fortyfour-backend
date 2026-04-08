package repository

import (
	"database/sql"
	"survey/internal/models"
)

type RespondenRepository struct {
	db *sql.DB
}

func NewRespondenRepository(db *sql.DB) *RespondenRepository {
	return &RespondenRepository{db: db}
}

func (r *RespondenRepository) Create(m models.Responden) error {
	query := `INSERT INTO responden
		(nama_lengkap, jabatan, perusahaan, email, no_telepon,
		 sektor, sektor_lainnya, sertifikat_training)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		m.NamaLengkap,
		m.Jabatan,
		m.Perusahaan,
		m.Email,
		m.NoTelepon,
		m.Sektor,
		m.SektorLainnya,
		m.SertifikatTraining,
	)

	return err
}

func (r *RespondenRepository) GetAll() ([]models.Responden, error) {
	rows, err := r.db.Query(`
		SELECT id, nama_lengkap, jabatan, perusahaan, email, no_telepon,
		       sektor, sektor_lainnya, sertifikat_training, created_at, updated_at
		FROM responden`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Responden
	for rows.Next() {
		var m models.Responden
		rows.Scan(
			&m.ID,
			&m.NamaLengkap,
			&m.Jabatan,
			&m.Perusahaan,
			&m.Email,
			&m.NoTelepon,
			&m.Sektor,
			&m.SektorLainnya,
			&m.SertifikatTraining,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		result = append(result, m)
	}
	return result, nil
}

func (r *RespondenRepository) GetByID(id int) (*models.Responden, error) {
	row := r.db.QueryRow(`
		SELECT id, nama_lengkap, jabatan, perusahaan, email, no_telepon,
		       sektor, sektor_lainnya, sertifikat_training, created_at, updated_at
		FROM responden WHERE id = ?`, id)

	var m models.Responden
	if err := row.Scan(
		&m.ID,
		&m.NamaLengkap,
		&m.Jabatan,
		&m.Perusahaan,
		&m.Email,
		&m.NoTelepon,
		&m.Sektor,
		&m.SektorLainnya,
		&m.SertifikatTraining,
		&m.CreatedAt,
		&m.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &m, nil
}

func (r *RespondenRepository) Update(id int, m models.Responden) error {
	_, err := r.db.Exec(`
		UPDATE responden SET
		nama_lengkap=?,
		jabatan=?,
		perusahaan=?,
		email=?,
		no_telepon=?,
		sektor=?,
		sektor_lainnya=?,
		sertifikat_training=?,
		updated_at=NOW()
		WHERE id=?`,
		m.NamaLengkap,
		m.Jabatan,
		m.Perusahaan,
		m.Email,
		m.NoTelepon,
		m.Sektor,
		m.SektorLainnya,
		m.SertifikatTraining,
		id,
	)

	return err
}

func (r *RespondenRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM responden WHERE id=?`, id)
	return err
}
