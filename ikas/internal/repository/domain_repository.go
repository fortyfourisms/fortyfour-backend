package repository

import (
	"database/sql"
	"ikas/internal/dto"
	"strings"

	"fortyfour-backend/pkg/logger"
)

type DomainRepository struct {
	db *sql.DB
}

func NewDomainRepository(db *sql.DB) *DomainRepository {
	return &DomainRepository{db: db}
}

func (r *DomainRepository) Create(req dto.CreateDomainRequest, id string) error {
	query := `INSERT INTO domain (id, nama_domain) VALUES (?, ?)`
	_, err := r.db.Exec(query, id, req.NamaDomain)
	if err != nil {
		logger.Error(err, "operation failed")
	}
	return err
}

func (r *DomainRepository) GetAll() ([]dto.DomainResponse, error) {
	query := `SELECT id, nama_domain, created_at, updated_at FROM domain ORDER BY nama_domain ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}
	defer rows.Close()

	var result []dto.DomainResponse
	for rows.Next() {
		var item dto.DomainResponse
		if err := rows.Scan(&item.ID, &item.NamaDomain, &item.CreatedAt, &item.UpdatedAt); err != nil {
			logger.Error(err, "operation failed")
			continue
		}
		result = append(result, item)
	}
	return result, nil
}

func (r *DomainRepository) GetByID(id string) (*dto.DomainResponse, error) {
	query := `SELECT id, nama_domain, created_at, updated_at FROM domain WHERE id=?`

	var item dto.DomainResponse
	err := r.db.QueryRow(query, id).
		Scan(&item.ID, &item.NamaDomain, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}
	return &item, nil
}

func (r *DomainRepository) Update(id string, req dto.UpdateDomainRequest) error {
	query := "UPDATE domain SET "
	args := []interface{}{}
	updates := []string{}

	if req.NamaDomain != nil {
		updates = append(updates, "nama_domain=?")
		args = append(args, *req.NamaDomain)
	}

	if len(updates) == 0 {
		return nil
	}

	query += strings.Join(updates, ", ")
	query += " WHERE id=?"
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		logger.Error(err, "operation failed")
	}
	return err
}

func (r *DomainRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM domain WHERE id=?`, id)
	if err != nil {
		logger.Error(err, "operation failed")
	}
	return err
}

func (r *DomainRepository) CheckDuplicateName(nama string, excludeID string) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM domain WHERE LOWER(TRIM(nama_domain)) = LOWER(TRIM(?))`
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
