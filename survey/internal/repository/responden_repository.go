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
// HELPER (ANTI NULL PANIC)
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
func (r *RespondenRepository) Create(m models.Responden) error {
	_, err := r.db.Exec(`
		INSERT INTO responden
		(user_id, no_telepon, sektor, sektor_lainnya, sertifikat_training)
		VALUES (?, ?, ?, ?, ?)
	`,
		m.UserID,
		m.NoTelepon,
		m.Sektor,
		m.SektorLainnya,
		m.SertifikatTraining,
	)

	return err
}

// ========================
// BASE QUERY (REUSABLE)
// ========================
const baseDetailQuery = `
SELECT 
	r.id,
	r.user_id,
	u.display_name,
	u.email,
	j.nama_jabatan,
	p.nama_perusahaan,
	p.id AS perusahaan_id,
	r.no_telepon,
	r.sektor,
	r.sektor_lainnya,
	r.sertifikat_training,
	r.created_at,
	r.updated_at
FROM responden r
LEFT JOIN users u ON r.user_id = u.id
LEFT JOIN jabatan j ON u.id_jabatan = j.id
LEFT JOIN perusahaan p ON u.id_perusahaan = p.id
`

// ========================
// GET ALL
// ========================
func (r *RespondenRepository) GetAllDetail() ([]models.RespondenDetail, error) {

	rows, err := r.db.Query(baseDetailQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.RespondenDetail

	for rows.Next() {

		var (
			m models.RespondenDetail

			displayName sql.NullString
			email       sql.NullString
			jabatan     sql.NullString
			perusahaan  sql.NullString
			perusahaanID sql.NullString
			sektorLainnya sql.NullString
		)

		err := rows.Scan(
			&m.ID,
			&m.UserID,
			&displayName,
			&email,
			&jabatan,
			&perusahaan,
			&perusahaanID,
			&m.NoTelepon,
			&m.Sektor,
			&sektorLainnya,
			&m.SertifikatTraining,
			&m.CreatedAt,
			&m.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		m.NamaLengkap = nullToString(displayName)
		m.Email = nullToString(email)
		m.Jabatan = nullToString(jabatan)
		m.NamaPerusahaan = nullToString(perusahaan)
		m.PerusahaanID = nullToString(perusahaanID)
		m.SektorLainnya = nullToString(sektorLainnya)

		result = append(result, m)
	}

	return result, nil
}

// ========================
// GET BY ID
// ========================
func (r *RespondenRepository) GetDetailByID(id int) (*models.RespondenDetail, error) {

	row := r.db.QueryRow(baseDetailQuery + " WHERE r.id = ?", id)

	var (
		m models.RespondenDetail

		displayName sql.NullString
		email       sql.NullString
		jabatan     sql.NullString
		perusahaan  sql.NullString
		perusahaanID sql.NullString
		sektorLainnya sql.NullString
	)

	err := row.Scan(
		&m.ID,
		&m.UserID,
		&displayName,
		&email,
		&jabatan,
		&perusahaan,
		&perusahaanID,
		&m.NoTelepon,
		&m.Sektor,
		&sektorLainnya,
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

	m.NamaLengkap = nullToString(displayName)
	m.Email = nullToString(email)
	m.Jabatan = nullToString(jabatan)
	m.NamaPerusahaan = nullToString(perusahaan)
	m.PerusahaanID = nullToString(perusahaanID)
	m.SektorLainnya = nullToString(sektorLainnya)

	return &m, nil
}

// ========================
// GET BY USER ID
// ========================
func (r *RespondenRepository) GetDetailByUserID(userID string) (*models.RespondenDetail, error) {

	row := r.db.QueryRow(baseDetailQuery+" WHERE r.user_id = ?", userID)

	var (
		m models.RespondenDetail

		displayName sql.NullString
		email       sql.NullString
		jabatan     sql.NullString
		perusahaan  sql.NullString
		perusahaanID sql.NullString
		sektorLainnya sql.NullString
	)

	err := row.Scan(
		&m.ID,
		&m.UserID,
		&displayName,
		&email,
		&jabatan,
		&perusahaan,
		&perusahaanID,
		&m.NoTelepon,
		&m.Sektor,
		&sektorLainnya,
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

	m.NamaLengkap = nullToString(displayName)
	m.Email = nullToString(email)
	m.Jabatan = nullToString(jabatan)
	m.NamaPerusahaan = nullToString(perusahaan)
	m.PerusahaanID = nullToString(perusahaanID)
	m.SektorLainnya = nullToString(sektorLainnya)

	return &m, nil
}

// ========================
// GET BASIC BY ID
// ========================
func (r *RespondenRepository) GetByID(id int) (*models.Responden, error) {

	row := r.db.QueryRow(`
		SELECT id, user_id, no_telepon, sektor, sektor_lainnya, sertifikat_training, created_at, updated_at
		FROM responden WHERE id = ?
	`, id)

	var m models.Responden

	err := row.Scan(
		&m.ID,
		&m.UserID,
		&m.NoTelepon,
		&m.Sektor,
		&m.SektorLainnya,
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
			no_telepon = ?,
			sektor = ?,
			sektor_lainnya = ?,
			sertifikat_training = ?,
			updated_at = NOW()
		WHERE id = ?
	`,
		m.NoTelepon,
		m.Sektor,
		m.SektorLainnya,
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

// ========================
// EXISTS
// ========================
func (r *RespondenRepository) Exists(id int) (bool, error) {

	var exists bool

	err := r.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM responden WHERE id = ?)
	`, id).Scan(&exists)

	return exists, err
}