package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/utils"

)

type SeCsirtRepository struct {
	db *sql.DB
}

func NewSeCsirtRepository(db *sql.DB) *SeCsirtRepository {
	return &SeCsirtRepository{db: db}
}

func (r *SeCsirtRepository) Create(req dto.CreateSeCsirtRequest, id string) error {
	_, err := r.db.Exec(`
		INSERT INTO se_csirt
		(id, id_csirt, nama_se, ip_se, as_number_se, pengelola_se, fitur_se, kategori_se)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		id,
		req.IdCsirt,
		utils.ValueOrEmpty(req.NamaSe),
		utils.ValueOrEmpty(req.IpSe),
		utils.ValueOrEmpty(req.AsNumberSe),
		utils.ValueOrEmpty(req.PengelolaSe),
		utils.ValueOrEmpty(req.FiturSe),
		utils.ValueOrEmpty(req.KategoriSe),
	)
	return err
}

func (r *SeCsirtRepository) GetAll() ([]dto.SeCsirtResponse, error) {
	rows, err := r.db.Query(`
		SELECT
			se.id, se.nama_se, se.ip_se, se.as_number_se,
			se.pengelola_se, se.fitur_se, se.kategori_se,
			se.created_at, se.updated_at,
			c.id, c.nama_csirt, c.web_csirt, c.telepon_csirt
		FROM se_csirt se
		LEFT JOIN csirt c ON se.id_csirt = c.id
		ORDER BY se.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.SeCsirtResponse

	for rows.Next() {
		var se dto.SeCsirtResponse
		var csirt dto.CsirtMiniResponse
		var csirtID sql.NullString

		if err := rows.Scan(
			&se.ID,
			&se.NamaSe,
			&se.IpSe,
			&se.AsNumberSe,
			&se.Pengelola,
			&se.FiturSe,
			&se.KategoriSe,
			&se.CreatedAt,
			&se.UpdatedAt,
			&csirtID,
			&csirt.NamaCsirt,
			&csirt.WebCsirt,
			&csirt.TeleponCsirt,
		); err != nil {
			return nil, err
		}

		if csirtID.Valid {
			csirt.ID = csirtID.String
			se.Csirt = &csirt
		}

		result = append(result, se)
	}

	return result, nil
}

func (r *SeCsirtRepository) GetByID(id string) (*dto.SeCsirtResponse, error) {
	var se dto.SeCsirtResponse
	var csirt dto.CsirtMiniResponse
	var csirtID sql.NullString

	err := r.db.QueryRow(`
		SELECT
			se.id, se.nama_se, se.ip_se, se.as_number_se,
			se.pengelola_se, se.fitur_se, se.kategori_se,
			se.created_at, se.updated_at,
			c.id, c.nama_csirt, c.web_csirt, c.telepon_csirt
		FROM se_csirt se
		LEFT JOIN csirt c ON se.id_csirt = c.id
		WHERE se.id = ?
	`, id).Scan(
		&se.ID,
		&se.NamaSe,
		&se.IpSe,
		&se.AsNumberSe,
		&se.Pengelola,
		&se.FiturSe,
		&se.KategoriSe,
		&se.CreatedAt,
		&se.UpdatedAt,
		&csirtID,
		&csirt.NamaCsirt,
		&csirt.WebCsirt,
		&csirt.TeleponCsirt,
	)

	if err != nil {
		return nil, err
	}

	if csirtID.Valid {
		csirt.ID = csirtID.String
		se.Csirt = &csirt
	}

	return &se, nil
}

func (r *SeCsirtRepository) Update(id string, s dto.SeCsirtResponse) error {
	_, err := r.db.Exec(`
		UPDATE se_csirt SET
			nama_se = ?, ip_se = ?, as_number_se = ?,
			pengelola_se = ?, fitur_se = ?, kategori_se = ?
		WHERE id = ?
	`,
		s.NamaSe,
		s.IpSe,
		s.AsNumberSe,
		s.Pengelola,
		s.FiturSe,
		s.KategoriSe,
		id,
	)
	return err
}

func (r *SeCsirtRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM se_csirt WHERE id = ?`, id)
	return err
}
