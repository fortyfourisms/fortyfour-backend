package repository

import (
	"database/sql"
	"errors"

	"survey/internal/models"
)

type RespondenRepository struct {
	db *sql.DB
}

func NewRespondenRepository(db *sql.DB) *RespondenRepository {
	return &RespondenRepository{db: db}
}

// ========================
// HELPER
// ========================
func nullToString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

// ========================
// CREATE
// ========================
func (r *RespondenRepository) Create(m models.Responden) (int64, error) {

	res, err := r.db.Exec(`
		INSERT INTO responden
		(id_perusahaan, nama_lengkap, jabatan, email, no_telepon, sertifikat_training)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		m.IdPerusahaan,
		m.NamaLengkap,
		m.Jabatan,
		m.Email,
		m.NoTelepon,
		m.SertifikatTraining,
	)

	if err != nil {
		return 0, err
	}

	insertID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return insertID, nil
}

// ========================
// BASE QUERY (FIXED)
// ========================
const baseDetailQuery = `
SELECT 
	r.id,
	r.id_perusahaan,
	r.nama_lengkap,
	r.jabatan,
	r.email,
	r.no_telepon,
	r.sertifikat_training,
	p.nama_perusahaan,
	ss.nama_sub_sektor,
	s.nama_sektor,
	r.created_at,
	r.updated_at
FROM responden r
LEFT JOIN perusahaan p ON r.id_perusahaan = p.id
LEFT JOIN sub_sektor ss ON p.id_sub_sektor = ss.id
LEFT JOIN sektor s ON ss.id_sektor = s.id
`

// ========================
// GET ALL
// ========================
func (r *RespondenRepository) GetAllDetail() ([]models.RespondenDetail, error) {

	rows, err := r.db.Query(baseDetailQuery + " ORDER BY r.id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.RespondenDetail

	for rows.Next() {

		var (
			m models.RespondenDetail

			namaPerusahaan sql.NullString
			subSektor      sql.NullString
			sektor         sql.NullString
		)

		err := rows.Scan(
			&m.ID,
			&m.IdPerusahaan,
			&m.NamaLengkap,
			&m.Jabatan,
			&m.Email,
			&m.NoTelepon,
			&m.SertifikatTraining,
			&namaPerusahaan,
			&subSektor,
			&sektor,
			&m.CreatedAt,
			&m.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		m.NamaPerusahaan = nullToString(namaPerusahaan)
		m.NamaSubSektor = nullToString(subSektor)
		m.NamaSektor = nullToString(sektor)

		result = append(result, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// ========================
// GET BY ID
// ========================
func (r *RespondenRepository) GetDetailByID(id int) (*models.RespondenDetail, error) {

	row := r.db.QueryRow(baseDetailQuery+" WHERE r.id = ?", id)

	var (
		m models.RespondenDetail

		namaPerusahaan sql.NullString
		subSektor      sql.NullString
		sektor         sql.NullString
	)

	err := row.Scan(
		&m.ID,
		&m.IdPerusahaan,
		&m.NamaLengkap,
		&m.Jabatan,
		&m.Email,
		&m.NoTelepon,
		&m.SertifikatTraining,
		&namaPerusahaan,
		&subSektor,
		&sektor,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}

	m.NamaPerusahaan = nullToString(namaPerusahaan)
	m.NamaSubSektor = nullToString(subSektor)
	m.NamaSektor = nullToString(sektor)

	return &m, nil
}

// ========================
// GET BASIC BY ID
// ========================
func (r *RespondenRepository) GetByID(id int) (*models.Responden, error) {

	row := r.db.QueryRow(`
		SELECT id, id_perusahaan, nama_lengkap, jabatan, email, no_telepon, sertifikat_training, created_at, updated_at
		FROM responden
		WHERE id = ?
	`, id)

	var m models.Responden

	err := row.Scan(
		&m.ID,
		&m.IdPerusahaan,
		&m.NamaLengkap,
		&m.Jabatan,
		&m.Email,
		&m.NoTelepon,
		&m.SertifikatTraining,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}

	return &m, nil
}

// ========================
// UPDATE
// ========================
func (r *RespondenRepository) Update(id int, m models.Responden) error {

	res, err := r.db.Exec(`
		UPDATE responden SET
			id_perusahaan = ?,
			nama_lengkap = ?,
			jabatan = ?,
			email = ?,
			no_telepon = ?,
			sertifikat_training = ?,
			updated_at = NOW()
		WHERE id = ?
	`,
		m.IdPerusahaan,
		m.NamaLengkap,
		m.Jabatan,
		m.Email,
		m.NoTelepon,
		m.SertifikatTraining,
		id,
	)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("data tidak ditemukan")
	}

	return nil
}

// EXISTS
func (r *RespondenRepository) Exists(id int) (bool, error) {

	var exists bool

	err := r.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM responden WHERE id = ?)
	`, id).Scan(&exists)

	return exists, err
}