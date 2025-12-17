package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
)

type CsirtRepository struct {
	db *sql.DB
}

func NewCsirtRepository(db *sql.DB) *CsirtRepository {
	return &CsirtRepository{db: db}
}

func (r *CsirtRepository) Create(req dto.CreateCsirtRequest, id string) error {
	_, err := r.db.Exec(`
		INSERT INTO csirt (
			id, id_perusahaan, nama_csirt, web_csirt, photo_csirt, file_rfc2350, file_public_key_pgp
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id,
		req.IdPerusahaan,
		req.NamaCsirt,
		req.WebCsirt,
		req.PhotoCsirt,
		req.FileRFC2350,
		req.FilePublicKeyPGP,
	)
	return err
}

func (r *CsirtRepository) GetAll() ([]models.Csirt, error) {
	rows, err := r.db.Query(`
		SELECT id, id_perusahaan, nama_csirt, web_csirt, photo_csirt, file_rfc2350, file_public_key_pgp
		FROM csirt
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Csirt
	for rows.Next() {
		var c models.Csirt
		rows.Scan(
			&c.ID,
			&c.IdPerusahaan,
			&c.NamaCsirt,
			&c.WebCsirt,
			&c.PhotoCsirt,
			&c.FileRFC2350,
			&c.FilePublicKeyPGP,
		)
		result = append(result, c)
	}
	return result, nil
}

func (r *CsirtRepository) GetByID(id string) (*models.Csirt, error) {
	row := r.db.QueryRow(`
		SELECT id, id_perusahaan, nama_csirt, web_csirt, photo_csirt, file_rfc2350, file_public_key_pgp
		FROM csirt WHERE id = ?`, id)

	var c models.Csirt
	err := row.Scan(
		&c.ID,
		&c.IdPerusahaan,
		&c.NamaCsirt,
		&c.WebCsirt,
		&c.PhotoCsirt,
		&c.FileRFC2350,
		&c.FilePublicKeyPGP,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CsirtRepository) Update(id string, c models.Csirt) error {
	_, err := r.db.Exec(`
		UPDATE csirt SET
			nama_csirt=?,
			web_csirt=?,
			photo_csirt=?,
			file_rfc2350=?,
			file_public_key_pgp=?
		WHERE id=?`,
		c.NamaCsirt,
		c.WebCsirt,
		c.PhotoCsirt,
		c.FileRFC2350,
		c.FilePublicKeyPGP,
		id,
	)
	return err
}

func (r *CsirtRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM csirt WHERE id=?`, id)
	return err
}
