package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"

	"github.com/rollbar/rollbar-go"
)

type CsirtRepository struct {
	db *sql.DB
}

func NewCsirtRepository(db *sql.DB) *CsirtRepository {
	return &CsirtRepository{db: db}
}

/*
========================
CREATE
========================
*/
func (r *CsirtRepository) Create(req dto.CreateCsirtRequest, id string) error {
	_, err := r.db.Exec(`
		INSERT INTO csirt (
			id, id_perusahaan, nama_csirt, web_csirt, telepon_csirt,
			photo_csirt, file_rfc2350, file_public_key_pgp
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id,
		req.IdPerusahaan,
		req.NamaCsirt,
		req.WebCsirt,
		req.TeleponCsirt,
		req.PhotoCsirt,
		req.FileRFC2350,
		req.FilePublicKeyPGP,
	)
	return err
}

/*
========================
GET ALL
========================
*/
func (r *CsirtRepository) GetAll() ([]models.Csirt, error) {
	rows, err := r.db.Query(`
		SELECT id, id_perusahaan, nama_csirt, web_csirt, telepon_csirt,
		       photo_csirt, file_rfc2350, file_public_key_pgp
		FROM csirt
	`)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var result []models.Csirt
	for rows.Next() {
		var c models.Csirt
		err := rows.Scan(
			&c.ID,
			&c.IdPerusahaan,
			&c.NamaCsirt,
			&c.WebCsirt,
			&c.TeleponCsirt,
			&c.PhotoCsirt,
			&c.FileRFC2350,
			&c.FilePublicKeyPGP,
		)
		if err != nil {
			rollbar.Error(err)
			return nil, err
		}
		result = append(result, c)
	}
	return result, nil
}

/*
========================
GET BY ID
========================
*/
func (r *CsirtRepository) GetByID(id string) (*models.Csirt, error) {
	row := r.db.QueryRow(`
		SELECT id, id_perusahaan, nama_csirt, web_csirt, telepon_csirt,
		       photo_csirt, file_rfc2350, file_public_key_pgp
		FROM csirt WHERE id = ?`, id)

	var c models.Csirt
	err := row.Scan(
		&c.ID,
		&c.IdPerusahaan,
		&c.NamaCsirt,
		&c.WebCsirt,
		&c.TeleponCsirt,
		&c.PhotoCsirt,
		&c.FileRFC2350,
		&c.FilePublicKeyPGP,
	)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	return &c, nil
}

/*
========================
GET ALL + PERUSAHAAN
========================
*/
func (r *CsirtRepository) GetAllWithPerusahaan() ([]dto.CsirtResponse, error) {
	rows, err := r.db.Query(`
		SELECT 
			c.id, c.nama_csirt, c.web_csirt, c.telepon_csirt, 
			c.photo_csirt, c.file_rfc2350, c.file_public_key_pgp,
			p.id, p.photo, p.nama_perusahaan, 
			p.alamat, p.telepon, p.email, p.website,
			p.created_at, p.updated_at,
			ss.id, ss.nama_sub_sektor, ss.id_sektor, ss.created_at, ss.updated_at,
			s.nama_sektor
		FROM csirt c
		JOIN perusahaan p ON c.id_perusahaan = p.id
		LEFT JOIN sub_sektor ss ON p.id_sub_sektor = ss.id
		LEFT JOIN sektor s ON ss.id_sektor = s.id
	`)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var result []dto.CsirtResponse

	for rows.Next() {
		var csirt dto.CsirtResponse
		var perusahaan dto.PerusahaanResponse
		var subID, namaSubSektor, idSektor, namaSektor, subCreatedAt, subUpdatedAt sql.NullString

		err := rows.Scan(
			&csirt.ID,
			&csirt.NamaCsirt,
			&csirt.WebCsirt,
			&csirt.TeleponCsirt,
			&csirt.PhotoCsirt,
			&csirt.FileRFC2350,
			&csirt.FilePublicKeyPGP,
			&perusahaan.ID,
			&perusahaan.Photo,
			&perusahaan.NamaPerusahaan,
			&perusahaan.Alamat,
			&perusahaan.Telepon,
			&perusahaan.Email,
			&perusahaan.Website,
			&perusahaan.CreatedAt,
			&perusahaan.UpdatedAt,
			&subID,
			&namaSubSektor,
			&idSektor,
			&subCreatedAt,
			&subUpdatedAt,
			&namaSektor,
		)
		if err != nil {
			rollbar.Error(err)
			return nil, err
		}

		// Tambahkan info sub sektor jika ada
		if subID.Valid {
			perusahaan.SubSektor = &dto.SubSektorResponse{
				ID:            subID.String,
				NamaSubSektor: namaSubSektor.String,
				IDSektor:      idSektor.String,
				NamaSektor:    namaSektor.String,
				CreatedAt:     subCreatedAt.String,
				UpdatedAt:     subUpdatedAt.String,
			}
		}

		csirt.Perusahaan = perusahaan
		result = append(result, csirt)
	}

	return result, nil
}

/*
========================
GET BY ID + PERUSAHAAN
========================
*/
func (r *CsirtRepository) GetByIDWithPerusahaan(id string) (*dto.CsirtResponse, error) {
	row := r.db.QueryRow(`
		SELECT 
			c.id, c.nama_csirt, c.web_csirt, c.telepon_csirt, 
			c.photo_csirt, c.file_rfc2350, c.file_public_key_pgp,
			p.id, p.photo, p.nama_perusahaan,
			p.alamat, p.telepon, p.email, p.website,
			p.created_at, p.updated_at,
			ss.id, ss.nama_sub_sektor, ss.id_sektor, ss.created_at, ss.updated_at,
			s.nama_sektor
		FROM csirt c
		JOIN perusahaan p ON c.id_perusahaan = p.id
		LEFT JOIN sub_sektor ss ON p.id_sub_sektor = ss.id
		LEFT JOIN sektor s ON ss.id_sektor = s.id
		WHERE c.id = ?
	`, id)

	var csirt dto.CsirtResponse
	var perusahaan dto.PerusahaanResponse
	var subID, namaSubSektor, idSektor, namaSektor, subCreatedAt, subUpdatedAt sql.NullString

	err := row.Scan(
		&csirt.ID,
		&csirt.NamaCsirt,
		&csirt.WebCsirt,
		&csirt.TeleponCsirt,
		&csirt.PhotoCsirt,
		&csirt.FileRFC2350,
		&csirt.FilePublicKeyPGP,
		&perusahaan.ID,
		&perusahaan.Photo,
		&perusahaan.NamaPerusahaan,
		&perusahaan.Alamat,
		&perusahaan.Telepon,
		&perusahaan.Email,
		&perusahaan.Website,
		&perusahaan.CreatedAt,
		&perusahaan.UpdatedAt,
		&subID,
		&namaSubSektor,
		&idSektor,
		&subCreatedAt,
		&subUpdatedAt,
		&namaSektor,
	)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	// Tambahkan info sub sektor jika ada
	if subID.Valid {
		perusahaan.SubSektor = &dto.SubSektorResponse{
			ID:            subID.String,
			NamaSubSektor: namaSubSektor.String,
			IDSektor:      idSektor.String,
			NamaSektor:    namaSektor.String,
			CreatedAt:     subCreatedAt.String,
			UpdatedAt:     subUpdatedAt.String,
		}
	}

	csirt.Perusahaan = perusahaan
	return &csirt, nil
}

/*
========================
UPDATE
========================
*/
func (r *CsirtRepository) Update(id string, c models.Csirt) error {
	_, err := r.db.Exec(`
		UPDATE csirt SET
			nama_csirt = ?,
			web_csirt = ?,
			telepon_csirt = ?,
			photo_csirt = ?,
			file_rfc2350 = ?,
			file_public_key_pgp = ?
		WHERE id = ?`,
		c.NamaCsirt,
		c.WebCsirt,
		c.TeleponCsirt,
		c.PhotoCsirt,
		c.FileRFC2350,
		c.FilePublicKeyPGP,
		id,
	)
	return err
}

/*
========================
DELETE
========================
*/
func (r *CsirtRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM csirt WHERE id = ?`, id)
	return err
}