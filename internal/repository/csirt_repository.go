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

// nullStr safely converts sql.NullString to string
func nullStr(n sql.NullString) string {
	if n.Valid {
		return n.String
	}
	return ""
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
		return nil, err
	}
	return &c, nil
}

// scanCsirtWithPerusahaan is a shared helper to scan a CSIRT row that includes perusahaan JOIN data.
// Handles NULL values for all nullable columns.
func scanCsirtWithPerusahaan(scanner interface {
	Scan(dest ...any) error
}) (dto.CsirtResponse, error) {
	var csirt dto.CsirtResponse
	var perusahaan dto.PerusahaanResponse

	var (
		webCsirt, teleponCsirt, photoCsirt, fileRFC, filePGP sql.NullString
		photoPerusahaan, alamat, telepon, email, website      sql.NullString
		subID, namaSubSektor, idSektor, namaSektor            sql.NullString
		subCreatedAt, subUpdatedAt                            sql.NullString
	)

	err := scanner.Scan(
		&csirt.ID,
		&csirt.NamaCsirt,
		&webCsirt,
		&teleponCsirt,
		&photoCsirt,
		&fileRFC,
		&filePGP,
		&perusahaan.ID,
		&photoPerusahaan,
		&perusahaan.NamaPerusahaan,
		&alamat,
		&telepon,
		&email,
		&website,
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
		return dto.CsirtResponse{}, err
	}

	csirt.WebCsirt = nullStr(webCsirt)
	csirt.PhotoCsirt = nullStr(photoCsirt)
	csirt.FileRFC2350 = nullStr(fileRFC)
	csirt.FilePublicKeyPGP = nullStr(filePGP)
	if teleponCsirt.Valid {
		csirt.TeleponCsirt = &teleponCsirt.String
	}

	perusahaan.Photo = nullStr(photoPerusahaan)
	perusahaan.Alamat = nullStr(alamat)
	perusahaan.Telepon = nullStr(telepon)
	perusahaan.Email = nullStr(email)
	perusahaan.Website = nullStr(website)

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
	return csirt, nil
}

const csirtWithPerusahaanQuery = `
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
	LEFT JOIN sektor s ON ss.id_sektor = s.id`

/*
========================
GET ALL + PERUSAHAAN
========================
*/
func (r *CsirtRepository) GetAllWithPerusahaan() ([]dto.CsirtResponse, error) {
	rows, err := r.db.Query(csirtWithPerusahaanQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.CsirtResponse
	for rows.Next() {
		csirt, err := scanCsirtWithPerusahaan(rows)
		if err != nil {
			return nil, err
		}
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
	row := r.db.QueryRow(csirtWithPerusahaanQuery+` WHERE c.id = ?`, id)
	csirt, err := scanCsirtWithPerusahaan(row)
	if err != nil {
		return nil, err
	}
	return &csirt, nil
}

/*
========================
GET BY PERUSAHAAN
========================
*/
func (r *CsirtRepository) GetByPerusahaan(idPerusahaan string) ([]dto.CsirtResponse, error) {
	rows, err := r.db.Query(csirtWithPerusahaanQuery+` WHERE c.id_perusahaan = ?`, idPerusahaan)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.CsirtResponse
	for rows.Next() {
		csirt, err := scanCsirtWithPerusahaan(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, csirt)
	}
	return result, nil
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