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
	var idSubSektor interface{}
	if req.IDSubSektor != nil && *req.IDSubSektor != "" {
		idSubSektor = *req.IDSubSektor
	} else {
		idSubSektor = nil
	}

	_, err := r.db.Exec(`INSERT INTO perusahaan
        (id, photo, nama_perusahaan, id_sub_sektor, alamat, telepon, email, website)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id,
		valueOrNull(req.Photo),
		valueOrEmpty(req.NamaPerusahaan),
		idSubSektor,
		valueOrNull(req.Alamat),
		valueOrNull(req.Telepon),
		valueOrNull(req.Email),
		valueOrNull(req.Website),
	)
	return err
}

func (r *PerusahaanRepository) GetAll() ([]dto.PerusahaanResponse, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.photo, p.nama_perusahaan, p.alamat, p.telepon, p.email, p.website, p.created_at, p.updated_at,
		       ss.id, ss.nama_sub_sektor, ss.id_sektor, ss.created_at, ss.updated_at,
		       s.nama_sektor
		FROM perusahaan p
		LEFT JOIN sub_sektor ss ON p.id_sub_sektor = ss.id
		LEFT JOIN sektor s ON ss.id_sektor = s.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.PerusahaanResponse
	for rows.Next() {
		var p dto.PerusahaanResponse
		var photo, alamat, telepon, email, website sql.NullString
		var subID, namaSubSektor, idSektor, namaSektor, subCreatedAt, subUpdatedAt sql.NullString

		err := rows.Scan(
			&p.ID, &photo, &p.NamaPerusahaan,
			&alamat, &telepon, &email, &website,
			&p.CreatedAt, &p.UpdatedAt,
			&subID, &namaSubSektor, &idSektor, &subCreatedAt, &subUpdatedAt,
			&namaSektor,
		)
		if err != nil {
			continue
		}

		p.Photo = photo.String
		p.Alamat = alamat.String
		p.Telepon = telepon.String
		p.Email = email.String
		p.Website = website.String

		// Tambahkan info sub sektor jika ada
		if subID.Valid {
			p.SubSektor = &dto.SubSektorResponse{
				ID:            subID.String,
				NamaSubSektor: namaSubSektor.String,
				IDSektor:      idSektor.String,
				NamaSektor:    namaSektor.String,
				CreatedAt:     subCreatedAt.String,
				UpdatedAt:     subUpdatedAt.String,
			}
		}

		result = append(result, p)
	}
	return result, nil
}

func (r *PerusahaanRepository) GetByID(id string) (*dto.PerusahaanResponse, error) {
	row := r.db.QueryRow(`
		SELECT p.id, p.photo, p.nama_perusahaan, p.alamat, p.telepon, p.email, p.website, p.created_at, p.updated_at,
		       ss.id, ss.nama_sub_sektor, ss.id_sektor, ss.created_at, ss.updated_at,
		       s.nama_sektor
		FROM perusahaan p
		LEFT JOIN sub_sektor ss ON p.id_sub_sektor = ss.id
		LEFT JOIN sektor s ON ss.id_sektor = s.id
		WHERE p.id=?
	`, id)

	var p dto.PerusahaanResponse
	var photo, alamat, telepon, email, website sql.NullString
	var subID, namaSubSektor, idSektor, namaSektor, subCreatedAt, subUpdatedAt sql.NullString

	err := row.Scan(
		&p.ID, &photo, &p.NamaPerusahaan,
		&alamat, &telepon, &email, &website,
		&p.CreatedAt, &p.UpdatedAt,
		&subID, &namaSubSektor, &idSektor, &subCreatedAt, &subUpdatedAt,
		&namaSektor,
	)
	if err != nil {
		return nil, err
	}

	p.Photo = photo.String
	p.Alamat = alamat.String
	p.Telepon = telepon.String
	p.Email = email.String
	p.Website = website.String

	// Tambahkan info sub sektor jika ada
	if subID.Valid {
		p.SubSektor = &dto.SubSektorResponse{
			ID:            subID.String,
			NamaSubSektor: namaSubSektor.String,
			IDSektor:      idSektor.String,
			NamaSektor:    namaSektor.String,
			CreatedAt:     subCreatedAt.String,
			UpdatedAt:     subUpdatedAt.String,
		}
	}

	return &p, nil
}

func (r *PerusahaanRepository) GetByNama(nama string) (*dto.PerusahaanResponse, error) {
	row := r.db.QueryRow(`
		SELECT p.id, p.photo, p.nama_perusahaan, p.alamat, p.telepon, p.email, p.website, p.created_at, p.updated_at,
		       ss.id, ss.nama_sub_sektor, ss.id_sektor, ss.created_at, ss.updated_at,
		       s.nama_sektor
		FROM perusahaan p
		LEFT JOIN sub_sektor ss ON p.id_sub_sektor = ss.id
		LEFT JOIN sektor s ON ss.id_sektor = s.id
		WHERE LOWER(p.nama_perusahaan) = LOWER(?)
	`, nama)

	var p dto.PerusahaanResponse
	var photo, alamat, telepon, email, website sql.NullString
	var subID, namaSubSektor, idSektor, namaSektor, subCreatedAt, subUpdatedAt sql.NullString

	err := row.Scan(
		&p.ID, &photo, &p.NamaPerusahaan,
		&alamat, &telepon, &email, &website,
		&p.CreatedAt, &p.UpdatedAt,
		&subID, &namaSubSektor, &idSektor, &subCreatedAt, &subUpdatedAt,
		&namaSektor,
	)
	if err != nil {
		return nil, err
	}

	p.Photo = photo.String
	p.Alamat = alamat.String
	p.Telepon = telepon.String
	p.Email = email.String
	p.Website = website.String

	if subID.Valid {
		p.SubSektor = &dto.SubSektorResponse{
			ID:            subID.String,
			NamaSubSektor: namaSubSektor.String,
			IDSektor:      idSektor.String,
			NamaSektor:    namaSektor.String,
			CreatedAt:     subCreatedAt.String,
			UpdatedAt:     subUpdatedAt.String,
		}
	}

	return &p, nil
}

func (r *PerusahaanRepository) Update(id string, p dto.PerusahaanResponse) error {
	var idSubSektor interface{}
	if p.SubSektor != nil {
		idSubSektor = p.SubSektor.ID
	} else {
		idSubSektor = nil
	}

	_, err := r.db.Exec(`UPDATE perusahaan SET
        photo=?, nama_perusahaan=?, id_sub_sektor=?, alamat=?, telepon=?, email=?, website=?, updated_at=CURRENT_TIMESTAMP
        WHERE id=?`,
		stringOrNull(p.Photo),
		p.NamaPerusahaan,
		idSubSektor,
		stringOrNull(p.Alamat),
		stringOrNull(p.Telepon),
		stringOrNull(p.Email),
		stringOrNull(p.Website),
		id,
	)
	return err
}

func (r *PerusahaanRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM perusahaan WHERE id=?`, id)
	return err
}

// Helper function for INSERT: return NULL if pointer is nil or empty
func valueOrNull(ptr *string) interface{} {
	if ptr == nil || *ptr == "" {
		return nil
	}
	return *ptr
}

// Helper function for INSERT: return empty string if pointer is nil (for required fields)
func valueOrEmpty(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// Helper function for UPDATE: return NULL if string is empty
func stringOrNull(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
