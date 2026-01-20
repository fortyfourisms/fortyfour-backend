package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
)

type seRepository struct {
	db *sql.DB
}

func NewSERepository(db *sql.DB) SERepositoryInterface {
	return &seRepository{db: db}
}

/* ================= CREATE ================= */

func (r *seRepository) Create(
	req dto.CreateSERequest,
	id string,
	totalBobot int,
	kategori string,
) error {
	_, err := r.db.Exec(`
		INSERT INTO se (
			id,
			id_perusahaan,
			id_sub_sektor,
			q1, q2, q3, q4, q5,
			q6, q7, q8, q9, q10,
			total_bobot,
			kategori_se
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	`,
		id,
		*req.IDPerusahaan,
		*req.IDSubSektor,
		req.Q1, req.Q2, req.Q3, req.Q4, req.Q5,
		req.Q6, req.Q7, req.Q8, req.Q9, req.Q10,
		totalBobot,
		kategori,
	)

	return err
}

/* ================= GET ALL ================= */

func (r *seRepository) GetAll() ([]dto.SEResponse, error) {
	rows, err := r.db.Query(`
		SELECT
			se.id,
			se.q1, se.q2, se.q3, se.q4, se.q5,
			se.q6, se.q7, se.q8, se.q9, se.q10,
			se.total_bobot,
			se.kategori_se,
			se.created_at,
			se.updated_at,

			p.id,
			p.nama_perusahaan,

			ss.id,
			ss.nama_sub_sektor,
			s.id,
			s.nama_sektor
		FROM se
		JOIN perusahaan p ON se.id_perusahaan = p.id
		JOIN sub_sektor ss ON se.id_sub_sektor = ss.id
		JOIN sektor s ON ss.id_sektor = s.id
		ORDER BY se.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.SEResponse

	for rows.Next() {
		var se dto.SEResponse
		se.Perusahaan = &dto.PerusahaanMiniResponse{}
		se.SubSektor = &dto.SubSektorMiniResponse{}

		err := rows.Scan(
			&se.ID,
			&se.Q1, &se.Q2, &se.Q3, &se.Q4, &se.Q5,
			&se.Q6, &se.Q7, &se.Q8, &se.Q9, &se.Q10,
			&se.TotalBobot,
			&se.KategoriSE,
			&se.CreatedAt,
			&se.UpdatedAt,

			&se.Perusahaan.ID,
			&se.Perusahaan.NamaPerusahaan,

			&se.SubSektor.ID,
			&se.SubSektor.NamaSubSektor,
			&se.SubSektor.IDSektor,
			&se.SubSektor.NamaSektor,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, se)
	}

	return result, nil
}

/* ================= GET BY ID ================= */

func (r *seRepository) GetByID(id string) (*dto.SEResponse, error) {
	row := r.db.QueryRow(`
		SELECT
			se.id,
			se.q1, se.q2, se.q3, se.q4, se.q5,
			se.q6, se.q7, se.q8, se.q9, se.q10,
			se.total_bobot,
			se.kategori_se,
			se.created_at,
			se.updated_at,

			p.id,
			p.nama_perusahaan,

			ss.id,
			ss.nama_sub_sektor,
			s.id,
			s.nama_sektor
		FROM se
		JOIN perusahaan p ON se.id_perusahaan = p.id
		JOIN sub_sektor ss ON se.id_sub_sektor = ss.id
		JOIN sektor s ON ss.id_sektor = s.id
		WHERE se.id = ?
	`, id)

	var se dto.SEResponse
	se.Perusahaan = &dto.PerusahaanMiniResponse{}
	se.SubSektor = &dto.SubSektorMiniResponse{}

	err := row.Scan(
		&se.ID,
		&se.Q1, &se.Q2, &se.Q3, &se.Q4, &se.Q5,
		&se.Q6, &se.Q7, &se.Q8, &se.Q9, &se.Q10,
		&se.TotalBobot,
		&se.KategoriSE,
		&se.CreatedAt,
		&se.UpdatedAt,

		&se.Perusahaan.ID,
		&se.Perusahaan.NamaPerusahaan,

		&se.SubSektor.ID,
		&se.SubSektor.NamaSubSektor,
		&se.SubSektor.IDSektor,
		&se.SubSektor.NamaSektor,
	)
	if err != nil {
		return nil, err
	}

	return &se, nil
}

/* ================= UPDATE ================= */

func (r *seRepository) Update(
	id string,
	req dto.UpdateSERequest,
	totalBobot int,
	kategori string,
) error {
	_, err := r.db.Exec(`
		UPDATE se SET
			q1 = ?, q2 = ?, q3 = ?, q4 = ?, q5 = ?,
			q6 = ?, q7 = ?, q8 = ?, q9 = ?, q10 = ?,
			total_bobot = ?,
			kategori_se = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`,
		req.Q1, req.Q2, req.Q3, req.Q4, req.Q5,
		req.Q6, req.Q7, req.Q8, req.Q9, req.Q10,
		totalBobot,
		kategori,
		id,
	)

	return err
}

/* ================= DELETE ================= */

func (r *seRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM se WHERE id = ?`, id)
	return err
}
