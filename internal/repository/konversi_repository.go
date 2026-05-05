package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
)

type KonversiRepositoryInterface interface {
	GetAllKonversi(perusahaanID string) ([]dto.KonversiResponse, error)
}

type konversiRepository struct {
	db *sql.DB
}

func NewKonversiRepository(db *sql.DB) KonversiRepositoryInterface {
	return &konversiRepository{db: db}
}

func (r *konversiRepository) GetAllKonversi(perusahaanID string) ([]dto.KonversiResponse, error) {
	query := `
		SELECT 
			p.id as perusahaan_id,
			p.nama_perusahaan,
			EXISTS(SELECT 1 FROM ikas i WHERE i.id_perusahaan = p.id) as has_ikas,
			EXISTS(SELECT 1 FROM se s WHERE s.id_perusahaan = p.id) as has_kse,
			EXISTS(SELECT 1 FROM responden r WHERE r.id_perusahaan = p.id) as has_survey,
			EXISTS(SELECT 1 FROM csirt c WHERE c.id_perusahaan = p.id) as has_csirt
		FROM perusahaan p
	`

	var args []interface{}
	if perusahaanID != "" {
		query += " WHERE p.id = ?"
		args = append(args, perusahaanID)
	}

	query += " ORDER BY p.nama_perusahaan ASC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []dto.KonversiResponse
	for rows.Next() {
		var res dto.KonversiResponse
		var ikas, kse, survey, csirt bool

		err := rows.Scan(
			&res.PerusahaanID,
			&res.NamaPerusahaan,
			&ikas,
			&kse,
			&survey,
			&csirt,
		)
		if err != nil {
			return nil, err
		}

		// Convert boolean to 1/0
		if ikas {
			res.PoinIkas = 1
		}
		if kse {
			res.PoinKse = 1
		}
		if survey {
			res.PoinSurvey = 1
		}
		if csirt {
			res.PoinCsirt = 1
		}

		results = append(results, res)
	}

	return results, nil
}
