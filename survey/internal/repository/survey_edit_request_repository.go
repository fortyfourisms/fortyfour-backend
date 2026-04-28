package repository

import (
	"database/sql"
	"survey/internal/models"
)

type EditRequestRepository struct {
	db *sql.DB
}

func NewEditRequestRepository(db *sql.DB) *EditRequestRepository {
	return &EditRequestRepository{db: db}
}

// CREATE
func (r *EditRequestRepository) Create(req *models.EditRequest) error {
	_, err := r.db.Exec(`
		INSERT INTO se_edit_request 
		(id, responden_id, risiko_id, id_user, status, catatan_user, data_perubahan, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`,
		req.ID,
		req.RespondenID,
		req.RisikoID,
		req.UserID,
		req.Status,
		req.CatatanUser,
		req.DataPerubahan,
	)

	return err
}

// FIND BY ID
func (r *EditRequestRepository) FindByID(id string) (*models.EditRequest, error) {
	row := r.db.QueryRow(`
		SELECT id, responden_id, risiko_id, id_user, status, catatan_user, catatan, data_perubahan, created_at, updated_at
		FROM se_edit_request
		WHERE id = ?
	`, id)

	return scanEditRequest(row)
}

// FIND ALL PENDING
func (r *EditRequestRepository) FindAllPending() ([]models.EditRequest, error) {
	rows, err := r.db.Query(`
		SELECT id, responden_id, risiko_id, id_user, status, catatan_user, catatan, data_perubahan, created_at, updated_at
		FROM se_edit_request
		WHERE status = 'pending'
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanEditRequestRows(rows)
}

// FIND BY USER
func (r *EditRequestRepository) FindByUser(userID string) ([]models.EditRequest, error) {
	rows, err := r.db.Query(`
		SELECT id, responden_id, risiko_id, id_user, status, catatan_user, catatan, data_perubahan, created_at, updated_at
		FROM se_edit_request
		WHERE id_user = ?
		ORDER BY created_at DESC
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanEditRequestRows(rows)
}

// FIND PENDING BY RESPONDEN + RISIKO
func (r *EditRequestRepository) FindPendingByRespondenRisiko(respondenID, risikoID int) ([]models.EditRequest, error) {
	rows, err := r.db.Query(`
		SELECT id, responden_id, risiko_id, id_user, status, catatan_user, catatan, data_perubahan, created_at, updated_at
		FROM se_edit_request
		WHERE responden_id = ? AND risiko_id = ? AND status = 'pending'
	`, respondenID, risikoID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanEditRequestRows(rows)
}

// UPDATE STATUS
func (r *EditRequestRepository) UpdateStatus(id string, status models.EditRequestStatus, catatan *string) error {
	_, err := r.db.Exec(`
		UPDATE se_edit_request 
		SET status = ?, catatan = ?, updated_at = NOW()
		WHERE id = ?
	`, status, catatan, id)

	return err
}

// HELPER

func scanEditRequest(row *sql.Row) (*models.EditRequest, error) {
	var req models.EditRequest
	var catatanUser, catatan sql.NullString

	err := row.Scan(
		&req.ID,
		&req.RespondenID,
		&req.RisikoID,
		&req.UserID,
		&req.Status,
		&catatanUser,
		&catatan,
		&req.DataPerubahan,
		&req.CreatedAt,
		&req.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if catatanUser.Valid {
		req.CatatanUser = &catatanUser.String
	}
	if catatan.Valid {
		req.Catatan = &catatan.String
	}

	return &req, nil
}

func scanEditRequestRows(rows *sql.Rows) ([]models.EditRequest, error) {
	var result []models.EditRequest

	for rows.Next() {
		var req models.EditRequest
		var catatanUser, catatan sql.NullString

		err := rows.Scan(
			&req.ID,
			&req.RespondenID,
			&req.RisikoID,
			&req.UserID,
			&req.Status,
			&catatanUser,
			&catatan,
			&req.DataPerubahan,
			&req.CreatedAt,
			&req.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		if catatanUser.Valid {
			req.CatatanUser = &catatanUser.String
		}
		if catatan.Valid {
			req.Catatan = &catatan.String
		}

		result = append(result, req)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
