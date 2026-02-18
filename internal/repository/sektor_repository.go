package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
)

type SektorRepository struct {
	db *sql.DB
}

func NewSektorRepository(db *sql.DB) *SektorRepository {
	return &SektorRepository{db: db}
}

func (r *SektorRepository) GetAll() ([]dto.SektorResponse, error) {
	// Query sektor dengan sub sektor
	rows, err := r.db.Query(`
		SELECT s.id, s.nama_sektor, s.created_at, s.updated_at,
		       ss.id as sub_id, ss.nama_sub_sektor, ss.id_sektor, ss.created_at as sub_created_at, ss.updated_at as sub_updated_at
		FROM sektor s
		LEFT JOIN sub_sektor ss ON s.id = ss.id_sektor
		ORDER BY s.nama_sektor, ss.nama_sub_sektor
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sektorMap := make(map[string]*dto.SektorResponse)

	for rows.Next() {
		var sektorID, namaSektor, createdAt, updatedAt string
		var subID, namaSubSektor, idSektor, subCreatedAt, subUpdatedAt sql.NullString

		err := rows.Scan(&sektorID, &namaSektor, &createdAt, &updatedAt,
			&subID, &namaSubSektor, &idSektor, &subCreatedAt, &subUpdatedAt)
		if err != nil {
			return nil, err
		}

		// Cek apakah sektor sudah ada di map
		if _, exists := sektorMap[sektorID]; !exists {
			sektorMap[sektorID] = &dto.SektorResponse{
				ID:         sektorID,
				NamaSektor: namaSektor,
				CreatedAt:  createdAt,
				UpdatedAt:  updatedAt,
				SubSektor:  []dto.SubSektorResponse{},
			}
		}

		// Tambahkan sub sektor jika ada
		if subID.Valid {
			sektorMap[sektorID].SubSektor = append(sektorMap[sektorID].SubSektor, dto.SubSektorResponse{
				ID:            subID.String,
				NamaSubSektor: namaSubSektor.String,
				IDSektor:      idSektor.String,
				CreatedAt:     subCreatedAt.String,
				UpdatedAt:     subUpdatedAt.String,
			})
		}
	}

	// Convert map to slice
	result := make([]dto.SektorResponse, 0, len(sektorMap))
	for _, sektor := range sektorMap {
		result = append(result, *sektor)
	}

	return result, nil
}

func (r *SektorRepository) GetByID(id string) (*dto.SektorResponse, error) {
	row := r.db.QueryRow(`SELECT id, nama_sektor, created_at, updated_at FROM sektor WHERE id=?`, id)

	var sektor dto.SektorResponse
	err := row.Scan(&sektor.ID, &sektor.NamaSektor, &sektor.CreatedAt, &sektor.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Get sub sektor
	rows, err := r.db.Query(`SELECT id, nama_sub_sektor, id_sektor, created_at, updated_at FROM sub_sektor WHERE id_sektor=?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sektor.SubSektor = []dto.SubSektorResponse{}
	for rows.Next() {
		var sub dto.SubSektorResponse
		rows.Scan(&sub.ID, &sub.NamaSubSektor, &sub.IDSektor, &sub.CreatedAt, &sub.UpdatedAt)
		sektor.SubSektor = append(sektor.SubSektor, sub)
	}

	return &sektor, nil
}
