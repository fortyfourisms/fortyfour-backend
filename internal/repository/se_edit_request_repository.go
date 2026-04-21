package repository

import (
	"database/sql"
	"fortyfour-backend/internal/models"
)

type SEEditRequestRepository struct {
	db *sql.DB
}

func NewSEEditRequestRepository(db *sql.DB) *SEEditRequestRepository {
	return &SEEditRequestRepository{db: db}
}

var _ SEEditRequestRepositoryInterface = (*SEEditRequestRepository)(nil)

func (r *SEEditRequestRepository) Create(req *models.SEEditRequest) error {
	_, err := r.db.Exec(
		`INSERT INTO se_edit_request (id, id_se, id_user, status, catatan_user, data_perubahan, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`,
		req.ID, req.IDSE, req.IDUser, req.Status, req.CatatanUser, req.DataPerubahan,
	)
	return err
}

func (r *SEEditRequestRepository) FindByID(id string) (*models.SEEditRequest, error) {
	row := r.db.QueryRow(
		`SELECT id, id_se, id_user, status, catatan_user, catatan, data_perubahan, created_at, updated_at
		 FROM se_edit_request WHERE id = ?`, id,
	)
	return scanSEEditRequest(row)
}

func (r *SEEditRequestRepository) FindPendingBySE(idSE string) ([]models.SEEditRequest, error) {
	rows, err := r.db.Query(
		`SELECT id, id_se, id_user, status, catatan_user, catatan, data_perubahan, created_at, updated_at
		 FROM se_edit_request WHERE id_se = ? AND status = 'pending' ORDER BY created_at DESC`, idSE,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSEEditRequestRows(rows)
}

func (r *SEEditRequestRepository) FindAllPending() ([]models.SEEditRequest, error) {
	rows, err := r.db.Query(
		`SELECT id, id_se, id_user, status, catatan_user, catatan, data_perubahan, created_at, updated_at
		 FROM se_edit_request WHERE status = 'pending' ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSEEditRequestRows(rows)
}

func (r *SEEditRequestRepository) FindByUser(idUser string) ([]models.SEEditRequest, error) {
	rows, err := r.db.Query(
		`SELECT id, id_se, id_user, status, catatan_user, catatan, data_perubahan, created_at, updated_at
		 FROM se_edit_request WHERE id_user = ? ORDER BY created_at DESC`, idUser,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSEEditRequestRows(rows)
}

func (r *SEEditRequestRepository) UpdateStatus(id string, status models.SEEditRequestStatus, catatan *string) error {
	_, err := r.db.Exec(
		`UPDATE se_edit_request SET status = ?, catatan = ?, updated_at = NOW() WHERE id = ?`,
		status, catatan, id,
	)
	return err
}

// ── helpers ───────────────────────────────────────────────────────────────────

func scanSEEditRequest(row *sql.Row) (*models.SEEditRequest, error) {
	var req models.SEEditRequest
	var catatanUser, catatan sql.NullString
	err := row.Scan(
		&req.ID, &req.IDSE, &req.IDUser, &req.Status, &catatanUser, &catatan,
		&req.DataPerubahan, &req.CreatedAt, &req.UpdatedAt,
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

func scanSEEditRequestRows(rows *sql.Rows) ([]models.SEEditRequest, error) {
	var result []models.SEEditRequest
	for rows.Next() {
		var req models.SEEditRequest
		var catatanUser, catatan sql.NullString
		if err := rows.Scan(
			&req.ID, &req.IDSE, &req.IDUser, &req.Status, &catatanUser, &catatan,
			&req.DataPerubahan, &req.CreatedAt, &req.UpdatedAt,
		); err != nil {
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
	return result, nil
}
