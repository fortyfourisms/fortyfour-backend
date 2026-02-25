package repository

import (
	"database/sql"
	"ikas/internal/dto"
	"strings"

	"fortyfour-backend/pkg/logger"
)

type RuangLingkupRepository struct {
	db *sql.DB
}

func NewRuangLingkupRepository(db *sql.DB) *RuangLingkupRepository {
	return &RuangLingkupRepository{db: db}
}

func (r *RuangLingkupRepository) Create(req dto.CreateRuangLingkupRequest, id string) error {
	// Prepared statement sudah aman dari SQL injection
	query := `INSERT INTO ruang_lingkup (id, nama_ruang_lingkup) VALUES (?, ?)`

	_, err := r.db.Exec(query, id, req.NamaRuangLingkup)
	if err != nil {
		logger.Error(err, "operation failed")
		return err
	}

	return nil
}

func (r *RuangLingkupRepository) GetAll() ([]dto.RuangLingkupResponse, error) {
	query := `SELECT id, nama_ruang_lingkup, created_at, updated_at FROM ruang_lingkup ORDER BY nama_ruang_lingkup ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}
	defer rows.Close()

	var result []dto.RuangLingkupResponse

	for rows.Next() {
		var item dto.RuangLingkupResponse
		if err := rows.Scan(&item.ID, &item.NamaRuangLingkup, &item.CreatedAt, &item.UpdatedAt); err != nil {
			logger.Error(err, "operation failed")
			continue
		}
		result = append(result, item)
	}

	return result, nil
}

func (r *RuangLingkupRepository) GetByID(id string) (*dto.RuangLingkupResponse, error) {
	// Prepared statement dengan placeholder ? untuk mencegah SQL injection
	query := `SELECT id, nama_ruang_lingkup, created_at, updated_at FROM ruang_lingkup WHERE id = ?`

	var item dto.RuangLingkupResponse
	err := r.db.QueryRow(query, id).Scan(&item.ID, &item.NamaRuangLingkup, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}

	return &item, nil
}

func (r *RuangLingkupRepository) Update(id string, req dto.UpdateRuangLingkupRequest) error {
	query := "UPDATE ruang_lingkup SET "
	args := []interface{}{}
	updates := []string{}

	if req.NamaRuangLingkup != nil {
		// Menggunakan placeholder ? dan parameter terpisah
		updates = append(updates, "nama_ruang_lingkup=?")
		args = append(args, *req.NamaRuangLingkup)
	}

	if len(updates) == 0 {
		return nil
	}

	query += strings.Join(updates, ", ")
	query += " WHERE id=?"
	args = append(args, id)

	// Prepared statement dengan args array
	_, err := r.db.Exec(query, args...)
	if err != nil {
		logger.Error(err, "operation failed")
		return err
	}

	return nil
}

func (r *RuangLingkupRepository) Delete(id string) error {
	// Prepared statement untuk DELETE
	_, err := r.db.Exec(`DELETE FROM ruang_lingkup WHERE id=?`, id)
	if err != nil {
		logger.Error(err, "operation failed")
		return err
	}

	return nil
}

// Cek duplikasi nama (case-insensitive, trim whitespace)
func (r *RuangLingkupRepository) CheckDuplicateName(nama string, excludeID string) (bool, error) {
	var count int

	// Prepared statement dengan placeholder ?
	query := `SELECT COUNT(*) FROM ruang_lingkup WHERE LOWER(TRIM(nama_ruang_lingkup)) = LOWER(TRIM(?))`
	args := []interface{}{nama}

	if excludeID != "" {
		query += ` AND id != ?`
		args = append(args, excludeID)
	}

	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		logger.Error(err, "operation failed")
		return false, err
	}

	return count > 0, nil
}
