package repository

import (
	"database/sql"
	"errors"

	"survey/internal/models"
	"github.com/google/uuid"
)

type RespondenRepository struct {
	db *sql.DB
}

func NewRespondenRepository(db *sql.DB) *RespondenRepository {
	return &RespondenRepository{db: db}
}

// HELPER
func nullToString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

// CREATE
func (r *RespondenRepository) Create(m models.Responden) (string, error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}

	_, err := r.db.Exec(`
		INSERT INTO responden
		(id, user_id, id_perusahaan, nama_lengkap, jabatan, email, no_telepon, sertifikat_training)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		m.ID,
		m.UserID,
		m.IdPerusahaan,
		m.NamaLengkap,
		m.Jabatan,
		m.Email,
		m.NoTelepon,
		m.SertifikatTraining,
	)

	if err != nil {
		return "", err
	}

	return m.ID, nil
}

// UPSERT BY USER ID
func (r *RespondenRepository) UpsertByUserID(userID string, m models.Responden) error {
	newID := uuid.New().String()
	_, err := r.db.Exec(`
		INSERT INTO responden
		(id, user_id, id_perusahaan, nama_lengkap, jabatan, email, no_telepon, sertifikat_training)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			id_perusahaan = VALUES(id_perusahaan),
			nama_lengkap = VALUES(nama_lengkap),
			jabatan = VALUES(jabatan),
			email = VALUES(email),
			no_telepon = VALUES(no_telepon),
			sertifikat_training = VALUES(sertifikat_training),
			updated_at = NOW()
	`,
		newID,
		userID,
		m.IdPerusahaan,
		m.NamaLengkap,
		m.Jabatan,
		m.Email,
		m.NoTelepon,
		m.SertifikatTraining,
	)

	return err
}

func (r *RespondenRepository) CanEditByUserID(userID string) (bool, string, error) {
	row := r.db.QueryRow(`
		SELECT
			COALESCE(sp.selesai, false),
			COALESCE(sp.status, 'draft')
		FROM responden r
		LEFT JOIN survey_progress sp ON sp.responden_id = r.id
		WHERE r.user_id = ?
	`, userID)

	var selesai bool
	var status string
	if err := row.Scan(&selesai, &status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, "draft", nil
		}
		return false, "", err
	}

	if status == "" {
		status = "draft"
	}

	if !selesai || status == "draft" || status == "edit_approved" {
		return true, status, nil
	}

	return false, status, nil
}

// BASE QUERY
const baseDetailQuery = `
SELECT 
	r.id,
	r.user_id,
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

// GET ALL DETAIL (ADMIN)
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

			sertifikatTraining sql.NullString
			namaPerusahaan     sql.NullString
			subSektor          sql.NullString
			sektor             sql.NullString
		)

		err := rows.Scan(
			&m.ID,
			&m.UserID,
			&m.IdPerusahaan,
			&m.NamaLengkap,
			&m.Jabatan,
			&m.Email,
			&m.NoTelepon,
			&sertifikatTraining,
			&namaPerusahaan,
			&subSektor,
			&sektor,
			&m.CreatedAt,
			&m.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		m.SertifikatTraining = nullToString(sertifikatTraining)
		m.NamaPerusahaan = nullToString(namaPerusahaan)
		m.NamaSubSektor = nullToString(subSektor)
		m.NamaSektor = nullToString(sektor)

		result = append(result, m)
	}

	return result, nil
}

// GET BY ID (ADMIN)
func (r *RespondenRepository) GetDetailByID(id string) (*models.RespondenDetail, error) {

	row := r.db.QueryRow(baseDetailQuery+" WHERE r.id = ?", id)

	var (
		m models.RespondenDetail

		sertifikatTraining sql.NullString
		namaPerusahaan     sql.NullString
		subSektor          sql.NullString
		sektor             sql.NullString
	)

	err := row.Scan(
		&m.ID,
		&m.UserID,
		&m.IdPerusahaan,
		&m.NamaLengkap,
		&m.Jabatan,
		&m.Email,
		&m.NoTelepon,
		&sertifikatTraining,
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

	m.SertifikatTraining = nullToString(sertifikatTraining)
	m.NamaPerusahaan = nullToString(namaPerusahaan)
	m.NamaSubSektor = nullToString(subSektor)
	m.NamaSektor = nullToString(sektor)

	return &m, nil
}

// GET BY USER ID (USER)
func (r *RespondenRepository) GetByUserID(userID string) (*models.RespondenDetail, error) {

	row := r.db.QueryRow(baseDetailQuery+" WHERE r.user_id = ?", userID)

	var (
		m models.RespondenDetail

		sertifikatTraining sql.NullString
		namaPerusahaan     sql.NullString
		subSektor          sql.NullString
		sektor             sql.NullString
	)

	err := row.Scan(
		&m.ID,
		&m.UserID,
		&m.IdPerusahaan,
		&m.NamaLengkap,
		&m.Jabatan,
		&m.Email,
		&m.NoTelepon,
		&sertifikatTraining,
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

	m.SertifikatTraining = nullToString(sertifikatTraining)
	m.NamaPerusahaan = nullToString(namaPerusahaan)
	m.NamaSubSektor = nullToString(subSektor)
	m.NamaSektor = nullToString(sektor)

	return &m, nil
}

// BASIC GET
func (r *RespondenRepository) GetByID(id string) (*models.Responden, error) {

	row := r.db.QueryRow(`
		SELECT id, user_id, id_perusahaan, nama_lengkap, jabatan, email, no_telepon, sertifikat_training, created_at, updated_at
		FROM responden
		WHERE id = ?
	`, id)

	var (
		m                  models.Responden
		sertifikatTraining sql.NullString
	)

	err := row.Scan(
		&m.ID,
		&m.UserID,
		&m.IdPerusahaan,
		&m.NamaLengkap,
		&m.Jabatan,
		&m.Email,
		&m.NoTelepon,
		&sertifikatTraining,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}

	m.SertifikatTraining = nullToString(sertifikatTraining)

	return &m, nil
}
