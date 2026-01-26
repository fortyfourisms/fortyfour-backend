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
        (id, photo, nama_perusahaan, id_sub_sektor, alamat, telepon, email, website)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id,
		utils.ValueOrEmpty(req.Photo),
		utils.ValueOrEmpty(req.NamaPerusahaan),
		utils.ValueOrEmpty(req.IDSubSektor),
		utils.ValueOrEmpty(req.Alamat),
		utils.ValueOrEmpty(req.Telepon),
		utils.ValueOrEmpty(req.Email),
		utils.ValueOrEmpty(req.Website),
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
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var result []dto.PerusahaanResponse
	for rows.Next() {
		var p dto.PerusahaanResponse
		var subID, namaSubSektor, idSektor, namaSektor, subCreatedAt, subUpdatedAt sql.NullString
		
		err := rows.Scan(
			&p.ID, &p.Photo, &p.NamaPerusahaan, 
			&p.Alamat, &p.Telepon, &p.Email, &p.Website, 
			&p.CreatedAt, &p.UpdatedAt,
			&subID, &namaSubSektor, &idSektor, &subCreatedAt, &subUpdatedAt, 
			&namaSektor,
		)
		if err != nil {
			continue
		}

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
	var subID, namaSubSektor, idSektor, namaSektor, subCreatedAt, subUpdatedAt sql.NullString
	
	err := row.Scan(
		&p.ID, &p.Photo, &p.NamaPerusahaan, 
		&p.Alamat, &p.Telepon, &p.Email, &p.Website, 
		&p.CreatedAt, &p.UpdatedAt,
		&subID, &namaSubSektor, &idSektor, &subCreatedAt, &subUpdatedAt,
		&namaSektor,
	)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

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

func (r *PerusahaanRepository) Update(id string, p dto.PerusahaanResponse) error {
	idSubSektor := ""
	if p.SubSektor != nil {
		idSubSektor = p.SubSektor.ID
	}
	
	_, err := r.db.Exec(`UPDATE perusahaan SET
        photo=?, nama_perusahaan=?, id_sub_sektor=?, alamat=?, telepon=?, email=?, website=?, updated_at=CURRENT_TIMESTAMP
        WHERE id=?`,
		p.Photo, p.NamaPerusahaan, idSubSektor, p.Alamat, p.Telepon, p.Email, p.Website, id,
	)
	return err
}

func (r *PerusahaanRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM perusahaan WHERE id=?`, id)
	return err
}